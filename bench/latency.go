package bench

import (
	"fmt"
	"io"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/unedtamps/orbook/book"
)

const (
	histMinNs = 1
	histMaxNs = 50_000_000 // 50 ms ceiling
	histSig   = 3
)

// OpPercentiles summarizes latency percentiles (nanoseconds) for one operation category.
type OpPercentiles struct {
	Op    string
	Count int64
	Min   int64
	P50   int64
	P90   int64
	P99   int64
	P999  int64
	Max   int64
	Mean  float64
}

// VariantLatencyReport captures throughput and HDR latency distributions for one Variant.
type VariantLatencyReport struct {
	Variant       Variant
	OpsCount      int
	ThroughputOps float64 // pure untimed-per-op throughput (ops/sec)
	AvgNsPerOp    float64
	All           OpPercentiles
	Limit         OpPercentiles
	Cancel        OpPercentiles
	Market        OpPercentiles
}

func summarizeHist(op string, h *hdrhistogram.Histogram) OpPercentiles {
	return OpPercentiles{
		Op:    op,
		Count: h.TotalCount(),
		Min:   h.Min(),
		P50:   h.ValueAtQuantile(50.0),
		P90:   h.ValueAtQuantile(90.0),
		P99:   h.ValueAtQuantile(99.0),
		P999:  h.ValueAtQuantile(99.9),
		Max:   h.Max(),
		Mean:  h.Mean(),
	}
}

func recordNs(h *hdrhistogram.Histogram, ns int64) {
	if ns < histMinNs {
		ns = histMinNs
	} else if ns > histMaxNs {
		ns = histMaxNs
	}
	_ = h.RecordValue(ns)
}

// MeasureVariantLatency runs w on v twice:
//  1. A pure untimed-per-op pass to measure accurate steady-state throughput (ops/sec)
//     without clock-read overhead on every op.
//  2. A per-operation timed pass recording nanosecond latencies into HDR histograms
//     for All, Limit, Cancel, and Market operations.
func MeasureVariantLatency(v Variant, w *Workload) (VariantLatencyReport, error) {
	evBuf := make([]book.Event, 0, 64)
	prewarm := len(w.Warmup) + 8192

	// Pass 1: Pure throughput measurement (no per-op time.Now overhead).
	runtime.GC()
	b1, err := v.NewBook(w.Config, prewarm)
	if err != nil {
		return VariantLatencyReport{}, err
	}
	evBuf = w.ApplyWarmup(b1, evBuf)

	tStart := time.Now()
	for i := range w.Ops {
		evBuf = ReplayOp(b1, &w.Ops[i], evBuf)
	}
	elapsed := time.Since(tStart)
	throughput := float64(len(w.Ops)) / elapsed.Seconds()
	avgNs := float64(elapsed.Nanoseconds()) / float64(len(w.Ops))

	// Pass 2: Per-operation HDR histogram recording.
	runtime.GC()
	b2, err := v.NewBook(w.Config, prewarm)
	if err != nil {
		return VariantLatencyReport{}, err
	}
	evBuf = w.ApplyWarmup(b2, evBuf)

	hAll := hdrhistogram.New(histMinNs, histMaxNs, histSig)
	hLimit := hdrhistogram.New(histMinNs, histMaxNs, histSig)
	hCancel := hdrhistogram.New(histMinNs, histMaxNs, histSig)
	hMarket := hdrhistogram.New(histMinNs, histMaxNs, histSig)

	for i := range w.Ops {
		op := &w.Ops[i]
		t0 := time.Now()
		evBuf = ReplayOp(b2, op, evBuf)
		ns := time.Since(t0).Nanoseconds()

		recordNs(hAll, ns)
		switch op.Kind {
		case OpLimit:
			recordNs(hLimit, ns)
		case OpCancel:
			recordNs(hCancel, ns)
		case OpMarket:
			recordNs(hMarket, ns)
		}
	}

	// Keep evBuf live to prevent dead-code elimination.
	_ = len(evBuf)

	return VariantLatencyReport{
		Variant:       v,
		OpsCount:      len(w.Ops),
		ThroughputOps: throughput,
		AvgNsPerOp:    avgNs,
		All:           summarizeHist("ALL", hAll),
		Limit:         summarizeHist("LIMIT", hLimit),
		Cancel:        summarizeHist("CANCEL", hCancel),
		Market:        summarizeHist("MARKET", hMarket),
	}, nil
}

// PrintLatencyReports writes a formatted HDR latency summary table to out.
func PrintLatencyReports(out io.Writer, w *Workload, reports []VariantLatencyReport, detailed bool) {
	fmt.Fprintf(out, "Workload: seed=%d warmup=%d ops=%d (limit=%d cancel=%d market=%d trades=%d stp=%d)\n\n",
		w.Seed, w.Stats.WarmupOps, w.Stats.TotalOps,
		w.Stats.LimitOps, w.Stats.CancelOps, w.Stats.MarketOps,
		w.Stats.Trades, w.Stats.STPCancels)

	tw := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	if !detailed {
		fmt.Fprintln(tw, "VARIANT\tTHROUGHPUT (ops/s)\tAVG (ns)\tp50 (ns)\tp90 (ns)\tp99 (ns)\tp99.9 (ns)\tMAX (ns)")
		fmt.Fprintln(tw, "-------\t------------------\t--------\t--------\t--------\t--------\t----------\t--------")
		for _, r := range reports {
			fmt.Fprintf(tw, "%s\t%.0f\t%.1f\t%d\t%d\t%d\t%d\t%d\n",
				r.Variant.Name,
				r.ThroughputOps,
				r.AvgNsPerOp,
				r.All.P50,
				r.All.P90,
				r.All.P99,
				r.All.P999,
				r.All.Max,
			)
		}
		_ = tw.Flush()
		return
	}

	fmt.Fprintln(tw, "VARIANT\tOP\tCOUNT\tp50 (ns)\tp90 (ns)\tp99 (ns)\tp99.9 (ns)\tMAX (ns)\tMEAN (ns)")
	fmt.Fprintln(tw, "-------\t--\t-----\t--------\t--------\t--------\t----------\t--------\t---------")
	for _, r := range reports {
		for _, p := range []OpPercentiles{r.All, r.Limit, r.Cancel, r.Market} {
			fmt.Fprintf(tw, "%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%.1f\n",
				r.Variant.Name, p.Op, p.Count, p.P50, p.P90, p.P99, p.P999, p.Max, p.Mean)
		}
	}
	_ = tw.Flush()
}
