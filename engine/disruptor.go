package engine

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/unedtamps/orbook/book"
)

// Disruptor is a true LMAX-style lock-free single-producer single-consumer (SPSC)
// ring-buffer engine runner. Cache-line padding separates producer and consumer
// atomic cursors to eliminate cross-core false sharing, and the consumer drains
// available bursts before updating its head cursor.
type Disruptor struct {
	book     *book.Book
	buf      []Command
	mask     uint64
	handler  EventHandler
	eventBuf []book.Event

	_          [64]byte
	tail       atomic.Uint64
	_          [56]byte
	cachedHead uint64
	_          [56]byte

	head       atomic.Uint64
	_          [56]byte
	cachedTail uint64
	_          [56]byte

	closed    atomic.Bool
	done      chan struct{}
	replyPool sync.Pool
	startOnce sync.Once
	stopOnce  sync.Once
}

// NewDisruptor creates a lock-free SPSC Disruptor runner owning b.
// QueueSize is rounded up to a power of 2 (default 8192).
func NewDisruptor(b *book.Book, opts Options) *Disruptor {
	size := uint64(8192)
	if opts.QueueSize > 0 {
		size = 1
		for size < uint64(opts.QueueSize) {
			size <<= 1
		}
	}
	evCap := opts.EventBufCap
	if evCap <= 0 {
		evCap = 64
	}
	d := &Disruptor{
		book:     b,
		buf:      make([]Command, size),
		mask:     size - 1,
		handler:  opts.EventHandler,
		eventBuf: make([]book.Event, 0, evCap),
		done:     make(chan struct{}),
	}
	d.replyPool.New = func() any {
		return make(chan Response, 1)
	}
	return d
}

// Start launches the single consumer goroutine.
func (d *Disruptor) Start() {
	d.startOnce.Do(func() {
		go d.run()
	})
}

// Stop signals completion (after the single producer stops publishing) and waits
// until all published commands have been drained and processed.
func (d *Disruptor) Stop() {
	d.stopOnce.Do(func() {
		d.closed.Store(true)
		<-d.done
	})
}

// Publish enqueues cmd from the single producer goroutine without locks.
func (d *Disruptor) Publish(cmd Command) {
	t := d.tail.Load()
	capRing := d.mask + 1
	for t-d.cachedHead >= capRing {
		d.cachedHead = d.head.Load()
		if t-d.cachedHead >= capRing {
			runtime.Gosched()
		}
	}
	d.buf[t&d.mask] = cmd
	d.tail.Store(t + 1)
}

// PlaceOrderAsync publishes an asynchronous PlaceOrder command.
func (d *Disruptor) PlaceOrderAsync(req book.OrderRequest) {
	d.Publish(Command{Type: CmdPlaceOrder, Order: req})
}

// CancelOrderAsync publishes an asynchronous CancelOrder command.
func (d *Disruptor) CancelOrderAsync(id book.OrderID) {
	d.Publish(Command{Type: CmdCancelOrder, OrderID: id})
}

// Stats queries book stats synchronously (also acts as a drain barrier).
func (d *Disruptor) Stats() book.BookStats {
	ch := d.replyPool.Get().(chan Response)
	d.Publish(Command{Type: CmdStats, Reply: ch})
	resp := <-ch
	d.replyPool.Put(ch)
	return resp.Stats
}

func (d *Disruptor) run() {
	defer close(d.done)
	mask := d.mask
	h := uint64(0)
	spin := 0

	for {
		avail := d.cachedTail
		if h == avail {
			avail = d.tail.Load()
			d.cachedTail = avail
			if h == avail {
				if d.closed.Load() {
					avail = d.tail.Load()
					if h == avail {
						return
					}
					d.cachedTail = avail
				} else {
					spin++
					if spin > 64 {
						runtime.Gosched()
						spin = 0
					}
					continue
				}
			}
		}
		spin = 0

		for h < avail {
			d.process(&d.buf[h&mask])
			h++
		}
		d.head.Store(h)
	}
}

func (d *Disruptor) process(cmd *Command) {
	switch cmd.Type {
	case CmdPlaceOrder:
		d.eventBuf = d.book.PlaceOrder(cmd.Order, d.eventBuf[:0])
		if d.handler != nil {
			d.handler(d.eventBuf)
		}
		if cmd.Reply != nil {
			dst := append(cmd.EventDst[:0], d.eventBuf...)
			cmd.Reply <- Response{Events: dst}
		}
	case CmdCancelOrder:
		d.eventBuf = d.book.CancelOrder(cmd.OrderID, d.eventBuf[:0])
		if d.handler != nil {
			d.handler(d.eventBuf)
		}
		if cmd.Reply != nil {
			dst := append(cmd.EventDst[:0], d.eventBuf...)
			cmd.Reply <- Response{Events: dst}
		}
	case CmdTopOfBook:
		if cmd.Reply != nil {
			cmd.Reply <- Response{Top: d.book.TopOfBook()}
		}
	case CmdDepth:
		if cmd.Reply != nil {
			bids, asks := d.book.Depth(cmd.DepthN, cmd.BidsDst, cmd.AsksDst)
			cmd.Reply <- Response{Bids: bids, Asks: asks}
		}
	case CmdStats:
		if cmd.Reply != nil {
			cmd.Reply <- Response{Stats: d.book.Stats()}
		}
	}
}
