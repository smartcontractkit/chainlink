package migration1598005829

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  ALTER TABLE initiators ADD COLUMN "irita_service_name" varchar(255);
	  ALTER TABLE initiators ADD COLUMN "irita_service_provider" varchar(255);
	`).Error
}
