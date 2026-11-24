package orderbook

import "math/bits"

const defaultInitialCapacity = 4096

// Compile-time check that *Book implements the exported OrderBook interface.
var _ OrderBook = (*Book)(nil)

// arenaOrder is exactly 64 bytes (one x86-64 L1 cache line) and contains zero
// Go pointers, so Go's GC marks the backing []arenaOrder slice as noscan.
type arenaOrder struct {
	ID       OrderID   // 8 B
	TraderID TraderID  // 8 B
	Price    int64     // 8 B
	Qty      int64     // 8 B
	Leaves   int64     // 8 B
	Seq      uint64    // 8 B
	PrevIdx  uint32    // 4 B (0 = nil)
	NextIdx  uint32    // 4 B (0 = nil; also intrusive free-list link)
	TickIdx  uint32    // 4 B (direct ladder slot index)
	Side     Side      // 1 B
	Type     OrderType // 1 B
	_        [2]byte   // pad to 64 B
}

// Book is a single-threaded, deterministic, zero-allocation matching engine for
// one ticker. It combines a tick-indexed Ladder, intrusive LinkedList FIFO queues,
// a 64-byte pointer-free (GC noscan) Order Arena, and a flat Fibonacci-hashed OrderMap.
//
// Book is NOT safe for concurrent use; run one *Book per ticker on a dedicated goroutine.
type Book struct {
	cfg        Config
	bids       ladder
	asks       ladder
	orders     orderMap
	arena      []arenaOrder // index 0 is reserved nil sentinel
	freeHead   uint32
	highMark   uint32
	orderSeq   uint64
	eventSeq   uint64
	tradeSeq   uint64
	inCallback bool
}

func (b *Book) assertNotInCallback() {
	if b.inCallback {
		panic("orderbook: re-entrant mutation during callback is not allowed")
	}
}

// New constructs a new single-ticker *Book validated against cfg.
func New(cfg Config) (*Book, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	initCap := cfg.InitialCapacity
	if initCap <= 0 {
		initCap = defaultInitialCapacity
	}
	return &Book{
		cfg:      cfg,
		bids:     newLadder(Buy, cfg),
		asks:     newLadder(Sell, cfg),
		orders:   newOrderMap(initCap),
		arena:    make([]arenaOrder, initCap+1), // slot 0 = nil sentinel
		highMark: 1,
	}, nil
}

func (b *Book) allocOrder() uint32 {
	if b.freeHead != 0 {
		idx := b.freeHead
		b.freeHead = b.arena[idx].NextIdx
		return idx
	}
	if int(b.highMark) >= len(b.arena) {
		newArena := make([]arenaOrder, len(b.arena)*2)
		copy(newArena, b.arena)
		b.arena = newArena
	}
	idx := b.highMark
	b.highMark++
	return idx
}

func (b *Book) freeOrder(idx uint32) {
	b.arena[idx].NextIdx = b.freeHead
	b.freeHead = idx
}

func (b *Book) validatePrice(price int64) (uint32, RejectReason) {
	if price <= 0 {
		return 0, RejectInvalidPrice
	}
	if price < b.cfg.MinPrice || price > b.cfg.MaxPrice {
		return 0, RejectOutOfPriceBand
	}
	delta := price - b.cfg.MinPrice
	if b.cfg.TickSize == 1 {
		return uint32(delta), RejectNone
	}
	tickIdx := delta / b.cfg.TickSize
	if tickIdx*b.cfg.TickSize != delta {
		return 0, RejectOffTickPrice
	}
	return uint32(tickIdx), RejectNone
}

func (b *Book) validateQty(qty int64) RejectReason {
	if qty <= 0 {
		return RejectInvalidQty
	}
	if b.cfg.LotSize != 1 && qty%b.cfg.LotSize != 0 {
		return RejectNonLotQty
	}
	return RejectNone
}

func (b *Book) validateOrder(req OrderRequest) (uint32, RejectReason) {
	if req.ID == 0 {
		return 0, RejectInvalidOrderID
	}
	if req.Side != Buy && req.Side != Sell {
		return 0, RejectInvalidSide
	}
	if req.Type != Limit && req.Type != Market {
		return 0, RejectInvalidType
	}
	if req.TimeInForce > PostOnly {
		return 0, RejectInvalidTimeInForce
	}
	if req.Type == Market && req.TimeInForce == PostOnly {
		return 0, RejectInvalidTimeInForce
	}
	if r := b.validateQty(req.Qty); r != RejectNone {
		return 0, r
	}
	var tickIdx uint32
	if req.Type == Limit {
		var r RejectReason
		tickIdx, r = b.validatePrice(req.Price)
		if r != RejectNone {
			return 0, r
		}
	}
	if b.orders.Has(req.ID) {
		return 0, RejectDuplicateOrderID
	}
	return tickIdx, RejectNone
}

