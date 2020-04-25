package migration1587591248

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`
ALTER TABLE initiators ADD COLUMN "value_triggers" jsonb;
ALTER TABLE initiators DROP COLUMN "threshold"
	`).Error
	return errors.Wrapf(err, "while migrating initiators to include "+
		"value_triggers column")
	// XXX: Add actual migration, here
}
