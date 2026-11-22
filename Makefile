.PHONY: test test-short bench bench-quick bench-cmp profile-cpu profile-mem latency latency-detailed demo fuzz report report-quick report-one report-serve clean

PROFILE_DIR  := profiles
BENCH        ?= BenchmarkMixedWorkload/ladder_linkedlist_pool
BENCHTIME    ?= 1s
PROFILE_MODE ?= all
PROFILE_TIME ?= 350ms
COUNT        ?= 1
THRESHOLD    ?= 2.0
PORT         ?= 8080

test:
	go test -v -race ./...

test-short:
	go test -race ./...

bench:
	go test -run=^$$ -bench=. -benchmem -benchtime=$(BENCHTIME) -count=$(COUNT) ./bench/...

bench-quick:
	go test -run=^$$ -bench=. -benchmem -benchtime=150ms -count=1 ./bench/...

bench-cmp:
	@mkdir -p $(PROFILE_DIR)
	@echo "Running mixed-workload benchmarks (5x count) -> $(PROFILE_DIR)/bench_mixed.txt"
	go test -run=^$$ -bench=BenchmarkMixedWorkload -benchmem -benchtime=500ms -count=5 ./bench/... | tee $(PROFILE_DIR)/bench_mixed.txt
	@if command -v benchstat >/dev/null 2>&1; then \
		echo "benchstat available: compare two runs with 'benchstat old.txt new.txt'"; \
	else \
		echo "Install benchstat for statistical diffs: go install golang.org/x/perf/cmd/benchstat@latest"; \
	fi

profile-cpu:
	@mkdir -p $(PROFILE_DIR)
	go test -run=^$$ -bench='^$(BENCH)$$' -benchmem -benchtime=2s -cpuprofile=$(PROFILE_DIR)/cpu.prof ./bench/...
	go tool pprof -top -nodecount=25 $(PROFILE_DIR)/cpu.prof

profile-mem:
	@mkdir -p $(PROFILE_DIR)
	go test -run=^$$ -bench='^$(BENCH)$$' -benchmem -benchtime=2s -memprofile=$(PROFILE_DIR)/mem.prof ./bench/...
	@echo "=== alloc_space ==="
	go tool pprof -top -alloc_space -nodecount=20 $(PROFILE_DIR)/mem.prof
	@echo "=== alloc_objects ==="
	go tool pprof -top -alloc_objects -nodecount=20 $(PROFILE_DIR)/mem.prof

latency:
	go run ./cmd/orbook latency -ops=200000

latency-detailed:
	go run ./cmd/orbook latency -ops=200000 -detailed

demo:
	go run ./cmd/orbook demo

fuzz:
	go test -fuzz=FuzzOrderBookDifferential -fuzztime=15s ./book

# Centralized reporting (JSON + Markdown for agents, HTML dashboard for humans)
# Full run: all benchmarks + CPU/Mem pprof for ALL variants + HDR latency (higher precision)
report:
	go run ./cmd/report -mode=full -bench=. -benchtime=500ms -profile-mode=$(PROFILE_MODE) -profile-time=$(PROFILE_TIME) -ops=200000 -threshold=$(THRESHOLD)

# Faster full run: still runs & profiles ALL variants, just with shorter benchtime/ops
report-quick:
	go run ./cmd/report -mode=full -bench=. -benchtime=100ms -profile-mode=$(PROFILE_MODE) -profile-time=150ms -ops=100000 -threshold=$(THRESHOLD)

# Single benchmark (e.g. BENCH=BenchmarkAddLimit/ladder_linkedlist_pool or BENCH=BenchmarkAddLimit)
# Runs timing + CPU/Mem pprof for every matched variant; creates a new version only if diff >= THRESHOLD
report-one:
	go run ./cmd/report -mode=single -bench='$(BENCH)' -benchtime=500ms -profile-time=500ms -threshold=$(THRESHOLD)

# Serve the centralized HTML dashboard & version history locally
report-serve:
	go run ./cmd/report serve -addr=:$(PORT) -dir=reports

clean:
	rm -rf $(PROFILE_DIR) *.test
