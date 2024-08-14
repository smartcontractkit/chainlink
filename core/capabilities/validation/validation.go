package validation

import (
	"encoding/hex"
	"errors"
)

const (
	validWorkflowIDLen = 64
)

// Workflow IDs and Execution IDs are 32-byte hex-encoded strings
func ValidateWorkflowOrExecutionID(id string) error {
	if len(id) != validWorkflowIDLen {
		return errors.New("must be 32 bytes long")
	}
	_, err := hex.DecodeString(id)
	if err != nil {
		return errors.New("must be a hex-encoded string")
	}

	return nil
}
