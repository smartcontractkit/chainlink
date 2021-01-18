package migration1609963213

import "github.com/jinzhu/gorm"

// Migrate renames eth_request_event_specs to direct_request_specs and adds on_chain_job_spec_id
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE flux_monitor_specs (
			id SERIAL PRIMARY KEY,
			contract_address bytea NOT NULL CHECK (octet_length(contract_address) = 20),
			precision integer,
			threshold real, 
			absolute_threshold real,
			poll_timer_period bigint,
			poll_timer_disabled boolean,
			CHECK (poll_timer_disabled OR poll_timer_period > 0),
			idle_timer_period bigint,
			idle_timer_disabled boolean,
			CHECK (idle_timer_disabled OR idle_timer_period > 0),
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL
		);

		ALTER TABLE jobs ADD COLUMN flux_monitor_spec_id INT REFERENCES flux_monitor_specs(id),
			DROP CONSTRAINT chk_only_one_spec,
			ADD CONSTRAINT chk_only_one_spec CHECK (
				num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id) = 1
		);
	`).Error
}
