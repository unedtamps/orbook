package engine

import (
	"sync"

	"github.com/unedtamps/orbook/book"
)

// BatchEngine is a multi-producer-safe channel runner that amortizes Go channel
// lock contention across slices of Commands (e.g. 32-128 commands per send).
type BatchEngine struct {
	book     *book.Book
	in       chan []Command
	done     chan struct{}
	handler  EventHandler
	eventBuf []book.Event

	startOnce sync.Once
	stopOnce  sync.Once
}

// NewBatchEngine constructs a batched channel runner owning b.
func NewBatchEngine(b *book.Book, opts Options) *BatchEngine {
	qSize := opts.QueueSize
	if qSize <= 0 {
		qSize = 256
	}
	evCap := opts.EventBufCap
	if evCap <= 0 {
		evCap = 64
	}
	return &BatchEngine{
		book:     b,
		in:       make(chan []Command, qSize),
		done:     make(chan struct{}),
		handler:  opts.EventHandler,
		eventBuf: make([]book.Event, 0, evCap),
	}
}

func (e *BatchEngine) Start() {
	e.startOnce.Do(func() {
		go e.run()
	})
}

func (e *BatchEngine) Stop() {
	e.stopOnce.Do(func() {
		close(e.in)
		<-e.done
	})
}

// SubmitBatch sends a slice of commands in a single channel operation.
// Callers must not mutate the slice after sending until processed (or pass a fresh/ring slice).
func (e *BatchEngine) SubmitBatch(batch []Command) {
	e.in <- batch
}

func (e *BatchEngine) run() {
	defer close(e.done)
	for batch := range e.in {
		for i := range batch {
			cmd := &batch[i]
			switch cmd.Type {
			case CmdPlaceOrder:
				e.eventBuf = e.book.PlaceOrder(cmd.Order, e.eventBuf[:0])
				if e.handler != nil {
					e.handler(e.eventBuf)
				}
			case CmdCancelOrder:
				e.eventBuf = e.book.CancelOrder(cmd.OrderID, e.eventBuf[:0])
				if e.handler != nil {
					e.handler(e.eventBuf)
				}
			case CmdStats:
				if cmd.Reply != nil {
					cmd.Reply <- Response{Stats: e.book.Stats()}
				}
			}
		}
	}
}
