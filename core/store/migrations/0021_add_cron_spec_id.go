package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up21 = `
		CREATE TABLE cron_specs (
			id BIGSERIAL PRIMARY KEY,
			cron_schedule text NOT NULL,
			to_address bytea NOT NULL,
			oracle_payment numeric(78, 0),
			eth_gas_limit integer,
			created_at timestamp with time zone NOT NULL,
			updated_at timestamp with time zone NOT NULL,
			on_chain_job_spec_id bytea NOT NULL,
			CONSTRAINT direct_request_specs_on_chain_job_spec_id_check CHECK ((octet_length(on_chain_job_spec_id) = 32)),
			CONSTRAINT cron_specs_from_address_check CHECK ((octet_length(to_address) = 20))
		);

		ALTER TABLE jobs ADD COLUMN cron_request_spec_id INT REFERENCES cron_specs(id),
		DROP CONSTRAINT chk_only_one_spec,
		ADD CONSTRAINT chk_only_one_spec CHECK (
			num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_request_spec_id) = 1
		);
	`

const down21 = `
		DROP TABLE IF EXISTS cron_specs;
	
		ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
		ADD CONSTRAINT chk_only_one_spec CHECK (
			num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id) = 1
		);
	
		ALTER TABLE jobs DROP COLUMN cron_request_spec_id integer;
	`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0021_add_cron_spec_tables",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up21).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down21).Error
		},
	})
}
