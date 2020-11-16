package migration1605218542

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE offchainreporting_oracle_specs DROP CONSTRAINT offchainreporting_oracle_specs_encrypted_ocr_key_bundle_id_fkey;
        ALTER TABLE offchainreporting_oracle_specs DROP CONSTRAINT offchainreporting_oracle_specs_transmitter_address_fkey;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN transmitter_address DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN encrypted_ocr_key_bundle_id DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ADD CONSTRAINT encrypted_ocr_key_bundle_id_not_null CHECK((is_bootstrap_peer) OR (encrypted_ocr_key_bundle_id IS NOT NULL));
    `).Error
}
