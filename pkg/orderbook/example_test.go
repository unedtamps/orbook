package orderbook_test

import (
	"fmt"

	"github.com/unedtamps/orbook/pkg/orderbook"
)

func ExampleNew() {
	bk, err := orderbook.New(orderbook.Config{
		Symbol:   "BBCA",
		TickSize: 25,
		LotSize:  100,
		MinPrice: 5000,
		MaxPrice: 15000,
	})
	if err != nil {
		panic(err)
	}

	// Reusable event buffer for 0 allocs/op on the hot path.
	events := make([]orderbook.Event, 0, 16)

	events = bk.PlaceOrder(orderbook.OrderRequest{
		ID: 1, TraderID: 10, Side: orderbook.Sell, Type: orderbook.Limit, Price: 10000, Qty: 500,
	}, events[:0])

	events = bk.PlaceOrder(orderbook.OrderRequest{
		ID: 2, TraderID: 20, Side: orderbook.Buy, Type: orderbook.Limit, Price: 10000, Qty: 200,
	}, events[:0])

	for _, ev := range events {
		if ev.Type == orderbook.EventTrade {
			fmt.Printf("TRADE qty=%d @ %d (makerLeaves=%d)\n", ev.Qty, ev.Price, ev.MakerLeaves)
		}
	}

	askPrice, askQty, _ := bk.BestAsk()
	fmt.Printf("BEST ASK: %d x %d\n", askPrice, askQty)

	// Output:
	// TRADE qty=200 @ 10000 (makerLeaves=300)
	// BEST ASK: 10000 x 300
}

func ExampleBook_PlaceOrderWithCallback() {
	bk, _ := orderbook.New(orderbook.Config{
		Symbol:   "ACME",
		TickSize: 1,
		LotSize:  100,
		MinPrice: 1000,
		MaxPrice: 5000,
	})

	bk.PlaceOrderWithCallback(orderbook.OrderRequest{
		ID: 1, TraderID: 1, Side: orderbook.Sell, Type: orderbook.Limit, Price: 2500, Qty: 100,
	}, nil)

	bk.PlaceOrderWithCallback(orderbook.OrderRequest{
		ID: 2, TraderID: 2, Side: orderbook.Buy, Type: orderbook.Market, Qty: 100,
	}, func(ev orderbook.Event) {
		if ev.Type == orderbook.EventTrade {
			fmt.Printf("FILLED %d @ %d\n", ev.Qty, ev.Price)
		}
	})

	// Output:
	// FILLED 100 @ 2500
}
