//go:build dev

package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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

// GetOrphanEvents returns events that were emitted without workflow context
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

// GetWorkflowEvents returns workflow-level events
// GET /v2/debug/workflow/:workflowID/events
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
