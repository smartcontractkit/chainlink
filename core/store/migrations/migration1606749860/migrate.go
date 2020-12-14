package migration1606749860

import "github.com/jinzhu/gorm"

// Migrate adds the eth_request_event_spec_id to jobs table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE eth_request_event_specs (
			id SERIAL PRIMARY KEY,
			contract_address bytea NOT NULL CHECK (octet_length(contract_address) = 20),
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL
		);

		ALTER TABLE jobs ADD COLUMN eth_request_event_spec_id INT REFERENCES eth_request_event_specs (id),
			ALTER COLUMN offchainreporting_oracle_spec_id SET DEFAULT NULL,
			DROP CONSTRAINT chk_valid,
			ADD CONSTRAINT chk_only_one_spec CHECK (
				(offchainreporting_oracle_spec_id IS NOT NULL AND eth_request_event_spec_id IS NULL)
				OR
				(offchainreporting_oracle_spec_id IS NULL AND eth_request_event_spec_id IS NOT NULL)
		);
		
		CREATE UNIQUE INDEX idx_jobs_unique_eth_request_event_spec_id ON jobs (eth_request_event_spec_id);
    `).Error
}
