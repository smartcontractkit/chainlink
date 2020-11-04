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

// Index lists all OCR job specs.
// Example:
// "GET <application>/ocr/specs"
func (ocrjsc *OCRJobSpecsController) Index(c *gin.Context) {
	jobs, err := ocrjsc.App.GetStore().ORM.OffChainReportingJobs()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, jobs, "offChainReportingJobSpec")
}

// Show returns the details of a OCR job spec.
// Example:
// "GET <application>/ocr/specs/:ID"
func (ocrjsc *OCRJobSpecsController) Show(c *gin.Context) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobSpec, err = ocrjsc.App.GetStore().ORM.FindOffChainReportingJob(jobSpec.ID)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("OCR job spec not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, jobSpec, "offChainReportingJobSpec")
}

// Create validates, saves and starts a new OCR job spec.
// Example:
// "POST <application>/ocr/specs"
func (ocrjsc *OCRJobSpecsController) Create(c *gin.Context) {
	request := models.CreateOCRJobSpecRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	jobSpec, err := services.ValidatedOracleSpec(request.TOML)
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
// "DELETE <application>/ocr/specs/:ID"
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

	jsonAPIResponseWithStatus(c, nil, "offChainReportingJobSpec", http.StatusNoContent)
}
