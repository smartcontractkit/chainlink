package migration1596021087

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the encrypted_p2p_keys table and renames encryped_secret_keys to the more specific encrypted_vrf_keys
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE encrypted_secret_keys RENAME TO encrypted_vrf_keys;
		CREATE TABLE encrypted_p2p_keys (
			id SERIAL PRIMARY KEY,
			peer_id text NOT NULL,
			pub_key bytea NOT NULL,
			encrypted_priv_key jsonb NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL
		);

		CREATE UNIQUE INDEX idx_unique_pub_keys ON encrypted_p2p_keys (pub_key);
		CREATE UNIQUE INDEX idx_unique_peer_ids ON encrypted_p2p_keys (peer_id);
		ALTER TABLE encrypted_p2p_keys ADD CONSTRAINT chk_pub_key_length CHECK (
			octet_length(pub_key) = 32
		);
	`).Error
}
