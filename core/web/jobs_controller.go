package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gopkg.in/guregu/null.v4"
)

// JobsController manages jobs
type JobsController struct {
	App chainlink.Application
}

// Index lists all jobs
// Example:
// "GET <application>/jobs"
func (jc *JobsController) Index(c *gin.Context) {
	jobs, err := jc.App.GetStore().ORM.JobsV2()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	fmt.Println("Balls", jobs)

	jsonAPIResponse(c, jobs, "jobs")
}

// Show returns the details of a job
// Example:
// "GET <application>/jobs/:ID"
func (jc *JobsController) Show(c *gin.Context) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobSpec, err = jc.App.GetStore().ORM.FindOffChainReportingJob(jobSpec.ID)
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

type GenericJobSpec struct {
	Type          string      `toml:"type"`
	SchemaVersion uint32      `toml:"schemaVersion"`
	Name          null.String `toml:"name"`
}

// Create validates, saves and starts a new job.
// Example:
// "POST <application>/jobs"
func (jc *JobsController) Create(c *gin.Context) {
	request := models.CreateJobSpecRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	genericJS := GenericJobSpec{}
	err := toml.Unmarshal([]byte(request.TOML), &genericJS)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	}

	switch genericJS.Type {
	case string(offchainreporting.JobType):
		jc.createOCR(c, request.TOML)
	case string(models.EthRequestEventJobType):
		jc.createEthRequestEvent(c, request.TOML)
	default:
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("unknown job type: %s", genericJS.Type))
	}

}

func (jc *JobsController) createOCR(c *gin.Context, toml string) {
	jobSpec, err := services.ValidatedOracleSpecToml(toml)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	config := jc.App.GetStore().Config
	if jobSpec.JobType() == offchainreporting.JobType && !config.Dev() && !config.FeatureOffchainReporting() {
		jsonAPIError(c, http.StatusNotImplemented, errors.New("The Offchain Reporting feature is disabled by configuration"))
		return
	}

	jobID, err := jc.App.AddJobV2(c.Request.Context(), jobSpec, jobSpec.Name)
	if err != nil {
		if errors.Cause(err) == job.ErrNoSuchKeyBundle || errors.Cause(err) == job.ErrNoSuchPeerID || errors.Cause(err) == job.ErrNoSuchTransmitterAddress {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	job, err := jc.App.GetStore().ORM.FindOffChainReportingJob(jobID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, job, "offChainReportingJobSpec")
}

func (jc *JobsController) createEthRequestEvent(c *gin.Context, toml string) {
	jobSpec, err := services.ValidatedEthRequestEventSpec(toml)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	jobID, err := jc.App.AddJobV2(c.Request.Context(), jobSpec, jobSpec.Name)
	if err != nil {
		if errors.Cause(err) == job.ErrNoSuchKeyBundle || errors.Cause(err) == job.ErrNoSuchPeerID || errors.Cause(err) == job.ErrNoSuchTransmitterAddress {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	job, err := jc.App.GetStore().ORM.FindOffChainReportingJob(jobID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, job, "ethRequestEventSpec")
}

// Delete soft deletes an OCR job spec.
// Example:
// "DELETE <application>/specs/:ID"
func (jc *JobsController) Delete(c *gin.Context) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = jc.App.DeleteJobV2(c.Request.Context(), jobSpec.ID)
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
