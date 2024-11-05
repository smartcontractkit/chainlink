package platform

// Observability keys
const (
	CapabilityIDKey        = "capabilityID"
	TriggerIDKey           = "triggerID"
	WorkflowIDKey          = "workflowID"
	WorkflowExecutionIDKey = "workflowExecutionID"
	WorkflowNameKey        = "workflowName"
	WorkflowOwnerKey       = "workflowOwner"
	StepIDKey              = "stepID"
	StepRefKey             = "stepRef"
)

var OrderedLabelKeys = []string{StepRefKey, StepIDKey, TriggerIDKey, CapabilityIDKey, WorkflowExecutionIDKey, WorkflowIDKey}
