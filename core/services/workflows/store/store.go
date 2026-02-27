package store

import (
	"context"
	"errors"
)

// ErrDuplicateExecution is returned by Add when the execution ID already exists.
var ErrDuplicateExecution = errors.New("duplicate execution")

type Store interface {
	Add(ctx context.Context, steps map[string]*WorkflowExecutionStep,
		executionID string, workflowID string, status string) (WorkflowExecution, error)
	UpsertStep(ctx context.Context, step *WorkflowExecutionStep) (WorkflowExecution, error)
	FinishExecution(ctx context.Context, executionID string, status string) (WorkflowExecution, error)
	Get(ctx context.Context, executionID string) (WorkflowExecution, error)
	DeleteByWorkflowID(ctx context.Context, workflowID string) error
}

var _ Store = (*InMemoryStore)(nil)
