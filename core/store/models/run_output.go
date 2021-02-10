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
// contains a result and preserves the resultCollection.
func NewRunOutputCompleteWithResult(resultVal interface{}, resultCollection gjson.Result) RunOutput {
	data, err := JSON{}.Add(ResultKey, resultVal)
	if err != nil {
		panic(fmt.Sprintf("invariant violated, add should not fail on empty JSON %v", err))
	}
	if resultCollection.String() != "" {
		collectionCopy := make([]interface{}, 0)
		for _, k := range resultCollection.Array() {
			collectionCopy = append(collectionCopy, k.Value())
		}
		data, err = data.Add(ResultCollectionKey, collectionCopy)
		if err != nil {
			return NewRunOutputError(err)
		}
	}
	return NewRunOutputComplete(data)
}

// NewRunOutputComplete returns a new RunOutput that is complete and contains
// raw data
func NewRunOutputComplete(data JSON) RunOutput {
	return RunOutput{status: RunStatusCompleted, data: data}
}

// NewRunOutputPendingOutgoingConfirmationsWithData returns a new RunOutput that
// indicates the task is pending outgoing confirmations but also has some data that
// needs to be fed in on next invocation
func NewRunOutputPendingOutgoingConfirmationsWithData(data JSON) RunOutput {
	return RunOutput{status: RunStatusPendingOutgoingConfirmations, data: data}
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

func (ro RunOutput) ResultCollection() gjson.Result {
	return ro.data.Get(ResultCollectionKey)
}

// Result returns the result as a gjson object
func (ro RunOutput) Result() gjson.Result {
	return ro.data.Get(ResultKey)
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
