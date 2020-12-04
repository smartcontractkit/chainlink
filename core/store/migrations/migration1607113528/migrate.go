package migration1607113528

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE offchainreporting_oracle_specs drop not null;
    `).Error
}