func wouldCross(opp *ladder, side Side, price int64) bool {
	if opp.bestIdx < 0 {
		return false
	}
	bestPrice := opp.slots[opp.bestIdx].price
	if side == Buy {
		return price >= bestPrice
	}
	return price <= bestPrice
}

// canFillFOK checks without mutating the book whether at least needQty can be
// filled immediately against opp without hitting a self-trade (STP) blocker.
func (b *Book) canFillFOK(opp *ladder, side Side, oType OrderType, limitPrice, needQty int64, traderID TraderID) bool {
	if opp.bestIdx < 0 {
		return false
	}
	var avail int64

	checkLevel := func(lvl *priceLevel) (done bool, ok bool) {
		if oType == Limit {
			if side == Buy && lvl.price > limitPrice {
				return true, false
			}
			if side == Sell && lvl.price < limitPrice {
				return true, false
			}
		}
		if traderID == 0 {
			avail += lvl.totalQty
			if avail >= needQty {
				return true, true
			}
			return false, false
		}
		for idx := lvl.headIdx; idx != 0; idx = b.arena[idx].NextIdx {
			maker := &b.arena[idx]
			if maker.TraderID == traderID {
				// STP Cancel-Newest would halt matching before reaching needQty.
				return true, false
			}
			avail += maker.Leaves
			if avail >= needQty {
				return true, true
			}
		}
		return false, false
	}

	if opp.side == Buy {
		for w := opp.bestIdx >> 6; w >= 0; w-- {
			word := opp.active[w]
			for word != 0 {
				bit := 63 - bits.LeadingZeros64(word)
				word &^= uint64(1) << uint(bit)
				done, ok := checkLevel(&opp.slots[(w<<6)|bit])
				if done {
					return ok
				}
			}
		}
	} else {
		for w := opp.bestIdx >> 6; w < len(opp.active); w++ {
			word := opp.active[w]
			for word != 0 {
				bit := bits.TrailingZeros64(word)
				word &^= uint64(1) << uint(bit)
				done, ok := checkLevel(&opp.slots[(w<<6)|bit])
				if done {
					return ok
				}
			}
		}
	}
	return avail >= needQty
}

// PlaceOrder validates, matches, and (for unfilled GTC/PostOnly limit remainders) rests req.
// Supports GTC, IOC, FOK, and PostOnly TimeInForce.
// Produced events are appended to dst and the extended slice is returned.
func (b *Book) PlaceOrder(req OrderRequest, dst []Event) []Event {
	b.assertNotInCallback()
	tickIdx, reason := b.validateOrder(req)
	if reason != RejectNone {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        req.Price,
			Qty:          req.Qty,
			RejectReason: reason,
		})
	}

	price := req.Price
	if req.Type == Market {
		price = 0
	}

	b.orderSeq++
	orderSeq := b.orderSeq
	leaves := req.Qty

	b.eventSeq++
	dst = append(dst, Event{
		Seq:       b.eventSeq,
		Type:      EventAccepted,
		OrderID:   req.ID,
		TraderID:  req.TraderID,
		Side:      req.Side,
		OrderType: req.Type,
		Price:     price,
		Qty:       req.Qty,
		Leaves:    leaves,
	})

	var own, opp *ladder
	if req.Side == Buy {
		own, opp = &b.bids, &b.asks
	} else {
		own, opp = &b.asks, &b.bids
	}

	// PostOnly: cancel immediately if order would cross the spread.
	if req.TimeInForce == PostOnly {
		if wouldCross(opp, req.Side, price) {
			b.eventSeq++
			return append(dst, Event{
				Seq:          b.eventSeq,
				Type:         EventCancelled,
				OrderID:      req.ID,
				TraderID:     req.TraderID,
				Side:         req.Side,
				OrderType:    req.Type,
				Price:        price,
				Qty:          leaves,
				Leaves:       0,
				CancelReason: CancelPostOnly,
			})
		}
		b.restLimit(own, tickIdx, req.ID, req.TraderID, price, req.Qty, leaves, orderSeq)
		return dst
	}

	// FOK: pre-check atomic fillability before executing any trades.
	if req.TimeInForce == FOK && !b.canFillFOK(opp, req.Side, req.Type, price, leaves, req.TraderID) {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventCancelled,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        price,
			Qty:          leaves,
			Leaves:       0,
			CancelReason: CancelFOK,
		})
	}

	if req.Type == Market || wouldCross(opp, req.Side, price) {
		var stpTriggered bool
		leaves, stpTriggered, dst = b.matchSlice(opp, req.ID, req.TraderID, req.Side, req.Type, price, leaves, dst)
		if stpTriggered {
			return dst
		}

		if leaves == 0 {
			b.eventSeq++
			return append(dst, Event{
				Seq:       b.eventSeq,
				Type:      EventDone,
				OrderID:   req.ID,
				TraderID:  req.TraderID,
				Side:      req.Side,
				OrderType: req.Type,
				Price:     price,
				Qty:       req.Qty,
				Leaves:    0,
			})
		}
	}

	if req.Type == Market {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventCancelled,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        0,
			Qty:          leaves,
			Leaves:       0,
			CancelReason: CancelMarketUnfilled,
		})
	}

	if req.TimeInForce == IOC {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventCancelled,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        price,
			Qty:          leaves,
			Leaves:       0,
			CancelReason: CancelIOC,
		})
	}

	b.restLimit(own, tickIdx, req.ID, req.TraderID, price, req.Qty, leaves, orderSeq)
	return dst
}

