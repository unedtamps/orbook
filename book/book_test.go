package book_test

import (
	"testing"

	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/book/levels"
	"github.com/unedtamps/orbook/book/queue"
	"github.com/unedtamps/orbook/pool"
)

type variantSpec struct {
	name      string
	levelKind string
	queueKind string
	usePool   bool
}

func allVariants() []variantSpec {
	levelKinds := []string{"sortedslice", "btree", "rbtree", "ladder"}
	queueKinds := []string{"linkedlist", "ringbuffer"}
	poolOpts := []bool{true, false}

	out := make([]variantSpec, 0, len(levelKinds)*len(queueKinds)*len(poolOpts))
	for _, lk := range levelKinds {
		for _, qk := range queueKinds {
			for _, p := range poolOpts {
				pName := "nopool"
				if p {
					pName = "pool"
				}
				out = append(out, variantSpec{
					name:      lk + "/" + qk + "/" + pName,
					levelKind: lk,
					queueKind: qk,
					usePool:   p,
				})
			}
		}
	}
	return out
}

func buildBook(t *testing.T, cfg book.Config, v variantSpec) *book.Book {
	t.Helper()
	var qf book.QueueFactory
	switch v.queueKind {
	case "linkedlist":
		qf = queue.NewLinkedList
	case "ringbuffer":
		qf = queue.NewRingBuffer
	default:
		t.Fatalf("unknown queue kind %q", v.queueKind)
	}

	makeSide := func(side book.Side) book.PriceLevels {
		switch v.levelKind {
		case "sortedslice":
			return levels.NewSortedSlice(side, qf)
		case "btree":
			return levels.NewBTree(side, qf)
		case "rbtree":
			return levels.NewRBTree(side, qf)
		case "ladder":
			l, err := levels.NewLadder(side, cfg, qf)
			if err != nil {
				t.Fatalf("NewLadder: %v", err)
			}
			return l
		default:
			t.Fatalf("unknown level kind %q", v.levelKind)
			return nil
		}
	}

	var op book.OrderPool
	if v.usePool {
		op = pool.NewFreeList(1024, 1024)
	} else {
		op = pool.NewNoPool()
	}

	b, err := book.New(cfg, makeSide(book.Buy), makeSide(book.Sell), op)
	if err != nil {
		t.Fatalf("book.New: %v", err)
	}
	return b
}

func defaultCfg() book.Config {
	return book.Config{
		Symbol:   "BBCA",
		TickSize: 25,
		LotSize:  100,
		MinPrice: 5000,
		MaxPrice: 15000,
	}
}

func TestConfigValidation(t *testing.T) {
	badConfigs := []book.Config{
		{TickSize: 0, LotSize: 100, MinPrice: 1000, MaxPrice: 2000},
		{TickSize: 25, LotSize: 0, MinPrice: 1000, MaxPrice: 2000},
		{TickSize: 25, LotSize: 100, MinPrice: 0, MaxPrice: 2000},
		{TickSize: 25, LotSize: 100, MinPrice: 2000, MaxPrice: 1000},
		{TickSize: 25, LotSize: 100, MinPrice: 1000, MaxPrice: 1010},
	}
	for i, cfg := range badConfigs {
		if err := cfg.Validate(); err == nil {
			t.Errorf("case %d: expected error for config %+v", i, cfg)
		}
	}
}

