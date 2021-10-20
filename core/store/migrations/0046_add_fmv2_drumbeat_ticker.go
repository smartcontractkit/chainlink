package migrations

import (
	"gorm.io/gorm"
)

const up46 = `
    ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_enabled boolean NOT NULL DEFAULT false;
    ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_schedule text;
`

const down46 = `
    ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_enabled;
    ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_schedule;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0046_add_fmv2_drumbeat_ticker",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up46).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down46).Error
		},
	})
}
