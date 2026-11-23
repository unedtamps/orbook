package book

import "fmt"

// CheckInvariants verifies all structural and financial invariants of the book:
//   - Book is never crossed (BestBid < BestAsk).
//   - Price levels are strictly ordered best-to-worst and match Best() / Get() / Len().
//   - Every level is non-empty with accurate TotalQty and OrderCount.
//   - Orders inside each level are in strictly increasing Seq (FIFO) order, have valid
//     price/side/leaves/back-pointer, and match the b.orders index 1-to-1.
func (b *Book) CheckInvariants() error {
	bestBid, _, hasBid := b.BestBid()
	bestAsk, _, hasAsk := b.BestAsk()
	if hasBid && hasAsk && bestBid >= bestAsk {
		return fmt.Errorf("invariant violated: book is crossed (bestBid=%d >= bestAsk=%d)", bestBid, bestAsk)
	}

	seen := 0
	if err := b.checkSide(Buy, b.bids, &seen); err != nil {
		return err
	}
	if err := b.checkSide(Sell, b.asks, &seen); err != nil {
		return err
	}

	if seen != b.orders.Len() {
		return fmt.Errorf("invariant violated: orders in levels (%d) != len(b.orders) (%d)", seen, b.orders.Len())
	}
	return nil
}

func (b *Book) checkSide(side Side, levels PriceLevels, totalSeen *int) error {
	levelCount := 0
	var prevPrice int64
	var firstLvl *PriceLevel
	var checkErr error

	levels.ForEach(func(lvl *PriceLevel) bool {
		if levelCount == 0 {
			firstLvl = lvl
		} else {
			if side == Buy && lvl.Price >= prevPrice {
				checkErr = fmt.Errorf("invariant violated: %s levels not strictly descending (%d after %d)", side, lvl.Price, prevPrice)
				return false
			}
			if side == Sell && lvl.Price <= prevPrice {
				checkErr = fmt.Errorf("invariant violated: %s levels not strictly ascending (%d after %d)", side, lvl.Price, prevPrice)
				return false
			}
		}
		prevPrice = lvl.Price
		levelCount++

		if got := levels.Get(lvl.Price); got != lvl {
			checkErr = fmt.Errorf("invariant violated: %s Get(%d) returned %p, want %p", side, lvl.Price, got, lvl)
			return false
		}
		if lvl.OrderCount <= 0 || lvl.TotalQty <= 0 {
			checkErr = fmt.Errorf("invariant violated: %s level %d empty or non-positive (count=%d, qty=%d)", side, lvl.Price, lvl.OrderCount, lvl.TotalQty)
			return false
		}
		if lvl.Queue.Len() != lvl.OrderCount {
			checkErr = fmt.Errorf("invariant violated: %s level %d Queue.Len()=%d != OrderCount=%d", side, lvl.Price, lvl.Queue.Len(), lvl.OrderCount)
			return false
		}

		var qtySum int64
		orderCount := 0
		var prevSeq uint64
		var frontOrder *Order

		lvl.Queue.ForEach(func(o *Order) bool {
			if orderCount == 0 {
				frontOrder = o
			} else if o.Seq <= prevSeq {
				checkErr = fmt.Errorf("invariant violated: %s level %d FIFO order broken (seq %d after %d)", side, lvl.Price, o.Seq, prevSeq)
				return false
			}
			prevSeq = o.Seq
			orderCount++

			if o.Side != side {
				checkErr = fmt.Errorf("invariant violated: order %d has side %s inside %s book", o.ID, o.Side, side)
				return false
			}
			if o.Price != lvl.Price {
				checkErr = fmt.Errorf("invariant violated: order %d price %d != level price %d", o.ID, o.Price, lvl.Price)
				return false
			}
			if o.Level != lvl {
				checkErr = fmt.Errorf("invariant violated: order %d back-pointer %p != level %p", o.ID, o.Level, lvl)
				return false
			}
			if o.Leaves <= 0 || o.Leaves > o.Qty {
				checkErr = fmt.Errorf("invariant violated: order %d invalid leaves %d (qty=%d)", o.ID, o.Leaves, o.Qty)
				return false
			}
			indexed, ok := b.orders.Get(o.ID)
			if !ok || indexed != o {
				checkErr = fmt.Errorf("invariant violated: order %d missing or mismatched in b.orders", o.ID)
				return false
			}
			qtySum += o.Leaves
			*totalSeen++
			return true
		})

		if checkErr != nil {
			return false
		}
		if lvl.Front() != frontOrder {
			checkErr = fmt.Errorf("invariant violated: %s level %d Front()=%p != first iterated order %p", side, lvl.Price, lvl.Front(), frontOrder)
			return false
		}
		if orderCount != lvl.OrderCount {
			checkErr = fmt.Errorf("invariant violated: %s level %d iterated %d orders != OrderCount %d", side, lvl.Price, orderCount, lvl.OrderCount)
			return false
		}
		if qtySum != lvl.TotalQty {
			checkErr = fmt.Errorf("invariant violated: %s level %d sum(Leaves)=%d != TotalQty=%d", side, lvl.Price, qtySum, lvl.TotalQty)
			return false
		}
		return true
	})

	if checkErr != nil {
		return checkErr
	}
	if levelCount != levels.Len() {
		return fmt.Errorf("invariant violated: %s iterated %d levels != Len() %d", side, levelCount, levels.Len())
	}
	if levels.Best() != firstLvl {
		return fmt.Errorf("invariant violated: %s Best()=%p != first level %p", side, levels.Best(), firstLvl)
	}
	return nil
}
