package adapters

import (
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// NoOp adapter type holds no fields
type NoOp struct{}

// Perform returns the empty RunResult
func (noa *NoOp) Perform(input models.JSON, result models.RunResult, _ *store.Store) models.RunResult {
	result.Status = models.RunStatusCompleted
	return result
}

// NoOpPend adapter type holds no fields
type NoOpPend struct{}

// Perform on this adapter type returns an empty RunResult with an
// added field for the status to indicate the task is Pending.
func (noa *NoOpPend) Perform(_ models.JSON, result models.RunResult, _ *store.Store) models.RunResult {
	result.MarkPendingConfirmations()
	return result
}
