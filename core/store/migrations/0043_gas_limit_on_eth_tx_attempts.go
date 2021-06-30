package migrations

import (
	"gorm.io/gorm"
)

const up43 = `
ALTER TABLE eth_tx_attempts ADD COLUMN chain_specific_gas_limit bigint;
UPDATE eth_tx_attempts
SET chain_specific_gas_limit = eth_txes.gas_limit
FROM eth_txes
WHERE eth_txes.id = eth_tx_attempts.eth_tx_id;
ALTER TABLE eth_tx_attempts ALTER COLUMN chain_specific_gas_limit SET NOT NULL;
`
const down43 = `
ALTER TABLE eth_tx_attempts DROP COLUMN chain_specific_gas_limit;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0043_gas_limit_on_eth_tx_attempts.go",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up43).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down43).Error
		},
	})
}
