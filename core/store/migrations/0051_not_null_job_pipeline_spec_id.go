package migrations

import (
	"gorm.io/gorm"
)

const up51 = `
ALTER TABLE jobs ALTER COLUMN pipeline_spec_id SET NOT NULL;
`

const down51 = `
ALTER TABLE jobs ALTER COLUMN pipeline_spec_id DEFAULT NULL;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0051_not_null_job_pipeline_spec_id",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up51).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down51).Error
		},
	})
}
