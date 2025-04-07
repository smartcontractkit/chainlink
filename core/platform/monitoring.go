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
	KeyWorkflowVersion     = "workflowVersion"
	KeyWorkflowOwner       = "workflowOwner"
	KeyStepID              = "stepID"
	KeyStepRef             = "stepRef"
	KeyDonID               = "DonID"
	KeyDonF                = "F"
	KeyDonN                = "N"
	KeyDonQ                = "Q"
	KeyP2PID               = "p2pID"
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
