package migration1560886530

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Migrate converts the heads table to use a surrogate ID and binary hash
func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(`
DROP INDEX IF EXISTS idx_heads_number;
DROP INDEX IF EXISTS idx_heads_hash;
ALTER TABLE heads RENAME TO heads_archive;`).Error; err != nil {
		return errors.Wrap(err, "failed to drop heads")
	}

	if err := tx.AutoMigrate(&Head{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Head")
	}

	err := tx.Exec(`
INSERT INTO heads ("hash", "number")
SELECT decode("hash"::text, 'hex'), "number"
FROM heads_archive;
DROP TABLE heads_archive;`).Error
	return errors.Wrap(err, "failed to migrate old Heads")
}

// Head represents a BlockNumber, BlockHash.
type Head struct {
	ID     uint64      `gorm:"primary_key;auto_increment"`
	Hash   common.Hash `gorm:"not null"`
	Number int64       `gorm:"index;not null"`
}
