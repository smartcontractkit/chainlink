package models

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// JobRun tracks the status of a job by holding its TaskRuns and the
// Result of each Run.
type JobRun struct {
	ID             *ID          `json:"id" gorm:"primary_key;not null"`
	JobSpecID      *ID          `json:"jobId" gorm:"index;not null;type:varchar(36) REFERENCES job_specs(id)"`
	Result         RunResult    `json:"result"`
	ResultID       uint         `json:"-"`
	RunRequest     RunRequest   `json:"-"`
	RunRequestID   uint         `json:"-"`
	Status         RunStatus    `json:"status" gorm:"index"`
	TaskRuns       []TaskRun    `json:"taskRuns"`
	CreatedAt      time.Time    `json:"createdAt" gorm:"index"`
	FinishedAt     null.Time    `json:"finishedAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Initiator      Initiator    `json:"initiator" gorm:"association_autoupdate:false;association_autocreate:false"`
	InitiatorID    uint         `json:"-"`
	CreationHeight *Big         `json:"creationHeight"`
	ObservedHeight *Big         `json:"observedHeight"`
	Overrides      JSON         `json:"overrides"`
	DeletedAt      null.Time    `json:"-" gorm:"index"`
	Payment        *assets.Link `json:"payment,omitempty"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (jr JobRun) GetID() string {
	return jr.ID.String()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (jr JobRun) GetName() string {
	return "runs"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (jr *JobRun) SetID(value string) error {
	return jr.ID.UnmarshalText([]byte(value))
}

// ForLogger formats the JobRun for a common formatting in the log.
func (jr JobRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", jr.JobSpecID.String(),
		"run", jr.ID.String(),
		"status", jr.Status,
	}

	if jr.CreationHeight != nil {
		output = append(output, "creation_height", jr.CreationHeight.ToInt())
	}

	if jr.ObservedHeight != nil {
		output = append(output, "observed_height", jr.ObservedHeight.ToInt())
	}

	if jr.HasError() {
		output = append(output, "job_error", jr.ErrorString())
	}

	if jr.Status == "completed" {
		output = append(output, "link_earned", jr.Payment)
	}

	return append(kvs, output...)
}

// HasError returns true if this JobRun has errored
func (jr JobRun) HasError() bool {
	return jr.Status.Errored()
}

// NextTaskRunIndex returns the position of the next unfinished task
func (jr *JobRun) NextTaskRunIndex() (int, bool) {
	for index, tr := range jr.TaskRuns {
		if tr.Status.CanStart() {
			return index, true
		}
	}
	return 0, false
}

// NextTaskRun returns the next immediate TaskRun in the list
// of unfinished TaskRuns.
func (jr *JobRun) NextTaskRun() *TaskRun {
	nextTaskIndex, runnable := jr.NextTaskRunIndex()
	if runnable {
		return &jr.TaskRuns[nextTaskIndex]
	}
	return nil
}

// TasksRemain returns true if there are unfinished tasks left for this job run
func (jr *JobRun) TasksRemain() bool {
	_, runnable := jr.NextTaskRunIndex()
	return runnable
}

// SetError sets this job run to failed and saves the error message
func (jr *JobRun) SetError(err error) {
	jr.Result.ErrorMessage = null.StringFrom(err.Error())
	jr.Status = RunStatusErrored
	jr.FinishedAt = null.TimeFrom(time.Now())
}

// ApplyOutput updates the JobRun's Result and Status
func (jr *JobRun) ApplyOutput(result RunOutput) {
	if result.HasError() {
		jr.SetError(result.Error())
		return
	}
	jr.Result.Data = result.Data()
	jr.setStatus(result.Status())
}

// ApplyBridgeRunResult saves the input from a BridgeAdapter
func (jr *JobRun) ApplyBridgeRunResult(result BridgeRunResult) {
	if result.HasError() {
		jr.SetError(result.GetError())
	}
	jr.Result.Data = result.Data
	jr.setStatus(result.Status)
}

func (jr *JobRun) setStatus(status RunStatus) {
	jr.Status = status
	if jr.Status.Completed() && jr.TasksRemain() {
		jr.Status = RunStatusInProgress
	} else if jr.Status.Finished() {
		jr.FinishedAt = null.TimeFrom(time.Now())
	}
}

// ResultString returns the "result" value as a string if possible
func (jr *JobRun) ResultString() (string, error) {
	val := jr.Result.Data.Get("result")
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// ErrorString returns the error as a string if present, otherwise "".
func (jr *JobRun) ErrorString() string {
	return jr.Result.ErrorMessage.ValueOrZero()
}

// RunRequest stores the fields used to initiate the parent job run.
type RunRequest struct {
	ID        uint `gorm:"primary_key"`
	RequestID *string
	TxHash    *common.Hash
	BlockHash *common.Hash
	Requester *common.Address
	CreatedAt time.Time
	Payment   *assets.Link
}

// NewRunRequest returns a new RunRequest instance.
func NewRunRequest() RunRequest {
	return RunRequest{CreatedAt: time.Now()}
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	ID                   *ID           `json:"id" gorm:"primary_key;not null"`
	JobRunID             *ID           `json:"-" gorm:"index;not null;type:varchar(36) REFERENCES job_runs(id) ON DELETE CASCADE"`
	Result               RunResult     `json:"result"`
	ResultID             uint          `json:"-"`
	Status               RunStatus     `json:"status"`
	TaskSpec             TaskSpec      `json:"task" gorm:"association_autoupdate:false;association_autocreate:false"`
	TaskSpecID           uint          `json:"-" gorm:"index;not null REFERENCES task_specs(id)"`
	MinimumConfirmations clnull.Uint32 `json:"minimumConfirmations"`
	Confirmations        clnull.Uint32 `json:"confirmations"`
	CreatedAt            time.Time     `json:"-" gorm:"index"`
}

// String returns info on the TaskRun as "ID,Type,Status,Result".
func (tr TaskRun) String() string {
	return fmt.Sprintf("TaskRun(%v,%v,%v,%v)", tr.ID.String(), tr.TaskSpec.Type, tr.Status, tr.Result)
}

// SetError sets this task run to failed and saves the error message
func (tr *TaskRun) SetError(err error) {
	tr.Result.ErrorMessage = null.StringFrom(err.Error())
	tr.Status = RunStatusErrored
}

// ApplyBridgeRunResult updates the TaskRun's Result and Status
func (tr *TaskRun) ApplyBridgeRunResult(result BridgeRunResult) {
	if result.HasError() {
		tr.SetError(result.GetError())
	}
	tr.Result.Data = result.Data
	tr.Status = result.Status
}

// ApplyOutput updates the TaskRun's Result and Status
func (tr *TaskRun) ApplyOutput(result RunOutput) {
	if result.HasError() {
		tr.SetError(result.Error())
		return
	}
	tr.Result.Data = result.Data()
	tr.Status = result.Status()
}

// ResultString returns the "result" value as a string if possible
func (tr *TaskRun) ResultString() (string, error) {
	val := tr.Result.Data.Get("result")
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// RunResult keeps track of the outcome of a TaskRun or JobRun. It stores the
// Data and ErrorMessage.
type RunResult struct {
	ID           uint        `json:"-" gorm:"primary_key;auto_increment"`
	Data         JSON        `json:"data" gorm:"type:text"`
	ErrorMessage null.String `json:"error"`
}