// PlaceOrderWithCallback streams every produced Event directly to fn with zero
// intermediate slice buffering. If fn is nil, events are discarded.
func (b *Book) PlaceOrderWithCallback(req OrderRequest, fn func(Event)) {
	b.assertNotInCallback()
	if fn != nil {
		b.inCallback = true
		defer func() { b.inCallback = false }()
	}
	emit := func(ev Event) {
		if fn != nil {
			fn(ev)
		}
	}

	tickIdx, reason := b.validateOrder(req)
	if reason != RejectNone {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        req.Price,
			Qty:          req.Qty,
			RejectReason: reason,
		})
		return
	}

	price := req.Price
	if req.Type == Market {
		price = 0
	}

	b.orderSeq++
	orderSeq := b.orderSeq
	leaves := req.Qty

	b.eventSeq++
	emit(Event{
		Seq:       b.eventSeq,
		Type:      EventAccepted,
		OrderID:   req.ID,
		TraderID:  req.TraderID,
		Side:      req.Side,
		OrderType: req.Type,
		Price:     price,
		Qty:       req.Qty,
		Leaves:    leaves,
	})

	var own, opp *ladder
	if req.Side == Buy {
		own, opp = &b.bids, &b.asks
	} else {
		own, opp = &b.asks, &b.bids
	}

	if req.TimeInForce == PostOnly {
		if wouldCross(opp, req.Side, price) {
			b.eventSeq++
			emit(Event{
				Seq:          b.eventSeq,
				Type:         EventCancelled,
				OrderID:      req.ID,
				TraderID:     req.TraderID,
				Side:         req.Side,
				OrderType:    req.Type,
				Price:        price,
				Qty:          leaves,
				Leaves:       0,
				CancelReason: CancelPostOnly,
			})
			return
		}
		b.restLimit(own, tickIdx, req.ID, req.TraderID, price, req.Qty, leaves, orderSeq)
		return
	}

	if req.TimeInForce == FOK && !b.canFillFOK(opp, req.Side, req.Type, price, leaves, req.TraderID) {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventCancelled,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        price,
			Qty:          leaves,
			Leaves:       0,
			CancelReason: CancelFOK,
		})
		return
	}

	if req.Type == Market || wouldCross(opp, req.Side, price) {
		var stpTriggered bool
		leaves, stpTriggered = b.matchCallback(opp, req.ID, req.TraderID, req.Side, req.Type, price, leaves, fn)
		if stpTriggered {
			return
		}

		if leaves == 0 {
			b.eventSeq++
			emit(Event{
				Seq:       b.eventSeq,
				Type:      EventDone,
				OrderID:   req.ID,
				TraderID:  req.TraderID,
				Side:      req.Side,
				OrderType: req.Type,
				Price:     price,
				Qty:       req.Qty,
				Leaves:    0,
			})
			return
		}
	}

	if req.Type == Market {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventCancelled,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        0,
			Qty:          leaves,
			Leaves:       0,
			CancelReason: CancelMarketUnfilled,
		})
		return
	}

	if req.TimeInForce == IOC {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventCancelled,
			OrderID:      req.ID,
			TraderID:     req.TraderID,
			Side:         req.Side,
			OrderType:    req.Type,
			Price:        price,
			Qty:          leaves,
			Leaves:       0,
			CancelReason: CancelIOC,
		})
		return
	}

	b.restLimit(own, tickIdx, req.ID, req.TraderID, price, req.Qty, leaves, orderSeq)
}

