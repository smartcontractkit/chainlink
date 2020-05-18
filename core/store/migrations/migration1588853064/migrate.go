package migration1588853064

import (
	"github.com/jinzhu/gorm"
)

// Migrate makes the nonce column on txes unique per account
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		-- where one of the txes is confirmed, keep the latest confirmed one
		-- if none are confirmed, keep the latest one
		WITH duplicate_txes_by_nonce AS (
			SELECT nonce, "from"
			FROM txes
			GROUP BY nonce, "from"
			HAVING count(id) > 1
		),
		txes_to_keep AS (
			SELECT DISTINCT ON (txes.nonce) txes.nonce, id, confirmed, txes."from"
			FROM txes
			JOIN duplicate_txes_by_nonce
			ON duplicate_txes_by_nonce.nonce = txes.nonce
			AND duplicate_txes_by_nonce.from = txes.from
			ORDER BY nonce, "from", confirmed desc, id desc
		),
		txes_to_delete AS (
			SELECT id
			FROM txes
			JOIN duplicate_txes_by_nonce
			ON duplicate_txes_by_nonce.nonce = txes.nonce
			AND duplicate_txes_by_nonce.from = txes.from
			WHERE id NOT IN (SELECT id FROM txes_to_keep)
		)
		DELETE FROM txes WHERE id IN (SELECT id FROM txes_to_delete);

		DROP INDEX idx_txes_nonce;
	  	CREATE UNIQUE INDEX idx_txes_unique_nonces_per_account ON txes(nonce, "from");
	`).Error
}
