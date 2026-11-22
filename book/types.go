package book

import "fmt"

// OrderID uniquely identifies an order while active on the book.
type OrderID uint64

// TraderID identifies the owner of an order for self-trade prevention (STP).
// A TraderID of 0 disables STP checks for that order; any non-zero TraderID
// enforces Cancel-Newest STP against resting orders with the same TraderID.
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

// OrderType defines the execution behaviour of an incoming order.
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

// EventType identifies the kind of output event produced by the matching engine.
type EventType uint8

const (
	EventAccepted EventType = iota + 1
	EventRejected
	EventTrade
	EventCancelled
	EventDone
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
	default:
		return fmt.Sprintf("EventType(%d)", e)
	}
}

// RejectReason explains why an order or cancel request was rejected.
type RejectReason uint8

const (
	RejectNone RejectReason = iota
	RejectInvalidOrderID
	RejectDuplicateOrderID
	RejectInvalidSide
	RejectInvalidType
	RejectInvalidQty
	RejectNonLotQty
	RejectInvalidPrice
	RejectOutOfPriceBand
	RejectOffTickPrice
	RejectOrderNotFound
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
	default:
		return fmt.Sprintf("RejectReason(%d)", r)
	}
}

// CancelReason explains why an order or its remaining quantity was cancelled.
type CancelReason uint8

const (
	CancelNone CancelReason = iota
	CancelUser
	CancelSTP
	CancelMarketUnfilled
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
	default:
		return fmt.Sprintf("CancelReason(%d)", r)
	}
}

// Config defines validation rules and price boundaries for a single ticker book.
type Config struct {
	Symbol   string
	TickSize int64
	LotSize  int64
	MinPrice int64
	MaxPrice int64
}

// Validate checks that Config defines a valid positive tick/lot and price band.
func (c Config) Validate() error {
	if c.TickSize <= 0 {
		return fmt.Errorf("book: TickSize must be > 0, got %d", c.TickSize)
	}
	if c.LotSize <= 0 {
		return fmt.Errorf("book: LotSize must be > 0, got %d", c.LotSize)
	}
	if c.MinPrice <= 0 {
		return fmt.Errorf("book: MinPrice must be > 0, got %d", c.MinPrice)
	}
	if c.MaxPrice < c.MinPrice {
		return fmt.Errorf("book: MaxPrice (%d) must be >= MinPrice (%d)", c.MaxPrice, c.MinPrice)
	}
	if (c.MaxPrice-c.MinPrice)%c.TickSize != 0 {
		return fmt.Errorf("book: price range [%d, %d] is not a multiple of TickSize %d", c.MinPrice, c.MaxPrice, c.TickSize)
	}
	return nil
}

// NumTicks returns the total number of discrete valid price ticks in [MinPrice, MaxPrice].
func (c Config) NumTicks() int {
	return int((c.MaxPrice-c.MinPrice)/c.TickSize) + 1
}

// OrderRequest is the input payload to PlaceOrder.
type OrderRequest struct {
	ID       OrderID
	TraderID TraderID
	Side     Side
	Type     OrderType
	Price    int64 // Required for Limit orders; ignored for Market orders.
	Qty      int64
}

// Order represents a resting limit order inside a PriceLevel's FIFO queue.
// It carries intrusive fields for both LinkedList and RingBuffer queues as well
// as a back-pointer to its owning PriceLevel for O(1) cancellation.
type Order struct {
	ID       OrderID
	TraderID TraderID
	Side     Side
	Type     OrderType
	Price    int64
	Qty      int64
	Leaves   int64
	Seq      uint64

	// Back-pointer to owning PriceLevel for O(1) cancel lookup.
	Level *PriceLevel

	// Intrusive doubly-linked list links (used by queue.LinkedList and pool.FreeList).
	Prev *Order
	Next *Order

	// Logical slot position inside queue.RingBuffer for O(1) tombstone cancel.
	QueuePos uint64
}

// Event is a value-typed output record appended to a caller-supplied buffer.
// Keeping Event pointer-free ensures zero allocations and zero GC write barriers.
type Event struct {
	Seq           uint64
	Type          EventType
	TradeID       uint64
	OrderID       OrderID
	TraderID      TraderID
	MakerOrderID  OrderID
	MakerTraderID TraderID
	Side          Side
	OrderType     OrderType
	Price         int64
	Qty           int64
	Leaves        int64
	MakerLeaves   int64
	RejectReason  RejectReason
	CancelReason  CancelReason
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

// BookStats captures deterministic counters and current size of the book.
type BookStats struct {
	OrderSeq     uint64
	EventSeq     uint64
	TradeSeq     uint64
	ActiveOrders int
	BidLevels    int
	AskLevels    int
}
