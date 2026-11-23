package levels

import (
	"github.com/google/btree"
	"github.com/unedtamps/orbook/book"
)

const defaultBTreeDegree = 32

// BTree stores price levels in a generic google/btree BTreeG ordered best-to-worst
// (Min() is always the highest-priority level). It caches the best level pointer
// for O(1) Best() queries and recycles PriceLevel structs and B-tree nodes.
type BTree struct {
	side     book.Side
	qFactory book.QueueFactory
	tree     *btree.BTreeG[*book.PriceLevel]
	best     *book.PriceLevel
	probe    book.PriceLevel
	free     []*book.PriceLevel
}

// NewBTree constructs a BTree price-level store with default degree 32.
func NewBTree(side book.Side, qFactory book.QueueFactory) book.PriceLevels {
	return NewBTreeWithDegree(side, qFactory, defaultBTreeDegree)
}

// NewBTreeWithDegree constructs a BTree price-level store with a custom degree.
func NewBTreeWithDegree(side book.Side, qFactory book.QueueFactory, degree int) book.PriceLevels {
	var less btree.LessFunc[*book.PriceLevel]
	if side == book.Buy {
		// Higher price has higher priority (comes first in Ascend / Min).
		less = func(a, b *book.PriceLevel) bool { return a.Price > b.Price }
	} else {
		// Lower price has higher priority (comes first in Ascend / Min).
		less = func(a, b *book.PriceLevel) bool { return a.Price < b.Price }
	}
	fl := btree.NewFreeListG[*book.PriceLevel](64)
	return &BTree{
		side:     side,
		qFactory: qFactory,
		tree:     btree.NewWithFreeListG(degree, less, fl),
		free:     make([]*book.PriceLevel, 0, 64),
	}
}

func (b *BTree) better(p1, p2 int64) bool {
	if b.side == book.Buy {
		return p1 > p2
	}
	return p1 < p2
}

func (b *BTree) acquireLevel(price int64) *book.PriceLevel {
	if n := len(b.free); n > 0 {
		lvl := b.free[n-1]
		b.free[n-1] = nil
		b.free = b.free[:n-1]
		lvl.Price = price
		lvl.TotalQty = 0
		lvl.OrderCount = 0
		return lvl
	}
	return &book.PriceLevel{
		Price: price,
		Queue: b.qFactory(),
	}
}

func (b *BTree) releaseLevel(lvl *book.PriceLevel) {
	lvl.Queue.Reset()
	lvl.Price = 0
	lvl.TotalQty = 0
	lvl.OrderCount = 0
	b.free = append(b.free, lvl)
}

func (b *BTree) Get(price int64) *book.PriceLevel {
	if b.best != nil && b.best.Price == price {
		return b.best
	}
	b.probe.Price = price
	lvl, ok := b.tree.Get(&b.probe)
	if !ok {
		return nil
	}
	return lvl
}

func (b *BTree) GetOrCreate(price int64) *book.PriceLevel {
	if b.best != nil && b.best.Price == price {
		return b.best
	}
	b.probe.Price = price
	if lvl, ok := b.tree.Get(&b.probe); ok {
		return lvl
	}
	lvl := b.acquireLevel(price)
	b.tree.ReplaceOrInsert(lvl)
	if b.best == nil || b.better(price, b.best.Price) {
		b.best = lvl
	}
	return lvl
}

func (b *BTree) Remove(price int64) {
	b.probe.Price = price
	removed, ok := b.tree.Delete(&b.probe)
	if !ok {
		return
	}
	if removed == b.best {
		if next, ok := b.tree.Min(); ok {
			b.best = next
		} else {
			b.best = nil
		}
	}
	b.releaseLevel(removed)
}

func (b *BTree) Best() *book.PriceLevel {
	return b.best
}

func (b *BTree) Len() int {
	return b.tree.Len()
}

func (b *BTree) TopN(n int, dst []book.DepthLevel) []book.DepthLevel {
	if n <= 0 || b.tree.Len() == 0 {
		return dst
	}
	count := 0
	b.tree.Ascend(func(lvl *book.PriceLevel) bool {
		dst = append(dst, book.DepthLevel{
			Price:      lvl.Price,
			TotalQty:   lvl.TotalQty,
			OrderCount: lvl.OrderCount,
		})
		count++
		return count < n
	})
	return dst
}

func (b *BTree) ForEach(fn func(*book.PriceLevel) bool) {
	b.tree.Ascend(fn)
}

func (b *BTree) Reset() {
	b.tree.Ascend(func(lvl *book.PriceLevel) bool {
		b.releaseLevel(lvl)
		return true
	})
	b.tree.Clear(true)
	b.best = nil
}
