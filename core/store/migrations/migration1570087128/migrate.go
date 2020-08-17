package migration1570087128

import (
	"github.com/smartcontractkit/chainlink/core/store/dbutil"

	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	if dbutil.IsPostgres(tx) {
		return tx.Exec(`
ALTER TABLE initiators ADD COLUMN "from_block" numeric(78, 0);
ALTER TABLE initiators ADD COLUMN "to_block" numeric(78, 0);
ALTER TABLE initiators ADD COLUMN "topics" text;
`).Error
	}

	return tx.Exec(`
ALTER TABLE initiators ADD COLUMN "from_block" varchar(255);
ALTER TABLE initiators ADD COLUMN "to_block" varchar(255);
ALTER TABLE initiators ADD COLUMN "topics" text;
    `).Error
}
