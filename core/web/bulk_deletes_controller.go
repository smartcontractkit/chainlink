package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// BulkDeletesController manages background tasks that delete resources given a query
type BulkDeletesController struct {
	App services.Application
}

// Delete removes all runs given a query
// Example:
//  "<application>/bulk_delete_runs"
func (bdc *BulkDeletesController) Delete(c *gin.Context) {
	request := &models.BulkDeleteRunRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		publicError(c, http.StatusBadRequest, err)
	} else if err := models.ValidateBulkDeleteRunRequest(request); err != nil {
		publicError(c, http.StatusUnprocessableEntity, err)
	} else if err := bdc.App.GetStore().BulkDeleteRuns(request); err != nil {
		publicError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponseWithStatus(c, nil, "nil", http.StatusNoContent)
	}
}
