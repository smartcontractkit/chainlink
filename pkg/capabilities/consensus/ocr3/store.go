package ocr3

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

var (
	errQueueEmpty = errors.New("queue is empty")
)

type store struct {
	requestIDs []string
	requests   map[string]*request

	clock          clockwork.Clock
	requestTimeout time.Duration

	mu sync.RWMutex
}

func newStore(requestTimeout time.Duration, clock clockwork.Clock) *store {
	return &store{
		requestIDs:     []string{},
		requests:       map[string]*request{},
		requestTimeout: requestTimeout,
		clock:          clock,
	}
}

func (s *store) add(ctx context.Context, req *request) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.clock.Now()
	req.ExpiresAt = now.Add(s.requestTimeout)
	s.requestIDs = append(s.requestIDs, req.WorkflowExecutionID)
	s.requests[req.WorkflowExecutionID] = req
	return nil
}

func (s *store) get(ctx context.Context, requestID string) (*request, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.requests[requestID]
	if !ok {
		return nil, fmt.Errorf("request with id %s not found", requestID)
	}

	return r, nil
}

func (s *store) getN(ctx context.Context, requestIDs []string) ([]*request, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	o := []*request{}
	for _, r := range requestIDs {
		gr, ok := s.requests[r]
		if !ok {
			return nil, fmt.Errorf("request with id %s not found", r)
		}

		o = append(o, gr)
	}

	return o, nil
}

func (s *store) firstN(ctx context.Context, batchSize int) ([]*request, error) {
	if batchSize == 0 {
		return nil, errors.New("batchsize cannot be 0")
	}
	if len(s.requestIDs) < 1 {
		return nil, errQueueEmpty
	}

	got := []*request{}
	newRequestIDs := []string{}
	lastIdx := 0
	for idx, r := range s.requestIDs {
		gr, ok := s.requests[r]
		if !ok {
			continue
		}

		got = append(got, gr)
		newRequestIDs = append(newRequestIDs, r)
		lastIdx = idx
		if len(got) == batchSize {
			break
		}
	}

	s.requestIDs = append(newRequestIDs, s.requestIDs[lastIdx+1:]...)
	if len(got) == 0 {
		return nil, errQueueEmpty
	}
	return got, nil
}

func (s *store) evict(ctx context.Context, requestID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.requests[requestID]
	if !ok {
		return false
	}

	delete(s.requests, requestID)
	return true
}
