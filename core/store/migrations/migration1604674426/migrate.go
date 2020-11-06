package migration1604003825

import "github.com/jinzhu/gorm"

// Migrate makes key deletion into a soft delete.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE pipeline_runs
		ADD COLUMN errors jsonb,
		ADD COLUMN result jsonb,
		ADD CHECK ((result is null and errors is null and finished_at is null) or (result is not null and errors is not null and finished_at is not null));
    `).Error
}
