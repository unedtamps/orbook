# Orbook — Benchmark & Profiling Analysis

> **Source of truth:** [reports/latest.json](reports/latest.json) (`v007_20261123_031648`)
> **Environment:** `Intel(R) Core(TM) i9-14900HX` · `linux/amd64` · `go1.27.3-X:nodwarf5`
> **Config:** 1 ticker (`ACME`), `TickSize=1`, `LotSize=100`, price band `[10000..30000]` (20,001 ticks)
> **Mixed Workload:** Seed `20261122`, 2,000 warmup resting orders + 100k/200k ops (60% Limit / 30% Cancel / 10% Market across 16 traders with random-walk mid & STP)

---

## Executive Summary & Recommended Production Stack

| Scenario | Recommended Stack | Mixed Workload | Throughput | HDR p50 / p99 / p99.9 | Steady-State Allocs |
|---|---|---:|---:|---:|---:|
| **Equities with daily price band** *(IDX / NYSE / NASDAQ)* | **`Ladder` + `LinkedList` + `FreeList`** | **44.4 ns/op** | **19.9M ops/s** | **57 / 148 / 235 ns** | **0 B/op, 0 allocs/op** |
| **Unbounded price range / deep book** *(crypto / wide ticks)* | **`RBTree` + `LinkedList` + `FreeList`** | **62.5 ns/op** | **13.8M ops/s** | **76 / 202 / 305 ns** | **0 B/op, 0 allocs/op** |
| **Unbounded price range / shallow book** *(<200 active levels)* | **`SortedSlice` + `LinkedList` + `FreeList`** | **54.7 ns/op** | **15.7M ops/s** | **67 / 156 / 244 ns** | **0 B/op, 0 allocs/op** |
| **With LMAX Channel Runner** (`engine.Engine`) | **`Ladder` + `LinkedList` + `FreeList`** | **90.6 ns/op** | **11.0M ops/s** | — | **0 B/op, 0 allocs/op** |

**Key conclusion:** **`LinkedList` + `FreeList` pool** is unambiguously the best queue & memory combination across every price-level store — delivering **0 allocs/op**, **3–5× tighter p99.9 tail latency** than `RingBuffer`, and **up to 7.4× tighter p99.9 tail latency** than `NoPool`.

---

## 1. Price-Level Store Comparison (`LinkedList` + `FreeList`, all `0 B/op, 0 allocs/op`)

| Method | `Ladder` | `SortedSlice` | `RBTree` *(hand-rolled)* | `BTree` *(`google/btree`)* | Winner |
|---|---:|---:|---:|---:|---|
| `AddLimit` (resting, 200 lvls) | **29.86 ns** | 59.29 ns | 58.47 ns | 79.39 ns | `Ladder` (**2.0×** vs Slice/RB, **2.7×** vs BTree) |
| `Cancel` (random resting order) | 39.23 ns | **39.02 ns** | 40.36 ns | 39.86 ns | Tie (~39–40 ns; map-bound) |
| `MatchSingle` (1 level FIFO) | 60.25 ns | **56.22 ns** | 63.34 ns | 60.43 ns | Tie (~56–63 ns; queue/map-bound) |
| `MatchSweep` (10 levels) | 368.10 ns | **365.60 ns** | 503.30 ns | 795.60 ns | `SortedSlice` / `Ladder` (**2.2×** vs BTree) |
| `Market` (10 levels) | 362.30 ns | **356.40 ns** | 512.80 ns | 781.30 ns | `SortedSlice` / `Ladder` (**2.2×** vs BTree) |
| `BestBidAsk` (top of book) | 2.31 ns | 3.13 ns | 4.15 ns | **2.17 ns** | `BTree` / `Ladder` (O(1) cached) |
| `Depth` (top 10 bids + asks) | 22.27 ns | **16.02 ns** | 39.26 ns | 60.11 ns | `SortedSlice` (**3.8×** vs BTree) |
| `MixedWorkload` (60/30/10) | **44.44 ns** | 54.65 ns | 62.53 ns | 68.16 ns | `Ladder` (**1.23×** vs Slice, **1.53×** vs BTree) |

### Why Each Store Behaved This Way (Backed by `pprof`)

1. **`Ladder` (`book/levels/ladder.go`) — Best overall for bounded equities**
   - O(1) tick arithmetic (`idx = price - minPrice`) eliminates search entirely: in `BenchmarkAddLimit`, `Ladder.GetOrCreate` doesn't even show up in the top CPU hotspots, leaving `Book.PlaceOrder` (`14.8%` flat) and Go `map` ops as the dominant work.
   - The `[]uint64` active bitset with `bits.LeadingZeros64` / `bits.TrailingZeros64` keeps `Ladder.Remove` to only **`5.7%` flat CPU** during 10-level sweeps (`368 ns`) and `Depth(10)` to **`22.3 ns`**.
   - **Trade-off:** Fixed upfront memory (`~1.9 MB` for 20,001 ticks × 2 sides) and requires a bounded daily `[MinPrice, MaxPrice]` band.

