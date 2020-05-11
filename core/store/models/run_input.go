package models

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// RunInput represents the input for performing a Task
type RunInput struct {
	jobRunID  ID
	taskRunID ID
	data      JSON
	status    RunStatus
}

// NewRunInput creates a new RunInput with arbitrary data
func NewRunInput(jobRunID *ID, taskRunID ID, data JSON, status RunStatus) *RunInput {
	return &RunInput{
		jobRunID:  *jobRunID,
		taskRunID: taskRunID,
		data:      data,
		status:    status,
	}
}

// NewRunInputWithResult creates a new RunInput with a value in the "result" field
func NewRunInputWithResult(jobRunID *ID, taskRunID ID, value interface{}, status RunStatus) *RunInput {
	data, err := JSON{}.Add("result", value)
	if err != nil {
		panic(fmt.Sprintf("invariant violated, add should not fail on empty JSON %v", err))
	}
	return &RunInput{
		jobRunID:  *jobRunID,
		taskRunID: taskRunID,
		data:      data,
		status:    status,
	}
}

// Result returns the result as a gjson object
func (ri RunInput) Result() gjson.Result {
	return ri.data.Get("result")
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
func (ri RunInput) JobRunID() *ID {
	return &ri.jobRunID
}

// TaskRunID returns this RunInput's TaskRunID
func (ri RunInput) TaskRunID() ID {
	return ri.taskRunID
}

func (ri RunInput) CloneWithData(data JSON) RunInput {
	ri.data = data
	return ri
}
