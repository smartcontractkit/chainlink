package changeset

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"

	workflowUtils "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

func FeedIDsToBytes16(feedIDs []string) ([][16]byte, error) {
	dataIDs := make([][16]byte, len(feedIDs))
	for i, feedID := range feedIDs {
		err := ValidateFeedID(feedID)
		if err != nil {
			return nil, err
		}
		dataIDs[i], err = ConvertHexToBytes16(feedID)
		if err != nil {
			return nil, err
		}
	}

	return dataIDs, nil
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

// HashedWorkflowName returns first 10 bytes of the sha256(workflow_name)
func HashedWorkflowName(name string) [10]byte {
	nameHash := workflowUtils.HashTruncateName(name)
	var result [10]byte
	copy(result[:], nameHash)
	return result
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
