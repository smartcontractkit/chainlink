package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/gin-gonic/gin"
)

// BulkDeletesController manages background tasks that delete resources given a query
type BulkDeletesController struct {
	App chainlink.Application
}

// Delete removes all runs given a query
// Example:
//  "<application>/bulk_delete_runs"
func (bdc *BulkDeletesController) Delete(c *gin.Context) {
	request := &models.BulkDeleteRunRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err := models.ValidateBulkDeleteRunRequest(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if err := postgres.BulkDeleteRuns(bdc.App.GetStore().DB, request); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "nil", http.StatusNoContent)
}
