package models

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// RunInput represents the input for performing a Task
type RunInput struct {
	JobRunID     ID
	Data         JSON
	Status       RunStatus
	ErrorMessage null.String
}

// Result returns the result as a gjson object
func (ri RunInput) Result() gjson.Result {
	return ri.Data.Get("result")
}

// ResultString returns the string result of the Data JSON field.
func (ri RunInput) ResultString() (string, error) {
	val := ri.Result()
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// HasError returns true if the status is errored or the error message is set
func (ri RunInput) HasError() bool {
	return ri.Status == RunStatusErrored || ri.ErrorMessage.Valid
}

// GetError returns the error of a RunResult if it is present.
func (ri RunInput) GetError() error {
	if ri.HasError() {
		return errors.New(ri.ErrorMessage.ValueOrZero())
	}
	return nil
}

// SetError marks the result as errored and saves the specified error message
func (ri *RunInput) SetError(err error) {
	ri.ErrorMessage = null.StringFrom(err.Error())
	ri.Status = RunStatusErrored
}
