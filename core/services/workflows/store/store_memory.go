package store

import (
	"context"
	"fmt"
	"sync"
)

// `InMemoryStore` is a temporary in-memory
// equivalent of the database table that should persist
// workflow progress.
type InMemoryStore struct {
	idToState map[string]*WorkflowExecution
	mu        sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{idToState: map[string]*WorkflowExecution{}}
}

// Add adds a new execution state under the given executionID
func (s *InMemoryStore) Add(ctx context.Context, state *WorkflowExecution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.idToState[state.ExecutionID]
	if ok {
		return fmt.Errorf("execution ID %s already exists in store", state.ExecutionID)
	}

	s.idToState[state.ExecutionID] = state
	return nil
}

// UpsertStep updates a step for the given executionID
func (s *InMemoryStore) UpsertStep(ctx context.Context, step *WorkflowExecutionStep) (WorkflowExecution, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.idToState[step.ExecutionID]
	if !ok {
		return WorkflowExecution{}, fmt.Errorf("could not find execution %s", step.ExecutionID)
	}

	state.Steps[step.Ref] = step
	return *state, nil
}

// UpdateStatus updates the status for the given executionID
func (s *InMemoryStore) UpdateStatus(ctx context.Context, executionID string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.idToState[executionID]
	if !ok {
		return fmt.Errorf("could not find execution %s", executionID)
	}

	state.Status = status
	return nil
}

// Get gets the state for the given executionID
func (s *InMemoryStore) Get(ctx context.Context, executionID string) (WorkflowExecution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.idToState[executionID]
	if !ok {
		return WorkflowExecution{}, fmt.Errorf("could not find execution %s", executionID)
	}

	return *state, nil
}

// GetUnfinished gets the states for execution that are in a started state
// Offset and limit are ignored for the in-memory store.
func (s *InMemoryStore) GetUnfinished(ctx context.Context, offset, limit int) ([]WorkflowExecution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	states := []WorkflowExecution{}
	for _, s := range s.idToState {
		if s.Status == StatusStarted {
			states = append(states, *s)
		}
	}

	return states, nil
}
