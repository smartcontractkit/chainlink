package migrations

import (
	"gorm.io/gorm"
)

func init() {
	// Note this is the 1612225637 v1 migration which never got run
	// because the v2 migrations started after 1611847145.
	Migrations = append(Migrations, &Migration{
		ID: "0004_cleanup_tx_state",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(`
				UPDATE eth_tx_attempts SET state = 'broadcast', broadcast_before_block_num = eth_receipts.block_number
				FROM eth_receipts
				WHERE eth_tx_attempts.state = 'in_progress' AND eth_tx_attempts.hash = eth_receipts.tx_hash
			`).Error
		},
		Rollback: func(db *gorm.DB) error {
			return nil
		},
	})
}
