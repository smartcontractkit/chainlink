package migrations

import (
	"gorm.io/gorm"
)

const up27 = `
	CREATE TABLE webhook_specs (
		id SERIAL PRIMARY KEY,
        on_chain_job_spec_id bytea NOT NULL,
		created_at timestamp with time zone NOT NULL,
		updated_at timestamp with time zone NOT NULL
	);

	ALTER TABLE jobs ADD COLUMN webhook_spec_id INT REFERENCES webhook_specs(id),
	DROP CONSTRAINT chk_only_one_spec,
	ADD CONSTRAINT chk_only_one_spec CHECK (
		num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_spec_id, webhook_spec_id) = 1
	);
`

const down27 = `
	ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
	ADD CONSTRAINT chk_only_one_spec CHECK (
		num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_spec_id) = 1
	);

	ALTER TABLE jobs DROP COLUMN webhook_spec_id;

	DROP TABLE IF EXISTS webhook_specs;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0027_add_webhook_spec_tables",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up27).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down27).Error
		},
	})
}
