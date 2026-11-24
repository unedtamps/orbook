package book

import (
	"math"
	"math/bits"
)

type fastLevel struct {
	price      int64
	totalQty   int64
	orderCount int
	head       *Order
	tail       *Order
}

type fastLadder struct {
	side     Side
	minPrice int64
	maxPrice int64
	tickSize int64
	slots    []fastLevel
	active   []uint64
	bestIdx  int
	count    int
}

func newFastLadder(side Side, cfg Config) fastLadder {
	numTicks := cfg.NumTicks()
	numWords := (numTicks + 63) / 64
	return fastLadder{
		side:     side,
		minPrice: cfg.MinPrice,
		maxPrice: cfg.MaxPrice,
		tickSize: cfg.TickSize,
		slots:    make([]fastLevel, numTicks),
		active:   make([]uint64, numWords),
		bestIdx:  -1,
	}
}

func (l *fastLadder) tickIdx(price int64) int {
	if l.tickSize == 1 {
		return int(price - l.minPrice)
	}
	return int((price - l.minPrice) / l.tickSize)
}

func (l *fastLadder) getOrCreate(price int64) (*fastLevel, int) {
	idx := l.tickIdx(price)
	w := idx >> 6
	bit := uint64(1) << uint(idx&63)
	lvl := &l.slots[idx]
	if (l.active[w] & bit) == 0 {
		l.active[w] |= bit
		lvl.price = price
		lvl.totalQty = 0
		lvl.orderCount = 0
		lvl.head = nil
		lvl.tail = nil
		l.count++
		if l.bestIdx < 0 {
			l.bestIdx = idx
		} else if l.side == Buy {
			if idx > l.bestIdx {
				l.bestIdx = idx
			}
		} else {
			if idx < l.bestIdx {
				l.bestIdx = idx
			}
		}
	}
	return lvl, idx
}

func (l *fastLadder) removeAt(idx int) {
	w := idx >> 6
	bit := uint64(1) << uint(idx&63)
	l.active[w] &^= bit
	lvl := &l.slots[idx]
	lvl.totalQty = 0
	lvl.orderCount = 0
	lvl.head = nil
	lvl.tail = nil
	l.count--

	if l.count == 0 {
		l.bestIdx = -1
		return
	}
	if idx == l.bestIdx {
		if l.side == Buy {
			l.bestIdx = l.findPrevActive(idx - 1)
		} else {
			l.bestIdx = l.findNextActive(idx + 1)
		}
	}
}

func (l *fastLadder) findPrevActive(maxIdx int) int {
	if maxIdx < 0 {
		return -1
	}
	w := maxIdx >> 6
	bit := uint(maxIdx & 63)
	mask := uint64(math.MaxUint64) >> (63 - bit)
	word := l.active[w] & mask
	for word == 0 {
		w--
		if w < 0 {
			return -1
		}
		word = l.active[w]
	}
	return (w << 6) | (63 - bits.LeadingZeros64(word))
}

func (l *fastLadder) findNextActive(minIdx int) int {
	if minIdx >= len(l.slots) {
		return -1
	}
	w := minIdx >> 6
	bit := uint(minIdx & 63)
	mask := ^uint64(0) << bit
	word := l.active[w] & mask
	for word == 0 {
		w++
		if w >= len(l.active) {
			return -1
		}
		word = l.active[w]
	}
	return (w << 6) | bits.TrailingZeros64(word)
}

func (l *fastLadder) topN(n int, dst []DepthLevel) []DepthLevel {
	if n <= 0 || l.bestIdx < 0 {
		return dst
	}
	count := 0
	if l.side == Buy {
		for w := l.bestIdx >> 6; w >= 0; w-- {
			word := l.active[w]
			for word != 0 {
				b := 63 - bits.LeadingZeros64(word)
				word &^= uint64(1) << uint(b)
				lvl := &l.slots[(w<<6)|b]
				dst = append(dst, DepthLevel{
					Price:      lvl.price,
					TotalQty:   lvl.totalQty,
					OrderCount: lvl.orderCount,
				})
				count++
				if count >= n {
					return dst
				}
			}
		}
	} else {
		for w := l.bestIdx >> 6; w < len(l.active); w++ {
			word := l.active[w]
			for word != 0 {
				b := bits.TrailingZeros64(word)
				word &^= uint64(1) << uint(b)
				lvl := &l.slots[(w<<6)|b]
				dst = append(dst, DepthLevel{
					Price:      lvl.price,
					TotalQty:   lvl.totalQty,
					OrderCount: lvl.orderCount,
				})
				count++
				if count >= n {
					return dst
				}
			}
		}
	}
	return dst
}

func (l *fastLadder) reset() {
	if l.count == 0 {
		return
	}
	for w, word := range l.active {
		for word != 0 {
			b := bits.TrailingZeros64(word)
			word &^= uint64(1) << uint(b)
			lvl := &l.slots[(w<<6)|b]
			for cur := lvl.head; cur != nil; {
				next := cur.Next
				cur.Prev = nil
				cur.Next = nil
				cur = next
			}
			lvl.head = nil
			lvl.tail = nil
			lvl.totalQty = 0
			lvl.orderCount = 0
		}
		l.active[w] = 0
	}
	l.bestIdx = -1
	l.count = 0
}

// FastBook is a concrete, zero-interface specialization of Book combining
// Ladder + inline intrusive LinkedList + inline FreeList + flat orderMap so
// every hot-path helper inlines without virtual dispatch.
type FastBook struct {
	cfg      Config
	bids     fastLadder
	asks     fastLadder
	orders   orderMap
	freeHead *Order
	chunkSz  int
	orderSeq uint64
	eventSeq uint64
	tradeSeq uint64
}

