package migrations

import (
	"gorm.io/gorm"
)

const up19 = `
		ALTER TABLE upkeep_registrations ADD COLUMN last_run_block_height BIGINT NOT NULL DEFAULT 0;
	`

const down19 = `
		ALTER TABLE upkeep_registrations DROP COLUMN last_run_block_height;
	`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0019_last_run_height_column_to_keeper_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up19).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down19).Error
		},
	})
}
