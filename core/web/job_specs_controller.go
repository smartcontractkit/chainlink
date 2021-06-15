package web

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
)

// JobSpecsController manages JobSpec requests.
type JobSpecsController struct {
	App chainlink.Application
}

// requireImplented verifies if a Job Spec's feature is enabled according to
// configured policy.
func (jsc *JobSpecsController) requireImplemented(js models.JobSpec) error {
	cfg := jsc.App.GetStore().Config
	if !cfg.Dev() && !cfg.FeatureFluxMonitor() {
		if intrs := js.InitiatorsFor(models.InitiatorFluxMonitor); len(intrs) > 0 {
			return errors.New("The Flux Monitor feature is disabled by configuration")
		}
	}
	if !cfg.EnableLegacyJobPipeline() {
		return errors.New("legacy job pipeline is disabled, cannot add job spec")
	}
	return nil
}

// getAndCheckJobSpec(c) returns a validated job spec from c, or errors. The
// httpStatus return value is only meaningful on error, and in that case
// reflects the type of failure to be reported back to the client.
func (jsc *JobSpecsController) getAndCheckJobSpec(c *gin.Context) (js models.JobSpec, httpStatus int, err error) {
	var jsr models.JobSpecRequest
	if err := c.ShouldBindJSON(&jsr); err != nil {
		// TODO(alx): Better parsing and more specific error messages
		// https://www.pivotaltracker.com/story/show/171164115
		return models.JobSpec{}, http.StatusBadRequest, err
	}
	js = models.NewJobFromRequest(jsr)
	if err := jsc.requireImplemented(js); err != nil {
		return models.JobSpec{}, http.StatusNotImplemented, err
	}
	if err := services.ValidateJob(js, jsc.App.GetStore(), jsc.App.GetKeyStore()); err != nil {
		return models.JobSpec{}, http.StatusBadRequest, err
	}
	return js, 0, nil
}

func showJobPresenter(jsc *JobSpecsController, job models.JobSpec) presenters.JobSpec {
	store := jsc.App.GetStore()
	jobLinkEarned, _ := store.LinkEarnedFor(&job)
	return presenters.JobSpec{JobSpec: job, Errors: job.Errors, Earnings: jobLinkEarned}
}

// Index lists JobSpecs, one page at a time.
// Example:
//  "<application>/specs?size=1&page=2"
func (jsc *JobSpecsController) Index(c *gin.Context, size, page, offset int) {
	var order orm.SortType
	if c.Query("sort") == "-createdAt" {
		order = orm.Descending
	} else {
		order = orm.Ascending
	}

	jobs, count, err := jsc.App.GetStore().JobsSorted(order, offset, size)
	pjs := make([]presenters.JobSpec, len(jobs))
	for i, j := range jobs {
		pjs[i] = presenters.JobSpec{JobSpec: j}
	}

	paginatedResponse(c, "Jobs", size, page, pjs, count, err)
}

// Create adds validates, saves, and starts a new JobSpec.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) Create(c *gin.Context) {
	// NOTE: getAndCheckJobSpec will generate and assign
	// a new UUID to the spec.
	js, httpStatus, err := jsc.getAndCheckJobSpec(c)
	if err != nil {
		jsonAPIError(c, httpStatus, err)
		return
	}

	if err := jsc.App.GetExternalInitiatorManager().Notify(js); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if err := jsc.App.AddJob(js); err != nil {
		switch err := err.(type) {
		case *pgconn.PgError:
			var apiErr error
			if err.ConstraintName == "job_specs_name_key" {
				apiErr = fmt.Errorf("name '%s' already taken", js.Name)
			} else {
				apiErr = err
			}
			jsonAPIError(c, http.StatusConflict, apiErr)
			return
		default:
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	// TODO: https://www.pivotaltracker.com/story/show/171169052
	jsonAPIResponse(c, presenters.JobSpec{JobSpec: js}, "job")
}

// Show returns the details of a JobSpec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Show(c *gin.Context) {
	id, err := models.NewJobIDFromString(c.Param("SpecID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	j, err := jsc.App.GetStore().FindJobWithErrors(id)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, showJobPresenter(jsc, j), "job")
}

// Destroy soft deletes a job spec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Destroy(c *gin.Context) {
	id, err := models.NewJobIDFromString(c.Param("SpecID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = jsc.App.ArchiveJob(id)
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
