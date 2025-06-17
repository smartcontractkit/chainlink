package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
)

const (
	// defaultPruneInterval is the default interval between pruning completed executions
	defaultPruneInterval = 30 * time.Second

	// maximumExecutionAge is the default maximum age of an execution before it is considered expired and eligible for pruning
	// regardless of its status
	maximumExecutionAge = 1 * time.Hour
)

// InMemoryStore is an in-memory implementation of the Store interface used to store workflow execution states.
// The store always returns a copy of the current workflow execution state in the store such that it is effectively an
// immutable object as state modification only occurs within the store.
// TODO make the WorkflowExecution type immutable to reflect the latter fact and prevent unexpected side effects from
// TODO code being added that modifies WorkflowExecution objects outside of the store. (https://smartcontract-it.atlassian.net/browse/CAPPL-682)
type InMemoryStore struct {
	lggr logger.Logger
	commonservices.StateMachine
	idToExecution     map[string]*WorkflowExecution
	mu                sync.RWMutex
	shutdownWaitGroup sync.WaitGroup
	chStop            commonservices.StopChan

	clock clockwork.Clock

	// pruneInterval is the interval between pruning completed (and expired) executions
	pruneInterval time.Duration

	// maximumExecutionAge is the maximum age of an execution before it is considered expired and eligible for pruning
	// regardless of its status
	maximumExecutionAge time.Duration
}

func NewInMemoryStore(lggr logger.Logger, clock clockwork.Clock) *InMemoryStore {
	return NewInMemoryStoreWithPruneConfiguration(lggr, clock, defaultPruneInterval, maximumExecutionAge)
}

func NewInMemoryStoreWithPruneConfiguration(lggr logger.Logger, clock clockwork.Clock, pruneFrequency time.Duration,
	maximumExecutionAge time.Duration) *InMemoryStore {
	return &InMemoryStore{lggr: lggr, idToExecution: map[string]*WorkflowExecution{}, clock: clock, chStop: make(chan struct{}),
		pruneInterval: pruneFrequency, maximumExecutionAge: maximumExecutionAge}
}

// Add adds a new execution state under the given executionID
func (s *InMemoryStore) Add(ctx context.Context, steps map[string]*WorkflowExecutionStep,
	executionID string, workflowID string, status string) (WorkflowExecution, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.idToExecution[executionID]
	if ok {
		return WorkflowExecution{}, fmt.Errorf("execution ID %s already exists in store", executionID)
	}

	now := s.clock.Now()
	execution := &WorkflowExecution{
		Steps:       steps,
		WorkflowID:  workflowID,
		ExecutionID: executionID,
		Status:      status,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	s.idToExecution[execution.ExecutionID] = execution

	return execution.DeepCopy(), nil
}

// UpsertStep updates a step for the given executionID
func (s *InMemoryStore) UpsertStep(ctx context.Context, step *WorkflowExecutionStep) (WorkflowExecution, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	execution, ok := s.idToExecution[step.ExecutionID]
	if !ok {
		return WorkflowExecution{}, fmt.Errorf("could not find execution %s", step.ExecutionID)
	}

	now := s.clock.Now()
	execution.UpdatedAt = &now

	execution.Steps[step.Ref] = step
	return execution.DeepCopy(), nil
}

// FinishExecution marks the execution as finished with the given status
func (s *InMemoryStore) FinishExecution(ctx context.Context, executionID string, status string) (WorkflowExecution, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	execution, ok := s.idToExecution[executionID]
	if !ok {
		return WorkflowExecution{}, fmt.Errorf("could not find execution %s", executionID)
	}

	if !isCompletedStatus(status) {
		return WorkflowExecution{}, fmt.Errorf("invalid status for a finished execution %s", status)
	}

	now := s.clock.Now()
	execution.UpdatedAt = &now
	execution.Status = status
	execution.FinishedAt = &now

	return execution.DeepCopy(), nil
}

func isCompletedStatus(status string) bool {
	switch status {
	case StatusCompleted, StatusErrored, StatusTimeout, StatusCompletedEarlyExit:
		return true
	}
	return false
}

// Get gets the state for the given executionID
func (s *InMemoryStore) Get(ctx context.Context, executionID string) (WorkflowExecution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	execution, ok := s.idToExecution[executionID]
	if !ok {
		return WorkflowExecution{}, fmt.Errorf("could not find execution %s", executionID)
	}

	return execution.DeepCopy(), nil
}

func (s *InMemoryStore) Start(context.Context) error {
	return s.StartOnce("InMemoryStore", func() error {
		s.shutdownWaitGroup.Add(1)
		go s.pruneExpiredExecutionEntries()
		return nil
	})
}

func (s *InMemoryStore) Close() error {
	return s.StopOnce("InMemoryStore", func() error {
		close(s.chStop)
		s.shutdownWaitGroup.Wait()
		return nil
	})
}

func (s *InMemoryStore) Ready() error {
	return nil
}

func (s *InMemoryStore) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *InMemoryStore) Name() string {
	return "WorkflowStore"
}

func (s *InMemoryStore) pruneExpiredExecutionEntries() {
	defer s.shutdownWaitGroup.Done()
	ticker := s.clock.NewTicker(s.pruneInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.chStop:
			return
		case <-ticker.Chan():
			expirationTime := s.clock.Now().Add(-s.maximumExecutionAge)
			s.mu.Lock()
			for id, state := range s.idToExecution {
				if isCompletedStatus(state.Status) {
					delete(s.idToExecution, id)
				}
			}

			// Prune non-terminated executions that are older than the maximum expiration time
			// This shouldn't be necessary - erring on the side of caution for now as this pruning logic
			// existed in the old store.
			var prunedNonTerminatedExecutionIDs []string
			for id, state := range s.idToExecution {
				if state.UpdatedAt.Before(expirationTime) {
					delete(s.idToExecution, id)
					prunedNonTerminatedExecutionIDs = append(prunedNonTerminatedExecutionIDs, id)
				}
			}
			s.mu.Unlock()
			if len(prunedNonTerminatedExecutionIDs) > 0 {
				s.lggr.Warnw("Found and pruned non completed workflow executions older than the maximum execution age",
					"maximumExecutionAge", s.maximumExecutionAge, "pruned execution ids", prunedNonTerminatedExecutionIDs)
			}
		}
	}
}
