package ocr3

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type store struct {
	requestIDs []string
	requests   map[string]*request

	// for testing
	evictedCh chan *request

	mu sync.RWMutex
}

func newStore() *store {
	return &store{
		requestIDs: []string{},
		requests:   map[string]*request{},
	}
}

func (s *store) add(ctx context.Context, req *request) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.requests[req.WorkflowExecutionID]; ok {
		return fmt.Errorf("request with id %s already exists", req.WorkflowExecutionID)
	}
	s.requestIDs = append(s.requestIDs, req.WorkflowExecutionID)
	s.requests[req.WorkflowExecutionID] = req
	return nil
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if batchSize == 0 {
		return nil, errors.New("batchsize cannot be 0")
	}
	got := []*request{}
	if len(s.requestIDs) == 0 {
		return got, nil
	}

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

	// remove the ones that didn't have corresponding requests
	s.requestIDs = append(newRequestIDs, s.requestIDs[lastIdx+1:]...)
	return got, nil
}

func (s *store) evict(ctx context.Context, requestID string) (*request, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.requests[requestID]
	if !ok {
		return nil, false
	}

	delete(s.requests, requestID)

	// for testing
	if s.evictedCh != nil {
		s.evictedCh <- r
	}

	return r, true
}
