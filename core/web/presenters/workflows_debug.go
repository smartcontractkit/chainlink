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

// WorkflowDebugStatsResource represents statistics about the dev observability store
type WorkflowDebugStatsResource struct {
	JAID
	TotalWorkflows int            `json:"totalWorkflows"`
	TotalEvents    int            `json:"totalEvents"`
	WorkflowStats  map[string]int `json:"workflowStats"`
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
