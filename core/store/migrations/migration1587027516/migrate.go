package migration1587027516

import (
	"github.com/jinzhu/gorm"
)

// Migrate converts text columns to enums where appropriate
func Migrate(tx *gorm.DB) error {
	// Create two different types for job runs/task runs
	// These are currently identical because right now they are represented by
	// one type in the code but may have different meanings in future and
	// arguably should be split up
	return tx.Exec(`
	CREATE TYPE job_run_status AS ENUM ('unstarted', 'in_progress', 'pending_confirmations', 'pending_connection', 'pending_bridge', 'pending_sleep', 'errored', 'completed', 'cancelled');
	CREATE TYPE task_run_status AS ENUM ('unstarted', 'in_progress', 'pending_confirmations', 'pending_connection', 'pending_bridge', 'pending_sleep', 'errored', 'completed', 'cancelled');

	-- It's no longer used as of deb84dbfc
	ALTER TABLE run_results DROP COLUMN status;

	UPDATE job_runs SET status = 'unstarted' WHERE status = '' OR status IS NULL;
	UPDATE task_runs SET status = 'unstarted' WHERE status = '' OR status IS NULL;

	DROP INDEX idx_job_runs_status;
	ALTER TABLE job_runs ALTER COLUMN status TYPE job_run_status USING status::job_run_status, ALTER COLUMN status SET DEFAULT 'unstarted'::job_run_status, ALTER COLUMN status SET NOT NULL;
	CREATE INDEX idx_job_runs_status ON job_runs (status) WHERE status != 'completed'::job_run_status;

	ALTER TABLE task_runs ALTER COLUMN status TYPE task_run_status USING status::task_run_status, ALTER COLUMN status SET DEFAULT 'unstarted'::task_run_status, ALTER COLUMN status SET NOT NULL;
	CREATE INDEX idx_task_runs_status ON task_runs (status) WHERE status != 'completed'::task_run_status;
	`).Error
}
