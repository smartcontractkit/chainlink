package migration1566498796

import (
	"github.com/jinzhu/gorm"
)

// Migrate optimizes the JobRuns table to reduce the cost of IDs
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
ALTER TABLE run_results DROP COLUMN IF EXISTS cached_job_run_id;
ALTER TABLE run_results DROP COLUMN IF EXISTS cached_task_run_id;
	`).Error
}
