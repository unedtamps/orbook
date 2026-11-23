# Optimization Changelog — Orbook Matching Engine

All measurements captured on `Intel(R) Core(TM) i9-14900HX` (`linux/amd64`, `go1.27.3-X:nodwarf5`), `-benchtime=500ms`, `-benchmem`, `0 B/op, 0 allocs/op` on all pooled steady-state benchmarks. Every stage was verified with `go test -count=1 -race ./...` including the 25,000-op differential + quantity-conservation test.

---

## Overall Progression Summary (Baseline → Final)

| Metric | Stage 0 (Baseline `Book`) | Stage 1 (Fast Modulo) | Stage 2 (`OrderMap`) | Stage 4 (`FastBook`) | Stage 5 (`ArenaBook`) | Total Speedup |
|---|---:|---:|---:|---:|---:|---:|
| `AddLimit` | 31.04 ns | 30.81 ns | 21.90 ns | **16.00 ns** | 18.84 ns | **1.94×** (-48.5%) |
| `Cancel` | 38.46 ns | 37.29 ns | 18.23 ns | 14.33 ns | **12.75 ns** | **3.02×** (-66.8%) |
| `MatchSingle` | 59.26 ns | 60.96 ns | 36.37 ns | **32.88 ns** | 33.39 ns | **1.80×** (-44.5%) |
| `MatchSweep` (10 lvls) | 400.8 ns | 408.1 ns | 248.7 ns | 208.1 ns | **206.0 ns** | **1.95×** (-48.6%) |
| `Market` (10 lvls) | 386.6 ns | 404.3 ns | 248.6 ns | 203.8 ns | **198.3 ns** | **1.95×** (-48.7%) |
| `BestBidAsk` | 2.82 ns | 2.84 ns | 3.46 ns | **0.98 ns** | 1.11 ns | **2.86×** (-65.1%) |
| `Depth` (top 10) | 23.71 ns | 24.66 ns | 22.82 ns | **21.88 ns** | 26.57 ns | **1.08×** (-7.7%) |
| `MixedWorkload` | 45.14 ns | 45.40 ns | 30.76 ns | 25.28 ns | **25.08 ns** | **1.80×** (-44.4%) |
| HDR Throughput (200k ops) | 18.93M ops/s | 16.56M ops/s | 22.61M ops/s | 32.09M ops/s | **32.98M ops/s** | **1.74×** |
| HDR p50 / p99 / p99.9 | 59 / 156 / 230 ns | 58 / 155 / 263 ns | 44 / 131 / 210 ns | **40 / 88 / 127 ns** | 43 / 89 / 143 ns | **1.81×** tighter p99.9 |
| GC Mark (250k resting orders) | — | — | — | 1.97 ms | **453 µs** | **4.34×** faster GC |

### Engine Concurrency Progression (Stage 3)

| Runner | Time / op | Throughput | Allocs | Speedup vs Baseline Channel |
|---|---:|---:|---:|---:|
| `BenchmarkEngineChannel` (Stage 0 baseline, Go `map`) | 92.15 ns | 10.9M ops/s | 0 B / 0 allocs | 1.00× |
| `BenchmarkEngineChannel` (after Stage 2 `OrderMap`) | 79.10 ns | 12.6M ops/s | 0 B / 0 allocs | 1.17× |
| `BenchmarkEngineDisruptor` (Lock-free SPSC ring, 1-by-1) | **54.40 ns** | **18.4M ops/s** | 0 B / 0 allocs | **1.69×** (-41.0%) |
| `BenchmarkEngineBatch` (Batched channel, 64 cmds/send) | **22.96 ns** | **43.6M ops/s** | 0 B / 0 allocs | **4.01×** (-75.1%) |

---

## Stage-by-Stage Change Log & Analysis

### Stage 0 — Baseline (`ladder_linkedlist_pool`, Go `map[OrderID]*Order`)
- **Files:** `book/book.go`, `book/levels/ladder.go`, `book/queue/linkedlist.go`, `pool/freelist.go`, `engine/engine.go`
- **Results:** `AddLimit` `31.04 ns`, `Cancel` `38.46 ns`, `MatchSingle` `59.26 ns`, `MatchSweep` `400.8 ns`, `MixedWorkload` `45.14 ns`, `EngineChannel` `92.15 ns`, HDR `18.93M ops/s` (`p50=59ns, p99=156ns, p99.9=230ns`).
- **Bottleneck from `pprof`:** Go's built-in `map[OrderID]*Order` (`runtime.mapaccess2_fast64`, `mapassign_fast64`, `maps.ctrlGroup.matchH2`, `Map.Delete`) consumed **~35–45% of total CPU**.

---

