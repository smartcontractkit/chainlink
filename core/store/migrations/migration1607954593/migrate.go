package migration1607954593

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE p2p_peers ADD COLUMN peer_id text REFERENCES encrypted_p2p_keys (peer_id) DEFERRABLE INITIALLY IMMEDIATE;
		
		UPDATE p2p_peers SET peer_id = offchainreporting_oracle_specs.p2p_peer_id
		FROM offchainreporting_oracle_specs
		JOIN jobs ON jobs.offchainreporting_oracle_spec_id = offchainreporting_oracle_specs.id
		WHERE jobs.id = p2p_peers.job_id;

		ALTER TABLE p2p_peers ALTER COLUMN peer_id SET NOT NULL, DROP COLUMN job_id;

		CREATE INDEX p2p_peers_peer_id ON p2p_peers (peer_id);
    `).Error
}
