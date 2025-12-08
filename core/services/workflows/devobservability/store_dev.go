//go:build dev

package devobservability

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const maxExecutions = 1000

type devStore struct {
	cache          *lru.Cache[string, *ExecutionData]
	workflowIndex  map[string][]string
	workflowEvents map[string][]EventEntry // Events with workflowID but no executionID
	orphanEvents   []EventEntry
	indexMu        sync.RWMutex
}

func init() {
	fmt.Println("[DevObservability] Initializing dev store...")
	cache, err := lru.NewWithEvict(maxExecutions, func(executionID string, data *ExecutionData) {
		if globalStore != nil {
			globalStore.(*devStore).removeFromIndex(data.WorkflowID, executionID)
		}
	})
	if err != nil {
		panic("failed to create LRU cache: " + err.Error())
	}

	store := &devStore{
		cache:          cache,
		workflowIndex:  make(map[string][]string),
		workflowEvents: make(map[string][]EventEntry),
		orphanEvents:   make([]EventEntry, 0),
	}

	globalStore = store
	fmt.Println("[DevObservability] Dev store initialized successfully")
}

func (s *devStore) getOrCreateExecution(workflowID, executionID string) *ExecutionData {
	if data, ok := s.cache.Get(executionID); ok {
		return data
	}

	fmt.Printf("[DevObservability] getOrCreateExecution: workflowID=%s, executionID=%s\n", workflowID, executionID)

	data := &ExecutionData{
		WorkflowID:  workflowID,
		ExecutionID: executionID,
		Status:      store.StatusStarted,
		StartTime:   time.Now(),
		Events:      make([]EventRecord, 0),
	}

	s.cache.Add(executionID, data)
	s.addToIndex(workflowID, executionID)

	return data
}

func (s *devStore) addToIndex(workflowID, executionID string) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	executions := s.workflowIndex[workflowID]
	for _, id := range executions {
		if id == executionID {
			return
		}
	}
	s.workflowIndex[workflowID] = append(executions, executionID)
}

func (s *devStore) removeFromIndex(workflowID, executionID string) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	executions := s.workflowIndex[workflowID]
	for i, id := range executions {
		if id == executionID {
			executions[i] = executions[len(executions)-1]
			s.workflowIndex[workflowID] = executions[:len(executions)-1]
			break
		}
	}

	if len(s.workflowIndex[workflowID]) == 0 {
		delete(s.workflowIndex, workflowID)
	}
}

func (s *devStore) storeRawEvent(ctx context.Context, workflowID, executionID, eventType string, payload []byte) {
	data := s.getOrCreateExecution(workflowID, executionID)

	data.mu.Lock()
	defer data.mu.Unlock()

	data.Events = append(data.Events, EventRecord{
		Timestamp: time.Now(),
		EventType: eventType,
		Data:      payload,
	})

	fmt.Printf("[DevObservability] Event stored! workflow=%s, execution=%s, type=%s, total_events=%d\n", workflowID, executionID, eventType, len(data.Events))
}

func (s *devStore) storeOrphanEvent(ctx context.Context, eventType string, payload []byte) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	s.orphanEvents = append(s.orphanEvents, EventEntry{
		Type:      eventType,
		Timestamp: time.Now(),
		Message:   payload,
	})

	// Keep only last 1000 orphan events
	if len(s.orphanEvents) > 1000 {
		s.orphanEvents = s.orphanEvents[len(s.orphanEvents)-1000:]
	}

	fmt.Printf("[DevObservability] Orphan event stored! type=%s, total_orphan_events=%d\n", eventType, len(s.orphanEvents))
}

func (s *devStore) storeWorkflowEvent(ctx context.Context, workflowID, eventType string, payload []byte) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	fmt.Printf("[DevObservability] storeWorkflowEvent called with workflowID=%s, eventType=%s\n", workflowID, eventType)

	entry := EventEntry{
		Type:      eventType,
		Timestamp: time.Now(),
		Message:   payload,
	}

	s.workflowEvents[workflowID] = append(s.workflowEvents[workflowID], entry)

	// Keep only last 100 events per workflow
	if len(s.workflowEvents[workflowID]) > 100 {
		s.workflowEvents[workflowID] = s.workflowEvents[workflowID][len(s.workflowEvents[workflowID])-100:]
	}

	fmt.Printf("[DevObservability] Workflow-level event stored! workflow=%s, type=%s, total_workflow_events=%d\n", workflowID, eventType, len(s.workflowEvents[workflowID]))
}

