// Package orderbook provides a zero-allocation, pointer-free (GC noscan),
// single-ticker stock order book matching engine using a tick-indexed Ladder,
// intrusive LinkedList FIFO queues, and a flat Fibonacci-hashed OrderMap.
//
// Book is single-threaded and lock-free by design: run one *Book per ticker
// inside its own goroutine/actor loop.
package orderbook

import "fmt"

// OrderID uniquely identifies an active resting order on the book (must be >= 1).
type OrderID uint64

// TraderID identifies the owner of an order for Cancel-Newest Self-Trade Prevention.
// A TraderID of 0 disables STP checks for that order; any non-zero TraderID cancels
// the incoming taker remainder if it would match a resting order with the same TraderID.
type TraderID uint64

// Side indicates whether an order is buying (Bid) or selling (Ask).
type Side uint8

const (
	Buy Side = iota + 1
	Sell
)

func (s Side) String() string {
	switch s {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return fmt.Sprintf("Side(%d)", s)
	}
}

// OrderType defines execution behaviour for an incoming order.
type OrderType uint8

const (
	Limit OrderType = iota + 1
	Market
)

func (t OrderType) String() string {
	switch t {
	case Limit:
		return "LIMIT"
	case Market:
		return "MARKET"
	default:
		return fmt.Sprintf("OrderType(%d)", t)
	}
}

// TimeInForce controls how long an incoming order remains active.
// Zero value is GTC (Good-Til-Cancelled) for full backward compatibility.
type TimeInForce uint8

const (
	GTC      TimeInForce = iota // Good-Til-Cancelled: unfilled limit remainder rests on the book
	IOC                         // Immediate-Or-Cancel: fill what you can now, cancel remainder
	FOK                         // Fill-Or-Kill: fill 100% immediately or cancel with 0 trades
	PostOnly                    // Maker-Only: cancel immediately if it would cross the spread
)

func (tif TimeInForce) String() string {
	switch tif {
	case GTC:
		return "GTC"
	case IOC:
		return "IOC"
	case FOK:
		return "FOK"
	case PostOnly:
		return "POST_ONLY"
	default:
		return fmt.Sprintf("TimeInForce(%d)", tif)
	}
}

// EventType identifies the kind of output event produced by the matching engine.
type EventType uint8

const (
	EventAccepted EventType = iota + 1
	EventRejected
	EventTrade
	EventCancelled
	EventDone
	EventModified
)

func (e EventType) String() string {
	switch e {
	case EventAccepted:
		return "ACCEPTED"
	case EventRejected:
		return "REJECTED"
	case EventTrade:
		return "TRADE"
	case EventCancelled:
		return "CANCELLED"
	case EventDone:
		return "DONE"
	case EventModified:
		return "MODIFIED"
	default:
		return fmt.Sprintf("EventType(%d)", e)
	}
}

// RejectReason explains why a place, modify, or cancel request was rejected.
type RejectReason uint8

const (
	RejectNone RejectReason = iota
	RejectInvalidOrderID
	RejectDuplicateOrderID
	RejectInvalidSide
	RejectInvalidType
	RejectInvalidTimeInForce
	RejectInvalidQty
	RejectNonLotQty
	RejectInvalidPrice
	RejectOutOfPriceBand
	RejectOffTickPrice
	RejectOrderNotFound
	RejectNoChange
)

func (r RejectReason) String() string {
	switch r {
	case RejectNone:
		return "NONE"
	case RejectInvalidOrderID:
		return "INVALID_ORDER_ID"
	case RejectDuplicateOrderID:
		return "DUPLICATE_ORDER_ID"
	case RejectInvalidSide:
		return "INVALID_SIDE"
	case RejectInvalidType:
		return "INVALID_TYPE"
	case RejectInvalidTimeInForce:
		return "INVALID_TIME_IN_FORCE"
	case RejectInvalidQty:
		return "INVALID_QTY"
	case RejectNonLotQty:
		return "NON_LOT_QTY"
	case RejectInvalidPrice:
		return "INVALID_PRICE"
	case RejectOutOfPriceBand:
		return "OUT_OF_PRICE_BAND"
	case RejectOffTickPrice:
		return "OFF_TICK_PRICE"
	case RejectOrderNotFound:
		return "ORDER_NOT_FOUND"
	case RejectNoChange:
		return "NO_CHANGE"
	default:
		return fmt.Sprintf("RejectReason(%d)", r)
	}
}

// CancelReason explains why an order or its unfilled remainder was cancelled.
type CancelReason uint8

const (
	CancelNone CancelReason = iota
	CancelUser
	CancelSTP
	CancelMarketUnfilled
	CancelIOC
	CancelFOK
	CancelPostOnly
)

func (r CancelReason) String() string {
	switch r {
	case CancelNone:
		return "NONE"
	case CancelUser:
		return "USER"
	case CancelSTP:
		return "SELF_TRADE_PREVENTION"
	case CancelMarketUnfilled:
		return "MARKET_UNFILLED"
	case CancelIOC:
		return "IOC_UNFILLED"
	case CancelFOK:
		return "FOK_UNFILLED"
	case CancelPostOnly:
		return "POST_ONLY_CROSS"
	default:
		return fmt.Sprintf("CancelReason(%d)", r)
	}
}

