package changeset

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
)

func FeedIDsToBytes16(feedIDs []string) ([][16]byte, error) {
	dataIds := make([][16]byte, len(feedIDs))
	for i, feedID := range feedIDs {
		err := ValidateFeedID(feedID)
		if err != nil {
			return nil, err
		}
		dataIds[i], err = ConvertHexToBytes16(feedID)
		if err != nil {
			return nil, err
		}
	}

	return dataIds, nil
}

func ConvertHexToBytes16(hexStr string) ([16]byte, error) {
	if hexStr[:2] == "0x" {
		hexStr = hexStr[2:] // Remove "0x" prefix
	}
	decodedBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed to decode hex string: %w", err)
	}

	var result [16]byte
	copy(result[:], decodedBytes[:16])

	return result, nil
}

func HashedWorkflowName(name string) [10]byte {
	// Compute SHA-256 hash of the input string
	hash := sha256.Sum256([]byte(name))

	// Encode as hex to ensure UTF8
	var hashBytes = hash[:]
	resultHex := hex.EncodeToString(hashBytes)

	// Truncate to 10 bytes
	var truncated [10]byte
	copy(truncated[:], []byte(resultHex)[:10])

	return truncated
}

func LoadJSON[T any](pth string, fs fs.ReadFileFS) (T, error) {
	var dflt T
	f, err := fs.ReadFile(pth)
	if err != nil {
		return dflt, fmt.Errorf("failed to read %s: %w", pth, err)
	}
	var v T
	err = json.Unmarshal(f, &v)
	if err != nil {
		return dflt, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return v, nil
}
