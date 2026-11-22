package levels

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nedo/orbook/book"
	"github.com/nedo/orbook/book/queue"
)

func (t *RBTree) checkRBProperties() error {
	if t.nilNode.color != rbBlack {
		return fmt.Errorf("sentinel is not black")
	}
	if t.root == t.nilNode {
		if t.count != 0 || t.best != t.nilNode {
			return fmt.Errorf("empty tree has count=%d best=%p", t.count, t.best)
		}
		return nil
	}
	if t.root.color != rbBlack {
		return fmt.Errorf("root is not black")
	}
	if t.root.parent != t.nilNode {
		return fmt.Errorf("root.parent != nilNode")
	}

	var visited int
	var check func(n *rbNode) (int, error)
	check = func(n *rbNode) (int, error) {
		if n == t.nilNode {
			return 1, nil
		}
		visited++
		if n.left != t.nilNode {
			if n.left.parent != n {
				return 0, fmt.Errorf("left.parent mismatch at price %d", n.level.Price)
			}
			if !t.better(n.left.level.Price, n.level.Price) {
				return 0, fmt.Errorf("BST violation: left %d not better than %d", n.left.level.Price, n.level.Price)
			}
			if n.color == rbRed && n.left.color == rbRed {
				return 0, fmt.Errorf("red-red violation at %d -> left %d", n.level.Price, n.left.level.Price)
			}
		}
		if n.right != t.nilNode {
			if n.right.parent != n {
				return 0, fmt.Errorf("right.parent mismatch at price %d", n.level.Price)
			}
			if !t.better(n.level.Price, n.right.level.Price) {
				return 0, fmt.Errorf("BST violation: %d not better than right %d", n.level.Price, n.right.level.Price)
			}
			if n.color == rbRed && n.right.color == rbRed {
				return 0, fmt.Errorf("red-red violation at %d -> right %d", n.level.Price, n.right.level.Price)
			}
		}
		lh, err := check(n.left)
		if err != nil {
			return 0, err
		}
		rh, err := check(n.right)
		if err != nil {
			return 0, err
		}
		if lh != rh {
			return 0, fmt.Errorf("black-height mismatch at %d: left=%d right=%d", n.level.Price, lh, rh)
		}
		if n.color == rbBlack {
			return lh + 1, nil
		}
		return lh, nil
	}

	if _, err := check(t.root); err != nil {
		return err
	}
	if visited != t.count {
		return fmt.Errorf("visited=%d != count=%d", visited, t.count)
	}
	if min := t.minNode(t.root); min != t.best {
		return fmt.Errorf("best=%v != minNode=%v", t.best.level.Price, min.level.Price)
	}
	return nil
}

func TestRBTreeCLRSInvariantsStress(t *testing.T) {
	for _, side := range []book.Side{book.Buy, book.Sell} {
		t.Run(side.String(), func(t *testing.T) {
			rb := NewRBTree(side, queue.NewLinkedList).(*RBTree)
			rng := rand.New(rand.NewSource(1337))

			for step := 0; step < 5000; step++ {
				p := int64(1000 + rng.Intn(150))
				if rng.Intn(100) < 55 {
					lvl := rb.GetOrCreate(p)
					lvl.TotalQty += 100
					lvl.OrderCount++
				} else {
					rb.Remove(p)
				}
				if step%25 == 0 {
					if err := rb.checkRBProperties(); err != nil {
						t.Fatalf("step %d: %v", step, err)
					}
				}
				if step == 2500 {
					rb.Reset()
					if err := rb.checkRBProperties(); err != nil {
						t.Fatalf("after Reset: %v", err)
					}
				}
			}
			if err := rb.checkRBProperties(); err != nil {
				t.Fatalf("final: %v", err)
			}
		})
	}
}
