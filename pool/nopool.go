package pool

import "github.com/nedo/orbook/book"

// NoPool allocates a fresh book.Order on the heap for every Get and relies on GC
// on Put. Used as a baseline to benchmark FreeList pooling impact in -benchmem and pprof.
type NoPool struct{}

// NewNoPool returns a heap-allocating OrderPool baseline.
func NewNoPool() *NoPool {
	return &NoPool{}
}

func (NoPool) Get() *book.Order {
	return new(book.Order)
}

func (NoPool) Put(_ *book.Order) {}
