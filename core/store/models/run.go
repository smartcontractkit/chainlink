package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// JobRun tracks the status of a job by holding its TaskRuns and the
// Result of each Run.
type JobRun struct {
	ID             string     `json:"id" gorm:"primary_key;not null"`
	JobSpecID      string     `json:"jobId" gorm:"index;not null;type:varchar(36) REFERENCES job_specs(id)"`
	Result         RunResult  `json:"result"`
	ResultID       uint       `json:"-"`
	RunRequest     RunRequest `json:"-"`
	RunRequestID   uint       `json:"-"`
	Status         RunStatus  `json:"status" gorm:"index"`
	TaskRuns       []TaskRun  `json:"taskRuns"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"index"`
	CompletedAt    null.Time  `json:"completedAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	Initiator      Initiator  `json:"initiator" gorm:"association_autoupdate:false;association_autocreate:false"`
	InitiatorID    uint       `json:"-"`
	CreationHeight *Big       `json:"creationHeight" gorm:"type:varchar(255)"`
	ObservedHeight *Big       `json:"observedHeight" gorm:"type:varchar(255)"`
	Overrides      RunResult  `json:"overrides"`
	OverridesID    uint       `json:"-"`
	DeletedAt      null.Time  `json:"-" gorm:"index"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (jr JobRun) GetID() string {
	return jr.ID
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (jr JobRun) GetName() string {
	return "runs"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (jr *JobRun) SetID(value string) error {
	jr.ID = value
	return nil
}

// ForLogger formats the JobRun for a common formatting in the log.
func (jr JobRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", jr.JobSpecID,
		"run", jr.ID,
		"status", jr.Status,
	}

	if jr.CreationHeight != nil {
		output = append(output, "creation_height", jr.CreationHeight.ToInt())
	}

	if jr.ObservedHeight != nil {
		output = append(output, "observed_height", jr.ObservedHeight.ToInt())
	}

	if jr.Result.HasError() {
		output = append(output, "job_error", jr.Result.Error())
	}

	return append(kvs, output...)
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

// PreviousTaskRun returns the last task to be processed, if it exists
func (jr *JobRun) PreviousTaskRun() *TaskRun {
	index, runnable := jr.NextTaskRunIndex()
	if runnable && index > 0 {
		return &jr.TaskRuns[index-1]
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
	jr.Result.Status = RunStatusErrored
	jr.Status = jr.Result.Status
}

// ApplyResult updates the JobRun's Result and Status
func (jr *JobRun) ApplyResult(result RunResult) {
	jr.Result = result
	jr.Status = result.Status
	if jr.Status.Completed() {
		jr.CompletedAt = null.Time{Time: time.Now(), Valid: true}
	}
}

// MarkCompleted sets the JobRun's status to completed and records the
// completed at time.
func (jr *JobRun) MarkCompleted() {
	jr.Status = RunStatusCompleted
	jr.Result.Status = RunStatusCompleted
	jr.CompletedAt = null.Time{Time: time.Now(), Valid: true}
}

// JobRunsWithStatus filters passed job runs returning those that have
// the desired status, entirely in memory.
func JobRunsWithStatus(runs []JobRun, status RunStatus) []JobRun {
	rval := []JobRun{}
	for _, r := range runs {
		if r.Status == status {
			rval = append(rval, r)
		}
	}
	return rval
}

// RunRequest stores the fields used to initiate the parent job run.
type RunRequest struct {
	ID        uint `gorm:"primary_key"`
	RequestID *string
	TxHash    *common.Hash
	Requester *common.Address
	CreatedAt time.Time
}

// NewRunRequest returns a new RunRequest instance.
func NewRunRequest() RunRequest {
	return RunRequest{CreatedAt: time.Now()}
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	ID                   string    `json:"id" gorm:"primary_key;not null"`
	JobRunID             string    `json:"-" gorm:"index;not null;type:varchar(36) REFERENCES job_runs(id) ON DELETE CASCADE"`
	Result               RunResult `json:"result"`
	ResultID             uint      `json:"-"`
	Status               RunStatus `json:"status"`
	TaskSpec             TaskSpec  `json:"task" gorm:"association_autoupdate:false;association_autocreate:false"`
	TaskSpecID           uint      `json:"-" gorm:"index;not null REFERENCES task_specs(id)"`
	MinimumConfirmations uint64    `json:"minimumConfirmations"`
	CreatedAt            time.Time `json:"-" gorm:"index"`
}

// String returns info on the TaskRun as "ID,Type,Status,Result".
func (tr TaskRun) String() string {
	return fmt.Sprintf("TaskRun(%v,%v,%v,%v)", tr.ID, tr.TaskSpec.Type, tr.Status, tr.Result)
}

// ForLogger formats the TaskRun info for a common formatting in the log.
func (tr *TaskRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"type", tr.TaskSpec.Type,
		"params", tr.TaskSpec.Params,
		"taskrun", tr.ID,
		"status", tr.Status,
	}

	if tr.Result.HasError() {
		output = append(output, "error", tr.Result.Error())
	}

	return append(kvs, output...)
}

// SetError sets this task run to failed and saves the error message
func (tr *TaskRun) SetError(err error) {
	tr.Result.ErrorMessage = null.StringFrom(err.Error())
	tr.Result.Status = RunStatusErrored
	tr.Status = tr.Result.Status
}

// ApplyResult updates the TaskRun's Result and Status
func (tr *TaskRun) ApplyResult(result RunResult) {
	tr.Result = result
	tr.Status = result.Status
}

// MarkCompleted marks the task's status as completed.
func (tr *TaskRun) MarkCompleted() {
	tr.Status = RunStatusCompleted
	tr.Result.Status = RunStatusCompleted
}

// MarkPendingConfirmations marks the task's status as blocked.
func (tr *TaskRun) MarkPendingConfirmations() {
	tr.Status = RunStatusPendingConfirmations
	tr.Result.Status = RunStatusPendingConfirmations
}

// RunResult keeps track of the outcome of a TaskRun or JobRun. It stores the
// Data and ErrorMessage, and contains a field to track the status.
type RunResult struct {
	ID              uint         `json:"-" gorm:"primary_key;auto_increment"`
	CachedJobRunID  string       `json:"jobRunId"`
	CachedTaskRunID string       `json:"taskRunId"`
	Data            JSON         `json:"data" gorm:"type:text"`
	Status          RunStatus    `json:"status"`
	ErrorMessage    null.String  `json:"error"`
	Amount          *assets.Link `json:"amount,omitempty" gorm:"type:varchar(255)"`
}

// ApplyResult saves a value to a RunResult and marks it as completed
func (rr *RunResult) ApplyResult(val interface{}) {
	rr.Status = RunStatusCompleted
	rr.Add("result", val)
}

// Add adds a key and result to the RunResult's JSON payload.
func (rr *RunResult) Add(key string, result interface{}) {
	data, err := rr.Data.Add(key, result)
	if err != nil {
		rr.SetError(err)
		return
	}
	rr.Data = data
}

// SetError marks the result as errored and saves the specified error message
func (rr *RunResult) SetError(err error) {
	rr.ErrorMessage = null.StringFrom(err.Error())
	rr.Status = RunStatusErrored
}

// MarkPendingBridge sets the status to pending_bridge
func (rr *RunResult) MarkPendingBridge() {
	rr.Status = RunStatusPendingBridge
}

// MarkPendingConfirmations sets the status to pending_confirmations.
func (rr *RunResult) MarkPendingConfirmations() {
	rr.Status = RunStatusPendingConfirmations
}

// MarkPendingConnection sets the status to pending_connection.
func (rr *RunResult) MarkPendingConnection() {
	rr.Status = RunStatusPendingConnection
}

// Get searches for and returns the JSON at the given path.
func (rr *RunResult) Get(path string) gjson.Result {
	return rr.Data.Get(path)
}

// ResultString returns the string result of the Data JSON field.
func (rr *RunResult) ResultString() (string, error) {
	val := rr.Result()
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// Result returns the result as a gjson object
func (rr *RunResult) Result() gjson.Result {
	return rr.Get("result")
}

// HasError returns true if the ErrorMessage is present.
func (rr *RunResult) HasError() bool {
	return rr.ErrorMessage.Valid
}

// Error returns the string value of the ErrorMessage field.
func (rr *RunResult) Error() string {
	return rr.ErrorMessage.String
}

// GetError returns the error of a RunResult if it is present.
func (rr *RunResult) GetError() error {
	if rr.HasError() {
		return errors.New(rr.ErrorMessage.ValueOrZero())
	}
	return nil
}

// Merge saves the specified result's data onto the receiving RunResult. The
// input result's data takes preference over the receivers'.
func (rr *RunResult) Merge(in RunResult) error {
	var err error
	rr.Data, err = rr.Data.Merge(in.Data)
	if err != nil {
		return err
	}
	rr.ErrorMessage = in.ErrorMessage
	rr.Amount = in.Amount
	rr.Status = in.Status
	return nil
}

// BridgeRunResult handles the parsing of RunResults from external adapters.
type BridgeRunResult struct {
	RunResult
	ExternalPending bool   `json:"pending"`
	AccessToken     string `json:"accessToken"`
}

// UnmarshalJSON parses the given input and updates the BridgeRunResult in the
// external adapter format.
func (brr *BridgeRunResult) UnmarshalJSON(input []byte) error {
	type biAlias BridgeRunResult
	var anon biAlias
	err := json.Unmarshal(input, &anon)
	*brr = BridgeRunResult(anon)

	if brr.Status.Errored() || brr.HasError() {
		brr.Status = RunStatusErrored
	} else if brr.ExternalPending || brr.Status.PendingBridge() {
		brr.Status = RunStatusPendingBridge
	} else {
		brr.Status = RunStatusCompleted
	}

	return err
}
