package engine

import (
	"sync"

	"github.com/nedo/orbook/book"
)

// CommandType identifies an operation sent to the single-threaded Engine goroutine.
type CommandType uint8

const (
	CmdPlaceOrder CommandType = iota + 1
	CmdCancelOrder
	CmdTopOfBook
	CmdDepth
	CmdStats
)

// Response carries the result of a synchronous command back to the caller.
type Response struct {
	Events []book.Event
	Top    book.TopOfBook
	Bids   []book.DepthLevel
	Asks   []book.DepthLevel
	Stats  book.BookStats
}

// Command is a single message pushed onto the engine's LMAX-style input channel.
type Command struct {
	Type     CommandType
	Order    book.OrderRequest
	OrderID  book.OrderID
	DepthN   int
	EventDst []book.Event
	BidsDst  []book.DepthLevel
	AsksDst  []book.DepthLevel
	Reply    chan Response // Optional: nil for fire-and-forget async commands.
}

// EventHandler is invoked on the engine goroutine after each mutating command.
// The provided slice is backed by the engine's reusable buffer and must not be
// retained after the callback returns.
type EventHandler func(events []book.Event)

// Options configures the LMAX-style Engine runner.
type Options struct {
	QueueSize    int          // Input channel buffer size (default 4096).
	EventBufCap  int          // Initial capacity of engine's reusable event buffer (default 64).
	EventHandler EventHandler // Optional callback for all produced events.
}

// Engine wraps a single-threaded *book.Book in an LMAX-style actor goroutine fed
// by one input channel. The book itself uses zero locks; all mutations and reads
// are serialized through the input channel.
type Engine struct {
	book     *book.Book
	in       chan Command
	done     chan struct{}
	handler  EventHandler
	eventBuf []book.Event

	replyPool sync.Pool
	startOnce sync.Once
	stopOnce  sync.Once
	started   bool
	mu        sync.Mutex // guards started flag only (never touched on hot path)
}

// New creates an Engine owning b. Call Start() before submitting commands, and Stop()
// after all producers have stopped submitting to drain and shut down cleanly.
func New(b *book.Book, opts Options) *Engine {
	qSize := opts.QueueSize
	if qSize <= 0 {
		qSize = 4096
	}
	evCap := opts.EventBufCap
	if evCap <= 0 {
		evCap = 64
	}
	e := &Engine{
		book:     b,
		in:       make(chan Command, qSize),
		done:     make(chan struct{}),
		handler:  opts.EventHandler,
		eventBuf: make([]book.Event, 0, evCap),
	}
	e.replyPool.New = func() any {
		return make(chan Response, 1)
	}
	return e
}

// Start launches the single matching goroutine. Calling Start more than once is a no-op.
func (e *Engine) Start() {
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.mu.Unlock()
		go e.run()
	})
}

// Stop closes the input channel and blocks until all buffered commands have been
// processed and the engine goroutine has exited. Producers must stop calling Submit
// before Stop is called.
func (e *Engine) Stop() {
	e.stopOnce.Do(func() {
		close(e.in)
		e.mu.Lock()
		started := e.started
		e.mu.Unlock()
		if started {
			<-e.done
		}
	})
}

func (e *Engine) run() {
	defer close(e.done)
	for cmd := range e.in {
		e.process(cmd)
	}
}

func (e *Engine) process(cmd Command) {
	switch cmd.Type {
	case CmdPlaceOrder:
		e.eventBuf = e.book.PlaceOrder(cmd.Order, e.eventBuf[:0])
		if e.handler != nil {
			e.handler(e.eventBuf)
		}
		if cmd.Reply != nil {
			dst := append(cmd.EventDst[:0], e.eventBuf...)
			cmd.Reply <- Response{Events: dst}
		}

	case CmdCancelOrder:
		e.eventBuf = e.book.CancelOrder(cmd.OrderID, e.eventBuf[:0])
		if e.handler != nil {
			e.handler(e.eventBuf)
		}
		if cmd.Reply != nil {
			dst := append(cmd.EventDst[:0], e.eventBuf...)
			cmd.Reply <- Response{Events: dst}
		}

	case CmdTopOfBook:
		if cmd.Reply != nil {
			cmd.Reply <- Response{Top: e.book.TopOfBook()}
		}

	case CmdDepth:
		if cmd.Reply != nil {
			bids, asks := e.book.Depth(cmd.DepthN, cmd.BidsDst, cmd.AsksDst)
			cmd.Reply <- Response{Bids: bids, Asks: asks}
		}

	case CmdStats:
		if cmd.Reply != nil {
			cmd.Reply <- Response{Stats: e.book.Stats()}
		}
	}
}

// Submit pushes a raw Command onto the input channel.
func (e *Engine) Submit(cmd Command) {
	e.in <- cmd
}

// PlaceOrderAsync submits a non-blocking (until buffer full) PlaceOrder command.
func (e *Engine) PlaceOrderAsync(req book.OrderRequest) {
	e.in <- Command{
		Type:  CmdPlaceOrder,
		Order: req,
	}
}

// CancelOrderAsync submits an asynchronous CancelOrder command.
func (e *Engine) CancelOrderAsync(id book.OrderID) {
	e.in <- Command{
		Type:    CmdCancelOrder,
		OrderID: id,
	}
}

// PlaceOrder submits req and waits for the resulting events, appending into dst[:0].
func (e *Engine) PlaceOrder(req book.OrderRequest, dst []book.Event) []book.Event {
	ch := e.replyPool.Get().(chan Response)
	e.in <- Command{
		Type:     CmdPlaceOrder,
		Order:    req,
		EventDst: dst,
		Reply:    ch,
	}
	resp := <-ch
	e.replyPool.Put(ch)
	return resp.Events
}

// CancelOrder cancels id and waits for the resulting events, appending into dst[:0].
func (e *Engine) CancelOrder(id book.OrderID, dst []book.Event) []book.Event {
	ch := e.replyPool.Get().(chan Response)
	e.in <- Command{
		Type:     CmdCancelOrder,
		OrderID:  id,
		EventDst: dst,
		Reply:    ch,
	}
	resp := <-ch
	e.replyPool.Put(ch)
	return resp.Events
}

// TopOfBook queries the current best bid and best ask through the engine queue.
func (e *Engine) TopOfBook() book.TopOfBook {
	ch := e.replyPool.Get().(chan Response)
	e.in <- Command{
		Type:  CmdTopOfBook,
		Reply: ch,
	}
	resp := <-ch
	e.replyPool.Put(ch)
	return resp.Top
}

// Depth queries the top n bid and ask levels through the engine queue.
func (e *Engine) Depth(n int, bidsDst, asksDst []book.DepthLevel) ([]book.DepthLevel, []book.DepthLevel) {
	ch := e.replyPool.Get().(chan Response)
	e.in <- Command{
		Type:    CmdDepth,
		DepthN:  n,
		BidsDst: bidsDst,
		AsksDst: asksDst,
		Reply:   ch,
	}
	resp := <-ch
	e.replyPool.Put(ch)
	return resp.Bids, resp.Asks
}

// Stats queries current deterministic counters and sizes through the engine queue.
func (e *Engine) Stats() book.BookStats {
	ch := e.replyPool.Get().(chan Response)
	e.in <- Command{
		Type:  CmdStats,
		Reply: ch,
	}
	resp := <-ch
	e.replyPool.Put(ch)
	return resp.Stats
}
