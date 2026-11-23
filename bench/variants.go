package bench

import (
	"fmt"

	"github.com/unedtamps/orbook/book"
	"github.com/unedtamps/orbook/book/levels"
	"github.com/unedtamps/orbook/book/queue"
	"github.com/unedtamps/orbook/pool"
)

// Variant identifies one configuration in the benchmark matrix:
// {sortedslice, btree, rbtree, ladder} x {linkedlist, ringbuffer} x {pool, nopool}.
type Variant struct {
	Name   string
	Levels string
	Queue  string
	Pool   bool
}

// DefaultConfig returns the standard single-ticker config used for benchmarks:
// $100.00 - $300.00 in 1-cent ticks (20,001 discrete ticks), 100-share lots.
func DefaultConfig() book.Config {
	return book.Config{
		Symbol:   "ACME",
		TickSize: 1,
		LotSize:  100,
		MinPrice: 10000,
		MaxPrice: 30000,
	}
}

// AllVariants returns all 16 combinations in deterministic matrix order.
func AllVariants() []Variant {
	levelKinds := []string{"sortedslice", "btree", "rbtree", "ladder"}
	queueKinds := []string{"linkedlist", "ringbuffer"}
	poolOpts := []bool{true, false}

	out := make([]Variant, 0, len(levelKinds)*len(queueKinds)*len(poolOpts))
	for _, lk := range levelKinds {
		for _, qk := range queueKinds {
			for _, usePool := range poolOpts {
				pLabel := "nopool"
				if usePool {
					pLabel = "pool"
				}
				out = append(out, Variant{
					Name:   fmt.Sprintf("%s_%s_%s", lk, qk, pLabel),
					Levels: lk,
					Queue:  qk,
					Pool:   usePool,
				})
			}
		}
	}
	return out
}

// NewBook instantiates a fresh *book.Book configured according to v.
func (v Variant) NewBook(cfg book.Config, poolPrewarm int) (*book.Book, error) {
	var qf book.QueueFactory
	switch v.Queue {
	case "linkedlist":
		qf = queue.NewLinkedList
	case "ringbuffer":
		qf = queue.NewRingBuffer
	default:
		return nil, fmt.Errorf("unknown queue type %q", v.Queue)
	}

	makeSide := func(side book.Side) (book.PriceLevels, error) {
		switch v.Levels {
		case "sortedslice":
			return levels.NewSortedSlice(side, qf), nil
		case "btree":
			return levels.NewBTree(side, qf), nil
		case "rbtree":
			return levels.NewRBTree(side, qf), nil
		case "ladder":
			return levels.NewLadder(side, cfg, qf)
		default:
			return nil, fmt.Errorf("unknown levels type %q", v.Levels)
		}
	}

	bids, err := makeSide(book.Buy)
	if err != nil {
		return nil, err
	}
	asks, err := makeSide(book.Sell)
	if err != nil {
		return nil, err
	}

	var op book.OrderPool
	if v.Pool {
		op = pool.NewFreeList(4096, poolPrewarm)
	} else {
		op = pool.NewNoPool()
	}
	return book.New(cfg, bids, asks, op)
}
