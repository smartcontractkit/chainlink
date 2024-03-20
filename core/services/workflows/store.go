package workflows

import (
	"context"
	"fmt"
	"sync"
)

type store struct {
	idToState map[string]*executionState
	mu        sync.RWMutex
}

func newStore() *store {
	return &store{idToState: map[string]*executionState{}}
}

// add adds a new execution state under the given executionID
func (s *store) add(ctx context.Context, state *executionState) error {
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
func (s *store) updateStep(ctx context.Context, step *stepState) (executionState, error) {
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
func (s *store) updateStatus(ctx context.Context, executionID string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.idToState[executionID]
	if !ok {
		return fmt.Errorf("could not find execution %s", executionID)
	}

	state.status = status
	return nil
}
