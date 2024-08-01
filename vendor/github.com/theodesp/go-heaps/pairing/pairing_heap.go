// Package pairing implements a Pairing heap Data structure
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Pairing_heap
package pairing

import (
	heap "github.com/theodesp/go-heaps"
	"fmt"
)

// PairHeap implements the Extended interface
var _ heap.Extended = (*PairHeap)(nil)

// PairHeap is an implementation of a Pairing Heap.
// The zero value for PairHeap Root is an empty Heap.
type PairHeap struct {
	root       *node
}

// node contains the current item and the list if the sub-heaps
type node struct {
	// for use by client; untouched by this library
	item heap.Item
	// List of children nodes all containing values less than the Top of the heap
	children []*node
	// A reference to the parent Heap Node
	parent *node
}

func (n *node) detach() []*node {
	if n.parent == nil {
		return nil // avoid detaching root
	}
	for _, node := range n.children {
		node.parent = nil
	}
	var idx int
	for i, node := range n.parent.children {
		if node == n {
			idx = i
			break
		}
	}
	n.parent.children = append(n.parent.children[:idx], n.parent.children[idx+1:]...)
	n.parent = nil
	return n.children
}

func (n *node) iterItem(iter heap.ItemIterator) {
	if !iter(n.item) {
		return
	}
	n.iterChildren(n.children, iter)
}

func (n *node) iterChildren(children []*node, iter heap.ItemIterator) {
	if len(children) == 0 {
		return
	}
	for _, node := range children {
		if !iter(node.item) {
			return
		}
		n.iterChildren(node.children, iter)
	}
}

func (n *node) findNode(item heap.Item) *node {
	if n.item.Compare(item) == 0 {
		return n
	} else {
		return n.findInChildren(n.children, item)
	}
}

func (n *node) findInChildren(children []*node, item heap.Item) *node {
	if len(children) == 0 {
		return nil
	}
	var node *node
loop:
	for _, child := range children {
		node = child.findNode(item)
		if node != nil {
			break loop
		}
	}
	return node
}

// Init initializes or clears the PairHeap
func (p *PairHeap) Init() *PairHeap {
	p.root = &node{}
	return p
}

// New returns an initialized PairHeap.
func New() *PairHeap { return new(PairHeap).Init() }

// IsEmpty returns true if PairHeap p is empty.
// The complexity is O(1).
func (p *PairHeap) IsEmpty() bool {
	return p.root.item == nil
}

// Resets the current PairHeap
func (p *PairHeap) Clear() {
	p.Init()
}

// Find the smallest item in the priority queue.
// The complexity is O(1).
func (p *PairHeap) FindMin() heap.Item {
	if p.IsEmpty() {
		return nil
	}
	return p.root.item
}

// Inserts the value to the PairHeap and returns the item
// The complexity is O(1).
func (p *PairHeap) Insert(item heap.Item) heap.Item {
	p.root = merge(p.root, &node{item: item})
	return item
}


// toDelete details what item to remove in a node call.
type toDelete int

const (
	removeItem toDelete = iota   // removes the given item
	removeMin                  // removes min item in the heap
)

// DeleteMin removes the top most value from the PairHeap and returns it
// The complexity is O(log n) amortized.
func (p *PairHeap) DeleteMin() heap.Item {
	return p.deleteItem(nil, removeMin)
}

// Deletes a node from the heap and returns the item
// The complexity is O(log n) amortized.
func (p *PairHeap) Delete(item heap.Item) heap.Item {
	return p.deleteItem(item, removeItem)
}

func (p *PairHeap) deleteItem(item heap.Item, typ toDelete) heap.Item {
	var result node

	if len(p.root.children) == 0 {
		result = *p.root
		p.root.item = nil
	} else {
		switch typ {
		case removeMin:
			result = *p.root
			p.root = mergePairs(p.root, p.root.children)
		case removeItem:
			node := p.root.findNode(item)
			if node == nil {
				return nil
			} else {
				children := node.detach()
				p.root.children = append(p.root.children, children...)
				result = *node
			}
		default:
			panic("invalid type")
		}
	}

	return result.item
}

// Adjusts the value to the node item and returns it
// The complexity is O(n) amortized.
func (p *PairHeap) Adjust(item, new heap.Item) heap.Item {
	n := p.root.findNode(item)
	if n == nil {
		return nil
	}

	if n == p.root {
		p.DeleteMin()
		return p.Insert(new)
	} else {
		children := n.detach()
		p.Insert(new)
		p.root.children = append(p.root.children, children...)
		return n.item
	}
}

// Exhausting search of the element that matches item and returns it
// The complexity is O(n) amortized.
func (p *PairHeap) Find(item heap.Item) heap.Item {
	if p.IsEmpty() {
		return nil
	}
	var found heap.Item
	p.root.iterItem(func(i heap.Item) bool {
		if item.Compare(i) == 0 {
			found = i
			return false
		} else {
			return true
		}
	})
	return found
}


// Do calls function cb on each element of the PairingHeap, in order of appearance.
// The behavior of Do is undefined if cb changes *p.
func (p *PairHeap) Do(it heap.ItemIterator) {
	if p.IsEmpty() {
		return
	}
	p.root.iterItem(it)
}

// Return the heap formed by taking the union of the item disjoint
// current heap and a that is of the same type
func (p *PairHeap) Meld(a heap.Interface) heap.Interface {
	if a == nil {
		return p
	}
	switch a.(type) {
	case *PairHeap:
		h := a.(*PairHeap)
		if h.IsEmpty() {
			return p
		}
		if p.IsEmpty() {
			p.root = h.root
			h.Clear()
			return p
		}
		if p.FindMin().Compare(h.FindMin()) > 0 {
			h.root = merge(h.root, p.root)
			p.root = h.root
			h.Clear()
		} else {
			p.root = merge(p.root, h.root)
		}

	default:
		panic(fmt.Sprintf("unexpected type %T", a))
	}

	return p
}

func merge(a, b *node) *node {
	if a.item == nil { // Case when root is empty
		a = b
		return a
	}

	if a.item.Compare(b.item) < 0 {
		// put 'second' as the first child of 'first' and update the parent
		a.children = append([]*node{b}, a.children...)
		b.parent = a
		return a
	} else {
		// put 'first' as the first child of 'second' and update the parent
		b.children = append([]*node{a}, b.children...)
		a.parent = b
		return b
	}
}

// Merges heaps together
func mergePairs(root *node, heaps []*node) *node {
	if len(heaps) == 1 {
		root = heaps[0]
		heaps[0].parent = nil
		return root
	}
	var merged *node
	for { // iteratively merge heaps
		if len(heaps) == 0 {
			break
		}
		if merged == nil {
			merged = merge(heaps[0], heaps[1])
			heaps = heaps[2:]
		} else {
			merged = merge(merged, heaps[0])
			heaps = heaps[1:]
		}
	}
	root = merged
	merged.parent = nil

	return root
}
