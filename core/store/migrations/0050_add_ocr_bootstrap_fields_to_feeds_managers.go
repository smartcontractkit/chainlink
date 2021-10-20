package migrations

import (
	"gorm.io/gorm"
)

const up50 = `
ALTER TABLE feeds_managers
DROP COLUMN network,
ADD COLUMN ocr_bootstrap_peer_multiaddr VARCHAR,
ADD CONSTRAINT chk_ocr_bootstrap_peer_multiaddr CHECK ( NOT (
	is_ocr_bootstrap_peer AND
	(
		ocr_bootstrap_peer_multiaddr IS NULL OR
		ocr_bootstrap_peer_multiaddr = ''
	)
));
`

const down50 = `
ALTER TABLE feeds_managers
ADD COLUMN network VARCHAR (100),
DROP CONSTRAINT chk_ocr_bootstrap_peer_multiaddr,
DROP COLUMN ocr_bootstrap_peer_multiaddr;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0050_add_ocr_bootstrap_fields_to_feeds_managers",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up50).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down50).Error
		},
	})
}
