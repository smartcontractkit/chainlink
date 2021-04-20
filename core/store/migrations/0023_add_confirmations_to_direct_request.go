package migrations

import (
	"gorm.io/gorm"
)

const (
	up23 = `
ALTER TABLE direct_request_specs ADD COLUMN num_confirmations bigint DEFAULT NULL;
`
	down23 = `
ALTER TABLE direct_request_specs DROP COLUMN num_confirmations;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0023_add_confirmations_to_direct_request",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up23).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down23).Error
		},
	})
}
