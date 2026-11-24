package pool

import "github.com/unedtamps/orbook/book"

const defaultChunkSize = 4096

// FreeList is a single-threaded, lock-free slab-backed intrusive free-list pool
// for book.Order structs. Because the matching core is single-threaded, it avoids
// sync.Pool mutex/per-P overhead and GC sweep eviction.
type FreeList struct {
	head      *book.Order
	chunkSize int
	allocated int
	freeCount int
}

// NewFreeList creates a FreeList pool with slabs of chunkSize orders (defaults to 4096).
// If prewarm > 0, it pre-allocates slabs to hold at least prewarm orders upfront.
func NewFreeList(chunkSize, prewarm int) *FreeList {
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}
	p := &FreeList{
		chunkSize: chunkSize,
	}
	for p.freeCount < prewarm {
		p.grow()
	}
	return p
}

func (p *FreeList) grow() {
	slab := make([]book.Order, p.chunkSize)
	for i := 0; i < len(slab)-1; i++ {
		slab[i].PoolNext = &slab[i+1]
	}
	slab[len(slab)-1].PoolNext = p.head
	p.head = &slab[0]
	p.allocated += len(slab)
	p.freeCount += len(slab)
}

// Get returns a freshly zeroed *book.Order from the free list, growing by one slab if empty.
func (p *FreeList) Get() *book.Order {
	if p.head == nil {
		p.grow()
	}
	o := p.head
	p.head = o.PoolNext
	p.freeCount--
	*o = book.Order{}
	return o
}

// Put clears pointer fields on o and returns it to the head of the free list.
func (p *FreeList) Put(o *book.Order) {
	if o == nil {
		return
	}
	o.Level = nil
	o.Prev = nil
	o.Next = nil
	o.PoolNext = p.head
	p.head = o
	p.freeCount++
}

// Allocated returns the total number of Order slots allocated across all slabs.
func (p *FreeList) Allocated() int {
	return p.allocated
}

// FreeCount returns the number of Order slots currently available in the free list.
func (p *FreeList) FreeCount() int {
	return p.freeCount
}
