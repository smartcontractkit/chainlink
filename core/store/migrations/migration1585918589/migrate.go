package migration1585918589

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates and optimizes table indexes
func Migrate(tx *gorm.DB) error {
	// I can't believe we were missing this one
	// Need `if not exists` because I created it manually on the kovan util node
	err := tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_task_runs_job_run_id ON task_runs(job_run_id);
	`).Error
	if err != nil {
		return err
	}

	err = tx.Exec(`
	CREATE INDEX idx_job_runs_job_spec_id ON job_runs(job_spec_id);
	`).Error
	if err != nil {
		return err
	}

	// FIXME: This should ideally be unique but there exists non-unique data out in the wild so need to handle that...
	// Same for tx_attempts
	err = tx.Exec(`
	CREATE INDEX idx_txs_hash ON txes(hash);
	`).Error
	if err != nil {
		return err
	}

	// The majority of runs are completed so there is no point in indexing those ones.
	// We can reduce the size of the index by excluding this status
	err = tx.Exec(`
	DROP INDEX idx_job_runs_status;
	CREATE INDEX idx_job_runs_status ON job_runs(status) WHERE status != 'completed';
	`).Error
	if err != nil {
		return err
	}

	// Brin indexes offer much more efficient storage for time series data on large tables
	return tx.Exec(`
	DROP INDEX idx_task_runs_created_at;
	CREATE INDEX idx_task_runs_created_at ON task_runs USING BRIN (created_at);

	DROP INDEX idx_job_runs_created_at;
	CREATE INDEX idx_job_runs_created_at ON job_runs USING BRIN (created_at);

	CREATE INDEX idx_job_runs_updated_at ON job_runs USING BRIN (updated_at);

	CREATE INDEX idx_job_runs_finished_at ON job_runs USING BRIN (finished_at);

	DROP INDEX idx_sessions_last_used;
	CREATE INDEX idx_sessions_last_used ON sessions USING BRIN (last_used);

	DROP INDEX idx_sessions_created_at;
	CREATE INDEX idx_sessions_created_at ON sessions USING BRIN (created_at);

	DROP INDEX idx_tx_attempts_created_at;
	CREATE INDEX idx_tx_attempts_created_at ON tx_attempts USING BRIN (created_at);

	CREATE INDEX idx_run_requests_created_at ON run_requests USING BRIN (created_at);

	CREATE INDEX idx_task_specs_created_at ON task_specs USING BRIN (created_at);
	CREATE INDEX idx_task_specs_updated_at ON task_specs USING BRIN (updated_at);
	`).Error
}
