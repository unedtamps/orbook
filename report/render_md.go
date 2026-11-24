package report

import (
	"fmt"
	"strings"
)

// RenderMarkdown produces a concise, agent- and PR-friendly Markdown report
// from a Report snapshot.
func RenderMarkdown(r *Report) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# Orbook Benchmark & Profile Report — `%s`\n\n", r.RunID)
	fmt.Fprintf(&sb, "- **Version:** `v%03d` | **Mode:** `%s`", r.Version, r.Mode)
	if r.Filter != "" && r.Filter != "." {
		fmt.Fprintf(&sb, " | **Filter:** `%s`", r.Filter)
	}
	fmt.Fprintf(&sb, "\n- **Time:** `%s` | **Commit:** `%s` | **Go:** `%s` | **CPU:** `%s`\n\n",
		r.Meta.Timestamp.Format("2006-01-02 15:04:05 MST"),
		r.Meta.GitCommit, r.Meta.GoVersion, r.Meta.CPU)

	if r.Diff != nil && len(r.Diff.Items) > 0 {
		fmt.Fprintf(&sb, "## Diff vs Previous (`%s`, threshold ±%.1f%%)\n\n", r.Diff.BaseRunID, r.Diff.ThresholdPct)
		sb.WriteString("| Benchmark | Status | Old (ns/op) | New (ns/op) | Delta | Allocs |\n")
		sb.WriteString("|---|---|---:|---:|---:|---|\n")
		for _, d := range r.Diff.Items {
			if d.Status == "unchanged" {
				continue
			}
			allocStr := fmt.Sprintf("%d→%d B, %d→%d allocs", d.OldBytes, d.NewBytes, d.OldAllocs, d.NewAllocs)
			if d.Status == "added" {
				fmt.Fprintf(&sb, "| `%s` | `added` | — | %.2f | — | %d B, %d allocs |\n",
					d.FullName, d.NewNs, d.NewBytes, d.NewAllocs)
			} else {
				fmt.Fprintf(&sb, "| `%s` | `%s` | %.2f | %.2f | %+.2f%% | %s |\n",
					d.FullName, d.Status, d.OldNs, d.NewNs, d.DeltaNsPct, allocStr)
			}
		}
		sb.WriteString("\n")
	}

	if len(r.Benchmarks) > 0 {
		sb.WriteString("## Benchmarks\n\n")
		sb.WriteString("| Method | Variant | ns/op | B/op | allocs/op |\n")
		sb.WriteString("|---|---|---:|---:|---:|\n")
		for _, b := range r.Benchmarks {
			fmt.Fprintf(&sb, "| `%s` | `%s` | %.2f | %d | %d |\n",
				b.Method, b.Variant, b.NsPerOp, b.BytesPerOp, b.AllocsOp)
		}
		sb.WriteString("\n")
	}

	if len(r.Latency) > 0 {
		sb.WriteString("## HDR Histogram Latency (Mixed Workload)\n\n")
		sb.WriteString("| Variant | Throughput (ops/s) | Avg (ns) | p50 (ns) | p90 (ns) | p99 (ns) | p99.9 (ns) | Max (ns) |\n")
		sb.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|\n")
		for _, l := range r.Latency {
			fmt.Fprintf(&sb, "| `%s` | %.0f | %.1f | %d | %d | %d | %d | %d |\n",
				l.Variant, l.ThroughputOps, l.AvgNs, l.P50Ns, l.P90Ns, l.P99Ns, l.P999Ns, l.MaxNs)
		}
		sb.WriteString("\n")
	}

	if len(r.Profiles) > 0 {
		fmt.Fprintf(&sb, "## Profiling Hotspots (pprof — %d profiled benchmarks)\n\n", len(r.Profiles))
		if len(r.Profiles) > 15 {
			sb.WriteString("| Benchmark | Top CPU Hotspot | Top Mem Alloc Space |\n")
			sb.WriteString("|---|---|---|\n")
			for _, p := range r.Profiles {
				cpuTop := "—"
				if len(p.CPUTop) > 0 {
					cpuTop = fmt.Sprintf("`%s` (%.1f%% flat)", p.CPUTop[0].Func, p.CPUTop[0].FlatPct)
				}
				memTop := "—"
				if len(p.MemSpaceTop) > 0 {
					memTop = fmt.Sprintf("`%s` (%s, %.1f%%)", p.MemSpaceTop[0].Func, p.MemSpaceTop[0].Flat, p.MemSpaceTop[0].FlatPct)
				}
				fmt.Fprintf(&sb, "| `%s` | %s | %s |\n", p.Benchmark, cpuTop, memTop)
			}
			sb.WriteString("\n*(Full top-15 CPU & memory tables for all profiled variants are in `latest.json` and `latest.html`.)*\n\n")
		} else {
			for _, p := range r.Profiles {
				fmt.Fprintf(&sb, "### `%s`\n\n", p.Benchmark)
				if len(p.CPUTop) > 0 {
					sb.WriteString("**CPU Top Hotspots**\n\n")
					sb.WriteString("| Flat | Flat% | Cum | Cum% | Function |\n")
					sb.WriteString("|---:|---:|---:|---:|---|\n")
					for _, e := range p.CPUTop {
						fmt.Fprintf(&sb, "| %s | %.2f%% | %s | %.2f%% | `%s` |\n",
							e.Flat, e.FlatPct, e.Cum, e.CumPct, e.Func)
					}
					sb.WriteString("\n")
				}
				if len(p.MemSpaceTop) > 0 {
					sb.WriteString("**Memory Alloc Space Top**\n\n")
					sb.WriteString("| Flat | Flat% | Cum | Cum% | Function |\n")
					sb.WriteString("|---:|---:|---:|---:|---|\n")
					for _, e := range p.MemSpaceTop {
						fmt.Fprintf(&sb, "| %s | %.2f%% | %s | %.2f%% | `%s` |\n",
							e.Flat, e.FlatPct, e.Cum, e.CumPct, e.Func)
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String()
}
