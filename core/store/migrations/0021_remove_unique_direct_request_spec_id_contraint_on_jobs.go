package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up21 = `
	DROP INDEX idx_jobs_unique_direct_request_spec_id;
`

const down21 = `
	CREATE UNIQUE INDEX idx_jobs_unique_direct_request_spec_id ON jobs USING btree (direct_request_spec_id);
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0021_remove_unique_direct_request_spec_id_contraint_on_jobs",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up21).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down21).Error
		},
	})
}
