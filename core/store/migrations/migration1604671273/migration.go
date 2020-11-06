package migration1604671273

import "github.com/jinzhu/gorm"

// Migrate removes an old column that was never used
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE pipeline_runs DROP COLUMN finished_at;
    `).Error
}
