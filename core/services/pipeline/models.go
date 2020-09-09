package pipeline

import (
	"github.com/smartcontractkit/chainlink/core/store/models"
	"time"
)

type (
	PipelineSpec struct {
		ID           models.Sha256Hash `gorm:"primary_key"`
		SourceDotDag string
		CreatedAt    time.Time
	}

	PipelineTaskSpec struct {
		ID           int64 `gorm:"primary_key"`
		PipelineSpec PipelineSpec
		TaskSpec     models.JSON
		CreatedAt    time.Time
	}

	PipelineRun struct {
		ID           int64 `gorm:"primary_key"`
		PipelineSpec PipelineSpec
		CreatedAt    time.Time
	}

	PipelineTaskRun struct {
		ID               int64 `gorm:"primary_key"`
		PipelineRun      PipelineRun
		Output           models.JSON
		Error            string
		PipelineTaskSpec PipelineTaskSpec
		CreatedAt        time.Time
		FinishedAt       time.Time
	}
)
