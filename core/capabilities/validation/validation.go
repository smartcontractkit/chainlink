package validation

import (
	"encoding/hex"
	"errors"
	"unicode"
)

const (
	validWorkflowIDLen = 64
	maxIDLen           = 128
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

// Trigger event IDs and message IDs can only contain printable characters and must be non-empty
func IsValidID(id string) bool {
	if len(id) == 0 || len(id) > maxIDLen {
		return false
	}
	for i := 0; i < len(id); i++ {
		if !unicode.IsPrint(rune(id[i])) {
			return false
		}
	}
	return true
}
