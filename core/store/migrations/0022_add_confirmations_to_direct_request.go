package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	up22 = `
ALTER TABLE direct_request_specs ADD COLUMN num_confirmations bigint DEFAULT 1 NOT NULL;
`
	down22 = `
ALTER TABLE direct_request_specs DROP COLUMN num_confirmations;
`
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0022_add_confirmations_to_direct_request",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up22).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down22).Error
		},
	})
}
