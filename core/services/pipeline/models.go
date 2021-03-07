package pipeline

import (
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	Spec struct {
		ID                int32           `gorm:"primary_key"`
		DotDagSource      string          `json:"dotDagSource"`
		CreatedAt         time.Time       `json:"-"`
		MaxTaskDuration   models.Interval `json:"-"`
		PipelineTaskSpecs []TaskSpec      `json:"-" gorm:"foreignkey:PipelineSpecID;->"`
	}

	TaskSpec struct {
		ID             int32             `json:"-" gorm:"primary_key"`
		DotID          string            `json:"dotId"`
		PipelineSpecID int32             `json:"-"`
		PipelineSpec   Spec              `json:"-"`
		Type           TaskType          `json:"-"`
		JSON           JSONSerializable  `json:"-" gorm:"type:jsonb"`
		Index          int32             `json:"-"`
		SuccessorID    null.Int          `json:"-"`
		CreatedAt      time.Time         `json:"-"`
		BridgeName     *string           `json:"-"`
		Bridge         models.BridgeType `json:"-" gorm:"foreignKey:BridgeName;->"`
	}

	Run struct {
		ID               int64            `json:"-" gorm:"primary_key"`
		PipelineSpecID   int32            `json:"-"`
		PipelineSpec     Spec             `json:"pipelineSpec"`
		Meta             JSONSerializable `json:"meta"`
		Errors           JSONSerializable `json:"errors" gorm:"type:jsonb"`
		Outputs          JSONSerializable `json:"outputs" gorm:"type:jsonb"`
		CreatedAt        time.Time        `json:"createdAt"`
		FinishedAt       *time.Time       `json:"finishedAt"`
		PipelineTaskRuns []TaskRun        `json:"taskRuns" gorm:"foreignkey:PipelineRunID;->"`
	}

	TaskRun struct {
		ID                 int64             `json:"-" gorm:"primary_key"`
		Type               TaskType          `json:"type"`
		PipelineRun        Run               `json:"-"`
		PipelineRunID      int64             `json:"-"`
		Output             *JSONSerializable `json:"output" gorm:"type:jsonb"`
		Error              null.String       `json:"error"`
		PipelineTaskSpecID int32             `json:"-"`
		PipelineTaskSpec   TaskSpec          `json:"taskSpec" gorm:"foreignkey:PipelineTaskSpecID;->"`
		CreatedAt          time.Time         `json:"createdAt"`
		FinishedAt         *time.Time        `json:"finishedAt"`
	}
)

func (Spec) TableName() string     { return "pipeline_specs" }
func (Run) TableName() string      { return "pipeline_runs" }
func (TaskSpec) TableName() string { return "pipeline_task_specs" }
func (TaskRun) TableName() string  { return "pipeline_task_runs" }

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

func (r Run) HasErrors() bool {
	return r.FinalErrors().HasErrors()
}

func (r Run) FinalErrors() (f FinalErrors) {
	f, _ = r.Errors.Val.(FinalErrors)
	return f
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
