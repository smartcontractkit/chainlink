package testsinkstateful

import (
	"sync"
	"time"
)

// CapturedEvent is the test-facing representation of a ChipIngress event.
type CapturedEvent struct {
	Timestamp time.Time      `json:"timestamp"`
	Domain    string         `json:"domain"`
	Schema    string         `json:"schema"`
	Entity    string         `json:"entity"` // may be empty (orphan)
	Body      []byte         `json:"-"`      // raw protobuf payload
	Attrs     map[string]any `json:"attrs"`  // attributes, as a simple map
}

type Store struct {
	mu       sync.Mutex
	capacity int
	events   []CapturedEvent
	head     int // index of the oldest element
	count    int // number of stored elements
}

func NewStore(size int) (*Store, error) {
	if size <= 0 {
		size = 10000
	}
	return &Store{
		capacity: size,
		events:   make([]CapturedEvent, size),
	}, nil
}

func (s *Store) Add(ev CapturedEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// copy body to avoid aliasing
	if ev.Body != nil {
		bodyCopy := make([]byte, len(ev.Body))
		copy(bodyCopy, ev.Body)
		ev.Body = bodyCopy
	}

	ev.Timestamp = time.Now()

	// write new event right after the logical tail, wrapping around when needed
	idx := (s.head + s.count) % s.capacity
	s.events[idx] = ev
	if s.count == s.capacity {
		// buffer full: advance head to drop oldest element (FIFO eviction)
		s.head = (s.head + 1) % s.capacity
	} else {
		s.count++
	}
}

func (s *Store) All() []CapturedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]CapturedEvent, 0, s.count)
	// iterate from oldest to newest, wrapping via modulo to preserve order
	for i := 0; i < s.count; i++ {
		idx := (s.head + i) % s.capacity
		out = append(out, s.events[idx])
	}
	return out
}
