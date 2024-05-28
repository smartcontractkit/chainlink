package store

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	StatusStarted   = "started"
	StatusErrored   = "errored"
	StatusTimeout   = "timeout"
	StatusCompleted = "completed"
)

type StepOutput struct {
	Err   error
	Value values.Value
}

type WorkflowExecutionStep struct {
	ExecutionID string
	Ref         string
	Status      string

	Inputs  *values.Map
	Outputs *StepOutput

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
