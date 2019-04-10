package migration1551816486

import (
	"github.com/jinzhu/gorm"
)

type Migration struct{}

// Migrate creates a new bridge_types table with the correct primary key
// because sqlite does not allow you to modify the primary key
// after table creation.
func (m Migration) Migrate(tx *gorm.DB) error {
	tx = tx.Begin()
	err := tx.Exec(`
		CREATE TABLE "bridge_types_with_pk" ("name" varchar(255),"url" varchar(255),"confirmations" bigint,"incoming_token" varchar(255),"outgoing_token" varchar(255),"minimum_contract_payment" varchar(255) , PRIMARY KEY ("name"));
		INSERT INTO "bridge_types_with_pk" SELECT * FROM "bridge_types";
		DROP TABLE "bridge_types";
		ALTER TABLE "bridge_types_with_pk" RENAME TO "bridge_types";
	`).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

type BridgeType struct {
	Name                   string
	URL                    string
	Confirmations          uint64
	IncomingToken          string
	OutgoingToken          string
	MinimumContractPayment string
}
