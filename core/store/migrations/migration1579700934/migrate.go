package migration1579700934

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the polling_interval duration (nanoseconds) to support the Flux Monitor.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE initiators ADD COLUMN "polling_interval" BigInt;
	`).Error
}
