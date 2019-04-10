package web

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// JobRunsController manages JobRun requests in the node.
type JobRunsController struct {
	App services.Application
}

// Index returns paginated JobRuns for a given JobSpec
// Example:
//  "<application>/runs?jobSpecId=:jobSpecId&size=1&page=2"
func (jrc *JobRunsController) Index(c *gin.Context, size, page, offset int) {
	id := c.Query("jobSpecId")

	order := orm.Ascending
	if c.Query("sort") == "-createdAt" {
		order = orm.Descending
	}

	store := jrc.App.GetStore()
	var runs []models.JobRun
	var count int
	var err error
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
		publicError(c, http.StatusNotFound, errors.New("Job not found"))
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if !j.WebAuthorized() {
		publicError(c, http.StatusForbidden, errors.New("Job not available on web API, recreate with web initiator"))
	} else if data, err := getRunData(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if jr, err := services.ExecuteJob(j, j.InitiatorsFor(models.InitiatorWeb)[0], models.RunResult{Data: data}, nil, jrc.App.GetStore()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if doc, err := jsonapi.Marshal(presenters.JobRun{JobRun: *jr}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.Data(http.StatusOK, MediaType, doc)
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
		publicError(c, http.StatusNotFound, errors.New("Job Run not found"))
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if doc, err := jsonapi.Marshal(presenters.JobRun{JobRun: jr}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.Data(http.StatusOK, MediaType, doc)
	}
}

// Update allows external adapters to resume a JobRun, reporting the result of
// the task and marking it no longer pending.
// Example:
//  "<application>/runs/:RunID"
func (jrc *JobRunsController) Update(c *gin.Context) {
	id := c.Param("RunID")
	authToken := utils.StripBearer(c.Request.Header.Get("Authorization"))

	var brr models.BridgeRunResult

	unscoped := jrc.App.GetStore().Unscoped()
	if jr, err := unscoped.FindJobRun(id); err == orm.ErrorNotFound {
		publicError(c, http.StatusNotFound, errors.New("Job Run not found"))
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if !jr.Result.Status.PendingBridge() {
		publicError(c, http.StatusMethodNotAllowed, errors.New("Cannot resume a job run that isn't pending"))
	} else if err := c.ShouldBindJSON(&brr); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if bt, err := unscoped.PendingBridgeType(jr); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if _, err := models.AuthenticateBridgeType(&bt, authToken); err != nil {
		publicError(c, http.StatusUnauthorized, err)
	} else if err = services.ResumePendingTask(&jr, unscoped, brr.RunResult); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"id": jr.ID})
	}
}
