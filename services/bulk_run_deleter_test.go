package services_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteJobRuns(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, initiator := cltest.NewJobWithWebInitiator()

	// matches updated before but none of the statuses
	oldIncompleteRun := job.NewRun(initiator)
	oldIncompleteRun.Status = models.RunStatusInProgress
	oldIncompleteRun.UpdatedAt = cltest.ParseISO8601("2018-01-01T00:00:00Z")
	err := store.ORM.DB.Save(&oldIncompleteRun)
	assert.NoError(t, err)

	// matches one of the statuses and the updated before
	oldCompletedRun := job.NewRun(initiator)
	oldCompletedRun.Status = models.RunStatusCompleted
	oldCompletedRun.UpdatedAt = cltest.ParseISO8601("2018-01-01T00:00:00Z")
	err = store.ORM.DB.Save(&oldCompletedRun)
	assert.NoError(t, err)

	// matches one of the statuses but not the updated before
	newCompletedRun := job.NewRun(initiator)
	newCompletedRun.Status = models.RunStatusCompleted
	newCompletedRun.UpdatedAt = cltest.ParseISO8601("2018-01-30T00:00:00Z")
	err = store.ORM.DB.Save(&newCompletedRun)
	assert.NoError(t, err)

	// matches nothing
	newIncompleteRun := job.NewRun(initiator)
	newIncompleteRun.Status = models.RunStatusCompleted
	newIncompleteRun.UpdatedAt = cltest.ParseISO8601("2018-01-30T00:00:00Z")
	err = store.ORM.DB.Save(&newIncompleteRun)
	assert.NoError(t, err)

	err = services.DeleteJobRuns(store.ORM, &models.BulkDeleteRunRequest{
		Status:        []models.RunStatus{models.RunStatusCompleted},
		UpdatedBefore: cltest.ParseISO8601("2018-01-15T00:00:00Z"),
	})

	runCount, err := store.ORM.DB.Count(&models.JobRun{})
	assert.NoError(t, err)
	assert.Equal(t, 3, runCount)
}

func TestRunPendingTask(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	task, err := models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{})
	assert.NoError(t, err)
	err = services.RunPendingTask(store.ORM, task)
	assert.NoError(t, err)
	assert.Equal(t, string(models.BulkTaskStatusCompleted), string(task.Status))
}

func TestRunPendingTask_Error(t *testing.T) {
	store, cleanup := cltest.NewStore()
	// Close store immediately to trigger error
	cleanup()

	task, err := models.NewBulkDeleteRunTask(models.BulkDeleteRunRequest{})
	assert.NoError(t, err)
	err = services.RunPendingTask(store.ORM, task)
	assert.Error(t, err)
	assert.Equal(t, string(models.BulkTaskStatusErrored), string(task.Status))
	assert.NotNil(t, task.Error)
}
