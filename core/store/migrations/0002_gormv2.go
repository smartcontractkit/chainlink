package migrations

import (
	"gorm.io/gorm"
)

const up2 = `
UPDATE offchainreporting_oracle_specs SET contract_config_confirmations = 0 where contract_config_confirmations is NULL;
ALTER TABLE offchainreporting_oracle_specs
	ALTER COLUMN contract_config_confirmations SET NOT NULL;
ALTER TABLE external_initiators ADD CONSTRAINT "access_key_unique" UNIQUE ("access_key");
`

const down2 = `
ALTER TABLE offchainreporting_oracle_specs 
	ALTER COLUMN contract_config_confirmations DROP NOT NULL;
ALTER TABLE external_initiators DROP CONSTRAINT "access_key_unique";
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0002_gormv2",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up2).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down2).Error
		},
	})
}
