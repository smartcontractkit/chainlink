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
func (c *Copy) Perform(input models.RunResult, store *store.Store) models.RunResult {
	jp := JSONParse{Path: c.CopyPath}

	data, err := input.Data.Add("result", input.Data.String())
	if err != nil {
		input.SetError(err)
		return input
	}
	input.Data = data

	return jp.Perform(input, store)
}
