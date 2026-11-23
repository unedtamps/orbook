package pool_test

import (
	"testing"

	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/pool"
)

func TestFreeListReuseAndZeroing(t *testing.T) {
	p := pool.NewFreeList(4, 4)
	if p.Allocated() != 4 || p.FreeCount() != 4 {
		t.Fatalf("unexpected initial counts: alloc=%d free=%d", p.Allocated(), p.FreeCount())
	}

	got := make([]*book.Order, 6)
	for i := range got {
		o := p.Get()
		if o.ID != 0 || o.Prev != nil || o.Next != nil || o.Level != nil {
			t.Fatalf("expected zeroed order, got %+v", o)
		}
		o.ID = book.OrderID(i + 1)
		o.Price = 1000
		o.Leaves = 500
		got[i] = o
	}
	if p.Allocated() != 8 || p.FreeCount() != 2 {
		t.Fatalf("after 6 gets: alloc=%d free=%d, want 8 and 2", p.Allocated(), p.FreeCount())
	}

	for _, o := range got {
		p.Put(o)
	}
	if p.FreeCount() != 8 {
		t.Fatalf("after putting all back: free=%d, want 8", p.FreeCount())
	}

	// Next 8 gets must not grow slabs and must return freshly zeroed structs.
	for i := 0; i < 8; i++ {
		o := p.Get()
		if o.ID != 0 || o.Price != 0 || o.Leaves != 0 {
			t.Fatalf("reused order not zeroed: %+v", o)
		}
	}
	if p.Allocated() != 8 {
		t.Fatalf("pool grew unexpectedly to %d", p.Allocated())
	}
}

func TestFreeListZeroAllocsSteadyState(t *testing.T) {
	p := pool.NewFreeList(64, 64)
	allocs := testing.AllocsPerRun(1000, func() {
		o := p.Get()
		o.ID = 1
		p.Put(o)
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs/op, got %v", allocs)
	}
}
