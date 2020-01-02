package migration1576022702

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE initiators ADD COLUMN "request_data" text;
		ALTER TABLE initiators ADD COLUMN "feeds" text;
		ALTER TABLE initiators ADD COLUMN "threshold" float;
		ALTER TABLE initiators ADD COLUMN "precision" smallint;
	`).Error
}
