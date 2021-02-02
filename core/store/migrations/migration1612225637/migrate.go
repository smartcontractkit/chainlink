package migration1612225637

import "github.com/jinzhu/gorm"

// Migrate fixes the cases where eth_tx_attempts might be erroneously left forever in in_progress state
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	UPDATE eth_tx_attempts SET state = 'broadcast', broadcast_before_block_num = eth_receipts.block_number
	FROM eth_receipts
	WHERE eth_tx_attempts.state = 'in_progress' AND eth_tx_attempts.hash = eth_receipts.tx_hash
	`).Error
}
