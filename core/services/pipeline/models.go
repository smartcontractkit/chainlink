package pipeline

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"gopkg.in/guregu/null.v4"
)

type Spec struct {
	ID              int32           `gorm:"primary_key"`
	DotDagSource    string          `json:"dotDagSource"`
	CreatedAt       time.Time       `json:"-"`
	MaxTaskDuration models.Interval `json:"-"`

	JobID   int32  `gorm:"-" json:"-"`
	JobName string `gorm:"-" json:"-"`
}

func (Spec) TableName() string {
	return "pipeline_specs"
}

func (s Spec) Pipeline() (*Pipeline, error) {
	return Parse(s.DotDagSource)
}

type Run struct {
	ID             int64            `json:"-" gorm:"primary_key"`
	PipelineSpecID int32            `json:"-"`
	PipelineSpec   Spec             `json:"pipelineSpec"`
	Meta           JSONSerializable `json:"meta"`
	// The errors are only ever strings
	// DB example: [null, null, "my error"]
	Errors RunErrors `json:"errors" gorm:"type:jsonb"`
	// The outputs can be anything.
	// DB example: [1234, {"a": 10}, null]
	Outputs          JSONSerializable `json:"outputs" gorm:"type:jsonb"`
	CreatedAt        time.Time        `json:"createdAt"`
	FinishedAt       *time.Time       `json:"finishedAt"`
	PipelineTaskRuns []TaskRun        `json:"taskRuns" gorm:"foreignkey:PipelineRunID;->"`
}

func (Run) TableName() string {
	return "pipeline_runs"
}

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
	for _, err := range r.Errors {
		if !err.IsZero() {
			return true
		}
	}
	return false
}

// Status determines the status of the run.
func (r *Run) Status() RunStatus {
	if r.HasErrors() {
		return RunStatusErrored
	} else if r.FinishedAt != nil {
		return RunStatusCompleted
	}

	return RunStatusInProgress
}

type RunErrors []null.String

func (re *RunErrors) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("RunErrors#Scan received a value of type %T", value)
	}
	return json.Unmarshal(bytes, re)
}

func (re RunErrors) Value() (driver.Value, error) {
	if len(re) == 0 {
		return nil, nil
	}
	return json.Marshal(re)
}

func (re RunErrors) HasError() bool {
	for _, e := range re {
		if !e.IsZero() {
			return true
		}
	}
	return false
}

type TaskRun struct {
	ID            int64             `json:"-" gorm:"primary_key"`
	Type          TaskType          `json:"type"`
	PipelineRun   Run               `json:"-"`
	PipelineRunID int64             `json:"-"`
	Output        *JSONSerializable `json:"output" gorm:"type:jsonb"`
	Error         null.String       `json:"error"`
	CreatedAt     time.Time         `json:"createdAt"`
	FinishedAt    *time.Time        `json:"finishedAt"`
	Index         int32             `json:"index"`
	DotID         string            `json:"dotId"`
}

func (TaskRun) TableName() string {
	return "pipeline_task_runs"
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

func (tr TaskRun) GetDotID() string {
	return tr.DotID
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

// RunStatus represents the status of a run
type RunStatus int

const (
	// RunStatusUnknown is the when the run status cannot be determined.
	RunStatusUnknown RunStatus = iota
	// RunStatusInProgress is used for when a run is actively being executed.
	RunStatusInProgress
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted
)

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Finished returns true if the status is final and can't be changed.
func (s RunStatus) Finished() bool {
	return s.Completed() || s.Errored()
}