func (b *Book) matchSlice(opp *ladder, orderID OrderID, traderID TraderID, side Side, oType OrderType, price, leaves int64, dst []Event) (int64, bool, []Event) {
	stpTriggered := false
	for leaves > 0 {
		bestIdx := opp.bestIdx
		if bestIdx < 0 {
			break
		}
		bestLvl := &opp.slots[bestIdx]
		if oType == Limit {
			if side == Buy {
				if bestLvl.price > price {
					break
				}
			} else {
				if bestLvl.price < price {
					break
				}
			}
		}

		for leaves > 0 && bestLvl.headIdx != 0 {
			makerIdx := bestLvl.headIdx
			maker := &b.arena[makerIdx]

			if traderID != 0 && maker.TraderID == traderID {
				stpTriggered = true
				b.eventSeq++
				dst = append(dst, Event{
					Seq:          b.eventSeq,
					Type:         EventCancelled,
					OrderID:      orderID,
					TraderID:     traderID,
					Side:         side,
					OrderType:    oType,
					Price:        price,
					Qty:          leaves,
					Leaves:       0,
					CancelReason: CancelSTP,
				})
				break
			}

			matchQty := leaves
			if maker.Leaves < matchQty {
				matchQty = maker.Leaves
			}

			leaves -= matchQty
			maker.Leaves -= matchQty
			bestLvl.totalQty -= matchQty

			b.tradeSeq++
			b.eventSeq++
			dst = append(dst, Event{
				Seq:           b.eventSeq,
				Type:          EventTrade,
				TradeID:       b.tradeSeq,
				OrderID:       orderID,
				TraderID:      traderID,
				MakerOrderID:  maker.ID,
				MakerTraderID: maker.TraderID,
				Side:          side,
				OrderType:     oType,
				Price:         maker.Price,
				Qty:           matchQty,
				Leaves:        leaves,
				MakerLeaves:   maker.Leaves,
			})

			if maker.Leaves == 0 {
				nextIdx := maker.NextIdx
				bestLvl.headIdx = nextIdx
				if nextIdx != 0 {
					b.arena[nextIdx].PrevIdx = 0
				} else {
					bestLvl.tailIdx = 0
				}
				bestLvl.orderCount--
				b.orders.Delete(maker.ID)

				b.eventSeq++
				dst = append(dst, Event{
					Seq:       b.eventSeq,
					Type:      EventDone,
					OrderID:   maker.ID,
					TraderID:  maker.TraderID,
					Side:      maker.Side,
					OrderType: maker.Type,
					Price:     maker.Price,
					Qty:       maker.Qty,
					Leaves:    0,
				})
				b.freeOrder(makerIdx)
			}
		}

		if bestLvl.orderCount == 0 {
			opp.removeAt(bestIdx)
		}
		if stpTriggered {
			break
		}
	}
	return leaves, stpTriggered, dst
}

func (b *Book) matchCallback(opp *ladder, orderID OrderID, traderID TraderID, side Side, oType OrderType, price, leaves int64, fn func(Event)) (int64, bool) {
	stpTriggered := false
	for leaves > 0 {
		bestIdx := opp.bestIdx
		if bestIdx < 0 {
			break
		}
		bestLvl := &opp.slots[bestIdx]
		if oType == Limit {
			if side == Buy {
				if bestLvl.price > price {
					break
				}
			} else {
				if bestLvl.price < price {
					break
				}
			}
		}

		for leaves > 0 && bestLvl.headIdx != 0 {
			makerIdx := bestLvl.headIdx
			maker := &b.arena[makerIdx]

			if traderID != 0 && maker.TraderID == traderID {
				stpTriggered = true
				b.eventSeq++
				if fn != nil {
					fn(Event{
						Seq:          b.eventSeq,
						Type:         EventCancelled,
						OrderID:      orderID,
						TraderID:     traderID,
						Side:         side,
						OrderType:    oType,
						Price:        price,
						Qty:          leaves,
						Leaves:       0,
						CancelReason: CancelSTP,
					})
				}
				break
			}

			matchQty := leaves
			if maker.Leaves < matchQty {
				matchQty = maker.Leaves
			}

			leaves -= matchQty
			maker.Leaves -= matchQty
			bestLvl.totalQty -= matchQty

			makerID := maker.ID
			makerTraderID := maker.TraderID
			makerSide := maker.Side
			makerType := maker.Type
			makerPrice := maker.Price
			makerQty := maker.Qty
			makerLeaves := maker.Leaves

			// Commit all book mutations before invoking user callbacks so read-only
			// queries inside fn always observe a consistent book state.
			if makerLeaves == 0 {
				nextIdx := maker.NextIdx
				bestLvl.headIdx = nextIdx
				if nextIdx != 0 {
					b.arena[nextIdx].PrevIdx = 0
				} else {
					bestLvl.tailIdx = 0
				}
				bestLvl.orderCount--
				b.orders.Delete(makerID)
				b.freeOrder(makerIdx)
				if bestLvl.orderCount == 0 {
					opp.removeAt(bestIdx)
				}
			}

			b.tradeSeq++
			b.eventSeq++
			if fn != nil {
				fn(Event{
					Seq:           b.eventSeq,
					Type:          EventTrade,
					TradeID:       b.tradeSeq,
					OrderID:       orderID,
					TraderID:      traderID,
					MakerOrderID:  makerID,
					MakerTraderID: makerTraderID,
					Side:          side,
					OrderType:     oType,
					Price:         makerPrice,
					Qty:           matchQty,
					Leaves:        leaves,
					MakerLeaves:   makerLeaves,
				})
			}

			if makerLeaves == 0 {
				b.eventSeq++
				if fn != nil {
					fn(Event{
						Seq:       b.eventSeq,
						Type:      EventDone,
						OrderID:   makerID,
						TraderID:  makerTraderID,
						Side:      makerSide,
						OrderType: makerType,
						Price:     makerPrice,
						Qty:       makerQty,
						Leaves:    0,
					})
				}
			}
		}

		if stpTriggered {
			break
		}
	}
	return leaves, stpTriggered
}

