//go:build !dev

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// WorkflowsDebugController is not available in production builds
type WorkflowsDebugController struct {
	App chainlink.Application
}

func (wdc *WorkflowsDebugController) GetWorkflows(c *gin.Context)      {}
func (wdc *WorkflowsDebugController) GetExecutions(c *gin.Context)     {}
func (wdc *WorkflowsDebugController) GetEvents(c *gin.Context)         {}
func (wdc *WorkflowsDebugController) GetExecution(c *gin.Context)      {}
func (wdc *WorkflowsDebugController) GetStats(c *gin.Context)          {}
func (wdc *WorkflowsDebugController) Clear(c *gin.Context)             {}
func (wdc *WorkflowsDebugController) GetOrphanEvents(c *gin.Context)   {}
func (wdc *WorkflowsDebugController) GetWorkflowEvents(c *gin.Context) {}
