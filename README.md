# Orbook

`orbook` is a deterministic, single-ticker stock limit order book matching engine written in Go.

The repository started as a systematic benchmark lab comparing 16 data-structure combinations across price-level stores (`Ladder`, `SortedSlice`, `RBTree`, `BTree`), FIFO queues (`LinkedList`, `RingBuffer`), and memory allocators (`FreeList`, `NoPool`), followed by iterative profiling-driven optimizations. The winning architecture—a devirtualized, pointer-free (`noscan`) arena engine—is packaged as a standalone, zero-dependency Go library in [`pkg/orderbook`](pkg/orderbook).

---

## Benchmark Summary

Measured on `Intel(R) Core(TM) i9-14900HX` (`linux/amd64`, `go1.27.3`), `TickSize=1`, `LotSize=100`, 20,001 price ticks, steady-state **`0 B/op, 0 allocs/op`**:

| Operation | Baseline (`Book` + Go `map`) | `ArenaBook` (core of `pkg/orderbook`) | Speedup |
|---|---:|---:|---:|
| `AddLimit` (resting, 200 lvls) | 31.04 ns/op | **19.36 ns/op** | **1.60×** |
| `Cancel` (random resting order) | 38.46 ns/op | **13.38 ns/op** | **2.87×** |
| `MatchSingle` (1-level FIFO) | 59.26 ns/op | **31.82 ns/op** | **1.86×** |
| `MatchSweep` (10 levels) | 400.80 ns/op | **200.50 ns/op** | **2.00×** |
| `BestBidAsk` (top of book) | 2.82 ns/op | **1.06 ns/op** | **2.66×** |
| `MixedWorkload` (60% Limit / 30% Cancel / 10% Market) | 45.14 ns/op | **25.33 ns/op** (~39.5M ops/s) | **1.78×** |
| HDR Latency (200k mixed ops) | 59 / 156 / 230 ns *(p50/p99/p99.9)* | **43 / 89 / 143 ns** *(32.98M ops/s)* | **1.61×** tighter p99.9 |
| `runtime.GC()` pause (250k resting orders) | 1.97 ms *(pointer-based `FastBook`)* | **453 µs** | **4.34×** faster GC |

> For the full 16-variant comparison matrix, `pprof` hotspot breakdowns, and stage-by-stage optimization history, see [`BENCHMARK.md`](BENCHMARK.md), [`CHANGELOG.md`](CHANGELOG.md), and [`reports/latest.md`](reports/latest.md).

---

## Inside `pkg/orderbook`

`pkg/orderbook` is a self-contained, **stdlib-only** package designed to be embedded directly in matching services, market simulators, or exchange gateways.

```bash
go get github.com/unedtamps/orbook/pkg/orderbook
```

### 1. Concurrency & Ownership Model

`*orderbook.Book` is **single-threaded, deterministic, and lock-free by design**. Rather than paying mutex contention per order, run one `*Book` instance per ticker inside its own dedicated goroutine (actor / LMAX-style loop).

`orderbook.New(cfg)` returns the concrete `*Book` for zero-overhead direct calls, while also satisfying the exported `orderbook.OrderBook` interface for dependency injection and testing.

### 2. Pointer-Free (`noscan`) Memory Layout

In Go, even a zero-allocation engine suffers GC mark latency if its resting structures contain pointers—the garbage collector must traverse every live `*Order` pointer on the heap during mark phases.

`pkg/orderbook` eliminates Go pointers from all internal data structures by replacing pointers with compact `uint32` arena and tick indices. Every backing slice (`[]arenaOrder`, `[]priceLevel`, `[]mapSlot`, `[]uint64`) is classified as **`noscan`** by the Go runtime, keeping GC mark overhead O(1) regardless of how many orders rest on the book.

Each core struct is sized to an exact divisor or multiple of an x86-64 64-byte L1 cache line (verified by `TestStructSizesNoscan`):

| Struct | Size | Cache Packing | Role |
|---|---:|---|---|
| `arenaOrder` | **64 B** | 1 per cache line | Order state, `PrevIdx`/`NextIdx` intrusive FIFO & free-list links, `TickIdx` ladder back-reference |
| `priceLevel` | **32 B** | 2 per cache line | Aggregated `price`, `totalQty`, `orderCount`, and `headIdx`/`tailIdx` FIFO endpoints |
| `mapSlot` | **16 B** | 4 per cache line | Flat hash table entry (`OrderID` key + `uint32` arena index) |

