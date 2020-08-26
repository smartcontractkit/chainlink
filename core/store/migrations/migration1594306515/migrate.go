package migration1594306515

import (
	"github.com/jinzhu/gorm"
)

// Migrate ensures that heads are unique and adds parent hash for use in reorg detection
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE keys ALTER COLUMN next_nonce DROP NOT NULL;
		ALTER TABLE keys ALTER COLUMN next_nonce SET DEFAULT NULL;
	`).Error
}
