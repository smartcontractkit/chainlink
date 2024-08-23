package store

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// Note: any update to the enum below should be reflected in
// ValidStatuses and the database enum `workflow_status`.
const (
	StatusStarted            = "started"
	StatusErrored            = "errored"
	StatusTimeout            = "timeout"
	StatusCompleted          = "completed"
	StatusCompletedEarlyExit = "completed_early_exit"
)

var ValidStatuses = map[string]bool{
	StatusStarted:            true,
	StatusErrored:            true,
	StatusTimeout:            true,
	StatusCompleted:          true,
	StatusCompletedEarlyExit: true,
}

type StepOutput struct {
	Err   error
	Value values.Value
}

type WorkflowExecutionStep struct {
	ExecutionID string
	Ref         string
	Status      string

	Inputs  *values.Map
	Outputs StepOutput

	UpdatedAt *time.Time
}

type WorkflowExecution struct {
	Steps       map[string]*WorkflowExecutionStep
	ExecutionID string
	WorkflowID  string

	Status     string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	FinishedAt *time.Time
}
