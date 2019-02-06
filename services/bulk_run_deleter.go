package services

import (
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"go.uber.org/multierr"
)

// NewBulkRunDeleter creates a task runner that is responsible for executing bulk tasks.
func NewBulkRunDeleter(store *store.Store) SleeperTask {
	return NewSleeperTask(&bulkRunDeleter{
		store: store,
	})
}

type bulkRunDeleter struct {
	store *store.Store
}

func (btr *bulkRunDeleter) Work() {
	deleteTasks, err := btr.store.BulkDeletesInProgress()
	if err != nil {
		logger.Errorw("Error querying bulk tasks", "error", err)
		return
	}

	for _, task := range deleteTasks {
		err := RunPendingTask(btr.store.ORM, &task)
		if err != nil {
			logger.Errorw("Error deleting runs for bulk task", "task_id", task.ID, "error", err)
		}
	}
}

// RunPendingTask executes bulk run tasks
func RunPendingTask(orm *orm.ORM, task *models.BulkDeleteRunTask) error {
	logger.Infow("Processing bulk run delete task",
		"task_id", task.ID,
		"statuses", task.Query.Status,
		"updated_before", task.Query.UpdatedBefore,
	)
	err := DeleteJobRuns(orm, &task.Query)
	if err != nil {
		task.ErrorMessage = err.Error()
		task.Status = models.BulkTaskStatusErrored
	} else {
		task.Status = models.BulkTaskStatusCompleted
	}
	return multierr.Append(err, orm.DB.Save(task).Error)
}

// DeleteJobRuns removes runs that match a query
func DeleteJobRuns(orm *orm.ORM, bulkQuery *models.BulkDeleteRunRequest) error {
	return orm.DB.
		Where("status IN (?)", bulkQuery.Status.ToStrings()).
		Where("updated_at < ?", bulkQuery.UpdatedBefore).
		Delete(&[]models.JobRun{}).
		Error
}
