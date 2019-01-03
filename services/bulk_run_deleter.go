package services

import (
	"fmt"

	"github.com/asdine/storm/q"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
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
	query := btr.store.Select(q.Eq("Status", models.BulkTaskStatusInProgress)).OrderBy("CreatedAt")
	err := query.Each(&models.BulkDeleteRunTask{}, func(r interface{}) error {
		task := r.(*models.BulkDeleteRunTask)
		logger.Infow("Processing bulk run delete task",
			"task_id", task.ID,
			"statuses", task.Query.Status,
			"updated_before", task.Query.UpdatedBefore,
		)

		err := RunPendingTask(btr.store.ORM, task)
		if err != nil {
			logger.Errorw("Error deleting runs for bulk task", "task_id", task.ID, "error", err)
		}
		return err
	})
	if err != nil && err != orm.ErrorNotFound {
		logger.Errorw("Error querying bulk tasks", "error", err)
	}
}

// RunPendingTask executes bulk run tasks
func RunPendingTask(orm *orm.ORM, task *models.BulkDeleteRunTask) error {
	err := DeleteJobRuns(orm, &task.Query)
	if err != nil {
		task.Error = err
		task.Status = models.BulkTaskStatusErrored
	} else {
		task.Status = models.BulkTaskStatusCompleted
	}
	return orm.DB.Save(task)
}

// DeleteJobRuns removes runs that match a query
func DeleteJobRuns(orm *orm.ORM, bulkQuery *models.BulkDeleteRunRequest) error {
	query := orm.Select(
		q.And(
			q.In("Status", bulkQuery.Status),
			q.Lt("UpdatedAt", bulkQuery.UpdatedBefore),
		),
	).OrderBy("CompletedAt")

	return query.Each(&models.JobRun{}, func(r interface{}) error {
		run := r.(*models.JobRun)
		logger.Debugw("Deleting run", "run_id", run.ID, "status", run.Status, "updated_at", run.UpdatedAt)
		err := orm.DeleteStruct(run)
		if err != nil {
			return fmt.Errorf("error deleting run %s: %+v", run.ID, err)
		}
		return nil
	})
}
