package migrations

import (
	"gorm.io/gorm"
)

// Note this automatically updates the FK in jobs
const up54 = `
	UPDATE jobs SET type='ethlog' WHERE type='directrequest';
	ALTER TABLE jobs RENAME COLUMN direct_request_spec_id TO eth_log_spec_id;
    ALTER TABLE direct_request_specs RENAME TO eth_log_specs;
`

const down54 = `
	UPDATE jobs SET type='directrequest' WHERE type='ethlog';
	ALTER TABLE jobs RENAME COLUMN  eth_log_spec_id TO direct_request_spec_id;
    ALTER TABLE eth_log_specs RENAME TO direct_request_specs;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0054_direct_request_rename_eth_log",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up54).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down54).Error
		},
	})
}
