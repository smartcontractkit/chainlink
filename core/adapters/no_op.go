package adapters

import (
	"chainlink/core/store"
	"chainlink/core/store/models"
)

// NoOp adapter type holds no fields
type NoOp struct{}

// Perform returns the input
func (noa *NoOp) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	val := input.Result().Value()
	return models.NewRunOutputCompleteWithResult(val)
}

// NoOpPend adapter type holds no fields
type NoOpPend struct{}

// Perform on this adapter type returns an empty RunResult with an
// added field for the status to indicate the task is Pending.
func (noa *NoOpPend) Perform(_ models.RunInput, _ *store.Store) models.RunOutput {
	return models.NewRunOutputPendingConfirmations()
}
