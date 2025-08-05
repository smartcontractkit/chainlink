package environment

import (
	"encoding/json"
	"os"
)

// WriteJSONFile marshals data into pretty JSON and writes it at path.
func WriteJSONFile(path string, data any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0600)
}
