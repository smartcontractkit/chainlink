package main

import (
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

// CapturedEvent is the test-facing representation of a ChipIngress event.
type CapturedEvent struct {
	Sequence uint64         `json:"sequence"`
	Domain   string         `json:"domain"`
	Schema   string         `json:"schema"`
	Entity   string         `json:"entity"` // may be empty (orphan)
	Body     []byte         `json:"-"`      // raw protobuf payload
	Attrs    map[string]any `json:"attrs"`  // attributes, as a simple map
}

type Store struct {
	mu       sync.Mutex
	cache    *lru.Cache[uint64, CapturedEvent]
	sequence uint64
}

func NewStore(size int) (*Store, error) {
	if size <= 0 {
		size = 10000
	}
	cache, err := lru.New[uint64, CapturedEvent](size)
	if err != nil {
		return nil, fmt.Errorf("create LRU cache: %w", err)
	}
	return &Store{
		cache: cache,
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

	// Increment sequence and assign it to the event
	s.sequence++
	ev.Sequence = s.sequence

	// Add to LRU cache
	s.cache.Add(s.sequence, ev)
}

func (s *Store) All() []CapturedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := s.cache.Keys()
	out := make([]CapturedEvent, 0, len(keys))
	for _, key := range keys {
		if ev, ok := s.cache.Get(key); ok {
			out = append(out, ev)
		}
	}
	return out
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache.Purge()
	s.sequence = 0
}
