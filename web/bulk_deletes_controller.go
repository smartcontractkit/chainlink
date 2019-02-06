package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// BulkDeletesController manages background tasks that delete resources given a query
type BulkDeletesController struct {
	App services.Application
}

// Delete removes all runs given a query
// Example:
//  "<application>/bulk_delete_runs"
func (c *BulkDeletesController) Delete(ctx *gin.Context) {
	request := &models.BulkDeleteRunRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithError(400, err)
	} else if err := models.ValidateBulkDeleteRunRequest(request); err != nil {
		ctx.AbortWithError(422, err)
	} else if err := c.App.GetStore().BulkDeleteRuns(request); err != nil {
		ctx.AbortWithError(500, err)
	} else {
		ctx.Status(204)
	}
}
