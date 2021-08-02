package migrations

import (
	"gorm.io/gorm"
)

const up47 = `
    ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_jitter_within_period bigint NOT NULL DEFAULT 0;
`

const down47 = `
    ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_jitter_within_period;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0047_add_fmv2_drumbeat_jitter",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up47).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down47).Error
		},
	})
}