// ModifyOrder amends a resting limit order's price and/or remaining quantity (Leaves).
//   - Same price & NewQty < current Leaves: modifies in-place, preserving FIFO priority.
//   - Price change or NewQty > current Leaves: loses priority, re-matches if crossing, and rests remainder.
func (b *Book) ModifyOrder(req ModifyRequest, dst []Event) []Event {
	b.assertNotInCallback()
	idx, ok := b.orders.Get(req.ID)
	if !ok {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			RejectReason: RejectOrderNotFound,
		})
	}

	o := &b.arena[idx]
	newPrice := req.NewPrice
	newTickIdx := o.TickIdx
	if newPrice == 0 {
		newPrice = o.Price
	}
	newLeaves := req.NewQty
	if newLeaves == 0 {
		newLeaves = o.Leaves
	}

	if newPrice != o.Price {
		var r RejectReason
		newTickIdx, r = b.validatePrice(newPrice)
		if r != RejectNone {
			b.eventSeq++
			return append(dst, Event{
				Seq:          b.eventSeq,
				Type:         EventRejected,
				OrderID:      req.ID,
				TraderID:     o.TraderID,
				Side:         o.Side,
				OrderType:    Limit,
				Price:        newPrice,
				Qty:          newLeaves,
				RejectReason: r,
			})
		}
	}
	if r := b.validateQty(newLeaves); r != RejectNone {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			TraderID:     o.TraderID,
			Side:         o.Side,
			OrderType:    Limit,
			Price:        newPrice,
			Qty:          newLeaves,
			RejectReason: r,
		})
	}
	if newPrice == o.Price && newLeaves == o.Leaves {
		b.eventSeq++
		return append(dst, Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			TraderID:     o.TraderID,
			Side:         o.Side,
			OrderType:    Limit,
			Price:        newPrice,
			Qty:          newLeaves,
			RejectReason: RejectNoChange,
		})
	}

	var own, opp *ladder
	if o.Side == Buy {
		own, opp = &b.bids, &b.asks
	} else {
		own, opp = &b.asks, &b.bids
	}

	// Fast paths for same-price modifications (already uncrossed, no map/arena churn).
	if newPrice == o.Price {
		lvl := &own.slots[o.TickIdx]
		if newLeaves < o.Leaves {
			// Quantity reduced -> keep FIFO time priority in-place.
			lvl.totalQty -= o.Leaves - newLeaves
			o.Qty = newLeaves
			o.Leaves = newLeaves
		} else {
			// Quantity increased -> lose FIFO time priority by moving to tail of the same level.
			lvl.totalQty += newLeaves - o.Leaves
			b.orderSeq++
			o.Qty = newLeaves
			o.Leaves = newLeaves
			o.Seq = b.orderSeq
			if lvl.tailIdx != idx {
				prevIdx := o.PrevIdx
				nextIdx := o.NextIdx
				if prevIdx != 0 {
					b.arena[prevIdx].NextIdx = nextIdx
				} else {
					lvl.headIdx = nextIdx
				}
				b.arena[nextIdx].PrevIdx = prevIdx

				oldTail := lvl.tailIdx
				b.arena[oldTail].NextIdx = idx
				o.PrevIdx = oldTail
				o.NextIdx = 0
				lvl.tailIdx = idx
			}
		}

		b.eventSeq++
		return append(dst, Event{
			Seq:       b.eventSeq,
			Type:      EventModified,
			OrderID:   o.ID,
			TraderID:  o.TraderID,
			Side:      o.Side,
			OrderType: Limit,
			Price:     o.Price,
			Qty:       newLeaves,
			Leaves:    newLeaves,
		})
	}

	// Price changed -> lose time priority, unlink, re-match if crossing.
	traderID := o.TraderID
	side := o.Side
	b.unlinkAndFree(own, idx)
	b.orders.Delete(req.ID)

	b.orderSeq++
	orderSeq := b.orderSeq

	b.eventSeq++
	dst = append(dst, Event{
		Seq:       b.eventSeq,
		Type:      EventModified,
		OrderID:   req.ID,
		TraderID:  traderID,
		Side:      side,
		OrderType: Limit,
		Price:     newPrice,
		Qty:       newLeaves,
		Leaves:    newLeaves,
	})

	leaves := newLeaves
	if wouldCross(opp, side, newPrice) {
		var stpTriggered bool
		leaves, stpTriggered, dst = b.matchSlice(opp, req.ID, traderID, side, Limit, newPrice, newLeaves, dst)
		if stpTriggered {
			return dst
		}
		if leaves == 0 {
			b.eventSeq++
			return append(dst, Event{
				Seq:       b.eventSeq,
				Type:      EventDone,
				OrderID:   req.ID,
				TraderID:  traderID,
				Side:      side,
				OrderType: Limit,
				Price:     newPrice,
				Qty:       newLeaves,
				Leaves:    0,
			})
		}
	}

	b.restLimit(own, newTickIdx, req.ID, traderID, newPrice, newLeaves, leaves, orderSeq)
	return dst
}

