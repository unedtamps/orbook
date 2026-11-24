package orderbook

import (
	"math"
	"math/rand"
	"testing"
	"unsafe"
)

func testCfg() Config {
	return Config{
		Symbol:   "BBCA",
		TickSize: 25,
		LotSize:  100,
		MinPrice: 5000,
		MaxPrice: 15000,
	}
}

func TestStructSizesNoscan(t *testing.T) {
	if sz := unsafe.Sizeof(arenaOrder{}); sz != 64 {
		t.Fatalf("arenaOrder size = %d, want 64", sz)
	}
	if sz := unsafe.Sizeof(priceLevel{}); sz != 32 {
		t.Fatalf("priceLevel size = %d, want 32", sz)
	}
	if sz := unsafe.Sizeof(mapSlot{}); sz != 16 {
		t.Fatalf("mapSlot size = %d, want 16", sz)
	}
	if sz := unsafe.Sizeof(Event{}); sz != 88 {
		t.Fatalf("Event size = %d, want 88", sz)
	}
}

func TestConfigValidation(t *testing.T) {
	bad := []Config{
		{TickSize: 0, LotSize: 100, MinPrice: 1000, MaxPrice: 2000},
		{TickSize: 25, LotSize: 0, MinPrice: 1000, MaxPrice: 2000},
		{TickSize: 25, LotSize: 100, MinPrice: 0, MaxPrice: 2000},
		{TickSize: 25, LotSize: 100, MinPrice: 2000, MaxPrice: 1000},
		{TickSize: 25, LotSize: 100, MinPrice: 1000, MaxPrice: 1010},
	}
	for i, c := range bad {
		if _, err := New(c); err == nil {
			t.Errorf("case %d: expected error for %+v", i, c)
		}
	}
}

func TestMatchingSTPGetOrderAndReset(t *testing.T) {
	b, err := New(testCfg())
	if err != nil {
		t.Fatal(err)
	}
	var evs []Event

	// Reject off-tick & non-lot.
	evs = b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Limit, Price: 10010, Qty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].RejectReason != RejectOffTickPrice {
		t.Fatalf("expected RejectOffTickPrice, got %+v", evs)
	}
	evs = b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Limit, Price: 10000, Qty: 150}, evs[:0])
	if len(evs) != 1 || evs[0].RejectReason != RejectNonLotQty {
		t.Fatalf("expected RejectNonLotQty, got %+v", evs)
	}

	// Rest two asks at 10000 (Trader 1 then Trader 2) and one at 10025.
	b.PlaceOrder(OrderRequest{ID: 10, TraderID: 1, Side: Sell, Type: Limit, Price: 10000, Qty: 200}, evs[:0])
	b.PlaceOrder(OrderRequest{ID: 11, TraderID: 2, Side: Sell, Type: Limit, Price: 10000, Qty: 300}, evs[:0])
	b.PlaceOrder(OrderRequest{ID: 12, TraderID: 3, Side: Sell, Type: Limit, Price: 10025, Qty: 400}, evs[:0])

	// Verify GetOrder.
	snap, ok := b.GetOrder(11)
	if !ok || snap.ID != 11 || snap.Price != 10000 || snap.Leaves != 300 || snap.TraderID != 2 {
		t.Fatalf("unexpected GetOrder(11): %+v ok=%v", snap, ok)
	}
	if _, ok := b.GetOrder(999); ok {
		t.Fatalf("expected GetOrder(999) == false")
	}

	// Trader 2 aggressively buys 400 at 10025 -> fills 200 vs Trader 1 (#10), then hits STP on own #11!
	evs = b.PlaceOrder(OrderRequest{ID: 20, TraderID: 2, Side: Buy, Type: Limit, Price: 10025, Qty: 400}, evs[:0])
	if len(evs) != 4 || evs[1].Type != EventTrade || evs[1].MakerOrderID != 10 || evs[1].Qty != 200 ||
		evs[2].Type != EventDone || evs[3].Type != EventCancelled || evs[3].CancelReason != CancelSTP || evs[3].Qty != 200 {
		t.Fatalf("unexpected STP events: %+v", evs)
	}
	if b.HasOrder(10) || !b.HasOrder(11) || b.HasOrder(20) {
		t.Fatalf("unexpected resting state after STP")
	}

	// Market Buy for 1000 sweeps #11 (300) + #12 (400), cancels 300 unfilled.
	evs = b.PlaceOrder(OrderRequest{ID: 21, TraderID: 9, Side: Buy, Type: Market, Qty: 1000}, evs[:0])
	last := evs[len(evs)-1]
	if last.Type != EventCancelled || last.CancelReason != CancelMarketUnfilled || last.Qty != 300 {
		t.Fatalf("unexpected market unfilled tail: %+v", last)
	}
	if b.OrderCount() != 0 {
		t.Fatalf("expected empty book, got %d", b.OrderCount())
	}

	// Reset and confirm clean counters.
	b.Reset()
	if st := b.Stats(); st.OrderSeq != 0 || st.EventSeq != 0 || st.TradeSeq != 0 || st.ActiveOrders != 0 {
		t.Fatalf("expected zeroed stats after Reset, got %+v", st)
	}
}

