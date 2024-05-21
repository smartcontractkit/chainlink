package workflows

import (
	"context"
	"fmt"
	"sync"
)

// `inMemoryStore` is a temporary in-memory
// equivalent of the database table that should persist
// workflow progress.
type inMemoryStore struct {
	idToState map[string]*executionState
	mu        sync.RWMutex
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{idToState: map[string]*executionState{}}
}

// add adds a new execution state under the given executionID
func (s *inMemoryStore) add(ctx context.Context, state *executionState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.idToState[state.executionID]
	if ok {
		return fmt.Errorf("execution ID %s already exists in store", state.executionID)
	}

	s.idToState[state.executionID] = state
	return nil
}

// updateStep updates a step for the given executionID
func (s *inMemoryStore) updateStep(ctx context.Context, step *stepState) (executionState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.idToState[step.executionID]
	if !ok {
		return executionState{}, fmt.Errorf("could not find execution %s", step.executionID)
	}

	state.steps[step.ref] = step
	return *state, nil
}

// updateStatus updates the status for the given executionID
func (s *inMemoryStore) updateStatus(ctx context.Context, executionID string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.idToState[executionID]
	if !ok {
		return fmt.Errorf("could not find execution %s", executionID)
	}

	state.status = status
	return nil
}

// get gets the state for the given executionID
func (s *inMemoryStore) get(ctx context.Context, executionID string) (executionState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.idToState[executionID]
	if !ok {
		return executionState{}, fmt.Errorf("could not find execution %s", executionID)
	}

	return *state, nil
}
