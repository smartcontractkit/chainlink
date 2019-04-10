package migration1554855314

import (
	"github.com/jinzhu/gorm"
)

// Migration is the singleton type for this migration
type Migration struct{}

// Migrate adds the sync_events table
func (m Migration) Migrate(tx *gorm.DB) error {
	tx = tx.Begin()
	if err := tx.Exec(`
ALTER TABLE "bridge_types" ADD COLUMN "incoming_token_hash" VARCHAR(32);
ALTER TABLE "bridge_types" ADD COLUMN "salt" VARCHAR(32);
`).Error; err != nil {
		tx.Rollback()
		return err
	}

	// TODO: migrate passwords

	//if err := tx.Exec(`
	//ALTER TABLE "bridge_types" DROP COLUMN "incoming_token";
	//`).Error; err != nil {
	//tx.Rollback()
	//return err
	//}

	return tx.Commit().Error
}
