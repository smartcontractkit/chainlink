package devobservability

import (
	"context"
	"sync"
	"time"
)

// Store is the interface for dev observability storage
type Store interface {
	UpdateExecutionStatus(ctx context.Context, status string)

	GetExecution(executionID string) *ExecutionData
	GetExecutions(workflowID string, statusFilter string) []ExecutionSummary
	GetWorkflows() []string
	GetEvents(workflowID, executionID string) ([]EventEntry, error)

	Clear()
	Stats() map[string]interface{}
	GetOrphanEvents(limit int) []EventEntry
	GetWorkflowEvents(workflowID string, limit int) []EventEntry
}

type ExecutionData struct {
	WorkflowID  string
	ExecutionID string
	Status      string
	StartTime   time.Time
	EndTime     *time.Time

	Events []EventRecord

	mu sync.RWMutex
}

type ExecutionSummary struct {
	ExecutionID string    `json:"executionId"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"startTime"`
}

type EventRecord struct {
	Timestamp time.Time
	EventType string
	Data      []byte
}

type EventEntry struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   []byte    `json:"message"` // protobuf bytes
}

var globalStore Store

// GetStore returns the global store instance
func GetStore() Store {
	return globalStore
}
