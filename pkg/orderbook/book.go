package orderbook

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
	cfg      Config
	bids     ladder
	asks     ladder
	orders   orderMap
	arena    []arenaOrder // index 0 is reserved nil sentinel
	freeHead uint32
	highMark uint32
	orderSeq uint64
	eventSeq uint64
	tradeSeq uint64
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

func (b *Book) validateOrder(req OrderRequest) RejectReason {
	if req.ID == 0 {
		return RejectInvalidOrderID
	}
	if req.Side != Buy && req.Side != Sell {
		return RejectInvalidSide
	}
	if req.Type != Limit && req.Type != Market {
		return RejectInvalidType
	}
	if req.Qty <= 0 {
		return RejectInvalidQty
	}
	if b.cfg.LotSize != 1 && req.Qty%b.cfg.LotSize != 0 {
		return RejectNonLotQty
	}
	if req.Type == Limit {
		if req.Price <= 0 {
			return RejectInvalidPrice
		}
		if req.Price < b.cfg.MinPrice || req.Price > b.cfg.MaxPrice {
			return RejectOutOfPriceBand
		}
		if b.cfg.TickSize != 1 && (req.Price-b.cfg.MinPrice)%b.cfg.TickSize != 0 {
			return RejectOffTickPrice
		}
	}
	if b.orders.Has(req.ID) {
		return RejectDuplicateOrderID
	}
	return RejectNone
}

// PlaceOrder validates, matches, and (for unfilled limit remainders) rests req.
// Produced events are appended to dst and the extended slice is returned.
// Reuse dst[:0] across calls for 0 allocs/op, or pass nil for convenience.
func (b *Book) PlaceOrder(req OrderRequest, dst []Event) []Event {
	if reason := b.validateOrder(req); reason != RejectNone {
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

	stpTriggered := false
	for leaves > 0 {
		bestIdx := opp.bestIdx
		if bestIdx < 0 {
			break
		}
		bestLvl := &opp.slots[bestIdx]
		if req.Type == Limit {
			if req.Side == Buy {
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

			if req.TraderID != 0 && maker.TraderID == req.TraderID {
				stpTriggered = true
				b.eventSeq++
				dst = append(dst, Event{
					Seq:          b.eventSeq,
					Type:         EventCancelled,
					OrderID:      req.ID,
					TraderID:     req.TraderID,
					Side:         req.Side,
					OrderType:    req.Type,
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
				OrderID:       req.ID,
				TraderID:      req.TraderID,
				MakerOrderID:  maker.ID,
				MakerTraderID: maker.TraderID,
				Side:          req.Side,
				OrderType:     req.Type,
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

	b.restLimit(own, req, price, leaves, orderSeq)
	return dst
}

// PlaceOrderWithCallback streams every produced Event directly to fn with zero
// intermediate slice buffering. If fn is nil, events are discarded.
func (b *Book) PlaceOrderWithCallback(req OrderRequest, fn func(Event)) {
	emit := func(ev Event) {
		if fn != nil {
			fn(ev)
		}
	}

	if reason := b.validateOrder(req); reason != RejectNone {
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

	stpTriggered := false
	for leaves > 0 {
		bestIdx := opp.bestIdx
		if bestIdx < 0 {
			break
		}
		bestLvl := &opp.slots[bestIdx]
		if req.Type == Limit {
			if req.Side == Buy {
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

			if req.TraderID != 0 && maker.TraderID == req.TraderID {
				stpTriggered = true
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
			emit(Event{
				Seq:           b.eventSeq,
				Type:          EventTrade,
				TradeID:       b.tradeSeq,
				OrderID:       req.ID,
				TraderID:      req.TraderID,
				MakerOrderID:  maker.ID,
				MakerTraderID: maker.TraderID,
				Side:          req.Side,
				OrderType:     req.Type,
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
				emit(Event{
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

	b.restLimit(own, req, price, leaves, orderSeq)
}

func (b *Book) restLimit(own *ladder, req OrderRequest, price, leaves int64, orderSeq uint64) {
	idx := b.allocOrder()
	lvl, tickIdx := own.getOrCreate(price)
	tailIdx := lvl.tailIdx

	b.arena[idx] = arenaOrder{
		ID:       req.ID,
		TraderID: req.TraderID,
		Price:    price,
		Qty:      req.Qty,
		Leaves:   leaves,
		Seq:      orderSeq,
		PrevIdx:  tailIdx,
		NextIdx:  0,
		TickIdx:  tickIdx,
		Side:     req.Side,
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

	b.orders.Set(req.ID, idx)
}

func (b *Book) cancelInternal(id OrderID) Event {
	idx, ok := b.orders.Delete(id)
	if !ok {
		b.eventSeq++
		return Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      id,
			RejectReason: RejectOrderNotFound,
		}
	}

	o := &b.arena[idx]
	var side *ladder
	if o.Side == Buy {
		side = &b.bids
	} else {
		side = &b.asks
	}
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

	b.eventSeq++
	ev := Event{
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

	b.freeOrder(idx)
	return ev
}

// CancelOrder cancels an active resting order by ID, appending either
// EventCancelled(CancelUser) or EventRejected(RejectOrderNotFound) to dst.
func (b *Book) CancelOrder(id OrderID, dst []Event) []Event {
	return append(dst, b.cancelInternal(id))
}

// CancelOrderWithCallback cancels an active resting order by ID and streams the
// resulting event directly to fn.
func (b *Book) CancelOrderWithCallback(id OrderID, fn func(Event)) {
	ev := b.cancelInternal(id)
	if fn != nil {
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
	b.bids.reset()
	b.asks.reset()
	b.orders.Clear()
	b.freeHead = 0
	b.highMark = 1
	b.orderSeq = 0
	b.eventSeq = 0
	b.tradeSeq = 0
}
