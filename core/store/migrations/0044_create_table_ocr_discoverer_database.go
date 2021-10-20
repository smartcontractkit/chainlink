package migrations

import (
	"gorm.io/gorm"
)

const up44 = `
CREATE TABLE offchainreporting_discoverer_announcements (
	local_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id) DEFERRABLE INITIALLY IMMEDIATE,
	remote_peer_id text NOT NULL,
	ann bytea NOT NULL,
	created_at timestamptz not null,
	updated_at timestamptz not null,
	PRIMARY KEY(local_peer_id, remote_peer_id)
);
`
const down44 = `
DROP TABLE offchainreporting_discoverer_announcements;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0044_create_table_offchainreporting_discoverer_announcements",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up44).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down44).Error
		},
	})
}
