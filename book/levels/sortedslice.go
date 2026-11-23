package levels

import "github.com/unedtamps/orbook/book"

// SortedSlice stores price levels in a contiguous slice sorted worst-to-best:
//   - Buy (Bids):  ascending price, so Best() is at index len-1 (highest price).
//   - Sell (Asks): descending price, so Best() is at index len-1 (lowest price).
//
// Placing the best price at the tail makes Best(), best-level removal, and
// BBO-improving insertions O(1) without shifting elements.
type SortedSlice struct {
	side     book.Side
	qFactory book.QueueFactory
	levels   []*book.PriceLevel
	free     []*book.PriceLevel
}

// NewSortedSlice constructs a SortedSlice price-level store for the given side.
func NewSortedSlice(side book.Side, qFactory book.QueueFactory) book.PriceLevels {
	return &SortedSlice{
		side:     side,
		qFactory: qFactory,
		levels:   make([]*book.PriceLevel, 0, 64),
		free:     make([]*book.PriceLevel, 0, 64),
	}
}

func (s *SortedSlice) search(price int64) (int, bool) {
	lo, hi := 0, len(s.levels)
	if s.side == book.Buy {
		// Ascending order: find smallest i where levels[i].Price >= price.
		for lo < hi {
			mid := int(uint(lo+hi) >> 1)
			if s.levels[mid].Price < price {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
	} else {
		// Descending order: find smallest i where levels[i].Price <= price.
		for lo < hi {
			mid := int(uint(lo+hi) >> 1)
			if s.levels[mid].Price > price {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
	}
	if lo < len(s.levels) && s.levels[lo].Price == price {
		return lo, true
	}
	return lo, false
}

func (s *SortedSlice) acquireLevel(price int64) *book.PriceLevel {
	if n := len(s.free); n > 0 {
		lvl := s.free[n-1]
		s.free[n-1] = nil
		s.free = s.free[:n-1]
		lvl.Price = price
		lvl.TotalQty = 0
		lvl.OrderCount = 0
		return lvl
	}
	return &book.PriceLevel{
		Price: price,
		Queue: s.qFactory(),
	}
}

func (s *SortedSlice) releaseLevel(lvl *book.PriceLevel) {
	lvl.Queue.Reset()
	lvl.Price = 0
	lvl.TotalQty = 0
	lvl.OrderCount = 0
	s.free = append(s.free, lvl)
}

func (s *SortedSlice) Get(price int64) *book.PriceLevel {
	idx, ok := s.search(price)
	if !ok {
		return nil
	}
	return s.levels[idx]
}

func (s *SortedSlice) GetOrCreate(price int64) *book.PriceLevel {
	// Fast path: check tail (best level) first.
	n := len(s.levels)
	if n > 0 {
		last := s.levels[n-1]
		if last.Price == price {
			return last
		}
		if (s.side == book.Buy && price > last.Price) || (s.side == book.Sell && price < last.Price) {
			lvl := s.acquireLevel(price)
			s.levels = append(s.levels, lvl)
			return lvl
		}
	}

	idx, ok := s.search(price)
	if ok {
		return s.levels[idx]
	}
	lvl := s.acquireLevel(price)
	s.levels = append(s.levels, nil)
	copy(s.levels[idx+1:], s.levels[idx:])
	s.levels[idx] = lvl
	return lvl
}

func (s *SortedSlice) Remove(price int64) {
	n := len(s.levels)
	if n == 0 {
		return
	}
	// Fast path: removing the best level (tail) requires zero shifting.
	if s.levels[n-1].Price == price {
		lvl := s.levels[n-1]
		s.levels[n-1] = nil
		s.levels = s.levels[:n-1]
		s.releaseLevel(lvl)
		return
	}

	idx, ok := s.search(price)
	if !ok {
		return
	}
	lvl := s.levels[idx]
	copy(s.levels[idx:], s.levels[idx+1:])
	s.levels[n-1] = nil
	s.levels = s.levels[:n-1]
	s.releaseLevel(lvl)
}

func (s *SortedSlice) Best() *book.PriceLevel {
	n := len(s.levels)
	if n == 0 {
		return nil
	}
	return s.levels[n-1]
}

func (s *SortedSlice) Len() int {
	return len(s.levels)
}

func (s *SortedSlice) TopN(n int, dst []book.DepthLevel) []book.DepthLevel {
	if n <= 0 || len(s.levels) == 0 {
		return dst
	}
	count := 0
	for i := len(s.levels) - 1; i >= 0 && count < n; i-- {
		lvl := s.levels[i]
		dst = append(dst, book.DepthLevel{
			Price:      lvl.Price,
			TotalQty:   lvl.TotalQty,
			OrderCount: lvl.OrderCount,
		})
		count++
	}
	return dst
}

func (s *SortedSlice) ForEach(fn func(*book.PriceLevel) bool) {
	for i := len(s.levels) - 1; i >= 0; i-- {
		if !fn(s.levels[i]) {
			return
		}
	}
}

func (s *SortedSlice) Reset() {
	for i, lvl := range s.levels {
		s.releaseLevel(lvl)
		s.levels[i] = nil
	}
	s.levels = s.levels[:0]
}
