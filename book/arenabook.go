package book

import (
	"math"
	"math/bits"
)

type arenaIdx uint32

const nilArenaIdx arenaIdx = 0

type arenaOrder struct {
	ID       OrderID
	TraderID TraderID
	Price    int64
	Qty      int64
	Leaves   int64
	Seq      uint64
	PrevIdx  arenaIdx
	NextIdx  arenaIdx
	TickIdx  uint32
	Side     Side
	Type     OrderType
	_        [2]byte
}

type arenaLevel struct {
	price      int64
	totalQty   int64
	orderCount int32
	headIdx    arenaIdx
	tailIdx    arenaIdx
	_          uint32
}

type arenaMapSlot struct {
	key OrderID
	idx arenaIdx
	_   uint32
}

type arenaOrderMap struct {
	slots []arenaMapSlot
	mask  uint64
	shift uint8
	count int
}

func newArenaOrderMap(initCap int) arenaOrderMap {
	capPow2 := orderMapInitCap
	for capPow2 < initCap {
		capPow2 <<= 1
	}
	return arenaOrderMap{
		slots: make([]arenaMapSlot, capPow2),
		mask:  uint64(capPow2 - 1),
		shift: uint8(64 - bits.TrailingZeros64(uint64(capPow2))),
	}
}

func (m *arenaOrderMap) hash(id OrderID) uint64 {
	x := uint64(id)
	return ((x ^ (x >> 32)) * fibMul64) >> m.shift
}

func (m *arenaOrderMap) Has(id OrderID) bool {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			return true
		}
		if k == 0 {
			return false
		}
		i = (i + 1) & mask
	}
}

func (m *arenaOrderMap) Set(id OrderID, idx arenaIdx) {
	if m.count*2 >= len(m.slots) {
		m.grow()
	}
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == 0 {
			m.slots[i] = arenaMapSlot{key: id, idx: idx}
			m.count++
			return
		}
		if k == id {
			m.slots[i].idx = idx
			return
		}
		i = (i + 1) & mask
	}
}

func (m *arenaOrderMap) Delete(id OrderID) (arenaIdx, bool) {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			idx := m.slots[i].idx
			m.count--
			hole := i
			j := (i + 1) & mask
			for {
				jk := m.slots[j].key
				if jk == 0 {
					m.slots[hole] = arenaMapSlot{}
					return idx, true
				}
				home := m.hash(jk)
				if ((hole - home) & mask) < ((j - home) & mask) {
					m.slots[hole] = m.slots[j]
					hole = j
				}
				j = (j + 1) & mask
			}
		}
		if k == 0 {
			return nilArenaIdx, false
		}
		i = (i + 1) & mask
	}
}

func (m *arenaOrderMap) grow() {
	oldSlots := m.slots
	newCap := len(oldSlots) * 2
	m.slots = make([]arenaMapSlot, newCap)
	m.mask = uint64(newCap - 1)
	m.shift = uint8(64 - bits.TrailingZeros64(uint64(newCap)))

	mask := m.mask
	for _, s := range oldSlots {
		if s.key == 0 {
			continue
		}
		i := m.hash(s.key)
		for m.slots[i].key != 0 {
			i = (i + 1) & mask
		}
		m.slots[i] = s
	}
}

func (m *arenaOrderMap) Clear() {
	if m.count == 0 {
		return
	}
	clear(m.slots)
	m.count = 0
}

type arenaLadder struct {
	side     Side
	minPrice int64
	maxPrice int64
	tickSize int64
	slots    []arenaLevel
	active   []uint64
	bestIdx  int
	count    int
}

func newArenaLadder(side Side, cfg Config) arenaLadder {
	numTicks := cfg.NumTicks()
	numWords := (numTicks + 63) / 64
	return arenaLadder{
		side:     side,
		minPrice: cfg.MinPrice,
		maxPrice: cfg.MaxPrice,
		tickSize: cfg.TickSize,
		slots:    make([]arenaLevel, numTicks),
		active:   make([]uint64, numWords),
		bestIdx:  -1,
	}
}

func (l *arenaLadder) tickIdx(price int64) int {
	if l.tickSize == 1 {
		return int(price - l.minPrice)
	}
	return int((price - l.minPrice) / l.tickSize)
}

