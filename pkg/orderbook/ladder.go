package orderbook

import (
	"math"
	"math/bits"
)

type priceLevel struct {
	price      int64
	totalQty   int64
	orderCount int32
	headIdx    arenaIdx
	tailIdx    arenaIdx
	_          uint32
}

type ladder struct {
	side     Side
	minPrice int64
	maxPrice int64
	tickSize int64
	slots    []priceLevel
	active   []uint64
	bestIdx  int
	count    int
}

func newLadder(side Side, cfg Config) ladder {
	numTicks := cfg.NumTicks()
	numWords := (numTicks + 63) / 64
	return ladder{
		side:     side,
		minPrice: cfg.MinPrice,
		maxPrice: cfg.MaxPrice,
		tickSize: cfg.TickSize,
		slots:    make([]priceLevel, numTicks),
		active:   make([]uint64, numWords),
		bestIdx:  -1,
	}
}

func (l *ladder) getOrCreateAt(tickIdx uint32, price int64) *priceLevel {
	idx := int(tickIdx)
	lvl := &l.slots[idx]
	if (l.active[idx>>6] & (uint64(1) << uint(idx&63))) == 0 {
		l.activate(idx, price)
	}
	return lvl
}

func (l *ladder) activate(idx int, price int64) {
	l.active[idx>>6] |= uint64(1) << uint(idx&63)
	lvl := &l.slots[idx]
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

func (l *ladder) removeAt(idx int) {
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

func (l *ladder) findPrevActive(maxIdx int) int {
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

func (l *ladder) findNextActive(minIdx int) int {
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

func (l *ladder) topN(n int, dst []DepthLevel) []DepthLevel {
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

func (l *ladder) reset() {
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