// MaxTicks is the maximum number of discrete price ticks permitted by Config.Validate
// to guard against unbounded upfront ladder allocations (~64 B per tick across both sides).
const MaxTicks = 1_000_000

// Config defines tick/lot validation rules, price band bounds, and initial arena
// capacity for a single ticker book.
type Config struct {
	Symbol          string
	TickSize        int64
	LotSize         int64
	MinPrice        int64
	MaxPrice        int64
	InitialCapacity int // Optional pre-allocated order capacity (defaults to 4096).
}

// Validate checks that Config defines a valid positive tick/lot and bounded price band.
func (c Config) Validate() error {
	if c.TickSize <= 0 {
		return fmt.Errorf("orderbook: TickSize must be > 0, got %d", c.TickSize)
	}
	if c.LotSize <= 0 {
		return fmt.Errorf("orderbook: LotSize must be > 0, got %d", c.LotSize)
	}
	if c.MinPrice <= 0 {
		return fmt.Errorf("orderbook: MinPrice must be > 0, got %d", c.MinPrice)
	}
	if c.MaxPrice < c.MinPrice {
		return fmt.Errorf("orderbook: MaxPrice (%d) must be >= MinPrice (%d)", c.MaxPrice, c.MinPrice)
	}
	span := c.MaxPrice - c.MinPrice
	if span%c.TickSize != 0 {
		return fmt.Errorf("orderbook: price range [%d, %d] is not a multiple of TickSize %d", c.MinPrice, c.MaxPrice, c.TickSize)
	}
	if span/c.TickSize >= MaxTicks {
		return fmt.Errorf("orderbook: tick count exceeds MaxTicks (%d)", MaxTicks)
	}
	return nil
}

// NumTicks returns the number of discrete valid price ticks in [MinPrice, MaxPrice].
func (c Config) NumTicks() int {
	return int((c.MaxPrice-c.MinPrice)/c.TickSize) + 1
}

// OrderRequest is the input payload for PlaceOrder / PlaceOrderWithCallback.
type OrderRequest struct {
	ID          OrderID
	TraderID    TraderID
	Side        Side
	Type        OrderType
	TimeInForce TimeInForce // Defaults to GTC (0).
	Price       int64       // Required for Limit orders; ignored for Market orders.
	Qty         int64
}

// ModifyRequest amends a resting limit order's price and/or remaining quantity (Leaves).
//   - Decreasing NewQty at the same price preserves FIFO time priority.
//   - Increasing NewQty or changing NewPrice loses time priority and re-matches if crossing.
//   - Setting NewPrice == 0 keeps the current price; setting NewQty == 0 keeps current Leaves.
type ModifyRequest struct {
	ID       OrderID
	NewPrice int64
	NewQty   int64
}

// Event is a pointer-free value record emitted by matching operations.
// Fields are ordered largest-to-smallest to eliminate internal alignment padding (88 B).
type Event struct {
	Seq           uint64
	TradeID       uint64
	OrderID       OrderID
	TraderID      TraderID
	MakerOrderID  OrderID
	MakerTraderID TraderID
	Price         int64
	Qty           int64
	Leaves        int64
	MakerLeaves   int64
	Type          EventType
	Side          Side
	OrderType     OrderType
	RejectReason  RejectReason
	CancelReason  CancelReason
}

// OrderSnapshot is a read-only value copy of a resting order returned by GetOrder.
type OrderSnapshot struct {
	ID       OrderID
	TraderID TraderID
	Side     Side
	Type     OrderType
	Price    int64
	Qty      int64
	Leaves   int64
	Seq      uint64
}

// DepthLevel represents aggregated volume at a single price level for L2 snapshots.
type DepthLevel struct {
	Price      int64
	TotalQty   int64
	OrderCount int
}

// TopOfBook captures the current best bid and best ask.
type TopOfBook struct {
	BidPrice int64
	BidQty   int64
	HasBid   bool
	AskPrice int64
	AskQty   int64
	HasAsk   bool
}

// Stats captures deterministic sequence counters and current book sizes.
type Stats struct {
	OrderSeq     uint64
	EventSeq     uint64
	TradeSeq     uint64
	ActiveOrders int
	BidLevels    int
	AskLevels    int
}

// OrderBook is the public interface satisfied by *Book for dependency injection
// and mocking. For maximum hot-path performance, callers may hold *Book directly.
type OrderBook interface {
	PlaceOrder(req OrderRequest, dst []Event) []Event
	ModifyOrder(req ModifyRequest, dst []Event) []Event
	CancelOrder(id OrderID, dst []Event) []Event
	PlaceOrderWithCallback(req OrderRequest, fn func(Event))
	ModifyOrderWithCallback(req ModifyRequest, fn func(Event))
	CancelOrderWithCallback(id OrderID, fn func(Event))
	BestBid() (price, qty int64, ok bool)
	BestAsk() (price, qty int64, ok bool)
	TopOfBook() TopOfBook
	Depth(n int, bidsDst, asksDst []DepthLevel) (bids, asks []DepthLevel)
	GetOrder(id OrderID) (OrderSnapshot, bool)
	HasOrder(id OrderID) bool
	OrderCount() int
	Stats() Stats
	Config() Config
	Reset()
}
