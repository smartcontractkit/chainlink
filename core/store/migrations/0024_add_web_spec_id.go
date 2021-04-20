package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up24 = `
		CREATE TABLE web_specs (
			id SERIAL PRIMARY KEY,
			created_at timestamp with time zone NOT NULL,
			updated_at timestamp with time zone NOT NULL
		);

		ALTER TABLE jobs ADD COLUMN web_spec_id INT REFERENCES web_specs(id),
		DROP CONSTRAINT chk_only_one_spec,
		ADD CONSTRAINT chk_only_one_spec CHECK (
			num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, web_spec_id) = 1
		);
	`

const down24 = `	
		ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
		ADD CONSTRAINT chk_only_one_spec CHECK (
			num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id) = 1
		);
	
		ALTER TABLE jobs DROP COLUMN web_spec_id;

		DROP TABLE IF EXISTS web_specs;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0024_add_web_spec_tables",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up24).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down24).Error
		},
	})
}
