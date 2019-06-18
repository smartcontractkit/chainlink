package migration1560886530

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// Migrate converts the heads table to use a surrogate ID and binary hash
func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(`
DROP INDEX IF EXISTS idx_heads_number;
DROP INDEX IF EXISTS idx_heads_hash;
ALTER TABLE heads RENAME TO heads_archive;`).Error; err != nil {
		return errors.Wrap(err, "failed to drop heads")
	}

	if err := tx.AutoMigrate(&models.Head{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Head")
	}

	var err error
	if dbutil.IsPostgres(tx) {
		err = tx.Exec(`
INSERT INTO heads ("hash", "number")
SELECT decode(convert_from("hash", 'utf-8'), 'hex'), "number"
FROM heads_archive;
DROP TABLE heads_archive;`).Error
	} else {
		// SQLite doesn't support decoding at the SQL level
		err = orm.Batch(1000, func(offset, limit uint) (uint, error) {
			var heads []migration0.Head
			err := tx.
				Table("heads_archive").
				Limit(limit).
				Offset(offset).
				Order("number").
				Find(&heads).Error
			if err != nil {
				return 0, err
			}

			for _, head := range heads {
				migratedHead := models.Head{
					Hash:   common.HexToHash(head.HashRaw),
					Number: head.Number,
				}
				err = tx.Create(&migratedHead).Error
				if err != nil {
					return 0, err
				}
			}

			return uint(len(heads)), err
		})
	}
	return errors.Wrap(err, "failed to migrate old Heads")
}
