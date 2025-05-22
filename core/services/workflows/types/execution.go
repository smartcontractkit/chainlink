package types

import (
	"crypto/sha256"
	"encoding/hex"
)

// hash of (workflowID, triggerEventID)
// TODO(CAPPL-838): improve for V2
func GenerateExecutionID(workflowID, triggerEventID string) (string, error) {
	s := sha256.New()
	_, err := s.Write([]byte(workflowID))
	if err != nil {
		return "", err
	}

	_, err = s.Write([]byte(triggerEventID))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}