func TestCallbackMatchesSliceIdentically(t *testing.T) {
	cfg := testCfg()
	bSlice, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	bCB, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewSource(42))
	var sliceEvs, cbEvs []Event
	var nextID OrderID = 1

	for step := 0; step < 5000; step++ {
		sliceEvs = sliceEvs[:0]
		cbEvs = cbEvs[:0]

		roll := rng.Intn(100)
		if roll < 25 && nextID > 1 {
			cid := OrderID(rng.Intn(int(nextID-1)) + 1)
			sliceEvs = bSlice.CancelOrder(cid, sliceEvs)
			bCB.CancelOrderWithCallback(cid, func(ev Event) {
				cbEvs = append(cbEvs, ev)
			})
		} else if roll < 45 && nextID > 1 {
			mid := OrderID(rng.Intn(int(nextID-1)) + 1)
			mreq := ModifyRequest{
				ID:       mid,
				NewPrice: 9900 + int64(rng.Intn(9))*25,
				NewQty:   int64(rng.Intn(5)+1) * 100,
			}
			sliceEvs = bSlice.ModifyOrder(mreq, sliceEvs)
			bCB.ModifyOrderWithCallback(mreq, func(ev Event) {
				cbEvs = append(cbEvs, ev)
			})
		} else {
			side := Buy
			if rng.Intn(2) == 1 {
				side = Sell
			}
			oType := Limit
			tif := TimeInForce(rng.Intn(4)) // GTC, IOC, FOK, PostOnly
			if rng.Intn(10) == 0 {
				oType = Market
				if tif == PostOnly {
					tif = IOC
				}
			}
			req := OrderRequest{
				ID:          nextID,
				TraderID:    TraderID(rng.Intn(6) + 1),
				Side:        side,
				Type:        oType,
				TimeInForce: tif,
				Price:       9900 + int64(rng.Intn(9))*25,
				Qty:         int64(rng.Intn(5)+1) * 100,
			}
			nextID++
			sliceEvs = bSlice.PlaceOrder(req, sliceEvs)
			bCB.PlaceOrderWithCallback(req, func(ev Event) {
				cbEvs = append(cbEvs, ev)
			})
		}

		if len(sliceEvs) != len(cbEvs) {
			t.Fatalf("step %d: len mismatch %d vs %d", step, len(sliceEvs), len(cbEvs))
		}
		for i := range sliceEvs {
			if sliceEvs[i] != cbEvs[i] {
				t.Fatalf("step %d: event[%d] mismatch:\nslice=%+v\ncb=%+v", step, i, sliceEvs[i], cbEvs[i])
			}
		}
	}

	if bSlice.TopOfBook() != bCB.TopOfBook() || bSlice.Stats() != bCB.Stats() {
		t.Fatalf("final state mismatch: top=%+v/%+v stats=%+v/%+v",
			bSlice.TopOfBook(), bCB.TopOfBook(), bSlice.Stats(), bCB.Stats())
	}
}

