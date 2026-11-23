package book_test

import (
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/unedtamps/orbook/book"
)

type diffOp struct {
	isCancel bool
	req      book.OrderRequest
	cancelID book.OrderID
}

type orderLedger struct {
	origQty      int64
	filledQty    int64
	cancelledQty int64
	terminal     bool
}

func TestDifferentialAllImplementations(t *testing.T) {
	cfg := book.Config{
		Symbol:   "BBCA",
		TickSize: 25,
		LotSize:  100,
		MinPrice: 8000,
		MaxPrice: 12000,
	}

	variants := allVariants()
	books := make([]*book.Book, len(variants))
	evBufs := make([][]book.Event, len(variants))
	for i, v := range variants {
		books[i] = buildBook(t, cfg, v)
		evBufs[i] = make([]book.Event, 0, 32)
	}
	fastBk, err := book.NewFastBook(cfg, 4096)
	if err != nil {
		t.Fatal(err)
	}
	fastEvs := make([]book.Event, 0, 32)
	arenaBk, err := book.NewArenaBook(cfg, 4096)
	if err != nil {
		t.Fatal(err)
	}
	arenaEvs := make([]book.Event, 0, 32)

	rng := rand.New(rand.NewSource(20261122))
	const totalOps = 25000

	mid := int64(10000)
	var nextID book.OrderID = 1
	activeIDs := make([]book.OrderID, 0, 2048)
	activePos := make(map[book.OrderID]int, 2048)
	ledger := make(map[book.OrderID]*orderLedger, totalOps)

	removeActive := func(id book.OrderID) {
		idx, ok := activePos[id]
		if !ok {
			return
		}
		last := activeIDs[len(activeIDs)-1]
		activeIDs[idx] = last
		activePos[last] = idx
		activeIDs = activeIDs[:len(activeIDs)-1]
		delete(activePos, id)
	}
	addActive := func(id book.OrderID) {
		activePos[id] = len(activeIDs)
		activeIDs = append(activeIDs, id)
	}

	for step := 0; step < totalOps; step++ {
		// Random-walk mid-price within interior of price band.
		if step%10 == 0 {
			delta := int64(rng.Intn(3)-1) * cfg.TickSize
			if mid+delta >= cfg.MinPrice+500 && mid+delta <= cfg.MaxPrice-500 {
				mid += delta
			}
		}

		var op diffOp
		roll := rng.Intn(100)
		switch {
		case roll < 2:
			// Occasional invalid order (off-tick, non-lot, out-of-band, or duplicate).
			op = diffOp{
				req: book.OrderRequest{
					ID:       nextID,
					TraderID: book.TraderID(rng.Intn(8) + 1),
					Side:     book.Buy,
					Type:     book.Limit,
					Price:    mid + 7, // off-tick
					Qty:      100,
				},
			}
			nextID++
		case roll < 32 && len(activeIDs) > 0:
			// 30% Cancel (mix of active orders and occasional unknown ID).
			if rng.Intn(20) == 0 {
				op = diffOp{isCancel: true, cancelID: 999_999_999}
			} else {
				pick := activeIDs[rng.Intn(len(activeIDs))]
				op = diffOp{isCancel: true, cancelID: pick}
			}
		case roll < 42:
			// 10% Market order.
			side := book.Buy
			if rng.Intn(2) == 0 {
				side = book.Sell
			}
			op = diffOp{
				req: book.OrderRequest{
					ID:       nextID,
					TraderID: book.TraderID(rng.Intn(8) + 1),
					Side:     side,
					Type:     book.Market,
					Qty:      int64(rng.Intn(5)+1) * cfg.LotSize,
				},
			}
			nextID++
		default:
			// ~58% Limit order distributed around mid (both passive and crossing).
			side := book.Buy
			if rng.Intn(2) == 0 {
				side = book.Sell
			}
			offsetTicks := int64(rng.Intn(21) - 10) // [-10, +10] ticks around mid
			price := mid + offsetTicks*cfg.TickSize
			if price < cfg.MinPrice {
				price = cfg.MinPrice
			}
			if price > cfg.MaxPrice {
				price = cfg.MaxPrice
			}
			op = diffOp{
				req: book.OrderRequest{
					ID:       nextID,
					TraderID: book.TraderID(rng.Intn(8) + 1), // 8 traders => frequent realistic STP collisions
					Side:     side,
					Type:     book.Limit,
					Price:    price,
					Qty:      int64(rng.Intn(6)+1) * cfg.LotSize,
				},
			}
			nextID++
		}

		// Execute op across all 16 Book variants + FastBook.
		for i, b := range books {
			if op.isCancel {
				evBufs[i] = b.CancelOrder(op.cancelID, evBufs[i][:0])
			} else {
				evBufs[i] = b.PlaceOrder(op.req, evBufs[i][:0])
			}
		}
		if op.isCancel {
			fastEvs = fastBk.CancelOrder(op.cancelID, fastEvs[:0])
			arenaEvs = arenaBk.CancelOrder(op.cancelID, arenaEvs[:0])
		} else {
			fastEvs = fastBk.PlaceOrder(op.req, fastEvs[:0])
			arenaEvs = arenaBk.PlaceOrder(op.req, arenaEvs[:0])
		}

		// Assert identical event streams across all implementations.
		refEvs := evBufs[0]
		if len(fastEvs) != len(refEvs) {
			t.Fatalf("step %d: FastBook produced %d events, ref produced %d", step, len(fastEvs), len(refEvs))
		}
		if len(arenaEvs) != len(refEvs) {
			t.Fatalf("step %d: ArenaBook produced %d events, ref produced %d", step, len(arenaEvs), len(refEvs))
		}
		for j := range refEvs {
			if fastEvs[j] != refEvs[j] {
				t.Fatalf("step %d: FastBook event[%d] mismatch:\ngot=%+v\nref=%+v", step, j, fastEvs[j], refEvs[j])
			}
			if arenaEvs[j] != refEvs[j] {
				t.Fatalf("step %d: ArenaBook event[%d] mismatch:\ngot=%+v\nref=%+v", step, j, arenaEvs[j], refEvs[j])
			}
		}
		for i := 1; i < len(books); i++ {
			got := evBufs[i]
			if len(got) != len(refEvs) {
				t.Fatalf("step %d: variant %s produced %d events, ref (%s) produced %d\ngot=%+v\nref=%+v",
					step, variants[i].name, len(got), variants[0].name, len(refEvs), got, refEvs)
			}
			for j := range refEvs {
				if got[j] != refEvs[j] {
					t.Fatalf("step %d: variant %s event[%d] mismatch:\ngot=%+v\nref=%+v",
						step, variants[i].name, j, got[j], refEvs[j])
				}
			}
		}

		// Update conservation ledger & active ID tracker from reference events.
		for _, ev := range refEvs {
			switch ev.Type {
			case book.EventAccepted:
				ledger[ev.OrderID] = &orderLedger{origQty: ev.Qty}
				if ev.OrderType == book.Limit {
					addActive(ev.OrderID)
				}
			case book.EventTrade:
				if ev.Qty <= 0 {
					t.Fatalf("step %d: non-positive trade qty %+v", step, ev)
				}
				ledger[ev.OrderID].filledQty += ev.Qty
				ledger[ev.MakerOrderID].filledQty += ev.Qty
			case book.EventCancelled:
				l := ledger[ev.OrderID]
				if l.terminal {
					t.Fatalf("step %d: order %d cancelled after already terminal", step, ev.OrderID)
				}
				l.cancelledQty += ev.Qty
				l.terminal = true
				removeActive(ev.OrderID)
				if l.filledQty+l.cancelledQty != l.origQty {
					t.Fatalf("step %d: conservation failed on cancel for order %d: filled=%d + cancelled=%d != orig=%d",
						step, ev.OrderID, l.filledQty, l.cancelledQty, l.origQty)
				}
			case book.EventDone:
				l := ledger[ev.OrderID]
				if l.terminal {
					t.Fatalf("step %d: order %d Done after already terminal", step, ev.OrderID)
				}
				l.terminal = true
				removeActive(ev.OrderID)
				if l.filledQty != l.origQty || l.cancelledQty != 0 {
					t.Fatalf("step %d: conservation failed on Done for order %d: filled=%d cancelled=%d orig=%d",
						step, ev.OrderID, l.filledQty, l.cancelledQty, l.origQty)
				}
			}
		}

		// Periodically check full structural invariants on every implementation.
		if step%500 == 0 {
			for i, b := range books {
				if err := b.CheckInvariants(); err != nil {
					t.Fatalf("step %d: variant %s invariant error: %v", step, variants[i].name, err)
				}
			}
		}
	}

	// Final invariant, depth, top-of-book, stats, and conservation verification.
	refTop := books[0].TopOfBook()
	refStats := books[0].Stats()
	refBids, refAsks := books[0].Depth(1000, nil, nil)

	if fastBk.TopOfBook() != refTop || fastBk.Stats() != refStats {
		t.Fatalf("final FastBook top/stats mismatch: top=%+v stats=%+v", fastBk.TopOfBook(), fastBk.Stats())
	}
	if arenaBk.TopOfBook() != refTop || arenaBk.Stats() != refStats {
		t.Fatalf("final ArenaBook top/stats mismatch: top=%+v stats=%+v", arenaBk.TopOfBook(), arenaBk.Stats())
	}

	for i, b := range books {
		if err := b.CheckInvariants(); err != nil {
			t.Fatalf("final variant %s invariant error: %v", variants[i].name, err)
		}
		if b.TopOfBook() != refTop {
			t.Fatalf("final TopOfBook mismatch on %s: got %+v want %+v", variants[i].name, b.TopOfBook(), refTop)
		}
		if b.Stats() != refStats {
			t.Fatalf("final Stats mismatch on %s: got %+v want %+v", variants[i].name, b.Stats(), refStats)
		}
		gotBids, gotAsks := b.Depth(1000, nil, nil)
		if len(gotBids) != len(refBids) || len(gotAsks) != len(refAsks) {
			t.Fatalf("final Depth len mismatch on %s", variants[i].name)
		}
		for j := range refBids {
			if gotBids[j] != refBids[j] {
				t.Fatalf("final bids[%d] mismatch on %s: %+v vs %+v", j, variants[i].name, gotBids[j], refBids[j])
			}
		}
		for j := range refAsks {
			if gotAsks[j] != refAsks[j] {
				t.Fatalf("final asks[%d] mismatch on %s: %+v vs %+v", j, variants[i].name, gotAsks[j], refAsks[j])
			}
		}
	}

	// Cancel all remaining active orders on ref book and confirm total quantity conservation!
	remaining := append([]book.OrderID(nil), activeIDs...)
	for _, id := range remaining {
		evBufs[0] = books[0].CancelOrder(id, evBufs[0][:0])
		if len(evBufs[0]) != 1 || evBufs[0][0].Type != book.EventCancelled {
			t.Fatalf("expected active order %d to cancel cleanly, got %+v", id, evBufs[0])
		}
		l := ledger[id]
		l.cancelledQty += evBufs[0][0].Qty
		l.terminal = true
		if l.filledQty+l.cancelledQty != l.origQty {
			t.Fatalf("final conservation failed for order %d: filled=%d + cancelled=%d != orig=%d",
				id, l.filledQty, l.cancelledQty, l.origQty)
		}
	}
	if books[0].OrderCount() != 0 {
		t.Fatalf("expected 0 orders left after cancelling all tracked active IDs, got %d", books[0].OrderCount())
	}
	for id, l := range ledger {
		if !l.terminal {
			t.Fatalf("order %d was never terminated", id)
		}
	}
}

