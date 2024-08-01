package ringbuffer

import "fmt"

// RingBuffer implements a fixed capacity ringbuffer for items of type
// T
type RingBuffer[T any] struct {
	start  int
	length int
	buffer []T
}

func NewRingBuffer[T any](cap int) *RingBuffer[T] {
	if cap <= 0 {
		panic(fmt.Sprintf("NewRingBuffer: cap must be positive, got %d", cap))
	}
	return &RingBuffer[T]{
		0,
		0,
		make([]T, cap),
	}
}

func (rb *RingBuffer[T]) Length() int {
	return rb.length
}

// Peek at the front item. Panics if there isn't one.
func (rb *RingBuffer[T]) Peek() T {
	if rb.length == 0 {
		panic("Peek: buffer empty")
	}
	return rb.buffer[rb.start]
}

// Pop front item. Panics if there isn't one.
func (rb *RingBuffer[T]) Pop() T {
	result := rb.Peek()
	var zero T
	rb.buffer[rb.start] = zero
	rb.start = (rb.start + 1) % len(rb.buffer)
	rb.length--
	return result
}

// Push new item to back. If the additional item would lead
// to the capacity being exceeded, remove the front item first
func (rb *RingBuffer[T]) Push(item T) {
	if rb.length == len(rb.buffer) {
		rb.Pop()
	}
	rb.buffer[(rb.start+rb.length)%len(rb.buffer)] = item
	rb.length++
}
