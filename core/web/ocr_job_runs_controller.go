package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// OCRJobRunsController manages OCR job run requests.
type OCRJobRunsController struct {
	App chainlink.Application
}

// Index returns all pipeline runs for an OCR job.
// Example:
// "GET <application>/ocr/specs/:ID/runs"
func (ocrjrc *OCRJobRunsController) Index(c *gin.Context) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	var pipelineRuns []pipeline.Run
	err = preloadPipelineRunDependencies(ocrjrc.App.GetStore().DB).
		Joins("INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id").
		Where("jobs.offchainreporting_oracle_spec_id IS NOT NULL").
		Where("jobs.id = ?", jobSpec.ID).
		Order("created_at ASC, id ASC").
		Find(&pipelineRuns).Error

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, pipelineRuns, "offChainReportingJobRun")
}

// Show returns a specified pipeline run.
// Example:
// "GET <application>/ocr/specs/:ID/runs/:runID"
func (ocrjrc *OCRJobRunsController) Show(c *gin.Context) {
	pipelineRun := pipeline.Run{}
	err := pipelineRun.SetID(c.Param("runID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = preloadPipelineRunDependencies(ocrjrc.App.GetStore().DB).
		Where("pipeline_runs.id = ?", pipelineRun.ID).
		First(&pipelineRun).Error

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, pipelineRun, "offChainReportingJobRun")
}

// Create triggers a pipeline run for an OCR job.
// Example:
// "POST <application>/ocr/specs/:ID/runs"
func (ocrjrc *OCRJobRunsController) Create(c *gin.Context) {
	jobSpec := models.JobSpecV2{}
	err := jobSpec.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	jobRunID, err := ocrjrc.App.RunJobV2(c, jobSpec.ID, nil)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, models.OCRJobRun{ID: jobRunID}, "offChainReportingJobRun")
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