func (s *devStore) UpdateExecutionStatus(ctx context.Context, status string) {
	workflowID, executionID, ok := GetExecutionContext(ctx)
	if !ok {
		return
	}

	data := s.getOrCreateExecution(workflowID, executionID)
	data.mu.Lock()
	defer data.mu.Unlock()

	data.Status = status
	if status == store.StatusCompleted || status == store.StatusErrored || status == store.StatusTimeout || status == store.StatusCompletedEarlyExit {
		now := time.Now()
		data.EndTime = &now
	}
}

func (s *devStore) GetExecution(executionID string) *ExecutionData {
	data, _ := s.cache.Get(executionID)
	return data
}

func (s *devStore) GetExecutions(workflowID string, statusFilter string) []ExecutionSummary {
	s.indexMu.RLock()
	executionIDs := s.workflowIndex[workflowID]
	s.indexMu.RUnlock()

	summaries := make([]ExecutionSummary, 0, len(executionIDs))
	for _, execID := range executionIDs {
		if data, ok := s.cache.Get(execID); ok {
			data.mu.RLock()
			// Apply status filter if specified
			if statusFilter == "" || data.Status == statusFilter {
				summaries = append(summaries, ExecutionSummary{
					ExecutionID: data.ExecutionID,
					Status:      data.Status,
					StartTime:   data.StartTime,
				})
			}
			data.mu.RUnlock()
		}
	}

	return summaries
}

func (s *devStore) GetWorkflows() []string {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()

	// Use a map to deduplicate workflow IDs from both sources
	workflowMap := make(map[string]bool)

	// Add workflows that have executions
	for workflowID := range s.workflowIndex {
		workflowMap[workflowID] = true
	}

	// Add workflows that have workflow-level events
	for workflowID := range s.workflowEvents {
		workflowMap[workflowID] = true
	}

	workflowIDs := make([]string, 0, len(workflowMap))
	for workflowID := range workflowMap {
		workflowIDs = append(workflowIDs, workflowID)
	}

	return workflowIDs
}

func (s *devStore) GetEvents(workflowID, executionID string) ([]EventEntry, error) {
	data, ok := s.cache.Get(executionID)
	if !ok {
		return nil, errors.New("execution not found")
	}

	data.mu.RLock()
	defer data.mu.RUnlock()

	entries := make([]EventEntry, len(data.Events))
	for i, evt := range data.Events {
		entries[i] = EventEntry{
			Type:      evt.EventType,
			Timestamp: evt.Timestamp,
			Message:   evt.Data,
		}
	}

	return entries, nil
}

func (s *devStore) GetOrphanEvents(limit int) []EventEntry {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()

	if limit <= 0 || limit > len(s.orphanEvents) {
		limit = len(s.orphanEvents)
	}

	// Return most recent events
	start := len(s.orphanEvents) - limit
	if start < 0 {
		start = 0
	}

	result := make([]EventEntry, limit)
	copy(result, s.orphanEvents[start:])

	return result
}

func (s *devStore) GetWorkflowEvents(workflowID string, limit int) []EventEntry {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()

	events, ok := s.workflowEvents[workflowID]
	if !ok {
		return []EventEntry{}
	}

	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}

	// Return most recent events
	start := len(events) - limit
	if start < 0 {
		start = 0
	}

	result := make([]EventEntry, limit)
	copy(result, events[start:])

	return result
}

func (s *devStore) Clear() {
	s.cache.Purge()
	s.indexMu.Lock()
	s.workflowIndex = make(map[string][]string)
	s.workflowEvents = make(map[string][]EventEntry)
	s.orphanEvents = make([]EventEntry, 0)
	s.indexMu.Unlock()
}

func (s *devStore) Stats() map[string]interface{} {
	s.indexMu.RLock()
	workflowCount := len(s.workflowIndex)
	workflowStats := make(map[string]int)
	for wfID, execIDs := range s.workflowIndex {
		workflowStats[wfID] = len(execIDs)
	}
	s.indexMu.RUnlock()

	totalEvents := 0
	totalExecutions := s.cache.Len()

	// Count all events across all executions
	keys := s.cache.Keys()
	for _, execID := range keys {
		if data, ok := s.cache.Peek(execID); ok {
			data.mu.RLock()
			totalEvents += len(data.Events)
			data.mu.RUnlock()
		}
	}

	s.indexMu.RLock()
	orphanEventsCount := len(s.orphanEvents)
	workflowEventsCount := 0
	for _, events := range s.workflowEvents {
		workflowEventsCount += len(events)
	}
	s.indexMu.RUnlock()

	return map[string]interface{}{
		"total_executions": totalExecutions,
		"total_workflows":  workflowCount,
		"total_events":     totalEvents,
		"workflow_events":  workflowEventsCount,
		"orphan_events":    orphanEventsCount,
		"workflow_stats":   workflowStats,
		"capacity":         maxExecutions,
	}
}
