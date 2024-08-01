package graph

import (
	"container/heap"
	"errors"
)

// priorityQueue implements a minimum priority queue using a minimum binary heap
// that prioritizes smaller values over larger values.
type priorityQueue[T comparable] struct {
	items *minHeap[T]
	cache map[T]*priorityItem[T]
}

// priorityItem is an item on the binary heap consisting of a priority value and
// an actual payload value.
type priorityItem[T comparable] struct {
	value    T
	priority float64
	index    int
}

func newPriorityQueue[T comparable]() *priorityQueue[T] {
	return &priorityQueue[T]{
		items: &minHeap[T]{},
		cache: map[T]*priorityItem[T]{},
	}
}

// Len returns the total number of items in the priority queue.
func (p *priorityQueue[T]) Len() int {
	return p.items.Len()
}

// Push pushes a new item with the given priority into the queue. This operation
// may cause a re-balance of the heap and thus scales with O(log n).
func (p *priorityQueue[T]) Push(item T, priority float64) {
	if _, ok := p.cache[item]; ok {
		return
	}

	newItem := &priorityItem[T]{
		value:    item,
		priority: priority,
		index:    0,
	}

	heap.Push(p.items, newItem)
	p.cache[item] = newItem
}

// Pop returns and removes the item with the lowest priority. This operation may
// cause a re-balance of the heap and thus scales with O(log n).
func (p *priorityQueue[T]) Pop() (T, error) {
	if len(*p.items) == 0 {
		var empty T
		return empty, errors.New("priority queue is empty")
	}

	item := heap.Pop(p.items).(*priorityItem[T])
	delete(p.cache, item.value)

	return item.value, nil
}

// UpdatePriority updates the priority of a given item and sets it to the given
// priority. If the item doesn't exist, nothing happens. This operation may
// cause a re-balance of the heap and this scales with O(log n).
func (p *priorityQueue[T]) UpdatePriority(item T, priority float64) {
	targetItem, ok := p.cache[item]
	if !ok {
		return
	}

	targetItem.priority = priority
	heap.Fix(p.items, targetItem.index)
}

// minHeap is a minimum binary heap that implements heap.Interface.
type minHeap[T comparable] []*priorityItem[T]

func (m *minHeap[T]) Len() int {
	return len(*m)
}

func (m *minHeap[T]) Less(i, j int) bool {
	return (*m)[i].priority < (*m)[j].priority
}

func (m *minHeap[T]) Swap(i, j int) {
	(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
	(*m)[i].index = i
	(*m)[j].index = j
}

func (m *minHeap[T]) Push(item interface{}) {
	i := item.(*priorityItem[T])
	i.index = len(*m)
	*m = append(*m, i)
}

func (m *minHeap[T]) Pop() interface{} {
	old := *m
	item := old[len(old)-1]
	*m = old[:len(old)-1]

	return item
}

type stack[T any] interface {
	push(T)
	pop() (T, error)
	top() (T, error)
	isEmpty() bool
	// forEach iterate the stack from bottom to top
	forEach(func(T))
}

func newStack[T any]() stack[T] {
	return &stackImpl[T]{
		elements: make([]T, 0),
	}
}

type stackImpl[T any] struct {
	elements []T
}

func (s *stackImpl[T]) push(t T) {
	s.elements = append(s.elements, t)
}

func (s *stackImpl[T]) pop() (T, error) {
	e, err := s.top()
	if err != nil {
		var defaultValue T
		return defaultValue, err
	}

	s.elements = s.elements[:len(s.elements)-1]
	return e, nil
}

func (s *stackImpl[T]) top() (T, error) {
	if s.isEmpty() {
		var defaultValue T
		return defaultValue, errors.New("no element in stack")
	}

	return s.elements[len(s.elements)-1], nil
}

func (s *stackImpl[T]) isEmpty() bool {
	return len(s.elements) == 0
}

func (s *stackImpl[T]) forEach(f func(T)) {
	for _, e := range s.elements {
		f(e)
	}
}
