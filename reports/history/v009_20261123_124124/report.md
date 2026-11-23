# Orbook Benchmark & Profile Report — `v009_20261123_124124`

- **Version:** `v009` | **Mode:** `full`
- **Time:** `2026-11-23 12:41:24 WIB` | **Commit:** `4e16c62` | **Go:** `go1.27.3-X:nodwarf5` | **CPU:** `Intel(R) Core(TM) i9-14900HX`

## Diff vs Previous (`v008_20261123_041545`, threshold ±2.0%)

| Benchmark | Status | Old (ns/op) | New (ns/op) | Delta | Allocs |
|---|---|---:|---:|---:|---|
| `BenchmarkAddLimit/sortedslice_linkedlist_pool` | `faster` | 59.29 | 47.12 | -20.53% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `faster` | 95.32 | 83.80 | -12.09% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/sortedslice_ringbuffer_pool` | `faster` | 59.13 | 48.10 | -18.65% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `faster` | 104.20 | 78.56 | -24.61% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `faster` | 79.39 | 68.19 | -14.11% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `faster` | 114.70 | 95.58 | -16.67% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/btree_ringbuffer_pool` | `faster` | 82.56 | 66.33 | -19.66% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/btree_ringbuffer_nopool` | `faster` | 116.10 | 102.70 | -11.54% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/rbtree_linkedlist_pool` | `faster` | 58.47 | 41.00 | -29.88% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/rbtree_linkedlist_nopool` | `faster` | 92.91 | 68.06 | -26.75% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/rbtree_ringbuffer_pool` | `faster` | 61.77 | 43.10 | -30.23% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/rbtree_ringbuffer_nopool` | `faster` | 99.58 | 76.15 | -23.53% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `faster` | 29.86 | 20.37 | -31.78% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_linkedlist_nopool` | `faster` | 59.79 | 46.71 | -21.88% | 96→96 B, 1→1 allocs |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `faster` | 33.22 | 22.98 | -30.82% | 0→0 B, 0→0 allocs |
| `BenchmarkAddLimit/ladder_ringbuffer_nopool` | `faster` | 60.21 | 49.62 | -17.59% | 96→96 B, 1→1 allocs |
| `BenchmarkCancel/sortedslice_linkedlist_pool` | `faster` | 39.02 | 20.08 | -48.54% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/sortedslice_linkedlist_nopool` | `faster` | 42.14 | 34.42 | -18.32% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `faster` | 38.01 | 20.37 | -46.41% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/sortedslice_ringbuffer_nopool` | `faster` | 49.72 | 28.82 | -42.04% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_linkedlist_pool` | `faster` | 39.86 | 21.42 | -46.26% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_linkedlist_nopool` | `faster` | 42.70 | 36.87 | -13.65% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_ringbuffer_pool` | `faster` | 39.67 | 20.15 | -49.21% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/btree_ringbuffer_nopool` | `faster` | 54.27 | 28.13 | -48.17% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_linkedlist_pool` | `faster` | 40.36 | 19.78 | -50.99% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_linkedlist_nopool` | `faster` | 47.97 | 34.63 | -27.81% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_ringbuffer_pool` | `faster` | 40.55 | 19.92 | -50.88% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/rbtree_ringbuffer_nopool` | `faster` | 61.20 | 25.62 | -58.14% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_linkedlist_pool` | `faster` | 39.23 | 17.87 | -54.45% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_linkedlist_nopool` | `faster` | 44.52 | 25.36 | -43.04% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `faster` | 37.73 | 17.79 | -52.85% | 0→0 B, 0→0 allocs |
| `BenchmarkCancel/ladder_ringbuffer_nopool` | `faster` | 45.35 | 26.10 | -42.45% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_linkedlist_pool` | `faster` | 56.22 | 37.83 | -32.71% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_linkedlist_nopool` | `faster` | 64.13 | 43.01 | -32.93% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_pool` | `faster` | 60.57 | 35.38 | -41.59% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `faster` | 65.04 | 38.99 | -40.05% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_linkedlist_pool` | `faster` | 60.43 | 35.05 | -42.00% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_linkedlist_nopool` | `faster` | 64.05 | 44.28 | -30.87% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_ringbuffer_pool` | `faster` | 59.31 | 38.77 | -34.63% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/btree_ringbuffer_nopool` | `faster` | 64.78 | 41.47 | -35.98% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/rbtree_linkedlist_pool` | `faster` | 63.34 | 35.74 | -43.57% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/rbtree_linkedlist_nopool` | `faster` | 72.29 | 43.21 | -40.23% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/rbtree_ringbuffer_pool` | `faster` | 64.37 | 38.16 | -40.72% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/rbtree_ringbuffer_nopool` | `faster` | 69.30 | 40.26 | -41.90% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_linkedlist_pool` | `faster` | 60.25 | 36.85 | -38.84% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_linkedlist_nopool` | `faster` | 63.95 | 45.58 | -28.73% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_ringbuffer_pool` | `faster` | 60.18 | 36.55 | -39.27% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSingle/ladder_ringbuffer_nopool` | `faster` | 66.26 | 42.51 | -35.84% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_linkedlist_pool` | `faster` | 365.60 | 258.60 | -29.27% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_linkedlist_nopool` | `faster` | 375.40 | 248.20 | -33.88% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_pool` | `faster` | 384.60 | 252.60 | -34.32% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_nopool` | `faster` | 420.80 | 281.50 | -33.10% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `faster` | 795.60 | 642.40 | -19.26% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_linkedlist_nopool` | `faster` | 827.50 | 617.40 | -25.39% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_ringbuffer_pool` | `faster` | 836.60 | 632.40 | -24.41% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/btree_ringbuffer_nopool` | `faster` | 1078.00 | 759.10 | -29.58% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_linkedlist_pool` | `faster` | 503.30 | 317.10 | -37.00% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_linkedlist_nopool` | `faster` | 531.50 | 345.80 | -34.94% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_ringbuffer_pool` | `faster` | 518.10 | 334.40 | -35.46% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/rbtree_ringbuffer_nopool` | `faster` | 626.00 | 420.80 | -32.78% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `faster` | 368.10 | 246.70 | -32.98% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_linkedlist_nopool` | `faster` | 419.60 | 295.20 | -29.65% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_ringbuffer_pool` | `faster` | 385.90 | 261.00 | -32.37% | 0→0 B, 0→0 allocs |
| `BenchmarkMatchSweep/ladder_ringbuffer_nopool` | `faster` | 462.70 | 325.90 | -29.57% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `faster` | 356.40 | 247.80 | -30.47% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_linkedlist_nopool` | `faster` | 367.20 | 260.00 | -29.19% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `faster` | 364.70 | 262.00 | -28.16% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `faster` | 429.30 | 276.60 | -35.57% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_linkedlist_pool` | `faster` | 781.30 | 618.00 | -20.90% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_linkedlist_nopool` | `faster` | 787.20 | 620.30 | -21.20% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_ringbuffer_pool` | `faster` | 822.40 | 618.60 | -24.78% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/btree_ringbuffer_nopool` | `faster` | 1038.00 | 800.10 | -22.92% | 0→1 B, 0→0 allocs |
| `BenchmarkMarket/rbtree_linkedlist_pool` | `faster` | 512.80 | 300.40 | -41.42% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/rbtree_linkedlist_nopool` | `faster` | 529.30 | 354.60 | -33.01% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/rbtree_ringbuffer_pool` | `faster` | 526.70 | 319.90 | -39.26% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/rbtree_ringbuffer_nopool` | `faster` | 639.40 | 410.90 | -35.74% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_linkedlist_pool` | `faster` | 362.30 | 239.80 | -33.81% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_linkedlist_nopool` | `faster` | 417.10 | 311.10 | -25.41% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_ringbuffer_pool` | `faster` | 371.60 | 250.70 | -32.53% | 0→0 B, 0→0 allocs |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `faster` | 464.30 | 320.50 | -30.97% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `slower` | 3.24 | 3.39 | +4.66% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `slower` | 3.13 | 3.31 | +5.42% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `slower` | 2.17 | 2.22 | +2.44% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_linkedlist_nopool` | `slower` | 2.14 | 2.19 | +2.20% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `slower` | 2.14 | 2.24 | +4.72% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `slower` | 2.19 | 2.29 | +4.56% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_linkedlist_pool` | `faster` | 4.15 | 3.38 | -18.53% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_linkedlist_nopool` | `faster` | 4.08 | 3.24 | -20.57% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_pool` | `faster` | 3.98 | 3.37 | -15.42% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_nopool` | `faster` | 3.90 | 3.36 | -13.83% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `slower` | 2.31 | 3.12 | +35.07% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `slower` | 2.36 | 3.30 | +39.80% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_ringbuffer_pool` | `slower` | 2.35 | 3.14 | +33.47% | 0→0 B, 0→0 allocs |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `slower` | 2.29 | 3.15 | +37.20% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_linkedlist_pool` | `slower` | 16.02 | 16.77 | +4.68% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `slower` | 16.42 | 18.28 | +11.33% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_ringbuffer_pool` | `slower` | 15.92 | 17.37 | +9.11% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `slower` | 16.31 | 16.85 | +3.31% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_ringbuffer_pool` | `faster` | 61.68 | 58.62 | -4.96% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/btree_ringbuffer_nopool` | `slower` | 57.49 | 59.22 | +3.01% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/rbtree_linkedlist_pool` | `faster` | 39.26 | 38.46 | -2.04% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/rbtree_linkedlist_nopool` | `faster` | 41.93 | 38.73 | -7.63% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/rbtree_ringbuffer_pool` | `faster` | 42.04 | 40.32 | -4.09% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/rbtree_ringbuffer_nopool` | `faster` | 42.09 | 40.00 | -4.97% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `faster` | 22.64 | 21.53 | -4.90% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `slower` | 21.40 | 22.02 | +2.90% | 0→0 B, 0→0 allocs |
| `BenchmarkDepth/ladder_ringbuffer_nopool` | `faster` | 21.47 | 20.87 | -2.79% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `faster` | 54.65 | 40.33 | -26.20% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_nopool` | `faster` | 67.70 | 57.92 | -14.45% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `faster` | 54.08 | 43.09 | -20.32% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_nopool` | `faster` | 69.68 | 55.69 | -20.08% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `faster` | 68.16 | 54.80 | -19.60% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `faster` | 78.60 | 75.08 | -4.48% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_ringbuffer_pool` | `faster` | 67.02 | 56.45 | -15.77% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/btree_ringbuffer_nopool` | `faster` | 81.59 | 68.60 | -15.92% | 51→52 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_linkedlist_pool` | `faster` | 62.53 | 39.76 | -36.41% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_linkedlist_nopool` | `faster` | 75.51 | 55.64 | -26.31% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_pool` | `faster` | 64.79 | 41.03 | -36.67% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_nopool` | `faster` | 79.04 | 59.19 | -25.11% | 51→52 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `faster` | 44.44 | 30.68 | -30.96% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_linkedlist_nopool` | `faster` | 58.20 | 41.23 | -29.16% | 51→51 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_ringbuffer_pool` | `faster` | 43.61 | 31.18 | -28.50% | 0→0 B, 0→0 allocs |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `faster` | 60.98 | 48.87 | -19.86% | 51→51 B, 0→0 allocs |
| `BenchmarkEngineChannel` | `slower` | 78.91 | 84.01 | +6.46% | 0→0 B, 0→0 allocs |
| `BenchmarkEngineDisruptor` | `slower` | 51.98 | 53.74 | +3.39% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/MatchSweep` | `faster` | 213.10 | 207.50 | -2.63% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/BestBidAsk` | `slower` | 0.94 | 0.98 | +4.21% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/Depth` | `faster` | 22.87 | 22.34 | -2.32% | 0→0 B, 0→0 allocs |
| `BenchmarkFastBook/MixedWorkload` | `slower` | 24.57 | 26.09 | +6.19% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/Cancel` | `faster` | 13.82 | 13.32 | -3.62% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/MatchSingle` | `faster` | 35.12 | 32.57 | -7.26% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/MatchSweep` | `faster` | 212.80 | 201.20 | -5.45% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/Market` | `faster` | 203.10 | 195.80 | -3.59% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/BestBidAsk` | `slower` | 1.16 | 1.19 | +2.93% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/Depth` | `faster` | 27.54 | 26.80 | -2.69% | 0→0 B, 0→0 allocs |
| `BenchmarkArenaBook/MixedWorkload` | `faster` | 27.24 | 25.46 | -6.53% | 0→0 B, 0→0 allocs |

