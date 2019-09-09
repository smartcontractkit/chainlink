package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	} else if runID, err := models.NewIDFromString(id); err == nil {
		runs, count, err = store.JobRunsSortedFor(runID, order, offset, size)
	} else {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	}

	paginatedResponse(c, "JobRuns", size, page, runs, count, err)
}

// Create starts a new Run for the requested JobSpec.
// Example:
//  "<application>/specs/:SpecID/runs"
func (jrc *JobRunsController) Create(c *gin.Context) {
	if id, err := models.NewIDFromString(c.Param("SpecID")); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if j, err := jrc.App.GetStore().FindJob(id); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if initiator, err := getAuthenticatedInitiator(c, j); err != nil {
		jsonAPIError(c, http.StatusForbidden, err)
	} else if data, err := getRunData(c); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if jr, err := services.ExecuteJob(j, *initiator, models.RunResult{Data: data}, nil, jrc.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.JobRun{JobRun: *jr}, "job run")
	}
}

// getInitiator returns the Job Spec's initiator for the given web context.
func getAuthenticatedInitiator(c *gin.Context, js models.JobSpec) (*models.Initiator, error) {
	if _, ok := authenticatedUser(c); ok {
		webInitiators := js.InitiatorsFor(models.InitiatorWeb)
		if len(webInitiators) == 0 {
			return nil, errors.New("Job not available on web API, recreate with initiator type 'InitiatorWeb'")
		}
		return &webInitiators[0], nil
	} else if ei, ok := authenticatedEI(c); ok {
		var initiator *models.Initiator
		for _, i := range js.InitiatorsFor(models.InitiatorExternal) {
			if strings.ToLower(i.Name) == ei.Name {
				initiator = &i
				break
			}
		}
		if initiator == nil {
			return nil, fmt.Errorf("Job not available via External Initiator '%s'", ei.Name)
		}
		return initiator, nil
	} else {
		return nil, errors.New("authentication required")
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
	if id, err := models.NewIDFromString(c.Param("RunID")); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if jr, err := jrc.App.GetStore().FindJobRun(id); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job run not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.JobRun{JobRun: jr}, "job run")
	}
}

// Update allows external adapters to resume a JobRun, reporting the result of
// the task and marking it no longer pending.
// Example:
//  "<application>/runs/:RunID"
func (jrc *JobRunsController) Update(c *gin.Context) {
	var brr models.BridgeRunResult

	authToken := utils.StripBearer(c.Request.Header.Get("Authorization"))
	unscoped := jrc.App.GetStore().Unscoped()

	if id, err := models.NewIDFromString(c.Param("RunID")); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if jr, err := unscoped.FindJobRun(id); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job Run not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if !jr.Result.Status.PendingBridge() {
		jsonAPIError(c, http.StatusMethodNotAllowed, errors.New("Cannot resume a job run that isn't pending"))
	} else if err := c.ShouldBindJSON(&brr); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if bt, err := unscoped.PendingBridgeType(jr); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if ok, err := models.AuthenticateBridgeType(&bt, authToken); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else if err = services.ResumePendingTask(&jr, unscoped, brr.RunResult); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, jr, "job run")
	}
}
