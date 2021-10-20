package migrations

import (
	"gorm.io/gorm"
)

const up47 = `
ALTER TABLE job_proposals
ADD COLUMN remote_uuid UUID NOT NULL;

CREATE UNIQUE INDEX idx_job_proposals_remote_uuid ON job_proposals(remote_uuid);
`

const down47 = `
ALTER TABLE job_proposals
DROP COLUMN remote_uuid;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0047_add_uuid_to_job_proposals",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up47).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down47).Error
		},
	})
}