### 3. Core Data Structures

#### Tick-Indexed Price Ladder (`pkg/orderbook/ladder.go`)
Because equity markets operate within discrete tick increments and daily price bands (`[MinPrice, MaxPrice]`), each side (`bids` and `asks`) is stored as a contiguous `[]priceLevel` slice indexed directly by tick offset:

```
tickIdx = (price - minPrice) / tickSize   // subtraction only when TickSize == 1
```

A compact `[]uint64` active bitset tracks non-empty price levels (64 ticks per word). When the best price level is exhausted during matching or cancellation, the ladder locates the next active level in hardware using `bits.LeadingZeros64` (highest bid) or `bits.TrailingZeros64` (lowest ask), skipping empty price gaps 64 ticks at a time. The same bitset drives `Depth(n)` and `FOK` fillability scans.

#### Contiguous Order Arena & Intrusive FIFO Queue (`pkg/orderbook/book.go`)
Resting orders live in a flat `[]arenaOrder` slice where index `0` is reserved as the nil sentinel.
- **Allocation:** `allocOrder()` pops from an intrusive singly-linked free-list (`freeHead`, threaded through `arenaOrder.NextIdx`) in O(1), only bumping `highMark` or doubling the arena when peak concurrent resting orders exceed capacity.
- **FIFO Queue:** Each `priceLevel` maintains `headIdx` and `tailIdx` over an intrusive doubly-linked list (`PrevIdx`/`NextIdx`) embedded inside `arenaOrder`. Enqueueing, matching from the head, and cancelling an arbitrary resting order all execute in strict worst-case O(1) with zero per-level slice allocations.
- **Direct Level Access:** Each `arenaOrder` stores its `TickIdx`, allowing `CancelOrder` and in-place `ModifyOrder` to jump straight to the owning `priceLevel` without re-computing or searching for the price.

#### Flat Open-Addressed `orderMap` (`pkg/orderbook/ordermap.go`)
Active orders are indexed by `OrderID` in a custom power-of-2 open-addressed linear-probing hash table (maximum 50% load factor):
- **Fibonacci Hashing:** `((x ^ (x >> 32)) * 11400714819323198485) >> shift` avoids Go runtime map hashing overhead and inlines cleanly into hot paths.
- **Backward-Shift Deletion:** Removing a key shifts subsequent probe-chain entries back toward their home slots instead of leaving tombstones, keeping probe lengths short over long-running sessions.
- **Single-Probe Cancel:** `orders.Delete(id)` returns `(arenaIdx, ok)` in one lookup, avoiding the `Get` + `delete` double-probe required by Go's built-in `map`.

### 4. Matching Rules & Order Lifecycle

`pkg/orderbook` enforces deterministic **Price-Time (FIFO)** priority, executing all trades at the resting **maker price**:

- **Validation:** Every incoming request is checked against `Config` (`OrderID >= 1`, unique active `OrderID`, valid `Side`/`OrderType`/`TimeInForce`, positive `Qty` aligned to `LotSize`, and `Price` within `[MinPrice, MaxPrice]` aligned to `TickSize`). Invalid requests emit `EventRejected` with a typed `RejectReason` without mutating book state.
- **Order Types & Time-in-Force (`TimeInForce`):**
  - `Limit` + `GTC` (default): matches aggressively up to limit price, then rests any unfilled remainder.
  - `Limit` + `IOC`: fills immediately available liquidity up to limit price, then emits `EventCancelled(CancelIOC)` for any unfilled remainder.
  - `Limit` / `Market` + `FOK`: runs a read-only `canFillFOK` scan across active bitset levels (including STP checks) before executing; if 100% of `Qty` cannot fill immediately, emits `EventCancelled(CancelFOK)` with zero trades.
  - `Limit` + `PostOnly`: checks `wouldCross` against the opposite best price; if crossing, emits `EventCancelled(CancelPostOnly)` immediately instead of taking liquidity.
  - `Market`: sweeps opposite levels until filled or book exhausted; any unfilled remainder emits `EventCancelled(CancelMarketUnfilled)`.
