package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// PipelineRunsController manages JobRun requests in the node.
type PipelineRunsController struct {
	App chainlink.Application
}

// Update allows external adapters to resume a JobRun, reporting the result of
// the task and marking it no longer pending.
// Example:
//  "<application>/runs/:RunID"
func (prc *PipelineRunsController) Update(c *gin.Context) {
	authToken := utils.StripBearer(c.Request.Header.Get("Authorization"))
	unscoped := prc.App.GetStore().Unscoped()

	runID, err := models.NewIDFromString(c.Param("RunID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jr, err := unscoped.FindJobRun(runID)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job Run not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if !jr.GetStatus().PendingBridge() {
		jsonAPIError(c, http.StatusMethodNotAllowed, errors.New("Cannot resume a job run that isn't pending"))
		return
	}

	var brr models.BridgeRunResult
	if e := c.ShouldBindJSON(&brr); e != nil {
		jsonAPIError(c, http.StatusInternalServerError, e)
		return
	}

	bt, err := unscoped.PendingBridgeType(jr)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ok, err := models.AuthenticateBridgeType(&bt, authToken)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = prc.App.ResumePendingBridge(runID, brr); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job Run not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, jr, "job run")
}
