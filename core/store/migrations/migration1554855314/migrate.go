package migration1554855314

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1551816486"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Migrate adds the sync_events table
func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(`
CREATE TABLE "bridge_types_with_incoming_token_hash" (
	"name" varchar(255),
	"url" varchar(255),
	"confirmations" bigint,
	"incoming_token_hash" VARCHAR(64),
	"salt" VARCHAR(32),
	"outgoing_token" varchar(255),
	"minimum_contract_payment" varchar(255),
	PRIMARY KEY ("name")
);`).Error; err != nil {
		return errors.Wrap(err, "error creating temporary bridge_types_with_incoming_token_hash table")
	}

	// CAUTION: This could be quite slow with a big enough database
	if err := orm.Batch(1000, func(offset, limit uint) (uint, error) {
		var bts []migration1551816486.BridgeType
		err := tx.Limit(limit).Offset(offset).Find(&bts).Error
		if err != nil {
			return 0, errors.Wrap(err, "error loading old bridge types")
		}

		for _, bt := range bts {
			newBT := BridgeType{
				Name:                   bt.Name,
				URL:                    bt.URL,
				Confirmations:          bt.Confirmations,
				OutgoingToken:          bt.OutgoingToken,
				MinimumContractPayment: bt.MinimumContractPayment,
				Salt:                   utils.NewBytes32ID(),
			}

			hash, err := utils.Sha256(bt.IncomingToken)
			if err != nil {
				return 0, errors.Wrap(err, "error generating new bridge type token hash")
			}
			newBT.IncomingTokenHash = hash

			if err := tx.Table("bridge_types_with_incoming_token_hash").Create(&newBT).Error; err != nil {
				return 0, errors.Wrap(err, "error saving new bridge type")
			}
		}

		return uint(len(bts)), nil
	}); err != nil {
		return errors.Wrap(err, "error migrating hash for bridge_types")
	}

	if err := tx.Exec(`
DROP TABLE "bridge_types";
ALTER TABLE "bridge_types_with_incoming_token_hash" RENAME TO "bridge_types";
	`).Error; err != nil {
		return errors.Wrap(err, "error renaming temporary bridge_types_with_incoming_token_hash table")
	}

	return nil
}

// BridgeType is for migrating between bridge type tables
type BridgeType struct {
	Name                   string `gorm:"primary_key"`
	URL                    string
	Confirmations          uint64
	IncomingTokenHash      string
	Salt                   string
	OutgoingToken          string
	MinimumContractPayment string `gorm:"type:varchar(255)"`
}