func (l *arenaLadder) getOrCreate(price int64) (*arenaLevel, uint32) {
	idx := l.tickIdx(price)
	w := idx >> 6
	bit := uint64(1) << uint(idx&63)
	lvl := &l.slots[idx]
	if (l.active[w] & bit) == 0 {
		l.active[w] |= bit
		lvl.price = price
		lvl.totalQty = 0
		lvl.orderCount = 0
		lvl.headIdx = nilArenaIdx
		lvl.tailIdx = nilArenaIdx
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
	return lvl, uint32(idx)
}

func (l *arenaLadder) removeAt(idx int) {
	w := idx >> 6
	bit := uint64(1) << uint(idx&63)
	l.active[w] &^= bit
	lvl := &l.slots[idx]
	lvl.totalQty = 0
	lvl.orderCount = 0
	lvl.headIdx = nilArenaIdx
	lvl.tailIdx = nilArenaIdx
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

func (l *arenaLadder) findPrevActive(maxIdx int) int {
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

func (l *arenaLadder) findNextActive(minIdx int) int {
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

func (l *arenaLadder) topN(n int, dst []DepthLevel) []DepthLevel {
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
					OrderCount: int(lvl.orderCount),
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
					OrderCount: int(lvl.orderCount),
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

func (l *arenaLadder) reset() {
	if l.count == 0 {
		return
	}
	for w, word := range l.active {
		for word != 0 {
			b := bits.TrailingZeros64(word)
			word &^= uint64(1) << uint(b)
			lvl := &l.slots[(w<<6)|b]
			lvl.headIdx = nilArenaIdx
			lvl.tailIdx = nilArenaIdx
			lvl.totalQty = 0
			lvl.orderCount = 0
		}
		l.active[w] = 0
	}
	l.bestIdx = -1
	l.count = 0
}

// ArenaBook is a zero-interface, pointer-free (GC noscan) matching engine core.
// Every internal table ([]arenaOrder, []arenaLevel, []arenaMapSlot, []uint64)
// contains zero Go pointers, shrinking orders to 64 bytes (1 cache line) and
// reducing Go GC mark scan work over resting orders to near zero.
type ArenaBook struct {
	cfg      Config
	bids     arenaLadder
	asks     arenaLadder
	orders   arenaOrderMap
	arena    []arenaOrder
	freeHead arenaIdx
	highMark arenaIdx
	orderSeq uint64
	eventSeq uint64
	tradeSeq uint64
}

// NewArenaBook constructs a pointer-free noscan ArenaBook pre-sized for initCap orders.
func NewArenaBook(cfg Config, initCap int) (*ArenaBook, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if initCap < 4096 {
		initCap = 4096
	}
	return &ArenaBook{
		cfg:      cfg,
		bids:     newArenaLadder(Buy, cfg),
		asks:     newArenaLadder(Sell, cfg),
		orders:   newArenaOrderMap(initCap * 2),
		arena:    make([]arenaOrder, initCap+1),
		highMark: 1,
	}, nil
}

func (b *ArenaBook) allocOrder() arenaIdx {
	if b.freeHead != nilArenaIdx {
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

func (b *ArenaBook) freeOrder(idx arenaIdx) {
	b.arena[idx].NextIdx = b.freeHead
	b.freeHead = idx
}

func (b *ArenaBook) validateOrder(req OrderRequest) RejectReason {
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

// PlaceOrder validates, matches, and rests an order using 64-byte pointer-free arena slots.
func (b *ArenaBook) PlaceOrder(req OrderRequest, events []Event) []Event {
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

	var own, opp *arenaLadder
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

		for leaves > 0 && bestLvl.headIdx != nilArenaIdx {
			makerIdx := bestLvl.headIdx
			maker := &b.arena[makerIdx]

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
				nextIdx := maker.NextIdx
				bestLvl.headIdx = nextIdx
				if nextIdx != nilArenaIdx {
					b.arena[nextIdx].PrevIdx = nilArenaIdx
				} else {
					bestLvl.tailIdx = nilArenaIdx
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
		NextIdx:  nilArenaIdx,
		TickIdx:  tickIdx,
		Side:     req.Side,
		Type:     Limit,
	}

	if tailIdx != nilArenaIdx {
		b.arena[tailIdx].NextIdx = idx
	} else {
		lvl.headIdx = idx
	}
	lvl.tailIdx = idx
	lvl.totalQty += leaves
	lvl.orderCount++

	b.orders.Set(req.ID, idx)
	return events
}

// CancelOrder cancels a resting order with zero pointer chasing and single-probe map delete.
func (b *ArenaBook) CancelOrder(id OrderID, events []Event) []Event {
	idx, ok := b.orders.Delete(id)
	if !ok {
		b.eventSeq++
		return append(events, Event{
			Seq:          b.eventSeq,
			Type:         EventRejected,
			OrderID:      id,
			RejectReason: RejectOrderNotFound,
		})
	}

	o := &b.arena[idx]
	var side *arenaLadder
	if o.Side == Buy {
		side = &b.bids
	} else {
		side = &b.asks
	}
	lvl := &side.slots[o.TickIdx]

	prevIdx := o.PrevIdx
	nextIdx := o.NextIdx
	if prevIdx != nilArenaIdx {
		b.arena[prevIdx].NextIdx = nextIdx
	} else {
		lvl.headIdx = nextIdx
	}
	if nextIdx != nilArenaIdx {
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

	b.freeOrder(idx)
	return events
}

func (b *ArenaBook) BestBid() (price, qty int64, ok bool) {
	if b.bids.bestIdx < 0 {
		return 0, 0, false
	}
	lvl := &b.bids.slots[b.bids.bestIdx]
	return lvl.price, lvl.totalQty, true
}

func (b *ArenaBook) BestAsk() (price, qty int64, ok bool) {
	if b.asks.bestIdx < 0 {
		return 0, 0, false
	}
	lvl := &b.asks.slots[b.asks.bestIdx]
	return lvl.price, lvl.totalQty, true
}

func (b *ArenaBook) TopOfBook() TopOfBook {
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

func (b *ArenaBook) Depth(n int, bidsDst, asksDst []DepthLevel) ([]DepthLevel, []DepthLevel) {
	return b.bids.topN(n, bidsDst[:0]), b.asks.topN(n, asksDst[:0])
}

func (b *ArenaBook) OrderCount() int { return b.orders.count }

func (b *ArenaBook) Stats() BookStats {
	return BookStats{
		OrderSeq:     b.orderSeq,
		EventSeq:     b.eventSeq,
		TradeSeq:     b.tradeSeq,
		ActiveOrders: b.orders.count,
		BidLevels:    b.bids.count,
		AskLevels:    b.asks.count,
	}
}

// Reset clears the entire arena in O(active levels + map size) without walking orders.
func (b *ArenaBook) Reset() {
	b.bids.reset()
	b.asks.reset()
	b.orders.Clear()
	b.freeHead = nilArenaIdx
	b.highMark = 1
	b.orderSeq = 0
	b.eventSeq = 0
	b.tradeSeq = 0
}
