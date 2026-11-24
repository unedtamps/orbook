package bench_test

import (
	"bytes"
	"math/rand"
	"sync"
	"testing"

	"github.com/unedtamps/orbook/bench"
	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/engine"
)

var (
	sinkInt64 int64
	sinkBool  bool
	sinkInt   int

	sharedWorkloadOnce sync.Once
	sharedWorkload     *bench.Workload
)

func getSharedWorkload() *bench.Workload {
	sharedWorkloadOnce.Do(func() {
		sharedWorkload = bench.GenerateMixedWorkload(bench.DefaultConfig(), 20261122, 2000, 100_000)
	})
	return sharedWorkload
}

func BenchmarkAddLimit(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const batchSize = 4096

	rng := rand.New(rand.NewSource(101))
	reqs := make([]book.OrderRequest, batchSize)
	for i := 0; i < batchSize; i++ {
		side := book.Buy
		offset := int64(1 + rng.Intn(200))
		price := mid - offset
		if i%2 == 1 {
			side = book.Sell
			price = mid + offset
		}
		reqs[i] = book.OrderRequest{
			ID:       book.OrderID(10_000 + i),
			TraderID: book.TraderID(1 + i%16),
			Side:     side,
			Type:     book.Limit,
			Price:    price,
			Qty:      100,
		}
	}

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, batchSize+512)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 16)

			for k := 1; k <= 100; k++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k*2), Qty: 100,
				}, evBuf[:0])
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k*2), Qty: 100,
				}, evBuf[:0])
			}

			for i := 0; i < batchSize; i++ {
				evBuf = bk.PlaceOrder(reqs[i], evBuf[:0])
			}
			for i := 0; i < batchSize; i++ {
				evBuf = bk.CancelOrder(reqs[i].ID, evBuf[:0])
			}

			idx := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if idx == batchSize {
					b.StopTimer()
					for j := 0; j < batchSize; j++ {
						evBuf = bk.CancelOrder(reqs[j].ID, evBuf[:0])
					}
					idx = 0
					b.StartTimer()
				}
				evBuf = bk.PlaceOrder(reqs[idx], evBuf[:0])
				idx++
			}
			sinkInt = len(evBuf)
		})
	}
}

func BenchmarkCancel(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const batchSize = 4096

	rng := rand.New(rand.NewSource(202))
	reqs := make([]book.OrderRequest, batchSize)
	for i := 0; i < batchSize; i++ {
		side := book.Buy
		offset := int64(1 + rng.Intn(200))
		price := mid - offset
		if i%2 == 1 {
			side = book.Sell
			price = mid + offset
		}
		reqs[i] = book.OrderRequest{
			ID:       book.OrderID(10_000 + i),
			TraderID: book.TraderID(1 + i%16),
			Side:     side,
			Type:     book.Limit,
			Price:    price,
			Qty:      100,
		}
	}
	cancelOrder := make([]book.OrderID, batchSize)
	for i, p := range rng.Perm(batchSize) {
		cancelOrder[i] = reqs[p].ID
	}

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, batchSize+512)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 16)

			for k := 1; k <= 100; k++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k*2), Qty: 100,
				}, evBuf[:0])
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k*2), Qty: 100,
				}, evBuf[:0])
			}
			for i := 0; i < batchSize; i++ {
				evBuf = bk.PlaceOrder(reqs[i], evBuf[:0])
			}

			idx := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if idx == batchSize {
					b.StopTimer()
					for j := 0; j < batchSize; j++ {
						evBuf = bk.PlaceOrder(reqs[j], evBuf[:0])
					}
					idx = 0
					b.StartTimer()
				}
				evBuf = bk.CancelOrder(cancelOrder[idx], evBuf[:0])
				idx++
			}
			sinkInt = len(evBuf)
		})
	}
}

