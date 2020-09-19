package migration1599691818

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the encrypted_ocr_key_bundles table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE encrypted_ocr_key_bundles (
			id bytea PRIMARY KEY,
			on_chain_signing_address bytea NOT NULL,
			off_chain_public_key bytea NOT NULL,
			encrypted_private_keys jsonb NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL
			);
			`).Error
}
