package migration1551895034

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/store/dbutil"
)

type Migration struct{}

// Migrate creates a new indexable_block_numbers table with
// 1. the correct primary key because sqlite does not allow you to modify
// the primary key after table creation.
// 2. number backed by int64 instead of string.
func (m Migration) Migrate(tx *gorm.DB) error {
	if !tx.HasTable("indexable_block_numbers") {
		return nil
	}

	// db specific bytes -> hexadecimal conversion operation
	conversion := "hex(hash)" // sqlite default
	if dbutil.IsPostgres(tx) {
		conversion = "encode(hash::bytea, 'hex')"
	}

	return tx.Exec(fmt.Sprintf(`
		BEGIN TRANSACTION;
		CREATE TABLE "heads" (
			"number" bigint NOT NULL,
			"hash" varchar,
			PRIMARY KEY (hash));
		INSERT INTO "heads"
			SELECT CAST(number as bigint) as number, LOWER(%s) as hash
			FROM "indexable_block_numbers";
		DROP TABLE "indexable_block_numbers";
		CREATE INDEX idx_heads_number ON "heads"("number");
		COMMIT TRANSACTION;
	`, conversion)).Error
}
