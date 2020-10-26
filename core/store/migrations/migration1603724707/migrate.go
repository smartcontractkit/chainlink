package migration1603724707

import "github.com/jinzhu/gorm"

const up = `
ALTER TABLE encrypted_ocr_key_bundles ADD COLUMN config_public_key bytea NOT NULL;
`

// Migrate adds config_public_key to encrypted_ocr_key_bundles
// This migration will fail if there are any keys already.
// We do not automatically delete them, but instead will require the node operator to do that manually.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}
