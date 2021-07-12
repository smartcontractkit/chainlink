package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// JobRunsController manages JobRun requests in the node.
type JobRunsController struct {
	App chainlink.Application
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
	var completedCount int
	var erroredCount int
	var err error
	if id == "" {
		runs, count, err = store.JobRunsSorted(order, offset, size)
	} else {
		var runID models.JobID
		runID, err = models.NewJobIDFromString(id)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}

		runs, count, completedCount, erroredCount, err = store.JobRunsSortedFor(runID, order, offset, size)
	}
	meta := make(map[string]interface{})
	meta["completed"] = completedCount
	meta["errored"] = erroredCount
	paginatedResponseWithMeta(c, "JobRuns", size, page, runs, count, err, meta)
}

// Create starts a new Run for the requested JobSpec.
// Example:
//  "<application>/specs/:SpecID/runs"
func (jrc *JobRunsController) Create(c *gin.Context) {
	id, err := models.NewJobIDFromString(c.Param("SpecID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	j, err := jrc.App.GetStore().Unscoped().FindJobSpec(id)

	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job not found"))
		return
	}
	if j.DeletedAt.Valid {
		jsonAPIError(c, http.StatusGone, errors.New("Job spec not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	initiator, err := getAuthenticatedInitiator(c, j)
	if err != nil {
		jsonAPIError(c, http.StatusForbidden, err)
		return
	}

	data, err := getRunData(c)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jr, err := jrc.App.Create(j.ID, initiator, nil, &models.RunRequest{RequestParams: data})
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.JobRun{JobRun: *jr}, "job run")
}

// getInitiator returns the Job Spec's initiator for the given web context.
func getAuthenticatedInitiator(c *gin.Context, js models.JobSpec) (*models.Initiator, error) {
	if _, ok := authenticatedUser(c); ok {
		webInitiators := js.InitiatorsFor(models.InitiatorWeb)
		if len(webInitiators) == 0 {
			return nil, errors.New("Job not available on web API, recreate with initiator type 'InitiatorWeb'")
		}
		return &webInitiators[0], nil
	}
	if ei, ok := authenticatedEI(c); ok {
		initiator := js.InitiatorExternal(ei.Name)
		if initiator == nil {
			return nil, fmt.Errorf("job not available via External Initiator '%s'", ei.Name)
		}
		return initiator, nil
	}
	return nil, errors.New("authentication required")
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
	id, err := models.NewJobIDFromString(c.Param("RunID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jr, err := jrc.App.GetStore().FindJobRun(id.UUID())
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job run not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.JobRun{JobRun: jr}, "job run")
}

// Update allows external adapters to resume a JobRun, reporting the result of
// the task and marking it no longer pending.
// Example:
//  "<application>/runs/:RunID"
func (jrc *JobRunsController) Update(c *gin.Context) {
	authToken := utils.StripBearer(c.Request.Header.Get("Authorization"))

	runID, err := models.NewJobIDFromString(c.Param("RunID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jr, err := jrc.App.GetStore().FindJobRunIncludingArchived(runID.UUID())
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job Run not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if !jr.GetStatus().PendingBridge() {
		jsonAPIError(c, http.StatusMethodNotAllowed, errors.New("Cannot resume a job run that isn't pending"))
		return
	}

	var brr models.BridgeRunResult
	if e := c.ShouldBindJSON(&brr); e != nil {
		jsonAPIError(c, http.StatusInternalServerError, e)
		return
	}

	bt, err := jrc.App.GetStore().PendingBridgeType(jr)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ok, err := models.AuthenticateBridgeType(&bt, authToken)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = jrc.App.ResumePendingBridge(runID.UUID(), brr); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job Run not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, jr, "job run")
}

// Cancel stops a Run from continuing.
// Example:
//  "<application>/runs/:RunID/cancellation"
func (jrc *JobRunsController) Cancel(c *gin.Context) {
	id, err := models.NewJobIDFromString(c.Param("RunID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jr, err := jrc.App.Cancel(id.UUID())
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("Job run not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.JobRun{JobRun: *jr}, "job run")
}