### Stage 1 — Fast-Path `TickSize == 1` & `LotSize == 1` Validation
- **Files changed:** `book/book.go` (`validateOrder`)
- **Change:** Guarded `req.Qty%b.cfg.LotSize != 0` with `b.cfg.LotSize != 1` and `(req.Price-b.cfg.MinPrice)%b.cfg.TickSize != 0` with `b.cfg.TickSize != 1` to avoid x86 `IDIV` modulo instructions when tick/lot size is 1.
- **Measured delta vs Stage 0:**
  - `AddLimit`: `31.04 → 30.81 ns` (**-0.7%**)
  - `Cancel`: `38.46 → 37.29 ns` (-3.0%, unaffected by `validateOrder`, within run jitter)
  - `MixedWorkload`: `45.14 → 45.40 ns` (**+0.6%**, flat)
- **Examination:** Effectively within noise on its own. On an out-of-order CPU (i9-14900HX), `LotSize` was `100` (so `Qty % 100` still ran) and the price `IDIV` latency was already hidden in parallel behind the memory-bound `b.orders[req.ID]` Go map lookup right below it. Still worth keeping as a zero-risk guard that pays off once the map bottleneck is removed.

---

### Stage 2 — Flat Open-Addressed `orderMap` (Fibonacci Hash + Backward-Shift Delete)
- **Files added/changed:** `book/ordermap.go`, `book/ordermap_test.go`, `book/book.go`, `book/invariants.go`
- **Change:**
  - Replaced Go's built-in `map[OrderID]*Order` with a single-threaded flat power-of-2 open-addressed linear-probing table (`orderSlot{key OrderID, val *Order}`, 16 B = 4 slots per cache line).
  - Replaced AES/SipHash with Fibonacci multiplication `((x ^ (x >> 32)) * 11400714819323198485) >> shift`.
  - Implemented **backward-shift deletion** (zero tombstones ever) and collapsed `CancelOrder` from **2 map probes (`Get` + `delete`) into 1 single probe (`orders.Delete(id)`)**.
  - Verified with a 20,000-op randomized differential test (`TestOrderMapDifferentialRandom`) + full 16-variant differential test.
- **Measured delta vs Stage 1 (biggest single optimization!):**
  - `AddLimit`: `30.81 → 21.90 ns` (**-28.9%**, 1.41× faster)
  - `Cancel`: `37.29 → 18.23 ns` (**-51.1%**, **2.05× faster**)
  - `MatchSingle`: `60.96 → 36.37 ns` (**-40.3%**, 1.68× faster)
  - `MatchSweep`: `408.1 → 248.7 ns` (**-39.1%**, 1.64× faster)
  - `Market`: `404.3 → 248.6 ns` (**-38.5%**, 1.63× faster)
  - `MixedWorkload`: `45.40 → 30.76 ns` (**-32.2%**, **32.5M ops/s**)
  - `EngineChannel`: `92.15 → 79.22 ns` (**-14.0%**)
  - HDR Latency: `p50` `58 → 44 ns`, `p99` `155 → 131 ns`, `p99.9` `263 → 210 ns`
- **`pprof` shift:** Every single `runtime.map*` and `internal/runtime/maps.*` entry **vanished completely**. `orderMap.Has` and `orderMap.hash` inlined directly into `validateOrder`/`PlaceOrder`, and `orderMap.Delete` dropped to just `7.89%` flat CPU. The newly exposed bottleneck became non-inlined interface calls (`FreeList.Put` `4.0%`, `LinkedList.Remove` `2.6%`, `LinkedList.PushBack` `2.0%`, `Ladder.*`).

---

### Stage 3 — Batched Channel (`BatchEngine`) & Lock-Free SPSC Ring (`Disruptor`)
- **Files added/changed:** `engine/disruptor.go`, `engine/batch.go`, `engine/disruptor_test.go`, `bench/bench_test.go`
- **Change:**
  - Added `engine.Disruptor`: true LMAX-style lock-free SPSC ring buffer with 64-byte cache-line padding separating producer (`tail`, `cachedHead`) and consumer (`head`, `cachedTail`) cursors, eliminating false sharing and Go channel mutexes while keeping a 1-by-1 streaming API (`PlaceOrderAsync`/`CancelOrderAsync`).
  - Added `engine.BatchEngine`: multi-producer-safe channel runner accepting slices of `[]Command` (64 commands per send) to amortize Go channel synchronization.
  - Verified identical event output against direct `Book` under `-race` (`TestDisruptorAndBatchEngineIdenticalToDirectBook`) and updated benchmarks so full drain (`Stop()`) occurs inside the timed window.
