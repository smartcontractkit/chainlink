package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBulkDeleteRunTask(t *testing.T) {
	task, err := NewBulkDeleteRunTask(BulkDeleteRunRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	task, err = NewBulkDeleteRunTask(BulkDeleteRunRequest{Status: []RunStatus{RunStatusCompleted}})
	assert.NoError(t, err)

	task, err = NewBulkDeleteRunTask(BulkDeleteRunRequest{Status: []RunStatus{""}})
	assert.Error(t, err)

	task, err = NewBulkDeleteRunTask(BulkDeleteRunRequest{Status: []RunStatus{RunStatusInProgress}})
	assert.Error(t, err)
}
