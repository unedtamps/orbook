package engine_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/book/levels"
	"github.com/unedtamps/orbook/book/queue"
	"github.com/unedtamps/orbook/engine"
	"github.com/unedtamps/orbook/pool"
)

func newTestBook(t *testing.T) *book.Book {
	t.Helper()
	cfg := book.Config{
		Symbol:   "BBCA",
		TickSize: 25,
		LotSize:  100,
		MinPrice: 5000,
		MaxPrice: 15000,
	}
	bids, err := levels.NewLadder(book.Buy, cfg, queue.NewLinkedList)
	if err != nil {
		t.Fatal(err)
	}
	asks, err := levels.NewLadder(book.Sell, cfg, queue.NewLinkedList)
	if err != nil {
		t.Fatal(err)
	}
	b, err := book.New(cfg, bids, asks, pool.NewFreeList(1024, 1024))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestEngineSyncAndAsync(t *testing.T) {
	b := newTestBook(t)
	var tradeCount atomic.Int64
	var eventCount atomic.Int64

	eng := engine.New(b, engine.Options{
		QueueSize: 1024,
		EventHandler: func(events []book.Event) {
			eventCount.Add(int64(len(events)))
			for _, ev := range events {
				if ev.Type == book.EventTrade {
					tradeCount.Add(1)
				}
			}
		},
	})
	eng.Start()

	var evBuf []book.Event
	evBuf = eng.PlaceOrder(book.OrderRequest{
		ID: 1, TraderID: 1, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 500,
	}, evBuf)
	if len(evBuf) != 1 || evBuf[0].Type != book.EventAccepted {
		t.Fatalf("unexpected events: %+v", evBuf)
	}

	evBuf = eng.PlaceOrder(book.OrderRequest{
		ID: 2, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 200,
	}, evBuf)
	if len(evBuf) != 3 || evBuf[1].Type != book.EventTrade || evBuf[1].Qty != 200 || evBuf[2].Type != book.EventDone {
		t.Fatalf("unexpected match events: %+v", evBuf)
	}

	top := eng.TopOfBook()
	if top.HasBid || !top.HasAsk || top.AskPrice != 10000 || top.AskQty != 300 {
		t.Fatalf("unexpected TopOfBook: %+v", top)
	}

	bids, asks := eng.Depth(5, nil, nil)
	if len(bids) != 0 || len(asks) != 1 || asks[0].TotalQty != 300 {
		t.Fatalf("unexpected Depth: bids=%+v asks=%+v", bids, asks)
	}

	// Cancel remaining ask asynchronously and drain via Stop.
	eng.CancelOrderAsync(1)
	eng.Stop()

	if tradeCount.Load() != 1 {
		t.Fatalf("got %d trades, want 1", tradeCount.Load())
	}
	if eventCount.Load() != 5 { // Accepted(1) + Accepted(2)+Trade+Done(2) + Cancelled(1)
		t.Fatalf("got %d total events, want 5", eventCount.Load())
	}
	if err := b.CheckInvariants(); err != nil {
		t.Fatalf("invariants failed after Stop: %v", err)
	}
	if b.OrderCount() != 0 {
		t.Fatalf("expected 0 resting orders, got %d", b.OrderCount())
	}
}

func TestEngineConcurrentProducersRaceSafe(t *testing.T) {
	b := newTestBook(t)
	var processed atomic.Int64

	eng := engine.New(b, engine.Options{
		QueueSize: 4096,
		EventHandler: func(events []book.Event) {
			processed.Add(int64(len(events)))
		},
	})
	eng.Start()

	const producers = 8
	const perProducer = 500
	var wg sync.WaitGroup
	wg.Add(producers)

	for p := 0; p < producers; p++ {
		go func(pid int) {
			defer wg.Done()
			baseID := uint64(pid*perProducer + 1)
			side := book.Buy
			price := int64(9900)
			if pid%2 == 1 {
				side = book.Sell
				price = 10000
			}
			for i := 0; i < perProducer; i++ {
				id := book.OrderID(baseID + uint64(i))
				eng.PlaceOrderAsync(book.OrderRequest{
					ID:       id,
					TraderID: book.TraderID(pid + 1),
					Side:     side,
					Type:     book.Limit,
					Price:    price + int64((i%5)*25),
					Qty:      100,
				})
				if i%3 == 0 {
					eng.CancelOrderAsync(id)
				}
			}
		}(p)
	}

	wg.Wait()
	eng.Stop()

	if processed.Load() == 0 {
		t.Fatalf("expected events to be processed")
	}
	if err := b.CheckInvariants(); err != nil {
		t.Fatalf("invariants violated after concurrent run: %v", err)
	}
}
