package store

import (
	"context"
)

type Store interface {
	Add(ctx context.Context, steps map[string]*WorkflowExecutionStep,
		executionID string, workflowID string, status string) (WorkflowExecution, error)
	UpsertStep(ctx context.Context, step *WorkflowExecutionStep) (WorkflowExecution, error)
	FinishExecution(ctx context.Context, executionID string, status string) (WorkflowExecution, error)
	Get(ctx context.Context, executionID string) (WorkflowExecution, error)
}

var _ Store = (*InMemoryStore)(nil)
