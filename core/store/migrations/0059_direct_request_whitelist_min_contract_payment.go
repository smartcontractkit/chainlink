package migrations

import (
	"gorm.io/gorm"
)

const up59 = `
ALTER TABLE direct_request_specs ADD COLUMN min_contract_payment numeric(78,0); 
`

const down59 = `
ALTER TABLE direct_request_specs DROP COLUMN min_contract_payment; 
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0059_direct_request_whitelist_min_contract_payment",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up59).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down59).Error
		},
	})
}
