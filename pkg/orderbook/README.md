# `pkg/orderbook`

A self-contained, zero-dependency (Go stdlib only), single-ticker stock order book matching engine.

## Architecture
- **Price Levels:** Tick-indexed fixed `Ladder` with a `[]uint64` active bitset (`bits.LeadingZeros64` / `bits.TrailingZeros64`).
- **FIFO Queue:** Intrusive doubly-linked list embedded directly in 64-byte arena order slots (O(1) push, pop, cancel).
- **Order Index:** Flat open-addressed Fibonacci-hashed `orderMap` with single-probe backward-shift deletion (zero tombstones).
- **Memory & GC:** Every internal table (`[]arenaOrder`, `[]priceLevel`, `[]mapSlot`, `[]uint64`) contains **zero Go pointers (`noscan`)**, achieving `0 B/op, 0 allocs/op` in steady state and O(1) GC mark time regardless of resting book depth.
- **Matching Rules:** Price-Time (FIFO) priority, executes at maker price, Limit & Market orders, **Time-in-Force (`GTC`, `IOC`, `FOK`, `PostOnly`)**, **Modify/Amend** (qty reduction preserves FIFO priority; qty increase or price change loses priority and re-matches if crossing), Cancel, and Cancel-Newest Self-Trade Prevention (for `TraderID != 0`).

## Install

```bash
go get github.com/unedtamps/orbook/pkg/orderbook
```

## Quick Start (One `*Book` per Ticker)

`Book` is single-threaded and lock-free by design. Run one `*Book` per ticker inside its own goroutine/actor:

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
        InitialCapacity: 4096, // optional (defaults to 4096)
    })
    if err != nil {
        panic(err)
    }

    // Reuse this slice across calls for 0 allocs/op.
    events := make([]orderbook.Event, 0, 16)

    events = bk.PlaceOrder(orderbook.OrderRequest{
        ID: 1, TraderID: 10, Side: orderbook.Sell, Type: orderbook.Limit, Price: 10000, Qty: 500,
    }, events[:0])

    events = bk.PlaceOrder(orderbook.OrderRequest{
        ID: 2, TraderID: 20, Side: orderbook.Buy, Type: orderbook.Limit, Price: 10000, Qty: 200,
    }, events[:0])

    for _, ev := range events {
        fmt.Printf("%s order=%d price=%d qty=%d\n", ev.Type, ev.OrderID, ev.Price, ev.Qty)
    }

    // Inspect a resting order by ID:
    if snap, ok := bk.GetOrder(1); ok {
        fmt.Printf("Resting order #1 leaves=%d\n", snap.Leaves)
    }

    // Or stream events directly with a zero-alloc callback:
    bk.CancelOrderWithCallback(1, func(ev orderbook.Event) {
        fmt.Printf("%s order=%d reason=%s\n", ev.Type, ev.OrderID, ev.CancelReason)
    })
}
```

## Public API Summary

`orderbook.New(cfg)` returns concrete `*orderbook.Book`, which also satisfies the exported `orderbook.OrderBook` interface (for DI/mocking):

- **Mutations:** `PlaceOrder(req, dst)`, `ModifyOrder(req, dst)`, `CancelOrder(id, dst)`, `PlaceOrderWithCallback(req, fn)`, `ModifyOrderWithCallback(req, fn)`, `CancelOrderWithCallback(id, fn)`, `Reset()`
- **Queries:** `BestBid()`, `BestAsk()`, `TopOfBook()`, `Depth(n, bidsDst, asksDst)`, `GetOrder(id)`, `HasOrder(id)`, `OrderCount()`, `Stats()`, `Config()`