## Benchmarks

| Method | Variant | ns/op | B/op | allocs/op |
|---|---|---:|---:|---:|
| `BenchmarkAddLimit` | `sortedslice_linkedlist_pool` | 47.12 | 0 | 0 |
| `BenchmarkAddLimit` | `sortedslice_linkedlist_nopool` | 83.80 | 96 | 1 |
| `BenchmarkAddLimit` | `sortedslice_ringbuffer_pool` | 48.10 | 0 | 0 |
| `BenchmarkAddLimit` | `sortedslice_ringbuffer_nopool` | 78.56 | 96 | 1 |
| `BenchmarkAddLimit` | `btree_linkedlist_pool` | 68.19 | 0 | 0 |
| `BenchmarkAddLimit` | `btree_linkedlist_nopool` | 95.58 | 96 | 1 |
| `BenchmarkAddLimit` | `btree_ringbuffer_pool` | 66.33 | 0 | 0 |
| `BenchmarkAddLimit` | `btree_ringbuffer_nopool` | 102.70 | 96 | 1 |
| `BenchmarkAddLimit` | `rbtree_linkedlist_pool` | 41.00 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_linkedlist_nopool` | 68.06 | 96 | 1 |
| `BenchmarkAddLimit` | `rbtree_ringbuffer_pool` | 43.10 | 0 | 0 |
| `BenchmarkAddLimit` | `rbtree_ringbuffer_nopool` | 76.15 | 96 | 1 |
| `BenchmarkAddLimit` | `ladder_linkedlist_pool` | 20.37 | 0 | 0 |
| `BenchmarkAddLimit` | `ladder_linkedlist_nopool` | 46.71 | 96 | 1 |
| `BenchmarkAddLimit` | `ladder_ringbuffer_pool` | 22.98 | 0 | 0 |
| `BenchmarkAddLimit` | `ladder_ringbuffer_nopool` | 49.62 | 96 | 1 |
| `BenchmarkCancel` | `sortedslice_linkedlist_pool` | 20.08 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_linkedlist_nopool` | 34.42 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_ringbuffer_pool` | 20.37 | 0 | 0 |
| `BenchmarkCancel` | `sortedslice_ringbuffer_nopool` | 28.82 | 0 | 0 |
| `BenchmarkCancel` | `btree_linkedlist_pool` | 21.42 | 0 | 0 |
| `BenchmarkCancel` | `btree_linkedlist_nopool` | 36.87 | 0 | 0 |
| `BenchmarkCancel` | `btree_ringbuffer_pool` | 20.15 | 0 | 0 |
| `BenchmarkCancel` | `btree_ringbuffer_nopool` | 28.13 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_linkedlist_pool` | 19.78 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_linkedlist_nopool` | 34.63 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_ringbuffer_pool` | 19.92 | 0 | 0 |
| `BenchmarkCancel` | `rbtree_ringbuffer_nopool` | 25.62 | 0 | 0 |
| `BenchmarkCancel` | `ladder_linkedlist_pool` | 17.87 | 0 | 0 |
| `BenchmarkCancel` | `ladder_linkedlist_nopool` | 25.36 | 0 | 0 |
| `BenchmarkCancel` | `ladder_ringbuffer_pool` | 17.79 | 0 | 0 |
| `BenchmarkCancel` | `ladder_ringbuffer_nopool` | 26.10 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_linkedlist_pool` | 37.83 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_linkedlist_nopool` | 43.01 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_ringbuffer_pool` | 35.38 | 0 | 0 |
| `BenchmarkMatchSingle` | `sortedslice_ringbuffer_nopool` | 38.99 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_linkedlist_pool` | 35.05 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_linkedlist_nopool` | 44.28 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_ringbuffer_pool` | 38.77 | 0 | 0 |
| `BenchmarkMatchSingle` | `btree_ringbuffer_nopool` | 41.47 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_linkedlist_pool` | 35.74 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_linkedlist_nopool` | 43.21 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_ringbuffer_pool` | 38.16 | 0 | 0 |
| `BenchmarkMatchSingle` | `rbtree_ringbuffer_nopool` | 40.26 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_linkedlist_pool` | 36.85 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_linkedlist_nopool` | 45.58 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_ringbuffer_pool` | 36.55 | 0 | 0 |
| `BenchmarkMatchSingle` | `ladder_ringbuffer_nopool` | 42.51 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_linkedlist_pool` | 258.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_linkedlist_nopool` | 248.20 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_ringbuffer_pool` | 252.60 | 0 | 0 |
| `BenchmarkMatchSweep` | `sortedslice_ringbuffer_nopool` | 281.50 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_linkedlist_pool` | 642.40 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_linkedlist_nopool` | 617.40 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_ringbuffer_pool` | 632.40 | 0 | 0 |
| `BenchmarkMatchSweep` | `btree_ringbuffer_nopool` | 759.10 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_linkedlist_pool` | 317.10 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_linkedlist_nopool` | 345.80 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_ringbuffer_pool` | 334.40 | 0 | 0 |
| `BenchmarkMatchSweep` | `rbtree_ringbuffer_nopool` | 420.80 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_linkedlist_pool` | 246.70 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_linkedlist_nopool` | 295.20 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_ringbuffer_pool` | 261.00 | 0 | 0 |
| `BenchmarkMatchSweep` | `ladder_ringbuffer_nopool` | 325.90 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_linkedlist_pool` | 247.80 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_linkedlist_nopool` | 260.00 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_ringbuffer_pool` | 262.00 | 0 | 0 |
| `BenchmarkMarket` | `sortedslice_ringbuffer_nopool` | 276.60 | 0 | 0 |
| `BenchmarkMarket` | `btree_linkedlist_pool` | 618.00 | 0 | 0 |
| `BenchmarkMarket` | `btree_linkedlist_nopool` | 620.30 | 0 | 0 |
| `BenchmarkMarket` | `btree_ringbuffer_pool` | 618.60 | 0 | 0 |
| `BenchmarkMarket` | `btree_ringbuffer_nopool` | 800.10 | 1 | 0 |
| `BenchmarkMarket` | `rbtree_linkedlist_pool` | 300.40 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_linkedlist_nopool` | 354.60 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_ringbuffer_pool` | 319.90 | 0 | 0 |
| `BenchmarkMarket` | `rbtree_ringbuffer_nopool` | 410.90 | 0 | 0 |
| `BenchmarkMarket` | `ladder_linkedlist_pool` | 239.80 | 0 | 0 |
| `BenchmarkMarket` | `ladder_linkedlist_nopool` | 311.10 | 0 | 0 |
| `BenchmarkMarket` | `ladder_ringbuffer_pool` | 250.70 | 0 | 0 |
| `BenchmarkMarket` | `ladder_ringbuffer_nopool` | 320.50 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_linkedlist_pool` | 3.14 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_linkedlist_nopool` | 3.39 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_ringbuffer_pool` | 3.31 | 0 | 0 |
| `BenchmarkBestBidAsk` | `sortedslice_ringbuffer_nopool` | 3.16 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_linkedlist_pool` | 2.22 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_linkedlist_nopool` | 2.19 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_ringbuffer_pool` | 2.24 | 0 | 0 |
| `BenchmarkBestBidAsk` | `btree_ringbuffer_nopool` | 2.29 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_linkedlist_pool` | 3.38 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_linkedlist_nopool` | 3.24 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_ringbuffer_pool` | 3.37 | 0 | 0 |
| `BenchmarkBestBidAsk` | `rbtree_ringbuffer_nopool` | 3.36 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_linkedlist_pool` | 3.12 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_linkedlist_nopool` | 3.30 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_ringbuffer_pool` | 3.14 | 0 | 0 |
| `BenchmarkBestBidAsk` | `ladder_ringbuffer_nopool` | 3.15 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_linkedlist_pool` | 16.77 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_linkedlist_nopool` | 18.28 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_ringbuffer_pool` | 17.37 | 0 | 0 |
| `BenchmarkDepth` | `sortedslice_ringbuffer_nopool` | 16.85 | 0 | 0 |
| `BenchmarkDepth` | `btree_linkedlist_pool` | 60.94 | 0 | 0 |
| `BenchmarkDepth` | `btree_linkedlist_nopool` | 59.20 | 0 | 0 |
| `BenchmarkDepth` | `btree_ringbuffer_pool` | 58.62 | 0 | 0 |
| `BenchmarkDepth` | `btree_ringbuffer_nopool` | 59.22 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_linkedlist_pool` | 38.46 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_linkedlist_nopool` | 38.73 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_ringbuffer_pool` | 40.32 | 0 | 0 |
| `BenchmarkDepth` | `rbtree_ringbuffer_nopool` | 40.00 | 0 | 0 |
| `BenchmarkDepth` | `ladder_linkedlist_pool` | 22.37 | 0 | 0 |
| `BenchmarkDepth` | `ladder_linkedlist_nopool` | 21.53 | 0 | 0 |
| `BenchmarkDepth` | `ladder_ringbuffer_pool` | 22.02 | 0 | 0 |
| `BenchmarkDepth` | `ladder_ringbuffer_nopool` | 20.87 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_linkedlist_pool` | 40.33 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_linkedlist_nopool` | 57.92 | 51 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_ringbuffer_pool` | 43.09 | 0 | 0 |
| `BenchmarkMixedWorkload` | `sortedslice_ringbuffer_nopool` | 55.69 | 51 | 0 |
| `BenchmarkMixedWorkload` | `btree_linkedlist_pool` | 54.80 | 0 | 0 |
| `BenchmarkMixedWorkload` | `btree_linkedlist_nopool` | 75.08 | 51 | 0 |
| `BenchmarkMixedWorkload` | `btree_ringbuffer_pool` | 56.45 | 0 | 0 |
| `BenchmarkMixedWorkload` | `btree_ringbuffer_nopool` | 68.60 | 52 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_linkedlist_pool` | 39.76 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_linkedlist_nopool` | 55.64 | 51 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_ringbuffer_pool` | 41.03 | 0 | 0 |
| `BenchmarkMixedWorkload` | `rbtree_ringbuffer_nopool` | 59.19 | 52 | 0 |
| `BenchmarkMixedWorkload` | `ladder_linkedlist_pool` | 30.68 | 0 | 0 |
| `BenchmarkMixedWorkload` | `ladder_linkedlist_nopool` | 41.23 | 51 | 0 |
| `BenchmarkMixedWorkload` | `ladder_ringbuffer_pool` | 31.18 | 0 | 0 |
| `BenchmarkMixedWorkload` | `ladder_ringbuffer_nopool` | 48.87 | 51 | 0 |
| `BenchmarkEngineChannel` | `default` | 84.01 | 0 | 0 |
| `BenchmarkEngineBatch` | `default` | 22.73 | 0 | 0 |
| `BenchmarkEngineDisruptor` | `default` | 53.74 | 0 | 0 |
| `BenchmarkFastBook` | `AddLimit` | 15.91 | 0 | 0 |
| `BenchmarkFastBook` | `Cancel` | 15.04 | 0 | 0 |
| `BenchmarkFastBook` | `MatchSingle` | 33.70 | 0 | 0 |
| `BenchmarkFastBook` | `MatchSweep` | 207.50 | 0 | 0 |
| `BenchmarkFastBook` | `Market` | 211.50 | 0 | 0 |
| `BenchmarkFastBook` | `BestBidAsk` | 0.98 | 0 | 0 |
| `BenchmarkFastBook` | `Depth` | 22.34 | 0 | 0 |
| `BenchmarkFastBook` | `MixedWorkload` | 26.09 | 0 | 0 |
| `BenchmarkArenaBook` | `AddLimit` | 20.10 | 0 | 0 |
| `BenchmarkArenaBook` | `Cancel` | 13.32 | 0 | 0 |
| `BenchmarkArenaBook` | `MatchSingle` | 32.57 | 0 | 0 |
| `BenchmarkArenaBook` | `MatchSweep` | 201.20 | 0 | 0 |
| `BenchmarkArenaBook` | `Market` | 195.80 | 0 | 0 |
| `BenchmarkArenaBook` | `BestBidAsk` | 1.19 | 0 | 0 |
| `BenchmarkArenaBook` | `Depth` | 26.80 | 0 | 0 |
| `BenchmarkArenaBook` | `MixedWorkload` | 25.46 | 0 | 0 |

## HDR Histogram Latency (Mixed Workload)

| Variant | Throughput (ops/s) | Avg (ns) | p50 (ns) | p90 (ns) | p99 (ns) | p99.9 (ns) | Max (ns) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `sortedslice_linkedlist_pool` | 13322920 | 75.1 | 61 | 99 | 159 | 613 | 139391 |
| `sortedslice_linkedlist_nopool` | 16839292 | 59.4 | 63 | 92 | 153 | 520 | 71487 |
| `sortedslice_ringbuffer_pool` | 22131802 | 45.2 | 57 | 87 | 154 | 807 | 67775 |
| `sortedslice_ringbuffer_nopool` | 15703666 | 63.7 | 70 | 117 | 239 | 859 | 131967 |
| `btree_linkedlist_pool` | 18126654 | 55.2 | 62 | 98 | 167 | 278 | 211071 |
| `btree_linkedlist_nopool` | 12957680 | 77.2 | 79 | 121 | 230 | 891 | 73023 |
| `btree_ringbuffer_pool` | 17556818 | 57.0 | 65 | 105 | 197 | 936 | 226175 |
| `btree_ringbuffer_nopool` | 12328463 | 81.1 | 81 | 157 | 295 | 1490 | 154367 |
| `rbtree_linkedlist_pool` | 15268598 | 65.5 | 55 | 90 | 154 | 674 | 60287 |
| `rbtree_linkedlist_nopool` | 16975133 | 58.9 | 61 | 92 | 172 | 454 | 62815 |
| `rbtree_ringbuffer_pool` | 22486600 | 44.5 | 54 | 90 | 181 | 1009 | 228223 |
| `rbtree_ringbuffer_nopool` | 15487927 | 64.6 | 61 | 94 | 210 | 884 | 61823 |
| `ladder_linkedlist_pool` | 30217040 | 33.1 | 42 | 68 | 121 | 416 | 206335 |
| `ladder_linkedlist_nopool` | 18906748 | 52.9 | 49 | 85 | 216 | 1993 | 189567 |
| `ladder_ringbuffer_pool` | 21353099 | 46.8 | 42 | 69 | 135 | 482 | 58879 |
| `ladder_ringbuffer_nopool` | 16346806 | 61.2 | 53 | 102 | 232 | 2235 | 221823 |

## Profiling Hotspots (pprof — 147 profiled benchmarks)

| Benchmark | Top CPU Hotspot | Top Mem Alloc Space |
|---|---|---|
| `BenchmarkAddLimit/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (47.1% flat) | `runtime.mallocgc` (3594.25kB, 25.7%) |
| `BenchmarkAddLimit/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (19.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (321.53MB, 97.0%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (25.7% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (6372.08kB, 34.1%) |
| `BenchmarkAddLimit/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (296.53MB, 95.9%) |
| `BenchmarkAddLimit/btree_linkedlist_pool` | `sort.Search (inline)` (17.6% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4248.06kB, 36.0%) |
| `BenchmarkAddLimit/btree_linkedlist_nopool` | `github.com/google/btree.(*node[go.shape.*uint8]).get` (6.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (266.02MB, 95.9%) |
| `BenchmarkAddLimit/btree_ringbuffer_pool` | `github.com/google/btree.items[go.shape.*uint8].find` (11.8% flat) | `runtime.mallocgc` (4107.26kB, 26.8%) |
| `BenchmarkAddLimit/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (260.02MB, 95.7%) |
| `BenchmarkAddLimit/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (4248.06kB, 37.8%) |
| `BenchmarkAddLimit/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (14.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (316.03MB, 95.8%) |
| `BenchmarkAddLimit/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (21.6% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4248.06kB, 31.0%) |
| `BenchmarkAddLimit/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (22.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (308.03MB, 96.5%) |
| `BenchmarkAddLimit/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (32.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 43.9%) |
| `BenchmarkAddLimit/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (13.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (420.04MB, 96.5%) |
| `BenchmarkAddLimit/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (16.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 37.9%) |
| `BenchmarkAddLimit/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (12.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (431.54MB, 96.2%) |
| `BenchmarkCancel/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (26.8% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (2832.04kB, 25.2%) |
| `BenchmarkCancel/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (25.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (557.05MB, 98.1%) |
| `BenchmarkCancel/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.9% flat) | `runtime.mallocgc` (4619.31kB, 31.6%) |
| `BenchmarkCancel/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).search` (17.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (755.57MB, 98.6%) |
| `BenchmarkCancel/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (13.2% flat) | `runtime.mallocgc` (5651.41kB, 35.1%) |
| `BenchmarkCancel/btree_linkedlist_nopool` | `github.com/google/btree.items[go.shape.*uint8].find` (10.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (648.56MB, 98.4%) |
| `BenchmarkCancel/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.5% flat) | `runtime.mallocgc` (5130.42kB, 33.3%) |
| `BenchmarkCancel/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (9.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (672.06MB, 98.3%) |
| `BenchmarkCancel/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.7% flat) | `runtime.mallocgc` (5139.36kB, 39.5%) |
| `BenchmarkCancel/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (19.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (655.06MB, 98.6%) |
| `BenchmarkCancel/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (31.4% flat) | `runtime.mallocgc` (5653.30kB, 42.9%) |
| `BenchmarkCancel/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (708.56MB, 98.4%) |
| `BenchmarkCancel/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (34.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 41.3%) |
| `BenchmarkCancel/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (12.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (784.57MB, 97.6%) |
| `BenchmarkCancel/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (11003.72kB, 38.4%) |
| `BenchmarkCancel/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (788.57MB, 98.1%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (80.6% flat) | `runtime.mallocgc` (5643.60kB, 44.1%) |
| `BenchmarkMatchSingle/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (439.04MB, 97.8%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (69.4% flat) | `runtime.mallocgc` (4619.32kB, 38.4%) |
| `BenchmarkMatchSingle/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (31.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (477.54MB, 98.4%) |
| `BenchmarkMatchSingle/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (79.4% flat) | `runtime.mallocgc` (3594.26kB, 29.5%) |
| `BenchmarkMatchSingle/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (31.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (472.54MB, 98.5%) |
| `BenchmarkMatchSingle/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (60.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (5664.08kB, 36.5%) |
| `BenchmarkMatchSingle/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (28.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (441.54MB, 97.8%) |
| `BenchmarkMatchSingle/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (51.5% flat) | `runtime.mallocgc` (4107.26kB, 40.2%) |
| `BenchmarkMatchSingle/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (35.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (459.04MB, 98.2%) |
| `BenchmarkMatchSingle/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (75.7% flat) | `runtime.mallocgc` (5644.37kB, 41.0%) |
| `BenchmarkMatchSingle/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (42.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (450.54MB, 97.6%) |
| `BenchmarkMatchSingle/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (58.3% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 34.1%) |
| `BenchmarkMatchSingle/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (482.54MB, 97.0%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (61.1% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 34.2%) |
| `BenchmarkMatchSingle/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (34.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (445.54MB, 97.1%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_pool` | `runtime.memmove` (60.4% flat) | `runtime.mallocgc` (4106.31kB, 33.0%) |
| `BenchmarkMatchSweep/sortedslice_linkedlist_nopool` | `runtime.memmove` (46.5% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (596.55MB, 98.5%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_pool` | `runtime.memmove` (60.2% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4248.06kB, 18.7%) |
| `BenchmarkMatchSweep/sortedslice_ringbuffer_nopool` | `runtime.memmove` (48.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (577.05MB, 97.4%) |
| `BenchmarkMatchSweep/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.5% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (54.53MB, 51.4%) |
| `BenchmarkMatchSweep/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (266.02MB, 84.1%) |
| `BenchmarkMatchSweep/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (31.9% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt` (36.02MB, 52.8%) |
| `BenchmarkMatchSweep/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (237.52MB, 83.2%) |
| `BenchmarkMatchSweep/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.8% flat) | `runtime.mallocgc` (5131.36kB, 33.2%) |
| `BenchmarkMatchSweep/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (804.07MB, 98.3%) |
| `BenchmarkMatchSweep/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (29.4% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (4248.06kB, 19.6%) |
| `BenchmarkMatchSweep/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (21.1% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (741.57MB, 97.9%) |
| `BenchmarkMatchSweep/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (59.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (5001.69kB, 30.0%) |
| `BenchmarkMatchSweep/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (27.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1002.09MB, 98.1%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (52.8% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (4001.35kB, 22.6%) |
| `BenchmarkMatchSweep/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (36.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (905.08MB, 97.9%) |
| `BenchmarkMarket/sortedslice_linkedlist_pool` | `runtime.memmove` (61.7% flat) | `runtime.mallocgc` (4106.31kB, 27.1%) |
| `BenchmarkMarket/sortedslice_linkedlist_nopool` | `runtime.memmove` (43.8% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (1.14GB, 98.8%) |
| `BenchmarkMarket/sortedslice_ringbuffer_pool` | `runtime.memmove` (63.1% flat) | `runtime.mallocgc` (4621.16kB, 29.0%) |
| `BenchmarkMarket/sortedslice_ringbuffer_nopool` | `runtime.memmove` (41.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (673.06MB, 98.0%) |
| `BenchmarkMarket/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.0% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt` (37.52MB, 58.6%) |
| `BenchmarkMarket/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (15.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (476.04MB, 81.6%) |
| `BenchmarkMarket/btree_ringbuffer_pool` | `github.com/google/btree.items[go.shape.*uint8].find.func1 (inline)` (10.4% flat) | `github.com/google/btree.(*items[go.shape.*uint8]).insertAt (inline)` (25.51MB, 45.3%) |
| `BenchmarkMarket/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (19.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (212.02MB, 83.6%) |
| `BenchmarkMarket/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (37.7% flat) | `runtime.mallocgc` (5130.62kB, 39.0%) |
| `BenchmarkMarket/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (809.07MB, 98.3%) |
| `BenchmarkMarket/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (42.5% flat) | `runtime.mallocgc` (6166.30kB, 29.6%) |
| `BenchmarkMarket/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (20.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (826.58MB, 98.4%) |
| `BenchmarkMarket/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (56.4% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 33.4%) |
| `BenchmarkMarket/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (28.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (973.59MB, 98.2%) |
| `BenchmarkMarket/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (55.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.71kB, 30.5%) |
| `BenchmarkMarket/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (35.0% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (528.55MB, 96.4%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (34.4% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (2832.04kB, 28.6%) |
| `BenchmarkBestBidAsk/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (38.9% flat) | `runtime.mallocgc` (4106.31kB, 50.6%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (40.6% flat) | `runtime.mallocgc` (3090.19kB, 26.6%) |
| `BenchmarkBestBidAsk/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).Best` (40.0% flat) | `runtime.mallocgc` (4619.31kB, 67.5%) |
| `BenchmarkBestBidAsk/btree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (43.8% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (3540.05kB, 38.5%) |
| `BenchmarkBestBidAsk/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (52.9% flat) | `runtime.mallocgc` (5130.38kB, 56.9%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (60.6% flat) | `runtime.mallocgc` (4107.26kB, 40.2%) |
| `BenchmarkBestBidAsk/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (56.2% flat) | `runtime.mallocgc` (3595.21kB, 47.3%) |
| `BenchmarkBestBidAsk/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (44.1% flat) | `runtime.mallocgc` (4114.48kB, 42.5%) |
| `BenchmarkBestBidAsk/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (44.1% flat) | `runtime.mallocgc` (5130.88kB, 65.3%) |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (39.4% flat) | `runtime.mallocgc` (4108.21kB, 40.5%) |
| `BenchmarkBestBidAsk/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).BestBid (inline)` (42.1% flat) | `runtime.mallocgc` (3594.26kB, 48.8%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (53.3% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 45.3%) |
| `BenchmarkBestBidAsk/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (70.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (6002.03kB, 44.3%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (48.6% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 48.7%) |
| `BenchmarkBestBidAsk/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).Best` (50.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 51.5%) |
| `BenchmarkDepth/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (90.0% flat) | `runtime.mallocgc` (3604.14kB, 29.5%) |
| `BenchmarkDepth/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (80.0% flat) | `runtime.mallocgc` (4619.82kB, 50.6%) |
| `BenchmarkDepth/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (91.2% flat) | `runtime.mallocgc` (5140.29kB, 57.3%) |
| `BenchmarkDepth/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*SortedSlice).TopN` (93.9% flat) | `runtime.mallocgc` (3594.21kB, 52.3%) |
| `BenchmarkDepth/btree_linkedlist_pool` | `github.com/google/btree.(*node[go.shape.*uint8]).iterate` (48.0% flat) | `runtime.mallocgc` (4107.21kB, 43.2%) |
| `BenchmarkDepth/btree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (50.0% flat) | `runtime.mallocgc` (4619.31kB, 74.8%) |
| `BenchmarkDepth/btree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (58.3% flat) | `runtime.mallocgc` (5643.62kB, 46.9%) |
| `BenchmarkDepth/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*BTree).TopN.func1` (58.3% flat) | `runtime.mallocgc` (4107.43kB, 54.9%) |
| `BenchmarkDepth/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).TopN` (50.0% flat) | `runtime.mallocgc` (4108.40kB, 53.1%) |
| `BenchmarkDepth/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (45.5% flat) | `runtime.mallocgc` (2570.15kB, 40.2%) |
| `BenchmarkDepth/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (47.6% flat) | `runtime.mallocgc` (3594.21kB, 32.0%) |
| `BenchmarkDepth/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).successor (inline)` (50.0% flat) | `runtime.mallocgc` (4107.26kB, 63.3%) |
| `BenchmarkDepth/ladder_linkedlist_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (85.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (9003.04kB, 51.1%) |
| `BenchmarkDepth/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (94.7% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 46.2%) |
| `BenchmarkDepth/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (82.0% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (7002.37kB, 37.5%) |
| `BenchmarkDepth/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book/levels.(*Ladder).TopN` (89.5% flat) | `github.com/nedo/orbook/book/levels.NewLadder` (8002.70kB, 56.5%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (33.3% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (14868.20kB, 41.4%) |
| `BenchmarkMixedWorkload/sortedslice_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (17.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (265.52MB, 90.2%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (17.9% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (18408.25kB, 39.7%) |
| `BenchmarkMixedWorkload/sortedslice_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (11.4% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (246.02MB, 90.3%) |
| `BenchmarkMixedWorkload/btree_linkedlist_pool` | `sort.Search (inline)` (23.3% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16284.22kB, 44.1%) |
| `BenchmarkMixedWorkload/btree_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (16.7% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (206.52MB, 91.2%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (31.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16992.23kB, 40.1%) |
| `BenchmarkMixedWorkload/btree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (22.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (199.52MB, 87.3%) |
| `BenchmarkMixedWorkload/rbtree_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (31.0% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (16992.23kB, 40.6%) |
| `BenchmarkMixedWorkload/rbtree_linkedlist_nopool` | `github.com/nedo/orbook/book/levels.(*RBTree).GetOrCreate` (17.6% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (257.52MB, 92.0%) |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (23.1% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (17700.24kB, 36.9%) |
| `BenchmarkMixedWorkload/rbtree_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (13.2% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (234.02MB, 88.9%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (13.8% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow` (11328.15kB, 27.4%) |
| `BenchmarkMixedWorkload/ladder_linkedlist_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (14.3% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (296.03MB, 90.3%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_pool` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (25.9% flat) | `github.com/nedo/orbook/pool.(*FreeList).grow (inline)` (19824.27kB, 40.5%) |
| `BenchmarkMixedWorkload/ladder_ringbuffer_nopool` | `github.com/nedo/orbook/book.(*Book).CancelOrder` (21.9% flat) | `github.com/nedo/orbook/pool.NoPool.Get` (272.02MB, 88.8%) |
| `BenchmarkEngineChannel` | `runtime.procyieldAsm` (40.4% flat) | `github.com/nedo/orbook/engine.New (inline)` (45.02MB, 48.6%) |
| `BenchmarkEngineBatch` | `github.com/nedo/orbook/book.(*Book).PlaceOrder` (30.2% flat) | `github.com/nedo/orbook/bench_test.BenchmarkEngineBatch` (69.33MB, 60.1%) |
| `BenchmarkEngineDisruptor` | `github.com/nedo/orbook/engine.(*Disruptor).Publish` (20.0% flat) | `github.com/nedo/orbook/engine.NewDisruptor` (53.82MB, 51.8%) |
| `BenchmarkFastBook/AddLimit` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (35.7% flat) | `github.com/nedo/orbook/book.newFastLadder` (6515.34kB, 29.6%) |
| `BenchmarkFastBook/Cancel` | `github.com/nedo/orbook/book.(*FastBook).CancelOrder` (34.0% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (7002.37kB, 32.0%) |
| `BenchmarkFastBook/MatchSingle` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (62.9% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (9003.04kB, 41.0%) |
| `BenchmarkFastBook/MatchSweep` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (66.1% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (7002.37kB, 34.0%) |
| `BenchmarkFastBook/Market` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (66.7% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (6002.03kB, 33.0%) |
| `BenchmarkFastBook/BestBidAsk` | `github.com/nedo/orbook/book.(*FastBook).BestAsk (inline)` (55.6% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (9003.04kB, 45.2%) |
| `BenchmarkFastBook/Depth` | `github.com/nedo/orbook/book.(*fastLadder).topN` (85.7% flat) | `github.com/nedo/orbook/book.newFastLadder (inline)` (9003.04kB, 44.5%) |
| `BenchmarkFastBook/MixedWorkload` | `github.com/nedo/orbook/book.(*FastBook).PlaceOrder` (40.0% flat) | `github.com/nedo/orbook/book.(*FastBook).growPool` (15576.21kB, 34.9%) |
| `BenchmarkArenaBook/AddLimit` | `github.com/nedo/orbook/book.(*ArenaBook).CancelOrder` (39.5% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (7131.34kB, 37.9%) |
| `BenchmarkArenaBook/Cancel` | `github.com/nedo/orbook/book.(*ArenaBook).CancelOrder` (35.6% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (7131.34kB, 40.9%) |
| `BenchmarkArenaBook/MatchSingle` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (75.0% flat) | `runtime.mallocgc` (5644.37kB, 33.4%) |
| `BenchmarkArenaBook/MatchSweep` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (67.7% flat) | `runtime.mallocgc` (5643.42kB, 40.9%) |
| `BenchmarkArenaBook/Market` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (70.0% flat) | `runtime.mallocgc` (5130.37kB, 33.3%) |
| `BenchmarkArenaBook/BestBidAsk` | `github.com/nedo/orbook/book.(*ArenaBook).BestAsk (inline)` (48.3% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (8022.75kB, 45.2%) |
| `BenchmarkArenaBook/Depth` | `github.com/nedo/orbook/book.(*arenaLadder).topN` (84.2% flat) | `github.com/nedo/orbook/book.newArenaLadder (inline)` (7131.34kB, 46.5%) |
| `BenchmarkArenaBook/MixedWorkload` | `github.com/nedo/orbook/book.(*ArenaBook).PlaceOrder` (45.5% flat) | `github.com/nedo/orbook/book.NewArenaBook` (11037.44kB, 25.4%) |

*(Full top-15 CPU & memory tables for all profiled variants are in `latest.json` and `latest.html`.)*

