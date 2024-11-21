package syncer

import "container/heap"

type Heap interface {
	// Push adds a new item to the heap.
	Push(x WorkflowRegistryEventResponse)

	// Pop removes the smallest item from the heap and returns it.
	Pop() WorkflowRegistryEventResponse

	// Len returns the number of items in the heap.
	Len() int
}

// publicHeap is a wrapper around the heap.Interface that exposes the Push and Pop methods.
type publicHeap[T any] struct {
	heap heap.Interface
}

func (h *publicHeap[T]) Push(x T) {
	heap.Push(h.heap, x)
}

func (h *publicHeap[T]) Pop() T {
	return heap.Pop(h.heap).(T)
}

func (h *publicHeap[T]) Len() int {
	return h.heap.Len()
}

// blockHeightHeap is a heap.Interface that sorts WorkflowRegistryEventResponses by block height.
type blockHeightHeap []WorkflowRegistryEventResponse

// newBlockHeightHeap returns an initialized heap that sorts WorkflowRegistryEventResponses by block height.
func newBlockHeightHeap() Heap {
	h := blockHeightHeap(make([]WorkflowRegistryEventResponse, 0))
	heap.Init(&h)
	return &publicHeap[WorkflowRegistryEventResponse]{heap: &h}
}

func (h *blockHeightHeap) Len() int { return len(*h) }

func (h *blockHeightHeap) Less(i, j int) bool {
	return (*h)[i].Event.Head.Height < (*h)[j].Event.Head.Height
}

func (h *blockHeightHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *blockHeightHeap) Push(x any) {
	*h = append(*h, x.(WorkflowRegistryEventResponse))
}

func (h *blockHeightHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
