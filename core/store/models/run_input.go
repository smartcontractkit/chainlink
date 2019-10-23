package models

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// RunInput represents the input for performing a Task
type RunInput struct {
	jobRunID ID
	data     JSON
	status   RunStatus
	err      error
}

// NewRunInput creates a new RunInput with arbitrary data
func NewRunInput(jobRunID *ID, data JSON, status RunStatus) *RunInput {
	return &RunInput{
		jobRunID: *jobRunID,
		data:     data,
		status:   status,
	}
}

// NewRunInputWithResult creates a new RunInput with a value in the "result" field
func NewRunInputWithResult(jobRunID *ID, value interface{}, status RunStatus) *RunInput {
	data, _ := JSON{}.Add("result", value)
	return &RunInput{
		jobRunID: *jobRunID,
		data:     data,
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

// GetError returns the error of a RunResult if it is present.
func (ri RunInput) Error() error {
	return ri.err
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
