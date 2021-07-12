package migrations

import (
	"gorm.io/gorm"
)

const (
	up26 = `
ALTER TABLE eth_txes ADD COLUMN meta jsonb;
`
	down26 = `
ALTER TABLE eth_txes DROP COLUMN meta;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0026_eth_tx_meta",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up26).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down26).Error
		},
	})
}