2. **`SortedSlice` (`book/levels/sortedslice.go`) — Best for shallow books & `Depth(10)`**
   - Storing levels **worst-to-best** (best level at tail `len-1`) makes `Best()`, popping the best level during `MatchSweep`/`Market` (`356–366 ns`), and BBO-improving inserts O(1) with zero shifting.
   - `Depth(10)` simply reads 10 adjacent pointers backwards (`16.0 ns`, `88.6%` flat in `SortedSlice.TopN`).
   - **CPU profile catch:** `SortedSlice.search` takes **`15.7%` flat** on `AddLimit` (`59.3 ns`), and inserting deep levels across a 5,000-level book during `MatchSweep` setup triggers **`46.7%` flat in `runtime.memmove`** — confirming O(N) shift degradation on deep books.

3. **`RBTree` (`book/levels/rbtree.go`) — Beats `google/btree` across the board**
   - Outperforms `google/btree` on `AddLimit` (**1.36×**), `MatchSweep` (**1.58×**), `Market` (**1.52×**), `Depth(10)` (**1.53×**), `MixedWorkload` (**1.09×**), and HDR `p99.9` (**305 ns vs 348 ns**).
   - **Why (from `pprof`):**
     - Direct inline `int64` price comparisons instead of `google/btree`'s generic `less` function-pointer closure (`btree.items.find` + `sort.Search` consumed **`28.8%` cum CPU** on `AddLimit` and **`22.4%` cum** on `MixedWorkload`).
     - `PriceLevel` embedded directly inside `rbNode` (one less pointer dereference).
     - O(1) `nextBest = successor(z)` when deleting the leftmost node during sweeps (vs `google/btree` node shifting + `tree.Min()`).
     - Iterative parent-pointer `TopN` (`39.3 ns` vs `60.1 ns` for `BTree.Ascend` closure traversal).

4. **`BTree` (`book/levels/btree.go`, `google/btree` v1.1.3)**
   - Fastest `BestBidAsk` (**2.17 ns**) due to returning the cached `b.best` pointer directly.
   - Slowest on `MatchSweep` (`795.6 ns`) and `Depth(10)` (`60.1 ns`, where `BTree.TopN.func1` + `node.iterate` take **`87.5%` cum CPU**).

---

## 2. Queue Comparison: `LinkedList` vs `RingBuffer`

| Variant (all `FreeList` pooled) | `MixedWorkload` (ns/op) | HDR Throughput (ops/s) | p50 (ns) | p90 (ns) | p99 (ns) | **p99.9 (ns)** | Max (ns) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `ladder_linkedlist_pool` | 44.44 | **19,936,538** | **57** | **87** | **148** | **235** | 135,679 |
| `ladder_ringbuffer_pool` | **43.61** | 18,855,741 | 64 | 103 | 199 | **866** | 34,303 |
| `sortedslice_linkedlist_pool` | 54.65 | 15,748,969 | **67** | **102** | **156** | **244** | 84,223 |
| `sortedslice_ringbuffer_pool` | **54.08** | 16,527,072 | 68 | 109 | 190 | **1,191** | 178,815 |
| `rbtree_linkedlist_pool` | **62.53** | **13,754,218** | **76** | **124** | **202** | **305** | 28,863 |
| `rbtree_ringbuffer_pool` | 64.79 | 11,564,123 | 81 | 139 | 248 | **902** | 31,567 |
| `btree_linkedlist_pool` | 68.16 | **13,891,074** | **75** | **123** | **200** | **348** | 109,311 |
| `btree_ringbuffer_pool` | **67.02** | 13,597,491 | 79 | 130 | 224 | **1,270** | 38,143 |

### Analysis
- **Throughput:** Virtually identical (~1–2% difference in `BenchmarkMixedWorkload`).
- **Tail latency (decisive):** `LinkedList` beats `RingBuffer` by **3.0× to 4.9× on p99.9** across every single price-level store (`235–348 ns` vs `866–1,270 ns`).
- **Why:** Even though `RingBuffer` clears tombstones in O(1) and compacts in place with 0 allocs, periodic compaction scans and 2× buffer resizes on busy ticks create sub-microsecond latency spikes. Intrusive `LinkedList` executes a fixed 4-pointer unlink in strict worst-case O(1) and holds **zero spare slice memory per level**.

---

## 3. Memory Pooling Impact: `FreeList` vs `NoPool`