func BenchmarkMatchSingle(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const batchSize = 4096

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, batchSize+512)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 16)

			seedLevel := func() {
				for j := 0; j < batchSize; j++ {
					evBuf = bk.PlaceOrder(book.OrderRequest{
						ID:       book.OrderID(1 + j),
						TraderID: 1,
						Side:     book.Sell,
						Type:     book.Limit,
						Price:    mid,
						Qty:      100,
					}, evBuf[:0])
				}
			}
			seedLevel()

			idx := 0
			var takerID book.OrderID = 100_000
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if idx == batchSize {
					b.StopTimer()
					seedLevel()
					idx = 0
					b.StartTimer()
				}
				takerID++
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID:       takerID,
					TraderID: 2,
					Side:     book.Buy,
					Type:     book.Limit,
					Price:    mid,
					Qty:      100,
				}, evBuf[:0])
				idx++
			}
			sinkInt = len(evBuf)
		})
	}
}

func BenchmarkMatchSweep(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const sweepLevels = 10
	const numSweeps = 500
	const totalAsks = sweepLevels * numSweeps

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, totalAsks+512)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 64)

			seedAsks := func() {
				for j := 0; j < totalAsks; j++ {
					evBuf = bk.PlaceOrder(book.OrderRequest{
						ID:       book.OrderID(1 + j),
						TraderID: 1,
						Side:     book.Sell,
						Type:     book.Limit,
						Price:    mid + int64(j),
						Qty:      100,
					}, evBuf[:0])
				}
			}
			seedAsks()

			idx := 0
			var takerID book.OrderID = 1_000_000
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if idx == numSweeps {
					b.StopTimer()
					seedAsks()
					idx = 0
					b.StartTimer()
				}
				takerID++
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID:       takerID,
					TraderID: 2,
					Side:     book.Buy,
					Type:     book.Limit,
					Price:    mid + int64(totalAsks),
					Qty:      sweepLevels * 100,
				}, evBuf[:0])
				idx++
			}
			sinkInt = len(evBuf)
		})
	}
}

func BenchmarkMarket(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const sweepLevels = 10
	const numSweeps = 500
	const totalAsks = sweepLevels * numSweeps

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, totalAsks+512)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 64)

			seedAsks := func() {
				for j := 0; j < totalAsks; j++ {
					evBuf = bk.PlaceOrder(book.OrderRequest{
						ID:       book.OrderID(1 + j),
						TraderID: 1,
						Side:     book.Sell,
						Type:     book.Limit,
						Price:    mid + int64(j),
						Qty:      100,
					}, evBuf[:0])
				}
			}
			seedAsks()

			idx := 0
			var takerID book.OrderID = 1_000_000
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if idx == numSweeps {
					b.StopTimer()
					seedAsks()
					idx = 0
					b.StartTimer()
				}
				takerID++
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID:       takerID,
					TraderID: 2,
					Side:     book.Buy,
					Type:     book.Market,
					Qty:      sweepLevels * 100,
				}, evBuf[:0])
				idx++
			}
			sinkInt = len(evBuf)
		})
	}
}

func BenchmarkBestBidAsk(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, 1024)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 16)
			for k := 1; k <= 100; k++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k), Qty: 100,
				}, evBuf[:0])
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k), Qty: 100,
				}, evBuf[:0])
			}

			var bp, bq, ap, aq int64
			var bok, aok bool
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bp, bq, bok = bk.BestBid()
				ap, aq, aok = bk.BestAsk()
			}
			sinkInt64 = bp + bq + ap + aq
			sinkBool = bok && aok
		})
	}
}

func BenchmarkDepth(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(cfg, 1024)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 16)
			for k := 1; k <= 100; k++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k), Qty: 100,
				}, evBuf[:0])
				evBuf = bk.PlaceOrder(book.OrderRequest{
					ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k), Qty: 100,
				}, evBuf[:0])
			}

			bidsBuf := make([]book.DepthLevel, 0, 16)
			asksBuf := make([]book.DepthLevel, 0, 16)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bidsBuf, asksBuf = bk.Depth(10, bidsBuf, asksBuf)
			}
			sinkInt = len(bidsBuf) + len(asksBuf)
		})
	}
}

