package migration1587975059

import (
	"github.com/jinzhu/gorm"
)

// Migrate drops the LogCursor table
// This is already out in the wild as of 0.8.2 so we cannot simply delete the old migration
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS log_cursors`).Error
}
