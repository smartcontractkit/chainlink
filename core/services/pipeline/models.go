package pipeline

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"
)

type (
	PipelineSpec struct {
		ID           models.Sha256Hash `gorm:"primary_key"`
		SourceDotDag string
		CreatedAt    time.Time
	}

	PipelineTaskSpec struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID models.Sha256Hash
		// TODO: Task should be a special type?
		Task        Task
		SuccessorID null.Int
		CreatedAt   time.Time
	}

	PipelineRun struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID models.Sha256Hash
		CreatedAt      time.Time
	}

	PipelineTaskRun struct {
		ID                 int64 `gorm:"primary_key"`
		PipelineRunID      int64
		Output             *JSONSerializable
		Error              null.String
		PipelineTaskSpecID int64
		PipelineTaskSpec   PipelineTaskSpec
		CreatedAt          time.Time
		FinishedAt         time.Time
	}
)

func (ptRun PipelineTaskRun) ResultError() error {
	if !ptRun.Error.IsZero() {
		return nil
	}
	return errors.New(ptRun.Error.ValueOrZero())
}