func TestZeroAllocSteadyState(t *testing.T) {
	b, err := New(testCfg())
	if err != nil {
		t.Fatal(err)
	}
	evBuf := make([]Event, 0, 16)
	bidsBuf := make([]DepthLevel, 0, 16)
	asksBuf := make([]DepthLevel, 0, 16)

	var nextID OrderID = 1
	allocs := testing.AllocsPerRun(500, func() {
		id := nextID
		nextID++
		evBuf = b.PlaceOrder(OrderRequest{ID: id, TraderID: 1, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, evBuf[:0])
		evBuf = b.CancelOrder(id, evBuf[:0])
		_, _, _ = b.BestBid()
		bidsBuf, asksBuf = b.Depth(5, bidsBuf, asksBuf)
	})
	if allocs != 0 {
		t.Fatalf("PlaceOrder/CancelOrder steady-state allocs = %.2f, want 0", allocs)
	}

	var cbCount int
	cb := func(_ Event) { cbCount++ }
	cbAllocs := testing.AllocsPerRun(500, func() {
		id := nextID
		nextID++
		b.PlaceOrderWithCallback(OrderRequest{ID: id, TraderID: 1, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, cb)
		b.CancelOrderWithCallback(id, cb)
	})
	if cbAllocs != 0 {
		t.Fatalf("Callback API steady-state allocs = %.2f, want 0", cbAllocs)
	}
}

func BenchmarkBookAddCancel(b *testing.B) {
	cfg := Config{Symbol: "ACME", TickSize: 1, LotSize: 100, MinPrice: 10000, MaxPrice: 30000, InitialCapacity: 8192}
	bk, err := New(cfg)
	if err != nil {
		b.Fatal(err)
	}
	evBuf := make([]Event, 0, 16)
	var id OrderID = 1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evBuf = bk.PlaceOrder(OrderRequest{ID: id, TraderID: 1, Side: Buy, Type: Limit, Price: 20000, Qty: 100}, evBuf[:0])
		evBuf = bk.CancelOrder(id, evBuf[:0])
		id++
	}
}

func TestTimeInForceIOCFOKPostOnly(t *testing.T) {
	b, err := New(testCfg())
	if err != nil {
		t.Fatal(err)
	}
	var evs []Event

	// Market + PostOnly is invalid.
	evs = b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Market, TimeInForce: PostOnly, Qty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].Type != EventRejected || evs[0].RejectReason != RejectInvalidTimeInForce {
		t.Fatalf("expected RejectInvalidTimeInForce, got %+v", evs)
	}

	// Seed asks: 10000 (200 by Trader 1), 10025 (300 by Trader 2).
	b.PlaceOrder(OrderRequest{ID: 10, TraderID: 1, Side: Sell, Type: Limit, Price: 10000, Qty: 200}, evs[:0])
	b.PlaceOrder(OrderRequest{ID: 11, TraderID: 2, Side: Sell, Type: Limit, Price: 10025, Qty: 300}, evs[:0])

	// PostOnly Buy at 10000 crosses -> cancelled with CancelPostOnly, 0 trades, ask untouched.
	evs = b.PlaceOrder(OrderRequest{ID: 20, TraderID: 9, Side: Buy, Type: Limit, TimeInForce: PostOnly, Price: 10000, Qty: 100}, evs[:0])
	if len(evs) != 2 || evs[1].Type != EventCancelled || evs[1].CancelReason != CancelPostOnly || evs[1].Qty != 100 {
		t.Fatalf("expected CancelPostOnly, got %+v", evs)
	}
	if b.HasOrder(20) {
		t.Fatalf("PostOnly crossing order should not rest")
	}

	// PostOnly Buy at 9975 does NOT cross -> rests on the book.
	evs = b.PlaceOrder(OrderRequest{ID: 21, TraderID: 9, Side: Buy, Type: Limit, TimeInForce: PostOnly, Price: 9975, Qty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].Type != EventAccepted || !b.HasOrder(21) {
		t.Fatalf("expected passive PostOnly to rest, got %+v", evs)
	}

	// FOK Buy at 10025 for 600 (only 500 available) -> cancelled with CancelFOK, 0 trades!
	evs = b.PlaceOrder(OrderRequest{ID: 22, TraderID: 9, Side: Buy, Type: Limit, TimeInForce: FOK, Price: 10025, Qty: 600}, evs[:0])
	if len(evs) != 2 || evs[1].Type != EventCancelled || evs[1].CancelReason != CancelFOK || evs[1].Qty != 600 {
		t.Fatalf("expected CancelFOK, got %+v", evs)
	}

	// FOK Buy by Trader 1 at 10025 for 300 -> blocked by Trader 1's own ask #10 at 10000 (STP!) -> CancelFOK, 0 trades!
	evs = b.PlaceOrder(OrderRequest{ID: 23, TraderID: 1, Side: Buy, Type: Limit, TimeInForce: FOK, Price: 10025, Qty: 300}, evs[:0])
	if len(evs) != 2 || evs[1].Type != EventCancelled || evs[1].CancelReason != CancelFOK {
		t.Fatalf("expected STP-blocked FOK to cancel atomically with CancelFOK, got %+v", evs)
	}

	// FOK Buy by Trader 9 at 10025 for 400 -> fills 200 @ 10000 + 200 @ 10025, Done!
	evs = b.PlaceOrder(OrderRequest{ID: 24, TraderID: 9, Side: Buy, Type: Limit, TimeInForce: FOK, Price: 10025, Qty: 400}, evs[:0])
	if len(evs) != 5 || evs[1].Type != EventTrade || evs[3].Type != EventTrade || evs[4].Type != EventDone {
		t.Fatalf("expected FOK full fill, got %+v", evs)
	}

	// IOC Buy at 10025 for 300 (only 100 left on #11) -> fills 100, cancels 200 with CancelIOC, does not rest!
	evs = b.PlaceOrder(OrderRequest{ID: 25, TraderID: 9, Side: Buy, Type: Limit, TimeInForce: IOC, Price: 10025, Qty: 300}, evs[:0])
	last := evs[len(evs)-1]
	if last.Type != EventCancelled || last.CancelReason != CancelIOC || last.Qty != 200 || b.HasOrder(25) {
		t.Fatalf("expected IOC partial fill + CancelIOC(200), got %+v", evs)
	}
	if _, _, hasAsk := b.BestAsk(); hasAsk {
		t.Fatalf("expected asks to be empty")
	}
}

func TestModifyOrderPriorityAndCrossing(t *testing.T) {
	b, err := New(testCfg())
	if err != nil {
		t.Fatal(err)
	}
	var evs []Event

	// Two resting bids at 10000: #1 (300, older) then #2 (200, newer).
	b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Limit, Price: 10000, Qty: 300}, evs[:0])
	b.PlaceOrder(OrderRequest{ID: 2, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 200}, evs[:0])

	// Reject no-change and unknown ID.
	evs = b.ModifyOrder(ModifyRequest{ID: 999, NewPrice: 10000, NewQty: 100}, evs[:0])
	if evs[0].Type != EventRejected || evs[0].RejectReason != RejectOrderNotFound {
		t.Fatalf("expected RejectOrderNotFound, got %+v", evs)
	}
	evs = b.ModifyOrder(ModifyRequest{ID: 1, NewPrice: 10000, NewQty: 300}, evs[:0])
	if evs[0].Type != EventRejected || evs[0].RejectReason != RejectNoChange {
		t.Fatalf("expected RejectNoChange, got %+v", evs)
	}

	// 1. Reduce #1 qty from 300 -> 100 at same price: MUST KEEP FIFO PRIORITY ahead of #2!
	evs = b.ModifyOrder(ModifyRequest{ID: 1, NewQty: 100}, evs[:0])
	if len(evs) != 1 || evs[0].Type != EventModified || evs[0].Leaves != 100 {
		t.Fatalf("expected EventModified(100), got %+v", evs)
	}
	if _, qty, _ := b.BestBid(); qty != 300 {
		t.Fatalf("BestBid qty = %d, want 300 (100 + 200)", qty)
	}

	// Incoming Sell for 100 at 10000 MUST fill #1 first (proving #1 kept priority)!
	evs = b.PlaceOrder(OrderRequest{ID: 3, TraderID: 9, Side: Sell, Type: Limit, Price: 10000, Qty: 100}, evs[:0])
	if evs[1].Type != EventTrade || evs[1].MakerOrderID != 1 {
		t.Fatalf("expected #1 to fill first (kept FIFO priority), got %+v", evs)
	}

	// 2. Now rest #4 (200) behind #2 (200) at 10000, then INCREASE #2 qty from 200 -> 300:
	// #2 MUST LOSE PRIORITY and move behind #4!
	b.PlaceOrder(OrderRequest{ID: 4, TraderID: 4, Side: Buy, Type: Limit, Price: 10000, Qty: 200}, evs[:0])
	evs = b.ModifyOrder(ModifyRequest{ID: 2, NewQty: 300}, evs[:0])
	if len(evs) != 1 || evs[0].Type != EventModified || evs[0].Leaves != 300 {
		t.Fatalf("expected EventModified(300), got %+v", evs)
	}
	evs = b.PlaceOrder(OrderRequest{ID: 5, TraderID: 9, Side: Sell, Type: Limit, Price: 10000, Qty: 200}, evs[:0])
	if evs[1].Type != EventTrade || evs[1].MakerOrderID != 4 {
		t.Fatalf("expected #4 to fill before #2 after #2 increased qty, got %+v", evs)
	}

	// 3. Rest an Ask #6 at 10050 (200), then Modify Bid #2's price up to 10050 -> MUST CROSS & MATCH!
	b.PlaceOrder(OrderRequest{ID: 6, TraderID: 6, Side: Sell, Type: Limit, Price: 10050, Qty: 200}, evs[:0])
	evs = b.ModifyOrder(ModifyRequest{ID: 2, NewPrice: 10050}, evs[:0])
	// Expected: Modified(#2 @ 10050, 300), Trade(maker=#6 @ 10050, 200, leaves=100), Done(#6) -> 100 rests at 10050!
	if len(evs) != 3 || evs[0].Type != EventModified || evs[1].Type != EventTrade || evs[1].MakerOrderID != 6 || evs[2].Type != EventDone {
		t.Fatalf("expected crossing modify to match #6, got %+v", evs)
	}
	bidP, bidQ, ok := b.BestBid()
	if !ok || bidP != 10050 || bidQ != 100 {
		t.Fatalf("BestBid = (%d, %d, %v), want (10050, 100, true)", bidP, bidQ, ok)
	}
}

