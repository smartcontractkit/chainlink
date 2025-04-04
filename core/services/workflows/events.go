package workflows

const (
	EventsProtoPkg = "pb"
	// EventWorkflowExecutionStarted represents a workflow execution started event
	EventWorkflowExecutionStarted string = "WorkflowExecutionStarted"
	// EventWorkflowExecutionFinished represents a workflow execution finished event
	EventWorkflowExecutionFinished string = "WorkflowExecutionFinished"
	// EventCapabilityExecutionStarted represents a capability execution started event
	EventCapabilityExecutionStarted string = "CapabilityExecutionStarted"
	// EventCapabilityExecutionFinished represents a capability execution finished event
	EventCapabilityExecutionFinished string = "CapabilityExecutionFinished"

	// SchemaWorkflowStarted represents the schema for workflow started events
	SchemaWorkflowStarted string = "/cre-events-workflow-started/v1"
	// SchemaWorkflowFinished represents the schema for workflow finished events
	SchemaWorkflowFinished string = "/cre-events-workflow-finished/v1"
	// SchemaCapabilityStarted represents the schema for capability started events
	SchemaCapabilityStarted string = "/cre-events-capability-started/v1"
	// SchemaCapabilityFinished represents the schema for capability finished events
	SchemaCapabilityFinished string = "/cre-events-capability-finished/v1"
)
