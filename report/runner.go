package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/unedtamps/orbook/bench"
)

// RunOptions configures a report generation run (full or single-benchmark).
type RunOptions struct {
	OutDir       string  // Root output directory (default "reports").
	Mode         string  // "full" or "single".
	BenchFilter  string  // Regex passed to -bench (e.g. "." or "^BenchmarkAddLimit/ladder_linkedlist_pool$").
	BenchTime    string  // Passed to -benchtime for benchmark timing (default "150ms").
	ProfileMode  string  // "all", "methods", "single", or "none". In Mode="single", always profiles all matched variants.
	ProfileBench string  // Specific benchmark regex when ProfileMode="single".
	ProfileTime  string  // Passed to -benchtime during profiling pass (default "300ms" for all, "500ms" otherwise).
	IncludeHDR   bool    // Run HDR histogram latency suite (true in full mode).
	HDROps       int     // Number of ops for HDR latency (default 200000).
	HDRWarmup    int     // Warmup orders for HDR latency (default 2000).
	Seed         int64   // Workload RNG seed (default 20261122).
	ThresholdPct float64 // Min % ns/op change to count as a diff (default 2.0).
	Force        bool    // Always create a new version even if no diff.
	Progress     func(msg string)
}

// RunResult summarizes the outcome of Run().
type RunResult struct {
	Report     *Report
	Saved      bool
	VersionDir string
}

var (
	benchLineRe = regexp.MustCompile(`^(Benchmark\S+?)(?:-\d+)?\s+(\d+)\s+([0-9.]+)\s+ns/op(?:\s+(\d+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?`)
	cpuLineRe   = regexp.MustCompile(`^cpu:\s+(.+)$`)
	pprofRowRe  = regexp.MustCompile(`^\s*(\S+)\s+([0-9.]+)%\s+[0-9.]+%\s+(\S+)\s+([0-9.]+)%\s+(.+)$`)
)

// ExactBenchRegex converts a full benchmark name like
// "BenchmarkAddLimit/ladder_linkedlist_pool" into a Go testing -bench pattern
// that matches that exact sub-benchmark ("^BenchmarkAddLimit$/^ladder_linkedlist_pool$").
func ExactBenchRegex(fullName string) string {
	parts := strings.Split(fullName, "/")
	for i, p := range parts {
		parts[i] = "^" + regexp.QuoteMeta(p) + "$"
	}
	return strings.Join(parts, "/")
}

func safeFileName(fullName string) string {
	r := strings.NewReplacer("/", "__", " ", "_", ":", "_")
	return r.Replace(fullName)
}

