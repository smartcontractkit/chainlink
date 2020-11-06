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
		CreatedAt    time.Time
	}

	SpecError struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID int32
		Description    string
		Occurrences    uint
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	TaskSpec struct {
		ID             int32 `gorm:"primary_key"`
		DotID          string
		PipelineSpecID int32
		Type           TaskType
		JSON           JSONSerializable `gorm:"type:jsonb"`
		Index          int32
		SuccessorID    null.Int
		CreatedAt      time.Time
	}

	Run struct {
		ID             int64 `gorm:"primary_key"`
		PipelineSpecID int32
		Meta           JSONSerializable
		CreatedAt      time.Time
		FinishedAt     time.Time
	}

	TaskRun struct {
		ID                 int64 `gorm:"primary_key"`
		PipelineRun        Run
		PipelineRunID      int64
		Output             *JSONSerializable `gorm:"type:jsonb"`
		Error              null.String
		PipelineTaskSpecID int32
		PipelineTaskSpec   TaskSpec
		CreatedAt          time.Time
		FinishedAt         time.Time
	}
)

func (Spec) TableName() string      { return "pipeline_specs" }
func (SpecError) TableName() string { return "pipeline_spec_errors" }
func (Run) TableName() string       { return "pipeline_runs" }
func (TaskSpec) TableName() string  { return "pipeline_task_specs" }
func (TaskRun) TableName() string   { return "pipeline_task_runs" }

func (s TaskSpec) IsFinalPipelineOutput() bool {
	return s.SuccessorID.IsZero()
}

func (r TaskRun) DotID() string {
	return r.PipelineTaskSpec.DotID
}

func (r TaskRun) Result() Result {
	var result Result
	if !r.Error.IsZero() {
		result.Error = errors.New(r.Error.ValueOrZero())
	} else if r.Output != nil && r.Output.Val != nil {
		result.Value = r.Output.Val
	}
	return result
}
