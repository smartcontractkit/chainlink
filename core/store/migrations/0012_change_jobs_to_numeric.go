package migrations

import (
	"gorm.io/gorm"
)

const (
	up12 = `
ALTER TABLE job_specs ALTER COLUMN min_payment TYPE numeric(78, 0) USING min_payment::numeric;
ALTER TABLE flux_monitor_specs ALTER COLUMN min_payment TYPE numeric(78, 0) USING min_payment::numeric;
`
	down12 = `
ALTER TABLE job_specs ALTER COLUMN min_payment TYPE varchar(255) USING min_payment::varchar;
ALTER TABLE flux_monitor_specs ALTER COLUMN min_payment TYPE varchar(255) USING min_payment::varchar;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0012_change_jobs_to_numeric",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up12).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down12).Error
		},
	})
}