func BenchmarkBookModifyInPlace(b *testing.B) {
	cfg := Config{Symbol: "ACME", TickSize: 1, LotSize: 100, MinPrice: 10000, MaxPrice: 30000, InitialCapacity: 8192}
	bk, _ := New(cfg)
	evBuf := make([]Event, 0, 16)
	bk.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Limit, Price: 20000, Qty: 1_000_000_000}, evBuf[:0])
	qty := int64(1_000_000_000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qty -= 100
		if qty <= 100 {
			b.StopTimer()
			bk.ModifyOrder(ModifyRequest{ID: 1, NewQty: 1_000_000_000}, evBuf[:0])
			qty = 1_000_000_000 - 100
			b.StartTimer()
		}
		evBuf = bk.ModifyOrder(ModifyRequest{ID: 1, NewQty: qty}, evBuf[:0])
	}
}

func TestCallbackBookStateConsistent(t *testing.T) {
	b, err := New(testCfg())
	if err != nil {
		t.Fatal(err)
	}

	// 1. Seed a single resting ask #1 (100 @ 10000), then fully fill it via callback.
	b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Sell, Type: Limit, Price: 10000, Qty: 100}, nil)

	var seenTrade, seenDone bool
	var bidsBuf, asksBuf []DepthLevel
	b.PlaceOrderWithCallback(OrderRequest{ID: 2, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, func(ev Event) {
		switch ev.Type {
		case EventTrade:
			seenTrade = true
			if p, q, ok := b.BestAsk(); ok {
				t.Errorf("during EventTrade (full fill): BestAsk = (%d, %d, true), want empty", p, q)
			}
			if b.TopOfBook().HasAsk {
				t.Errorf("during EventTrade (full fill): TopOfBook.HasAsk = true, want false")
			}
			if snap, ok := b.GetOrder(1); ok {
				t.Errorf("during EventTrade (full fill): GetOrder(1) = %+v, want not found", snap)
			}
			_, asksBuf = b.Depth(5, bidsBuf, asksBuf)
			if len(asksBuf) != 0 {
				t.Errorf("during EventTrade (full fill): Depth asks = %+v, want empty", asksBuf)
			}
			if st := b.Stats(); st.ActiveOrders != 0 || st.AskLevels != 0 {
				t.Errorf("during EventTrade (full fill): Stats = %+v, want 0 active orders and 0 ask levels", st)
			}
		case EventDone:
			if ev.OrderID == 1 {
				seenDone = true
				if p, q, ok := b.BestAsk(); ok {
					t.Errorf("during EventDone: BestAsk = (%d, %d, true), want empty", p, q)
				}
				_, asksBuf = b.Depth(5, bidsBuf, asksBuf)
				if len(asksBuf) != 0 {
					t.Errorf("during EventDone: Depth asks = %+v, want empty", asksBuf)
				}
				if st := b.Stats(); st.ActiveOrders != 0 || st.AskLevels != 0 {
					t.Errorf("during EventDone: Stats = %+v, want 0 active orders and 0 ask levels", st)
				}
			}
		}
	})
	if !seenTrade || !seenDone {
		t.Fatalf("expected both EventTrade and maker EventDone, got trade=%v done=%v", seenTrade, seenDone)
	}

	// 2. Non-crossing priority-losing ModifyOrderWithCallback must expose the resting order inside EventModified.
	b.PlaceOrder(OrderRequest{ID: 10, TraderID: 1, Side: Buy, Type: Limit, Price: 9900, Qty: 200}, nil)
	var seenMod bool
	b.ModifyOrderWithCallback(ModifyRequest{ID: 10, NewPrice: 9950, NewQty: 300}, func(ev Event) {
		if ev.Type != EventModified {
			return
		}
		seenMod = true
		snap, ok := b.GetOrder(10)
		if !ok || snap.Price != 9950 || snap.Leaves != 300 {
			t.Errorf("during EventModified (non-crossing): GetOrder(10) = (%+v, %v), want Price=9950 Leaves=300", snap, ok)
		}
		if b.OrderCount() != 1 {
			t.Errorf("during EventModified (non-crossing): OrderCount = %d, want 1", b.OrderCount())
		}
	})
	if !seenMod {
		t.Fatalf("expected EventModified callback")
	}
}

