package migration1580827546

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the polling_interval duration (nanoseconds) to support the Flux Monitor.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE job_runs RENAME COLUMN overrides TO initial_params;
	`).Error
	return nil
}
