package migration1584153740

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the "failed" field to the txes table so that we stop resubmitting a tx in error cases.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE txes ADD COLUMN "failed" BOOL NOT NULL DEFAULT FALSE;
	`).Error
}
