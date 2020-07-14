package migration1594642891

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds last_used to keys
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE keys ADD COLUMN last_used timestamptz;
	`).Error
}
