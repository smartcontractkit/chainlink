package migration1560791143

import (
	"github.com/smartcontractkit/chainlink/core/store/dbutil"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	if !dbutil.IsPostgres(tx) {
		return nil
	}

	if err := tx.Exec(`SELECT SETVAL('txes_id_seq1', (SELECT MAX(id) FROM txes));`).Error; err != nil {
		return errors.Wrap(err, "failed to update sequence on tx")
	}
	return nil
}
