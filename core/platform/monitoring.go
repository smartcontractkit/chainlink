package platform

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

var OrderedLabelKeys = []string{KeyStepRef, KeyStepID, KeyTriggerID, KeyCapabilityID, KeyWorkflowExecutionID, KeyWorkflowID}
