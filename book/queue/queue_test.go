package queue_test

import (
	"math/rand"
	"testing"

	"github.com/nedo/orbook/book"
	"github.com/nedo/orbook/book/queue"
)

func TestQueuesBasicOperations(t *testing.T) {
	factories := map[string]book.QueueFactory{
		"LinkedList": queue.NewLinkedList,
		"RingBuffer": queue.NewRingBuffer,
	}

	for name, factory := range factories {
		t.Run(name, func(t *testing.T) {
			q := factory()
			if q.Len() != 0 || q.Front() != nil || q.PopFront() != nil {
				t.Fatalf("expected empty queue")
			}

			orders := make([]*book.Order, 5)
			for i := range orders {
				orders[i] = &book.Order{ID: book.OrderID(i + 1)}
				q.PushBack(orders[i])
			}
			if q.Len() != 5 {
				t.Fatalf("got len %d, want 5", q.Len())
			}
			if q.Front() != orders[0] {
				t.Fatalf("got front %v, want %v", q.Front(), orders[0])
			}

			// Remove middle, head, tail.
			q.Remove(orders[2]) // remaining: 1,2,4,5
			q.Remove(orders[0]) // remaining: 2,4,5
			q.Remove(orders[4]) // remaining: 2,4

			if q.Len() != 2 {
				t.Fatalf("got len %d, want 2", q.Len())
			}
			var got []book.OrderID
			q.ForEach(func(o *book.Order) bool {
				got = append(got, o.ID)
				return true
			})
			if len(got) != 2 || got[0] != 2 || got[1] != 4 {
				t.Fatalf("ForEach got %v, want [2 4]", got)
			}

			if p := q.PopFront(); p != orders[1] {
				t.Fatalf("PopFront got %v, want order 2", p)
			}
			if p := q.PopFront(); p != orders[3] {
				t.Fatalf("PopFront got %v, want order 4", p)
			}
			if q.Len() != 0 || q.Front() != nil || q.PopFront() != nil {
				t.Fatalf("expected empty queue after draining")
			}

			// Reuse after Reset.
			q.PushBack(orders[0])
			q.PushBack(orders[1])
			q.Reset()
			if q.Len() != 0 || q.Front() != nil {
				t.Fatalf("expected empty queue after Reset")
			}
			q.PushBack(orders[4])
			if q.Len() != 1 || q.Front() != orders[4] {
				t.Fatalf("expected order 5 after Reset+PushBack")
			}
		})
	}
}

func TestRingBufferWrapAroundAndCompaction(t *testing.T) {
	q := queue.NewRingBuffer()
	var live []*book.Order

	// Repeatedly push and remove middle elements to trigger in-place compaction
	// across wrapped head/tail boundaries without unbounded growth.
	nextID := book.OrderID(1)
	for round := 0; round < 50; round++ {
		for i := 0; i < 6; i++ {
			o := &book.Order{ID: nextID}
			nextID++
			q.PushBack(o)
			live = append(live, o)
		}
		// Remove middle orders to create tombstones.
		for len(live) > 3 {
			idx := len(live) / 2
			o := live[idx]
			q.Remove(o)
			live = append(live[:idx], live[idx+1:]...)
		}
		// Pop front once to advance head past wrap boundary.
		front := q.PopFront()
		if front != live[0] {
			t.Fatalf("round %d: PopFront got %v, want %v", round, front, live[0])
		}
		live = live[1:]

		if q.Len() != len(live) {
			t.Fatalf("round %d: len mismatch got %d want %d", round, q.Len(), len(live))
		}
		i := 0
		q.ForEach(func(o *book.Order) bool {
			if o != live[i] {
				t.Fatalf("round %d: ForEach[%d] got %v want %v", round, i, o, live[i])
			}
			i++
			return true
		})
	}
}

func TestQueuesDifferentialRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	ll := queue.NewLinkedList()
	rb := queue.NewRingBuffer()

	var llOrders []*book.Order
	var rbOrders []*book.Order
	var nextID book.OrderID = 1

	for step := 0; step < 5000; step++ {
		op := rng.Intn(100)
		switch {
		case op < 50 || len(llOrders) == 0:
			// PushBack
			o1 := &book.Order{ID: nextID}
			o2 := &book.Order{ID: nextID}
			nextID++
			ll.PushBack(o1)
			rb.PushBack(o2)
			llOrders = append(llOrders, o1)
			rbOrders = append(rbOrders, o2)
		case op < 75:
			// PopFront
			p1 := ll.PopFront()
			p2 := rb.PopFront()
			if p1.ID != p2.ID {
				t.Fatalf("step %d: PopFront mismatch %d vs %d", step, p1.ID, p2.ID)
			}
			llOrders = llOrders[1:]
			rbOrders = rbOrders[1:]
		default:
			// Remove random element
			idx := rng.Intn(len(llOrders))
			ll.Remove(llOrders[idx])
			rb.Remove(rbOrders[idx])
			llOrders = append(llOrders[:idx], llOrders[idx+1:]...)
			rbOrders = append(rbOrders[:idx], rbOrders[idx+1:]...)
		}

		if ll.Len() != rb.Len() || ll.Len() != len(llOrders) {
			t.Fatalf(
				"step %d: len mismatch ll=%d rb=%d ref=%d",
				step,
				ll.Len(),
				rb.Len(),
				len(llOrders),
			)
		}
		f1, f2 := ll.Front(), rb.Front()
		if (f1 == nil) != (f2 == nil) || (f1 != nil && f1.ID != f2.ID) {
			t.Fatalf("step %d: Front mismatch %v vs %v", step, f1, f2)
		}

		if step%100 == 0 {
			var ids1, ids2 []book.OrderID
			ll.ForEach(func(o *book.Order) bool { ids1 = append(ids1, o.ID); return true })
			rb.ForEach(func(o *book.Order) bool { ids2 = append(ids2, o.ID); return true })
			if len(ids1) != len(ids2) {
				t.Fatalf("step %d: ForEach len mismatch %d vs %d", step, len(ids1), len(ids2))
			}
			for i := range ids1 {
				if ids1[i] != ids2[i] {
					t.Fatalf("step %d: ForEach[%d] mismatch %d vs %d", step, i, ids1[i], ids2[i])
				}
			}
		}
	}
}
