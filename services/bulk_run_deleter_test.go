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

	db := store.ORM.DB
	job, initiator := cltest.NewJobWithWebInitiator()

	// matches updated before but none of the statuses
	oldIncompleteRun := job.NewRun(initiator)
	oldIncompleteRun.Status = models.RunStatusInProgress
	err := db.Save(&oldIncompleteRun).Error
	assert.NoError(t, err)
	db.Model(&oldIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-01T00:00:00Z"))

	// matches one of the statuses and the updated before
	oldCompletedRun := job.NewRun(initiator)
	oldCompletedRun.Status = models.RunStatusCompleted
	err = db.Save(&oldCompletedRun).Error
	assert.NoError(t, err)
	db.Model(&oldCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-01T00:00:00Z"))

	// matches one of the statuses but not the updated before
	newCompletedRun := job.NewRun(initiator)
	newCompletedRun.Status = models.RunStatusCompleted
	err = db.Save(&newCompletedRun).Error
	assert.NoError(t, err)
	db.Model(&newCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-30T00:00:00Z"))

	// matches nothing
	newIncompleteRun := job.NewRun(initiator)
	newIncompleteRun.Status = models.RunStatusCompleted
	err = db.Save(&newIncompleteRun).Error
	assert.NoError(t, err)
	db.Model(&newIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-30T00:00:00Z"))

	err = services.DeleteJobRuns(store.ORM, &models.BulkDeleteRunRequest{
		Status:        []models.RunStatus{models.RunStatusCompleted},
		UpdatedBefore: cltest.ParseISO8601("2018-01-15T00:00:00Z"),
	})

	assert.NoError(t, err)
	var runCount int
	err = db.Model(&models.JobRun{}).Count(&runCount).Error
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
	assert.NotNil(t, task.ErrorMessage)
}
