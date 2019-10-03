package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"chainlink/core/services"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/store/presenters"
)

// JobSpecsController manages JobSpec requests.
type JobSpecsController struct {
	App services.Application
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
	var jsr models.JobSpecRequest
	if err := c.ShouldBindJSON(&jsr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
	} else if js := models.NewJobFromRequest(jsr); false {
	} else if err := services.ValidateJob(js, jsc.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
	} else if err = jsc.App.AddJob(js); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.JobSpec{JobSpec: js}, "job")
	}
}

// Show returns the details of a JobSpec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Show(c *gin.Context) {
	if id, err := models.NewIDFromString(c.Param("SpecID")); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if j, err := jsc.App.GetStore().FindJob(id); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, jobPresenter(jsc, j), "job")
	}
}

// Destroy soft deletes a job spec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Destroy(c *gin.Context) {
	if id, err := models.NewIDFromString(c.Param("SpecID")); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if err := jsc.App.ArchiveJob(id); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
	}
}

func jobPresenter(jsc *JobSpecsController, job models.JobSpec) presenters.JobSpec {
	st := jsc.App.GetStore()
	jobLinkEarned, _ := st.LinkEarnedFor(&job)
	return presenters.JobSpec{JobSpec: job, Earnings: jobLinkEarned}
}
