package models

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/utils"
)

// BulkTaskStatus indicates what a bulk task is doing.
type BulkTaskStatus string

const (
	// BulkTaskStatusInProgress is the default state of any run status.
	BulkTaskStatusInProgress = BulkTaskStatus("")
	// BulkTaskStatusErrored means a bulk task stopped because it encountered an error.
	BulkTaskStatusErrored = BulkTaskStatus("errored")
	// BulkTaskStatusCompleted means a bulk task finished.
	BulkTaskStatusCompleted = BulkTaskStatus("completed")
)

// BulkDeleteRunRequest describes the query for deletion of runs
type BulkDeleteRunRequest struct {
	Status        []RunStatus `json:"status"`
	UpdatedBefore time.Time   `json:"updatedBefore"`
}

// BulkDeleteRunTask represents a task that is working to delete runs with a query
type BulkDeleteRunTask struct {
	ID     string               `json:"id" storm:"id,unique"`
	Query  BulkDeleteRunRequest `json:"query"`
	Status BulkTaskStatus       `json:"status"`
	Error  error                `json:"error"`
}

// NewBulkDeleteRunTask returns a task from a request to make a task
func NewBulkDeleteRunTask(request BulkDeleteRunRequest) (*BulkDeleteRunTask, error) {
	for _, status := range request.Status {
		if status != RunStatusCompleted && status != RunStatusErrored {
			return nil, fmt.Errorf("cannot delete Runs with status %s", status)
		}
	}

	return &BulkDeleteRunTask{
		ID:    utils.NewBytes32ID(),
		Query: request,
	}, nil
}

// GetID returns the ID of this structure for jsonapi serialization.
func (t BulkDeleteRunTask) GetID() string {
	return t.ID
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (t BulkDeleteRunTask) GetName() string {
	return "bulk_delete_runs_tasks"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (t *BulkDeleteRunTask) SetID(value string) error {
	t.ID = value
	return nil
}
