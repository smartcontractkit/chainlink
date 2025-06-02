package platform

import (
	"slices"

	"iter"
)

// Observability keys
const (
	KeyCapabilityID        = "capabilityID"
	KeyTriggerID           = "triggerID"
	KeyWorkflowID          = "workflowID"
	KeyWorkflowExecutionID = "workflowExecutionID"
	KeyWorkflowName        = "workflowName"
	KeyWorkflowVersion     = "workflowVersion"
	KeyWorkflowOwner       = "workflowOwner"
	KeyStepID              = "stepID"
	KeyStepRef             = "stepRef"
	KeyDonID               = "DonID"
	KeyDonF                = "F"
	KeyDonN                = "N"
	KeyDonQ                = "Q"
	KeyP2PID               = "p2pID"
	ValueWorkflowVersion   = "1.0.0"
)

func LabelKeysSorted() iter.Seq[string] {
	return slices.Values([]string{
		KeyStepRef,
		KeyStepID,
		KeyTriggerID,
		KeyCapabilityID,
		KeyWorkflowExecutionID,
		KeyWorkflowID,
	})
}
