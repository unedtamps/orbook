# Orbook Benchmark & Profile Report — `v010_20261123_125155`

- **Version:** `v010` | **Mode:** `full`
- **Time:** `2026-11-23 12:51:55 WIB` | **Commit:** `4e16c62` | **Go:** `go1.27.3-X:nodwarf5` | **CPU:** `Intel(R) Core(TM) i9-14900HX`

## Diff vs Previous (`v009_20261123_124124`, threshold ±2.0%)

| Benchmark | Status | Old (ns/op) | New (ns/op) | Delta | Allocs |
|---|---|---:|---:|---:|---|
| `BenchmarkAddLimit/sortedslice_linkedlist_pool` | `faster` | 47.12 | 44.34 | -5.90% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `faster` | 83.80 | 79.72 | -4.87% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `slower` | 78.56 | 82.18 | +4.61% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `faster` | 68.19 | 62.11 | -8.92% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `slower` | 95.58 | 98.65 | +3.21% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/btree_ringbuffer_nopool` | `faster` | 102.70 | 97.13 | -5.42% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/rbtree_linkedlist_pool` | `faster` | 41.00 | 39.55 | -3.54% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/rbtree_linkedlist_nopool` | `slower` | 68.06 | 71.11 | +4.48% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/rbtree_ringbuffer_pool` | `faster` | 43.10 | 40.67 | -5.64% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `slower` | 20.37 | 21.38 | +4.96% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_linkedlist_nopool` | `faster` | 46.71 | 44.25 | -5.27% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `faster` | 22.98 | 22.24 | -3.22% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_ringbuffer_nopool` | `faster` | 49.62 | 45.76 | -7.78% | 96→96 B, 1→1 allocs |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `faster` | 20.37 | 19.93 | -2.16% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_linkedlist_pool` | `faster` | 21.42 | 19.97 | -6.77% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_linkedlist_nopool` | `faster` | 36.87 | 34.43 | -6.62% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_linkedlist_pool` | `faster` | 19.78 | 19.36 | -2.12% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_linkedlist_nopool` | `faster` | 34.63 | 31.21 | -9.88% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_ringbuffer_nopool` | `slower` | 25.62 | 26.25 | +2.46% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_linkedlist_pool` | `faster` | 17.87 | 17.18 | -3.86% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_linkedlist_nopool` | `slower` | 25.36 | 28.57 | +12.66% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `faster` | 17.79 | 17.22 | -3.20% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_ringbuffer_nopool` | `faster` | 26.10 | 24.71 | -5.33% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_linkedlist_pool` | `faster` | 37.83 | 35.00 | -7.48% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `slower` | 38.99 | 41.98 | +7.67% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_linkedlist_pool` | `slower` | 35.05 | 35.90 | +2.43% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_linkedlist_nopool` | `faster` | 44.28 | 43.02 | -2.85% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_ringbuffer_pool` | `faster` | 38.77 | 37.08 | -4.36% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/rbtree_linkedlist_nopool` | `slower` | 43.21 | 44.72 | +3.49% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/rbtree_ringbuffer_nopool` | `slower` | 40.26 | 43.00 | +6.81% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_ringbuffer_nopool` | `slower` | 42.51 | 43.89 | +3.25% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_linkedlist_pool` | `faster` | 258.60 | 251.00 | -2.94% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_linkedlist_nopool` | `faster` | 248.20 | 242.80 | -2.18% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_nopool` | `faster` | 281.50 | 275.50 | -2.13% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `faster` | 642.40 | 582.20 | -9.37% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_ringbuffer_pool` | `faster` | 632.40 | 596.20 | -5.72% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_ringbuffer_nopool` | `slower` | 759.10 | 884.10 | +16.47% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_linkedlist_pool` | `faster` | 317.10 | 299.20 | -5.64% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_linkedlist_nopool` | `faster` | 345.80 | 338.80 | -2.02% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_ringbuffer_pool` | `faster` | 334.40 | 298.30 | -10.80% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_ringbuffer_nopool` | `faster` | 420.80 | 399.00 | -5.18% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `faster` | 246.70 | 240.90 | -2.35% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_ringbuffer_pool` | `faster` | 261.00 | 240.40 | -7.89% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `faster` | 247.80 | 240.00 | -3.15% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_linkedlist_nopool` | `faster` | 260.00 | 239.90 | -7.73% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `faster` | 262.00 | 251.00 | -4.20% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `faster` | 276.60 | 267.80 | -3.18% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_linkedlist_pool` | `faster` | 618.00 | 590.90 | -4.39% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_ringbuffer_pool` | `faster` | 618.60 | 588.10 | -4.93% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_ringbuffer_nopool` | `slower` | 800.10 | 849.90 | +6.22% | 1→0 B, 0→0 allocs |
| `BenchmarkMarket/rbtree_ringbuffer_pool` | `faster` | 319.90 | 305.10 | -4.63% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/rbtree_ringbuffer_nopool` | `faster` | 410.90 | 388.90 | -5.35% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_linkedlist_nopool` | `faster` | 311.10 | 278.80 | -10.38% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_ringbuffer_pool` | `faster` | 250.70 | 242.20 | -3.39% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `faster` | 320.50 | 300.90 | -6.12% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `faster` | 3.39 | 3.20 | -5.60% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `faster` | 3.31 | 3.13 | -5.14% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `faster` | 2.22 | 2.12 | -4.68% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `faster` | 2.24 | 2.13 | -5.09% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `faster` | 2.29 | 2.16 | -5.97% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_linkedlist_pool` | `faster` | 3.38 | 3.26 | -3.73% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_pool` | `faster` | 3.37 | 3.15 | -6.47% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_nopool` | `faster` | 3.36 | 3.15 | -6.10% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `slower` | 3.12 | 3.21 | +2.86% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `faster` | 3.30 | 3.11 | -5.88% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `faster` | 3.15 | 3.07 | -2.38% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `faster` | 18.28 | 16.81 | -8.04% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_ringbuffer_pool` | `faster` | 17.37 | 16.21 | -6.68% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `faster` | 16.85 | 16.17 | -4.04% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_linkedlist_pool` | `slower` | 60.94 | 62.22 | +2.10% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_linkedlist_nopool` | `faster` | 59.20 | 57.19 | -3.40% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_ringbuffer_pool` | `slower` | 58.62 | 62.22 | +6.14% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/rbtree_linkedlist_nopool` | `faster` | 38.73 | 37.86 | -2.25% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/rbtree_ringbuffer_nopool` | `faster` | 40.00 | 39.11 | -2.23% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_linkedlist_pool` | `faster` | 22.37 | 20.33 | -9.12% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `faster` | 21.53 | 20.49 | -4.83% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `faster` | 22.02 | 21.08 | -4.27% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `faster` | 40.33 | 38.68 | -4.09% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_nopool` | `faster` | 57.92 | 55.40 | -4.35% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `faster` | 43.09 | 40.82 | -5.27% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `faster` | 54.80 | 51.68 | -5.69% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `faster` | 75.08 | 68.81 | -8.35% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_ringbuffer_pool` | `faster` | 56.45 | 51.90 | -8.06% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_ringbuffer_nopool` | `alloc_changed` | 68.60 | 68.74 | +0.20% | 52→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_linkedlist_pool` | `faster` | 39.76 | 38.47 | -3.24% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_linkedlist_nopool` | `faster` | 55.64 | 53.56 | -3.74% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_pool` | `faster` | 41.03 | 39.65 | -3.36% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_nopool` | `faster` | 59.19 | 54.71 | -7.57% | 52→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `faster` | 30.68 | 28.58 | -6.84% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `faster` | 48.87 | 45.75 | -6.38% | 51→51 B, 0→0 allocs |
| `BenchmarkEngineChannel` | `faster` | 84.01 | 80.47 | -4.21% | 0→0 B, 0→0 allocs |
| `BenchmarkEngineBatch` | `faster` | 22.73 | 21.66 | -4.71% | 0→0 B, 0→0 allocs |
| `BenchmarkEngineDisruptor` | `faster` | 53.74 | 51.03 | -5.04% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/MatchSweep` | `faster` | 207.50 | 202.40 | -2.46% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/Market` | `faster` | 211.50 | 207.20 | -2.03% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/BestBidAsk` | `faster` | 0.98 | 0.92 | -6.82% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/AddLimit` | `faster` | 20.10 | 19.36 | -3.68% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/MatchSingle` | `faster` | 32.57 | 31.82 | -2.30% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/Market` | `slower` | 195.80 | 200.70 | +2.50% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/BestBidAsk` | `faster` | 1.19 | 1.06 | -10.98% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/Depth` | `faster` | 26.80 | 24.68 | -7.91% | 0→0 B, 0→0 allocs |

## Benchmarks

| Method | Variant | ns/op | B/op | allocs/op |
|---|---|---:|---:|---:|
| `BenchmarkAddLimit` | `sortedslice_linkedlist_pool` | 44.34 | 0 | 0 |
| `BenchmarkAddLimit` | `sortedslice_linkedlist_nopool` | 79.72 | 96 | 1 |
| `BenchmarkAddLimit` | `sortedslice_ringbuffer_pool` | 47.78 | 0 | 0 |
| `BenchmarkAddLimit` | `sortedslice_ringbuffer_nopool` | 82.18 | 96 | 1 |
| `BenchmarkAddLimit` | `btree_linkedlist_pool` | 62.11 | 0 | 0 |
| `BenchmarkAddLimit` | `btree_linkedlist_nopool` | 98.65 | 96 | 1 |
| `BenchmarkAddLimit` | `btree_ringbuffer_pool` | 67.28 | 0 | 0 |
| `BenchmarkAddLimit` | `btree_ringbuffer_nopool` | 97.13 | 96 | 1 |
| `BenchmarkAddLimit` | `rbtree_linkedlist_pool` | 39.55 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_linkedlist_nopool` | 71.11 | 96 | 1 |
| `BenchmarkAddLimit` | `rbtree_ringbuffer_pool` | 40.67 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_ringbuffer_nopool` | 77.07 | 96 | 1 |
| `BenchmarkAddLimit` | `ladder_linkedlist_pool` | 21.38 | 0 | 0 |
| `BenchmarkAddLimit` | `ladder_linkedlist_nopool` | 44.25 | 96 | 1 |
| `BenchmarkAddLimit` | `ladder_ringbuffer_pool` | 22.24 | 0 | 0 |
| `BenchmarkAddLimit` | `ladder_ringbuffer_nopool` | 45.76 | 96 | 1 |
| `BenchmarkCancel` | `sortedslice_linkedlist_pool` | 20.34 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_linkedlist_nopool` | 34.56 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_ringbuffer_pool` | 19.93 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_ringbuffer_nopool` | 28.39 | 0 | 0 |
| `BenchmarkCancel` | `btree_linkedlist_pool` | 19.97 | 0 | 0 |
| `BenchmarkCancel` | `btree_linkedlist_nopool` | 34.43 | 0 | 0 |
| `BenchmarkCancel` | `btree_ringbuffer_pool` | 19.82 | 0 | 0 |
| `BenchmarkCancel` | `btree_ringbuffer_nopool` | 27.79 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_linkedlist_pool` | 19.36 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_linkedlist_nopool` | 31.21 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_ringbuffer_pool` | 19.85 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_ringbuffer_nopool` | 26.25 | 0 | 0 |
| `BenchmarkCancel` | `ladder_linkedlist_pool` | 17.18 | 0 | 0 |
| `BenchmarkCancel` | `ladder_linkedlist_nopool` | 28.57 | 0 | 0 |
| `BenchmarkCancel` | `ladder_ringbuffer_pool` | 17.22 | 0 | 0 |
| `BenchmarkCancel` | `ladder_ringbuffer_nopool` | 24.71 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_linkedlist_pool` | 35.00 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_linkedlist_nopool` | 43.56 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_ringbuffer_pool` | 34.94 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_ringbuffer_nopool` | 41.98 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_linkedlist_pool` | 35.90 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_linkedlist_nopool` | 43.02 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_ringbuffer_pool` | 37.08 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_ringbuffer_nopool` | 41.31 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_linkedlist_pool` | 35.24 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_linkedlist_nopool` | 44.72 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_ringbuffer_pool` | 37.65 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_ringbuffer_nopool` | 43.00 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_linkedlist_pool` | 36.19 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_linkedlist_nopool` | 46.20 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_ringbuffer_pool` | 35.85 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_ringbuffer_nopool` | 43.89 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_linkedlist_pool` | 251.00 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_linkedlist_nopool` | 242.80 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_ringbuffer_pool` | 254.10 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_ringbuffer_nopool` | 275.50 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_linkedlist_pool` | 582.20 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_linkedlist_nopool` | 627.30 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_ringbuffer_pool` | 596.20 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_ringbuffer_nopool` | 884.10 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_linkedlist_pool` | 299.20 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_linkedlist_nopool` | 338.80 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_ringbuffer_pool` | 298.30 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_ringbuffer_nopool` | 399.00 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_linkedlist_pool` | 240.90 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_linkedlist_nopool` | 295.20 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_ringbuffer_pool` | 240.40 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_ringbuffer_nopool` | 320.80 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_linkedlist_pool` | 240.00 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_linkedlist_nopool` | 239.90 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_ringbuffer_pool` | 251.00 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_ringbuffer_nopool` | 267.80 | 0 | 0 |
| `BenchmarkMarket` | `btree_linkedlist_pool` | 590.90 | 0 | 0 |
| `BenchmarkMarket` | `btree_linkedlist_nopool` | 610.60 | 0 | 0 |
| `BenchmarkMarket` | `btree_ringbuffer_pool` | 588.10 | 0 | 0 |
| `BenchmarkMarket` | `btree_ringbuffer_nopool` | 849.90 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_linkedlist_pool` | 297.40 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_linkedlist_nopool` | 351.40 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_ringbuffer_pool` | 305.10 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_ringbuffer_nopool` | 388.90 | 0 | 0 |
| `BenchmarkMarket` | `ladder_linkedlist_pool` | 241.80 | 0 | 0 |
| `BenchmarkMarket` | `ladder_linkedlist_nopool` | 278.80 | 0 | 0 |
| `BenchmarkMarket` | `ladder_ringbuffer_pool` | 242.20 | 0 | 0 |
| `BenchmarkMarket` | `ladder_ringbuffer_nopool` | 300.90 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_linkedlist_pool` | 3.13 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_linkedlist_nopool` | 3.20 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_ringbuffer_pool` | 3.13 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_ringbuffer_nopool` | 3.21 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_linkedlist_pool` | 2.12 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_linkedlist_nopool` | 2.20 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_ringbuffer_pool` | 2.13 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_ringbuffer_nopool` | 2.16 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_linkedlist_pool` | 3.26 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_linkedlist_nopool` | 3.27 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_ringbuffer_pool` | 3.15 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_ringbuffer_nopool` | 3.15 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_linkedlist_pool` | 3.21 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_linkedlist_nopool` | 3.11 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_ringbuffer_pool` | 3.08 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_ringbuffer_nopool` | 3.07 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_linkedlist_pool` | 16.62 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_linkedlist_nopool` | 16.81 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_ringbuffer_pool` | 16.21 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_ringbuffer_nopool` | 16.17 | 0 | 0 |
| `BenchmarkDepth` | `btree_linkedlist_pool` | 62.22 | 0 | 0 |
| `BenchmarkDepth` | `btree_linkedlist_nopool` | 57.19 | 0 | 0 |
| `BenchmarkDepth` | `btree_ringbuffer_pool` | 62.22 | 0 | 0 |
| `BenchmarkDepth` | `btree_ringbuffer_nopool` | 60.01 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_linkedlist_pool` | 38.52 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_linkedlist_nopool` | 37.86 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_ringbuffer_pool` | 39.56 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_ringbuffer_nopool` | 39.11 | 0 | 0 |
| `BenchmarkDepth` | `ladder_linkedlist_pool` | 20.33 | 0 | 0 |
| `BenchmarkDepth` | `ladder_linkedlist_nopool` | 20.49 | 0 | 0 |
| `BenchmarkDepth` | `ladder_ringbuffer_pool` | 21.08 | 0 | 0 |
| `BenchmarkDepth` | `ladder_ringbuffer_nopool` | 21.23 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_linkedlist_pool` | 38.68 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_linkedlist_nopool` | 55.40 | 51 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_ringbuffer_pool` | 40.82 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_ringbuffer_nopool` | 55.97 | 51 | 0 |
| `BenchmarkMixedWorkload` | `btree_linkedlist_pool` | 51.68 | 0 | 0 |
| `BenchmarkMixedWorkload` | `btree_linkedlist_nopool` | 68.81 | 51 | 0 |
| `BenchmarkMixedWorkload` | `btree_ringbuffer_pool` | 51.90 | 0 | 0 |
| `BenchmarkMixedWorkload` | `btree_ringbuffer_nopool` | 68.74 | 51 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_linkedlist_pool` | 38.47 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_linkedlist_nopool` | 53.56 | 51 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_ringbuffer_pool` | 39.65 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_ringbuffer_nopool` | 54.71 | 51 | 0 |
| `BenchmarkMixedWorkload` | `ladder_linkedlist_pool` | 28.58 | 0 | 0 |
| `BenchmarkMixedWorkload` | `ladder_linkedlist_nopool` | 41.71 | 51 | 0 |
| `BenchmarkMixedWorkload` | `ladder_ringbuffer_pool` | 31.03 | 0 | 0 |
| `BenchmarkMixedWorkload` | `ladder_ringbuffer_nopool` | 45.75 | 51 | 0 |
| `BenchmarkEngineChannel` | `default` | 80.47 | 0 | 0 |
| `BenchmarkEngineBatch` | `default` | 21.66 | 0 | 0 |
| `BenchmarkEngineDisruptor` | `default` | 51.03 | 0 | 0 |
| `BenchmarkFastBook` | `AddLimit` | 15.63 | 0 | 0 |
| `BenchmarkFastBook` | `Cancel` | 14.84 | 0 | 0 |
| `BenchmarkFastBook` | `MatchSingle` | 33.28 | 0 | 0 |
| `BenchmarkFastBook` | `MatchSweep` | 202.40 | 0 | 0 |
| `BenchmarkFastBook` | `Market` | 207.20 | 0 | 0 |
| `BenchmarkFastBook` | `BestBidAsk` | 0.92 | 0 | 0 |
| `BenchmarkFastBook` | `Depth` | 22.71 | 0 | 0 |
| `BenchmarkFastBook` | `MixedWorkload` | 25.74 | 0 | 0 |
| `BenchmarkArenaBook` | `AddLimit` | 19.36 | 0 | 0 |
| `BenchmarkArenaBook` | `Cancel` | 13.38 | 0 | 0 |
| `BenchmarkArenaBook` | `MatchSingle` | 31.82 | 0 | 0 |
| `BenchmarkArenaBook` | `MatchSweep` | 200.50 | 0 | 0 |
| `BenchmarkArenaBook` | `Market` | 200.70 | 0 | 0 |
| `BenchmarkArenaBook` | `BestBidAsk` | 1.06 | 0 | 0 |
| `BenchmarkArenaBook` | `Depth` | 24.68 | 0 | 0 |
| `BenchmarkArenaBook` | `MixedWorkload` | 25.33 | 0 | 0 |

