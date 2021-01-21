package migration1611174407

import "github.com/jinzhu/gorm"

// Migrate makes the explicit the pre-existing
// implicit assumption that lowercase external initiator names are unique
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE pipeline_task_runs 
		ADD COLUMN predecessor_task_run_ids integer[] NOT NULL; 
	`).Error
}