func TestCallbackReentrancyPanics(t *testing.T) {
	b, err := New(testCfg())
	if err != nil {
		t.Fatal(err)
	}
	b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Sell, Type: Limit, Price: 10000, Qty: 200}, nil)

	assertPanics := func(name string, fn func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s: expected panic on re-entrant mutation during callback, got none", name)
			}
		}()
		fn()
	}

	assertPanics("PlaceOrder inside PlaceOrderWithCallback", func() {
		b.PlaceOrderWithCallback(OrderRequest{ID: 2, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, func(ev Event) {
			if ev.Type == EventTrade {
				b.PlaceOrder(OrderRequest{ID: 3, TraderID: 3, Side: Buy, Type: Limit, Price: 9900, Qty: 100}, nil)
			}
		})
	})

	// After recovering from the panic, the book must not remain locked in callback mode.
	evs := b.CancelOrder(1, nil)
	if len(evs) != 1 || evs[0].Type != EventCancelled {
		t.Fatalf("book should remain usable after recovered re-entrancy panic, got %+v", evs)
	}

	b.PlaceOrder(OrderRequest{ID: 10, TraderID: 1, Side: Buy, Type: Limit, Price: 9900, Qty: 200}, nil)
	assertPanics("CancelOrder inside ModifyOrderWithCallback", func() {
		b.ModifyOrderWithCallback(ModifyRequest{ID: 10, NewQty: 100}, func(ev Event) {
			b.CancelOrder(10, nil)
		})
	})

	assertPanics("Reset inside CancelOrderWithCallback", func() {
		b.CancelOrderWithCallback(10, func(ev Event) {
			b.Reset()
		})
	})
}