func FuzzOrderBookDifferential(f *testing.F) {
	f.Add([]byte{0, 1, 1, 0, 5, 2, 1, 2, 2, 0, 5, 3, 2, 1, 0, 0, 0, 0})
	f.Add([]byte{1, 1, 1, 1, 10, 4, 0, 2, 2, 0, 9, 2})

	cfg := book.Config{
		Symbol:   "FUZZ",
		TickSize: 10,
		LotSize:  10,
		MinPrice: 100,
		MaxPrice: 500,
	}
	numTicks := cfg.NumTicks()

	f.Fuzz(func(t *testing.T, data []byte) {
		variants := []variantSpec{
			{name: "sortedslice/linkedlist/pool", levelKind: "sortedslice", queueKind: "linkedlist", usePool: true},
			{name: "btree/ringbuffer/pool", levelKind: "btree", queueKind: "ringbuffer", usePool: true},
			{name: "rbtree/linkedlist/pool", levelKind: "rbtree", queueKind: "linkedlist", usePool: true},
			{name: "ladder/ringbuffer/nopool", levelKind: "ladder", queueKind: "ringbuffer", usePool: false},
			{name: "ladder/linkedlist/pool", levelKind: "ladder", queueKind: "linkedlist", usePool: true},
		}
		books := make([]*book.Book, len(variants))
		evBufs := make([][]book.Event, len(variants))
		for i, v := range variants {
			books[i] = buildBook(t, cfg, v)
		}

		var nextID book.OrderID = 1
		for off := 0; off+6 <= len(data); off += 6 {
			chunk := data[off : off+6]
			action := chunk[0] % 4
			side := book.Buy
			if chunk[1]%2 == 1 {
				side = book.Sell
			}
			trader := book.TraderID(chunk[2] % 5)
			tick := int64(int(chunk[3]) % numTicks)
			price := cfg.MinPrice + tick*cfg.TickSize
			qty := int64(chunk[4]%10+1) * cfg.LotSize

			if action == 0 && nextID > 1 {
				cancelID := book.OrderID(uint64(binary.LittleEndian.Uint16(chunk[4:6]))%uint64(nextID-1) + 1)
				for i, b := range books {
					evBufs[i] = b.CancelOrder(cancelID, evBufs[i][:0])
				}
			} else {
				oType := book.Limit
				if action == 1 {
					oType = book.Market
				}
				req := book.OrderRequest{
					ID:       nextID,
					TraderID: trader,
					Side:     side,
					Type:     oType,
					Price:    price,
					Qty:      qty,
				}
				nextID++
				for i, b := range books {
					evBufs[i] = b.PlaceOrder(req, evBufs[i][:0])
				}
			}

			for i := 1; i < len(books); i++ {
				if len(evBufs[i]) != len(evBufs[0]) {
					t.Fatalf("event count mismatch %s vs %s", variants[i].name, variants[0].name)
				}
				for j := range evBufs[0] {
					if evBufs[i][j] != evBufs[0][j] {
						t.Fatalf("event[%d] mismatch %s vs %s", j, variants[i].name, variants[0].name)
					}
				}
			}
		}

		for i, b := range books {
			if err := b.CheckInvariants(); err != nil {
				t.Fatalf("invariant failed on %s: %v", variants[i].name, err)
			}
		}
	})
}
