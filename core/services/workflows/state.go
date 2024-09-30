package workflows

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// copyState returns a deep copy of the input executionState
func copyState(es store.WorkflowExecution) store.WorkflowExecution {
	steps := map[string]*store.WorkflowExecutionStep{}
	for ref, step := range es.Steps {
		var mval *values.Map
		if step.Inputs != nil {
			mval = step.Inputs.CopyMap()
		}

		copiedov := values.Copy(step.Outputs.Value)

		newState := &store.WorkflowExecutionStep{
			ExecutionID: step.ExecutionID,
			Ref:         step.Ref,
			Status:      step.Status,

			Outputs: store.StepOutput{
				Err:   step.Outputs.Err,
				Value: copiedov,
			},

			Inputs: mval,
		}

		steps[ref] = newState
	}
	return store.WorkflowExecution{
		ExecutionID: es.ExecutionID,
		WorkflowID:  es.WorkflowID,
		Status:      es.Status,
		Steps:       steps,
	}
}
