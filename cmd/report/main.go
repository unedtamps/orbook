package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/unedtamps/orbook/report"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "serve" {
		runServe(os.Args[2:])
		return
	}
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "run" {
		args = args[1:]
	}
	runReport(args)
}

func runReport(args []string) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	mode := fs.String("mode", "full", "report mode: full or single")
	benchFilter := fs.String("bench", ".", "benchmark regex filter (e.g. . or ^BenchmarkAddLimit$ or ^BenchmarkAddLimit$/^ladder_linkedlist_pool$)")
	benchTime := fs.String("benchtime", "150ms", "go test -benchtime value for timing pass")
	profileMode := fs.String("profile-mode", "all", "profiling scope in full mode: all, methods, single, none")
	profileBench := fs.String("profile-bench", "", "benchmark name to profile when profile-mode=single")
	profileTime := fs.String("profile-time", "", "benchtime per profiled benchmark (default 250ms for all, 500ms otherwise)")
	includeHDR := fs.Bool("hdr", true, "include HDR histogram latency suite (all variants in full mode, matched variants in single mode)")
	hdrOps := fs.Int("ops", 200_000, "ops for HDR latency suite")
	threshold := fs.Float64("threshold", 2.0, "minimum % ns/op delta to count as a diff and create a new version")
	force := fs.Bool("force", false, "force creating a new version even if no diff")
	outDir := fs.String("out", "reports", "output directory for centralized reports")
	_ = fs.Parse(args)

	// Normalize simple "Method/Variant" filter into exact Go sub-benchmark regex if no regex metachars used.
	filter := *benchFilter
	if *mode == "single" && filter != "." && !strings.ContainsAny(filter, "^$*+?()[]|\\") {
		filter = report.ExactBenchRegex(filter)
	}

	opts := report.RunOptions{
		OutDir:       *outDir,
		Mode:         *mode,
		BenchFilter:  filter,
		BenchTime:    *benchTime,
		ProfileMode:  *profileMode,
		ProfileBench: *profileBench,
		ProfileTime:  *profileTime,
		IncludeHDR:   *includeHDR,
		HDROps:       *hdrOps,
		ThresholdPct: *threshold,
		Force:        *force,
		Progress: func(msg string) {
			fmt.Println("  " + msg)
		},
	}

	fmt.Printf("Running report (mode=%s, bench=%q, benchtime=%s, profile=%s, threshold=±%.1f%%)...\n",
		opts.Mode, opts.BenchFilter, opts.BenchTime, opts.ProfileMode, opts.ThresholdPct)

	res, err := report.Run(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "report failed: %v\n", err)
		os.Exit(1)
	}

	if res.Saved {
		fmt.Printf("Created new version: %s\n", res.VersionDir)
		fmt.Printf("Updated centralized outputs:\n  - %s/latest.json (agent)\n  - %s/latest.md   (agent / PR)\n  - %s/latest.html (human dashboard)\n",
			*outDir, *outDir, *outDir)
	} else {
		fmt.Printf("No significant diff vs %s (within ±%.1f%% threshold) — skipped creating new version.\n",
			res.Report.RunID, opts.ThresholdPct)
	}

	if res.Report.Diff != nil && len(res.Report.Diff.Items) > 0 {
		for _, d := range res.Report.Diff.Items {
			if d.Status == "unchanged" {
				continue
			}
			if d.Status == "added" {
				fmt.Printf("  [added]  %s: %.2f ns/op (%d B/op, %d allocs/op)\n",
					d.FullName, d.NewNs, d.NewBytes, d.NewAllocs)
			} else {
				fmt.Printf("  [%s] %s: %.2f -> %.2f ns/op (%+.2f%%) | %d->%d B/op, %d->%d allocs/op\n",
					d.Status, d.FullName, d.OldNs, d.NewNs, d.DeltaNsPct, d.OldBytes, d.NewBytes, d.OldAllocs, d.NewAllocs)
			}
		}
	}
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", ":8080", "listen address (e.g. :8080)")
	dir := fs.String("dir", "reports", "reports directory to serve")
	_ = fs.Parse(args)

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "serve failed: %v\n", err)
		os.Exit(1)
	}

	fsHandler := http.FileServer(http.Dir(absDir))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/latest.html", http.StatusFound)
			return
		}
		fsHandler.ServeHTTP(w, r)
	})

	fmt.Printf("Serving reports from %s\n", absDir)
	fmt.Printf("  Dashboard: http://localhost%s/latest.html\n", *addr)
	fmt.Printf("  JSON:      http://localhost%s/latest.json\n", *addr)
	fmt.Printf("  History:   http://localhost%s/history/\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
