package migration1586342453

import (
	"github.com/jinzhu/gorm"
)

// Migrate removes some columns that are no longer in use
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	ALTER TABLE job_runs DROP COLUMN overrides_id;
	ALTER TABLE task_runs DROP COLUMN confirmations_old1560433987;
	ALTER TABLE initiators ALTER COLUMN topics TYPE jsonb USING topics::jsonb;
	`).Error
}
