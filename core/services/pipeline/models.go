package pipeline

import (
	"fmt"
	"strconv"
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
		ID             int32            `json:"-" gorm:"primary_key"`
		DotID          string           `json:"dotId"`
		PipelineSpecID int32            `json:"-"`
		Type           TaskType         `json:"-"`
		JSON           JSONSerializable `json:"-" gorm:"type:jsonb"`
		Index          int32            `json:"-"`
		SuccessorID    null.Int         `json:"-"`
		CreatedAt      time.Time        `json:"-"`
	}

	Run struct {
		ID               int64             `json:"-" gorm:"primary_key"`
		PipelineSpecID   int32             `json:"-"`
		PipelineSpec     Spec              `json:"pipelineSpec"`
		Meta             JSONSerializable  `json:"meta"`
		Errors           *JSONSerializable `json:"errors" gorm:"type:jsonb"`
		Outputs          *JSONSerializable `json:"outputs" gorm:"type:jsonb"`
		CreatedAt        time.Time         `json:"createdAt"`
		FinishedAt       *time.Time        `json:"finishedAt"`
		PipelineTaskRuns []TaskRun         `json:"taskRuns" gorm:"foreignkey:PipelineRunID;association_autoupdate:false;association_autocreate:false"`
	}

	TaskRun struct {
		ID                 int64             `json:"-" gorm:"primary_key"`
		Type               TaskType          `json:"type"`
		PipelineRun        Run               `json:"-"`
		PipelineRunID      int64             `json:"-"`
		Output             *JSONSerializable `json:"output" gorm:"type:jsonb"`
		Error              null.String       `json:"error"`
		PipelineTaskSpecID int32             `json:"-"`
		PipelineTaskSpec   TaskSpec          `json:"taskSpec" gorm:"foreignkey:PipelineTaskSpecID;association_autoupdate:false;association_autocreate:false"`
		CreatedAt          time.Time         `json:"createdAt"`
		FinishedAt         *time.Time        `json:"finishedAt"`
	}
)

func (Spec) TableName() string      { return "pipeline_specs" }
func (SpecError) TableName() string { return "pipeline_spec_errors" }
func (Run) TableName() string       { return "pipeline_runs" }
func (TaskSpec) TableName() string  { return "pipeline_task_specs" }
func (TaskRun) TableName() string   { return "pipeline_task_runs" }

func (r Run) GetID() string {
	return fmt.Sprintf("%v", r.ID)
}

func (r *Run) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	r.ID = int64(ID)
	return nil
}

func (tr TaskRun) GetID() string {
	return fmt.Sprintf("%v", tr.ID)
}

func (tr *TaskRun) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	tr.ID = int64(ID)
	return nil
}

func (s TaskSpec) IsFinalPipelineOutput() bool {
	return s.SuccessorID.IsZero()
}

func (tr TaskRun) DotID() string {
	return tr.PipelineTaskSpec.DotID
}

func (tr TaskRun) Result() Result {
	var result Result
	if !tr.Error.IsZero() {
		result.Error = errors.New(tr.Error.ValueOrZero())
	} else if tr.Output != nil && tr.Output.Val != nil {
		result.Value = tr.Output.Val
	}
	return result
}
