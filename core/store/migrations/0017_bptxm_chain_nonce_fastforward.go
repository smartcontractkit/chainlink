package migrations

import (
	"gorm.io/gorm"
)

const (
	up17 = `
UPDATE keys SET next_nonce = 0 WHERE next_nonce IS NULL;
ALTER TABLE keys ALTER COLUMN next_nonce SET NOT NULL, ALTER COLUMN next_nonce SET DEFAULT 0;
`
	down17 = `
ALTER TABLE keys ALTER COLUMN next_nonce SET DEFAULT NULL;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0017_bptxm_chain_nonce_fastforward",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up17).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down17).Error
		},
	})
}
