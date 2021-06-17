package adapters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// ResultCollect receiver type
type ResultCollect struct{}

// TaskType returns the TaskTypeResultCollect adapter
func (r ResultCollect) TaskType() models.TaskType {
	return TaskTypeResultCollect
}

// Perform takes an input to run and returns the output
func (r ResultCollect) Perform(input models.RunInput, store *store.Store, _ *keystore.Master) models.RunOutput {
	updatedCollection := make([]interface{}, 0)
	for _, c := range input.ResultCollection().Array() {
		updatedCollection = append(updatedCollection, c.Value())
	}
	updatedCollection = append(updatedCollection, input.Result().Value())
	ro, err := input.Data().Add(models.ResultCollectionKey, updatedCollection)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	return models.NewRunOutputComplete(ro)
}
