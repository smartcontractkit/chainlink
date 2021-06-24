package migrations

import (
	"gorm.io/gorm"
)

const up40 = `
ALTER TABLE heads ADD COLUMN l1_block_number bigint;
`
const down40 = `
ALTER TABLE heads DROP COLUMN l1_block_number;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0040_heads_l1_block_number",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up40).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down40).Error
		},
	})
}
