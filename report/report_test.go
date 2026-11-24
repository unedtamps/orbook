package report_test

import (
	"strings"
	"testing"
	"time"

	"github.com/unedtamps/orbook/report"
)

func TestParseBenchOutput(t *testing.T) {
	raw := `goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i9-14900HX
BenchmarkAddLimit/ladder_linkedlist_pool-32    6156284    30.07 ns/op    0 B/op    0 allocs/op
BenchmarkEngineChannel-32                      3800502    92.25 ns/op    0 B/op    0 allocs/op
PASS`

	benches, cpu := report.ParseBenchOutput(raw)
	if cpu != "Intel(R) Core(TM) i9-14900HX" {
		t.Fatalf("got cpu %q", cpu)
	}
	if len(benches) != 2 {
		t.Fatalf("got %d benches, want 2", len(benches))
	}
	if benches[0].FullName != "BenchmarkAddLimit/ladder_linkedlist_pool" || benches[0].Method != "BenchmarkAddLimit" ||
		benches[0].Variant != "ladder_linkedlist_pool" || benches[0].NsPerOp != 30.07 {
		t.Fatalf("unexpected bench[0]: %+v", benches[0])
	}
	if benches[1].FullName != "BenchmarkEngineChannel" || benches[1].Variant != "default" {
		t.Fatalf("unexpected bench[1]: %+v", benches[1])
	}
}

func TestParsePprofTop(t *testing.T) {
	raw := `Showing top 2 nodes out of 50
      flat  flat%   sum%        cum   cum%
     600ms 21.43% 21.43%     1750ms 62.50%  github.com/unedtamps/orbook/book.(*Book).PlaceOrder
     220ms  7.86% 29.29%      730ms 26.07%  github.com/unedtamps/orbook/book.(*Book).CancelOrder`

	entries := report.ParsePprofTop(raw)
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}
	if entries[0].Flat != "600ms" || entries[0].FlatPct != 21.43 || entries[0].CumPct != 62.50 ||
		!strings.Contains(entries[0].Func, "PlaceOrder") {
		t.Fatalf("unexpected entry[0]: %+v", entries[0])
	}
}

func TestComputeDiffVersioning(t *testing.T) {
	prev := &report.Report{
		Version: 1,
		RunID:   "v001_20261123_010000",
		Benchmarks: []report.BenchResult{
			{FullName: "BenchmarkAddLimit/ladder_linkedlist_pool", NsPerOp: 30.0, BytesPerOp: 0, AllocsOp: 0},
			{FullName: "BenchmarkCancel/ladder_linkedlist_pool", NsPerOp: 40.0, BytesPerOp: 0, AllocsOp: 0},
		},
	}

	same := []report.BenchResult{
		{FullName: "BenchmarkAddLimit/ladder_linkedlist_pool", NsPerOp: 30.4, BytesPerOp: 0, AllocsOp: 0},
	}
	d1 := report.ComputeDiff(prev, same, 2.0)
	if d1.HasDiff {
		t.Fatalf("expected HasDiff=false for +1.3%% change, got %+v", d1)
	}

	changed := []report.BenchResult{
		{FullName: "BenchmarkAddLimit/ladder_linkedlist_pool", NsPerOp: 25.0, BytesPerOp: 0, AllocsOp: 0},
		{FullName: "BenchmarkCancel/ladder_linkedlist_pool", NsPerOp: 40.1, BytesPerOp: 96, AllocsOp: 1},
	}
	d2 := report.ComputeDiff(prev, changed, 2.0)
	if !d2.HasDiff || d2.ChangedCount != 2 {
		t.Fatalf("expected HasDiff=true ChangedCount=2, got %+v", d2)
	}
	if d2.Items[0].Status != "faster" || d2.Items[1].Status != "alloc_changed" {
		t.Fatalf("unexpected diff statuses: %+v", d2.Items)
	}
}

func TestRenderMarkdownAndHTML(t *testing.T) {
	rep := &report.Report{
		Version: 2,
		RunID:   "v002_20261123_015000",
		Mode:    "single",
		Filter:  "BenchmarkAddLimit/ladder_linkedlist_pool",
		Meta: report.Metadata{
			Timestamp: time.Now(),
			GitCommit: "abc1234",
			GoVersion: "go1.27",
			CPU:       "TestCPU",
			BenchTime: "150ms",
		},
		Benchmarks: []report.BenchResult{
			{FullName: "BenchmarkAddLimit/ladder_linkedlist_pool", Method: "BenchmarkAddLimit", Variant: "ladder_linkedlist_pool", NsPerOp: 30.1},
		},
	}

	md := report.RenderMarkdown(rep)
	if !strings.Contains(md, "v002_20261123_015000") || !strings.Contains(md, "30.10") {
		t.Fatalf("markdown missing expected content:\n%s", md)
	}

	htmlBytes, err := report.RenderHTML(rep)
	if err != nil {
		t.Fatalf("RenderHTML failed: %v", err)
	}
	if !strings.Contains(string(htmlBytes), "BenchmarkAddLimit/ladder_linkedlist_pool") {
		t.Fatalf("HTML missing embedded benchmark data")
	}
	if !strings.Contains(string(htmlBytes), "Side-by-Side Variant Comparison") {
		t.Fatalf("HTML missing side-by-side comparison section")
	}
}

func TestExactBenchRegex(t *testing.T) {
	if got := report.ExactBenchRegex("BenchmarkAddLimit/ladder_linkedlist_pool"); got != "^BenchmarkAddLimit$/^ladder_linkedlist_pool$" {
		t.Fatalf("got %q", got)
	}
	if got := report.ExactBenchRegex("BenchmarkAddLimit"); got != "^BenchmarkAddLimit$" {
		t.Fatalf("got %q", got)
	}
}
