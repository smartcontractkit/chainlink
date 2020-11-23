package migration1604707007

import (
	"github.com/jinzhu/gorm"
)

// Migrate drops the legacy txs and tx_attempts tables
func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`DROP TABLE IF EXISTS tx_attempts`).Error
	if err != nil {
		return err
	}
	return tx.Exec(`DROP TABLE IF EXISTS txes`).Error
}
