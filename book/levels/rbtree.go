package levels

import "github.com/unedtamps/orbook/book"

type rbColor uint8

const (
	rbBlack rbColor = iota
	rbRed
)

type rbNode struct {
	level  book.PriceLevel
	color  rbColor
	parent *rbNode
	left   *rbNode
	right  *rbNode
}

// RBTree is a hand-rolled CLRS Red-Black Tree storing PriceLevels ordered
// best-to-worst (leftmost node is always the highest-priority level).
// It embeds PriceLevel directly inside each node, caches the best level pointer
// for O(1) Best(), iterates TopN/ForEach in O(k) via parent pointers with zero
// allocations, and recycles nodes (with their LevelQueues) via a free list.
type RBTree struct {
	side     book.Side
	qFactory book.QueueFactory
	sentinel rbNode
	nilNode  *rbNode
	root     *rbNode
	best     *rbNode
	count    int
	free     []*rbNode
}

// NewRBTree constructs a Red-Black Tree price-level store for the given side.
func NewRBTree(side book.Side, qFactory book.QueueFactory) book.PriceLevels {
	t := &RBTree{
		side:     side,
		qFactory: qFactory,
		free:     make([]*rbNode, 0, 64),
	}
	t.nilNode = &t.sentinel
	t.sentinel.color = rbBlack
	t.sentinel.parent = t.nilNode
	t.sentinel.left = t.nilNode
	t.sentinel.right = t.nilNode
	t.root = t.nilNode
	t.best = t.nilNode
	return t
}

// better reports whether price a has strictly higher priority than price b.
func (t *RBTree) better(a, b int64) bool {
	if t.side == book.Buy {
		return a > b
	}
	return a < b
}

func (t *RBTree) acquireNode(price int64) *rbNode {
	var n *rbNode
	if k := len(t.free); k > 0 {
		n = t.free[k-1]
		t.free[k-1] = nil
		t.free = t.free[:k-1]
	} else {
		n = &rbNode{}
		n.level.Queue = t.qFactory()
	}
	n.level.Price = price
	n.level.TotalQty = 0
	n.level.OrderCount = 0
	n.color = rbRed
	n.parent = t.nilNode
	n.left = t.nilNode
	n.right = t.nilNode
	return n
}

func (t *RBTree) releaseNode(n *rbNode) {
	n.level.Queue.Reset()
	n.level.Price = 0
	n.level.TotalQty = 0
	n.level.OrderCount = 0
	n.parent = nil
	n.left = nil
	n.right = nil
	t.free = append(t.free, n)
}

func (t *RBTree) findNode(price int64) *rbNode {
	if t.best != t.nilNode && t.best.level.Price == price {
		return t.best
	}
	cur := t.root
	if t.side == book.Buy {
		for cur != t.nilNode {
			p := cur.level.Price
			if price == p {
				return cur
			}
			if price > p {
				cur = cur.left
			} else {
				cur = cur.right
			}
		}
	} else {
		for cur != t.nilNode {
			p := cur.level.Price
			if price == p {
				return cur
			}
			if price < p {
				cur = cur.left
			} else {
				cur = cur.right
			}
		}
	}
	return t.nilNode
}

func (t *RBTree) minNode(x *rbNode) *rbNode {
	for x.left != t.nilNode {
		x = x.left
	}
	return x
}