// Run executes benchmarks and profiles, diffs against reports/latest.json, and
// saves a new versioned snapshot if a diff is detected (or Force is set).
func Run(opts RunOptions) (*RunResult, error) {
	if opts.OutDir == "" {
		opts.OutDir = "reports"
	}
	if opts.Mode == "" {
		opts.Mode = "full"
	}
	if opts.BenchFilter == "" {
		opts.BenchFilter = "."
	}
	if opts.BenchTime == "" {
		opts.BenchTime = "150ms"
	}
	if opts.ProfileMode == "" {
		if opts.Mode == "single" {
			opts.ProfileMode = "all"
		} else {
			opts.ProfileMode = "all"
		}
	}
	if opts.ProfileTime == "" {
		if opts.ProfileMode == "all" && opts.Mode == "full" {
			opts.ProfileTime = "250ms"
		} else {
			opts.ProfileTime = "500ms"
		}
	}
	if opts.HDROps <= 0 {
		opts.HDROps = 200_000
	}
	if opts.HDRWarmup <= 0 {
		opts.HDRWarmup = 2000
	}
	if opts.Seed == 0 {
		opts.Seed = 20261122
	}
	if opts.ThresholdPct <= 0 {
		opts.ThresholdPct = 2.0
	}
	logf := func(format string, args ...any) {
		if opts.Progress != nil {
			opts.Progress(fmt.Sprintf(format, args...))
		}
	}

	prev, _ := Load(filepath.Join(opts.OutDir, "latest.json"))

	tmpDir, err := os.MkdirTemp("", "orbook-report-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Compile ./bench test binary ONCE to avoid repeated go toolchain startup overhead.
	logf("Compiling ./bench test binary...")
	benchBin := filepath.Join(tmpDir, "bench.test")
	if err := compileBenchBinary(benchBin); err != nil {
		return nil, err
	}

	// 1. Run clean timing benchmarks (no profiling overhead skewing ns/op).
	logf("Running benchmarks (filter=%q, benchtime=%s)...", opts.BenchFilter, opts.BenchTime)
	newBenches, cpuModel, err := runBenchBinary(benchBin, opts.BenchFilter, opts.BenchTime, "", "")
	if err != nil {
		return nil, err
	}
	if len(newBenches) == 0 {
		return nil, fmt.Errorf("no benchmarks matched filter %q", opts.BenchFilter)
	}

	// 2. Determine which benchmarks to profile.
	targets := selectProfileTargets(newBenches, opts)
	type profJob struct {
		bench   BenchResult
		cpuPath string
		memPath string
		safe    string
	}
	jobs := make([]profJob, 0, len(targets))

	if len(targets) > 0 {
		logf("Profiling %d benchmark(s) (benchtime=%s each)...", len(targets), opts.ProfileTime)
		for i, t := range targets {
			safe := safeFileName(t.FullName)
			cpuProf := filepath.Join(tmpDir, safe+".cpu.prof")
			memProf := filepath.Join(tmpDir, safe+".mem.prof")
			if _, _, err := runBenchBinary(benchBin, ExactBenchRegex(t.FullName), opts.ProfileTime, cpuProf, memProf); err != nil {
				return nil, fmt.Errorf("profiling %s failed: %w", t.FullName, err)
			}
			jobs = append(jobs, profJob{
				bench:   t,
				cpuPath: cpuProf,
				memPath: memProf,
				safe:    safe,
			})
			if (i+1)%12 == 0 || i+1 == len(targets) {
				logf("  profiled %d/%d (%s)", i+1, len(targets), t.FullName)
			}
		}
	}

	// Extract pprof top tables concurrently (read-only post-processing; does not affect measurements).
	newProfiles := make([]ProfileSummary, len(jobs))
	if len(jobs) > 0 {
		logf("Extracting pprof hotspots for %d profile(s)...", len(jobs))
		workers := 4
		if len(jobs) < workers {
			workers = len(jobs)
		}
		ch := make(chan int, len(jobs))
		for i := range jobs {
			ch <- i
		}
		close(ch)
		var wg sync.WaitGroup
		wg.Add(workers)
		for w := 0; w < workers; w++ {
			go func() {
				defer wg.Done()
				for idx := range ch {
					j := jobs[idx]
					newProfiles[idx] = extractProfile(benchBin, j.bench, j.cpuPath, j.memPath)
				}
			}()
		}
		wg.Wait()
	}

	// 3. Optional HDR latency suite (all variants in full mode; only matched variants in single mode).
	var newLatency []LatencySummary
	if opts.IncludeHDR {
		hdrVariants := selectHDRVariants(newBenches, opts.Mode)
		if len(hdrVariants) > 0 {
			logf("Running HDR histogram latency suite (%d ops x %d variants)...", opts.HDROps, len(hdrVariants))
			w := bench.GenerateMixedWorkload(bench.DefaultConfig(), opts.Seed, opts.HDRWarmup, opts.HDROps)
			for _, v := range hdrVariants {
				rep, err := bench.MeasureVariantLatency(v, w)
				if err != nil {
					return nil, err
				}
				newLatency = append(newLatency, LatencySummary{
					Variant:       rep.Variant.Name,
					ThroughputOps: rep.ThroughputOps,
					AvgNs:         rep.AvgNsPerOp,
					P50Ns:         rep.All.P50,
					P90Ns:         rep.All.P90,
					P99Ns:         rep.All.P99,
					P999Ns:        rep.All.P999,
					MaxNs:         rep.All.Max,
				})
			}
		}
	}

	// 4. Compute diff against previous latest.json.
	diff := ComputeDiff(prev, newBenches, opts.ThresholdPct)

	nextVersion := 1
	if prev != nil {
		nextVersion = prev.Version + 1
	}
	now := time.Now()
	runID := fmt.Sprintf("v%03d_%s", nextVersion, now.Format("20060102_150405"))

	// Merge into centralized dataset if single mode and prev exists.
	mergedBenches := newBenches
	mergedLatency := newLatency
	mergedProfiles := newProfiles
	if opts.Mode == "single" && prev != nil {
		mergedBenches = mergeBenchResults(prev.Benchmarks, newBenches)
		mergedLatency = mergeLatency(prev.Latency, newLatency)
		mergedProfiles = mergeProfiles(prev.Profiles, newProfiles)
	}

	rep := &Report{
		Version: nextVersion,
		RunID:   runID,
		Mode:    opts.Mode,
		Filter:  opts.BenchFilter,
		Meta: Metadata{
			Timestamp: now,
			GitCommit: gitShortCommit(),
			GoVersion: runtime.Version(),
			GOOS:      runtime.GOOS,
			GOARCH:    runtime.GOARCH,
			CPU:       cpuModel,
			BenchTime: opts.BenchTime,
			Seed:      opts.Seed,
			OpsCount:  opts.HDROps,
		},
		Benchmarks: mergedBenches,
		Latency:    mergedLatency,
		Profiles:   mergedProfiles,
		Diff:       diff,
	}

	// 5. Only create a new version if there is a diff (or first run, or Force).
	if prev != nil && !diff.HasDiff && !opts.Force {
		rep.Version = prev.Version
		rep.RunID = prev.RunID
		return &RunResult{Report: rep, Saved: false}, nil
	}

	versionDir := filepath.Join(opts.OutDir, "history", runID)
	profOutDir := filepath.Join(versionDir, "profiles")
	if err := os.MkdirAll(profOutDir, 0o755); err != nil {
		return nil, err
	}

	// Copy newly generated .prof files into history/<runID>/profiles/ and record relative paths.
	jobByName := make(map[string]profJob, len(jobs))
	for _, j := range jobs {
		jobByName[j.bench.FullName] = j
	}
	for i := range rep.Profiles {
		j, ok := jobByName[rep.Profiles[i].Benchmark]
		if !ok {
			continue
		}
		if _, err := os.Stat(j.cpuPath); err == nil {
			dst := filepath.Join(profOutDir, j.safe+".cpu.prof")
			if copyFile(j.cpuPath, dst) == nil {
				rep.Profiles[i].CPUFile = filepath.Join("history", runID, "profiles", j.safe+".cpu.prof")
			}
		}
		if _, err := os.Stat(j.memPath); err == nil {
			dst := filepath.Join(profOutDir, j.safe+".mem.prof")
			if copyFile(j.memPath, dst) == nil {
				rep.Profiles[i].MemFile = filepath.Join("history", runID, "profiles", j.safe+".mem.prof")
			}
		}
	}

	jsonBytes, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return nil, err
	}
	mdBytes := []byte(RenderMarkdown(rep))
	htmlBytes, err := RenderHTML(rep)
	if err != nil {
		return nil, err
	}

	files := map[string][]byte{
		filepath.Join(versionDir, "results.json"):  jsonBytes,
		filepath.Join(versionDir, "report.md"):     mdBytes,
		filepath.Join(versionDir, "report.html"):   htmlBytes,
		filepath.Join(opts.OutDir, "latest.json"):  jsonBytes,
		filepath.Join(opts.OutDir, "latest.md"):    mdBytes,
		filepath.Join(opts.OutDir, "latest.html"):  htmlBytes,
	}
	for path, data := range files {
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return nil, err
		}
	}

	return &RunResult{
		Report:     rep,
		Saved:      true,
		VersionDir: versionDir,
	}, nil
}

