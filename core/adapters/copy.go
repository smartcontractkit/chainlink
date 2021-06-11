package adapters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Copy obj keys refers to which value to copy inside `data`,
// each obj value refers to where to copy the value to inside `data`
type Copy struct {
	CopyPath JSONPath `json:"copyPath"`
}

// TaskType returns the type of Adapter.
func (c *Copy) TaskType() models.TaskType {
	return TaskTypeCopy
}

// Perform returns the copied values from the desired mapping within the `data` JSON object
func (c *Copy) Perform(input models.RunInput, store *store.Store, keyStore *keystore.Master) models.RunOutput {
	data, err := models.JSON{}.Add("result", input.Data().String())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	jp := JSONParse{Path: c.CopyPath}
	input = input.CloneWithData(data)
	return jp.Perform(input, store, keyStore)
}
