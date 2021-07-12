package models

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
)

// RunInput represents the input for performing a Task
type RunInput struct {
	jobRun    JobRun
	taskRunID uuid.UUID
	data      JSON
	status    RunStatus
}

// NewRunInput creates a new RunInput with arbitrary data
func NewRunInput(jobRun JobRun, taskRunID uuid.UUID, data JSON, status RunStatus) *RunInput {
	return &RunInput{
		jobRun:    jobRun,
		taskRunID: taskRunID,
		data:      data,
		status:    status,
	}
}

// NewRunInputWithResult creates a new RunInput with a value in the "result" field
func NewRunInputWithResult(jobRun JobRun, taskRunID uuid.UUID, value interface{}, status RunStatus) *RunInput {
	data, err := JSON{}.Add(ResultKey, value)
	if err != nil {
		panic(fmt.Sprintf("invariant violated, add should not fail on empty JSON %v", err))
	}
	return &RunInput{
		jobRun:    jobRun,
		taskRunID: taskRunID,
		data:      data,
		status:    status,
	}
}

func (ri RunInput) ResultCollection() gjson.Result {
	return ri.data.Get(ResultCollectionKey)
}

// Result returns the result as a gjson object
func (ri RunInput) Result() gjson.Result {
	return ri.data.Get(ResultKey)
}

// ResultString returns the string result of the Data JSON field.
func (ri RunInput) ResultString() (string, error) {
	val := ri.Result()
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// Status returns the RunInput's status
func (ri RunInput) Status() RunStatus {
	return ri.status
}

// Data returns the RunInput's data
func (ri RunInput) Data() JSON {
	return ri.data
}

// JobRunID returns this RunInput's JobRunID
func (ri RunInput) JobRunID() uuid.UUID {
	return ri.jobRun.ID
}

func (ri RunInput) JobRun() JobRun {
	return ri.jobRun
}

// TaskRunID returns this RunInput's TaskRunID
func (ri RunInput) TaskRunID() uuid.UUID {
	return ri.taskRunID
}

func (ri RunInput) CloneWithData(data JSON) RunInput {
	ri.data = data
	return ri
}
