package migrations

import (
	"gorm.io/gorm"
)

const up52 = `
ALTER TABLE jobs ALTER COLUMN pipeline_spec_id SET NOT NULL;
`

const down52 = `
ALTER TABLE jobs ALTER COLUMN pipeline_spec_id DEFAULT NULL;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0052_not_null_job_pipeline_spec_id",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up52).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down52).Error
		},
	})
}
