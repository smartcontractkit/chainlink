package adapters

import (
	"chainlink/core/store"
	"chainlink/core/store/models"
)

// NoOp adapter type holds no fields
type NoOp struct{}

// Perform returns the empty RunResult
func (noa *NoOp) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val := input.Result().Value()
	return models.RunResultComplete(val)
}

// NoOpPend adapter type holds no fields
type NoOpPend struct{}

// Perform on this adapter type returns an empty RunResult with an
// added field for the status to indicate the task is Pending.
func (noa *NoOpPend) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	var output models.RunResult
	output.MarkPendingConfirmations()
	return output
}
