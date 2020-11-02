package migration1604003825

import "github.com/jinzhu/gorm"

// Migrate makes key deletion into a soft delete.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE keys                      ADD COLUMN deleted_at timestamptz;
        ALTER TABLE encrypted_vrf_keys        ADD COLUMN deleted_at timestamptz;
        ALTER TABLE encrypted_p2p_keys        ADD COLUMN deleted_at timestamptz;
        ALTER TABLE encrypted_ocr_key_bundles ADD COLUMN deleted_at timestamptz;
    `).Error
}
