package migrations

import (
	"gorm.io/gorm"
)

const up53 = `
    ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_random_delay bigint NOT NULL DEFAULT 0;

		UPDATE flux_monitor_specs SET drumbeat_schedule = '' where drumbeat_schedule IS NULL;
		ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule SET DEFAULT '';
		ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule SET NOT NULL;
`

const down53 = `
    ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule SET NULL;
    ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule DROP DEFAULT;
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
