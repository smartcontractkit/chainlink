package migration1551895034

import (
	"fmt"

	ormpkg "github.com/smartcontractkit/chainlink/store/orm"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1551895034"
}

// Migrate creates a new indexable_block_numbers table with
// 1. the correct primary key because sqlite does not allow you to modify
// the primary key after table creation.
// 2. number backed by int64 instead of string.
func (m Migration) Migrate(orm *ormpkg.ORM) error {
	// db specific bytes -> hexadecimal conversion operation
	conversion := "hex(hash)" // sqlite default
	if orm.IsPostgres() {
		conversion = "encode(hash::bytea, 'hex')"
	}

	tx := orm.DB.Begin()
	err := tx.Exec(fmt.Sprintf(`
		CREATE TABLE "indexable_block_numbers_refactored_1551895034" (
			"number" bigint NOT NULL,
			"hash" varchar,
			PRIMARY KEY (hash));
		INSERT INTO "indexable_block_numbers_refactored_1551895034"
			SELECT CAST(number as bigint) as number, LOWER(%s) as hash
			FROM "indexable_block_numbers";
		DROP TABLE "indexable_block_numbers";
		ALTER TABLE "indexable_block_numbers_refactored_1551895034" RENAME TO "indexable_block_numbers";
		CREATE INDEX idx_indexable_block_numbers_number ON "indexable_block_numbers"("number");
	`, conversion)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
