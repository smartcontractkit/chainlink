package adapters

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Map holds holds a path to the desired field in the `data` JSON object,
// made up of an array of strings.
type Map struct {
	Path []string `json:"path"`
}

// Perform returns the value from a desired mapping within the `data` JSON object,
// this is specifically used for bridge adaptor returns
//
// For reference on how the parsing is done, refer to JsonParse
func (m *Map) Perform(input models.RunResult, store *store.Store) models.RunResult {
	jp := JSONParse{Path: m.Path}

	data, err := input.Data.Add("value", input.Data.String())
	if err != nil {
		return input.WithError(err)
	}
	input.Data = data

	return jp.Perform(input, store)
}
