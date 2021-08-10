package migrations

import (
	"gorm.io/gorm"
)

const up55 = `
ALTER TABLE job_proposals
ADD COLUMN multiaddrs TEXT[] DEFAULT NULL;
`

const down55 = `
ALTER TABLE job_proposals
DROP COLUMN multiaddrs;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0055_add_multiaddrs_to_job_proposals",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up55).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down55).Error
		},
	})
}