func selectProfileTargets(benches []BenchResult, opts RunOptions) []BenchResult {
	if opts.Mode == "single" {
		// In single mode, profile EVERY matched benchmark variant.
		return benches
	}
	switch opts.ProfileMode {
	case "none":
		return nil
	case "all":
		return benches
	case "methods":
		// Pick the fastest variant for each distinct method.
		bestByMethod := make(map[string]BenchResult)
		order := make([]string, 0, 16)
		for _, b := range benches {
			cur, ok := bestByMethod[b.Method]
			if !ok {
				order = append(order, b.Method)
				bestByMethod[b.Method] = b
			} else if b.NsPerOp < cur.NsPerOp {
				bestByMethod[b.Method] = b
			}
		}
		out := make([]BenchResult, 0, len(order))
		for _, m := range order {
			out = append(out, bestByMethod[m])
		}
		return out
	case "single":
		target := opts.ProfileBench
		if target == "" {
			target = "BenchmarkMixedWorkload/ladder_linkedlist_pool"
		}
		target = strings.Trim(target, "^$")
		for _, b := range benches {
			if b.FullName == target {
				return []BenchResult{b}
			}
		}
		if len(benches) > 0 {
			return []BenchResult{benches[0]}
		}
		return nil
	default:
		return benches
	}
}