// ModifyOrderWithCallback amends a resting limit order and streams events to fn.
func (b *Book) ModifyOrderWithCallback(req ModifyRequest, fn func(Event)) {
	b.assertNotInCallback()
	if fn != nil {
		b.inCallback = true
		defer func() { b.inCallback = false }()
	}
	emit := func(ev Event) {
		if fn != nil {
			fn(ev)
		}
	}

	idx, ok := b.orders.Get(req.ID)
	if !ok {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			RejectReason: RejectOrderNotFound,
		})
		return
	}

	o := &b.arena[idx]
	newPrice := req.NewPrice
	newTickIdx := o.TickIdx
	if newPrice == 0 {
		newPrice = o.Price
	}
	newLeaves := req.NewQty
	if newLeaves == 0 {
		newLeaves = o.Leaves
	}

	if newPrice != o.Price {
		var r RejectReason
		newTickIdx, r = b.validatePrice(newPrice)
		if r != RejectNone {
			b.eventSeq++
			emit(Event{
				Seq:          b.eventSeq,
				Type:         EventRejected,
				OrderID:      req.ID,
				TraderID:     o.TraderID,
				Side:         o.Side,
				OrderType:    Limit,
				Price:        newPrice,
				Qty:          newLeaves,
				RejectReason: r,
			})
			return
		}
	}
	if r := b.validateQty(newLeaves); r != RejectNone {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			TraderID:     o.TraderID,
			Side:         o.Side,
			OrderType:    Limit,
			Price:        newPrice,
			Qty:          newLeaves,
			RejectReason: r,
		})
		return
	}
	if newPrice == o.Price && newLeaves == o.Leaves {
		b.eventSeq++
		emit(Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      req.ID,
			TraderID:     o.TraderID,
			Side:         o.Side,
			OrderType:    Limit,
			Price:        newPrice,
			Qty:          newLeaves,
			RejectReason: RejectNoChange,
		})
		return
	}

	var own, opp *ladder
	if o.Side == Buy {
		own, opp = &b.bids, &b.asks
	} else {
		own, opp = &b.asks, &b.bids
	}

	if newPrice == o.Price {
		lvl := &own.slots[o.TickIdx]
		if newLeaves < o.Leaves {
			lvl.totalQty -= o.Leaves - newLeaves
			o.Qty = newLeaves
			o.Leaves = newLeaves
		} else {
			lvl.totalQty += newLeaves - o.Leaves
			b.orderSeq++
			o.Qty = newLeaves
			o.Leaves = newLeaves
			o.Seq = b.orderSeq
			if lvl.tailIdx != idx {
				prevIdx := o.PrevIdx
				nextIdx := o.NextIdx
				if prevIdx != 0 {
					b.arena[prevIdx].NextIdx = nextIdx
				} else {
					lvl.headIdx = nextIdx
				}
				b.arena[nextIdx].PrevIdx = prevIdx

				oldTail := lvl.tailIdx
				b.arena[oldTail].NextIdx = idx
				o.PrevIdx = oldTail
				o.NextIdx = 0
				lvl.tailIdx = idx
			}
		}

		b.eventSeq++
		emit(Event{
			Seq:       b.eventSeq,
			Type:      EventModified,
			OrderID:   o.ID,
			TraderID:  o.TraderID,
			Side:      o.Side,
			OrderType: Limit,
			Price:     o.Price,
			Qty:       newLeaves,
			Leaves:    newLeaves,
		})
		return
	}

	traderID := o.TraderID
	side := o.Side
	b.unlinkAndFree(own, idx)
	b.orders.Delete(req.ID)

	b.orderSeq++
	orderSeq := b.orderSeq

	// For non-crossing modifies, rest the amended order before emitting EventModified
	// so callbacks observe the resting order in a consistent state.
	if !wouldCross(opp, side, newPrice) {
		b.restLimit(own, newTickIdx, req.ID, traderID, newPrice, newLeaves, newLeaves, orderSeq)
		b.eventSeq++
		emit(Event{
			Seq:       b.eventSeq,
			Type:      EventModified,
			OrderID:   req.ID,
			TraderID:  traderID,
			Side:      side,
			OrderType: Limit,
			Price:     newPrice,
			Qty:       newLeaves,
			Leaves:    newLeaves,
		})
		return
	}

	b.eventSeq++
	emit(Event{
		Seq:       b.eventSeq,
		Type:      EventModified,
		OrderID:   req.ID,
		TraderID:  traderID,
		Side:      side,
		OrderType: Limit,
		Price:     newPrice,
		Qty:       newLeaves,
		Leaves:    newLeaves,
	})

	leaves, stpTriggered := b.matchCallback(opp, req.ID, traderID, side, Limit, newPrice, newLeaves, fn)
	if stpTriggered {
		return
	}
	if leaves == 0 {
		b.eventSeq++
		emit(Event{
			Seq:       b.eventSeq,
			Type:      EventDone,
			OrderID:   req.ID,
			TraderID:  traderID,
			Side:      side,
			OrderType: Limit,
			Price:     newPrice,
			Qty:       newLeaves,
			Leaves:    0,
		})
		return
	}

	b.restLimit(own, newTickIdx, req.ID, traderID, newPrice, newLeaves, leaves, orderSeq)
}

