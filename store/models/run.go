package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

type JobRun struct {
	ID        string    `json:"id" storm:"id,index,unique"`
	JobID     string    `json:"jobId" storm:"index"`
	Status    string    `json:"status" storm:"index"`
	CreatedAt time.Time `json:"createdAt" storm:"index"`
	Result    RunResult `json:"result" storm:"inline"`
	TaskRuns  []TaskRun `json:"taskRuns" storm:"inline"`
}

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

func (jr *JobRun) NextTaskRun() TaskRun {
	return jr.UnfinishedTaskRuns()[0]
}

type TaskRun struct {
	Task   Task      `json:"task"`
	ID     string    `json:"id" storm:"id,index,unique"`
	Status string    `json:"status"`
	Result RunResult `json:"result"`
}

func (tr TaskRun) Completed() bool {
	return tr.Status == StatusCompleted
}

func (tr TaskRun) Errored() bool {
	return tr.Status == StatusErrored
}

func (tr TaskRun) String() string {
	return fmt.Sprintf("TaskRun(%v,%v,%v,%v)", tr.ID, tr.Task.Type, tr.Status, tr.Result)
}

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

func (tr TaskRun) MergeTaskParams(j JSON) (TaskRun, error) {
	merged, err := tr.Task.Params.Merge(j)
	if err != nil {
		return TaskRun{}, fmt.Errorf("TaskRun#Merge merging outputs: %v", err.Error())
	}

	rval := tr
	rval.Task.Params = merged
	return rval, nil
}

type JSON struct {
	gjson.Result
}

func (j *JSON) UnmarshalJSON(b []byte) error {
	if !gjson.Valid(string(b)) {
		return fmt.Errorf("invalid JSON: %v", string(b))
	}
	*j = JSON{gjson.ParseBytes(b)}
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j.Exists() {
		return j.Bytes(), nil
	}
	return []byte("{}"), nil
}

func (j JSON) Merge(j2 JSON) (JSON, error) {
	body := j.Map()
	for key, value := range j2.Map() {
		body[key] = value
	}
	str, err := convertToJSON(body)
	if err != nil {
		return JSON{}, err
	}

	var rval JSON
	return rval, gjson.Unmarshal([]byte(str), &rval)
}

func (j JSON) Empty() bool {
	return !j.Exists()
}

func (j JSON) Bytes() []byte {
	return []byte(j.String())
}

func convertToJSON(body map[string]gjson.Result) (string, error) {
	str := "{"
	first := true

	for key, value := range body {
		if first {
			first = false
		} else {
			str += ","
		}
		b, err := json.Marshal(value.Value())
		if err != nil {
			return "", err
		}
		str += fmt.Sprintf(`"%v": %v`, key, string(b))
	}

	return (str + "}"), nil
}

type RunResult struct {
	Output       JSON        `json:"output"`
	ErrorMessage null.String `json:"error"`
	Pending      bool        `json:"pending"`
}

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

func RunResultWithError(err error) RunResult {
	return RunResult{
		ErrorMessage: null.StringFrom(err.Error()),
	}
}

func RunResultPending(input RunResult) RunResult {
	return RunResult{
		Output:       input.Output,
		ErrorMessage: input.ErrorMessage,
		Pending:      true,
	}
}

func (rr RunResult) Get(path string) (gjson.Result, error) {
	return rr.Output.Get(path), nil
}

func (rr RunResult) value() (gjson.Result, error) {
	return rr.Get("value")
}

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

func (rr RunResult) HasError() bool {
	return rr.ErrorMessage.Valid
}

func (rr RunResult) Error() string {
	return rr.ErrorMessage.String
}

func (rr RunResult) SetError(err error) {
	rr.ErrorMessage = null.StringFrom(err.Error())
}

func (rr RunResult) GetError() error {
	if rr.HasError() {
		return fmt.Errorf("Run Result: ", rr.Error())
	} else {
		return nil
	}
}
