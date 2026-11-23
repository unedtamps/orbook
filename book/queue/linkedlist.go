package queue

import "github.com/unedtamps/orbook/book"

// LinkedList is an intrusive doubly-linked FIFO queue backed by Order.Prev and Order.Next.
// All operations (PushBack, Front, PopFront, Remove) are O(1) with zero per-queue slice allocation.
type LinkedList struct {
	head *book.Order
	tail *book.Order
	len  int
}

// NewLinkedList creates an empty intrusive doubly-linked list queue.
func NewLinkedList() book.LevelQueue {
	return &LinkedList{}
}

func (l *LinkedList) PushBack(o *book.Order) {
	o.Prev = l.tail
	o.Next = nil
	if l.tail != nil {
		l.tail.Next = o
	} else {
		l.head = o
	}
	l.tail = o
	l.len++
}

func (l *LinkedList) Front() *book.Order {
	return l.head
}

func (l *LinkedList) PopFront() *book.Order {
	o := l.head
	if o == nil {
		return nil
	}
	l.head = o.Next
	if l.head != nil {
		l.head.Prev = nil
	} else {
		l.tail = nil
	}
	o.Prev = nil
	o.Next = nil
	l.len--
	return o
}

func (l *LinkedList) Remove(o *book.Order) {
	if o.Prev != nil {
		o.Prev.Next = o.Next
	} else {
		l.head = o.Next
	}
	if o.Next != nil {
		o.Next.Prev = o.Prev
	} else {
		l.tail = o.Prev
	}
	o.Prev = nil
	o.Next = nil
	l.len--
}

func (l *LinkedList) Len() int {
	return l.len
}

func (l *LinkedList) Reset() {
	for cur := l.head; cur != nil; {
		next := cur.Next
		cur.Prev = nil
		cur.Next = nil
		cur = next
	}
	l.head = nil
	l.tail = nil
	l.len = 0
}

func (l *LinkedList) ForEach(fn func(*book.Order) bool) {
	for cur := l.head; cur != nil; cur = cur.Next {
		if !fn(cur) {
			return
		}
	}
}
