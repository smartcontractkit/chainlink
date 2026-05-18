package store

import (
	"context"
)

type Store interface {
	Add(ctx context.Context, state *WorkflowExecution) error
	UpsertStep(ctx context.Context, step *WorkflowExecutionStep) (WorkflowExecution, error)
	UpdateStatus(ctx context.Context, executionID string, status string) error
	Get(ctx context.Context, executionID string) (WorkflowExecution, error)
	GetUnfinished(ctx context.Context, offset, limit int) ([]WorkflowExecution, error)
}

var _ Store = (*DBStore)(nil)
