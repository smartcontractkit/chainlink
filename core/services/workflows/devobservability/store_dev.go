//go:build dev

package devobservability

import (
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	maxWorkflows         = 100  // Maximum number of workflows to cache
	maxEventsPerWorkflow = 1000 // Maximum events per workflow
	maxOrphanEvents      = 1000 // Maximum orphan events to cache
)

type workflowEventData struct {
	mu           sync.RWMutex
	eventsCache  *lru.Cache[int64, EventEntry] // keyed by sequence number
	nextSequence int64
}

type devStore struct {
	cache              *lru.Cache[string, *workflowEventData]
	orphanCache        *lru.Cache[int64, EventEntry] // keyed by sequence number
	orphanMu           sync.RWMutex
	nextOrphanSequence int64
}

func init() {
	fmt.Println("[DevObservability] Initializing dev store...")
	cache, err := lru.New[string, *workflowEventData](maxWorkflows)
	if err != nil {
		panic("failed to create workflow LRU cache: " + err.Error())
	}

	orphanCache, err := lru.New[int64, EventEntry](maxOrphanEvents)
	if err != nil {
		panic("failed to create orphan LRU cache: " + err.Error())
	}

	store := &devStore{
		cache:              cache,
		orphanCache:        orphanCache,
		nextOrphanSequence: 1,
	}

	globalStore = store
	fmt.Println("[DevObservability] Dev store initialized successfully")
}

func (s *devStore) getOrCreateWorkflowData(workflowID string) *workflowEventData {
	if data, ok := s.cache.Get(workflowID); ok {
		return data
	}

	eventsCache, err := lru.New[int64, EventEntry](maxEventsPerWorkflow)
	if err != nil {
		panic("failed to create events LRU cache: " + err.Error())
	}

	data := &workflowEventData{
		eventsCache:  eventsCache,
		nextSequence: 1,
	}

	s.cache.Add(workflowID, data)
	fmt.Printf("[DevObservability] Created workflow data for: %s\n", workflowID)

	return data
}

func (s *devStore) storeWorkflowEvent(workflowID, eventType string, payload []byte) {
	data := s.getOrCreateWorkflowData(workflowID)

	data.mu.Lock()
	defer data.mu.Unlock()

	event := EventEntry{
		Sequence:  data.nextSequence,
		Type:      eventType,
		Timestamp: time.Now(),
		Message:   payload,
	}

	// Add to LRU cache - it will automatically evict oldest when at capacity
	data.eventsCache.Add(data.nextSequence, event)
	data.nextSequence++

	fmt.Printf("[DevObservability] Event stored! workflow=%s, type=%s, seq=%d, total_events=%d\n",
		workflowID, eventType, event.Sequence, data.eventsCache.Len())
}

func (s *devStore) storeOrphanEvent(eventType string, payload []byte) {
	s.orphanMu.Lock()
	defer s.orphanMu.Unlock()

	event := EventEntry{
		Sequence:  s.nextOrphanSequence,
		Type:      eventType,
		Timestamp: time.Now(),
		Message:   payload,
	}

	// Use sequence number as key
	// LRU will automatically evict oldest when capacity is reached
	s.orphanCache.Add(s.nextOrphanSequence, event)
	s.nextOrphanSequence++

	fmt.Printf("[DevObservability] Orphan event stored! type=%s, seq=%d, total_orphan_events=%d\n",
		eventType, event.Sequence, s.orphanCache.Len())
}

func (s *devStore) GetWorkflows() []string {
	keys := s.cache.Keys()
	workflowIDs := make([]string, len(keys))
	copy(workflowIDs, keys)
	return workflowIDs
}

func (s *devStore) GetOrphanEvents(limit int, minSequence int64) []EventEntry {
	s.orphanMu.RLock()
	defer s.orphanMu.RUnlock()

	// Get all sequence keys from cache
	keys := s.orphanCache.Keys()
	if len(keys) == 0 {
		return []EventEntry{}
	}

	// Get all events and filter by sequence
	allEvents := make([]EventEntry, 0, len(keys))
	for _, seq := range keys {
		if evt, ok := s.orphanCache.Peek(seq); ok {
			// Filter by minimum sequence if specified
			if minSequence > 0 && evt.Sequence < minSequence {
				continue
			}
			allEvents = append(allEvents, evt)
		}
	}

	// Sort by sequence to ensure chronological order
	for i := 0; i < len(allEvents)-1; i++ {
		for j := i + 1; j < len(allEvents); j++ {
			if allEvents[i].Sequence > allEvents[j].Sequence {
				allEvents[i], allEvents[j] = allEvents[j], allEvents[i]
			}
		}
	}

	// Return most recent N events
	if limit <= 0 || limit > len(allEvents) {
		limit = len(allEvents)
	}

	return allEvents[max(len(allEvents)-limit, 0):]
}

func (s *devStore) GetWorkflowEvents(workflowID string, limit int, minSequence int64) []EventEntry {
	data, ok := s.cache.Get(workflowID)
	if !ok {
		return []EventEntry{}
	}

	data.mu.RLock()
	defer data.mu.RUnlock()

	// Get all sequence keys from cache
	keys := data.eventsCache.Keys()
	if len(keys) == 0 {
		return []EventEntry{}
	}

	// Get all events and sort by sequence
	allEvents := make([]EventEntry, 0, len(keys))
	for _, seq := range keys {
		if evt, ok := data.eventsCache.Peek(seq); ok {
			// Filter by minimum sequence if specified
			if minSequence > 0 && evt.Sequence < minSequence {
				continue
			}
			allEvents = append(allEvents, evt)
		}
	}

	// Sort by sequence to ensure chronological order
	for i := 0; i < len(allEvents)-1; i++ {
		for j := i + 1; j < len(allEvents); j++ {
			if allEvents[i].Sequence > allEvents[j].Sequence {
				allEvents[i], allEvents[j] = allEvents[j], allEvents[i]
			}
		}
	}

	// Return most recent N events
	if limit <= 0 || limit > len(allEvents) {
		limit = len(allEvents)
	}

	return allEvents[max(len(allEvents)-limit, 0):]
}

func (s *devStore) Clear() {
	s.cache.Purge()
	s.orphanMu.Lock()
	s.orphanCache.Purge()
	s.nextOrphanSequence = 1
	s.orphanMu.Unlock()
}

func (s *devStore) Stats() map[string]interface{} {
	totalEvents := 0
	workflowStats := make(map[string]int)

	// Count all events across all workflows
	keys := s.cache.Keys()
	for _, workflowID := range keys {
		if data, ok := s.cache.Peek(workflowID); ok {
			data.mu.RLock()
			eventCount := data.eventsCache.Len()
			workflowStats[workflowID] = eventCount
			totalEvents += eventCount
			data.mu.RUnlock()
		}
	}

	s.orphanMu.RLock()
	orphanEventsCount := s.orphanCache.Len()
	s.orphanMu.RUnlock()

	return map[string]interface{}{
		"total_workflows":         len(keys),
		"total_events":            totalEvents,
		"orphan_events":           orphanEventsCount,
		"workflow_stats":          workflowStats,
		"max_workflows":           maxWorkflows,
		"max_events_per_workflow": maxEventsPerWorkflow,
		"max_orphan_events":       maxOrphanEvents,
	}
}
