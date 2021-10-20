package postgres

import (
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

// BatchSize is the default number of DB records to access in one batch
const BatchSize uint = 1000

// BatchFunc is the function to execute on each batch of records, should return the count of records affected
type BatchFunc func(offset, limit uint) (count uint, err error)

// Batch is an iterator for batches of records
func Batch(cb BatchFunc) error {
	offset := uint(0)
	limit := BatchSize

	for {
		count, err := cb(offset, limit)
		if err != nil {
			return err
		}

		if count < limit {
			return nil
		}

		offset += limit
	}
}

// BulkDeleteRuns removes JobRuns and their related records: TaskRuns and
// RunResults.
//
// RunResults and RunRequests are pointed at by JobRuns so we must use two CTEs
// to remove both parents in one hit.
//
// TaskRuns are removed by ON DELETE CASCADE when the JobRuns and RunResults
// are deleted.
func BulkDeleteRuns(db *gorm.DB, bulkQuery *models.BulkDeleteRunRequest) error {
	return Batch(func(_, limit uint) (count uint, err error) {
		res := db.Exec(`
WITH job_runs_to_delete AS (
	SELECT id FROM job_runs jr
	WHERE jr.status IN (?) AND jr.updated_at < ?
	ORDER BY jr.id ASC
	LIMIT ?
),
deleted_task_runs AS (
	DELETE FROM task_runs as tr
	WHERE tr.job_run_id IN (SELECT id FROM job_runs_to_delete)
	RETURNING tr.result_id
),
deleted_task_run_results AS (
	DELETE FROM run_results WHERE id IN (SELECT result_id FROM deleted_task_runs)
),
deleted_job_runs AS (
	DELETE FROM job_runs as jr
	WHERE jr.id IN (SELECT id FROM job_runs_to_delete)
	RETURNING jr.result_id, jr.run_request_id
),
deleted_job_run_results AS (
	DELETE FROM run_results WHERE id IN (SELECT result_id FROM deleted_job_runs)
)
DELETE FROM run_requests WHERE id IN (SELECT run_request_id FROM deleted_job_runs)
;
		`, bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore, limit)

		if res.Error != nil {
			return count, res.Error
		}
		return uint(res.RowsAffected), res.Error
	})
}

// Sessions returns all sessions limited by the parameters.
func Sessions(db *gorm.DB, offset, limit int) ([]models.Session, error) {
	var sessions []models.Session
	err := db.
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}
