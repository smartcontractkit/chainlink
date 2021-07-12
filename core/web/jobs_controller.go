package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/job"
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
	jobs, err := jc.App.JobORM().JobsV2()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewJobResources(jobs), "jobs")
}

// Show returns the details of a job
// Example:
// "GET <application>/jobs/:ID"
func (jc *JobsController) Show(c *gin.Context) {
	jobSpec := job.Job{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobSpec, err = jc.App.JobORM().FindJob(jobSpec.ID)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("job not found"))
		return
	}

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewJobResource(jobSpec), "jobs")
}

type GenericJobSpec struct {
	Type          job.Type    `toml:"type"`
	SchemaVersion uint32      `toml:"schemaVersion"`
	Name          null.String `toml:"name"`
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

	genericJS := GenericJobSpec{}
	err := toml.Unmarshal([]byte(request.TOML), &genericJS)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "failed to parse V2 job TOML. HINT: If you are trying to add a V1 job spec (json) via the CLI, try `job_specs create` instead"))
	}

	var js job.Job
	config := jc.App.GetStore().Config
	switch genericJS.Type {
	case job.OffchainReporting:
		js, err = offchainreporting.ValidatedOracleSpecToml(jc.App.GetStore().Config, request.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting() {
			jsonAPIError(c, http.StatusNotImplemented, errors.New("The Offchain Reporting feature is disabled by configuration"))
			return
		}
	case job.DirectRequest:
		js, err = directrequest.ValidatedDirectRequestSpec(request.TOML)
	case job.FluxMonitor:
		js, err = fluxmonitorv2.ValidatedFluxMonitorSpec(jc.App.GetStore().Config, request.TOML)
	case job.Keeper:
		js, err = keeper.ValidatedKeeperSpec(request.TOML)
	case job.Cron:
		js, err = cron.ValidatedCronSpec(request.TOML)
	case job.VRF:
		js, err = vrf.ValidatedVRFSpec(request.TOML)
	case job.Webhook:
		js, err = webhook.ValidatedWebhookSpec(request.TOML, jc.App.GetExternalInitiatorManager())
	default:
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("unknown job type: %s", genericJS.Type))
	}
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jobID, err := jc.App.AddJobV2(c.Request.Context(), js, js.Name)
	if err != nil {
		if errors.Cause(err) == job.ErrNoSuchKeyBundle || errors.Cause(err) == job.ErrNoSuchPeerID || errors.Cause(err) == job.ErrNoSuchTransmitterAddress {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	job, err := jc.App.JobORM().FindJob(jobID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewJobResource(job), job.Type.String())
}

// Delete hard deletes a job spec.
// Example:
// "DELETE <application>/specs/:ID"
func (jc *JobsController) Delete(c *gin.Context) {
	jobSpec := job.Job{}
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

	jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
}
