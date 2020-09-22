package migration1600504870

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds a new index necessary for EthConfirmer's FindEthTxsRequiringNewAttempt
// This index will be very tiny and fast since it only applies to unconfirmed eth_txes
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key ON eth_txes (nonce, from_address) WHERE state = 'unconfirmed';
	`).Error
}
