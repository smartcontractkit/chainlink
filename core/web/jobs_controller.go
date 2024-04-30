package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
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

	jobs, count, err := jc.App.JobORM().FindJobs(c.Request.Context(), offset, size)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var resources []presenters.JobResource
	for _, individualJob := range jobs {
		resources = append(resources, *presenters.NewJobResource(individualJob))
	}

	paginatedResponse(c, "jobs", size, page, resources, count, err)
}

// Show returns the details of a job
// :ID could be both job ID and external job ID
// Example:
// "GET <application>/jobs/:ID"
func (jc *JobsController) Show(c *gin.Context) {
	ctx := c.Request.Context()
	var err error
	jobSpec := job.Job{}
	if externalJobID, pErr := uuid.Parse(c.Param("ID")); pErr == nil {
		// Find a job by external job ID
		jobSpec, err = jc.App.JobORM().FindJobByExternalJobID(ctx, externalJobID)
	} else if pErr = jobSpec.SetID(c.Param("ID")); pErr == nil {
		// Find a job by job ID
		jobSpec, err = jc.App.JobORM().FindJobTx(ctx, jobSpec.ID)
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

	jb, status, err := jc.validateJobSpec(c.Request.Context(), request.TOML)
	if err != nil {
		jsonAPIError(c, status, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = jc.App.AddJobV2(ctx, &jb)
	if err != nil {
		if errors.Is(errors.Cause(err), job.ErrNoSuchKeyBundle) || errors.As(err, &keystore.KeyNotFoundError{}) || errors.Is(errors.Cause(err), job.ErrNoSuchTransmitterKey) || errors.Is(errors.Cause(err), job.ErrNoSuchSendingKey) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jbj, err := json.Marshal(jb)
	if err == nil {
		jc.App.GetAuditLogger().Audit(audit.JobCreated, map[string]interface{}{"job": string(jbj)})
	} else {
		jc.App.GetLogger().Errorf("Could not send audit log for JobCreation", "err", err)
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

	jc.App.GetAuditLogger().Audit(audit.JobDeleted, map[string]interface{}{"id": j.ID})
	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}

// UpdateJobRequest represents a request to update a job with new toml and start a job (V2).
type UpdateJobRequest struct {
	TOML string `json:"toml"`
}

// Update validates a new TOML for an existing job, stops and deletes existing job, saves and starts a new job.
// Example:
// "PUT <application>/jobs/:ID"
func (jc *JobsController) Update(c *gin.Context) {
	request := UpdateJobRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jb, status, err := jc.validateJobSpec(c.Request.Context(), request.TOML)
	if err != nil {
		jsonAPIError(c, status, err)
		return
	}

	err = jb.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// If the provided job id is not matching any job, delete will fail with 404 leaving state unchanged.
	err = jc.App.DeleteJob(ctx, jb.ID)
	// Error can be either come from ORM or from the activeJobs map.
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "job not found") {
			jsonAPIError(c, http.StatusNotFound, errors.Wrap(err, "failed to update job"))
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	err = jc.App.AddJobV2(ctx, &jb)
	if err != nil {
		if errors.Is(errors.Cause(err), job.ErrNoSuchKeyBundle) || errors.As(err, &keystore.KeyNotFoundError{}) || errors.Is(errors.Cause(err), job.ErrNoSuchTransmitterKey) || errors.Is(errors.Cause(err), job.ErrNoSuchSendingKey) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewJobResource(jb), jb.Type.String())
}

func (jc *JobsController) validateJobSpec(ctx context.Context, tomlString string) (jb job.Job, statusCode int, err error) {
	jobType, err := job.ValidateSpec(tomlString)
	if err != nil {
		return jb, http.StatusUnprocessableEntity, errors.Wrap(err, "failed to parse TOML")
	}
	config := jc.App.GetConfig()
	switch jobType {
	case job.OffchainReporting:
		jb, err = ocr.ValidatedOracleSpecToml(config, jc.App.GetRelayers().LegacyEVMChains(), tomlString)
		if !config.OCR().Enabled() {
			return jb, http.StatusNotImplemented, errors.New("The Offchain Reporting feature is disabled by configuration")
		}
	case job.OffchainReporting2:
		jb, err = validate.ValidatedOracleSpecToml(ctx, config.OCR2(), config.Insecure(), tomlString, jc.App.GetLoopRegistrarConfig())
		if !config.OCR2().Enabled() {
			return jb, http.StatusNotImplemented, errors.New("The Offchain Reporting 2 feature is disabled by configuration")
		}
	case job.DirectRequest:
		jb, err = directrequest.ValidatedDirectRequestSpec(tomlString)
	case job.FluxMonitor:
		jb, err = fluxmonitorv2.ValidatedFluxMonitorSpec(config.JobPipeline(), tomlString)
	case job.Keeper:
		jb, err = keeper.ValidatedKeeperSpec(tomlString)
	case job.Cron:
		jb, err = cron.ValidatedCronSpec(tomlString)
	case job.VRF:
		jb, err = vrfcommon.ValidatedVRFSpec(tomlString)
	case job.Webhook:
		jb, err = webhook.ValidatedWebhookSpec(ctx, tomlString, jc.App.GetExternalInitiatorManager())
	case job.BlockhashStore:
		jb, err = blockhashstore.ValidatedSpec(tomlString)
	case job.BlockHeaderFeeder:
		jb, err = blockheaderfeeder.ValidatedSpec(tomlString)
	case job.Bootstrap:
		jb, err = ocrbootstrap.ValidatedBootstrapSpecToml(tomlString)
	case job.Gateway:
		jb, err = gateway.ValidatedGatewaySpec(tomlString)
	case job.Stream:
		jb, err = streams.ValidatedStreamSpec(tomlString)
	case job.Workflow:
		jb, err = workflows.ValidatedWorkflowSpec(tomlString)
	default:
		return jb, http.StatusUnprocessableEntity, errors.Errorf("unknown job type: %s", jobType)
	}

	if err != nil {
		return jb, http.StatusBadRequest, err
	}
	return jb, 0, nil
}