func compileBenchBinary(outPath string) error {
	cmd := exec.Command("go", "test", "-c", "-o", outPath, "./bench")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compile bench binary failed: %w\n%s", err, out.String())
	}
	return nil
}

func runBenchBinary(binPath, filter, benchTime, cpuProf, memProf string) ([]BenchResult, string, error) {
	args := []string{
		"-test.run=^$",
		"-test.bench=" + filter,
		"-test.benchmem",
		"-test.benchtime=" + benchTime,
		"-test.count=1",
	}
	if cpuProf != "" {
		args = append(args, "-test.cpuprofile="+cpuProf)
	}
	if memProf != "" {
		args = append(args, "-test.memprofile="+memProf)
	}

	cmd := exec.Command(binPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, "", fmt.Errorf("bench run (%s) failed: %w\n%s", filter, err, out.String())
	}
	benches, cpu := ParseBenchOutput(out.String())
	return benches, cpu, nil
}

// Load reads a saved Report JSON file from disk.
func Load(path string) (*Report, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rep Report
	if err := json.Unmarshal(data, &rep); err != nil {
		return nil, err
	}
	return &rep, nil
}

// ComputeDiff compares updated benchmark results against prev.
func ComputeDiff(prev *Report, updated []BenchResult, thresholdPct float64) *DiffSummary {
	if prev == nil {
		items := make([]BenchDiff, 0, len(updated))
		for _, b := range updated {
			items = append(items, BenchDiff{
				FullName:  b.FullName,
				Status:    "added",
				NewNs:     b.NsPerOp,
				NewBytes:  b.BytesPerOp,
				NewAllocs: b.AllocsOp,
			})
		}
		return &DiffSummary{
			ThresholdPct: thresholdPct,
			HasDiff:      len(updated) > 0,
			AddedCount:   len(updated),
			Items:        items,
		}
	}

	prevMap := make(map[string]BenchResult, len(prev.Benchmarks))
	for _, b := range prev.Benchmarks {
		prevMap[b.FullName] = b
	}

	ds := &DiffSummary{
		BaseVersion:  prev.Version,
		BaseRunID:    prev.RunID,
		ThresholdPct: thresholdPct,
	}

	for _, cur := range updated {
		old, exists := prevMap[cur.FullName]
		if !exists {
			ds.HasDiff = true
			ds.AddedCount++
			ds.Items = append(ds.Items, BenchDiff{
				FullName:  cur.FullName,
				Status:    "added",
				NewNs:     cur.NsPerOp,
				NewBytes:  cur.BytesPerOp,
				NewAllocs: cur.AllocsOp,
			})
			continue
		}

		var deltaPct float64
		if old.NsPerOp > 0 {
			deltaPct = ((cur.NsPerOp - old.NsPerOp) / old.NsPerOp) * 100.0
		}
		allocChanged := cur.BytesPerOp != old.BytesPerOp || cur.AllocsOp != old.AllocsOp
		timeChanged := math.Abs(deltaPct) >= thresholdPct

		status := "unchanged"
		if timeChanged {
			if deltaPct < 0 {
				status = "faster"
			} else {
				status = "slower"
			}
		} else if allocChanged {
			status = "alloc_changed"
		}

		if timeChanged || allocChanged {
			ds.HasDiff = true
			ds.ChangedCount++
		}

		ds.Items = append(ds.Items, BenchDiff{
			FullName:   cur.FullName,
			Status:     status,
			OldNs:      old.NsPerOp,
			NewNs:      cur.NsPerOp,
			DeltaNsPct: deltaPct,
			OldBytes:   old.BytesPerOp,
			NewBytes:   cur.BytesPerOp,
			OldAllocs:  old.AllocsOp,
			NewAllocs:  cur.AllocsOp,
		})
	}
	return ds
}

// ParseBenchOutput parses standard `go test -bench -benchmem` text output.
func ParseBenchOutput(raw string) ([]BenchResult, string) {
	var results []BenchResult
	var cpu string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if m := cpuLineRe.FindStringSubmatch(line); m != nil {
			cpu = strings.TrimSpace(m[1])
			continue
		}
		m := benchLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		fullName := m[1]
		method := fullName
		variant := "default"
		if idx := strings.IndexByte(fullName, '/'); idx >= 0 {
			method = fullName[:idx]
			variant = fullName[idx+1:]
		}
		iters, _ := strconv.ParseInt(m[2], 10, 64)
		nsOp, _ := strconv.ParseFloat(m[3], 64)
		var bOp, aOp int64
		if m[4] != "" {
			bOp, _ = strconv.ParseInt(m[4], 10, 64)
		}
		if m[5] != "" {
			aOp, _ = strconv.ParseInt(m[5], 10, 64)
		}
		results = append(results, BenchResult{
			FullName:   fullName,
			Method:     method,
			Variant:    variant,
			Iterations: iters,
			NsPerOp:    nsOp,
			BytesPerOp: bOp,
			AllocsOp:   aOp,
		})
	}
	return results, cpu
}

