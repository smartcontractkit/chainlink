package web

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

// JobRunsController manages JobRun requests in the node.
type JobRunsController struct {
	App services.Application
}

// Index returns paginated JobRuns for a given JobSpec
// Example:
//  "<application>/runs?jobSpecId=:jobSpecId&size=1&page=2"
func (jrc *JobRunsController) Index(c *gin.Context) {
	id := c.Query("jobSpecId")
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	order := orm.Ascending
	if c.Query("sort") == "-createdAt" {
		order = orm.Descending
	}

	store := jrc.App.GetStore()
	var runs []models.JobRun
	var count int
	if id == "" {
		runs, count, err = store.JobRunsSorted(order, offset, size)
	} else {
		runs, count, err = store.JobRunsSortedFor(id, order, offset, size)
	}

	paginatedResponse(c, "JobRuns", size, page, runs, count, err)
}

// Create starts a new Run for the requested JobSpec.
// Example:
//  "<application>/specs/:SpecID/runs"
func (jrc *JobRunsController) Create(c *gin.Context) {
	id := c.Param("SpecID")

	if j, err := jrc.App.GetStore().FindJob(id); err == orm.ErrorNotFound {
		c.AbortWithError(404, errors.New("Job not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if !j.WebAuthorized() {
		c.AbortWithError(403, errors.New("Job not available on web API, recreate with web initiator"))
	} else if data, err := getRunData(c); err != nil {
		c.AbortWithError(500, err)
	} else if jr, err := services.ExecuteJob(j, j.InitiatorsFor(models.InitiatorWeb)[0], models.RunResult{Data: data}, nil, jrc.App.GetStore()); err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(presenters.JobRun{JobRun: *jr}); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, doc)
	}
}

func getRunData(c *gin.Context) (models.JSON, error) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return models.JSON{}, err
	}
	return models.ParseJSON(b)
}

// Show returns the details of a JobRun.
// Example:
//  "<application>/runs/:RunID"
func (jrc *JobRunsController) Show(c *gin.Context) {
	id := c.Param("RunID")
	if jr, err := jrc.App.GetStore().FindJobRun(id); err == orm.ErrorNotFound {
		c.AbortWithError(404, errors.New("Job Run not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(presenters.JobRun{JobRun: jr}); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, doc)
	}
}

// Update allows external adapters to resume a JobRun, reporting the result of
// the task and marking it no longer pending.
// Example:
//  "<application>/runs/:RunID"
func (jrc *JobRunsController) Update(c *gin.Context) {
	id := c.Param("RunID")
	var brr models.BridgeRunResult
	if jr, err := jrc.App.GetStore().FindJobRun(id); err == orm.ErrorNotFound {
		c.AbortWithError(404, errors.New("Job Run not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if !jr.Result.Status.PendingBridge() {
		c.AbortWithError(405, errors.New("Cannot resume a job run that isn't pending"))
	} else if err := c.ShouldBindJSON(&brr); err != nil {
		c.AbortWithError(500, err)
	} else if bt, err := jrc.App.GetStore().PendingBridgeType(jr); err != nil {
		c.AbortWithError(500, err)
	} else if _, err := bt.Authenticate(utils.StripBearer(c.Request.Header.Get("Authorization"))); err != nil {
		publicError(c, http.StatusUnauthorized, err)
	} else if _, err = services.ResumePendingTask(&jr, jrc.App.GetStore(), brr.RunResult); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, gin.H{"id": jr.ID})
	}
}
