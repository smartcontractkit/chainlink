package events

const (
	ProtoPkg = "workflows.v1"
	// WorkflowExecutionStarted represents a workflow execution started event
	WorkflowExecutionStarted string = "WorkflowExecutionStarted"
	// WorkflowExecutionFinished represents a workflow execution finished event
	WorkflowExecutionFinished string = "WorkflowExecutionFinished"
	// CapabilityExecutionStarted represents a capability execution started event
	CapabilityExecutionStarted string = "CapabilityExecutionStarted"
	// CapabilityExecutionFinished represents a capability execution finished event
	CapabilityExecutionFinished string = "CapabilityExecutionFinished"

	// SchemaWorkflowStarted represents the schema for workflow started events
	SchemaWorkflowStarted string = "/cre-events-workflow-started/v1"
	// SchemaWorkflowFinished represents the schema for workflow finished events
	SchemaWorkflowFinished string = "/cre-events-workflow-finished/v1"
	// SchemaCapabilityStarted represents the schema for capability started events
	SchemaCapabilityStarted string = "/cre-events-capability-started/v1"
	// SchemaCapabilityFinished represents the schema for capability finished events
	SchemaCapabilityFinished string = "/cre-events-capability-finished/v1"

	MeteringReportSchema string = "/workflows/v1/metering.proto"
	MeteringReportDomain string = "platform"
	MeteringReportEntity string = "MeteringReport"
)
