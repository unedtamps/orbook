package levels_test

import (
	"math/rand"
	"testing"

	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/book/levels"
	"github.com/unedtamps/orbook/book/queue"
)

func makeStores(t *testing.T, side book.Side, cfg book.Config) map[string]book.PriceLevels {
	t.Helper()
	ladder, err := levels.NewLadder(side, cfg, queue.NewLinkedList)
	if err != nil {
		t.Fatalf("NewLadder failed: %v", err)
	}
	return map[string]book.PriceLevels{
		"SortedSlice": levels.NewSortedSlice(side, queue.NewLinkedList),
		"BTree":       levels.NewBTree(side, queue.NewLinkedList),
		"RBTree":      levels.NewRBTree(side, queue.NewLinkedList),
		"Ladder":      ladder,
	}
}

func TestPriceLevelsWordBoundariesAndTickSize(t *testing.T) {
	cfg := book.Config{
		Symbol:   "BBCA",
		TickSize: 25,
		LotSize:  100,
		MinPrice: 1000,
		MaxPrice: 1000 + 25*300,
	}

	tickOffsets := []int64{0, 1, 63, 64, 65, 127, 128, 255, 300}

	for _, side := range []book.Side{book.Buy, book.Sell} {
		stores := makeStores(t, side, cfg)
		for name, st := range stores {
			t.Run(side.String()+"_"+name, func(t *testing.T) {
				if st.Best() != nil || st.Len() != 0 {
					t.Fatalf("expected empty store")
				}

				for _, off := range tickOffsets {
					p := cfg.MinPrice + off*cfg.TickSize
					lvl := st.GetOrCreate(p)
					lvl.TotalQty = off + 10
					lvl.OrderCount = 1
				}

				if st.Len() != len(tickOffsets) {
					t.Fatalf("Len() = %d, want %d", st.Len(), len(tickOffsets))
				}

				var top []book.DepthLevel
				top = st.TopN(len(tickOffsets), top)
				if len(top) != len(tickOffsets) {
					t.Fatalf("TopN len = %d, want %d", len(top), len(tickOffsets))
				}

				for i := 0; i < len(tickOffsets); i++ {
					var wantOff int64
					if side == book.Buy {
						wantOff = tickOffsets[len(tickOffsets)-1-i]
					} else {
						wantOff = tickOffsets[i]
					}
					wantPrice := cfg.MinPrice + wantOff*cfg.TickSize
					if top[i].Price != wantPrice {
						t.Fatalf("TopN[%d].Price = %d, want %d", i, top[i].Price, wantPrice)
					}
				}

				for i := 0; i < len(tickOffsets); i++ {
					var wantOff int64
					if side == book.Buy {
						wantOff = tickOffsets[len(tickOffsets)-1-i]
					} else {
						wantOff = tickOffsets[i]
					}
					wantPrice := cfg.MinPrice + wantOff*cfg.TickSize
					best := st.Best()
					if best == nil || best.Price != wantPrice {
						t.Fatalf("step %d: Best() = %v, want price %d", i, best, wantPrice)
					}
					st.Remove(wantPrice)
					if st.Get(wantPrice) != nil {
						t.Fatalf("step %d: Get(%d) non-nil after Remove", i, wantPrice)
					}
				}

				if st.Best() != nil || st.Len() != 0 {
					t.Fatalf("expected empty store after removing all levels")
				}
			})
		}
	}
}

func TestPriceLevelsDifferentialRandom(t *testing.T) {
	cfg := book.Config{
		Symbol:   "TEST",
		TickSize: 5,
		LotSize:  1,
		MinPrice: 500,
		MaxPrice: 500 + 5*200,
	}

	for _, side := range []book.Side{book.Buy, book.Sell} {
		t.Run(side.String(), func(t *testing.T) {
			stores := makeStores(t, side, cfg)
			rng := rand.New(rand.NewSource(99))
			names := []string{"SortedSlice", "BTree", "RBTree", "Ladder"}

			for step := 0; step < 3000; step++ {
				price := cfg.MinPrice + int64(rng.Intn(cfg.NumTicks()))*cfg.TickSize
				if rng.Intn(100) < 55 {
					qty := int64(rng.Intn(500) + 1)
					for _, name := range names {
						lvl := stores[name].GetOrCreate(price)
						lvl.TotalQty += qty
						lvl.OrderCount++
					}
				} else {
					for _, name := range names {
						stores[name].Remove(price)
					}
				}

				ref := stores["SortedSlice"]
				refBest := ref.Best()
				refTop := ref.TopN(10, nil)

				for _, name := range names[1:] {
					st := stores[name]
					if st.Len() != ref.Len() {
						t.Fatalf("step %d: %s Len()=%d != ref %d", step, name, st.Len(), ref.Len())
					}
					b := st.Best()
					if (b == nil) != (refBest == nil) || (b != nil && (b.Price != refBest.Price || b.TotalQty != refBest.TotalQty)) {
						t.Fatalf("step %d: %s Best()=%v != ref %v", step, name, b, refBest)
					}
					gotTop := st.TopN(10, nil)
					if len(gotTop) != len(refTop) {
						t.Fatalf("step %d: %s TopN len=%d != ref %d", step, name, len(gotTop), len(refTop))
					}
					for i := range refTop {
						if gotTop[i] != refTop[i] {
							t.Fatalf("step %d: %s TopN[%d]=%v != ref %v", step, name, i, gotTop[i], refTop[i])
						}
					}
				}
			}
		})
	}
}
