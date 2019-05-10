package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
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
	js := models.NewJob()
	if err := c.ShouldBindJSON(&js); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
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
	id := c.Param("SpecID")
	if j, err := jsc.App.GetStore().FindJob(id); err == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if runs, err := jsc.App.GetStore().JobRunsFor(j.ID); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, jobPresenter(j, runs), "job")
	}
}

// Destroy soft deletes a job spec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Destroy(c *gin.Context) {
	id := c.Param("SpecID")
	if err := jsc.App.ArchiveJob(id); err == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("JobSpec not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponseWithStatus(c, nil, "job", http.StatusNoContent)
	}
}

func jobPresenter(j models.JobSpec, runs []models.JobRun) presenters.JobSpec {
	pruns := make([]presenters.JobRun, len(runs))
	for i, r := range runs {
		pruns[i] = presenters.JobRun{r}
	}
	return presenters.JobSpec{JobSpec: j, Runs: pruns}
}
