package migration1574659987

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds VRF proving-key table
func Migrate(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE encrypted_secret_keys (
			public_key character varying(68) PRIMARY KEY,
			vrf_key text NOT NULL,
			created_at timestamp with time zone NOT NULL,
			updated_at timestamp with time zone NOT NULL
		);
	`).Error
}
