package levels

import (
	"math"
	"math/bits"

	"github.com/unedtamps/orbook/book"
)

// Ladder is a fixed price-ladder array indexed directly by tick offset from MinPrice.
// It maintains a []uint64 active bitset so Best() transitions and TopN/ForEach skip
// empty price gaps 64 ticks at a time using hardware LZCNT/TZCNT instructions.
type Ladder struct {
	side     book.Side
	minPrice int64
	maxPrice int64
	tickSize int64
	qFactory book.QueueFactory

	slots   []book.PriceLevel
	active  []uint64
	bestIdx int
	count   int
}

// NewLadder constructs a fixed price-ladder array bounded by cfg.MinPrice..cfg.MaxPrice.
func NewLadder(side book.Side, cfg book.Config, qFactory book.QueueFactory) (book.PriceLevels, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	numTicks := cfg.NumTicks()
	numWords := (numTicks + 63) / 64
	return &Ladder{
		side:     side,
		minPrice: cfg.MinPrice,
		maxPrice: cfg.MaxPrice,
		tickSize: cfg.TickSize,
		qFactory: qFactory,
		slots:    make([]book.PriceLevel, numTicks),
		active:   make([]uint64, numWords),
		bestIdx:  -1,
	}, nil
}

func (l *Ladder) index(price int64) (int, bool) {
	if price < l.minPrice || price > l.maxPrice {
		return 0, false
	}
	diff := price - l.minPrice
	if l.tickSize == 1 {
		return int(diff), true
	}
	if diff%l.tickSize != 0 {
		return 0, false
	}
	return int(diff / l.tickSize), true
}

func (l *Ladder) Get(price int64) *book.PriceLevel {
	idx, ok := l.index(price)
	if !ok {
		return nil
	}
	if (l.active[idx>>6] & (uint64(1) << uint(idx&63))) == 0 {
		return nil
	}
	return &l.slots[idx]
}

func (l *Ladder) GetOrCreate(price int64) *book.PriceLevel {
	idx, ok := l.index(price)
	if !ok {
		panic("levels.Ladder: price outside configured band or off-tick")
	}
	w := idx >> 6
	bit := uint64(1) << uint(idx&63)
	lvl := &l.slots[idx]
	if (l.active[w] & bit) == 0 {
		l.active[w] |= bit
		if lvl.Queue == nil {
			lvl.Queue = l.qFactory()
		}
		lvl.Price = price
		lvl.TotalQty = 0
		lvl.OrderCount = 0
		l.count++
		if l.bestIdx < 0 {
			l.bestIdx = idx
		} else if l.side == book.Buy {
			if idx > l.bestIdx {
				l.bestIdx = idx
			}
		} else {
			if idx < l.bestIdx {
				l.bestIdx = idx
			}
		}
	}
	return lvl
}

func (l *Ladder) Remove(price int64) {
	idx, ok := l.index(price)
	if !ok {
		return
	}
	w := idx >> 6
	bit := uint64(1) << uint(idx&63)
	if (l.active[w] & bit) == 0 {
		return
	}
	l.active[w] &^= bit
	lvl := &l.slots[idx]
	lvl.Queue.Reset()
	lvl.TotalQty = 0
	lvl.OrderCount = 0
	l.count--

	if l.count == 0 {
		l.bestIdx = -1
		return
	}
	if idx == l.bestIdx {
		if l.side == book.Buy {
			l.bestIdx = l.findPrevActive(idx - 1)
		} else {
			l.bestIdx = l.findNextActive(idx + 1)
		}
	}
}

// findPrevActive returns the highest active tick index <= maxIdx, or -1 if none.
func (l *Ladder) findPrevActive(maxIdx int) int {
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

// findNextActive returns the lowest active tick index >= minIdx, or -1 if none.
func (l *Ladder) findNextActive(minIdx int) int {
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

func (l *Ladder) Best() *book.PriceLevel {
	if l.bestIdx < 0 {
		return nil
	}
	return &l.slots[l.bestIdx]
}

func (l *Ladder) Len() int {
	return l.count
}

func (l *Ladder) TopN(n int, dst []book.DepthLevel) []book.DepthLevel {
	if n <= 0 || l.bestIdx < 0 {
		return dst
	}
	count := 0
	if l.side == book.Buy {
		for w := l.bestIdx >> 6; w >= 0; w-- {
			word := l.active[w]
			for word != 0 {
				b := 63 - bits.LeadingZeros64(word)
				word &^= uint64(1) << uint(b)
				lvl := &l.slots[(w<<6)|b]
				dst = append(dst, book.DepthLevel{
					Price:      lvl.Price,
					TotalQty:   lvl.TotalQty,
					OrderCount: lvl.OrderCount,
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
				dst = append(dst, book.DepthLevel{
					Price:      lvl.Price,
					TotalQty:   lvl.TotalQty,
					OrderCount: lvl.OrderCount,
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

func (l *Ladder) ForEach(fn func(*book.PriceLevel) bool) {
	if l.bestIdx < 0 {
		return
	}
	if l.side == book.Buy {
		for w := l.bestIdx >> 6; w >= 0; w-- {
			word := l.active[w]
			for word != 0 {
				b := 63 - bits.LeadingZeros64(word)
				word &^= uint64(1) << uint(b)
				if !fn(&l.slots[(w<<6)|b]) {
					return
				}
			}
		}
	} else {
		for w := l.bestIdx >> 6; w < len(l.active); w++ {
			word := l.active[w]
			for word != 0 {
				b := bits.TrailingZeros64(word)
				word &^= uint64(1) << uint(b)
				if !fn(&l.slots[(w<<6)|b]) {
					return
				}
			}
		}
	}
}

func (l *Ladder) Reset() {
	if l.count == 0 {
		return
	}
	for w, word := range l.active {
		for word != 0 {
			b := bits.TrailingZeros64(word)
			word &^= uint64(1) << uint(b)
			lvl := &l.slots[(w<<6)|b]
			lvl.Queue.Reset()
			lvl.TotalQty = 0
			lvl.OrderCount = 0
		}
		l.active[w] = 0
	}
	l.bestIdx = -1
	l.count = 0
}
