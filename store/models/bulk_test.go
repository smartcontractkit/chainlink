package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulkDeleteRunTask_NewBulkDeleteRunTask(t *testing.T) {
	task, err := models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	task, err = models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{
		Status: []models.RunStatus{models.RunStatusCompleted},
	})
	assert.NoError(t, err)

	task, err = models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{
		Status: []models.RunStatus{""},
	})
	assert.Error(t, err)

	task, err = models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{
		Status: []models.RunStatus{models.RunStatusInProgress},
	})
	assert.Error(t, err)
}

func TestBulkDeleteRunTask_RetrieveFromDB(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	task, err := models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{
		Status: []models.RunStatus{models.RunStatusCompleted, models.RunStatusErrored},
	})
	require.NoError(t, err)

	task.ErrorMessage = "got an error"
	require.NoError(t, store.SaveBulkDeleteRunTask(task))

	retrievedTask, err := store.FindBulkDeleteRunTask(task.ID)
	require.NoError(t, err)
	task.CreatedAt = retrievedTask.CreatedAt // ignore CreatedAt assertion
	assert.Equal(t, task, retrievedTask)
}