func extractProfile(benchBin string, b BenchResult, cpuProf, memProf string) ProfileSummary {
	ps := ProfileSummary{
		Benchmark: b.FullName,
		Method:    b.Method,
		Variant:   b.Variant,
	}
	if cpuProf != "" {
		if out, err := runPprofTop(benchBin, cpuProf, ""); err == nil {
			ps.CPUTop = ParsePprofTop(out)
		}
	}
	if memProf != "" {
		if out, err := runPprofTop(benchBin, memProf, "-alloc_space"); err == nil {
			ps.MemSpaceTop = ParsePprofTop(out)
		}
		if out, err := runPprofTop(benchBin, memProf, "-alloc_objects"); err == nil {
			ps.MemObjectTop = ParsePprofTop(out)
		}
	}
	return ps
}

func runPprofTop(benchBin, profPath, sampleFlag string) (string, error) {
	args := []string{"tool", "pprof", "-top", "-nodecount=15"}
	if sampleFlag != "" {
		args = append(args, sampleFlag)
	}
	if benchBin != "" {
		args = append(args, benchBin)
	}
	args = append(args, profPath)
	cmd := exec.Command("go", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ParsePprofTop parses `go tool pprof -top` table rows.
func ParsePprofTop(raw string) []PprofEntry {
	var entries []PprofEntry
	for _, line := range strings.Split(raw, "\n") {
		m := pprofRowRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		fp, _ := strconv.ParseFloat(m[2], 64)
		cp, _ := strconv.ParseFloat(m[4], 64)
		entries = append(entries, PprofEntry{
			Flat:    m[1],
			FlatPct: fp,
			Cum:     m[3],
			CumPct:  cp,
			Func:    strings.TrimSpace(m[5]),
		})
	}
	return entries
}

func mergeBenchResults(prev, updated []BenchResult) []BenchResult {
	byName := make(map[string]BenchResult, len(prev)+len(updated))
	order := make([]string, 0, len(prev)+len(updated))
	for _, b := range prev {
		byName[b.FullName] = b
		order = append(order, b.FullName)
	}
	for _, b := range updated {
		if _, exists := byName[b.FullName]; !exists {
			order = append(order, b.FullName)
		}
		byName[b.FullName] = b
	}
	out := make([]BenchResult, 0, len(order))
	for _, name := range order {
		out = append(out, byName[name])
	}
	return out
}

func selectHDRVariants(benches []BenchResult, mode string) []bench.Variant {
	all := bench.AllVariants()
	if mode != "single" {
		return all
	}
	matched := make(map[string]bool)
	for _, b := range benches {
		matched[b.Variant] = true
	}
	out := make([]bench.Variant, 0, len(matched))
	for _, v := range all {
		if matched[v.Name] {
			out = append(out, v)
		}
	}
	return out
}

func mergeLatency(prev, updated []LatencySummary) []LatencySummary {
	if len(updated) == 0 {
		return prev
	}
	byVar := make(map[string]LatencySummary, len(prev)+len(updated))
	order := make([]string, 0, len(prev)+len(updated))
	for _, l := range prev {
		byVar[l.Variant] = l
		order = append(order, l.Variant)
	}
	for _, l := range updated {
		if _, exists := byVar[l.Variant]; !exists {
			order = append(order, l.Variant)
		}
		byVar[l.Variant] = l
	}
	out := make([]LatencySummary, 0, len(order))
	for _, v := range order {
		out = append(out, byVar[v])
	}
	return out
}

func mergeProfiles(prev, updated []ProfileSummary) []ProfileSummary {
	byName := make(map[string]ProfileSummary, len(prev)+len(updated))
	for _, p := range prev {
		byName[p.Benchmark] = p
	}
	for _, p := range updated {
		byName[p.Benchmark] = p
	}
	keys := make([]string, 0, len(byName))
	for k := range byName {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]ProfileSummary, 0, len(keys))
	for _, k := range keys {
		out = append(out, byName[k])
	}
	return out
}

func gitShortCommit() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
