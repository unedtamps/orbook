package queue

import "github.com/nedo/orbook/book"

const defaultRingCap = 8

// RingBuffer is a contiguous power-of-2 circular FIFO queue using nil-slot
// tombstones for O(1) cancellation via Order.QueuePos.
//
// Tombstones are cleared immediately on Remove (so pooled Order structs can be
// safely recycled without dangling pointers), skipped lazily on PopFront, and
// compacted in-place without allocation when tombstones fill more than half the buffer.
type RingBuffer struct {
	buf   []*book.Order
	head  uint64
	tail  uint64
	count int
}

// NewRingBuffer creates an empty ring-buffer queue.
func NewRingBuffer() book.LevelQueue {
	return &RingBuffer{}
}

// NewRingBufferWithCap creates a ring-buffer queue pre-sized to a power-of-2 capacity >= initCap.
func NewRingBufferWithCap(initCap int) book.LevelQueue {
	c := defaultRingCap
	for c < initCap {
		c <<= 1
	}
	return &RingBuffer{
		buf: make([]*book.Order, c),
	}
}

func (r *RingBuffer) PushBack(o *book.Order) {
	if len(r.buf) == 0 {
		r.buf = make([]*book.Order, defaultRingCap)
	} else if r.tail-r.head == uint64(len(r.buf)) {
		r.compactOrGrow()
	}
	mask := uint64(len(r.buf) - 1)
	r.buf[int(r.tail&mask)] = o
	o.QueuePos = r.tail
	r.tail++
	r.count++
}

func (r *RingBuffer) Front() *book.Order {
	if r.count == 0 {
		return nil
	}
	return r.buf[int(r.head&uint64(len(r.buf)-1))]
}

func (r *RingBuffer) PopFront() *book.Order {
	if r.count == 0 {
		return nil
	}
	mask := uint64(len(r.buf) - 1)
	slot := int(r.head & mask)
	o := r.buf[slot]
	r.buf[slot] = nil
	r.head++
	r.count--
	if r.count == 0 {
		r.head = 0
		r.tail = 0
	} else {
		for r.head < r.tail && r.buf[int(r.head&mask)] == nil {
			r.head++
		}
	}
	return o
}

func (r *RingBuffer) Remove(o *book.Order) {
	mask := uint64(len(r.buf) - 1)
	slot := int(o.QueuePos & mask)
	r.buf[slot] = nil
	r.count--
	if r.count == 0 {
		r.head = 0
		r.tail = 0
		return
	}
	if o.QueuePos == r.head {
		for r.head < r.tail && r.buf[int(r.head&mask)] == nil {
			r.head++
		}
	}
}

func (r *RingBuffer) Len() int {
	return r.count
}

func (r *RingBuffer) Reset() {
	if len(r.buf) > 0 && r.head < r.tail {
		mask := uint64(len(r.buf) - 1)
		for i := r.head; i < r.tail; i++ {
			r.buf[int(i&mask)] = nil
		}
	}
	r.head = 0
	r.tail = 0
	r.count = 0
}

func (r *RingBuffer) ForEach(fn func(*book.Order) bool) {
	if r.count == 0 {
		return
	}
	mask := uint64(len(r.buf) - 1)
	for i := r.head; i < r.tail; i++ {
		o := r.buf[int(i&mask)]
		if o != nil {
			if !fn(o) {
				return
			}
		}
	}
}

func (r *RingBuffer) compactOrGrow() {
	capNow := len(r.buf)
	mask := uint64(capNow - 1)

	// If more than half of the slots are tombstones, compact in-place with zero allocation.
	if r.count*2 < capNow {
		write := r.head
		for read := r.head; read < r.tail; read++ {
			rslot := int(read & mask)
			o := r.buf[rslot]
			if o != nil {
				r.buf[rslot] = nil
				wslot := int(write & mask)
				r.buf[wslot] = o
				o.QueuePos = write
				write++
			}
		}
		r.tail = write
		return
	}

	// Otherwise grow 2x and pack live elements from index 0.
	newBuf := make([]*book.Order, capNow*2)
	var idx uint64
	for read := r.head; read < r.tail; read++ {
		o := r.buf[int(read&mask)]
		if o != nil {
			newBuf[idx] = o
			o.QueuePos = idx
			idx++
		}
	}
	r.buf = newBuf
	r.head = 0
	r.tail = idx
}
