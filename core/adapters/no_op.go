package adapters

import (
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// NoOp adapter type holds no fields
type NoOp struct{}

// TaskType returns the type of Adapter.
func (noa *NoOp) TaskType() models.TaskType {
	return TaskTypeNoOp
}

// Perform returns the input
func (noa *NoOp) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	val := input.Result().Value()
	return models.NewRunOutputCompleteWithResult(val)
}

// NoOpPend adapter type holds no fields
type NoOpPend struct{}

// TaskType returns the type of Adapter.
func (noa *NoOpPend) TaskType() models.TaskType {
	return TaskTypeNoOpPend
}

// Perform on this adapter type returns an empty RunResult with an
// added field for the status to indicate the task is Pending.
func (noa *NoOpPend) Perform(_ models.RunInput, _ *store.Store) models.RunOutput {
	return models.NewRunOutputPendingConfirmationsWithData(models.JSON{})
}
