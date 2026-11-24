package engine_test

import (
	"sync/atomic"
	"testing"

	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/engine"
)

func TestDisruptorAndBatchEngineIdenticalToDirectBook(t *testing.T) {
	bRef := newTestBook(t)
	bDis := newTestBook(t)
	bBat := newTestBook(t)

	var refEvents, disEvents, batEvents []book.Event

	dis := engine.NewDisruptor(bDis, engine.Options{
		QueueSize: 256,
		EventHandler: func(events []book.Event) {
			disEvents = append(disEvents, events...)
		},
	})
	dis.Start()

	bat := engine.NewBatchEngine(bBat, engine.Options{
		QueueSize: 32,
		EventHandler: func(events []book.Event) {
			batEvents = append(batEvents, events...)
		},
	})
	bat.Start()

	const total = 1000
	var evBuf []book.Event
	batch := make([]engine.Command, 0, 32)
	var batches [][]engine.Command

	for i := 1; i <= total; i++ {
		side := book.Buy
		price := int64(9900 + (i%5)*25)
		if i%2 == 0 {
			side = book.Sell
			price = int64(9950 + (i%5)*25)
		}
		req := book.OrderRequest{
			ID:       book.OrderID(i),
			TraderID: book.TraderID(1 + i%6),
			Side:     side,
			Type:     book.Limit,
			Price:    price,
			Qty:      100,
		}
		evBuf = bRef.PlaceOrder(req, evBuf[:0])
		refEvents = append(refEvents, evBuf...)

		dis.PlaceOrderAsync(req)

		batch = append(batch, engine.Command{Type: engine.CmdPlaceOrder, Order: req})
		if i%3 == 0 {
			cid := book.OrderID(i - 2)
			evBuf = bRef.CancelOrder(cid, evBuf[:0])
			refEvents = append(refEvents, evBuf...)

			dis.CancelOrderAsync(cid)
			batch = append(batch, engine.Command{Type: engine.CmdCancelOrder, OrderID: cid})
		}
		if len(batch) >= 32 {
			batches = append(batches, batch)
			bat.SubmitBatch(batch)
			batch = make([]engine.Command, 0, 32)
		}
	}
	if len(batch) > 0 {
		bat.SubmitBatch(batch)
	}

	_ = dis.Stats()
	dis.Stop()
	bat.Stop()

	if len(disEvents) != len(refEvents) {
		t.Fatalf("disruptor events len=%d != ref %d", len(disEvents), len(refEvents))
	}
	if len(batEvents) != len(refEvents) {
		t.Fatalf("batch events len=%d != ref %d", len(batEvents), len(refEvents))
	}
	for i := range refEvents {
		if disEvents[i] != refEvents[i] {
			t.Fatalf("disruptor event[%d] mismatch: %+v vs %+v", i, disEvents[i], refEvents[i])
		}
		if batEvents[i] != refEvents[i] {
			t.Fatalf("batch event[%d] mismatch: %+v vs %+v", i, batEvents[i], refEvents[i])
		}
	}
	if err := bDis.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
	if err := bBat.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}

func TestDisruptorHighVolumeDrain(t *testing.T) {
	b := newTestBook(t)
	var count atomic.Int64
	dis := engine.NewDisruptor(b, engine.Options{
		QueueSize: 1024,
		EventHandler: func(events []book.Event) {
			count.Add(int64(len(events)))
		},
	})
	dis.Start()
	for i := 1; i <= 5000; i++ {
		dis.PlaceOrderAsync(book.OrderRequest{
			ID: book.OrderID(i), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100,
		})
	}
	dis.Stop()
	if count.Load() != 5000 {
		t.Fatalf("got %d events, want 5000", count.Load())
	}
}
