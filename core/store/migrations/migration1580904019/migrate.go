package migration1580904019

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the InitiatorParams.IdleThreshold duration (nanoseconds) to support the Flux Monitor.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE initiators ADD COLUMN "idle_threshold" BigInt;
	`).Error
	return nil
}
