package migrations

import (
	"gorm.io/gorm"
)

const up39 = `
ALTER TABLE flux_monitor_specs DROP COLUMN precision;
`
const down39 = `
ALTER TABLE flux_monitor_specs ADD COLUMN precision integer;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0039_remove_fmv2_precision",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up39).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down39).Error
		},
	})
}
