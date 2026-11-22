package report

import "time"

// Report is the single source of truth saved to reports/latest.json and
// archived under reports/history/vNNN_<timestamp>/results.json.
type Report struct {
	Version    int                  `json:"version"`
	RunID      string               `json:"run_id"`
	Mode       string               `json:"mode"` // "full" or "single"
	Filter     string               `json:"filter,omitempty"`
	Meta       Metadata             `json:"meta"`
	Benchmarks []BenchResult        `json:"benchmarks"`
	Latency    []LatencySummary     `json:"latency,omitempty"`
	Profiles   []ProfileSummary     `json:"profiles,omitempty"`
	Diff       *DiffSummary         `json:"diff,omitempty"`
}

// Metadata captures environment and reproducibility details.
type Metadata struct {
	Timestamp time.Time `json:"timestamp"`
	GitCommit string    `json:"git_commit"`
	GoVersion string    `json:"go_version"`
	GOOS      string    `json:"goos"`
	GOARCH    string    `json:"goarch"`
	CPU       string    `json:"cpu"`
	BenchTime string    `json:"bench_time"`
	Seed      int64     `json:"seed,omitempty"`
	OpsCount  int       `json:"ops_count,omitempty"`
}

// BenchResult stores one parsed `go test -bench -benchmem` row.
type BenchResult struct {
	FullName   string  `json:"full_name"` // e.g. BenchmarkAddLimit/ladder_linkedlist_pool
	Method     string  `json:"method"`    // e.g. BenchmarkAddLimit
	Variant    string  `json:"variant"`   // e.g. ladder_linkedlist_pool (or "default")
	Iterations int64   `json:"iterations"`
	NsPerOp    float64 `json:"ns_per_op"`
	BytesPerOp int64   `json:"bytes_per_op"`
	AllocsOp   int64   `json:"allocs_per_op"`
}

// LatencySummary mirrors bench.VariantLatencyReport in a flat JSON-friendly form.
type LatencySummary struct {
	Variant       string  `json:"variant"`
	ThroughputOps float64 `json:"throughput_ops_per_sec"`
	AvgNs         float64 `json:"avg_ns"`
	P50Ns         int64   `json:"p50_ns"`
	P90Ns         int64   `json:"p90_ns"`
	P99Ns         int64   `json:"p99_ns"`
	P999Ns        int64   `json:"p999_ns"`
	MaxNs         int64   `json:"max_ns"`
}

// ProfileSummary stores top CPU and Memory hotspots for a profiled benchmark.
type ProfileSummary struct {
	Benchmark    string       `json:"benchmark"` // FullName e.g. BenchmarkAddLimit/ladder_linkedlist_pool
	Method       string       `json:"method"`
	Variant      string       `json:"variant"`
	CPUFile      string       `json:"cpu_file,omitempty"`
	MemFile      string       `json:"mem_file,omitempty"`
	CPUTop       []PprofEntry `json:"cpu_top,omitempty"`
	MemSpaceTop  []PprofEntry `json:"mem_space_top,omitempty"`
	MemObjectTop []PprofEntry `json:"mem_objects_top,omitempty"`
}

// PprofEntry is one parsed row from `go tool pprof -top`.
type PprofEntry struct {
	Flat    string  `json:"flat"`
	FlatPct float64 `json:"flat_pct"`
	Cum     string  `json:"cum"`
	CumPct  float64 `json:"cum_pct"`
	Func    string  `json:"func"`
}

// DiffSummary captures changes compared to the previous version.
type DiffSummary struct {
	BaseVersion   int         `json:"base_version"`
	BaseRunID     string      `json:"base_run_id"`
	ThresholdPct  float64     `json:"threshold_pct"`
	HasDiff       bool        `json:"has_diff"`
	ChangedCount  int         `json:"changed_count"`
	AddedCount    int         `json:"added_count"`
	Items         []BenchDiff `json:"items,omitempty"`
}

// BenchDiff describes the delta for a single benchmark between base and current.
type BenchDiff struct {
	FullName      string  `json:"full_name"`
	Status        string  `json:"status"` // "added", "faster", "slower", "alloc_changed", "unchanged"
	OldNs         float64 `json:"old_ns"`
	NewNs         float64 `json:"new_ns"`
	DeltaNsPct    float64 `json:"delta_ns_pct"`
	OldBytes      int64   `json:"old_bytes"`
	NewBytes      int64   `json:"new_bytes"`
	OldAllocs     int64   `json:"old_allocs"`
	NewAllocs     int64   `json:"new_allocs"`
}
