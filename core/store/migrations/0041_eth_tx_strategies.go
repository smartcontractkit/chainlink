package migrations

import (
	"gorm.io/gorm"
)

const up41 = `
ALTER TABLE eth_txes ADD COLUMN subject uuid;
CREATE INDEX idx_eth_txes_unstarted_subject_id ON eth_txes (subject, id) WHERE subject IS NOT NULL AND state = 'unstarted';
`
const down41 = `
ALTER TABLE eth_txes DROP COLUMN subject;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0041_eth_tx_strategies",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up41).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down41).Error
		},
	})
}
