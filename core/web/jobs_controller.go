package web

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/core/services/pg"
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

	jobs, count, err := jc.App.JobORM().FindJobs(offset, size)
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
		jobSpec, err = jc.App.JobORM().FindJobByExternalJobID(externalJobID, pg.WithParentCtx(c.Request.Context()))
	} else if pErr = jobSpec.SetID(c.Param("ID")); pErr == nil {
		// Find a job by job ID
		jobSpec, err = jc.App.JobORM().FindJobTx(jobSpec.ID)
	} else {
		jsonAPIError(c, http.StatusUnprocessableEntity, pErr)
		return
	}
	if err != nil {
		if errors.Is(errors.Cause(err), sql.ErrNoRows) {
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
		return
	}

	var jb job.Job
	config := jc.App.GetConfig()
	switch jobType {
	case job.OffchainReporting:
		jb, err = ocr.ValidatedOracleSpecToml(jc.App.GetChains().EVM, request.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting() {
			jsonAPIError(c, http.StatusNotImplemented, errors.New("The Offchain Reporting feature is disabled by configuration"))
			return
		}
	case job.OffchainReporting2:
		jb, err = validate.ValidatedOracleSpecToml(jc.App.GetConfig(), request.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting2() {
			jsonAPIError(c, http.StatusNotImplemented, errors.New("The Offchain Reporting 2 feature is disabled by configuration"))
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
	case job.BlockhashStore:
		jb, err = blockhashstore.ValidatedSpec(request.TOML)
	case job.Bootstrap:
		jb, err = ocrbootstrap.ValidatedBootstrapSpecToml(request.TOML)
	default:
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("unknown job type: %s", jobType))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = jc.App.AddJobV2(ctx, &jb)
	if err != nil {
		if errors.Is(errors.Cause(err), job.ErrNoSuchKeyBundle) || errors.As(err, &keystore.KeyNotFoundError{}) || errors.Is(errors.Cause(err), job.ErrNoSuchTransmitterKey) {
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
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))

		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)

		return
	}

	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}
