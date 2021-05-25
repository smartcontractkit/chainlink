package web

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
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

	pipelineRuns, count, err := prc.App.GetJobORM().PipelineRunsByJobID(jobSpec.ID, offset, size)
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

	err = preloadPipelineRunDependencies(prc.App.GetStore().DB).
		Where("pipeline_runs.id = ?", pipelineRun.ID).
		First(&pipelineRun).Error

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
		var pipelineRun pipeline.Run
		err := preloadPipelineRunDependencies(prc.App.GetStore().DB).
			Where("pipeline_runs.id = ?", jobRunID).
			First(&pipelineRun).Error
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
	pipelineInputs := []pipeline.Result{{Value: string(bodyBytes)}}
	idStr := c.Param("ID")

	// Is it a UUID? Then process it as a webhook job
	var jobUUID models.JobID
	err = jobUUID.UnmarshalText([]byte(idStr))
	if err == nil {
		jobRunID, err2 := prc.App.RunWebhookJobV2(context.Background(), jobUUID, pipelineInputs, pipeline.JSONSerializable{Null: true})
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

func preloadPipelineRunDependencies(db *gorm.DB) *gorm.DB {
	return db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		})
}