func (t *RBTree) successor(x *rbNode) *rbNode {
	if x.right != t.nilNode {
		return t.minNode(x.right)
	}
	y := x.parent
	for y != t.nilNode && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

func (t *RBTree) leftRotate(x *rbNode) {
	y := x.right
	x.right = y.left
	if y.left != t.nilNode {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == t.nilNode {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (t *RBTree) rightRotate(y *rbNode) {
	x := y.left
	y.left = x.right
	if x.right != t.nilNode {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == t.nilNode {
		t.root = x
	} else if y == y.parent.right {
		y.parent.right = x
	} else {
		y.parent.left = x
	}
	x.right = y
	y.parent = x
}

func (t *RBTree) insertFixup(z *rbNode) {
	for z.parent.color == rbRed {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == rbRed {
				z.parent.color = rbBlack
				y.color = rbBlack
				z.parent.parent.color = rbRed
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = rbBlack
				z.parent.parent.color = rbRed
				t.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y.color == rbRed {
				z.parent.color = rbBlack
				y.color = rbBlack
				z.parent.parent.color = rbRed
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = rbBlack
				z.parent.parent.color = rbRed
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = rbBlack
}

func (t *RBTree) transplant(u, v *rbNode) {
	if u.parent == t.nilNode {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func (t *RBTree) deleteFixup(x *rbNode) {
	for x != t.root && x.color == rbBlack {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == rbRed {
				w.color = rbBlack
				x.parent.color = rbRed
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == rbBlack && w.right.color == rbBlack {
				w.color = rbRed
				x = x.parent
			} else {
				if w.right.color == rbBlack {
					w.left.color = rbBlack
					w.color = rbRed
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = rbBlack
				w.right.color = rbBlack
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == rbRed {
				w.color = rbBlack
				x.parent.color = rbRed
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.right.color == rbBlack && w.left.color == rbBlack {
				w.color = rbRed
				x = x.parent
			} else {
				if w.left.color == rbBlack {
					w.right.color = rbBlack
					w.color = rbRed
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = rbBlack
				w.left.color = rbBlack
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = rbBlack
}

func (t *RBTree) Get(price int64) *book.PriceLevel {
	n := t.findNode(price)
	if n == t.nilNode {
		return nil
	}
	return &n.level
}

func (t *RBTree) GetOrCreate(price int64) *book.PriceLevel {
	if t.best != t.nilNode && t.best.level.Price == price {
		return &t.best.level
	}

	y := t.nilNode
	x := t.root
	goLeft := false
	if t.side == book.Buy {
		for x != t.nilNode {
			y = x
			p := x.level.Price
			if price == p {
				return &x.level
			}
			if price > p {
				x = x.left
				goLeft = true
			} else {
				x = x.right
				goLeft = false
			}
		}
	} else {
		for x != t.nilNode {
			y = x
			p := x.level.Price
			if price == p {
				return &x.level
			}
			if price < p {
				x = x.left
				goLeft = true
			} else {
				x = x.right
				goLeft = false
			}
		}
	}

	z := t.acquireNode(price)
	z.parent = y
	if y == t.nilNode {
		t.root = z
	} else if goLeft {
		y.left = z
	} else {
		y.right = z
	}

	t.insertFixup(z)
	t.count++

	if t.best == t.nilNode || t.better(price, t.best.level.Price) {
		t.best = z
	}
	return &z.level
}

func (t *RBTree) Remove(price int64) {
	z := t.findNode(price)
	if z == t.nilNode {
		return
	}

	// Because t.best is the leftmost node, its left child is always nilNode,
	// so its in-order successor is invariant under RB-tree rotations.
	removingBest := z == t.best
	var nextBest *rbNode
	if removingBest {
		nextBest = t.successor(z)
	}

	y := z
	yOrigColor := y.color
	var x *rbNode

	if z.left == t.nilNode {
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == t.nilNode {
		x = z.left
		t.transplant(z, z.left)
	} else {
		y = t.minNode(z.right)
		yOrigColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yOrigColor == rbBlack {
		t.deleteFixup(x)
	}
	// Restore sentinel parent cleanly.
	t.sentinel.parent = t.nilNode

	t.count--
	if t.count == 0 {
		t.best = t.nilNode
	} else if removingBest {
		t.best = nextBest
	}
	t.releaseNode(z)
}

func (t *RBTree) Best() *book.PriceLevel {
	if t.best == t.nilNode {
		return nil
	}
	return &t.best.level
}

func (t *RBTree) Len() int {
	return t.count
}

func (t *RBTree) TopN(n int, dst []book.DepthLevel) []book.DepthLevel {
	if n <= 0 || t.best == t.nilNode {
		return dst
	}
	count := 0
	for cur := t.best; cur != t.nilNode && count < n; cur = t.successor(cur) {
		dst = append(dst, book.DepthLevel{
			Price:      cur.level.Price,
			TotalQty:   cur.level.TotalQty,
			OrderCount: cur.level.OrderCount,
		})
		count++
	}
	return dst
}

func (t *RBTree) ForEach(fn func(*book.PriceLevel) bool) {
	for cur := t.best; cur != t.nilNode; cur = t.successor(cur) {
		if !fn(&cur.level) {
			return
		}
	}
}

func (t *RBTree) Reset() {
	if t.count == 0 {
		return
	}
	startIdx := len(t.free)
	for cur := t.best; cur != t.nilNode; cur = t.successor(cur) {
		t.free = append(t.free, cur)
	}
	for i := startIdx; i < len(t.free); i++ {
		n := t.free[i]
		n.level.Queue.Reset()
		n.level.Price = 0
		n.level.TotalQty = 0
		n.level.OrderCount = 0
		n.parent = nil
		n.left = nil
		n.right = nil
	}
	t.root = t.nilNode
	t.best = t.nilNode
	t.sentinel.parent = t.nilNode
	t.count = 0
}
