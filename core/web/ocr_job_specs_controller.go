package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// OCRJobSpecsController manages OCR job spec requests.
type OCRJobSpecsController struct {
	App chainlink.Application
}

// Create adds validates, saves, and starts a new OCR job spec.
// Example:
// "<application>/ocr/specs"
func (ocrjsc *OCRJobSpecsController) Create(c *gin.Context) {
	jobSpec, err := services.ValidatedOracleSpec(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	config := ocrjsc.App.GetStore().Config
	if jobSpec.JobType() == offchainreporting.JobType && !config.Dev() && !config.FeatureOffchainReporting() {
		jsonAPIError(c, http.StatusNotImplemented, errors.New("The Offchain Reporting feature is disabled by configuration"))
		return
	}

	jobID, err := ocrjsc.App.AddJobV2(c.Request.Context(), jobSpec)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, struct {
		JobID int32 `json:"jobID"`
	}{jobID})
}

// Delete soft deletes an OCR job spec.
// Example:
// "<application>/ocr/specs/:ID"
func (ocrjsc *OCRJobSpecsController) Delete(c *gin.Context) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = ocrjsc.App.DeleteJobV2(c.Request.Context(), jobSpec.ID)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}
