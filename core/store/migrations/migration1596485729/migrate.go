package migration1596485729

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds a job_run_id column to flux_monitor_round_stats
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE flux_monitor_round_stats ADD COLUMN "job_run_id" uuid REFERENCES job_runs(id) ON DELETE CASCADE;
	`).Error
}
