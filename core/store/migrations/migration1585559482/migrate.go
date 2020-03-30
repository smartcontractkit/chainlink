package migration1585559482

import (
	"github.com/jinzhu/gorm"
)

// Split pending_confirmations into pending_outgoing_confirmations and pending_incoming_confirmations
// For both task_runs and job_runs
// This is inherently ambiguous (that's the entire reason for splitting them out) but we can make a best
// guess by checking for presence of a transaction.
func Migrate(tx *gorm.DB) error {
	// HACK: Assume it's pending outgoing transactions if there exists a TX for the job run owning this task run
	err := tx.Exec(`
		UPDATE task_runs
		SET status = 'pending_outgoing_confirmations'
		FROM txes
		WHERE txes.surrogate_id::uuid = task_runs.job_run_id
		AND status = 'pending_confirmations'
	`).Error
	if err != nil {
		return err
	}
	// All remaining must be pending incoming confirmations
	err = tx.Exec(`
		UPDATE task_runs
		SET status = 'pending_incoming_confirmations'
		WHERE status = 'pending_confirmations'
	`).Error
	if err != nil {
		return err
	}

	// HACK: Assume it's pending outgoing transactions if there exists a TX for the job run
	err = tx.Exec(`
		UPDATE job_runs
		SET status = 'pending_outgoing_confirmations'
		FROM txes
		WHERE txes.surrogate_id::uuid = job_runs.id
		AND status = 'pending_confirmations'
	`).Error
	if err != nil {
		return err
	}

	// All remaining must be pending incoming confirmations
	return tx.Exec(`
		UPDATE task_runs
		SET status = 'pending_incoming_confirmations'
		WHERE status = 'pending_confirmations'
	`).Error
}
