package migrations

import (
	"gorm.io/gorm"
)

const up48 = `
ALTER TABLE feeds_managers
ADD COLUMN is_ocr_bootstrap_peer boolean NOT NULL DEFAULT false;
`

const down48 = `
ALTER TABLE feeds_managers
DROP COLUMN is_ocr_bootstrap_peer;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0048_add_ocr_bootstrap_node_to_feeds_manager",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up48).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down48).Error
		},
	})
}