| Variant | `AddLimit` | `MixedWorkload` | Steady-State Allocs | HDR Throughput | p99.9 Latency | 500ms `memprofile` `NoPool.Get` |
|---|---:|---:|---:|---:|---:|---:|
| `ladder_linkedlist_pool` | **29.86 ns** | **44.44 ns** | **0 B / 0 allocs** | **19.94M ops/s** | **235 ns** | 0 B |
| `ladder_linkedlist_nopool` | 59.79 ns *(2.0× slower)* | 58.20 ns *(1.31× slower)* | 96 B (Add) / 51 B (Mixed) | 14.39M ops/s | **1,733 ns** *(7.4× worse)* | **398.54 MB (93.7%)** |
| `sortedslice_linkedlist_pool` | **59.29 ns** | **54.65 ns** | **0 B / 0 allocs** | **15.75M ops/s** | **244 ns** | 0 B |
| `sortedslice_linkedlist_nopool` | 95.32 ns *(1.6× slower)* | 67.70 ns *(1.24× slower)* | 96 B (Add) / 51 B (Mixed) | 10.06M ops/s | **493 ns** *(2.0× worse)* | >400 MB |
| `rbtree_linkedlist_pool` | **58.47 ns** | **62.53 ns** | **0 B / 0 allocs** | **13.75M ops/s** | **305 ns** | 0 B |
| `rbtree_linkedlist_nopool` | 92.91 ns *(1.6× slower)* | 75.51 ns *(1.21× slower)* | 96 B (Add) / 51 B (Mixed) | 11.18M ops/s | **519 ns** *(1.7× worse)* | >650 MB |

### Analysis
- Eliminating per-order heap allocation (`96 B/op`) with a single shared slab-backed `FreeList` (`384 KB` per 4,096-order slab) cuts `AddLimit` time in **half** on `Ladder`, raises mixed-workload throughput by **+23% to +56%**, and eliminates GC-driven tail spikes (**7.4× lower p99.9** on `ladder_linkedlist`).
- All memory stays bounded by **peak concurrent resting orders** (~2–5 MB total for a ticker with 10k–30k resting orders), rather than total daily volume.
- Profiling also caught and prevented a real bug during development: `Book.Reset()` must reset intrusive level queues (`LinkedList`) *before* calling `FreeList.Put(o)`, since both share `Order.Next`.

---

## 4. CPU & Memory `pprof` Hotspots & Next Optimization Opportunities

### Hotspot #1: Go `map[OrderID]*Order` (~35–45% of total CPU in `ladder_linkedlist_pool`)
In `BenchmarkMixedWorkload/ladder_linkedlist_pool` (`44.4 ns/op`) and `BenchmarkCancel/ladder_linkedlist_pool` (`39.2 ns/op`), the price-level store and linked list are now so fast that Go's built-in SwissTable `map[OrderID]*Order` (`runtime.mapaccess2_fast64`, `runtime.mapassign_fast64`, `maps.ctrlGroup.matchH2`, `maps.(*Map).Delete`) is the single largest CPU consumer.
- This is also why `Cancel` takes **~39–40 ns identically across all 4 price-level stores**: `Order.Level` skips the price-level lookup on non-empty levels, so cancel is map-bound.
- **How to push further:** Replace `map[OrderID]*Order` with a flat open-addressed linear-probing hash table (power-of-2 slice, Fibonacci hash, zero runtime calls) or a direct ring/array index if order IDs are dense sequential integers from the gateway. Expected gain: **~15–20 ns/op** on `AddLimit`, `Cancel`, and `MixedWorkload` (pushing `Ladder` toward **~30–35M ops/s**).

### Hotspot #2: Channel Lock Contention in LMAX Runner (`BenchmarkEngineChannel`: `90.6 ns/op`)
CPU profile of `BenchmarkEngineChannel` (`11.0M ops/s`, `0 B/op, 0 allocs/op`):
- `runtime.lock2`: **34.0% cum**
- `runtime.procyieldAsm`: **23.0% flat**
- `runtime.unlock2`: **14.0% cum**
- `book.(*Book).PlaceOrder`: **23.0% cum**
- **How to push further:** Batch commands (e.g., send `[]Command` of 16–64 commands per channel op) or replace `chan Command` with an SPSC lock-free ring buffer (LMAX Disruptor pattern) to cut cross-goroutine transport overhead from ~46 ns to ~5–10 ns.

### Hotspot #3: GC Pointer Scanning on Large Resting Books
- `book.Order` currently holds 3 pointers (`Level *PriceLevel`, `Prev *Order`, `Next *Order`), so Go's GC still scans `FreeList` slabs during mark phases even with 0 allocations.
- **How to push further:** Store `uint32` arena indices instead of `*Order` pointers inside `Order` slabs so the order arena becomes pointer-free (`noscan`), reducing GC mark time to near zero even at 1M+ resting orders.

---

## 5. Reproducing & Viewing Reports

```bash
# Run all tests (unit, -race, zero-alloc, 16-variant differential + conservation, RBTree invariants)
make test

# Full matrix benchmark + CPU/Mem pprof for all 129 benchmarks + HDR latency
make report         # high precision (500ms bench / 350ms pprof)
make report-quick   # faster (100ms bench / 150ms pprof, still profiles all 129)

# Incremental run for a single variant or implementation (merges & versions only on diff)
make report-one BENCH='.*/rbtree_.*'
make report-one BENCH=BenchmarkAddLimit/ladder_linkedlist_pool

# Launch interactive side-by-side HTML dashboard at http://localhost:8080/latest.html
make report-serve
```
