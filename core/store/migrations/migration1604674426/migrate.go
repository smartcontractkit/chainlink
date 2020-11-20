package migration1604674426

import "github.com/jinzhu/gorm"

// Migrate makes key deletion into a soft delete.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE pipeline_runs
		ADD COLUMN errors jsonb,
		ADD COLUMN outputs jsonb,
		ADD CHECK ((outputs is null and errors is null and finished_at is null) or (outputs is not null and errors is not null and finished_at is not null));
    `).Error
}
