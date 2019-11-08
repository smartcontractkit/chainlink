package adapters

import (
	"chainlink/core/store"
	"chainlink/core/store/models"
)

// Copy obj keys refers to which value to copy inside `data`,
// each obj value refers to where to copy the value to inside `data`
type Copy struct {
	CopyPath JSONPath `json:"copyPath"`
}

// Perform returns the copied values from the desired mapping within the `data` JSON object
func (c *Copy) Perform(input models.RunInput, store *store.Store) models.RunOutput {
	data, err := models.JSON{}.Add("result", input.Data().String())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	jp := JSONParse{Path: c.CopyPath}
	input = *models.NewRunInput(input.JobRunID(), data, input.Status())
	return jp.Perform(input, store)
}
