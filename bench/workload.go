package bench

import (
	"math/rand"

	"github.com/unedtamps/orbook/book"
)

// OpKind distinguishes operations in a pre-generated Workload stream.
type OpKind uint8

const (
	OpLimit OpKind = iota + 1
	OpCancel
	OpMarket
)

// Op represents one pre-generated deterministic command for benchmarking or replay.
type Op struct {
	Kind     OpKind
	Req      book.OrderRequest
	CancelID book.OrderID
}

// WorkloadStats summarizes the composition and realized execution profile of a Workload.
type WorkloadStats struct {
	WarmupOps   int
	TotalOps    int
	LimitOps    int
	CancelOps   int
	MarketOps   int
	Trades      int
	STPCancels  int
	UserCancels int
	EndOrders   int
	EndBidLvls  int
	EndAskLvls  int
}

// Workload holds a seeded, reproducible warmup + main operation stream.
type Workload struct {
	Config book.Config
	Seed   int64
	Warmup []Op
	Ops    []Op
	Stats  WorkloadStats
}

// GenerateMixedWorkload produces a deterministic synthetic equity order flow:
//   - Warmup: resting passive limit orders across ~100 bid and ~100 ask levels around mid.
//   - Main stream: ~60% Limit (mix of passive depth and aggressive crossing around a
//     random-walk mid), ~30% Cancel (guaranteed to target currently resting orders),
//     and ~10% Market orders across 16 traders (producing realistic STP hits too).
//
// A lightweight reference book is driven during generation so every OpCancel targets
// a genuinely resting order at the exact step it executes.
func GenerateMixedWorkload(cfg book.Config, seed int64, warmupOrders, totalOps int) *Workload {
	rng := rand.New(rand.NewSource(seed))
	refVariant := Variant{Name: "ref", Levels: "ladder", Queue: "linkedlist", Pool: true}
	ref, err := refVariant.NewBook(cfg, warmupOrders+4096)
	if err != nil {
		panic(err)
	}

	numTicks := int64(cfg.NumTicks())
	midTick := numTicks / 2
	priceAt := func(tick int64) int64 {
		if tick < 0 {
			tick = 0
		}
		if tick >= numTicks {
			tick = numTicks - 1
		}
		return cfg.MinPrice + tick*cfg.TickSize
	}

	const numTraders = 16
	var nextID book.OrderID = 1
	evBuf := make([]book.Event, 0, 64)

	activeIDs := make([]book.OrderID, 0, warmupOrders+1024)
	activePos := make(map[book.OrderID]int, warmupOrders+1024)

	addActive := func(id book.OrderID) {
		activePos[id] = len(activeIDs)
		activeIDs = append(activeIDs, id)
	}
	removeActive := func(id book.OrderID) {
		idx, ok := activePos[id]
		if !ok {
			return
		}
		last := activeIDs[len(activeIDs)-1]
		activeIDs[idx] = last
		activePos[last] = idx
		activeIDs = activeIDs[:len(activeIDs)-1]
		delete(activePos, id)
	}

	var stats WorkloadStats
	stats.WarmupOps = warmupOrders
	stats.TotalOps = totalOps

	applyEvents := func(evs []book.Event, countStats bool) {
		for _, ev := range evs {
			switch ev.Type {
			case book.EventAccepted:
				if ev.OrderType == book.Limit {
					addActive(ev.OrderID)
				}
			case book.EventTrade:
				if countStats {
					stats.Trades++
				}
			case book.EventDone:
				removeActive(ev.OrderID)
			case book.EventCancelled:
				removeActive(ev.OrderID)
				if countStats {
					if ev.CancelReason == book.CancelSTP {
						stats.STPCancels++
					} else if ev.CancelReason == book.CancelUser {
						stats.UserCancels++
					}
				}
			}
		}
	}

	// Warmup: populate realistic two-sided book depth (1..100 ticks from mid).
	warmup := make([]Op, 0, warmupOrders)
	for i := 0; i < warmupOrders; i++ {
		side := book.Buy
		offset := int64(1 + rng.Intn(100))
		tick := midTick - offset
		if i%2 == 1 {
			side = book.Sell
			tick = midTick + offset
		}
		req := book.OrderRequest{
			ID:       nextID,
			TraderID: book.TraderID(1 + rng.Intn(numTraders)),
			Side:     side,
			Type:     book.Limit,
			Price:    priceAt(tick),
			Qty:      int64(1+rng.Intn(5)) * cfg.LotSize,
		}
		nextID++
		warmup = append(warmup, Op{Kind: OpLimit, Req: req})
		evBuf = ref.PlaceOrder(req, evBuf[:0])
		applyEvents(evBuf, false)
	}

	// Main mixed workload: 60% limit, 30% cancel, 10% market.
	ops := make([]Op, 0, totalOps)
	minSafeTick := int64(200)
	maxSafeTick := numTicks - 200
	if minSafeTick >= maxSafeTick {
		minSafeTick = 1
		maxSafeTick = numTicks - 2
	}

	for i := 0; i < totalOps; i++ {
		// Random-walk the latent mid-price every 8 ops.
		if i%8 == 0 {
			step := int64(rng.Intn(3) - 1) // -1, 0, +1 tick
			nextMid := midTick + step
			if nextMid >= minSafeTick && nextMid <= maxSafeTick {
				midTick = nextMid
			}
		}

		// Sync midTick softly toward current BBO mid if both sides exist.
		if i%64 == 0 {
			top := ref.TopOfBook()
			if top.HasBid && top.HasAsk {
				bboMid := ((top.BidPrice + top.AskPrice) / 2)
				t := (bboMid - cfg.MinPrice) / cfg.TickSize
				if t >= minSafeTick && t <= maxSafeTick {
					midTick = t
				}
			}
		}

		roll := rng.Intn(100)
		switch {
		case roll < 30 && len(activeIDs) > 200:
			// 30% Cancel an active resting order.
			pick := activeIDs[rng.Intn(len(activeIDs))]
			op := Op{Kind: OpCancel, CancelID: pick}
			ops = append(ops, op)
			stats.CancelOps++
			evBuf = ref.CancelOrder(pick, evBuf[:0])
			applyEvents(evBuf, true)

		case roll < 40:
			// 10% Market order.
			side := book.Buy
			if rng.Intn(2) == 1 {
				side = book.Sell
			}
			req := book.OrderRequest{
				ID:       nextID,
				TraderID: book.TraderID(1 + rng.Intn(numTraders)),
				Side:     side,
				Type:     book.Market,
				Qty:      int64(1+rng.Intn(4)) * cfg.LotSize,
			}
			nextID++
			ops = append(ops, Op{Kind: OpMarket, Req: req})
			stats.MarketOps++
			evBuf = ref.PlaceOrder(req, evBuf[:0])
			applyEvents(evBuf, true)

		default:
			// 60% Limit order (80% passive around 1..100 ticks behind mid, 20% crossing 0..5 ticks through mid).
			side := book.Buy
			if rng.Intn(2) == 1 {
				side = book.Sell
			}
			var tick int64
			if rng.Intn(100) < 20 {
				// Aggressive / crossing limit order.
				cross := int64(rng.Intn(6))
				if side == book.Buy {
					tick = midTick + cross
				} else {
					tick = midTick - cross
				}
			} else {
				// Passive resting limit order: concentrated near BBO (triangular-ish) + deeper book tail.
				var offset int64
				if rng.Intn(100) < 70 {
					offset = int64(1 + rng.Intn(25))
				} else {
					offset = int64(26 + rng.Intn(100))
				}
				if side == book.Buy {
					tick = midTick - offset
				} else {
					tick = midTick + offset
				}
			}
			req := book.OrderRequest{
				ID:       nextID,
				TraderID: book.TraderID(1 + rng.Intn(numTraders)),
				Side:     side,
				Type:     book.Limit,
				Price:    priceAt(tick),
				Qty:      int64(1+rng.Intn(5)) * cfg.LotSize,
			}
			nextID++
			ops = append(ops, Op{Kind: OpLimit, Req: req})
			stats.LimitOps++
			evBuf = ref.PlaceOrder(req, evBuf[:0])
			applyEvents(evBuf, true)
		}
	}

	stats.EndOrders = ref.OrderCount()
	stats.EndBidLvls = ref.BidLevelCount()
	stats.EndAskLvls = ref.AskLevelCount()

	return &Workload{
		Config: cfg,
		Seed:   seed,
		Warmup: warmup,
		Ops:    ops,
		Stats:  stats,
	}
}

// ApplyWarmup replays w.Warmup onto b using evBuf.
func (w *Workload) ApplyWarmup(b *book.Book, evBuf []book.Event) []book.Event {
	for i := range w.Warmup {
		op := &w.Warmup[i]
		evBuf = b.PlaceOrder(op.Req, evBuf[:0])
	}
	return evBuf
}

// ReplayOp executes one Op against b using evBuf.
func ReplayOp(b *book.Book, op *Op, evBuf []book.Event) []book.Event {
	if op.Kind == OpCancel {
		return b.CancelOrder(op.CancelID, evBuf[:0])
	}
	return b.PlaceOrder(op.Req, evBuf[:0])
}
