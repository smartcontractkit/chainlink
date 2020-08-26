package migration1590226486

import (
	"github.com/jinzhu/gorm"
)

// Migrate ensures that heads are unique and adds parent hash for use in reorg detection
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE UNIQUE INDEX idx_heads_hash ON heads (hash);
	ALTER TABLE heads
		ADD COLUMN parent_hash bytea,
		ADD COLUMN created_at timestamptz,
		ADD COLUMN timestamp timestamptz;
	UPDATE heads SET
		parent_hash = E'\\x0000000000000000000000000000000000000000000000000000000000000000',
		created_at = '2019-01-01',
		timestamp = '2019-01-01';
	ALTER TABLE heads
		ALTER COLUMN parent_hash SET NOT NULL,
		ALTER COLUMN created_at SET NOT NULL,
		ALTER COLUMN timestamp SET NOT NULL,
		ADD CONSTRAINT chk_hash_size CHECK (
			octet_length(hash) = 32
		),
		ADD CONSTRAINT chk_parent_hash_size CHECK (
			octet_length(parent_hash) = 32
		);
	`).Error
}
