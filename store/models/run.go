package models

import (
	"fmt"
	"time"

	null "gopkg.in/guregu/null.v3"
)

type JobRun struct {
	ID        string    `storm:"id,index,unique"`
	JobID     string    `storm:"index"`
	Status    string    `storm:"index"`
	CreatedAt time.Time `storm:"index"`
	Result    RunResult `storm:"inline"`
	TaskRuns  []TaskRun `storm:"inline"`
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

func (jr *JobRun) NextTaskRun() TaskRun { return jr.UnfinishedTaskRuns()[0] }

type TaskRun struct {
	Task
	ID     string `storm:"id,index,unique"`
	Status string
	Result RunResult
}

func (tr TaskRun) Completed() bool { return tr.Status == StatusCompleted }
func (tr TaskRun) Errored() bool   { return tr.Status == StatusErrored }

type Output map[string]null.String

type RunResult struct {
	Output       Output
	ErrorMessage null.String
	Pending      bool
}

func RunResultWithValue(val string) RunResult {
	return RunResult{
		Output: Output{"value": null.StringFrom(val)},
	}
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

func (rr RunResult) Value() string      { return rr.value().String }
func (rr RunResult) NullValue() bool    { return !rr.value().Valid }
func (rr RunResult) HasError() bool     { return rr.ErrorMessage.Valid }
func (rr RunResult) Error() string      { return rr.ErrorMessage.String }
func (rr RunResult) value() null.String { return rr.Output["value"] }

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
