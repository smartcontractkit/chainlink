package migrations

import (
	"gorm.io/gorm"
)

const up49 = `
ALTER TABLE job_proposals
RENAME COLUMN job_id TO external_job_id;

ALTER INDEX idx_job_proposals_job_id RENAME TO idx_job_proposals_external_job_id;
`

const down49 = `
ALTER TABLE job_proposals
RENAME COLUMN external_job_id TO job_id;

ALTER INDEX idx_job_proposals_external_job_id RENAME TO idx_job_proposals_job_id;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0049_rename_job_id_to_external_job_id_in_job_proposals",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up49).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down49).Error
		},
	})
}
