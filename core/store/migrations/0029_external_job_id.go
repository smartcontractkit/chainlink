package migrations

import "gorm.io/gorm"

const (
	up29 = `
               ALTER TABLE direct_request_specs DROP COLUMN on_chain_job_spec_id;
               ALTER TABLE jobs ADD COLUMN external_job_id uuid NOT NULL DEFAULT uuid_generate_v4();
               ADD CONSTRAINT external_job_id_uniq UNIQUE(external_job_id);
       `
	down29 = `
               ALTER TABLE direct_request_specs ADD COLUMN on_chain_job_spec_id bytea;
               ALTER TABLE jobs DROP CONSTRAINT external_job_id_uniq;
       `
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0029_external_job_id",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up29).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down29).Error
		},
	})
}