func (b *Book) unlinkAndFree(side *ladder, idx uint32) {
	o := &b.arena[idx]
	lvl := &side.slots[o.TickIdx]

	prevIdx := o.PrevIdx
	nextIdx := o.NextIdx
	if prevIdx != 0 {
		b.arena[prevIdx].NextIdx = nextIdx
	} else {
		lvl.headIdx = nextIdx
	}
	if nextIdx != 0 {
		b.arena[nextIdx].PrevIdx = prevIdx
	} else {
		lvl.tailIdx = prevIdx
	}

	lvl.totalQty -= o.Leaves
	lvl.orderCount--
	if lvl.orderCount == 0 {
		side.removeAt(int(o.TickIdx))
	}
	b.freeOrder(idx)
}

func (b *Book) restLimit(own *ladder, tickIdx uint32, id OrderID, traderID TraderID, price, qty, leaves int64, orderSeq uint64) {
	idx := b.allocOrder()
	lvl := own.getOrCreateAt(tickIdx, price)
	tailIdx := lvl.tailIdx

	b.arena[idx] = arenaOrder{
		ID:       id,
		TraderID: traderID,
		Price:    price,
		Qty:      qty,
		Leaves:   leaves,
		Seq:      orderSeq,
		PrevIdx:  tailIdx,
		NextIdx:  0,
		TickIdx:  tickIdx,
		Side:     own.side,
		Type:     Limit,
	}

	if tailIdx != 0 {
		b.arena[tailIdx].NextIdx = idx
	} else {
		lvl.headIdx = idx
	}
	lvl.tailIdx = idx
	lvl.totalQty += leaves
	lvl.orderCount++

	b.orders.Set(id, idx)
}

func (b *Book) cancelInto(id OrderID, out *Event) {
	idx, ok := b.orders.Delete(id)
	if !ok {
		b.eventSeq++
		*out = Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      id,
			RejectReason: RejectOrderNotFound,
		}
		return
	}

	o := &b.arena[idx]
	var side *ladder
	if o.Side == Buy {
		side = &b.bids
	} else {
		side = &b.asks
	}
	b.eventSeq++
	*out = Event{
		Seq:          b.eventSeq,
		Type:         EventCancelled,
		OrderID:      o.ID,
		TraderID:     o.TraderID,
		Side:         o.Side,
		OrderType:    o.Type,
		Price:        o.Price,
		Qty:          o.Leaves,
		Leaves:       0,
		CancelReason: CancelUser,
	}

	b.unlinkAndFree(side, idx)
}

