# Orbook Benchmark & Profile Report — `v006_20261123_030147`

- **Version:** `v006` | **Mode:** `full`
- **Time:** `2026-11-23 03:01:47 WIB` | **Commit:** `unknown` | **Go:** `go1.27.3-X:nodwarf5` | **CPU:** `Intel(R) Core(TM) i9-14900HX`

## Diff vs Previous (`v005_20261123_024831`, threshold ±2.0%)

| Benchmark | Status | Old (ns/op) | New (ns/op) | Delta | Allocs |
|---|---|---:|---:|---:|---|
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `slower` | 92.13 | 95.32 | +3.46% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/sortedslice_ringbuffer_pool` | `faster` | 61.68 | 59.13 | -4.13% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `slower` | 94.89 | 104.20 | +9.81% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `faster` | 83.19 | 79.39 | -4.57% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `slower` | 112.00 | 114.70 | +2.41% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/btree_ringbuffer_pool` | `slower` | 80.72 | 82.56 | +2.28% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `faster` | 30.91 | 29.86 | -3.40% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `slower` | 32.55 | 33.22 | +2.06% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/sortedslice_linkedlist_pool` | `slower` | 37.70 | 39.02 | +3.50% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `slower` | 37.20 | 38.01 | +2.18% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/sortedslice_ringbuffer_nopool` | `slower` | 47.90 | 49.72 | +3.80% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_ringbuffer_pool` | `slower` | 37.80 | 39.67 | +4.95% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_ringbuffer_nopool` | `slower` | 51.84 | 54.27 | +4.69% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_linkedlist_pool` | `slower` | 37.13 | 39.23 | +5.66% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `slower` | 36.38 | 37.73 | +3.71% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `slower` | 63.48 | 65.04 | +2.46% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_linkedlist_pool` | `slower` | 58.56 | 60.25 | +2.89% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_linkedlist_nopool` | `faster` | 67.18 | 63.95 | -4.81% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `faster` | 828.40 | 795.60 | -3.96% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `faster` | 393.70 | 368.10 | -6.50% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_linkedlist_nopool` | `faster` | 441.50 | 419.60 | -4.96% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `faster` | 377.40 | 356.40 | -5.56% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `faster` | 377.70 | 364.70 | -3.44% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `slower` | 406.40 | 429.30 | +5.63% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_linkedlist_nopool` | `faster` | 823.70 | 787.20 | -4.43% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `faster` | 480.00 | 464.30 | -3.27% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_pool` | `faster` | 3.52 | 3.13 | -11.10% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `faster` | 3.59 | 3.24 | -9.70% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `faster` | 3.68 | 3.13 | -14.76% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_nopool` | `faster` | 3.68 | 3.21 | -12.85% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `faster` | 2.40 | 2.17 | -9.62% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_linkedlist_nopool` | `faster` | 2.53 | 2.14 | -15.42% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `faster` | 2.36 | 2.14 | -9.52% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `faster` | 2.32 | 2.19 | -5.56% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `faster` | 2.61 | 2.31 | -11.51% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `faster` | 2.75 | 2.36 | -14.26% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_ringbuffer_pool` | `faster` | 2.75 | 2.35 | -14.49% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `faster` | 2.52 | 2.29 | -9.15% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_linkedlist_pool` | `faster` | 16.52 | 16.02 | -3.03% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `slower` | 15.95 | 16.42 | +2.95% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `slower` | 15.88 | 16.31 | +2.71% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_linkedlist_pool` | `slower` | 57.36 | 60.11 | +4.79% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_ringbuffer_pool` | `slower` | 58.45 | 61.68 | +5.53% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_linkedlist_pool` | `slower` | 21.62 | 22.27 | +3.01% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `slower` | 22.00 | 22.64 | +2.91% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `faster` | 21.88 | 21.40 | -2.19% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_ringbuffer_nopool` | `faster` | 21.98 | 21.47 | -2.32% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `slower` | 52.59 | 54.65 | +3.92% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `faster` | 55.19 | 54.08 | -2.01% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_nopool` | `faster` | 72.52 | 69.68 | -3.92% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `slower` | 65.61 | 68.16 | +3.89% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `faster` | 81.33 | 78.60 | -3.36% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `slower` | 43.45 | 44.44 | +2.28% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_linkedlist_nopool` | `slower` | 57.04 | 58.20 | +2.03% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `slower` | 58.67 | 60.98 | +3.94% | 51→51 B, 0→0 allocs |

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
| `BenchmarkEngineChannel` | `default` | 90.57 | 0 | 0 |

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

