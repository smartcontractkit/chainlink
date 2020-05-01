package migration1588293486

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE initiators ADD COLUMN "poll_timer" jsonb;
	  ALTER TABLE initiators ADD COLUMN "idle_timer" jsonb;
	  ALTER TABLE initiators DROP COLUMN "idle_threshold";
	`).Error
}
