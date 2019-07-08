package migration1562623854

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(`
ALTER TABLE tx_attempts ADD COLUMN status varchar(255) NOT NULL DEFAULT 'unconfirmed';
UPDATE tx_attempts SET status = 'confirmed' WHERE confirmed IS TRUE;`).Error; err != nil {
		return errors.Wrap(err, "failed to add status column to TxAttempts")
	}
	return nil
}