func TestModifyOrderQtyConsistency(t *testing.T) {
	runCase := func(t *testing.T, useCallback bool) {
		b, err := New(testCfg())
		if err != nil {
			t.Fatal(err)
		}
		var evs []Event

		// Seed Sell #1 (Qty=500 @ 10000), partially fill 200 -> Leaves=300.
		b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Sell, Type: Limit, Price: 10000, Qty: 500}, nil)
		b.PlaceOrder(OrderRequest{ID: 2, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 200}, nil)

		// In-place qty decrease: 300 -> 200.
		if useCallback {
			evs = evs[:0]
			b.ModifyOrderWithCallback(ModifyRequest{ID: 1, NewQty: 200}, func(ev Event) { evs = append(evs, ev) })
		} else {
			evs = b.ModifyOrder(ModifyRequest{ID: 1, NewQty: 200}, evs[:0])
		}
		if len(evs) != 1 || evs[0].Type != EventModified || evs[0].Qty != 200 || evs[0].Leaves != 200 {
			t.Fatalf("unexpected EventModified: %+v", evs)
		}

		snap, ok := b.GetOrder(1)
		if !ok || snap.Qty != 200 || snap.Leaves != 200 {
			t.Fatalf("after in-place modify: GetOrder(1) = %+v (ok=%v), want Qty=200 Leaves=200", snap, ok)
		}

		// Fill remaining 200 -> EventDone for maker #1 must also report Qty=200 (matching priority-losing modify).
		evs = b.PlaceOrder(OrderRequest{ID: 3, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 200}, evs[:0])
		var makerDone *Event
		for i := range evs {
			if evs[i].Type == EventDone && evs[i].OrderID == 1 {
				makerDone = &evs[i]
				break
			}
		}
		if makerDone == nil || makerDone.Qty != 200 || makerDone.Leaves != 0 {
			t.Fatalf("maker EventDone = %+v, want Qty=200 Leaves=0", makerDone)
		}
	}

	t.Run("Slice", func(t *testing.T) { runCase(t, false) })
	t.Run("Callback", func(t *testing.T) { runCase(t, true) })
}

