package utils

import (
	"sync"
)

// Mailbox contains a notify channel,
// a mutual exclusive lock,
// a queue of interfaces,
// and a queue capacity.
type Mailbox[T any] struct {
	chNotify chan struct{}
	mu       sync.Mutex
	queue    []T

	// capacity - number of items the mailbox can buffer
	// NOTE: if the capacity is 1, it's possible that an empty Retrieve may occur after a notification.
	capacity uint64
}

// NewHighCapacityMailbox create a new mailbox with a capacity
// that is better able to handle e.g. large log replays
func NewHighCapacityMailbox[T any]() *Mailbox[T] {
	return NewMailbox[T](100000)
}

// NewMailbox creates a new mailbox instance
func NewMailbox[T any](capacity uint64) *Mailbox[T] {
	queueCap := capacity
	if queueCap == 0 {
		queueCap = 100
	}
	return &Mailbox[T]{
		chNotify: make(chan struct{}, 1),
		queue:    make([]T, 0, queueCap),
		capacity: capacity,
	}
}

// Notify returns the contents of the notify channel
func (m *Mailbox[T]) Notify() chan struct{} {
	return m.chNotify
}

// Deliver appends to the queue
func (m *Mailbox[T]) Deliver(x T) (wasOverCapacity bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queue = append([]T{x}, m.queue...)
	if uint64(len(m.queue)) > m.capacity && m.capacity > 0 {
		m.queue = m.queue[:len(m.queue)-1]
		wasOverCapacity = true
	}

	select {
	case m.chNotify <- struct{}{}:
	default:
	}
	return
}

// Retrieve fetches from the queue
func (m *Mailbox[T]) Retrieve() (t T, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.queue) == 0 {
		return
	}
	t = m.queue[len(m.queue)-1]
	m.queue = m.queue[:len(m.queue)-1]
	ok = true
	return
}

// RetrieveLatestAndClear returns the latest value (or nil), and clears the queue.
func (m *Mailbox[T]) RetrieveLatestAndClear() (t T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.queue) == 0 {
		return
	}
	t = m.queue[0]
	m.queue = nil
	return
}
