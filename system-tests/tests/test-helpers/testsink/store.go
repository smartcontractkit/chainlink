package main

import "sync"

// CapturedEvent is the test-facing representation of a ChipIngress event.
type CapturedEvent struct {
	Domain string                 `json:"domain"`
	Schema string                 `json:"schema"`
	Entity string                 `json:"entity"` // may be empty (orphan)
	Body   []byte                 `json:"-"`      // raw protobuf payload
	Attrs  map[string]interface{} `json:"attrs"`  // attributes, as a simple map
}

type Store struct {
	mu     sync.Mutex
	events []CapturedEvent
}

func NewStore() *Store {
	return &Store{}
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
	s.events = append(s.events, ev)
}

func (s *Store) All() []CapturedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]CapturedEvent, len(s.events))
	copy(out, s.events)
	return out
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = nil
}
