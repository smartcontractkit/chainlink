package minheap

import (
	"container/heap"
)

// LessFn defines a less-than relation between a and b, i.e. a < b
type LessFn[T any] func(a, b T) bool

func NewMinHeap[T any](lessFn LessFn[T]) *MinHeap[T] {
	return &MinHeap[T]{minHeapInternal[T]{lessFn, nil}}
}

// Type-safe wrapper around minHeapInternal
type MinHeap[T any] struct {
	internal minHeapInternal[T]
}

func (h *MinHeap[T]) Push(item T) {
	heap.Push(&h.internal, item)
}

func (h *MinHeap[T]) Pop() T {
	return heap.Pop(&h.internal).(T)
}

func (h *MinHeap[T]) Peek() T {
	return h.internal.items[0]
}

func (h *MinHeap[T]) Len() int {
	return h.internal.Len()
}

// Implements heap.Interface and uses interface{} all over the place.
type minHeapInternal[T any] struct {
	lessFn LessFn[T]
	items  []T
}

var _ heap.Interface = new(minHeapInternal[struct{}])

func (hi *minHeapInternal[T]) Len() int { return len(hi.items) }

func (hi *minHeapInternal[T]) Less(i, j int) bool {
	return hi.lessFn(hi.items[i], hi.items[j])
}

func (hi *minHeapInternal[T]) Swap(i, j int) {
	hi.items[i], hi.items[j] = hi.items[j], hi.items[i]
}

func (hi *minHeapInternal[T]) Push(x interface{}) {
	item := x.(T)
	hi.items = append(hi.items, item)
}

func (hi *minHeapInternal[T]) Pop() interface{} {
	old := hi.items
	n := len(old)
	item := old[n-1]
	hi.items = old[0 : n-1]
	return item
}
