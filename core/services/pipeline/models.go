package pipeline

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	PipelineSpec struct {
		ID           int64 `gorm:"primary_key"`
		JobSpecID    *models.ID
		SourceDotDag string
		TaskSpecs    []PipelineTaskSpec
		CreatedAt    time.Time
	}

	PipelineTaskSpec struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID int64
		TaskJson       JSONSerializable `gorm:"type:jsonb"`
		SuccessorID    null.Int64
		CreatedAt      time.Time
	}

	PipelineRun struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID int64
		CreatedAt      time.Time
	}

	PipelineTaskRun struct {
		ID                 int64 `gorm:"primary_key"`
		PipelineRunID      int64
		Output             *JSONSerializable `gorm:"type:jsonb"`
		Error              null.String
		PipelineTaskSpecID int64
		PipelineTaskSpec   PipelineTaskSpec
		CreatedAt          time.Time
		FinishedAt         time.Time
	}
)

func (r PipelineTaskRun) Result() Result {
	var result Result
	if !r.Error.IsZero() {
		result.Error = errors.New(r.Error.ValueOrZero())
	} else if r.Output != nil && r.Output.Value != nil {
		result.Value = r.Output.Value
	}
	return result
}
