//go:build dev

package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/devobservability"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// WorkflowsDebugController provides dev-only endpoints for accessing workflow execution data
type WorkflowsDebugController struct {
	App chainlink.Application
}

// GetWorkflows returns all workflow IDs that have executions in the store
// GET /v2/debug/workflow
func (wdc *WorkflowsDebugController) GetWorkflows(c *gin.Context) {
	workflows := devobservability.GetStore().GetWorkflows()
	jsonAPIResponse(c, presenters.NewWorkflowDebugWorkflowsResource(workflows), "workflows")
}

// GetExecutions returns all execution IDs for a specific workflow
// GET /v2/debug/workflow/:workflowID/executions?status=<status>
func (wdc *WorkflowsDebugController) GetExecutions(c *gin.Context) {
	workflowID := c.Param("workflowID")
	statusFilter := c.Query("status") // Optional: filter by execution status

	executions := devobservability.GetStore().GetExecutions(workflowID, statusFilter)
	jsonAPIResponse(c, presenters.NewWorkflowDebugExecutionsResource(workflowID, executions), "executions")
}

// GetEvents returns all events for a specific workflow execution
// GET /v2/debug/workflow/:workflowID/executions/:executionID/events
func (wdc *WorkflowsDebugController) GetEvents(c *gin.Context) {
	workflowID := c.Param("workflowID")
	executionID := c.Param("executionID")

	events, err := devobservability.GetStore().GetEvents(workflowID, executionID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, errors.New("execution not found"))
		return
	}

	// Check if client wants formatted output (with decoded protobufs)
	if c.Query("format") == "decoded" {
		jsonAPIResponse(c, presenters.NewWorkflowDebugEventsFormattedResource(workflowID, executionID, events), "events")
	} else {
		jsonAPIResponse(c, presenters.NewWorkflowDebugEventsResource(workflowID, executionID, events), "events")
	}
}

// GetExecution returns detailed information about a specific execution
// GET /v2/debug/workflow/:workflowID/executions/:executionID
func (wdc *WorkflowsDebugController) GetExecution(c *gin.Context) {
	workflowID := c.Param("workflowID")
	executionID := c.Param("executionID")

	execution := devobservability.GetStore().GetExecution(executionID)
	if execution == nil || execution.WorkflowID != workflowID {
		jsonAPIError(c, http.StatusNotFound, errors.New("execution not found"))
		return
	}

	jsonAPIResponse(c, presenters.NewWorkflowDebugExecutionResource(workflowID, executionID, execution), "execution")
}

// GetStats returns statistics about the dev observability store
// GET /v2/debug/workflow/stats
func (wdc *WorkflowsDebugController) GetStats(c *gin.Context) {
	stats := devobservability.GetStore().Stats()
	jsonAPIResponse(c, presenters.NewWorkflowDebugStatsResource(stats), "stats")
}

// Clear clears all data from the dev observability store
// DELETE /v2/workflows_debug
func (wdc *WorkflowsDebugController) Clear(c *gin.Context) {
	devobservability.GetStore().Clear()
	c.JSON(http.StatusNoContent, nil)
}

// GetOrphanEvents returns events that were emitted without workflow/execution context
// GET /v2/debug/workflow/orphan_events
func (wdc *WorkflowsDebugController) GetOrphanEvents(c *gin.Context) {
	limit := 100 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	events := devobservability.GetStore().GetOrphanEvents(limit)

	// Check if client wants formatted output (with decoded protobufs)
	if c.Query("format") == "decoded" {
		jsonAPIResponse(c, presenters.NewWorkflowDebugOrphanEventsFormattedResource(events), "orphan_events")
	} else {
		jsonAPIResponse(c, presenters.NewWorkflowDebugOrphanEventsResource(events), "orphan_events")
	}
}

// GetWorkflowEvents returns workflow-level events (events with workflowID but no executionID)
// GET /v2/debug/workflow/:workflowID/workflow_events
func (wdc *WorkflowsDebugController) GetWorkflowEvents(c *gin.Context) {
	workflowID := c.Param("workflowID")

	limit := 100 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	events := devobservability.GetStore().GetWorkflowEvents(workflowID, limit)

	// Check if client wants formatted output (with decoded protobufs)
	if c.Query("format") == "decoded" {
		jsonAPIResponse(c, presenters.NewWorkflowDebugWorkflowEventsFormattedResource(workflowID, events), "workflow_events")
	} else {
		jsonAPIResponse(c, presenters.NewWorkflowDebugWorkflowEventsResource(workflowID, events), "workflow_events")
	}
}

func parseInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	return i, err
}
