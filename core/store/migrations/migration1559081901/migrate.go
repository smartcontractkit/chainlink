package migration1559081901

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(
		`DROP INDEX idx_txes_from;
		 DROP INDEX idx_txes_nonce;
		 DROP INDEX idx_tx_attempts_tx_id;
		 DROP INDEX idx_tx_attempts_created_at;
		 ALTER TABLE txes RENAME TO txes_archive;
		 ALTER TABLE tx_attempts RENAME TO tx_attempts_archive;
	`,
	).Error; err != nil {
		return errors.Wrap(err, "failed to drop txes and txattempts")
	}
	if err := tx.AutoMigrate(&models.Tx{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Tx")
	}
	if err := tx.AutoMigrate(&models.TxAttempt{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TxAttempt")
	}
	if err := tx.Exec(
		`INSERT INTO txes (
			"id", "from", "to", "data", "nonce", "value", "gas_limit", "hash", "gas_price", "confirmed", "sent_at", "signed_raw_tx"
		 )
		 SELECT
			"id", "from", "to", "data", "nonce", "value", "gas_limit", "hash", "gas_price", "confirmed", "sent_at", "hex"
		 FROM txes_archive;
		 INSERT INTO tx_attempts (
			"hash", "tx_id", "gas_price", "confirmed", "sent_at", "created_at", "signed_raw_tx"
		 )
		 SELECT
			"hash", "tx_id", "gas_price", "confirmed", "sent_at", "created_at", "hex"
		 FROM tx_attempts_archive;
		 DROP TABLE txes_archive;
		 DROP TABLE tx_attempts_archive;
		 `).Error; err != nil {
		return errors.Wrap(err, "failed to migrate old Txes, TxAttempts")
	}
	return nil
}