- **Order Amendment (`ModifyOrder`):**
  - Decreasing `NewQty` at the **same price** updates `Leaves` and `priceLevel.totalQty` in-place in O(1), **preserving FIFO time priority**.
  - Increasing `NewQty` or changing `NewPrice` unlinks the order, loses time priority, emits `EventModified`, re-matches if crossing, and rests any remaining quantity at the back of its new price level.
- **Self-Trade Prevention (STP):** When `TraderID != 0` and an incoming order would match a resting order with the same `TraderID`, **Cancel-Newest** triggers: the resting maker stays untouched and the incoming taker's remaining quantity is cancelled with `CancelSTP`.

### 5. Zero-Allocation Event Output

Every mutation supports two zero-allocation delivery styles that produce identical, monotonically sequenced `Event` records (`EventAccepted`, `EventRejected`, `EventTrade`, `EventModified`, `EventCancelled`, `EventDone`):

1. **Caller-owned slice buffer:** `PlaceOrder(req, dst)`, `ModifyOrder(req, dst)`, `CancelOrder(id, dst)` — pass `events[:0]` to reuse backing capacity across calls.
2. **Direct streaming callback:** `PlaceOrderWithCallback(req, fn)`, `ModifyOrderWithCallback(req, fn)`, `CancelOrderWithCallback(id, fn)` — streams each `Event` value straight to `fn` without buffering.

### 6. Quick Start

```go
package main

import (
	"fmt"
	"github.com/unedtamps/orbook/pkg/orderbook"
)

func main() {
	bk, err := orderbook.New(orderbook.Config{
		Symbol:          "BBCA",
		TickSize:        25,
		LotSize:         100,
		MinPrice:        5000,
		MaxPrice:        15000,
		InitialCapacity: 4096, // optional; defaults to 4096
	})
	if err != nil {
		panic(err)
	}

	// Reuse this slice across calls for 0 allocs/op.
	events := make([]orderbook.Event, 0, 16)

	// Rest a GTC limit sell order: 500 @ 10000
	events = bk.PlaceOrder(orderbook.OrderRequest{
		ID: 1, TraderID: 10, Side: orderbook.Sell, Type: orderbook.Limit, Price: 10000, Qty: 500,
	}, events[:0])

	// Aggressive buy: 200 @ 10000 -> executes at maker price 10000
	events = bk.PlaceOrder(orderbook.OrderRequest{
		ID: 2, TraderID: 20, Side: orderbook.Buy, Type: orderbook.Limit, Price: 10000, Qty: 200,
	}, events[:0])

	for _, ev := range events {
		if ev.Type == orderbook.EventTrade {
			fmt.Printf("TRADE #%d qty=%d @ %d (makerLeaves=%d)\n", ev.TradeID, ev.Qty, ev.Price, ev.MakerLeaves)
		}
	}

	// In-place qty reduction preserves FIFO priority (300 -> 200)
	events = bk.ModifyOrder(orderbook.ModifyRequest{ID: 1, NewQty: 200}, events[:0])

	// Query top of book and resting order snapshot
	top := bk.TopOfBook()
	fmt.Printf("Best Ask: %d x %d\n", top.AskPrice, top.AskQty)

	if snap, ok := bk.GetOrder(1); ok {
		fmt.Printf("Order #%d leaves=%d\n", snap.ID, snap.Leaves)
	}

	// Or stream events directly via callback
	bk.CancelOrderWithCallback(1, func(ev orderbook.Event) {
		fmt.Printf("%s order=%d reason=%s\n", ev.Type, ev.OrderID, ev.CancelReason)
	})
}
```

### 7. Public API Reference

| Category | Methods |
|---|---|
| **Mutations (Slice)** | `PlaceOrder(req, dst)`, `ModifyOrder(req, dst)`, `CancelOrder(id, dst)` |
| **Mutations (Callback)** | `PlaceOrderWithCallback(req, fn)`, `ModifyOrderWithCallback(req, fn)`, `CancelOrderWithCallback(id, fn)` |
| **Lifecycle** | `New(cfg)`, `Reset()` *(clears orders & counters while retaining pre-allocated storage)* |
| **Market Data Queries** | `BestBid()`, `BestAsk()`, `TopOfBook()`, `Depth(n, bidsDst, asksDst)` |
| **Order & State Queries** | `GetOrder(id)`, `HasOrder(id)`, `OrderCount()`, `Stats()`, `Config()` |
