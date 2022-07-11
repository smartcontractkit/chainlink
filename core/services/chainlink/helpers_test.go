package chainlink

import "encoding/json"

// FakeConfigDump is a test helper for configDump which loads dbData from a JSON string instead of a live ChainSet.
func FakeConfigDump(dbDataJSON []byte) (string, error) {
	var data dbData

	if len(dbDataJSON) > 0 {
		if err := json.Unmarshal(dbDataJSON, &data); err != nil {
			return "", err
		}
	}

	return configDump(data)
}
