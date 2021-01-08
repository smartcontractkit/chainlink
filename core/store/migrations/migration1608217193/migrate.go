package migration1608217193

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE pipeline_specs ADD COLUMN max_task_duration bigint;
    `).Error
}
