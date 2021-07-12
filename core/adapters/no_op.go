package adapters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore"
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
func (noa *NoOp) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	val := input.Result().Value()
	return models.NewRunOutputCompleteWithResult(val, input.ResultCollection())
}

// NoOpPendOutgoing adapter type holds no fields
type NoOpPendOutgoing struct{}

// TaskType returns the type of Adapter.
func (noa *NoOpPendOutgoing) TaskType() models.TaskType {
	return TaskTypeNoOpPendOutgoing
}

// Perform on this adapter type returns an empty RunResult with an
// added field for the status to indicate the task is Pending.
func (noa *NoOpPendOutgoing) Perform(_ models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	return models.NewRunOutputPendingOutgoingConfirmationsWithData(models.JSON{})
}
