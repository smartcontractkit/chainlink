package store

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/exec"
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

func (w WorkflowExecution) ResultForStep(s string) (*exec.Result, bool) {
	step, ok := w.Steps[s]
	if !ok {
		return &exec.Result{}, false
	}

	return &exec.Result{
		Inputs:  step.Inputs,
		Outputs: step.Outputs.Value,
		Error:   step.Outputs.Err,
	}, true
}

func (w WorkflowExecution) DeepCopy() WorkflowExecution {
	steps := map[string]*WorkflowExecutionStep{}
	for ref, step := range w.Steps {
		var mval *values.Map
		if step.Inputs != nil {
			mval = step.Inputs.CopyMap()
		}

		copiedov := values.Copy(step.Outputs.Value)

		newState := &WorkflowExecutionStep{
			ExecutionID: step.ExecutionID,
			Ref:         step.Ref,
			Status:      step.Status,

			Outputs: StepOutput{
				Err:   step.Outputs.Err,
				Value: copiedov,
			},

			Inputs:    mval,
			UpdatedAt: step.UpdatedAt,
		}

		steps[ref] = newState
	}
	return WorkflowExecution{
		ExecutionID: w.ExecutionID,
		WorkflowID:  w.WorkflowID,
		Status:      w.Status,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		FinishedAt:  w.FinishedAt,
		Steps:       steps,
	}
}

var _ exec.Results = WorkflowExecution{}
