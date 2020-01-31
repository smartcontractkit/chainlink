package migration1580290179

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the 'meta' field to the run_results table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
ALTER TABLE job_runs ADD COLUMN initial_meta text NOT NULL DEFAULT '{}';
	`).Error
}
