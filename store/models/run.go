package models

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// JobRun tracks the status of a job by holding its TaskRuns and the
// Result of each Run.
type JobRun struct {
	ID          string    `json:"id" storm:"id,unique"`
	JobID       string    `json:"jobId" storm:"index"`
	Status      string    `json:"status" storm:"index"`
	Result      RunResult `json:"result" storm:"inline"`
	TaskRuns    []TaskRun `json:"taskRuns" storm:"inline"`
	CreatedAt   time.Time `json:"createdAt" storm:"index"`
	CompletedAt null.Time `json:"completedAt"`
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
		if tr.Completed() {
			unfinished = unfinished[1:]
		} else if tr.Errored() {
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

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	Task   TaskSpec  `json:"task"`
	ID     string    `json:"id" storm:"id,unique"`
	Status string    `json:"status"`
	Result RunResult `json:"result"`
}

// Completed returns true if the TaskRun status is StatusCompleted.
func (tr TaskRun) Completed() bool {
	return tr.Status == StatusCompleted
}

// Errored returns true if the TaskRun status is StatusErrored.
func (tr TaskRun) Errored() bool {
	return tr.Status == StatusErrored
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

// RunResult keeps track of the outcome of a TaskRun. It stores
// the Data and ErrorMessage, if any of either, and contains
// a Pending field to track the status.
type RunResult struct {
	JobRunID     string      `json:"jobRunId"`
	Data         JSON        `json:"data"`
	ErrorMessage null.String `json:"error"`
	Pending      bool        `json:"pending"`
}

// WithValue returns a copy of the RunResult, overriding the "value" field of
// Data and setting Pending to false.
func (rr RunResult) WithValue(val string) RunResult {
	data, err := rr.Data.Add("value", val)
	if err != nil {
		return rr.WithError(err)
	}
	rr.Pending = false
	rr.Data = data
	return rr
}

// WithValue returns a copy of the RunResult, setting the error field
// and setting Pending to false.
func (rr RunResult) WithError(err error) RunResult {
	rr.ErrorMessage = null.StringFrom(err.Error())
	rr.Pending = false
	return rr
}

// MarkPending returns a copy of RunResult but with Pending set to true.
func (rr RunResult) MarkPending() RunResult {
	rr.Pending = true
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

// SetError stores the given error in the ErrorMessage field.
func (rr RunResult) SetError(err error) {
	rr.ErrorMessage = null.StringFrom(err.Error())
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
	if in.Pending || rr.Pending {
		in.Pending = true
	}
	return in, nil
}