func BenchmarkMixedWorkload(b *testing.B) {
	w := getSharedWorkload()
	prewarm := len(w.Warmup) + 32768

	for _, v := range bench.AllVariants() {
		v := v
		b.Run(v.Name, func(b *testing.B) {
			bk, err := v.NewBook(w.Config, prewarm)
			if err != nil {
				b.Fatal(err)
			}
			evBuf := make([]book.Event, 0, 64)
			ops := w.Ops
			nOps := len(ops)

			evBuf = w.ApplyWarmup(bk, evBuf)
			for i := 0; i < nOps; i++ {
				evBuf = bench.ReplayOp(bk, &ops[i], evBuf)
			}
			bk.Reset()
			evBuf = w.ApplyWarmup(bk, evBuf)

			idx := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if idx == nOps {
					b.StopTimer()
					bk.Reset()
					evBuf = w.ApplyWarmup(bk, evBuf)
					idx = 0
					b.StartTimer()
				}
				evBuf = bench.ReplayOp(bk, &ops[idx], evBuf)
				idx++
			}
			sinkInt = len(evBuf)
		})
	}
}

func BenchmarkEngineChannel(b *testing.B) {
	w := getSharedWorkload()
	v := bench.Variant{Name: "ladder_linkedlist_pool", Levels: "ladder", Queue: "linkedlist", Pool: true}
	bk, err := v.NewBook(w.Config, len(w.Warmup)+32768)
	if err != nil {
		b.Fatal(err)
	}
	evBuf := make([]book.Event, 0, 64)
	evBuf = w.ApplyWarmup(bk, evBuf)
	for i := range w.Ops {
		evBuf = bench.ReplayOp(bk, &w.Ops[i], evBuf)
	}
	bk.Reset()

	var totalEvents int64
	startEngine := func() *engine.Engine {
		eng := engine.New(bk, engine.Options{
			QueueSize: 8192,
			EventHandler: func(events []book.Event) {
				totalEvents += int64(len(events))
			},
		})
		eng.Start()
		for i := range w.Warmup {
			eng.PlaceOrderAsync(w.Warmup[i].Req)
		}
		_ = eng.Stats()
		return eng
	}

	eng := startEngine()
	ops := w.Ops
	nOps := len(ops)
	idx := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := &ops[idx]
		if op.Kind == bench.OpCancel {
			eng.CancelOrderAsync(op.CancelID)
		} else {
			eng.PlaceOrderAsync(op.Req)
		}
		idx++
		if idx == nOps {
			b.StopTimer()
			eng.Stop()
			bk.Reset()
			eng = startEngine()
			idx = 0
			b.StartTimer()
		}
	}
	eng.Stop()
	b.StopTimer()
	sinkInt64 = totalEvents
}

