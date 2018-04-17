package models

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
	Initiator      Initiator    `json:"initiator"`
	CreationHeight *hexutil.Big `json:"creationHeight"`
}

// ForLogger formats the JobRun for a common formatting in the log.
func (jr JobRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", jr.JobID,
		"run", jr.ID,
		"status", jr.Status,
	}

	if jr.Result.HasError() {
		output = append(output, "error", jr.Result.Error())
	}

	return append(kvs, output...)
}

// UnfinishedTaskRuns returns a list of TaskRuns for a JobRun
// which are not Completed or Errored.
func (jr JobRun) UnfinishedTaskRuns() []TaskRun {
	unfinished := jr.TaskRuns
	for _, tr := range jr.TaskRuns {
		if tr.Status.Completed() {
			unfinished = unfinished[1:]
		} else if tr.Status.Errored() {
			return []TaskRun{}
		} else {
			return unfinished
		}
	}
	return unfinished
}

// NextTaskRun returns the next immediate TaskRun in the list
// of unfinished TaskRuns.
func (jr JobRun) NextTaskRun() TaskRun {
	return jr.UnfinishedTaskRuns()[0]
}

// Runnable checks that the number of confirmations have passed since the
// job's creation height to determine if the JobRun can be started. Returns
// true for non-EthereumListener (runlog & ethlog) initiators.
func (jr JobRun) Runnable(currentHeight *IndexableBlockNumber, minConfs uint64) bool {
	if jr.CreationHeight == nil || currentHeight == nil {
		return true
	}

	diff := new(big.Int).Sub(currentHeight.ToInt(), jr.CreationHeight.ToInt())
	min := new(big.Int).SetUint64(minConfs)
	min = min.Sub(min, big.NewInt(1))
	return diff.Cmp(min) >= 0
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
	ID     string    `json:"id" storm:"id,unique"`
	Result RunResult `json:"result"`
	Status RunStatus `json:"status"`
	Task   TaskSpec  `json:"task"`
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

// MergeTaskParams merges the existing parameters on a TaskRun with the given JSON.
func (tr TaskRun) MergeTaskParams(j JSON) (TaskRun, error) {
	merged, err := tr.Task.Params.Merge(j)
	if err != nil {
		return tr, fmt.Errorf("TaskRun#Merge merging params: %v", err.Error())
	}

	tr.Task.Params = merged
	return tr, nil
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

// RunResult keeps track of the outcome of a TaskRun. It stores
// the Data and ErrorMessage, if any of either, and contains
// a Pending field to track the status.
type RunResult struct {
	JobRunID     string      `json:"jobRunId"`
	Data         JSON        `json:"data"`
	Status       RunStatus   `json:"status"`
	ErrorMessage null.String `json:"error"`
}

// WithValue returns a copy of the RunResult, overriding the "value" field of
// Data and setting the status to completed.
func (rr RunResult) WithValue(val string) RunResult {
	data, err := rr.Data.Add("value", val)
	if err != nil {
		return rr.WithError(err)
	}
	rr.Status = RunStatusCompleted
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

// Get searches for and returns the JSON at the given path.
func (rr RunResult) Get(path string) (gjson.Result, error) {
	return rr.Data.Get(path), nil
}

func (rr RunResult) value() (gjson.Result, error) {
	return rr.Get("value")
}

// Value returns the string value of the Data JSON field.
func (rr RunResult) Value() (string, error) {
	val, err := rr.value()
	if err != nil {
		return "", err
	}
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
	} else {
		return nil
	}
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
	ExternalPending bool `json:"pending"`
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
