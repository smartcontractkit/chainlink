package migrationsv2

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up2 = `
UPDATE offchainreporting_oracle_specs SET contract_config_confirmations = 0 where contract_config_confirmations is NULL;
ALTER TABLE offchainreporting_oracle_specs
	ALTER COLUMN contract_config_confirmations SET NOT NULL;
`

const down2 = `
ALTER TABLE offchainreporting_oracle_specs 
	ALTER COLUMN contract_config_confirmations DROP NOT NULL;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0002_ocr_confs_not_null",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up2).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down2).Error
		},
	})
}
