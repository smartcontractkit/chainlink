package migration1588757164

import (
	"github.com/jinzhu/gorm"
)

// Split pending_confirmations into pending_outgoing_confirmations and pending_incoming_confirmations
// For both task_runs and job_runs
// This is inherently ambiguous (that's the entire reason for splitting them out) but we can make a best
// guess by checking for presence of a transaction.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		-- Drop old partial index and defaults which will otherwise cause cast to fail
		DROP INDEX idx_job_runs_status;
		DROP INDEX idx_task_runs_status;
		ALTER TABLE job_runs ALTER COLUMN status SET DEFAULT NULL;
		ALTER TABLE task_runs ALTER COLUMN status SET DEFAULT NULL;

		-- Convert column type to text to allow for easy transformation
		ALTER TABLE task_runs ALTER COLUMN status TYPE text;
		ALTER TABLE job_runs ALTER COLUMN status TYPE text;

		-- HACK: Assume it's pending outgoing transactions if there exists a TX for the job run owning this task run
		UPDATE task_runs
		SET status = 'pending_outgoing_confirmations'
		FROM txes
		WHERE txes.surrogate_id::uuid = task_runs.job_run_id
		AND status = 'pending_confirmations';

		-- All remaining must be pending incoming confirmations
		UPDATE task_runs
		SET status = 'pending_incoming_confirmations'
		WHERE status = 'pending_confirmations';

		-- HACK: Assume it's pending outgoing transactions if there exists a TX for the job run
		UPDATE job_runs
		SET status = 'pending_outgoing_confirmations'
		FROM txes
		WHERE txes.surrogate_id::uuid = job_runs.id
		AND status = 'pending_confirmations';

		-- All remaining must be pending incoming confirmations
		UPDATE task_runs
		SET status = 'pending_incoming_confirmations'
		WHERE status = 'pending_confirmations';

		-- Create the new enum type
		CREATE TYPE run_status AS ENUM ('unstarted', 'in_progress', 'pending_incoming_confirmations', 'pending_outgoing_confirmations', 'pending_connection', 'pending_bridge', 'pending_sleep', 'errored', 'completed', 'cancelled');

		-- Cast the columns
		ALTER TABLE job_runs ALTER COLUMN status TYPE run_status USING status::run_status;	
		ALTER TABLE task_runs ALTER COLUMN status TYPE run_status USING status::run_status;	

		-- Drop the old types
		DROP TYPE job_run_status;
		DROP TYPE task_run_status;

		-- Recreate indexes and defaults
		CREATE INDEX idx_job_runs_status ON job_runs(status) WHERE status != 'completed'::run_status;
		CREATE INDEX idx_task_runs_status ON task_runs(status) WHERE status != 'completed'::run_status;
		ALTER TABLE job_runs ALTER COLUMN status SET DEFAULT 'unstarted';
		ALTER TABLE task_runs ALTER COLUMN status SET DEFAULT 'unstarted';
	`).Error
}
