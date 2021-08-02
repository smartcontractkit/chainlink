package migrations

import (
	"gorm.io/gorm"
)

const up53 = `
    ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_random_delay bigint;
`

const down53 = `
    ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_random_delay;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0053_add_fmv2_drumbeat_random_delay",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up53).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down53).Error
		},
	})
}