func BenchmarkEngineBatch(b *testing.B) {
	w := getSharedWorkload()
	v := bench.Variant{Name: "ladder_linkedlist_pool", Levels: "ladder", Queue: "linkedlist", Pool: true}
	bk, err := v.NewBook(w.Config, len(w.Warmup)+32768)
	if err != nil {
		b.Fatal(err)
	}
	evBuf := make([]book.Event, 0, 64)
	evBuf = w.ApplyWarmup(bk, evBuf)
	for i := range w.Ops {
		evBuf = bench.ReplayOp(bk, &w.Ops[i], evBuf)
	}
	bk.Reset()

	const batchSize = 64
	nOps := len(w.Ops)
	allCmds := make([]engine.Command, nOps)
	for i := range w.Ops {
		op := &w.Ops[i]
		if op.Kind == bench.OpCancel {
			allCmds[i] = engine.Command{Type: engine.CmdCancelOrder, OrderID: op.CancelID}
		} else {
			allCmds[i] = engine.Command{Type: engine.CmdPlaceOrder, Order: op.Req}
		}
	}
	warmupCmds := make([]engine.Command, len(w.Warmup))
	for i := range w.Warmup {
		warmupCmds[i] = engine.Command{Type: engine.CmdPlaceOrder, Order: w.Warmup[i].Req}
	}

	var totalEvents int64
	startBatchEngine := func() *engine.BatchEngine {
		be := engine.NewBatchEngine(bk, engine.Options{
			QueueSize: 512,
			EventHandler: func(events []book.Event) {
				totalEvents += int64(len(events))
			},
		})
		be.Start()
		be.SubmitBatch(warmupCmds)
		ch := make(chan engine.Response, 1)
		be.SubmitBatch([]engine.Command{{Type: engine.CmdStats, Reply: ch}})
		<-ch
		return be
	}

	be := startBatchEngine()
	idx := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; {
		rem := b.N - i
		end := idx + batchSize
		if end > nOps {
			end = nOps
		}
		if end-idx > rem {
			end = idx + rem
		}
		be.SubmitBatch(allCmds[idx:end])
		i += end - idx
		idx = end
		if idx == nOps && i < b.N {
			b.StopTimer()
			be.Stop()
			bk.Reset()
			be = startBatchEngine()
			idx = 0
			b.StartTimer()
		}
	}
	be.Stop()
	b.StopTimer()
	sinkInt64 = totalEvents
}

func BenchmarkEngineDisruptor(b *testing.B) {
	w := getSharedWorkload()
	v := bench.Variant{Name: "ladder_linkedlist_pool", Levels: "ladder", Queue: "linkedlist", Pool: true}
	bk, err := v.NewBook(w.Config, len(w.Warmup)+32768)
	if err != nil {
		b.Fatal(err)
	}
	evBuf := make([]book.Event, 0, 64)
	evBuf = w.ApplyWarmup(bk, evBuf)
	for i := range w.Ops {
		evBuf = bench.ReplayOp(bk, &w.Ops[i], evBuf)
	}
	bk.Reset()

	var totalEvents int64
	startDisruptor := func() *engine.Disruptor {
		d := engine.NewDisruptor(bk, engine.Options{
			QueueSize: 8192,
			EventHandler: func(events []book.Event) {
				totalEvents += int64(len(events))
			},
		})
		d.Start()
		for i := range w.Warmup {
			d.PlaceOrderAsync(w.Warmup[i].Req)
		}
		_ = d.Stats()
		return d
	}

	d := startDisruptor()
	ops := w.Ops
	nOps := len(ops)
	idx := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := &ops[idx]
		if op.Kind == bench.OpCancel {
			d.CancelOrderAsync(op.CancelID)
		} else {
			d.PlaceOrderAsync(op.Req)
		}
		idx++
		if idx == nOps {
			b.StopTimer()
			d.Stop()
			bk.Reset()
			d = startDisruptor()
			idx = 0
			b.StartTimer()
		}
	}
	d.Stop()
	b.StopTimer()
	sinkInt64 = totalEvents
}

func TestWorkloadAndLatencyHarness(t *testing.T) {
	w := bench.GenerateMixedWorkload(bench.DefaultConfig(), 42, 500, 5000)
	if w.Stats.LimitOps == 0 || w.Stats.CancelOps == 0 || w.Stats.MarketOps == 0 {
		t.Fatalf("expected all op types in workload: %+v", w.Stats)
	}
	if w.Stats.Trades == 0 {
		t.Fatalf("expected trades in mixed workload: %+v", w.Stats)
	}
	v := bench.AllVariants()[0]
	rep, err := bench.MeasureVariantLatency(v, w)
	if err != nil {
		t.Fatal(err)
	}
	if rep.All.Count != 5000 || rep.ThroughputOps <= 0 {
		t.Fatalf("unexpected latency report: %+v", rep)
	}
	var buf bytes.Buffer
	bench.PrintLatencyReports(&buf, w, []bench.VariantLatencyReport{rep}, false)
	if buf.Len() == 0 {
		t.Fatalf("expected formatted report output")
	}
}
