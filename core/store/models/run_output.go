package models

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// RunOutput represents the result of performing a Task
type RunOutput struct {
	data   JSON
	status RunStatus
	err    error
}

// NewRunOutputError returns a new RunOutput with an error
func NewRunOutputError(err error) RunOutput {
	return RunOutput{
		status: RunStatusErrored,
		err:    err,
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
	return RunOutput{status: RunStatusCompleted, data: data}
}

// NewRunOutputPendingSleep returns a new RunOutput that indicates the task is
// sleeping
func NewRunOutputPendingSleep() RunOutput {
	return RunOutput{status: RunStatusPendingSleep}
}

// NewRunOutputPendingConfirmations returns a new RunOutput that indicates the
// task is pending confirmations
func NewRunOutputPendingConfirmations() RunOutput {
	return RunOutput{status: RunStatusPendingConfirmations}
}

// NewRunOutputPendingConfirmationsWithData returns a new RunOutput that
// indicates the task is pending confirmations but also has some data that
// needs to be fed in on next invocation
func NewRunOutputPendingConfirmationsWithData(data JSON) RunOutput {
	return RunOutput{status: RunStatusPendingConfirmations, data: data}
}

// NewRunOutputPendingConnection returns a new RunOutput that indicates the
// task got disconnected
func NewRunOutputPendingConnection() RunOutput {
	return RunOutput{status: RunStatusPendingConnection}
}

// NewRunOutputPendingConnectionWithData returns a new RunOutput that
// indicates the task got disconnected but also has some data that needs to be
// fed in on next invocation
func NewRunOutputPendingConnectionWithData(data JSON) RunOutput {
	return RunOutput{status: RunStatusPendingConnection, data: data}
}

// NewRunOutputInProgress returns a new RunOutput that indicates the
// task is still in progress
func NewRunOutputInProgress(data JSON) RunOutput {
	return RunOutput{status: RunStatusInProgress, data: data}
}

// NewRunOutputPendingBridge returns a new RunOutput that indicates the
// task is still in progress
func NewRunOutputPendingBridge() RunOutput {
	return RunOutput{status: RunStatusPendingBridge}
}

// HasError returns true if the status is errored or the error message is set
func (ro RunOutput) HasError() bool {
	return ro.status == RunStatusErrored
}

// Result returns the result as a gjson object
func (ro RunOutput) Result() gjson.Result {
	return ro.data.Get("result")
}

// ResultString returns the "result" value as a string if possible
func (ro RunOutput) ResultString() (string, error) {
	val := ro.Result()
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string result")
	}
	return val.String(), nil
}

// Get searches for and returns the JSON at the given path.
func (ro RunOutput) Get(path string) gjson.Result {
	return ro.data.Get(path)
}

// Error returns error for this RunOutput
func (ro RunOutput) Error() error {
	return ro.err
}

// Data returns the data held by this RunOutput
func (ro RunOutput) Data() JSON {
	return ro.data
}

// Status returns the status returned from a task
func (ro RunOutput) Status() RunStatus {
	return ro.status
}
