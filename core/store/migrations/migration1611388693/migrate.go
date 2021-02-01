package migration1611388693

import "github.com/jinzhu/gorm"

// Migrate adds the proper index to optimise the pipeline ORM for locking on runs instead of task runs
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE INDEX IF NOT EXISTS idx_pipeline_runs_unfinished_runs ON pipeline_runs (id) WHERE finished_at IS NULL;
		DROP INDEX IF EXISTS idx_pipeline_task_runs_optimise_find_predecessor_unfinished_runs;
		DROP INDEX IF EXISTS idx_pipeline_task_runs_unfinished;
	`).Error
}
