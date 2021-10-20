package migrations

import (
	"gorm.io/gorm"
)

const up32 = `ALTER TABLE direct_request_specs RENAME COLUMN num_confirmations TO min_incoming_confirmations;`
const down32 = `ALTER TABLE direct_request_specs RENAME COLUMN min_incoming_confirmations TO num_confirmations;`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0032_rename_direct_request_specs_num_confirmations",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up32).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down32).Error
		},
	})
}
