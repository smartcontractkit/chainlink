package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/services"
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
		ctx.AbortWithError(http.StatusBadRequest, err)
	} else if err := models.ValidateBulkDeleteRunRequest(request); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
	} else if err := c.App.GetStore().BulkDeleteRuns(request); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}
