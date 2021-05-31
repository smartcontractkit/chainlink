package utils

import (
	"sync"
)

// Mailbox contains a notify channel,
// a mutual exclusive lock,
// a queue of interfaces,
// and a queue capacity.
type Mailbox struct {
	chNotify chan struct{}
	mu       sync.Mutex
	queue    []interface{}
	capacity uint64
}

// NewMailbox creates a new mailbox instance
func NewMailbox(capacity uint64) *Mailbox {
	queueCap := capacity
	if queueCap == 0 {
		queueCap = 100
	}
	return &Mailbox{
		chNotify: make(chan struct{}, 1),
		queue:    make([]interface{}, 0, queueCap),
		capacity: capacity,
	}
}

// Notify returns the contents of the notify channel
func (m *Mailbox) Notify() chan struct{} {
	return m.chNotify
}

// Deliver appends an interface to the queue
func (m *Mailbox) Deliver(x interface{}) (wasOverCapacity bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queue = append([]interface{}{x}, m.queue...)
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

// Retrieve fetches an interface from the queue
func (m *Mailbox) Retrieve() (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.queue) == 0 {
		return nil, false
	}
	x := m.queue[len(m.queue)-1]
	m.queue = m.queue[:len(m.queue)-1]
	return x, true
}

func (m *Mailbox) RetrieveLatestAndClear() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.queue) == 0 {
		return nil
	}
	x := m.queue[0]
	m.queue = nil
	return x
}
