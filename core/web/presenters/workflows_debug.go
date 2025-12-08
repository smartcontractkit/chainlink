//go:build dev

package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/devobservability"
)

// WorkflowDebugWorkflowsResource represents the list of workflow IDs
type WorkflowDebugWorkflowsResource struct {
	JAID
	WorkflowIDs []string `json:"workflowIds"`
}

// NewWorkflowDebugWorkflowsResource creates a new workflows resource
func NewWorkflowDebugWorkflowsResource(workflowIDs []string) *WorkflowDebugWorkflowsResource {
	return &WorkflowDebugWorkflowsResource{
		JAID:        NewJAID("workflows_debug"),
		WorkflowIDs: workflowIDs,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugWorkflowsResource) GetName() string {
	return "workflows_debug"
}

// WorkflowDebugExecutionsResource represents the list of executions for a workflow
type WorkflowDebugExecutionsResource struct {
	JAID
	WorkflowID string                              `json:"workflowId"`
	Executions []devobservability.ExecutionSummary `json:"executions"`
}

// NewWorkflowDebugExecutionsResource creates a new executions resource
func NewWorkflowDebugExecutionsResource(workflowID string, executions []devobservability.ExecutionSummary) *WorkflowDebugExecutionsResource {
	return &WorkflowDebugExecutionsResource{
		JAID:       NewJAID("workflow_executions"),
		WorkflowID: workflowID,
		Executions: executions,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugExecutionsResource) GetName() string {
	return "workflow_executions"
}

// WorkflowDebugEventsResource represents the events for a specific execution
type WorkflowDebugEventsResource struct {
	JAID
	WorkflowID  string                        `json:"workflowId"`
	ExecutionID string                        `json:"executionId"`
	Events      []devobservability.EventEntry `json:"events"`
}

// NewWorkflowDebugEventsResource creates a new events resource
func NewWorkflowDebugEventsResource(workflowID, executionID string, events []devobservability.EventEntry) *WorkflowDebugEventsResource {
	return &WorkflowDebugEventsResource{
		JAID:        NewJAID("workflow_events"),
		WorkflowID:  workflowID,
		ExecutionID: executionID,
		Events:      events,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugEventsResource) GetName() string {
	return "workflow_events"
}

// WorkflowDebugExecutionResource represents detailed execution information
type WorkflowDebugExecutionResource struct {
	JAID
	WorkflowID  string                          `json:"workflowId"`
	ExecutionID string                          `json:"executionId"`
	Execution   *devobservability.ExecutionData `json:"execution"`
}

// NewWorkflowDebugExecutionResource creates a new execution resource
func NewWorkflowDebugExecutionResource(workflowID, executionID string, execution *devobservability.ExecutionData) *WorkflowDebugExecutionResource {
	return &WorkflowDebugExecutionResource{
		JAID:        NewJAID("workflow_execution"),
		WorkflowID:  workflowID,
		ExecutionID: executionID,
		Execution:   execution,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugExecutionResource) GetName() string {
	return "workflow_execution"
}

// WorkflowDebugStatsResource represents statistics about the dev observability store
type WorkflowDebugStatsResource struct {
	JAID
	TotalWorkflows  int            `json:"totalWorkflows"`
	TotalExecutions int            `json:"totalExecutions"`
	TotalEvents     int            `json:"totalEvents"`
	WorkflowStats   map[string]int `json:"workflowStats"`
}

// NewWorkflowDebugStatsResource creates a new stats resource
func NewWorkflowDebugStatsResource(stats map[string]interface{}) *WorkflowDebugStatsResource {
	resource := &WorkflowDebugStatsResource{
		JAID:          NewJAID("workflow_stats"),
		WorkflowStats: make(map[string]int),
	}

	if v, ok := stats["total_workflows"].(int); ok {
		resource.TotalWorkflows = v
	}
	if v, ok := stats["total_executions"].(int); ok {
		resource.TotalExecutions = v
	}
	if v, ok := stats["total_events"].(int); ok {
		resource.TotalEvents = v
	}
	if v, ok := stats["workflow_stats"].(map[string]int); ok {
		resource.WorkflowStats = v
	}

	return resource
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugStatsResource) GetName() string {
	return "workflow_stats"
}

// WorkflowDebugOrphanEventsResource represents events without workflow/execution context
type WorkflowDebugOrphanEventsResource struct {
	JAID
	Events []devobservability.EventEntry `json:"events"`
}

// NewWorkflowDebugOrphanEventsResource creates a new orphan events resource
func NewWorkflowDebugOrphanEventsResource(events []devobservability.EventEntry) *WorkflowDebugOrphanEventsResource {
	return &WorkflowDebugOrphanEventsResource{
		JAID:   NewJAID("orphan_events"),
		Events: events,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugOrphanEventsResource) GetName() string {
	return "orphan_events"
}

// WorkflowDebugWorkflowEventsResource represents workflow-level events (events with workflowID but no executionID)
type WorkflowDebugWorkflowEventsResource struct {
	JAID
	WorkflowID string                        `json:"workflowId"`
	Events     []devobservability.EventEntry `json:"events"`
}

// NewWorkflowDebugWorkflowEventsResource creates a new workflow events resource
func NewWorkflowDebugWorkflowEventsResource(workflowID string, events []devobservability.EventEntry) *WorkflowDebugWorkflowEventsResource {
	return &WorkflowDebugWorkflowEventsResource{
		JAID:       NewJAID("workflow_events"),
		WorkflowID: workflowID,
		Events:     events,
	}
}

// GetName implements the api2go EntityNamer interface
func (r WorkflowDebugWorkflowEventsResource) GetName() string {
	return "workflow_events"
}
