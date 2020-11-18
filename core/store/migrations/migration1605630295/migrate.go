package migration1605630295

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE offchainreporting_oracle_specs ADD CONSTRAINT offchainreporting_oracle_specs_encrypted_ocr_key_bundle_id_fkey FOREIGN KEY(encrypted_ocr_key_bundle_id) REFERENCES encrypted_ocr_key_bundles(id);
        ALTER TABLE offchainreporting_oracle_specs ADD CONSTRAINT offchainreporting_oracle_specs_transmitter_address_fkey FOREIGN KEY(transmitter_address) REFERENCES keys(address);
    `).Error
}
