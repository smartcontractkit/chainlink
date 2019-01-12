package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// JobRun tracks the status of a job by holding its TaskRuns and the
// Result of each Run.
type JobRun struct {
	ID             string       `json:"id" storm:"id,unique"`
	JobID          string       `json:"jobId" storm:"index"`
	Result         RunResult    `json:"result" storm:"inline"`
	Status         RunStatus    `json:"status" storm:"index"`
	TaskRuns       []TaskRun    `json:"taskRuns" storm:"inline"`
	CreatedAt      time.Time    `json:"createdAt" storm:"index"`
	CompletedAt    null.Time    `json:"completedAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Initiator      Initiator    `json:"initiator"`
	CreationHeight *hexutil.Big `json:"creationHeight"`
	ObservedHeight *hexutil.Big `json:"observedHeight"`
	Overrides      RunResult    `json:"overrides"`
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
		"job", jr.JobID,
		"run", jr.ID,
		"status", jr.Status,
	}

	if jr.CreationHeight != nil {
		output = append(output, "creation_height", jr.CreationHeight.ToInt())
	}

	if jr.Result.HasError() {
		output = append(output, "job_error", jr.Result.Error())
	}

	return append(kvs, output...)
}

// NextTaskRunIndex returns the position of the next unfinished task
func (jr JobRun) NextTaskRunIndex() (int, bool) {
	for index, tr := range jr.TaskRuns {
		if !(tr.Status.Completed() || tr.Status.Errored()) {
			return index, true
		}
	}
	return 0, false
}

// NextTaskRun returns the next immediate TaskRun in the list
// of unfinished TaskRuns.
func (jr JobRun) NextTaskRun() *TaskRun {
	nextTaskIndex, runnable := jr.NextTaskRunIndex()
	if runnable {
		return &jr.TaskRuns[nextTaskIndex]
	}
	return nil
}

// PreviousTaskRun returns the last task to be processed, if it exists
func (jr JobRun) PreviousTaskRun() *TaskRun {
	index, runnable := jr.NextTaskRunIndex()
	if runnable && index > 0 {
		return &jr.TaskRuns[index-1]
	}
	return nil
}

// TasksRemain returns true if there are unfinished tasks left for this job run
func (jr JobRun) TasksRemain() bool {
	_, runnable := jr.NextTaskRunIndex()
	return runnable
}

// ApplyResult updates the JobRun's Result and Status
func (jr JobRun) ApplyResult(result RunResult) JobRun {
	jr.Result = result
	jr.Status = result.Status
	if jr.Status.Completed() {
		jr.CompletedAt = null.Time{Time: time.Now(), Valid: true}
	}
	return jr
}

// MarkCompleted sets the JobRun's status to completed and records the
// completed at time.
func (jr JobRun) MarkCompleted() JobRun {
	jr.Status = RunStatusCompleted
	jr.Result.Status = RunStatusCompleted
	jr.CompletedAt = null.Time{Time: time.Now(), Valid: true}
	return jr
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	ID                   string    `json:"id" storm:"id,unique"`
	Result               RunResult `json:"result"`
	Status               RunStatus `json:"status"`
	Task                 TaskSpec  `json:"task"`
	MinimumConfirmations uint64    `json:"minimumConfirmations"`
}

// String returns info on the TaskRun as "ID,Type,Status,Result".
func (tr TaskRun) String() string {
	return fmt.Sprintf("TaskRun(%v,%v,%v,%v)", tr.ID, tr.Task.Type, tr.Status, tr.Result)
}

// ForLogger formats the TaskRun info for a common formatting in the log.
func (tr TaskRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"type", tr.Task.Type,
		"params", tr.Task.Params,
		"taskrun", tr.ID,
		"status", tr.Status,
	}

	if tr.Result.HasError() {
		output = append(output, "error", tr.Result.Error())
	}

	return append(kvs, output...)
}

// ApplyResult updates the TaskRun's Result and Status
func (tr TaskRun) ApplyResult(result RunResult) TaskRun {
	tr.Result = result
	tr.Status = result.Status
	return tr
}

// MarkCompleted marks the task's status as completed.
func (tr TaskRun) MarkCompleted() TaskRun {
	tr.Status = RunStatusCompleted
	tr.Result.Status = RunStatusCompleted
	return tr
}

// MarkPendingConfirmations marks the task's status as blocked.
func (tr TaskRun) MarkPendingConfirmations() TaskRun {
	tr.Status = RunStatusPendingConfirmations
	tr.Result.Status = RunStatusPendingConfirmations
	return tr
}

// RunResult keeps track of the outcome of a TaskRun or JobRun. It stores the
// Data and ErrorMessage, and contains a field to track the status.
type RunResult struct {
	JobRunID     string       `json:"jobRunId"`
	Data         JSON         `json:"data"`
	Status       RunStatus    `json:"status"`
	ErrorMessage null.String  `json:"error"`
	Amount       *assets.Link `json:"amount,omitempty"`
}

// WithValue returns a copy of the RunResult, overriding the "value" field of
// Data and setting the status to completed.
func (rr RunResult) WithValue(val interface{}) RunResult {
	rr.Status = RunStatusCompleted
	return rr.Add("value", val)
}

// Add adds a key and value to the RunResult's JSON payload.
func (rr RunResult) Add(key string, value interface{}) RunResult {
	data, err := rr.Data.Add(key, value)
	if err != nil {
		return rr.WithError(err)
	}
	rr.Data = data
	return rr
}

// WithNull returns a copy of the RunResult, overriding the "value" field of
// Data to null.
func (rr RunResult) WithNull() RunResult {
	data, err := rr.Data.Add("value", nil)
	if err != nil {
		return rr.WithError(err)
	}
	rr.Data = data
	return rr
}

// WithError returns a copy of the RunResult, setting the error field
// and setting the status to in progress.
func (rr RunResult) WithError(err error) RunResult {
	rr.ErrorMessage = null.StringFrom(err.Error())
	rr.Status = RunStatusErrored
	return rr
}

// MarkPendingBridge returns a copy of RunResult but with status set to pending_bridge.
func (rr RunResult) MarkPendingBridge() RunResult {
	rr.Status = RunStatusPendingBridge
	return rr
}

// MarkPendingConfirmations returns a copy of RunResult but with status set to pending_confirmations.
func (rr RunResult) MarkPendingConfirmations() RunResult {
	rr.Status = RunStatusPendingConfirmations
	return rr
}

// MarkPendingConnection returns a copy of RunResult but with status set to pending_connection.
func (rr RunResult) MarkPendingConnection() RunResult {
	rr.Status = RunStatusPendingConnection
	return rr
}

// Get searches for and returns the JSON at the given path.
func (rr RunResult) Get(path string) gjson.Result {
	return rr.Data.Get(path)
}

func (rr RunResult) value() gjson.Result {
	return rr.Get("value")
}

// Value returns the string value of the Data JSON field.
func (rr RunResult) Value() (string, error) {
	val := rr.value()
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string value")
	}
	return val.String(), nil
}

// HasError returns true if the ErrorMessage is present.
func (rr RunResult) HasError() bool {
	return rr.ErrorMessage.Valid
}

// Error returns the string value of the ErrorMessage field.
func (rr RunResult) Error() string {
	return rr.ErrorMessage.String
}

// GetError returns the error of a RunResult if it is present.
func (rr RunResult) GetError() error {
	if rr.HasError() {
		return fmt.Errorf("Run Result: %v", rr.Error())
	}
	return nil
}

// Merge returns a copy which is the result of joining the input RunResult
// with the instance it is called on, preferring the RunResult values passed in,
// but using the existing values if the input RunResult values are of their
// respective zero value.
//
// Returns an error if called on a RunResult that already has an error.
func (rr RunResult) Merge(in RunResult) (RunResult, error) {
	if rr.HasError() {
		err := fmt.Errorf("Cannot merge onto a RunResult with error: %v", rr.Error())
		return rr, err
	}

	merged, err := rr.Data.Merge(in.Data)
	if err != nil {
		return in, fmt.Errorf("TaskRun#Merge merging JSON: %v", err.Error())
	}
	in.Data = merged
	if len(in.JobRunID) == 0 {
		in.JobRunID = rr.JobRunID
	}
	if in.Status.Errored() || rr.Status.Errored() {
		in.Status = RunStatusErrored
	} else if in.Status.PendingBridge() || rr.Status.PendingBridge() {
		in = in.MarkPendingBridge()
	}
	return in, nil
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
