package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// PipelineRunsController manages V2 job run requests.
type PipelineRunsController struct {
	App chainlink.Application
}

// Index returns all pipeline runs for a job.
// Example:
// "GET <application>/jobs/:ID/runs"
func (prc *PipelineRunsController) Index(c *gin.Context, size, page, offset int) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	pipelineRuns, count, err := prc.App.GetStore().PipelineRunsByJobID(jobSpec.ID, offset, size)

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
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobRunID, err := prc.App.RunJobV2(c, jobSpec.ID, nil)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, models.PipelineRun{ID: jobRunID}, "offChainReportingPipelineRun")
}

func preloadPipelineRunDependencies(db *gorm.DB) *gorm.DB {
	return db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Where(`pipeline_task_runs.type != 'result'`).
				Order("created_at ASC, id ASC")
		}).
		Preload("PipelineTaskRuns.PipelineTaskSpec")
}