func TestConfigValidateMaxTicks(t *testing.T) {
	tooWide := []Config{
		{Symbol: "WIDE", TickSize: 1, LotSize: 100, MinPrice: 1, MaxPrice: 2_000_000},
		{Symbol: "MAXINT", TickSize: 1, LotSize: 100, MinPrice: 1, MaxPrice: math.MaxInt64},
	}
	for _, cfg := range tooWide {
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate(%+v) = nil, want error for exceeding MaxTicks", cfg)
		}
	}
}

func TestModifySamePriceIncreaseAllPositions(t *testing.T) {
	run := func(t *testing.T, useCB bool) {
		b, err := New(testCfg())
		if err != nil {
			t.Fatal(err)
		}
		var evs []Event
		modify := func(req ModifyRequest) []Event {
			evs = evs[:0]
			if useCB {
				b.ModifyOrderWithCallback(req, func(ev Event) {
					if ev.Type == EventModified {
						snap, ok := b.GetOrder(req.ID)
						if !ok || snap.Qty != req.NewQty || snap.Leaves != req.NewQty {
							t.Errorf("inside EventModified: GetOrder(%d) = (%+v, %v), want Qty=Leaves=%d",
								req.ID, snap, ok, req.NewQty)
						}
					}
					evs = append(evs, ev)
				})
			} else {
				evs = b.ModifyOrder(req, evs)
			}
			return evs
		}

		// Case 1: Sole order on a level (#1: 100 @ 10000 -> 200 @ 10000).
		b.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, nil)
		origSnap, _ := b.GetOrder(1)
		out := modify(ModifyRequest{ID: 1, NewQty: 200})
		if len(out) != 1 || out[0].Type != EventModified || out[0].Qty != 200 || out[0].Leaves != 200 {
			t.Fatalf("sole order modify: got %+v", out)
		}
		newSnap, ok := b.GetOrder(1)
		if !ok || newSnap.Qty != 200 || newSnap.Leaves != 200 || newSnap.Seq <= origSnap.Seq {
			t.Fatalf("sole order snapshot = %+v (ok=%v), prevSeq=%d", newSnap, ok, origSnap.Seq)
		}
		if _, qty, _ := b.BestBid(); qty != 200 {
			t.Fatalf("BestBid qty = %d, want 200", qty)
		}

		// Add #2 (100) and #3 (100) behind #1 -> FIFO order: [#1(200), #2(100), #3(100)].
		b.PlaceOrder(OrderRequest{ID: 2, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, nil)
		b.PlaceOrder(OrderRequest{ID: 3, TraderID: 3, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, nil)

		// Case 2: Increase head (#1: 200 -> 300) -> FIFO becomes [#2(100), #3(100), #1(300)].
		modify(ModifyRequest{ID: 1, NewQty: 300})
		if _, qty, _ := b.BestBid(); qty != 500 {
			t.Fatalf("after head increase: BestBid qty = %d, want 500", qty)
		}

		// Case 3: Increase middle (#3: 100 -> 200) -> FIFO becomes [#2(100), #1(300), #3(200)].
		modify(ModifyRequest{ID: 3, NewQty: 200})
		if _, qty, _ := b.BestBid(); qty != 600 {
			t.Fatalf("after middle increase: BestBid qty = %d, want 600", qty)
		}

		// Case 4: Increase tail (#3: 200 -> 300) -> FIFO stays [#2(100), #1(300), #3(300)].
		modify(ModifyRequest{ID: 3, NewQty: 300})
		if _, qty, _ := b.BestBid(); qty != 700 {
			t.Fatalf("after tail increase: BestBid qty = %d, want 700", qty)
		}

		// Sweep with Sell 700 @ 10000 and verify exact FIFO fill order: #2, then #1, then #3!
		fillEvs := b.PlaceOrder(OrderRequest{ID: 99, TraderID: 99, Side: Sell, Type: Limit, Price: 10000, Qty: 700}, nil)
		var makers []OrderID
		for _, ev := range fillEvs {
			if ev.Type == EventTrade {
				makers = append(makers, ev.MakerOrderID)
			}
		}
		if len(makers) != 3 || makers[0] != 2 || makers[1] != 1 || makers[2] != 3 {
			t.Fatalf("FIFO fill order = %v, want [2 1 3]", makers)
		}
		if b.OrderCount() != 0 {
			t.Fatalf("OrderCount = %d, want 0", b.OrderCount())
		}
	}

	t.Run("Slice", func(t *testing.T) { run(t, false) })
	t.Run("Callback", func(t *testing.T) { run(t, true) })
}

func BenchmarkBookAddCancelTick25(b *testing.B) {
	cfg := Config{Symbol: "BBCA", TickSize: 25, LotSize: 100, MinPrice: 5000, MaxPrice: 15000, InitialCapacity: 8192}
	bk, err := New(cfg)
	if err != nil {
		b.Fatal(err)
	}
	// Rest an opposite ask so wouldCross inspects a live level.
	bk.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Sell, Type: Limit, Price: 10050, Qty: 100}, nil)
	evBuf := make([]Event, 0, 16)
	var id OrderID = 2
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evBuf = bk.PlaceOrder(OrderRequest{ID: id, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, evBuf[:0])
		evBuf = bk.CancelOrder(id, evBuf[:0])
		id++
	}
}

