package migrations

import (
	"gorm.io/gorm"
)

const up27 = `
ALTER TABLE offchainreporting_latest_round_requested
DROP CONSTRAINT offchainreporting_latest_roun_offchainreporting_oracle_spe_fkey,
ADD CONSTRAINT offchainreporting_latest_roun_offchainreporting_oracle_spe_fkey
	FOREIGN KEY (offchainreporting_oracle_spec_id)
	REFERENCES offchainreporting_oracle_specs (id)
	ON DELETE CASCADE
	DEFERRABLE INITIALLY IMMEDIATE
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0027_cascade_ocr_latest_round_request",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up27).Error
		},
		Rollback: func(db *gorm.DB) error {
			return nil
		},
	})
}
