package pipeline

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	Spec struct {
		ID           int32 `gorm:"primary_key"`
		DotDagSource string
		TaskSpecs    []TaskSpec
		CreatedAt    time.Time
	}

	TaskSpec struct {
		ID             int32 `gorm:"primary_key"`
		PipelineSpecID int32
		Type           TaskType
		JSON           JSONSerializable `gorm:"type:jsonb"`
		SuccessorID    null.Int
		CreatedAt      time.Time
	}

	Run struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID int32
		CreatedAt      time.Time
	}

	TaskRun struct {
		ID            int64 `gorm:"primary_key"`
		PipelineRunID int64
		Output        *JSONSerializable `gorm:"type:jsonb"`
		Error         null.String
		TaskSpecID    int32
		TaskSpec      TaskSpec
		CreatedAt     time.Time
		FinishedAt    time.Time
	}
)

func (Spec) TableName() string     { return "pipeline_specs" }
func (Run) TableName() string      { return "pipeline_runs" }
func (TaskSpec) TableName() string { return "pipeline_task_specs" }
func (TaskRun) TableName() string  { return "pipeline_task_runs" }

func (r TaskRun) Result() Result {
	var result Result
	if !r.Error.IsZero() {
		result.Error = errors.New(r.Error.ValueOrZero())
	} else if r.Output != nil && r.Output.Value != nil {
		result.Value = r.Output.Value
	}
	return result
}
