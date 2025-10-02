package events

const (
	ProtoPkg = "workflows.v1"
	// WorkflowStatusChanged represents the Workflow Registry Syncer changing the status of a workflow
	WorkflowStatusChanged string = "WorkflowStatusChanged"
	// WorkflowActivated represents a workflow activated event
	WorkflowActivated string = "WorkflowActivated"
	// WorkflowPaused represents a workflow paused event
	WorkflowPaused string = "WorkflowPaused"
	// WorkflowDeleted represents a workflow deleted event
	WorkflowDeleted string = "WorkflowDeleted"
	// WorkflowExecutionStarted represents a workflow execution started event
	WorkflowExecutionStarted string = "WorkflowExecutionStarted"
	// WorkflowExecutionFinished represents a workflow execution finished event
	WorkflowExecutionFinished string = "WorkflowExecutionFinished"
	// CapabilityExecutionStarted represents a capability execution started event
	CapabilityExecutionStarted string = "CapabilityExecutionStarted"
	// CapabilityExecutionFinished represents a capability execution finished event
	CapabilityExecutionFinished string = "CapabilityExecutionFinished"
	// TriggerExecutionStarted represents a trigger execution started event
	TriggerExecutionStarted string = "TriggerExecutionStarted"
	// WorkflowUserLog represents user log events
	WorkflowUserLog string = "WorkflowUserLog"
	// UserLogs represents user logs events (v1 event type)
	UserLogs string = "UserLogs"

	// SchemaWorkflowStatusChanged represents the schema for workflow status changed events
	SchemaWorkflowStatusChanged string = "/cre-events-workflow-status-changed/v1"
	// SchemaWorkflowStarted represents the schema for workflow started events
	SchemaWorkflowStarted string = "/cre-events-workflow-started/v1"
	// SchemaWorkflowFinished represents the schema for workflow finished events
	SchemaWorkflowFinished string = "/cre-events-workflow-finished/v1"
	// SchemaCapabilityStarted represents the schema for capability started events
	SchemaCapabilityStarted string = "/cre-events-capability-started/v1"
	// SchemaCapabilityFinished represents the schema for capability finished events
	SchemaCapabilityFinished string = "/cre-events-capability-finished/v1"
	// SchemaUserLogs represents the schema for user log events
	SchemaUserLogs string = "/cre-events-user-logs/v1"

	// V2 schema constants
	SchemaWorkflowStartedV2    string = "/cre-events-workflow-started/v2"
	SchemaWorkflowFinishedV2   string = "/cre-events-workflow-finished/v2"
	SchemaCapabilityStartedV2  string = "/cre-events-capability-started/v2"
	SchemaCapabilityFinishedV2 string = "/cre-events-capability-finished/v2"
	SchemaTriggerStartedV2     string = "/cre-events-trigger-started/v2"
	SchemaUserLogsV2           string = "/cre-events-user-logs/v2"
	SchemaWorkflowActivatedV2  string = "/cre-events-workflow-activated/v2"
	SchemaWorkflowPausedV2     string = "/cre-events-workflow-paused/v2"
	SchemaWorkflowDeletedV2    string = "/cre-events-workflow-deleted/v2"

	MeteringReportSchema string = "/workflows/v1/metering.proto"
	MeteringReportDomain string = "platform"
	MeteringReportEntity string = "MeteringReport"
)
