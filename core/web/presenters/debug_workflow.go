//go:build dev

package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/devobservability"
)

// DebugWorkflowsResource represents the list of workflow IDs
type DebugWorkflowsResource struct {
	JAID
	WorkflowIDs []string `json:"workflowIds"`
}

// NewDebugWorkflowsResource creates a new workflows resource
func NewDebugWorkflowsResource(workflowIDs []string) *DebugWorkflowsResource {
	return &DebugWorkflowsResource{
		JAID:        NewJAID("workflows_debug"),
		WorkflowIDs: workflowIDs,
	}
}

// GetName implements the api2go EntityNamer interface
func (r DebugWorkflowsResource) GetName() string {
	return "workflows_debug"
}

// DebugWorkflowStatsResource represents statistics about the dev observability store
type DebugWorkflowStatsResource struct {
	JAID
	TotalWorkflows int            `json:"totalWorkflows"`
	TotalEvents    int            `json:"totalEvents"`
	WorkflowStats  map[string]int `json:"workflowStats"`
}

// NewDebugWorkflowStatsResource creates a new stats resource
func NewDebugWorkflowStatsResource(stats map[string]interface{}) *DebugWorkflowStatsResource {
	resource := &DebugWorkflowStatsResource{
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
func (r DebugWorkflowStatsResource) GetName() string {
	return "workflow_stats"
}

// DebugWorkflowOrphanEventsResource represents events without workflow/execution context
type DebugWorkflowOrphanEventsResource struct {
	JAID
	Events []devobservability.EventEntry `json:"events"`
}

// NewDebugWorkflowOrphanEventsResource creates a new orphan events resource
func NewDebugWorkflowOrphanEventsResource(events []devobservability.EventEntry) *DebugWorkflowOrphanEventsResource {
	return &DebugWorkflowOrphanEventsResource{
		JAID:   NewJAID("orphan_events"),
		Events: events,
	}
}

// GetName implements the api2go EntityNamer interface
func (r DebugWorkflowOrphanEventsResource) GetName() string {
	return "orphan_events"
}

// DebugWorkflowEventsResource represents workflow-level events (events with workflowID but no executionID)
type DebugWorkflowEventsResource struct {
	JAID
	WorkflowID string                        `json:"workflowId"`
	Events     []devobservability.EventEntry `json:"events"`
}

// NewDebugWorkflowEventsResource creates a new workflow events resource
func NewDebugWorkflowEventsResource(workflowID string, events []devobservability.EventEntry) *DebugWorkflowEventsResource {
	return &DebugWorkflowEventsResource{
		JAID:       NewJAID("workflow_events"),
		WorkflowID: workflowID,
		Events:     events,
	}
}

// GetName implements the api2go EntityNamer interface
func (r DebugWorkflowEventsResource) GetName() string {
	return "workflow_events"
}
