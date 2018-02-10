package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// JobRun tracks the status of a job by holding its TaskRuns and the
// Result of each Run.
type JobRun struct {
	ID        string    `json:"id" storm:"id,index,unique"`
	JobID     string    `json:"jobId" storm:"index"`
	Status    string    `json:"status" storm:"index"`
	CreatedAt time.Time `json:"createdAt" storm:"index"`
	Result    RunResult `json:"result" storm:"inline"`
	TaskRuns  []TaskRun `json:"taskRuns" storm:"inline"`
}

// ForLogger formats the JobRun for a common formatting in the log.
func (jr *JobRun) ForLogger(kvs ...interface{}) []interface{} {
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
func (jr *JobRun) UnfinishedTaskRuns() []TaskRun {
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
func (jr *JobRun) NextTaskRun() TaskRun {
	return jr.UnfinishedTaskRuns()[0]
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	Task   Task      `json:"task"`
	ID     string    `json:"id" storm:"id,index,unique"`
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
		return TaskRun{}, fmt.Errorf("TaskRun#Merge merging outputs: %v", err.Error())
	}

	rval := tr
	rval.Task.Params = merged
	return rval, nil
}

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
}

// UnmarshalJSON parses the JSON bytes and stores in the *JSON pointer.
func (j *JSON) UnmarshalJSON(b []byte) error {
	if !gjson.Valid(string(b)) {
		return fmt.Errorf("invalid JSON: %v", string(b))
	}
	*j = JSON{gjson.ParseBytes(b)}
	return nil
}

// MarshalJSON returns the JSON data if it already exists, returns
// an empty JSON object as bytes if not.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j.Exists() {
		return j.Bytes(), nil
	}
	return []byte("{}"), nil
}

// Merge combines the given JSON with the existing JSON.
func (j JSON) Merge(j2 JSON) (JSON, error) {
	body := j.Map()
	for key, value := range j2.Map() {
		body[key] = value
	}

	cleaned := map[string]interface{}{}
	for k, v := range body {
		cleaned[k] = v.Value()
	}

	b, err := json.Marshal(cleaned)
	if err != nil {
		return JSON{}, err
	}

	var rval JSON
	return rval, gjson.Unmarshal(b, &rval)
}

// Empty returns true if the JSON does not exist.
func (j JSON) Empty() bool {
	return !j.Exists()
}

// Bytes returns the raw JSON.
func (j JSON) Bytes() []byte {
	return []byte(j.String())
}

// RunResult keeps track of the outcome of a TaskRun. It stores
// the Output and ErrorMessage, if any of either, and contains
// a Pending field to track the status.
type RunResult struct {
	Output       JSON        `json:"output"`
	ErrorMessage null.String `json:"error"`
	Pending      bool        `json:"pending"`
}

// RunResultWithValue returns a new RunResult with the given string
// value as a JSON object.
func RunResultWithValue(val string) RunResult {
	b, err := json.Marshal(map[string]string{"value": val})
	if err != nil {
		return RunResultWithError(err)
	}

	var output JSON
	if err = json.Unmarshal(b, &output); err != nil {
		return RunResultWithError(err)
	}

	return RunResult{Output: output}
}

// RunResultWithError returns a new RunResult with the given
// error message.
func RunResultWithError(err error) RunResult {
	return RunResult{
		ErrorMessage: null.StringFrom(err.Error()),
	}
}

// RunResultPending returns a new RunResult keeping the same
// Output and ErrorMessage as given but with a Pending status
// set to true.
func RunResultPending(input RunResult) RunResult {
	return RunResult{
		Output:       input.Output,
		ErrorMessage: input.ErrorMessage,
		Pending:      true,
	}
}

// Get searches for and returns the JSON at the given path.
func (rr RunResult) Get(path string) (gjson.Result, error) {
	return rr.Output.Get(path), nil
}

func (rr RunResult) value() (gjson.Result, error) {
	return rr.Get("value")
}

// Value returns the string value of the Output JSON field.
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
		return fmt.Errorf("Run Result: ", rr.Error())
	} else {
		return nil
	}
}
