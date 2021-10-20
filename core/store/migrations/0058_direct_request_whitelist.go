package migrations

import (
	"gorm.io/gorm"
)

const up58 = `
ALTER TABLE direct_request_specs ADD COLUMN requesters TEXT; 
`

const down58 = `
ALTER TABLE direct_request_specs DROP COLUMN requesters; 
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0058_direct_request_whitelist",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up58).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down58).Error
		},
	})
}
