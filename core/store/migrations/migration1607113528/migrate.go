package migration1607113528

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN p2p_peer_id DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN p2p_bootstrap_peers DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN observation_timeout DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN blockchain_timeout DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN contract_config_tracker_poll_interval DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN contract_config_tracker_subscribe_interval DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs ALTER COLUMN contract_config_confirmations DROP NOT NULL;
        ALTER TABLE offchainreporting_oracle_specs DROP CONSTRAINT encrypted_ocr_key_bundle_id_not_null;
    `).Error
}
