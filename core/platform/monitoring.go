package platform

import (
	"iter"
	"slices"
)

// Observability keys
const (
	KeyCapabilityID        = "capabilityID"
	KeyTriggerID           = "triggerID"
	KeyWorkflowID          = "workflowID"
	KeyWorkflowExecutionID = "workflowExecutionID"
	KeyWorkflowName        = "workflowName"
	KeyWorkflowOwner       = "workflowOwner"
	KeyStepID              = "stepID"
	KeyStepRef             = "stepRef"
)

func KeysSorted() iter.Seq[string] {
	return slices.Values([]string{KeyStepRef, KeyStepID, KeyTriggerID, KeyCapabilityID, KeyWorkflowExecutionID, KeyWorkflowID})
}
