package migrations

import "gorm.io/gorm"

const (
	up27 = `
        CREATE TABLE vrf_specs (
            id BIGSERIAL PRIMARY KEY,
            public_key text NOT NULL,
			coordinator_address bytea NOT NULL,
			confirmations bigint NOT NULL,
            created_at timestamp with time zone NOT NULL,
            updated_at timestamp with time zone NOT NULL
        );
        ALTER TABLE jobs ADD COLUMN vrf_spec_id INT REFERENCES vrf_specs(id),
        DROP CONSTRAINT chk_only_one_spec,
        ADD CONSTRAINT chk_only_one_spec CHECK (
            num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_spec_id, vrf_spec_id) = 1
        );
    `
	down27 = `
        ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
        ADD CONSTRAINT chk_only_one_spec CHECK (
            num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_spec_id) = 1
        );

        ALTER TABLE jobs DROP COLUMN vrf_spec_id;
        DROP TABLE IF EXISTS vrf_specs;
    `
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0027_vrf_v2",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up27).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down27).Error
		},
	})
}
