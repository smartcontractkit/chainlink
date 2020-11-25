package migration1606320711

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE jobs ADD COLUMN max_task_duration bigint;
    `).Error
}
