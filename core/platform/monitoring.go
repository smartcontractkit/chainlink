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
	KeyWorkflowTag         = "workflowTag"
	KeyWorkflowVersion     = "workflowVersion"
	KeyWorkflowOwner       = "workflowOwner"
	KeyOrganizationID      = "orgID"
	KeyStepID              = "stepID"
	KeyStepRef             = "stepRef"
	KeyDonID               = "DonID"
	KeyDonF                = "F"
	KeyDonN                = "N"
	KeyDonQ                = "Q"
	KeyP2PID               = "p2pID"
	ValueWorkflowVersion   = "1.0.0"
	ValueWorkflowVersionV2 = "2.0.0"

	// Registry and version keys
	WorkflowRegistryAddress     = "workflowRegistryAddress"
	WorkflowRegistryVersion     = "workflowRegistryVersion"
	WorkflowRegistryChain       = "workflowRegistryChain"
	EngineVersion               = "engineVersion"
	CapabilitiesRegistryVersion = "capabilitiesRegistryVersion"
	DonVersion                  = "donVersion"
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
