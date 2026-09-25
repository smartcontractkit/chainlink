package triggers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
)

const _prefix = "trigger_reg"

// RegistrationID constructs a trigger registration ID from a workflow ID and trigger index.
func RegistrationID(workflowID string, triggerIndex int) string {
	return fmt.Sprintf("%s_%s_%d", _prefix, workflowID, triggerIndex)
}

// ParseWorkflowID extracts the workflow ID from a registration ID produced by RegistrationID.
func ParseWorkflowID(registrationID string) (string, error) {
	rest, ok := strings.CutPrefix(registrationID, _prefix+"_")
	if !ok {
		return "", fmt.Errorf("invalid trigger registration ID %q: missing prefix %q", registrationID, _prefix)
	}

	idx := strings.LastIndex(rest, "_")
	if idx == -1 {
		return "", fmt.Errorf("invalid trigger registration ID %q: missing trigger index", registrationID)
	}

	workflowID, triggerIndexStr := rest[:idx], rest[idx+1:]
	if _, err := strconv.Atoi(triggerIndexStr); err != nil {
		return "", fmt.Errorf("invalid trigger registration ID %q: invalid trigger index %q: %w", registrationID, triggerIndexStr, err)
	}

	if err := validation.ValidateWorkflowOrExecutionID(workflowID); err != nil {
		return "", fmt.Errorf("invalid trigger registration ID %q: invalid workflow ID %q: %w", registrationID, workflowID, err)
	}

	return workflowID, nil
}