## Profiling Hotspots (pprof — 97 profiled benchmarks)

| Benchmark | Top CPU Hotspot | Top Mem Alloc Space |
|---|---|---|
| `BenchmarkAddLimit/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (15.7% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4956.07kB, 35.0%) |
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (16.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (447.04MB, 98.1%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (12.9% flat) | `runtime.mallocgc` (5660.33kB, 36.7%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (17.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (477.04MB, 98.3%) |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (12.5% flat) | `runtime.mallocgc` (3083.15kB, 38.6%) |
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (8.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (428.54MB, 97.5%) |
| `BenchmarkAddLimit/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (8.9% flat) | `runtime.mallocgc` (4619.50kB, 31.4%) |
| `BenchmarkAddLimit/btree_ringbuffer_nopool` | `github.com/google/btree.items[go.shape.*uint8].find` (10.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (381.03MB, 97.8%) |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9516.36kB, 43.7%) |
| `BenchmarkAddLimit/ladder_linkedlist_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (14.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (746.57MB, 97.9%) |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.9% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (10003.38kB, 41.1%) |
| `BenchmarkAddLimit/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/queue.(*RingBuffer).Remove` (11.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (663.56MB, 97.6%) |
| `BenchmarkCancel/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (14.8% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 33.8%) |
| `BenchmarkCancel/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (16.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (902.58MB, 99.0%) |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (22.1% flat) | `runtime.mallocgc` (2569.21kB, 18.4%) |
| `BenchmarkCancel/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (16.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (900.08MB, 98.8%) |
| `BenchmarkCancel/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (9.7% flat) | `runtime.mallocgc` (4626.36kB, 37.8%) |
| `BenchmarkCancel/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (8.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (841.08MB, 98.9%) |
| `BenchmarkCancel/btree_ringbuffer_pool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (9.9% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 34.2%) |
| `BenchmarkCancel/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (8.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (821.58MB, 98.5%) |
| `BenchmarkCancel/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (20.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8516.02kB, 36.1%) |
| `BenchmarkCancel/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (14.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1004.59MB, 98.3%) |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 43.4%) |
| `BenchmarkCancel/ladder_ringbuffer_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (10.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1034.09MB, 98.4%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (39.5% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 40.5%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (648.56MB, 98.9%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.6% flat) | `runtime.mallocgc` (6676.47kB, 41.4%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (710.07MB, 98.6%) |
| `BenchmarkMatchSingle/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.5% flat) | `runtime.mallocgc` (4107.26kB, 43.1%) |
| `BenchmarkMatchSingle/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1131.10MB, 99.3%) |
| `BenchmarkMatchSingle/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (34.8% flat) | `runtime.mallocgc` (6170.51kB, 42.1%) |
| `BenchmarkMatchSingle/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (605.06MB, 98.2%) |
| `BenchmarkMatchSingle/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 34.7%) |
| `BenchmarkMatchSingle/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (635.06MB, 97.7%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (40.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 33.2%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (26.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (663.56MB, 97.5%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_pool` | `runtime.memmove` (46.7% flat) | `runtime.mallocgc` (6156.42kB, 46.5%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_nopool` | `runtime.memmove` (41.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (885.58MB, 98.9%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_pool` | `runtime.memmove` (47.3% flat) | `runtime.mallocgc` (3595.21kB, 27.8%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_nopool` | `runtime.memmove` (36.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (905.08MB, 98.7%) |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.8% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (46.52MB, 48.7%) |
| `BenchmarkMatchSweep/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (405.04MB, 84.1%) |
| `BenchmarkMatchSweep/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.5% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (48.02MB, 47.5%) |
| `BenchmarkMatchSweep/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (749.57MB, 85.6%) |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 40.2%) |
| `BenchmarkMatchSweep/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1682.65MB, 99.1%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 33.0%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1537.14MB, 99.0%) |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `runtime.memmove` (53.9% flat) | `runtime.mallocgc` (5138.41kB, 28.2%) |
| `BenchmarkMarket/sortedslice_linkedlist_nopool` | `runtime.memmove` (41.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (911.08MB, 98.8%) |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `runtime.memmove` (46.4% flat) | `runtime.mallocgc` (2570.15kB, 19.8%) |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `runtime.memmove` (39.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (945.09MB, 98.7%) |
| `BenchmarkMarket/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.0% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (47.52MB, 51.0%) |
| `BenchmarkMarket/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (410.54MB, 83.7%) |
| `BenchmarkMarket/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (17.5% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (47.02MB, 55.2%) |
| `BenchmarkMarket/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (745.57MB, 84.5%) |
| `BenchmarkMarket/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 36.6%) |
| `BenchmarkMarket/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.57GB, 99.2%) |
| `BenchmarkMarket/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 26.4%) |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1435.13MB, 98.9%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (38.7% flat) | `runtime.mallocgc` (3604.14kB, 37.5%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (44.6% flat) | `runtime.mallocgc` (3595.20kB, 52.4%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/bench_test.BenchmarkBestBidAsk.func1` (30.1% flat) | `runtime.mallocgc` (2570.15kB, 42.2%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (36.5% flat) | `runtime.mallocgc` (6.52MB, 84.9%) |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (62.9% flat) | `runtime.mallocgc` (5139.36kB, 48.6%) |
| `BenchmarkBestBidAsk/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (64.6% flat) | `runtime.mallocgc` (4618.37kB, 62.6%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (65.1% flat) | `runtime.mallocgc` (3082.21kB, 30.7%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (53.1% flat) | `runtime.mallocgc` (4620.26kB, 54.6%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (46.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 43.9%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (42.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (11003.72kB, 56.7%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (51.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (11003.72kB, 58.6%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (47.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 58.6%) |
| `BenchmarkDepth/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (88.6% flat) | `runtime.mallocgc` (4618.87kB, 51.7%) |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (90.5% flat) | `runtime.mallocgc` (3603.19kB, 47.4%) |
| `BenchmarkDepth/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (90.1% flat) | `runtime.mallocgc` (5131.88kB, 48.3%) |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (93.1% flat) | `runtime.mallocgc` (2058.10kB, 31.7%) |
| `BenchmarkDepth/btree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (52.1% flat) | `runtime.mallocgc` (5645.32kB, 53.1%) |
| `BenchmarkDepth/btree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (55.1% flat) | `runtime.mallocgc` (5138.41kB, 59.6%) |
| `BenchmarkDepth/btree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (55.1% flat) | `runtime.mallocgc` (4620.77kB, 46.5%) |
| `BenchmarkDepth/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (64.0% flat) | `runtime.mallocgc` (3082.20kB, 39.9%) |
| `BenchmarkDepth/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (92.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 47.7%) |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (88.9% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 56.6%) |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (91.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 65.1%) |
| `BenchmarkDepth/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (88.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 50.7%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (14868.20kB, 46.0%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (385.04MB, 95.1%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (13.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (15576.21kB, 38.6%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (15.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (377.53MB, 93.6%) |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (17700.24kB, 44.5%) |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `internal/runtime/maps.ctrlGroup.matchH2 (inline)` (12.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (338.03MB, 94.7%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.3% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (15576.21kB, 36.4%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (336.03MB, 92.5%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (28.3% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (16284.22kB, 36.2%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (398.54MB, 93.7%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.4% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (13452.18kB, 28.4%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (418.04MB, 92.8%) |
| `BenchmarkEngineChannel` | `runtime.procyieldAsm` (23.0% flat) | `github.com/nedo/orbook/engine.New (inline)` (64.13MB, 59.1%) |

*(Full top-15 CPU & memory tables for all profiled variants are in `latest.json` and `latest.html`.)*

