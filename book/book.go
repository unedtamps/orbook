package book

import "errors"

// Book is the single-threaded, lock-free deterministic matching engine core
// for one ticker. It uses Price-Time (FIFO) priority, executes trades at the
// resting maker's price, and enforces Cancel-Newest Self-Trade Prevention.
type Book struct {
	cfg      Config
	bids     PriceLevels
	asks     PriceLevels
	orders   orderMap
	pool     OrderPool
	orderSeq uint64
	eventSeq uint64
	tradeSeq uint64
}

// New constructs a Book with the given validation config, price-level stores
// for bids and asks, and order pool.
func New(cfg Config, bids, asks PriceLevels, orderPool OrderPool) (*Book, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if bids == nil {
		return nil, errors.New("book: bids PriceLevels must not be nil")
	}
	if asks == nil {
		return nil, errors.New("book: asks PriceLevels must not be nil")
	}
	if orderPool == nil {
		return nil, errors.New("book: orderPool must not be nil")
	}
	return &Book{
		cfg:    cfg,
		bids:   bids,
		asks:   asks,
		orders: newOrderMap(orderMapInitCap),
		pool:   orderPool,
	}, nil
}

// Config returns the immutable configuration of this book.
func (b *Book) Config() Config {
	return b.cfg
}

// PlaceOrder validates, matches, and (for unfilled limit remainders) rests an order.
// Produced events are appended to events and the resulting slice is returned,
// allowing callers to reuse their buffer for zero-allocation operation.
func (b *Book) PlaceOrder(req OrderRequest, events []Event) []Event {
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

	var own, opp PriceLevels
	if req.Side == Buy {
		own, opp = b.bids, b.asks
	} else {
		own, opp = b.asks, b.bids
	}

	stpTriggered := false
	for leaves > 0 {
		bestLvl := opp.Best()
		if bestLvl == nil {
			break
		}
		if req.Type == Limit {
			if req.Side == Buy {
				if bestLvl.Price > price {
					break
				}
			} else {
				if bestLvl.Price < price {
					break
				}
			}
		}

		for leaves > 0 && !bestLvl.IsEmpty() {
			maker := bestLvl.Front()

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
			bestLvl.TotalQty -= matchQty

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
				bestLvl.PopFront()
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
				b.pool.Put(maker)
			}
		}

		if bestLvl.IsEmpty() {
			opp.Remove(bestLvl.Price)
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

	o := b.pool.Get()
	o.ID = req.ID
	o.TraderID = req.TraderID
	o.Side = req.Side
	o.Type = Limit
	o.Price = price
	o.Qty = req.Qty
	o.Leaves = leaves
	o.Seq = orderSeq

	lvl := own.GetOrCreate(price)
	lvl.Add(o)
	b.orders.Set(o.ID, o)
	return events
}

// CancelOrder cancels an active resting order by ID.
// If the order is not found, an EventRejected(RejectOrderNotFound) is appended.
func (b *Book) CancelOrder(id OrderID, events []Event) []Event {
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

	lvl := o.Level
	lvl.Remove(o)
	if lvl.IsEmpty() {
		if o.Side == Buy {
			b.bids.Remove(o.Price)
		} else {
			b.asks.Remove(o.Price)
		}
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

	b.pool.Put(o)
	return events
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

// BestBid returns the highest bid price and its aggregated quantity.
func (b *Book) BestBid() (price int64, qty int64, ok bool) {
	lvl := b.bids.Best()
	if lvl == nil {
		return 0, 0, false
	}
	return lvl.Price, lvl.TotalQty, true
}

// BestAsk returns the lowest ask price and its aggregated quantity.
func (b *Book) BestAsk() (price int64, qty int64, ok bool) {
	lvl := b.asks.Best()
	if lvl == nil {
		return 0, 0, false
	}
	return lvl.Price, lvl.TotalQty, true
}

// TopOfBook returns the current best bid and best ask in one value struct.
func (b *Book) TopOfBook() TopOfBook {
	var top TopOfBook
	if lvl := b.bids.Best(); lvl != nil {
		top.BidPrice = lvl.Price
		top.BidQty = lvl.TotalQty
		top.HasBid = true
	}
	if lvl := b.asks.Best(); lvl != nil {
		top.AskPrice = lvl.Price
		top.AskQty = lvl.TotalQty
		top.HasAsk = true
	}
	return top
}

// Depth writes up to n best bid and ask price levels into bidsDst and asksDst
// (resetting their lengths to 0 first so callers can reuse buffers with 0 allocs).
func (b *Book) Depth(n int, bidsDst, asksDst []DepthLevel) (bids []DepthLevel, asks []DepthLevel) {
	bids = b.bids.TopN(n, bidsDst[:0])
	asks = b.asks.TopN(n, asksDst[:0])
	return bids, asks
}

// OrderCount returns the number of currently resting orders.
func (b *Book) OrderCount() int {
	return b.orders.Len()
}

// BidLevelCount returns the number of active bid price levels.
func (b *Book) BidLevelCount() int {
	return b.bids.Len()
}

// AskLevelCount returns the number of active ask price levels.
func (b *Book) AskLevelCount() int {
	return b.asks.Len()
}

// Stats returns a snapshot of deterministic counters and book sizes.
func (b *Book) Stats() BookStats {
	return BookStats{
		OrderSeq:     b.orderSeq,
		EventSeq:     b.eventSeq,
		TradeSeq:     b.tradeSeq,
		ActiveOrders: b.orders.Len(),
		BidLevels:    b.bids.Len(),
		AskLevels:    b.asks.Len(),
	}
}

// HasOrder reports whether an order with id is currently resting on the book.
func (b *Book) HasOrder(id OrderID) bool {
	return b.orders.Has(id)
}

// Reset clears all resting orders and resets sequence counters while retaining
// pre-allocated map buckets, price-level stores, and pooled Order structs.
func (b *Book) Reset() {
	b.bids.Reset()
	b.asks.Reset()
	b.orders.ForEach(func(o *Order) {
		b.pool.Put(o)
	})
	b.orders.Clear()
	b.orderSeq = 0
	b.eventSeq = 0
	b.tradeSeq = 0
}
