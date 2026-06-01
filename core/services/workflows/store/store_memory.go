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

	// defaultRecycleMapForGCInterval is the default interval between rebuilding the execution
	// map so its old bucket storage becomes eligible for GC, even when no entries are deleted.
	defaultRecycleMapForGCInterval = 1 * time.Hour

	// maximumExecutionAge is the default maximum age of an execution before it is considered expired and eligible for pruning
	// regardless of its status
	maximumExecutionAge = 24 * time.Hour
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

	// recycleMapForGCInterval is the interval between rebuilding the execution map so old
	// bucket storage becomes eligible for GC, independent of entry deletion.
	recycleMapForGCInterval time.Duration
}

func NewInMemoryStore(lggr logger.Logger, clock clockwork.Clock) *InMemoryStore {
	return NewInMemoryStoreWithPruneConfiguration(lggr, clock, defaultPruneInterval, maximumExecutionAge,
		defaultRecycleMapForGCInterval)
}

func NewInMemoryStoreWithPruneConfiguration(lggr logger.Logger, clock clockwork.Clock, pruneFrequency time.Duration,
	maximumExecutionAge time.Duration, recycleMapForGCInterval time.Duration) *InMemoryStore {
	return &InMemoryStore{lggr: lggr, idToExecution: map[string]*WorkflowExecution{}, clock: clock, chStop: make(chan struct{}),
		pruneInterval: pruneFrequency, maximumExecutionAge: maximumExecutionAge,
		recycleMapForGCInterval: recycleMapForGCInterval}
}

// Add adds a new execution state under the given executionID
func (s *InMemoryStore) Add(ctx context.Context, steps map[string]*WorkflowExecutionStep,
	executionID string, workflowID string, status string) (WorkflowExecution, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.idToExecution[executionID]
	if ok {
		return WorkflowExecution{}, fmt.Errorf("execution ID %s: %w", executionID, ErrDuplicateExecution)
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

func (s *InMemoryStore) DeleteByWorkflowID(_ context.Context, workflowID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Build a new map instead of deleting in-place so the old backing array
	// becomes eligible for GC.  Go maps never shrink their bucket storage
	// after delete, which can retain significant memory after bulk removals.
	newMap := make(map[string]*WorkflowExecution, len(s.idToExecution))
	for id, state := range s.idToExecution {
		if state.WorkflowID != workflowID {
			newMap[id] = state
		}
	}
	s.idToExecution = newMap
	return nil
}

// recycleMapLocked copies all entries into a fresh map and swaps it in.
// Caller must hold s.mu.
func (s *InMemoryStore) recycleMapLocked() {
	if len(s.idToExecution) == 0 {
		return
	}
	newMap := make(map[string]*WorkflowExecution, len(s.idToExecution))
	for id, state := range s.idToExecution {
		newMap[id] = state
	}
	s.idToExecution = newMap
}

func (s *InMemoryStore) pruneExpiredExecutionEntries() {
	defer s.shutdownWaitGroup.Done()
	pruneTicker := s.clock.NewTicker(s.pruneInterval)
	defer pruneTicker.Stop()
	recycleTicker := s.clock.NewTicker(s.recycleMapForGCInterval)
	defer recycleTicker.Stop()
	for {
		select {
		case <-s.chStop:
			return
		case <-pruneTicker.Chan():
			s.pruneCompletedAndExpiredExecutions()
		case <-recycleTicker.Chan():
			s.recycleExecutionMapForGC()
		}
	}
}

func (s *InMemoryStore) pruneCompletedAndExpiredExecutions() {
	expirationTime := s.clock.Now().Add(-s.maximumExecutionAge)
	s.mu.Lock()
	newMap := make(map[string]*WorkflowExecution, len(s.idToExecution))
	prunedCompletedCount := 0

	// Prune non-terminated executions that are older than the maximum expiration time
	// This shouldn't be necessary - erring on the side of caution for now as this pruning logic
	// existed in the old store.
	var prunedNonTerminatedExecutionIDs []string
	for id, state := range s.idToExecution {
		if isCompletedStatus(state.Status) {
			prunedCompletedCount++
			continue
		}
		if state.UpdatedAt.Before(expirationTime) {
			prunedNonTerminatedExecutionIDs = append(prunedNonTerminatedExecutionIDs, id)
			continue
		}
		newMap[id] = state
	}
	prunedAny := prunedCompletedCount > 0 || len(prunedNonTerminatedExecutionIDs) > 0
	if prunedAny {
		// Rebuild the map to let old buckets become GC-eligible after churn-heavy periods.
		s.idToExecution = newMap
	}
	remainingExecutions := len(s.idToExecution)
	s.mu.Unlock()
	if prunedAny {
		s.lggr.Infow("Pruned workflow execution entries",
			"prunedCompletedCount", prunedCompletedCount,
			"prunedNonCompletedCount", len(prunedNonTerminatedExecutionIDs),
			"remainingExecutions", remainingExecutions)
	}
	if len(prunedNonTerminatedExecutionIDs) > 0 {
		s.lggr.Warnw("Found and pruned non completed workflow executions older than the maximum execution age",
			"maximumExecutionAge", s.maximumExecutionAge, "pruned execution ids", prunedNonTerminatedExecutionIDs)
	}
}

func (s *InMemoryStore) recycleExecutionMapForGC() {
	s.mu.Lock()
	remainingExecutions := len(s.idToExecution)
	if remainingExecutions > 0 {
		s.recycleMapLocked()
	}
	s.mu.Unlock()
	if remainingExecutions > 0 {
		s.lggr.Infow("Recycled workflow execution map for GC",
			"remainingExecutions", remainingExecutions)
	}
}
