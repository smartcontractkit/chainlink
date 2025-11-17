package environment

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// WriteJSONFile marshals data into pretty JSON and writes it at path.
func WriteJSONFile(path string, data any) error {
	b, jsonErr := json.MarshalIndent(data, "", "  ")
	if jsonErr != nil {
		return errors.Wrap(jsonErr, "failed to marshal data to JSON")
	}

	return os.WriteFile(path, b, 0600)
}