// CancelOrder cancels an active resting order by ID, appending either
// EventCancelled(CancelUser) or EventRejected(RejectOrderNotFound) to dst.
func (b *Book) CancelOrder(id OrderID, dst []Event) []Event {
	b.assertNotInCallback()
	n := len(dst)
	if n < cap(dst) {
		dst = dst[:n+1]
	} else {
		dst = append(dst, Event{})
	}
	b.cancelInto(id, &dst[n])
	return dst
}

// CancelOrderWithCallback cancels an active resting order by ID and streams the
// resulting event directly to fn.
func (b *Book) CancelOrderWithCallback(id OrderID, fn func(Event)) {
	b.assertNotInCallback()
	var ev Event
	b.cancelInto(id, &ev)
	if fn != nil {
		b.inCallback = true
		defer func() { b.inCallback = false }()
		fn(ev)
	}
}

// BestBid returns the highest bid price and its aggregated quantity.
func (b *Book) BestBid() (price, qty int64, ok bool) {
	if b.bids.bestIdx < 0 {
		return 0, 0, false
	}
	lvl := &b.bids.slots[b.bids.bestIdx]
	return lvl.price, lvl.totalQty, true
}

// BestAsk returns the lowest ask price and its aggregated quantity.
func (b *Book) BestAsk() (price, qty int64, ok bool) {
	if b.asks.bestIdx < 0 {
		return 0, 0, false
	}
	lvl := &b.asks.slots[b.asks.bestIdx]
	return lvl.price, lvl.totalQty, true
}

// TopOfBook returns the current best bid and best ask in a single value struct.
func (b *Book) TopOfBook() TopOfBook {
	var top TopOfBook
	if b.bids.bestIdx >= 0 {
		lvl := &b.bids.slots[b.bids.bestIdx]
		top.BidPrice = lvl.price
		top.BidQty = lvl.totalQty
		top.HasBid = true
	}
	if b.asks.bestIdx >= 0 {
		lvl := &b.asks.slots[b.asks.bestIdx]
		top.AskPrice = lvl.price
		top.AskQty = lvl.totalQty
		top.HasAsk = true
	}
	return top
}

// Depth writes up to n best bid and ask price levels into bidsDst[:0] and asksDst[:0].
func (b *Book) Depth(n int, bidsDst, asksDst []DepthLevel) (bids, asks []DepthLevel) {
	return b.bids.topN(n, bidsDst[:0]), b.asks.topN(n, asksDst[:0])
}

// GetOrder returns a read-only value snapshot of a resting order by ID.
func (b *Book) GetOrder(id OrderID) (OrderSnapshot, bool) {
	idx, ok := b.orders.Get(id)
	if !ok {
		return OrderSnapshot{}, false
	}
	o := &b.arena[idx]
	return OrderSnapshot{
		ID:       o.ID,
		TraderID: o.TraderID,
		Side:     o.Side,
		Type:     o.Type,
		Price:    o.Price,
		Qty:      o.Qty,
		Leaves:   o.Leaves,
		Seq:      o.Seq,
	}, true
}

// HasOrder reports whether an order with id is currently resting on the book.
func (b *Book) HasOrder(id OrderID) bool {
	return b.orders.Has(id)
}

// OrderCount returns the number of currently resting orders.
func (b *Book) OrderCount() int {
	return b.orders.count
}

// Stats returns deterministic sequence counters and current book sizes.
func (b *Book) Stats() Stats {
	return Stats{
		OrderSeq:     b.orderSeq,
		EventSeq:     b.eventSeq,
		TradeSeq:     b.tradeSeq,
		ActiveOrders: b.orders.count,
		BidLevels:    b.bids.count,
		AskLevels:    b.asks.count,
	}
}

// Config returns the immutable configuration of this book.
func (b *Book) Config() Config {
	return b.cfg
}

// Reset clears all resting orders and resets sequence counters while retaining
// pre-allocated arena, ladder, and hash-map storage.
func (b *Book) Reset() {
	b.assertNotInCallback()
	b.bids.reset()
	b.asks.reset()
	b.orders.Clear()
	b.freeHead = 0
	b.highMark = 1
	b.orderSeq = 0
	b.eventSeq = 0
	b.tradeSeq = 0
}
