package devobservability

import (
	"time"
)

// Store is the interface for dev observability storage
type Store interface {
	GetOrphanEvents(limit int, minSequence int64) []EventEntry
	GetWorkflowEvents(workflowID string, limit int, minSequence int64) []EventEntry
	GetWorkflows() []string

	Clear()
	Stats() map[string]interface{}
}

type EventRecord struct {
	Timestamp time.Time
	EventType string
	Data      []byte
}

type EventEntry struct {
	Sequence  int64     `json:"sequence"` // sequence number within workflow
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   []byte    `json:"message"` // protobuf bytes
}

var globalStore Store

// GetStore returns the global store instance
func GetStore() Store {
	return globalStore
}
