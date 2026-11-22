package book

// LevelQueue is the FIFO queue of resting orders at a single price level.
type LevelQueue interface {
	// PushBack appends o to the tail (lowest time priority) of the queue.
	PushBack(o *Order)
	// Front returns the oldest resting order (highest time priority), or nil if empty.
	Front() *Order
	// PopFront removes and returns the oldest resting order, or nil if empty.
	PopFront() *Order
	// Remove unlinks an arbitrary resting order o currently in the queue.
	Remove(o *Order)
	// Len returns the number of live resting orders in the queue.
	Len() int
	// Reset removes all entries while optionally retaining underlying capacity.
	Reset()
	// ForEach iterates live orders in FIFO order until fn returns false.
	ForEach(fn func(*Order) bool)
}

// QueueFactory creates a fresh LevelQueue for a price level.
type QueueFactory func() LevelQueue

// PriceLevel aggregates all resting orders at one discrete price tick.
type PriceLevel struct {
	Price      int64
	TotalQty   int64
	OrderCount int
	Queue      LevelQueue
}

// Add appends o to this price level and updates level totals and back-pointer.
func (pl *PriceLevel) Add(o *Order) {
	o.Level = pl
	pl.Queue.PushBack(o)
	pl.TotalQty += o.Leaves
	pl.OrderCount++
}

// Remove deletes o from this price level and updates level totals.
func (pl *PriceLevel) Remove(o *Order) {
	pl.Queue.Remove(o)
	pl.TotalQty -= o.Leaves
	pl.OrderCount--
	o.Level = nil
}

// Front returns the highest time-priority order at this price level.
func (pl *PriceLevel) Front() *Order {
	return pl.Queue.Front()
}

// PopFront removes and returns the front order after its Leaves reached zero.
// Note: TotalQty has already been decremented during matching.
func (pl *PriceLevel) PopFront() *Order {
	o := pl.Queue.PopFront()
	if o != nil {
		pl.OrderCount--
		o.Level = nil
	}
	return o
}

// IsEmpty reports whether this price level has no remaining resting orders.
func (pl *PriceLevel) IsEmpty() bool {
	return pl.OrderCount == 0
}

// PriceLevels stores all active price levels for one side (Bids or Asks)
// ordered by price priority (highest first for Buy, lowest first for Sell).
type PriceLevels interface {
	// Get returns the active PriceLevel at price, or nil if none exists.
	Get(price int64) *PriceLevel
	// GetOrCreate returns the existing PriceLevel at price, or inserts a new one.
	GetOrCreate(price int64) *PriceLevel
	// Remove removes the PriceLevel at price (called when it becomes empty).
	Remove(price int64)
	// Best returns the highest-priority PriceLevel, or nil if the side is empty.
	Best() *PriceLevel
	// Len returns the number of active (non-empty) price levels.
	Len() int
	// TopN appends up to n depth levels in best-to-worst order to dst and returns it.
	TopN(n int, dst []DepthLevel) []DepthLevel
	// ForEach iterates active price levels in best-to-worst order until fn returns false.
	ForEach(fn func(*PriceLevel) bool)
	// Reset clears all active price levels while retaining underlying storage/freelists.
	Reset()
}

// OrderPool manages allocation and recycling of resting Order structs.
// The book is single-threaded, so implementations do not require synchronization.
type OrderPool interface {
	Get() *Order
	Put(o *Order)
}