func TestOrderValidationAndRejects(t *testing.T) {
	cfg := defaultCfg()
	b := buildBook(t, cfg, allVariants()[0])
	var evs []book.Event

	cases := []struct {
		name string
		req  book.OrderRequest
		want book.RejectReason
	}{
		{"zero ID", book.OrderRequest{ID: 0, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100}, book.RejectInvalidOrderID},
		{"bad side", book.OrderRequest{ID: 1, Side: 99, Type: book.Limit, Price: 10000, Qty: 100}, book.RejectInvalidSide},
		{"bad type", book.OrderRequest{ID: 1, Side: book.Buy, Type: 99, Price: 10000, Qty: 100}, book.RejectInvalidType},
		{"zero qty", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 0}, book.RejectInvalidQty},
		{"neg qty", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: -100}, book.RejectInvalidQty},
		{"non lot", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 150}, book.RejectNonLotQty},
		{"zero price", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 0, Qty: 100}, book.RejectInvalidPrice},
		{"below band", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 4975, Qty: 100}, book.RejectOutOfPriceBand},
		{"above band", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 15025, Qty: 100}, book.RejectOutOfPriceBand},
		{"off tick", book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 10010, Qty: 100}, book.RejectOffTickPrice},
	}

	for _, tc := range cases {
		evs = b.PlaceOrder(tc.req, evs[:0])
		if len(evs) != 1 || evs[0].Type != book.EventRejected || evs[0].RejectReason != tc.want {
			t.Errorf("%s: got %+v, want Rejected(%s)", tc.name, evs, tc.want)
		}
	}

	// Valid order + duplicate ID rejection.
	evs = b.PlaceOrder(book.OrderRequest{ID: 10, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].Type != book.EventAccepted {
		t.Fatalf("expected Accepted, got %+v", evs)
	}
	evs = b.PlaceOrder(book.OrderRequest{ID: 10, TraderID: 2, Side: book.Sell, Type: book.Limit, Price: 10100, Qty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].Type != book.EventRejected || evs[0].RejectReason != book.RejectDuplicateOrderID {
		t.Fatalf("expected RejectDuplicateOrderID, got %+v", evs)
	}

	// Cancel missing order.
	evs = b.CancelOrder(999, evs[:0])
	if len(evs) != 1 || evs[0].Type != book.EventRejected || evs[0].RejectReason != book.RejectOrderNotFound {
		t.Fatalf("expected RejectOrderNotFound, got %+v", evs)
	}

	// Cancel order 10, then verify order ID 10 can be reused.
	evs = b.CancelOrder(10, evs[:0])
	if len(evs) != 1 || evs[0].Type != book.EventCancelled || evs[0].CancelReason != book.CancelUser {
		t.Fatalf("expected CancelUser, got %+v", evs)
	}
	evs = b.PlaceOrder(book.OrderRequest{ID: 10, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].Type != book.EventAccepted {
		t.Fatalf("expected reuse of cancelled ID to succeed, got %+v", evs)
	}
}

func TestAllVariantsMatchingScenarios(t *testing.T) {
	cfg := defaultCfg()
	for _, v := range allVariants() {
		v := v
		t.Run(v.name, func(t *testing.T) {
			t.Run("FIFOAndMakerPriceExecution", func(t *testing.T) {
				b := buildBook(t, cfg, v)
				var evs []book.Event

				// Three asks at 10000 (FIFO: 1, 2, 3) and one ask at 10025 (4).
				evs = b.PlaceOrder(book.OrderRequest{ID: 1, TraderID: 1, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
				evs = b.PlaceOrder(book.OrderRequest{ID: 2, TraderID: 2, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 200}, evs[:0])
				evs = b.PlaceOrder(book.OrderRequest{ID: 3, TraderID: 3, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
				evs = b.PlaceOrder(book.OrderRequest{ID: 4, TraderID: 4, Side: book.Sell, Type: book.Limit, Price: 10025, Qty: 300}, evs[:0])

				// Aggressive Buy at 10050 for 500 sweeps all of 10000 (400) + 100 of 10025.
				evs = b.PlaceOrder(book.OrderRequest{ID: 10, TraderID: 9, Side: book.Buy, Type: book.Limit, Price: 10050, Qty: 500}, evs[:0])

				// Expected sequence:
				// Accepted(10),
				// Trade(maker=1, p=10000, q=100), Done(1),
				// Trade(maker=2, p=10000, q=200), Done(2),
				// Trade(maker=3, p=10000, q=100), Done(3),
				// Trade(maker=4, p=10025, q=100, makerLeaves=200),
				// Done(10)
				if len(evs) != 9 {
					t.Fatalf("got %d events (%+v), want 9", len(evs), evs)
				}
				if evs[1].MakerOrderID != 1 || evs[1].Price != 10000 || evs[1].Qty != 100 {
					t.Fatalf("trade 1 mismatch: %+v", evs[1])
				}
				if evs[3].MakerOrderID != 2 || evs[3].Price != 10000 || evs[3].Qty != 200 {
					t.Fatalf("trade 2 mismatch: %+v", evs[3])
				}
				if evs[5].MakerOrderID != 3 || evs[5].Price != 10000 || evs[5].Qty != 100 {
					t.Fatalf("trade 3 mismatch: %+v", evs[5])
				}
				if evs[7].MakerOrderID != 4 || evs[7].Price != 10025 || evs[7].Qty != 100 || evs[7].MakerLeaves != 200 {
					t.Fatalf("trade 4 mismatch: %+v", evs[7])
				}
				if evs[8].Type != book.EventDone || evs[8].OrderID != 10 {
					t.Fatalf("expected Done(10), got %+v", evs[8])
				}

				askP, askQ, ok := b.BestAsk()
				if !ok || askP != 10025 || askQ != 200 {
					t.Fatalf("BestAsk = (%d, %d, %v), want (10025, 200, true)", askP, askQ, ok)
				}
				if _, _, hasBid := b.BestBid(); hasBid {
					t.Fatalf("expected no resting bids")
				}
				if err := b.CheckInvariants(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("PartialFillRemainderRests", func(t *testing.T) {
				b := buildBook(t, cfg, v)
				var evs []book.Event

				b.PlaceOrder(book.OrderRequest{ID: 1, TraderID: 1, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 200}, evs[:0])
				evs = b.PlaceOrder(book.OrderRequest{ID: 2, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: 10025, Qty: 500}, evs[:0])

				// Accepted(2), Trade(10000, 200, leaves=300), Done(1) -> 300 rests at 10025.
				if len(evs) != 3 || evs[1].Type != book.EventTrade || evs[1].Leaves != 300 || evs[2].Type != book.EventDone {
					t.Fatalf("unexpected events: %+v", evs)
				}
				bidP, bidQ, ok := b.BestBid()
				if !ok || bidP != 10025 || bidQ != 300 {
					t.Fatalf("BestBid = (%d, %d, %v), want (10025, 300, true)", bidP, bidQ, ok)
				}
				if _, _, hasAsk := b.BestAsk(); hasAsk {
					t.Fatalf("expected empty ask side")
				}
				if err := b.CheckInvariants(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("MarketOrdersFullPartialAndEmpty", func(t *testing.T) {
				b := buildBook(t, cfg, v)
				var evs []book.Event

				// Market order on empty book -> Accepted + Cancelled(CancelMarketUnfilled).
				evs = b.PlaceOrder(book.OrderRequest{ID: 1, TraderID: 1, Side: book.Buy, Type: book.Market, Qty: 300}, evs[:0])
				if len(evs) != 2 || evs[0].Type != book.EventAccepted || evs[1].Type != book.EventCancelled ||
					evs[1].CancelReason != book.CancelMarketUnfilled || evs[1].Qty != 300 {
					t.Fatalf("empty market buy events: %+v", evs)
				}

				// Seed 2 bid levels: 10000 (100) and 9975 (200).
				b.PlaceOrder(book.OrderRequest{ID: 2, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 3, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: 9975, Qty: 200}, evs[:0])

				// Market Sell for 500 sweeps 300 and cancels 200 unfilled.
				evs = b.PlaceOrder(book.OrderRequest{ID: 4, TraderID: 3, Side: book.Sell, Type: book.Market, Qty: 500}, evs[:0])
				if len(evs) != 6 {
					t.Fatalf("got %d events (%+v), want 6", len(evs), evs)
				}
				last := evs[len(evs)-1]
				if last.Type != book.EventCancelled || last.CancelReason != book.CancelMarketUnfilled || last.Qty != 200 {
					t.Fatalf("expected 200 CancelMarketUnfilled, got %+v", last)
				}
				if b.OrderCount() != 0 {
					t.Fatalf("expected empty book, got %d orders", b.OrderCount())
				}
				if err := b.CheckInvariants(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("SelfTradePreventionCancelNewest", func(t *testing.T) {
				b := buildBook(t, cfg, v)
				var evs []book.Event

				// Ask at 10000: Trader 2 (100), then Trader 1 (200).
				b.PlaceOrder(book.OrderRequest{ID: 1, TraderID: 2, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 2, TraderID: 1, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 200}, evs[:0])

				// Trader 1 submits aggressive Buy at 10050 for 400:
				// Must trade 100 with Trader 2, then cancel remaining 300 due to CancelSTP
				// without resting and without touching Trader 1's resting Ask #2.
				evs = b.PlaceOrder(book.OrderRequest{ID: 3, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10050, Qty: 400}, evs[:0])
				if len(evs) != 4 {
					t.Fatalf("got %d events (%+v), want 4", len(evs), evs)
				}
				if evs[1].Type != book.EventTrade || evs[1].MakerOrderID != 1 || evs[1].Qty != 100 {
					t.Fatalf("expected trade with maker 1: %+v", evs[1])
				}
				if evs[2].Type != book.EventDone || evs[2].OrderID != 1 {
					t.Fatalf("expected Done(1): %+v", evs[2])
				}
				if evs[3].Type != book.EventCancelled || evs[3].OrderID != 3 ||
					evs[3].CancelReason != book.CancelSTP || evs[3].Qty != 300 {
					t.Fatalf("expected CancelSTP for 300: %+v", evs[3])
				}
				if !b.HasOrder(2) || b.HasOrder(3) {
					t.Fatalf("expected resting order 2 intact and taker 3 gone")
				}
				if err := b.CheckInvariants(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("CancelMiddleFrontTailAndPartiallyFilled", func(t *testing.T) {
				b := buildBook(t, cfg, v)
				var evs []book.Event

				b.PlaceOrder(book.OrderRequest{ID: 1, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 300}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 2, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 200}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 3, TraderID: 3, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 4, TraderID: 4, Side: book.Buy, Type: book.Limit, Price: 9975, Qty: 400}, evs[:0])

				// Partially fill front order #1 by 100 (leaves=200).
				b.PlaceOrder(book.OrderRequest{ID: 5, TraderID: 9, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 100}, evs[:0])

				// Cancel middle (#2), tail (#3), then partially-filled front (#1).
				evs = b.CancelOrder(2, evs[:0])
				if len(evs) != 1 || evs[0].Type != book.EventCancelled || evs[0].Qty != 200 {
					t.Fatalf("cancel #2: %+v", evs)
				}
				evs = b.CancelOrder(3, evs[:0])
				if len(evs) != 1 || evs[0].Type != book.EventCancelled || evs[0].Qty != 100 {
					t.Fatalf("cancel #3: %+v", evs)
				}
				evs = b.CancelOrder(1, evs[:0])
				if len(evs) != 1 || evs[0].Type != book.EventCancelled || evs[0].Qty != 200 {
					t.Fatalf("cancel #1 remaining leaves: %+v", evs)
				}

				// Level 10000 should now be gone; BestBid transitions to 9975.
				bidP, bidQ, ok := b.BestBid()
				if !ok || bidP != 9975 || bidQ != 400 {
					t.Fatalf("BestBid = (%d, %d, %v), want (9975, 400, true)", bidP, bidQ, ok)
				}
				if err := b.CheckInvariants(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("DepthSnapshot", func(t *testing.T) {
				b := buildBook(t, cfg, v)
				var evs []book.Event
				b.PlaceOrder(book.OrderRequest{ID: 1, Side: book.Buy, Type: book.Limit, Price: 9900, Qty: 100}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 2, Side: book.Buy, Type: book.Limit, Price: 9950, Qty: 200}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 3, Side: book.Buy, Type: book.Limit, Price: 9950, Qty: 300}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 4, Side: book.Buy, Type: book.Limit, Price: 9925, Qty: 400}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 5, Side: book.Sell, Type: book.Limit, Price: 10050, Qty: 100}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 6, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 200}, evs[:0])
				b.PlaceOrder(book.OrderRequest{ID: 7, Side: book.Sell, Type: book.Limit, Price: 10025, Qty: 300}, evs[:0])

				bids, asks := b.Depth(2, nil, nil)
				if len(bids) != 2 || bids[0] != (book.DepthLevel{Price: 9950, TotalQty: 500, OrderCount: 2}) ||
					bids[1] != (book.DepthLevel{Price: 9925, TotalQty: 400, OrderCount: 1}) {
					t.Fatalf("unexpected bids depth: %+v", bids)
				}
				if len(asks) != 2 || asks[0] != (book.DepthLevel{Price: 10000, TotalQty: 200, OrderCount: 1}) ||
					asks[1] != (book.DepthLevel{Price: 10025, TotalQty: 300, OrderCount: 1}) {
					t.Fatalf("unexpected asks depth: %+v", asks)
				}
			})
		})
	}
}

func TestZeroAllocHotPathWithPool(t *testing.T) {
	cfg := defaultCfg()
	// Check pooled variants achieve 0 allocs/op on steady-state Place+Cancel and Match.
	for _, lk := range []string{"sortedslice", "btree", "rbtree", "ladder"} {
		for _, qk := range []string{"linkedlist", "ringbuffer"} {
			v := variantSpec{name: lk + "/" + qk + "/pool", levelKind: lk, queueKind: qk, usePool: true}
			t.Run(v.name, func(t *testing.T) {
				b := buildBook(t, cfg, v)
				evBuf := make([]book.Event, 0, 16)
				bidsBuf := make([]book.DepthLevel, 0, 16)
				asksBuf := make([]book.DepthLevel, 0, 16)

				// Warm up levels, map buckets, and pool.
				for i := 0; i < 20; i++ {
					p := int64(9500 + i*25)
					b.PlaceOrder(book.OrderRequest{ID: book.OrderID(i + 1), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: p, Qty: 100}, evBuf[:0])
					b.PlaceOrder(book.OrderRequest{ID: book.OrderID(i + 100), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: p + 1000, Qty: 100}, evBuf[:0])
				}
				// Warm up target tick 10000.
				b.PlaceOrder(book.OrderRequest{ID: 999, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100}, evBuf[:0])
				b.CancelOrder(999, evBuf[:0])

				var nextID book.OrderID = 10000
				allocs := testing.AllocsPerRun(500, func() {
					id := nextID
					nextID++
					evBuf = b.PlaceOrder(book.OrderRequest{
						ID: id, TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 10000, Qty: 100,
					}, evBuf[:0])
					evBuf = b.CancelOrder(id, evBuf[:0])
					_, _, _ = b.BestBid()
					_, _, _ = b.BestAsk()
					bidsBuf, asksBuf = b.Depth(5, bidsBuf, asksBuf)
				})
				if allocs != 0 {
					t.Errorf("%s: expected 0 allocs/op on steady-state hot path, got %.2f", v.name, allocs)
				}
			})
		}
	}
}

func TestBookResetDeterministicReplay(t *testing.T) {
	cfg := defaultCfg()
	reqs := []book.OrderRequest{
		{ID: 1, TraderID: 1, Side: book.Sell, Type: book.Limit, Price: 10000, Qty: 200},
		{ID: 2, TraderID: 2, Side: book.Sell, Type: book.Limit, Price: 10025, Qty: 300},
		{ID: 3, TraderID: 3, Side: book.Buy, Type: book.Limit, Price: 10025, Qty: 400},
		{ID: 4, TraderID: 4, Side: book.Buy, Type: book.Limit, Price: 9975, Qty: 500},
	}

	for _, v := range allVariants() {
		v := v
		t.Run(v.name, func(t *testing.T) {
			b := buildBook(t, cfg, v)
			runOnce := func() []book.Event {
				var all []book.Event
				for _, r := range reqs {
					all = b.PlaceOrder(r, all)
				}
				return all
			}
			ev1 := runOnce()
			stats1 := b.Stats()

			b.Reset()
			if b.OrderCount() != 0 || b.BidLevelCount() != 0 || b.AskLevelCount() != 0 {
				t.Fatalf("book not empty after Reset: %+v", b.Stats())
			}
			if err := b.CheckInvariants(); err != nil {
				t.Fatalf("invariants failed right after Reset: %v", err)
			}

			ev2 := runOnce()
			if b.Stats() != stats1 {
				t.Fatalf("stats after Reset+replay mismatch: %+v vs %+v", b.Stats(), stats1)
			}
			if len(ev1) != len(ev2) {
				t.Fatalf("event count mismatch: %d vs %d", len(ev1), len(ev2))
			}
			for i := range ev1 {
				if ev1[i] != ev2[i] {
					t.Fatalf("event[%d] mismatch after Reset: %+v vs %+v", i, ev1[i], ev2[i])
				}
			}
			if err := b.CheckInvariants(); err != nil {
				t.Fatalf("invariants failed after Reset+replay: %v", err)
			}
		})
	}
}

func TestBookResetPreservesFreeListPool(t *testing.T) {
	cfg := defaultCfg()
	p := pool.NewFreeList(64, 64)
	bids, err := levels.NewLadder(book.Buy, cfg, queue.NewLinkedList)
	if err != nil {
		t.Fatal(err)
	}
	asks, err := levels.NewLadder(book.Sell, cfg, queue.NewLinkedList)
	if err != nil {
		t.Fatal(err)
	}
	b, err := book.New(cfg, bids, asks, p)
	if err != nil {
		t.Fatal(err)
	}

	var evBuf []book.Event
	for round := 0; round < 50; round++ {
		for i := 1; i <= 20; i++ {
			evBuf = b.PlaceOrder(book.OrderRequest{
				ID: book.OrderID(i), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: 9900, Qty: 100,
			}, evBuf[:0])
		}
		b.Reset()
		if p.FreeCount() != p.Allocated() {
			t.Fatalf("round %d: pool leaked orders across Reset: free=%d alloc=%d", round, p.FreeCount(), p.Allocated())
		}
	}
	if p.Allocated() != 64 {
		t.Fatalf("expected pool never to grow past 64 across resets, got %d", p.Allocated())
	}
}
