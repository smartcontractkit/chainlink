package migrations

import (
	"gorm.io/gorm"
)

const (
	up5 = `
DROP INDEX IF EXISTS idx_eth_tx_attempts_in_progress;
CREATE INDEX idx_eth_tx_attempts_unbroadcast ON eth_tx_attempts (state enum_ops) WHERE state != 'broadcast'::eth_tx_attempts_state;
DROP INDEX IF EXISTS idx_only_one_in_progress_attempt_per_eth_tx;
CREATE UNIQUE INDEX idx_only_one_unbroadcast_attempt_per_eth_tx ON eth_tx_attempts(eth_tx_id int8_ops) WHERE state != 'broadcast'::eth_tx_attempts_state;
DROP INDEX IF EXISTS idx_eth_txes_state;
CREATE INDEX idx_eth_txes_state_from_address ON eth_txes(state, from_address) WHERE state <> 'confirmed'::eth_txes_state;
`
	down5 = `
DROP INDEX IF EXISTS idx_eth_tx_attempts_unbroadcast;
CREATE INDEX idx_eth_tx_attempts_in_progress ON eth_tx_attempts(state enum_ops) WHERE state = 'in_progress'::eth_tx_attempts_state;
DROP INDEX IF EXISTS idx_only_one_unbroadcast_attempt_per_eth_tx;
CREATE UNIQUE INDEX idx_only_one_in_progress_attempt_per_eth_tx ON eth_tx_attempts(eth_tx_id int8_ops) WHERE state = 'in_progress'::eth_tx_attempts_state;
DROP INDEX IF EXISTS idx_eth_txes_state_from_address;
CREATE INDEX idx_eth_txes_state ON eth_txes(state enum_ops) WHERE state <> 'confirmed'::eth_txes_state;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0005_eth_tx_attempts_insufficient_eth_index",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up5).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down5).Error
		},
	})
}
