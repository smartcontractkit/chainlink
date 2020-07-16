package migration1586949323

import (
	"github.com/jinzhu/gorm"
)

// Migrate makes txes.hash a unique index
// Due to faulty logic, there are some duplicates out in the wild (even though this should never happen).
// In the case of duplicate hashes, we keep the earliest row for each hash (by ID) and delete all other duplicate rows.
// Foreign key cascades will also cause all related tx_attempts to be deleted as well.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	DROP INDEX IF EXISTS idx_txs_hash;

	WITH duplicate_txes_hash AS (
		SELECT hash FROM txes GROUP BY hash HAVING count(hash) > 1
	), txes_to_delete AS (
		SELECT id
		FROM txes
		JOIN duplicate_txes_hash dups ON dups.hash = txes.hash
		WHERE id NOT IN (
			SELECT DISTINCT ON (txes.hash) id
			FROM txes
			JOIN duplicate_txes_hash dups ON dups.hash = txes.hash
			ORDER BY txes.hash, id ASC
		)
	)
	DELETE FROM txes WHERE id IN (SELECT id FROM txes_to_delete);

	CREATE UNIQUE INDEX idx_txes_hash ON txes (hash);
	`).Error
}
