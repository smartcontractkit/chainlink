//go:build !dev

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// DebugWorkflowController is not available in production builds
type DebugWorkflowController struct {
	App chainlink.Application
}

func (wdc *DebugWorkflowController) GetWorkflows(c *gin.Context)      {}
func (wdc *DebugWorkflowController) GetExecutions(c *gin.Context)     {}
func (wdc *DebugWorkflowController) GetEvents(c *gin.Context)         {}
func (wdc *DebugWorkflowController) GetExecution(c *gin.Context)      {}
func (wdc *DebugWorkflowController) GetStats(c *gin.Context)          {}
func (wdc *DebugWorkflowController) Clear(c *gin.Context)             {}
func (wdc *DebugWorkflowController) GetOrphanEvents(c *gin.Context)   {}
func (wdc *DebugWorkflowController) GetWorkflowEvents(c *gin.Context) {}
