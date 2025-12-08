package devobservability

import (
	"context"
	"fmt"
)

type executionContextKey struct{}

type ExecutionContext struct {
	WorkflowID  string
	ExecutionID string
}

// WithExecutionContext adds workflow execution metadata to context
func WithExecutionContext(ctx context.Context, workflowID, executionID string) context.Context {
	fmt.Printf("[DevObservability] Setting execution context: workflowID=%s, executionID=%s\n", workflowID, executionID)
	return context.WithValue(ctx, executionContextKey{}, &ExecutionContext{
		WorkflowID:  workflowID,
		ExecutionID: executionID,
	})
}

// GetExecutionContext extracts workflow execution metadata from context
func GetExecutionContext(ctx context.Context) (workflowID, executionID string, ok bool) {
	execCtx, ok := ctx.Value(executionContextKey{}).(*ExecutionContext)
	if !ok || execCtx == nil {
		return "", "", false
	}
	return execCtx.WorkflowID, execCtx.ExecutionID, true
}

// WithWorkflowID adds only workflow ID to context (execution ID remains empty)
func WithWorkflowID(ctx context.Context, workflowID string) context.Context {
	// Check if there's already an execution context
	if execCtx, ok := ctx.Value(executionContextKey{}).(*ExecutionContext); ok && execCtx != nil {
		// Preserve existing execution ID if present
		return context.WithValue(ctx, executionContextKey{}, &ExecutionContext{
			WorkflowID:  workflowID,
			ExecutionID: execCtx.ExecutionID,
		})
	}
	return context.WithValue(ctx, executionContextKey{}, &ExecutionContext{
		WorkflowID: workflowID,
	})
}

// WithExecutionID adds only execution ID to context (preserving workflow ID if present)
func WithExecutionID(ctx context.Context, executionID string) context.Context {
	// Check if there's already an execution context with workflow ID
	if execCtx, ok := ctx.Value(executionContextKey{}).(*ExecutionContext); ok && execCtx != nil {
		// Preserve existing workflow ID
		return context.WithValue(ctx, executionContextKey{}, &ExecutionContext{
			WorkflowID:  execCtx.WorkflowID,
			ExecutionID: executionID,
		})
	}
	// No existing context, create new one with just execution ID
	return context.WithValue(ctx, executionContextKey{}, &ExecutionContext{
		ExecutionID: executionID,
	})
}