## HDR Histogram Latency (Mixed Workload)

| Variant | Throughput (ops/s) | Avg (ns) | p50 (ns) | p90 (ns) | p99 (ns) | p99.9 (ns) | Max (ns) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `sortedslice_linkedlist_pool` | 19767870 | 50.6 | 56 | 88 | 138 | 203 | 235775 |
| `sortedslice_linkedlist_nopool` | 15811393 | 63.2 | 67 | 109 | 198 | 479 | 212991 |
| `sortedslice_ringbuffer_pool` | 18082630 | 55.3 | 54 | 82 | 143 | 853 | 233855 |
| `sortedslice_ringbuffer_nopool` | 16560259 | 60.4 | 64 | 102 | 229 | 1840 | 129983 |
| `btree_linkedlist_pool` | 13814745 | 72.4 | 63 | 99 | 168 | 261 | 128063 |
| `btree_linkedlist_nopool` | 14088533 | 71.0 | 74 | 127 | 247 | 1045 | 127935 |
| `btree_ringbuffer_pool` | 16663345 | 60.0 | 70 | 128 | 241 | 1175 | 237951 |
| `btree_ringbuffer_nopool` | 12999557 | 76.9 | 73 | 117 | 269 | 1970 | 129471 |
| `rbtree_linkedlist_pool` | 16609018 | 60.2 | 52 | 82 | 140 | 233 | 134527 |
| `rbtree_linkedlist_nopool` | 13940376 | 71.7 | 65 | 108 | 206 | 1468 | 126911 |
| `rbtree_ringbuffer_pool` | 19453141 | 51.4 | 55 | 92 | 179 | 1027 | 248703 |
| `rbtree_ringbuffer_nopool` | 16144141 | 61.9 | 61 | 98 | 240 | 1783 | 194687 |
| `ladder_linkedlist_pool` | 17743234 | 56.4 | 42 | 67 | 115 | 169 | 233087 |
| `ladder_linkedlist_nopool` | 18485221 | 54.1 | 51 | 80 | 157 | 1408 | 435199 |
| `ladder_ringbuffer_pool` | 25642877 | 39.0 | 45 | 73 | 141 | 826 | 131327 |
| `ladder_ringbuffer_nopool` | 19393665 | 51.6 | 68 | 104 | 258 | 2273 | 1040383 |