func BenchmarkBookModifyQtyIncrease(b *testing.B) {
	cfg := Config{Symbol: "BBCA", TickSize: 25, LotSize: 100, MinPrice: 5000, MaxPrice: 15000, InitialCapacity: 8192}
	bk, _ := New(cfg)
	evBuf := make([]Event, 0, 16)
	// Two orders at 10000 so increasing #1 moves it behind #2 (exercising queue tail move).
	bk.PlaceOrder(OrderRequest{ID: 1, TraderID: 1, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, evBuf[:0])
	bk.PlaceOrder(OrderRequest{ID: 2, TraderID: 2, Side: Buy, Type: Limit, Price: 10000, Qty: 100}, evBuf[:0])
	bk.PlaceOrder(OrderRequest{ID: 3, TraderID: 3, Side: Sell, Type: Limit, Price: 10050, Qty: 100}, evBuf[:0])
	qty := int64(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qty += 100
		if qty > 1_000_000_000 {
			b.StopTimer()
			bk.ModifyOrder(ModifyRequest{ID: 1, NewQty: 100}, evBuf[:0])
			qty = 200
			b.StartTimer()
		}
		evBuf = bk.ModifyOrder(ModifyRequest{ID: 1, NewQty: qty}, evBuf[:0])
	}
}
