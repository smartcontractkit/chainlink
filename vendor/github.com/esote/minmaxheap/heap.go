// Package minmaxheap provides min-max heap operations for any type that
// implements heap.Interface. A min-max heap can be used to implement a
// double-ended priority queue.
//
// Min-max heap implementation from the 1986 paper "Min-Max Heaps and
// Generalized Priority Queues" by Atkinson et. al.
// https://doi.org/10.1145/6617.6621.
package minmaxheap

import (
	"container/heap"
	"math/bits"
)

// Interface copied from the heap package, so code that imports minmaxheap does
// not also have to import "container/heap".
type Interface = heap.Interface

func level(i int) int {
	// floor(log2(i + 1))
	return bits.Len(uint(i)+1) - 1
}

func isMinLevel(i int) bool {
	return level(i)%2 == 0
}

func lchild(i int) int {
	return i*2 + 1
}

func rchild(i int) int {
	return i*2 + 2
}

func parent(i int) int {
	return (i - 1) / 2
}

func hasParent(i int) bool {
	return i > 0
}

func hasGrandparent(i int) bool {
	return i > 2
}

func grandparent(i int) int {
	return parent(parent(i))
}

func down(h Interface, i, n int) bool {
	min := isMinLevel(i)
	i0 := i
	for {
		m := i

		l := lchild(i)
		if l >= n || l < 0 /* overflow */ {
			break
		}
		if h.Less(l, m) == min {
			m = l
		}

		r := rchild(i)
		if r < n && h.Less(r, m) == min {
			m = r
		}

		// grandchildren are contiguous i*4+3+{0,1,2,3}
		for g := lchild(l); g < n && g <= rchild(r); g++ {
			if h.Less(g, m) == min {
				m = g
			}
		}

		if m == i {
			break
		}

		h.Swap(i, m)

		if m == l || m == r {
			break
		}

		// m is grandchild
		p := parent(m)
		if h.Less(p, m) == min {
			h.Swap(m, p)
		}
		i = m
	}
	return i > i0
}

func up(h Interface, i int) {
	min := isMinLevel(i)

	if hasParent(i) {
		p := parent(i)
		if h.Less(p, i) == min {
			h.Swap(i, p)
			min = !min
			i = p
		}
	}

	for hasGrandparent(i) {
		g := grandparent(i)
		if h.Less(i, g) != min {
			return
		}

		h.Swap(i, g)
		i = g
	}
}

// Init establishes the heap invariants required by the other routines in this
// package. Init may be called whenever the heap invariants may have been
// invalidated.
// The complexity is O(n) where n = h.Len().
func Init(h Interface) {
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func Push(h Interface, x interface{}) {
	h.Push(x)
	up(h, h.Len()-1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
func Pop(h Interface) interface{} {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n)
	return h.Pop()
}

// PopMax removes and returns the maximum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
func PopMax(h Interface) interface{} {
	n := h.Len()

	i := 0
	l := lchild(0)
	if l < n && !h.Less(l, i) {
		i = l
	}

	r := rchild(0)
	if r < n && !h.Less(r, i) {
		i = r
	}

	h.Swap(i, n-1)
	down(h, i, n-1)
	return h.Pop()
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func Remove(h Interface, i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.Pop()
}

// Fix re-establishes the heap ordering after the element at index i has
// changed its value. Changing the value of the element at index i and then
// calling Fix is equivalent to, but less expensive than, calling Remove(h, i)
// followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func Fix(h Interface, i int) {
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}
