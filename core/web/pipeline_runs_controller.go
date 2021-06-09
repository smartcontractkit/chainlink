package web

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// PipelineRunsController manages V2 job run requests.
type PipelineRunsController struct {
	App chainlink.Application
}

// Index returns all pipeline runs for a job.
// Example:
// "GET <application>/jobs/:ID/runs"
func (prc *PipelineRunsController) Index(c *gin.Context, size, page, offset int) {
	jobSpec := job.Job{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	pipelineRuns, count, err := prc.App.JobORM().PipelineRunsByJobID(jobSpec.ID, offset, size)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	paginatedResponse(c, "offChainReportingPipelineRun", size, page, pipelineRuns, count, err)
}

// Show returns a specified pipeline run.
// Example:
// "GET <application>/jobs/:ID/runs/:runID"
func (prc *PipelineRunsController) Show(c *gin.Context) {
	pipelineRun := pipeline.Run{}
	err := pipelineRun.SetID(c.Param("runID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	pipelineRun, err = prc.App.PipelineORM().FindRun(pipelineRun.ID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, pipelineRun, "offChainReportingPipelineRun")
}

// Create triggers a pipeline run for a job.
// Example:
// "POST <application>/jobs/:ID/runs"
func (prc *PipelineRunsController) Create(c *gin.Context) {
	respondWithPipelineRun := func(jobRunID int64) {
		pipelineRun, err := prc.App.PipelineORM().FindRun(jobRunID)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		jsonAPIResponse(c, pipelineRun, "offChainReportingPipelineRun")
	}

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	idStr := c.Param("ID")

	// Is it a UUID? Then process it as a webhook job
	jobUUID, err := uuid.FromString(idStr)
	if err == nil {
		jobRunID, err2 := prc.App.RunWebhookJobV2(context.Background(), jobUUID, string(bodyBytes), pipeline.JSONSerializable{Null: true})
		if err2 != nil {
			jsonAPIError(c, http.StatusInternalServerError, err2)
			return
		}
		respondWithPipelineRun(jobRunID)
		return
	}

	// Is it an int32? Then process it regardless of type
	var jobID int32
	jobID64, err := strconv.ParseInt(idStr, 10, 32)
	if err == nil {
		jobID = int32(jobID64)
		jobRunID, err := prc.App.RunJobV2(context.Background(), jobID, nil)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		respondWithPipelineRun(jobRunID)
		return
	}

	jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("bad job ID"))
}

// Resume finishes a task and resumes the pipeline run.
// Example:
// "PATCH <application>/jobs/:ID/runs/:runID"
func (prc *PipelineRunsController) Resume(c *gin.Context) {
	// TODO: lookup by UUID
	// TODO: json payload needs to contain a task run result
	// TODO: we need to give enough data to the bridge to map to a specific task run, simply run id is not enough
	pipelineRun := pipeline.Run{}
	err := pipelineRun.SetID(c.Param("runID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	// TODO: turn TaskRun ID into pipelineRun ID
	// TODO: use SELECT FOR UPDATE in a transaction to lock the pipeline_runs row
	// update the task run entry for task with new data
	// re-trigger run if necessary
	// Q: -> how do we solve for cases where a run is already ongoing?
	// A: -> I guess some sort of a run registry
	// could work since we're guaranteed that if the run is running, it'll have to wait for the lock and can't stop
	// while we're changing things. This is why I also liked having a larger scheduler that's not per-run

	// equivalent on trying to suspend run:
	// grab select for update row lock
	// select pending task runs and see if any of them are written as complete in the db
	// if so, immediately update the result data and re-execute
	// else, write the updated task runs to the db and stop execution

	pipelineRun, err = prc.App.PipelineORM().FindRun(pipelineRun.ID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, pipelineRun, "offChainReportingPipelineRun")
}