## Profiling Hotspots (pprof — 147 profiled benchmarks)

| Benchmark | Top CPU Hotspot | Top Mem Alloc Space |
|---|---|---|
| `BenchmarkAddLimit/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (30.4% flat) | `runtime.mallocgc` (4618.88kB, 32.5%) |
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (13.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (569.55MB, 98.0%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (32.9% flat) | `runtime.mallocgc` (4620.26kB, 45.6%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (20.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (519.05MB, 97.7%) |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.4% flat) | `runtime.mallocgc` (5643.42kB, 47.0%) |
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `github.com/google/btree.items[go.shape.*uint8].find` (12.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (443.54MB, 97.3%) |
| `BenchmarkAddLimit/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.8% flat) | `runtime.mallocgc` (7191.10kB, 46.8%) |
| `BenchmarkAddLimit/btree_ringbuffer_nopool` | `github.com/google/btree.items[go.shape.*uint8].find` (9.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (476.04MB, 97.9%) |
| `BenchmarkAddLimit/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (22.7% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4956.07kB, 27.5%) |
| `BenchmarkAddLimit/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (25.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (567.55MB, 97.9%) |
| `BenchmarkAddLimit/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (28.4% flat) | `runtime.mallocgc` (4626.36kB, 29.6%) |
| `BenchmarkAddLimit/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (21.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (573.55MB, 98.2%) |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (32.3% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 27.4%) |
| `BenchmarkAddLimit/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (20.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (828.58MB, 97.6%) |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (33.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 27.3%) |
| `BenchmarkAddLimit/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (13.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (808.57MB, 98.2%) |
| `BenchmarkCancel/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (35.4% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 35.7%) |
| `BenchmarkCancel/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (21.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.29GB, 99.0%) |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (26.1% flat) | `runtime.mallocgc` (3594.25kB, 34.0%) |
| `BenchmarkCancel/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (23.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.52GB, 99.4%) |
| `BenchmarkCancel/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (12.6% flat) | `runtime.mallocgc` (6157.37kB, 36.5%) |
| `BenchmarkCancel/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.19GB, 99.0%) |
| `BenchmarkCancel/btree_ringbuffer_pool` | `github.com/google/btree.items[go.shape.*uint8].find` (12.3% flat) | `runtime.mallocgc` (4620.26kB, 38.1%) |
| `BenchmarkCancel/btree_ringbuffer_nopool` | `sort.Search (inline)` (9.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.35GB, 99.2%) |
| `BenchmarkCancel/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (31.0% flat) | `runtime.mallocgc` (4620.26kB, 30.7%) |
| `BenchmarkCancel/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (14.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1249.11MB, 98.7%) |
| `BenchmarkCancel/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (22.9% flat) | `runtime.mallocgc` (3602.25kB, 32.5%) |
| `BenchmarkCancel/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (15.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.42GB, 99.3%) |
| `BenchmarkCancel/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (29.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 38.7%) |
| `BenchmarkCancel/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (15.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (3.25GB, 99.4%) |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (28.6% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 35.2%) |
| `BenchmarkCancel/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (18.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1548.64MB, 98.8%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (56.9% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4248.06kB, 40.6%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (42.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (895.08MB, 99.1%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (62.7% flat) | `runtime.mallocgc` (3082.21kB, 29.2%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (40.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (982.09MB, 99.1%) |
| `BenchmarkMatchSingle/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (62.7% flat) | `runtime.mallocgc` (5642.43kB, 44.6%) |
| `BenchmarkMatchSingle/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (898.08MB, 99.3%) |
| `BenchmarkMatchSingle/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (68.4% flat) | `runtime.mallocgc` (3603.15kB, 29.3%) |
| `BenchmarkMatchSingle/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (41.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (950.59MB, 99.0%) |
| `BenchmarkMatchSingle/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (67.6% flat) | `runtime.mallocgc` (3595.21kB, 34.6%) |
| `BenchmarkMatchSingle/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (892.08MB, 98.8%) |
| `BenchmarkMatchSingle/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (56.8% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4956.07kB, 35.2%) |
| `BenchmarkMatchSingle/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (38.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (990.09MB, 98.5%) |
| `BenchmarkMatchSingle/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (72.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 36.0%) |
| `BenchmarkMatchSingle/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (958.59MB, 98.4%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (60.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (5001.69kB, 28.9%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (33.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (941.59MB, 98.3%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_pool` | `runtime.memmove` (60.5% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4248.06kB, 25.5%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_nopool` | `runtime.memmove` (51.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.41GB, 99.5%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_pool` | `runtime.memmove` (61.4% flat) | `runtime.mallocgc` (5642.61kB, 33.8%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_nopool` | `runtime.memmove` (48.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.33GB, 99.4%) |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (18.1% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (118.06MB, 55.0%) |
| `BenchmarkMatchSweep/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (627.56MB, 83.6%) |
| `BenchmarkMatchSweep/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (13.8% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (68.53MB, 51.9%) |
| `BenchmarkMatchSweep/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (10.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (449.54MB, 84.3%) |
| `BenchmarkMatchSweep/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (38.7% flat) | `runtime.mallocgc` (5132.32kB, 32.1%) |
| `BenchmarkMatchSweep/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.85GB, 99.5%) |
| `BenchmarkMatchSweep/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.5% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (5664.08kB, 26.4%) |
| `BenchmarkMatchSweep/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.85GB, 99.4%) |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (55.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 32.8%) |
| `BenchmarkMatchSweep/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.18GB, 99.2%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (61.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 32.9%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (32.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.11GB, 99.0%) |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `runtime.memmove` (65.2% flat) | `runtime.mallocgc` (5643.93kB, 29.3%) |
| `BenchmarkMarket/sortedslice_linkedlist_nopool` | `runtime.memmove` (46.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.43GB, 99.5%) |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `runtime.memmove` (61.3% flat) | `runtime.mallocgc` (5139.36kB, 25.4%) |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `runtime.memmove` (46.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.31GB, 99.3%) |
| `BenchmarkMarket/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (18.4% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (69.03MB, 59.4%) |
| `BenchmarkMarket/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (11.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1060.60MB, 84.8%) |
| `BenchmarkMarket/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (18.4% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (116.06MB, 57.9%) |
| `BenchmarkMarket/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (496.55MB, 83.1%) |
| `BenchmarkMarket/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.8% flat) | `runtime.mallocgc` (4106.31kB, 21.4%) |
| `BenchmarkMarket/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.81GB, 99.1%) |
| `BenchmarkMarket/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (35.9% flat) | `runtime.mallocgc` (5644.57kB, 23.7%) |
| `BenchmarkMarket/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (24.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.87GB, 99.3%) |
| `BenchmarkMarket/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (51.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (10003.38kB, 44.1%) |
| `BenchmarkMarket/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (34.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.25GB, 99.3%) |
| `BenchmarkMarket/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (62.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 34.6%) |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (2.09GB, 99.0%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/bench_test.BenchmarkBestBidAsk.func1` (37.0% flat) | `runtime.mallocgc` (5146.40kB, 48.1%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (40.3% flat) | `runtime.mallocgc` (6164.41kB, 60.2%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (31.1% flat) | `runtime.mallocgc` (4124.18kB, 48.8%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (35.5% flat) | `runtime.mallocgc` (5131.37kB, 63.4%) |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (63.6% flat) | `runtime.mallocgc` (4114.25kB, 40.6%) |
| `BenchmarkBestBidAsk/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (68.2% flat) | `runtime.mallocgc` (4106.31kB, 46.9%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (59.4% flat) | `runtime.mallocgc` (5131.50kB, 62.0%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (57.1% flat) | `runtime.mallocgc` (2570.14kB, 47.9%) |
| `BenchmarkBestBidAsk/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).Best` (39.5% flat) | `runtime.mallocgc` (4626.54kB, 49.4%) |
| `BenchmarkBestBidAsk/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (48.0% flat) | `runtime.mallocgc` (4107.44kB, 59.5%) |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (50.0% flat) | `runtime.mallocgc` (5138.59kB, 58.9%) |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (45.3% flat) | `runtime.mallocgc` (3603.19kB, 39.8%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (40.3% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (10003.38kB, 46.6%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (52.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 50.7%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (54.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7515.68kB, 44.3%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (48.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (11003.72kB, 56.6%) |
| `BenchmarkDepth/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (91.5% flat) | `runtime.mallocgc` (4107.26kB, 42.7%) |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (94.9% flat) | `runtime.mallocgc` (4620.26kB, 65.3%) |
| `BenchmarkDepth/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (89.8% flat) | `runtime.mallocgc` (2570.15kB, 38.2%) |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (87.8% flat) | `runtime.mallocgc` (3082.21kB, 35.1%) |
| `BenchmarkDepth/btree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (57.7% flat) | `runtime.mallocgc` (4107.26kB, 41.9%) |
| `BenchmarkDepth/btree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (61.7% flat) | `runtime.mallocgc` (3595.21kB, 75.5%) |
| `BenchmarkDepth/btree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (54.5% flat) | `runtime.mallocgc` (3090.19kB, 34.5%) |
| `BenchmarkDepth/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (55.8% flat) | `runtime.mallocgc` (2570.10kB, 31.9%) |
| `BenchmarkDepth/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (50.0% flat) | `runtime.mallocgc` (5131.37kB, 48.3%) |
| `BenchmarkDepth/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (48.8% flat) | `runtime.mallocgc` (6162.88kB, 61.7%) |
| `BenchmarkDepth/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (45.5% flat) | `runtime.mallocgc` (3595.21kB, 38.3%) |
| `BenchmarkDepth/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (55.6% flat) | `runtime.mallocgc` (4620.26kB, 58.2%) |
| `BenchmarkDepth/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (88.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 48.8%) |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (86.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 49.2%) |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (86.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 52.8%) |
| `BenchmarkDepth/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (86.6% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 51.1%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.9% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (14160.19kB, 40.7%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (440.04MB, 94.7%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (24.6% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (20532.27kB, 44.1%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (18.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (434.04MB, 94.3%) |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (14.6% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16284.22kB, 47.4%) |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (386.04MB, 94.2%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (24.1% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16992.23kB, 36.2%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (375.53MB, 92.5%) |
| `BenchmarkMixedWorkload/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (24.5% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16284.22kB, 42.5%) |
| `BenchmarkMixedWorkload/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (447.54MB, 95.7%) |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.2% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (14160.19kB, 40.2%) |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (445.54MB, 94.1%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.2% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (19116.26kB, 41.2%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (615.06MB, 95.5%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (46.1% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16284.22kB, 33.6%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (974.59MB, 96.4%) |
| `BenchmarkEngineChannel` | `runtime.procyieldAsm` (23.5% flat) | `github.com/nedo/orbook/engine.New (inline)` (79.73MB, 62.4%) |
| `BenchmarkEngineBatch` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.3% flat) | `github.com/nedo/orbook/bench_test.BenchmarkEngineBatch` (70.63MB, 58.7%) |
| `BenchmarkEngineDisruptor` | `github.com/nedo/orbook/engine.(*Disruptor).process` (19.8% flat) | `github.com/nedo/orbook/engine.NewDisruptor` (99.85MB, 69.3%) |
| `BenchmarkFastBook/AddLimit` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (37.6% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (7002.37kB, 34.5%) |
| `BenchmarkFastBook/Cancel` | `github.com/nedo/orbook/book.(*FastBook).CancelOrder` (39.6% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (8002.71kB, 31.9%) |
| `BenchmarkFastBook/MatchSingle` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (68.9% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (7002.37kB, 35.2%) |
| `BenchmarkFastBook/MatchSweep` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (72.7% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (8002.71kB, 33.6%) |
| `BenchmarkFastBook/Market` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (79.7% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (6002.03kB, 32.9%) |
| `BenchmarkFastBook/BestBidAsk` | `github.com/nedo/orbook/book.(*FastBook).BestAsk (inline)` (45.3% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (10003.38kB, 49.1%) |
| `BenchmarkFastBook/Depth` | `github.com/nedo/orbook/book.(*fastLadder).topN` (89.1% flat) | `github.com/nedo/orbook/book.newFastLadder` (6002.03kB, 35.8%) |
| `BenchmarkFastBook/MixedWorkload` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (35.2% flat) | `github.com/nedo/orbook/book.(*FastBook).growPool` (13452.18kB, 32.2%) |
| `BenchmarkArenaBook/AddLimit` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (43.2% flat) | `runtime.mallocgc` (4619.31kB, 27.1%) |
| `BenchmarkArenaBook/Cancel` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (35.5% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6753.23kB, 39.1%) |
| `BenchmarkArenaBook/MatchSingle` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (74.7% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 42.8%) |
| `BenchmarkArenaBook/MatchSweep` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (73.1% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 43.3%) |
| `BenchmarkArenaBook/Market` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (71.9% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (6239.92kB, 37.3%) |
| `BenchmarkArenaBook/BestBidAsk` | `github.com/nedo/orbook/book.(*ArenaBook).BestAsk (inline)` (59.3% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (8022.75kB, 44.9%) |
| `BenchmarkArenaBook/Depth` | `github.com/nedo/orbook/book.(*arenaLadder).topN` (80.0% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (7131.34kB, 38.9%) |
| `BenchmarkArenaBook/MixedWorkload` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (42.3% flat) | `github.com/nedo/orbook/book.NewArenaBook` (11037.44kB, 24.6%) |

*(Full top-15 CPU & memory tables for all profiled variants are in `latest.json` and `latest.html`.)*