// NewFastBook constructs a devirtualized concrete FastBook with prewarmed slabs.
func NewFastBook(cfg Config, poolPrewarm int) (*FastBook, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	fb := &FastBook{
		cfg:     cfg,
		bids:    newFastLadder(Buy, cfg),
		asks:    newFastLadder(Sell, cfg),
		orders:  newOrderMap(orderMapInitCap),
		chunkSz: 4096,
	}
	for allocated := 0; allocated < poolPrewarm; allocated += fb.chunkSz {
		fb.growPool()
	}
	return fb, nil
}

func (b *FastBook) growPool() {
	slab := make([]Order, b.chunkSz)
	for i := 0; i < len(slab)-1; i++ {
		slab[i].PoolNext = &slab[i+1]
	}
	slab[len(slab)-1].PoolNext = b.freeHead
	b.freeHead = &slab[0]
}

func (b *FastBook) getOrder() *Order {
	if b.freeHead == nil {
		b.growPool()
	}
	o := b.freeHead
	b.freeHead = o.PoolNext
	o.Prev = nil
	o.Next = nil
	o.PoolNext = nil
	return o
}

func (b *FastBook) putOrder(o *Order) {
	o.Prev = nil
	o.Next = nil
	o.Level = nil
	o.PoolNext = b.freeHead
	b.freeHead = o
}

func (b *FastBook) validateOrder(req OrderRequest) RejectReason {
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

// PlaceOrder validates, matches, and rests an order with zero interface dispatch.
func (b *FastBook) PlaceOrder(req OrderRequest, events []Event) []Event {
	if reason := b.validateOrder(req); reason != RejectNone {
		b.eventSeq++
		return append(events, Event{
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
	events = append(events, Event{
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

	var own, opp *fastLadder
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

		for leaves > 0 && bestLvl.head != nil {
			maker := bestLvl.head

			if req.TraderID != 0 && maker.TraderID == req.TraderID {
				stpTriggered = true
				b.eventSeq++
				events = append(events, Event{
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
			events = append(events, Event{
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
				bestLvl.head = maker.Next
				if bestLvl.head != nil {
					bestLvl.head.Prev = nil
				} else {
					bestLvl.tail = nil
				}
				bestLvl.orderCount--
				b.orders.Delete(maker.ID)

				b.eventSeq++
				events = append(events, Event{
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
				b.putOrder(maker)
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
		return events
	}

	if leaves == 0 {
		b.eventSeq++
		return append(events, Event{
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
		return append(events, Event{
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

	o := b.getOrder()
	o.ID = req.ID
	o.TraderID = req.TraderID
	o.Side = req.Side
	o.Type = Limit
	o.Price = price
	o.Qty = req.Qty
	o.Leaves = leaves
	o.Seq = orderSeq

	lvl, idx := own.getOrCreate(price)
	o.QueuePos = uint64(idx)
	o.Prev = lvl.tail
	o.Next = nil
	if lvl.tail != nil {
		lvl.tail.Next = o
	} else {
		lvl.head = o
	}
	lvl.tail = o
	lvl.totalQty += leaves
	lvl.orderCount++

	b.orders.Set(o.ID, o)
	return events
}

// CancelOrder cancels an active resting order with single-probe map delete and inline unlink.
func (b *FastBook) CancelOrder(id OrderID, events []Event) []Event {
	o, ok := b.orders.Delete(id)
	if !ok {
		b.eventSeq++
		return append(events, Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      id,
			RejectReason: RejectOrderNotFound,
		})
	}

	var side *fastLadder
	if o.Side == Buy {
		side = &b.bids
	} else {
		side = &b.asks
	}
	idx := int(o.QueuePos)
	lvl := &side.slots[idx]

	if o.Prev != nil {
		o.Prev.Next = o.Next
	} else {
		lvl.head = o.Next
	}
	if o.Next != nil {
		o.Next.Prev = o.Prev
	} else {
		lvl.tail = o.Prev
	}
	lvl.totalQty -= o.Leaves
	lvl.orderCount--
	if lvl.orderCount == 0 {
		side.removeAt(idx)
	}

	b.eventSeq++
	events = append(events, Event{
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
	})

	b.putOrder(o)
	return events
}

func (b *FastBook) BestBid() (price, qty int64, ok bool) {
	if b.bids.bestIdx < 0 {
		return 0, 0, false
	}
	lvl := &b.bids.slots[b.bids.bestIdx]
	return lvl.price, lvl.totalQty, true
}

func (b *FastBook) BestAsk() (price, qty int64, ok bool) {
	if b.asks.bestIdx < 0 {
		return 0, 0, false
	}
	lvl := &b.asks.slots[b.asks.bestIdx]
	return lvl.price, lvl.totalQty, true
}

func (b *FastBook) TopOfBook() TopOfBook {
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

func (b *FastBook) Depth(n int, bidsDst, asksDst []DepthLevel) ([]DepthLevel, []DepthLevel) {
	return b.bids.topN(n, bidsDst[:0]), b.asks.topN(n, asksDst[:0])
}

func (b *FastBook) OrderCount() int { return b.orders.Len() }

func (b *FastBook) Stats() BookStats {
	return BookStats{
		OrderSeq:     b.orderSeq,
		EventSeq:     b.eventSeq,
		TradeSeq:     b.tradeSeq,
		ActiveOrders: b.orders.Len(),
		BidLevels:    b.bids.count,
		AskLevels:    b.asks.count,
	}
}

func (b *FastBook) Reset() {
	b.bids.reset()
	b.asks.reset()
	b.orders.ForEach(func(o *Order) {
		b.putOrder(o)
	})
	b.orders.Clear()
	b.orderSeq = 0
	b.eventSeq = 0
	b.tradeSeq = 0
}
