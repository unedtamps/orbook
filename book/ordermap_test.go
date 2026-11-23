package book

import (
	"math/rand"
	"testing"
	"unsafe"
)

func TestArenaStructSizes(t *testing.T) {
	if sz := unsafe.Sizeof(arenaOrder{}); sz != 64 {
		t.Fatalf("arenaOrder size = %d, want 64 (1 cache line)", sz)
	}
	if sz := unsafe.Sizeof(arenaLevel{}); sz != 32 {
		t.Fatalf("arenaLevel size = %d, want 32", sz)
	}
	if sz := unsafe.Sizeof(arenaMapSlot{}); sz != 16 {
		t.Fatalf("arenaMapSlot size = %d, want 16", sz)
	}
}

func TestOrderMapDifferentialRandom(t *testing.T) {
	m := newOrderMap(16) // small initial cap? clamped to orderMapInitCap; test small cap directly too.
	m = orderMap{
		slots: make([]orderSlot, 16),
		mask:  15,
		shift: 60,
	}
	ref := make(map[OrderID]*Order)
	rng := rand.New(rand.NewSource(777))

	for step := 0; step < 20000; step++ {
		id := OrderID(rng.Intn(400) + 1)
		roll := rng.Intn(100)
		switch {
		case roll < 45:
			o := &Order{ID: id, Price: int64(step)}
			m.Set(id, o)
			ref[id] = o
		case roll < 80:
			v1, ok1 := m.Delete(id)
			v2, ok2 := ref[id]
			delete(ref, id)
			if ok1 != ok2 || v1 != v2 {
				t.Fatalf("step %d: Delete(%d) got (%v,%v), want (%v,%v)", step, id, v1, ok1, v2, ok2)
			}
		case roll < 98:
			v1, ok1 := m.Get(id)
			v2, ok2 := ref[id]
			if ok1 != ok2 || v1 != v2 || m.Has(id) != ok2 {
				t.Fatalf("step %d: Get(%d) got (%v,%v), want (%v,%v)", step, id, v1, ok1, v2, ok2)
			}
		default:
			if step%1000 == 0 {
				m.Clear()
				clear(ref)
			}
		}

		if m.Len() != len(ref) {
			t.Fatalf("step %d: Len()=%d != ref %d", step, m.Len(), len(ref))
		}
	}

	// Verify all remaining keys match 1-to-1.
	for k, want := range ref {
		got, ok := m.Get(k)
		if !ok || got != want {
			t.Fatalf("final Get(%d) = (%v,%v), want (%v,true)", k, got, ok, want)
		}
	}
	seen := 0
	m.ForEach(func(o *Order) {
		seen++
		if ref[o.ID] != o {
			t.Fatalf("ForEach yielded mismatched order %v", o)
		}
	})
	if seen != len(ref) {
		t.Fatalf("ForEach saw %d, want %d", seen, len(ref))
	}
}
