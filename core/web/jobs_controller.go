package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// JobsController manages jobs
type JobsController struct {
	App chainlink.Application
}

// Index lists all jobs
// Example:
// "GET <application>/jobs"
func (jc *JobsController) Index(c *gin.Context, size, page, offset int) {
	// Temporary: if no size is passed in, use a large page size. Remove once frontend can handle pagination
	if c.Query("size") == "" {
		size = 1000
	}

	jobs, count, err := jc.App.JobORM().JobsV2(offset, size)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var resources []presenters.JobResource
	for _, job := range jobs {
		resources = append(resources, *presenters.NewJobResource(job))
	}

	paginatedResponse(c, "jobs", size, page, resources, count, err)
}

// Show returns the details of a job
// :ID could be both job ID and external job ID
// Example:
// "GET <application>/jobs/:ID"
func (jc *JobsController) Show(c *gin.Context) {
	var err error
	jobSpec := job.Job{}
	if externalJobID, pErr := uuid.FromString(c.Param("ID")); pErr == nil {
		// Find a job by external job ID
		jobSpec, err = jc.App.JobORM().FindJobByExternalJobID(c.Request.Context(), externalJobID)
	} else if pErr = jobSpec.SetID(c.Param("ID")); pErr == nil {
		// Find a job by job ID
		jobSpec, err = jc.App.JobORM().FindJobTx(jobSpec.ID)
	} else {
		jsonAPIError(c, http.StatusUnprocessableEntity, pErr)
		return
	}
	if err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			jsonAPIError(c, http.StatusNotFound, errors.New("job not found"))
		} else {
			jsonAPIError(c, http.StatusInternalServerError, err)
		}
		return
	}

	jsonAPIResponse(c, presenters.NewJobResource(jobSpec), "jobs")
}

// CreateJobRequest represents a request to create and start a job (V2).
type CreateJobRequest struct {
	TOML string `json:"toml"`
}

// Create validates, saves and starts a new job.
// Example:
// "POST <application>/jobs"
func (jc *JobsController) Create(c *gin.Context) {
	request := CreateJobRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobType, err := job.ValidateSpec(request.TOML)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "failed to parse TOML"))
	}

	var jb job.Job
	config := jc.App.GetConfig()
	switch jobType {
	case job.OffchainReporting:
		jb, err = offchainreporting.ValidatedOracleSpecToml(jc.App.GetChainSet(), request.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting() {
			jsonAPIError(c, http.StatusNotImplemented, errors.New("The Offchain Reporting feature is disabled by configuration"))
			return
		}
	case job.DirectRequest:
		jb, err = directrequest.ValidatedDirectRequestSpec(request.TOML)
	case job.FluxMonitor:
		jb, err = fluxmonitorv2.ValidatedFluxMonitorSpec(jc.App.GetConfig(), request.TOML)
	case job.Keeper:
		jb, err = keeper.ValidatedKeeperSpec(request.TOML)
	case job.Cron:
		jb, err = cron.ValidatedCronSpec(request.TOML)
	case job.VRF:
		jb, err = vrf.ValidatedVRFSpec(request.TOML)
	case job.Webhook:
		jb, err = webhook.ValidatedWebhookSpec(request.TOML, jc.App.GetExternalInitiatorManager())
	default:
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("unknown job type: %s", jobType))
	}
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jb, err = jc.App.AddJobV2(c.Request.Context(), jb, jb.Name)
	if err != nil {
		if errors.Cause(err) == job.ErrNoSuchKeyBundle || errors.Cause(err) == job.ErrNoSuchPeerID || errors.Cause(err) == job.ErrNoSuchTransmitterAddress {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewJobResource(jb), jb.Type.String())
}

// Delete hard deletes a job spec.
// Example:
// "DELETE <application>/specs/:ID"
func (jc *JobsController) Delete(c *gin.Context) {
	j := job.Job{}
	err := j.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Delete the job
	err = jc.App.DeleteJob(c.Request.Context(), j.ID)
	if errors.Cause(err) == gorm.ErrRecordNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))

		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)

		return
	}

	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}
