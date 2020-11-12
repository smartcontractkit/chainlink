package migration1605213161

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        DROP INDEX IF EXISTS idx_offchainreporting_oracle_specs_unique_key_bundles;
        DROP INDEX IF EXISTS idx_offchainreporting_oracle_specs_unique_peer_ids;
        ALTER TABLE offchainreporting_oracle_specs ADD CONSTRAINT unique_contract_addr UNIQUE (contract_address);
    `).Error
}
