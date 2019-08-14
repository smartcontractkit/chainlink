package adapters

import (
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Copy obj keys refers to which value to copy inside `data`,
// each obj value refers to where to copy the value to inside `data`
type Copy struct {
	CopyPath JSONPath `json:"copyPath"`
}

// Perform returns the copied values from the desired mapping within the `data` JSON object
func (c *Copy) Perform(input models.JSON, result models.RunResult, store *store.Store) models.RunResult {
	jp := JSONParse{Path: c.CopyPath}

	data, err := input.Add("result", input.String())
	if err != nil {
		result.SetError(err)
		return result
	}

	return jp.Perform(data, result, store)
}