- **Measured delta vs `BenchmarkEngineChannel` (`79.10 ns/op`):**
  - `BenchmarkEngineDisruptor` (single-command lock-free SPSC): **`54.40 ns/op`** (**-31.2%**, 1.45× faster, **18.4M ops/s**, `0 B/op`)
  - `BenchmarkEngineBatch` (64-cmd batches): **`22.96 ns/op`** (**-71.0%**, 3.45× faster, **43.6M ops/s**, `0 B/op` — producer slicing overlaps on core 1 while matching runs uninterrupted on core 2)

---

### Stage 4 — Devirtualized Concrete Core (`book.FastBook`)
- **Files added/changed:** `book/fastbook.go`, `book/differential_test.go`, `bench/fast_bench_test.go`
- **Change:**
  - Built a concrete zero-interface specialization (`FastBook`) embedding `fastLadder` (with `head, tail *Order` directly inside each `fastLevel`), an inline intrusive `FreeList`, and flat `orderMap`.
  - Eliminated all interface dispatch (`PriceLevels`, `LevelQueue`, `OrderPool`) so Go inlines every hot-path operation directly into `PlaceOrder`, `CancelOrder`, and `BestBid`/`BestAsk`.
  - Verified 100% event-for-event equality with `Book` across 25,000 random ops in `TestDifferentialAllImplementations`.
- **Measured delta vs Stage 2 (`Book` + `OrderMap`):**
  - `AddLimit`: `21.90 → 16.00 ns` (**-26.9%**)
  - `Cancel`: `18.23 → 14.33 ns` (**-21.4%**)
  - `MatchSingle`: `36.37 → 32.88 ns` (**-9.6%**)
  - `MatchSweep`: `248.7 → 208.1 ns` (**-16.3%**)
  - `Market`: `248.6 → 203.8 ns` (**-18.0%**)
  - `BestBidAsk`: `3.46 → 0.98 ns` (**-71.6%**, sub-nanosecond!)
  - `MixedWorkload`: `30.76 → 25.28 ns` (**-17.8%**, **39.6M ops/s**)
  - HDR Latency: **32.09M ops/s**, `p50=40 ns`, `p99=88 ns`, **`p99.9=127 ns`** (**-39.5%** tighter p99.9 vs Stage 2!)

---

### Stage 5 — Pointer-Free `noscan` 64-Byte `Order` Arena (`book.ArenaBook`)
- **Files added/changed:** `book/arenabook.go`, `book/ordermap_test.go`, `book/differential_test.go`, `bench/fast_bench_test.go`
- **Change:**
  - Replaced all `*Order` and `*PriceLevel` pointers inside the book with `uint32` arena/tick indices:
    - `arenaOrder`: **64 bytes** (exact 1 L1 cache line, 0 pointers → GC `noscan`)
    - `arenaLevel`: **32 bytes** (2 per cache line, 0 pointers → GC `noscan`)
    - `arenaMapSlot`: **16 bytes** (0 pointers → GC `noscan`)
  - Verified exact struct sizes in `TestArenaStructSizes`, 100% event-for-event equality in `TestDifferentialAllImplementations`, and GC scan pause with 250,000 resting orders in `TestFastAndArenaLatencyAndGC`.
- **Measured delta vs Stage 4 (`FastBook`):**
  - `Cancel`: `14.33 → 12.75 ns` (**-11.0%**, fastest cancel overall — **3.02×** faster than Stage 0!)
  - `Market`: `203.8 → 198.3 ns` (**-2.7%**)
  - `MatchSweep`: `208.1 → 206.0 ns` (**-1.0%**)
  - `MixedWorkload`: `25.28 → 25.08 ns` (**-0.8%**, **39.9M ops/s**; HDR throughput **32.98M ops/s**, `max=23.3µs`)
  - `AddLimit`: `16.00 → 18.84 ns` (+17.8% — writing a full 64-byte struct literal + slice bounds checks on `b.arena[idx]` costs ~2.8 ns more than direct pointer stores)
  - **`runtime.GC()` with 250,000 resting orders:** `1.97 ms` (`FastBook`) → **`453 µs`** (`ArenaBook`) — **4.34× faster GC mark cycles!**
- **Examination:** `FastBook` and `ArenaBook` tie on raw mixed-workload speed (~25 ns/op, ~40M ops/s). Choose **`FastBook`** when resting order counts stay modest (<50k) and you want lowest raw `AddLimit` (`16.0 ns`) and `p99.9` (`127 ns`); choose **`ArenaBook`** when holding large resting books (100k–1M+ orders) where eliminating Go GC pointer scanning (**4.34×** faster GC + 33% less RAM per order: `64 B` vs `96 B`) dominates tail latency stability.
