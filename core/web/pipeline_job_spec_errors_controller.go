package web

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// PipelineJobSpecErrorsController manages PipelineJobSpecError requests
type PipelineJobSpecErrorsController struct {
	App chainlink.Application
}

// Destroy deletes a PipelineJobSpecError record from the database, effectively
// silencing the error notification
func (psec *PipelineJobSpecErrorsController) Destroy(c *gin.Context) {
	jobSpec := job.SpecError{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = psec.App.JobORM().DismissError(c.Request.Context(), jobSpec.ID)
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("PipelineJobSpecError not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	psec.App.GetAuditLogger().Audit(audit.JobErrorDismissed, map[string]interface{}{"id": jobSpec.ID})
	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}
