package models

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// RunOutput represents the result of performing a Task
type RunOutput struct {
	Data         JSON
	Status       RunStatus
	ErrorMessage null.String
}

// NewRunOutputError returns a new RunOutput with an error
func NewRunOutputError(err error) RunOutput {
	return RunOutput{
		Status:       RunStatusErrored,
		ErrorMessage: null.StringFrom(err.Error()),
	}
}

// NewRunOutputCompleteWithResult returns a new RunOutput that is complete and
// contains a result
func NewRunOutputCompleteWithResult(resultVal interface{}) RunOutput {
	var data JSON
	data, _ = data.Add("result", resultVal)
	return NewRunOutputComplete(data)
}

// NewRunOutputComplete returns a new RunOutput that is complete and contains
// raw data
func NewRunOutputComplete(data JSON) RunOutput {
	return RunOutput{
		Status: RunStatusCompleted,
		Data:   data,
	}
}

// NewRunOutputPendingSleep returns a new RunOutput that indicates the task is
// sleeping
func NewRunOutputPendingSleep() RunOutput {
	return RunOutput{
		Status: RunStatusPendingSleep,
	}
}

// NewRunOutputPendingConfirmations returns a new RunOutput that indicates the
// task is pending confirmations
func NewRunOutputPendingConfirmations() RunOutput {
	return RunOutput{
		Status: RunStatusPendingConfirmations,
	}
}

// NewRunOutputPendingConfirmationsWithData returns a new RunOutput that
// indicates the task is pending confirmations but also has some data that
// needs to be fed in on next invocation
func NewRunOutputPendingConfirmationsWithData(data JSON) RunOutput {
	return RunOutput{
		Status: RunStatusPendingConfirmations,
		Data:   data,
	}
}

// NewRunOutputPendingConnection returns a new RunOutput that indicates the
// task got disconnected
func NewRunOutputPendingConnection() RunOutput {
	return RunOutput{
		Status: RunStatusPendingConnection,
	}
}

// NewRunOutputPendingConfirmationsWithData returns a new RunOutput that
// indicates the task got disconnected but also has some data that needs to be
// fed in on next invocation
func NewRunOutputPendingConnectionWithData(data JSON) RunOutput {
	return RunOutput{
		Status: RunStatusPendingConnection,
		Data:   data,
	}
}

// NewRunOutputPendingConnection returns a new RunOutput that indicates the
// task is still in progress
func NewRunOutputInProgress(data JSON) RunOutput {
	return RunOutput{
		Status: RunStatusInProgress,
		Data:   data,
	}
}

// HasError returns true if the status is errored or the error message is set
func (ro RunOutput) HasError() bool {
	return ro.Status == RunStatusErrored || ro.ErrorMessage.Valid
}

// Result returns the result as a gjson object
func (ro RunOutput) Result() gjson.Result {
	return ro.Data.Get("result")
}

// GetError returns the error of a RunResult if it is present.
func (ro RunOutput) GetError() error {
	if ro.HasError() {
		return errors.New(ro.ErrorMessage.ValueOrZero())
	}
	return nil
}

// ResultString returns the string result of the Data JSON field.
func (ro RunOutput) ResultString() (string, error) {
	val := ro.Result()
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// Get searches for and returns the JSON at the given path.
func (ro RunOutput) Get(path string) gjson.Result {
	return ro.Data.Get(path)
}

// Error returns the string value of the ErrorMessage field.
func (ro RunOutput) Error() string {
	return ro.ErrorMessage.String
}
