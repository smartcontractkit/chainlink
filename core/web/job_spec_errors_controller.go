package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// JobSpecErrorsController manages JobSpecError requests
type JobSpecErrorsController struct {
	App chainlink.Application
}

// Destroy deletes a JobSpecError record from the database, effectively
// silencing the error notification
func (jsec *JobSpecErrorsController) Destroy(c *gin.Context) {
	id, err := models.NewIDFromString(c.Param("jobSpecErrorID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = jsec.App.GetStore().DeleteJobSpecError(id)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpecError not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}
