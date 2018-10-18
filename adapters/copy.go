package adapters

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Copy obj keys refers to which value to copy inside `data`,
// each obj value refers to where to copy the value to inside `data`
type Copy struct {
	CopyPath []string `json:"copyPath"`
	Path     JSONPath `json:"path"`
}

// Perform returns the copied values from the desired mapping within the `data` JSON object
func (c *Copy) Perform(input models.RunResult, store *store.Store) models.RunResult {
	var jp JSONParse
	if len(c.CopyPath) > 0 {
		jp = JSONParse{Path: c.CopyPath}
	} else {
		jp = JSONParse{Path: c.Path}
	}

	data, err := input.Data.Add("value", input.Data.String())
	if err != nil {
		return input.WithError(err)
	}
	input.Data = data

	return jp.Perform(input, store)
}
