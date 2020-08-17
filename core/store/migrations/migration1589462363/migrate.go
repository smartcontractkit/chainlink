package migration1589462363

import (
	"github.com/jinzhu/gorm"
)

// Migrate updates the keys table for the BulletproofTxManager
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE keys ADD COLUMN next_nonce BIGINT NOT NULL DEFAULT 0;
		CREATE UNIQUE INDEX idx_unique_keys_address ON keys (address);
		ALTER TABLE keys ADD COLUMN id SERIAL NOT NULL;
		ALTER TABLE keys ADD CONSTRAINT chk_address_length CHECK (
			octet_length(address) = 20
		);
	`).Error
}
