package models

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// RunInput represents the input for performing a Task
type RunInput struct {
	jobRunID ID
	data     JSON
	Meta     JSON
	status   RunStatus
}

// NewRunInput creates a new RunInput with arbitrary data
func NewRunInput(jobRunID *ID, data JSON, meta JSON, status RunStatus) *RunInput {
	return &RunInput{
		jobRunID: *jobRunID,
		data:     data,
		Meta:     meta,
		status:   status,
	}
}

// NewRunInputWithResult creates a new RunInput with a value in the "result" field
func NewRunInputWithResult(jobRunID *ID, value interface{}, meta JSON, status RunStatus) *RunInput {
	data, err := JSON{}.Add("result", value)
	if err != nil {
		panic(fmt.Sprintf("invariant violated, add should not fail on empty JSON %v", err))
	}
	return &RunInput{
		jobRunID: *jobRunID,
		data:     data,
		Meta:     meta,
		status:   status,
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
