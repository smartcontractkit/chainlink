//go:build !dev

package devobservability

import (
	"context"
	"errors"
)

type noopStore struct{}

func init() {
	store := &noopStore{}
	globalStore = store
}

func (s *noopStore) UpdateExecutionStatus(ctx context.Context, status string) {}

func (s *noopStore) GetExecution(executionID string) *ExecutionData { return nil }

func (s *noopStore) GetExecutions(workflowID string, statusFilter string) []ExecutionSummary {
	return nil
}

func (s *noopStore) GetWorkflows() []string { return nil }

func (s *noopStore) GetEvents(workflowID, executionID string) ([]EventEntry, error) {
	return nil, errors.New("not available in production")
}

func (s *noopStore) Clear() {}

func (s *noopStore) GetOrphanEvents(limit int) []EventEntry { return nil }

func (s *noopStore) GetWorkflowEvents(workflowID string, limit int) []EventEntry { return nil }

func (s *noopStore) Stats() map[string]interface{} { return nil }
