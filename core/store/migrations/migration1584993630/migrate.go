package migration1584993630

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the InitiatorParams.RequestTimeout duration (nanoseconds) to support the Flux Monitor.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE initiators ADD COLUMN "request_timeout" BigInt;
	`).Error
}
