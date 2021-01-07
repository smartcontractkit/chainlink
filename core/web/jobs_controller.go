package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
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
	jobs, err := jc.App.GetJobORM().JobsV2()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, jobs, "jobs")
}

// Show returns the details of a job
// Example:
// "GET <application>/jobs/:ID"
func (jc *JobsController) Show(c *gin.Context) {
	jobSpec := job.SpecDB{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobSpec, err = jc.App.GetJobORM().FindJob(jobSpec.ID)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("job not found"))
		return
	}

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, jobSpec, "offChainReportingJobSpec")
}

type GenericJobSpec struct {
	Type          job.Type    `toml:"type"`
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
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "failed to parse V2 job TOML. HINT: If you are trying to add a V1 job spec (json) via the CLI, try `job_specs create` instead"))
	}

	switch genericJS.Type {
	case job.OffchainReporting:
		jc.createOCR(c, request.TOML)
	case job.DirectRequest:
		jc.createDirectRequest(c, request.TOML)
	default:
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("unknown job type: %s", genericJS.Type))
	}

}

func (jc *JobsController) createOCR(c *gin.Context, toml string) {
	jobSpec, err := services.ValidatedOracleSpecToml(jc.App.GetStore().Config, toml)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	config := jc.App.GetStore().Config
	if jobSpec.Type == job.OffchainReporting && !config.Dev() && !config.FeatureOffchainReporting() {
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

	job, err := jc.App.GetJobORM().FindJob(jobID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, job, "offChainReportingJobSpec")
}

func (jc *JobsController) createDirectRequest(c *gin.Context, toml string) {
	jobSpec, err := services.ValidatedDirectRequestSpec(toml)
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

	job, err := jc.App.GetJobORM().FindJob(jobID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, job, "DirectRequestSpec")
}

// Delete soft deletes an OCR job spec.
// Example:
// "DELETE <application>/specs/:ID"
func (jc *JobsController) Delete(c *gin.Context) {
	jobSpec := job.SpecDB{}
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
