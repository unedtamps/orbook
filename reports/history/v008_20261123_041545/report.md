# Orbook Benchmark & Profile Report — `v008_20261123_041545`

- **Version:** `v008` | **Mode:** `single` | **Filter:** `^(BenchmarkFastBook|BenchmarkArenaBook|BenchmarkEngine)`
- **Time:** `2026-11-23 04:15:45 WIB` | **Commit:** `4e16c62` | **Go:** `go1.27.3-X:nodwarf5` | **CPU:** `Intel(R) Core(TM) i9-14900HX`

## Diff vs Previous (`v007_20261123_031648`, threshold ±2.0%)

| Benchmark | Status | Old (ns/op) | New (ns/op) | Delta | Allocs |
|---|---|---:|---:|---:|---|
| `BenchmarkEngineChannel` | `faster` | 90.57 | 78.91 | -12.87% | 0→0 B, 0→0 allocs |
| `BenchmarkEngineBatch` | `added` | — | 22.40 | — | 0 B, 0 allocs |
| `BenchmarkEngineDisruptor` | `added` | — | 51.98 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/AddLimit` | `added` | — | 15.79 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/Cancel` | `added` | — | 15.15 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/MatchSingle` | `added` | — | 34.03 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/MatchSweep` | `added` | — | 213.10 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/Market` | `added` | — | 211.50 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/BestBidAsk` | `added` | — | 0.94 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/Depth` | `added` | — | 22.87 | — | 0 B, 0 allocs |
| `BenchmarkFastBook/MixedWorkload` | `added` | — | 24.57 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/AddLimit` | `added` | — | 19.78 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/Cancel` | `added` | — | 13.82 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/MatchSingle` | `added` | — | 35.12 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/MatchSweep` | `added` | — | 212.80 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/Market` | `added` | — | 203.10 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/BestBidAsk` | `added` | — | 1.16 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/Depth` | `added` | — | 27.54 | — | 0 B, 0 allocs |
| `BenchmarkArenaBook/MixedWorkload` | `added` | — | 27.24 | — | 0 B, 0 allocs |

## Benchmarks

