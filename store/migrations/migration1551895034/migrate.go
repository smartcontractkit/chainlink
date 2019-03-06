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
	if !orm.DB.HasTable("indexable_block_numbers") {
		return nil
	}

	// db specific bytes -> hexadecimal conversion operation
	conversion := "hex(hash)" // sqlite default
	if orm.IsPostgres() {
		conversion = "encode(hash::bytea, 'hex')"
	}

	tx := orm.DB.Begin()
	err := tx.Exec(fmt.Sprintf(`
		CREATE TABLE "heads" (
			"number" bigint NOT NULL,
			"hash" varchar,
			PRIMARY KEY (hash));
		INSERT INTO "heads"
			SELECT CAST(number as bigint) as number, LOWER(%s) as hash
			FROM "indexable_block_numbers";
		DROP TABLE "indexable_block_numbers";
		CREATE INDEX idx_heads_number ON "heads"("number");
	`, conversion)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
