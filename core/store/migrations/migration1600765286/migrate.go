package migration1600765286

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates necessary tables for OCR database implementation
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		-- NOTE: Placeholder table necessary for foreign keys, more columns will be added later
		CREATE TABLE offchainreporting_oracle_specs (
			id SERIAL PRIMARY KEY
		);

		CREATE TABLE offchainreporting_persistent_states (
			offchainreporting_oracle_spec_id INT NOT NULL REFERENCES offchainreporting_oracle_specs (id),
			config_digest bytea NOT NULL CHECK (octet_length(config_digest) = 16),
			epoch bigint NOT NULL,
			highest_sent_epoch bigint NOT NULL,
			highest_received_epoch bigint[] NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL,
			PRIMARY KEY("offchainreporting_oracle_spec_id", "config_digest")
	   );

		CREATE TABLE offchainreporting_contract_configs (
			offchainreporting_oracle_spec_id INT NOT NULL REFERENCES offchainreporting_oracle_specs (id),
			config_digest bytea NOT NULL CHECK (octet_length(config_digest) = 16),
			signers bytea[],
			transmitters bytea[],
			threshold integer,
			encoded_config_version bigint,
			encoded bytea,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL,
			PRIMARY KEY(offchainreporting_oracle_spec_id)
		);

	   CREATE TABLE offchainreporting_pending_transmissions (
			offchainreporting_oracle_spec_id INT NOT NULL REFERENCES offchainreporting_oracle_specs (id),
			config_digest bytea NOT NULL CHECK (octet_length(config_digest) = 16),
			epoch bigint NOT NULL,
			round bigint NOT NULL,
			time timestamptz NOT NULL, 
			median numeric(78,0) NOT NULL,
			serialized_report bytea NOT NULL,
			rs bytea[] NOT NULL,
			ss bytea[] NOT NULL,
			vs bytea NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL,
			PRIMARY KEY(offchainreporting_oracle_spec_id, config_digest, epoch, round)
		);

		CREATE INDEX idx_offchainreporting_pending_transmissions_time ON offchainreporting_pending_transmissions (time);
	  `).Error

}