| Method | Variant | ns/op | B/op | allocs/op |
|---|---|---:|---:|---:|
| `BenchmarkAddLimit` | `sortedslice_linkedlist_pool` | 59.29 | 0 | 0 |
| `BenchmarkAddLimit` | `sortedslice_linkedlist_nopool` | 95.32 | 96 | 1 |
| `BenchmarkAddLimit` | `sortedslice_ringbuffer_pool` | 59.13 | 0 | 0 |
| `BenchmarkAddLimit` | `sortedslice_ringbuffer_nopool` | 104.20 | 96 | 1 |
| `BenchmarkAddLimit` | `btree_linkedlist_pool` | 79.39 | 0 | 0 |
| `BenchmarkAddLimit` | `btree_linkedlist_nopool` | 114.70 | 96 | 1 |
| `BenchmarkAddLimit` | `btree_ringbuffer_pool` | 82.56 | 0 | 0 |
| `BenchmarkAddLimit` | `btree_ringbuffer_nopool` | 116.10 | 96 | 1 |
| `BenchmarkAddLimit` | `ladder_linkedlist_pool` | 29.86 | 0 | 0 |
| `BenchmarkAddLimit` | `ladder_linkedlist_nopool` | 59.79 | 96 | 1 |
| `BenchmarkAddLimit` | `ladder_ringbuffer_pool` | 33.22 | 0 | 0 |
| `BenchmarkAddLimit` | `ladder_ringbuffer_nopool` | 60.21 | 96 | 1 |
| `BenchmarkCancel` | `sortedslice_linkedlist_pool` | 39.02 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_linkedlist_nopool` | 42.14 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_ringbuffer_pool` | 38.01 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_ringbuffer_nopool` | 49.72 | 0 | 0 |
| `BenchmarkCancel` | `btree_linkedlist_pool` | 39.86 | 0 | 0 |
| `BenchmarkCancel` | `btree_linkedlist_nopool` | 42.70 | 0 | 0 |
| `BenchmarkCancel` | `btree_ringbuffer_pool` | 39.67 | 0 | 0 |
| `BenchmarkCancel` | `btree_ringbuffer_nopool` | 54.27 | 0 | 0 |
| `BenchmarkCancel` | `ladder_linkedlist_pool` | 39.23 | 0 | 0 |
| `BenchmarkCancel` | `ladder_linkedlist_nopool` | 44.52 | 0 | 0 |
| `BenchmarkCancel` | `ladder_ringbuffer_pool` | 37.73 | 0 | 0 |
| `BenchmarkCancel` | `ladder_ringbuffer_nopool` | 45.35 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_linkedlist_pool` | 56.22 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_linkedlist_nopool` | 64.13 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_ringbuffer_pool` | 60.57 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_ringbuffer_nopool` | 65.04 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_linkedlist_pool` | 60.43 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_linkedlist_nopool` | 64.05 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_ringbuffer_pool` | 59.31 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_ringbuffer_nopool` | 64.78 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_linkedlist_pool` | 60.25 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_linkedlist_nopool` | 63.95 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_ringbuffer_pool` | 60.18 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_ringbuffer_nopool` | 66.26 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_linkedlist_pool` | 365.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_linkedlist_nopool` | 375.40 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_ringbuffer_pool` | 384.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_ringbuffer_nopool` | 420.80 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_linkedlist_pool` | 795.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_linkedlist_nopool` | 827.50 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_ringbuffer_pool` | 836.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_ringbuffer_nopool` | 1078.00 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_linkedlist_pool` | 368.10 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_linkedlist_nopool` | 419.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_ringbuffer_pool` | 385.90 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_ringbuffer_nopool` | 462.70 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_linkedlist_pool` | 356.40 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_linkedlist_nopool` | 367.20 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_ringbuffer_pool` | 364.70 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_ringbuffer_nopool` | 429.30 | 0 | 0 |
| `BenchmarkMarket` | `btree_linkedlist_pool` | 781.30 | 0 | 0 |
| `BenchmarkMarket` | `btree_linkedlist_nopool` | 787.20 | 0 | 0 |
| `BenchmarkMarket` | `btree_ringbuffer_pool` | 822.40 | 0 | 0 |
| `BenchmarkMarket` | `btree_ringbuffer_nopool` | 1038.00 | 0 | 0 |
| `BenchmarkMarket` | `ladder_linkedlist_pool` | 362.30 | 0 | 0 |
| `BenchmarkMarket` | `ladder_linkedlist_nopool` | 417.10 | 0 | 0 |
| `BenchmarkMarket` | `ladder_ringbuffer_pool` | 371.60 | 0 | 0 |
| `BenchmarkMarket` | `ladder_ringbuffer_nopool` | 464.30 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_linkedlist_pool` | 3.13 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_linkedlist_nopool` | 3.24 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_ringbuffer_pool` | 3.13 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_ringbuffer_nopool` | 3.21 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_linkedlist_pool` | 2.17 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_linkedlist_nopool` | 2.14 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_ringbuffer_pool` | 2.14 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_ringbuffer_nopool` | 2.19 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_linkedlist_pool` | 2.31 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_linkedlist_nopool` | 2.36 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_ringbuffer_pool` | 2.35 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_ringbuffer_nopool` | 2.29 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_linkedlist_pool` | 16.02 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_linkedlist_nopool` | 16.42 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_ringbuffer_pool` | 15.92 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_ringbuffer_nopool` | 16.31 | 0 | 0 |
| `BenchmarkDepth` | `btree_linkedlist_pool` | 60.11 | 0 | 0 |
| `BenchmarkDepth` | `btree_linkedlist_nopool` | 59.44 | 0 | 0 |
| `BenchmarkDepth` | `btree_ringbuffer_pool` | 61.68 | 0 | 0 |
| `BenchmarkDepth` | `btree_ringbuffer_nopool` | 57.49 | 0 | 0 |
| `BenchmarkDepth` | `ladder_linkedlist_pool` | 22.27 | 0 | 0 |
| `BenchmarkDepth` | `ladder_linkedlist_nopool` | 22.64 | 0 | 0 |
| `BenchmarkDepth` | `ladder_ringbuffer_pool` | 21.40 | 0 | 0 |
| `BenchmarkDepth` | `ladder_ringbuffer_nopool` | 21.47 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_linkedlist_pool` | 54.65 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_linkedlist_nopool` | 67.70 | 51 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_ringbuffer_pool` | 54.08 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_ringbuffer_nopool` | 69.68 | 51 | 0 |
| `BenchmarkMixedWorkload` | `btree_linkedlist_pool` | 68.16 | 0 | 0 |
| `BenchmarkMixedWorkload` | `btree_linkedlist_nopool` | 78.60 | 51 | 0 |
| `BenchmarkMixedWorkload` | `btree_ringbuffer_pool` | 67.02 | 0 | 0 |
| `BenchmarkMixedWorkload` | `btree_ringbuffer_nopool` | 81.59 | 51 | 0 |
| `BenchmarkMixedWorkload` | `ladder_linkedlist_pool` | 44.44 | 0 | 0 |
| `BenchmarkMixedWorkload` | `ladder_linkedlist_nopool` | 58.20 | 51 | 0 |
| `BenchmarkMixedWorkload` | `ladder_ringbuffer_pool` | 43.61 | 0 | 0 |
| `BenchmarkMixedWorkload` | `ladder_ringbuffer_nopool` | 60.98 | 51 | 0 |
| `BenchmarkEngineChannel` | `default` | 78.91 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_linkedlist_pool` | 58.47 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_linkedlist_nopool` | 92.91 | 96 | 1 |
| `BenchmarkAddLimit` | `rbtree_ringbuffer_pool` | 61.77 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_ringbuffer_nopool` | 99.58 | 96 | 1 |
| `BenchmarkCancel` | `rbtree_linkedlist_pool` | 40.36 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_linkedlist_nopool` | 47.97 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_ringbuffer_pool` | 40.55 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_ringbuffer_nopool` | 61.20 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_linkedlist_pool` | 63.34 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_linkedlist_nopool` | 72.29 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_ringbuffer_pool` | 64.37 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_ringbuffer_nopool` | 69.30 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_linkedlist_pool` | 503.30 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_linkedlist_nopool` | 531.50 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_ringbuffer_pool` | 518.10 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_ringbuffer_nopool` | 626.00 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_linkedlist_pool` | 512.80 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_linkedlist_nopool` | 529.30 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_ringbuffer_pool` | 526.70 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_ringbuffer_nopool` | 639.40 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_linkedlist_pool` | 4.15 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_linkedlist_nopool` | 4.08 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_ringbuffer_pool` | 3.98 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_ringbuffer_nopool` | 3.90 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_linkedlist_pool` | 39.26 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_linkedlist_nopool` | 41.93 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_ringbuffer_pool` | 42.04 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_ringbuffer_nopool` | 42.09 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_linkedlist_pool` | 62.53 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_linkedlist_nopool` | 75.51 | 51 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_ringbuffer_pool` | 64.79 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_ringbuffer_nopool` | 79.04 | 51 | 0 |
| `BenchmarkEngineBatch` | `default` | 22.40 | 0 | 0 |
| `BenchmarkEngineDisruptor` | `default` | 51.98 | 0 | 0 |
| `BenchmarkFastBook` | `AddLimit` | 15.79 | 0 | 0 |
| `BenchmarkFastBook` | `Cancel` | 15.15 | 0 | 0 |
| `BenchmarkFastBook` | `MatchSingle` | 34.03 | 0 | 0 |
| `BenchmarkFastBook` | `MatchSweep` | 213.10 | 0 | 0 |
| `BenchmarkFastBook` | `Market` | 211.50 | 0 | 0 |
| `BenchmarkFastBook` | `BestBidAsk` | 0.94 | 0 | 0 |
| `BenchmarkFastBook` | `Depth` | 22.87 | 0 | 0 |
| `BenchmarkFastBook` | `MixedWorkload` | 24.57 | 0 | 0 |
| `BenchmarkArenaBook` | `AddLimit` | 19.78 | 0 | 0 |
| `BenchmarkArenaBook` | `Cancel` | 13.82 | 0 | 0 |
| `BenchmarkArenaBook` | `MatchSingle` | 35.12 | 0 | 0 |
| `BenchmarkArenaBook` | `MatchSweep` | 212.80 | 0 | 0 |
| `BenchmarkArenaBook` | `Market` | 203.10 | 0 | 0 |
| `BenchmarkArenaBook` | `BestBidAsk` | 1.16 | 0 | 0 |
| `BenchmarkArenaBook` | `Depth` | 27.54 | 0 | 0 |
| `BenchmarkArenaBook` | `MixedWorkload` | 27.24 | 0 | 0 |

## HDR Histogram Latency (Mixed Workload)

| Variant | Throughput (ops/s) | Avg (ns) | p50 (ns) | p90 (ns) | p99 (ns) | p99.9 (ns) | Max (ns) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `sortedslice_linkedlist_pool` | 15748969 | 63.5 | 67 | 102 | 156 | 244 | 84223 |
| `sortedslice_linkedlist_nopool` | 10057084 | 99.4 | 76 | 117 | 197 | 493 | 79039 |
| `sortedslice_ringbuffer_pool` | 16527072 | 60.5 | 68 | 109 | 190 | 1191 | 178815 |
| `sortedslice_ringbuffer_nopool` | 12547316 | 79.7 | 76 | 119 | 244 | 2081 | 170879 |
| `btree_linkedlist_pool` | 13891074 | 72.0 | 75 | 123 | 200 | 348 | 109311 |
| `btree_linkedlist_nopool` | 11676795 | 85.6 | 85 | 135 | 241 | 818 | 26015 |
| `btree_ringbuffer_pool` | 13597491 | 73.5 | 79 | 130 | 224 | 1270 | 38143 |
| `btree_ringbuffer_nopool` | 11289860 | 88.6 | 91 | 147 | 287 | 2105 | 250111 |
| `ladder_linkedlist_pool` | 19936538 | 50.2 | 57 | 87 | 148 | 235 | 135679 |
| `ladder_linkedlist_nopool` | 14386437 | 69.5 | 68 | 107 | 211 | 1733 | 152959 |
| `ladder_ringbuffer_pool` | 18855741 | 53.0 | 64 | 103 | 199 | 866 | 34303 |
| `ladder_ringbuffer_nopool` | 15103409 | 66.2 | 65 | 104 | 220 | 1809 | 179455 |
| `rbtree_linkedlist_pool` | 13754218 | 72.7 | 76 | 124 | 202 | 305 | 28863 |
| `rbtree_linkedlist_nopool` | 11181340 | 89.4 | 91 | 145 | 260 | 519 | 44767 |
| `rbtree_ringbuffer_pool` | 11564123 | 86.5 | 81 | 139 | 248 | 902 | 31567 |
| `rbtree_ringbuffer_nopool` | 9850899 | 101.5 | 87 | 155 | 332 | 2473 | 127039 |

## Profiling Hotspots (pprof — 147 profiled benchmarks)

| Benchmark | Top CPU Hotspot | Top Mem Alloc Space |
|---|---|---|
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (8.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (428.54MB, 97.5%) |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (12.5% flat) | `runtime.mallocgc` (3083.15kB, 38.6%) |
| `BenchmarkAddLimit/btree_ringbuffer_nopool` | `github.com/google/btree.items[go.shape.*uint8].find` (10.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (381.03MB, 97.8%) |
| `BenchmarkAddLimit/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (8.9% flat) | `runtime.mallocgc` (4619.50kB, 31.4%) |
| `BenchmarkAddLimit/ladder_linkedlist_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (14.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (746.57MB, 97.9%) |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9516.36kB, 43.7%) |
| `BenchmarkAddLimit/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/queue.(*RingBuffer).Remove` (11.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (663.56MB, 97.6%) |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.9% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (10003.38kB, 41.1%) |
| `BenchmarkAddLimit/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (13.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (651.06MB, 98.7%) |
| `BenchmarkAddLimit/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.2% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4248.06kB, 33.2%) |
| `BenchmarkAddLimit/rbtree_ringbuffer_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (12.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (705.06MB, 98.3%) |
| `BenchmarkAddLimit/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.7% flat) | `runtime.mallocgc` (6171.96kB, 35.2%) |
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (16.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (447.04MB, 98.1%) |
| `BenchmarkAddLimit/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (15.7% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4956.07kB, 35.0%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (17.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (477.04MB, 98.3%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (12.9% flat) | `runtime.mallocgc` (5660.33kB, 36.7%) |
| `BenchmarkArenaBook/AddLimit` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (41.4% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 34.2%) |
| `BenchmarkArenaBook/BestBidAsk` | `github.com/nedo/orbook/book.(*ArenaBook).BestAsk (inline)` (49.0% flat) | `runtime.mallocgc` (6156.42kB, 37.7%) |
| `BenchmarkArenaBook/Cancel` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (34.5% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 34.6%) |
| `BenchmarkArenaBook/Depth` | `github.com/nedo/orbook/book.(*arenaLadder).topN` (89.5% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 36.8%) |
| `BenchmarkArenaBook/Market` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (75.6% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 36.2%) |
| `BenchmarkArenaBook/MatchSingle` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (79.0% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 35.9%) |
| `BenchmarkArenaBook/MatchSweep` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (73.7% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (7131.34kB, 36.2%) |
| `BenchmarkArenaBook/MixedWorkload` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (55.8% flat) | `github.com/nedo/orbook/book.NewArenaBook` (11037.44kB, 24.4%) |
| `BenchmarkBestBidAsk/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (64.6% flat) | `runtime.mallocgc` (4618.37kB, 62.6%) |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (62.9% flat) | `runtime.mallocgc` (5139.36kB, 48.6%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (53.1% flat) | `runtime.mallocgc` (4620.26kB, 54.6%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (65.1% flat) | `runtime.mallocgc` (3082.21kB, 30.7%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (42.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (11003.72kB, 56.7%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (46.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 43.9%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (47.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 58.6%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (51.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (11003.72kB, 58.6%) |
| `BenchmarkBestBidAsk/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).Best` (53.1% flat) | `runtime.mallocgc` (4106.51kB, 45.5%) |
| `BenchmarkBestBidAsk/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).Best` (52.0% flat) | `runtime.mallocgc` (4618.36kB, 42.5%) |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).Best` (45.5% flat) | `runtime.mallocgc` (5131.37kB, 61.3%) |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/bench_test.BenchmarkBestBidAsk.func1` (45.5% flat) | `runtime.mallocgc` (4618.36kB, 41.4%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (44.6% flat) | `runtime.mallocgc` (3595.20kB, 52.4%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (38.7% flat) | `runtime.mallocgc` (3604.14kB, 37.5%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (36.5% flat) | `runtime.mallocgc` (6.52MB, 84.9%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/bench_test.BenchmarkBestBidAsk.func1` (30.1% flat) | `runtime.mallocgc` (2570.15kB, 42.2%) |
| `BenchmarkCancel/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (8.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (841.08MB, 98.9%) |
| `BenchmarkCancel/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (9.7% flat) | `runtime.mallocgc` (4626.36kB, 37.8%) |
| `BenchmarkCancel/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (8.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (821.58MB, 98.5%) |
| `BenchmarkCancel/btree_ringbuffer_pool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (9.9% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 34.2%) |
| `BenchmarkCancel/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (14.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1004.59MB, 98.3%) |
| `BenchmarkCancel/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (20.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8516.02kB, 36.1%) |
| `BenchmarkCancel/ladder_ringbuffer_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (10.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1034.09MB, 98.4%) |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 43.4%) |
| `BenchmarkCancel/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (11.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1220.61MB, 99.3%) |
| `BenchmarkCancel/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (16.7% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4956.07kB, 32.6%) |
| `BenchmarkCancel/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (8.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1312.62MB, 99.3%) |
| `BenchmarkCancel/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.5% flat) | `runtime.mallocgc` (5644.37kB, 36.6%) |
| `BenchmarkCancel/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (16.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (902.58MB, 99.0%) |
| `BenchmarkCancel/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (14.8% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 33.8%) |
| `BenchmarkCancel/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (16.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (900.08MB, 98.8%) |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (22.1% flat) | `runtime.mallocgc` (2569.21kB, 18.4%) |
| `BenchmarkDepth/btree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (55.1% flat) | `runtime.mallocgc` (5138.41kB, 59.6%) |
| `BenchmarkDepth/btree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (52.1% flat) | `runtime.mallocgc` (5645.32kB, 53.1%) |
| `BenchmarkDepth/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (64.0% flat) | `runtime.mallocgc` (3082.20kB, 39.9%) |
| `BenchmarkDepth/btree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (55.1% flat) | `runtime.mallocgc` (4620.77kB, 46.5%) |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (88.9% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 56.6%) |
| `BenchmarkDepth/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (92.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 47.7%) |
| `BenchmarkDepth/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (88.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 50.7%) |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (91.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 65.1%) |
| `BenchmarkDepth/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (50.8% flat) | `runtime.mallocgc` (4619.31kB, 52.6%) |
| `BenchmarkDepth/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (41.9% flat) | `runtime.mallocgc` (4122.29kB, 39.8%) |
| `BenchmarkDepth/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (53.7% flat) | `runtime.mallocgc` (2570.15kB, 36.9%) |
| `BenchmarkDepth/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (60.3% flat) | `runtime.mallocgc` (2571.10kB, 41.5%) |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (90.5% flat) | `runtime.mallocgc` (3603.19kB, 47.4%) |
| `BenchmarkDepth/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (88.6% flat) | `runtime.mallocgc` (4618.87kB, 51.7%) |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (93.1% flat) | `runtime.mallocgc` (2058.10kB, 31.7%) |
| `BenchmarkDepth/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (90.1% flat) | `runtime.mallocgc` (5131.88kB, 48.3%) |
| `BenchmarkEngineBatch` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (18.5% flat) | `github.com/nedo/orbook/bench_test.BenchmarkEngineBatch` (70.63MB, 61.8%) |
| `BenchmarkEngineChannel` | `runtime.procyieldAsm` (40.0% flat) | `github.com/nedo/orbook/engine.New (inline)` (70.93MB, 61.0%) |
| `BenchmarkEngineDisruptor` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.6% flat) | `github.com/nedo/orbook/engine.NewDisruptor` (98.59MB, 65.3%) |
| `BenchmarkFastBook/AddLimit` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (34.5% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (7002.37kB, 30.9%) |
| `BenchmarkFastBook/BestBidAsk` | `github.com/nedo/orbook/book.(*FastBook).BestAsk (inline)` (43.5% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (10003.38kB, 42.8%) |
| `BenchmarkFastBook/Cancel` | `github.com/nedo/orbook/book.(*FastBook).CancelOrder` (39.4% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (9003.04kB, 39.4%) |
| `BenchmarkFastBook/Depth` | `github.com/nedo/orbook/book.(*fastLadder).topN` (97.4% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (9003.04kB, 48.4%) |
| `BenchmarkFastBook/Market` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (66.4% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (6002.03kB, 29.9%) |
| `BenchmarkFastBook/MatchSingle` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (71.2% flat) | `github.com/nedo/orbook/book.newFastLadder` (10003.38kB, 41.6%) |
| `BenchmarkFastBook/MatchSweep` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (68.9% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (8002.70kB, 34.5%) |
| `BenchmarkFastBook/MixedWorkload` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (38.1% flat) | `github.com/nedo/orbook/book.(*FastBook).growPool` (17700.24kB, 37.2%) |
| `BenchmarkMarket/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (410.54MB, 83.7%) |
| `BenchmarkMarket/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.0% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (47.52MB, 51.0%) |
| `BenchmarkMarket/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (745.57MB, 84.5%) |
| `BenchmarkMarket/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (17.5% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (47.02MB, 55.2%) |
| `BenchmarkMarket/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.57GB, 99.2%) |
| `BenchmarkMarket/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 36.6%) |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1435.13MB, 98.9%) |
| `BenchmarkMarket/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 26.4%) |
| `BenchmarkMarket/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (13.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (928.59MB, 98.8%) |
| `BenchmarkMarket/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.2% flat) | `runtime.mallocgc` (5650.65kB, 34.2%) |
| `BenchmarkMarket/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (922.08MB, 98.5%) |
| `BenchmarkMarket/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (17.4% flat) | `runtime.mallocgc` (4618.37kB, 31.7%) |
| `BenchmarkMarket/sortedslice_linkedlist_nopool` | `runtime.memmove` (41.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (911.08MB, 98.8%) |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `runtime.memmove` (53.9% flat) | `runtime.mallocgc` (5138.41kB, 28.2%) |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `runtime.memmove` (39.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (945.09MB, 98.7%) |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `runtime.memmove` (46.4% flat) | `runtime.mallocgc` (2570.15kB, 19.8%) |
| `BenchmarkMatchSingle/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1131.10MB, 99.3%) |
| `BenchmarkMatchSingle/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.5% flat) | `runtime.mallocgc` (4107.26kB, 43.1%) |
| `BenchmarkMatchSingle/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (605.06MB, 98.2%) |
| `BenchmarkMatchSingle/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (34.8% flat) | `runtime.mallocgc` (6170.51kB, 42.1%) |
| `BenchmarkMatchSingle/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (635.06MB, 97.7%) |
| `BenchmarkMatchSingle/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 34.7%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (26.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (663.56MB, 97.5%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (40.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 33.2%) |
| `BenchmarkMatchSingle/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (848.58MB, 98.8%) |
| `BenchmarkMatchSingle/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (41.9% flat) | `runtime.mallocgc` (5139.86kB, 41.2%) |
| `BenchmarkMatchSingle/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (790.57MB, 99.1%) |
| `BenchmarkMatchSingle/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (38.6% flat) | `runtime.mallocgc` (5652.36kB, 40.4%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (648.56MB, 98.9%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (39.5% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 40.5%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (710.07MB, 98.6%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.6% flat) | `runtime.mallocgc` (6676.47kB, 41.4%) |
| `BenchmarkMatchSweep/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (405.04MB, 84.1%) |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.8% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (46.52MB, 48.7%) |
| `BenchmarkMatchSweep/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (749.57MB, 85.6%) |
| `BenchmarkMatchSweep/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.5% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (48.02MB, 47.5%) |
| `BenchmarkMatchSweep/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1682.65MB, 99.1%) |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 40.2%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1537.14MB, 99.0%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 33.0%) |
| `BenchmarkMatchSweep/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (905.58MB, 99.0%) |
| `BenchmarkMatchSweep/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (17.5% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (3540.05kB, 29.7%) |
| `BenchmarkMatchSweep/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (917.58MB, 98.5%) |
| `BenchmarkMatchSweep/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.5% flat) | `runtime.mallocgc` (3595.21kB, 19.6%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_nopool` | `runtime.memmove` (41.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (885.58MB, 98.9%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_pool` | `runtime.memmove` (46.7% flat) | `runtime.mallocgc` (6156.42kB, 46.5%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_nopool` | `runtime.memmove` (36.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (905.08MB, 98.7%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_pool` | `runtime.memmove` (47.3% flat) | `runtime.mallocgc` (3595.21kB, 27.8%) |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (12.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (338.03MB, 94.7%) |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (17700.24kB, 44.5%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (336.03MB, 92.5%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.3% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (15576.21kB, 36.4%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (398.54MB, 93.7%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (28.3% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (16284.22kB, 36.2%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (418.04MB, 92.8%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.4% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (13452.18kB, 28.4%) |
| `BenchmarkMixedWorkload/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (18.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (508.55MB, 96.2%) |
| `BenchmarkMixedWorkload/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.6% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (17700.24kB, 56.0%) |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (479.04MB, 95.2%) |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (17.1% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (13452.18kB, 36.2%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (385.04MB, 95.1%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (14868.20kB, 46.0%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (15.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (377.53MB, 93.6%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (13.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (15576.21kB, 38.6%) |

*(Full top-15 CPU & memory tables for all profiled variants are in `latest.json` and `latest.html`.)*

