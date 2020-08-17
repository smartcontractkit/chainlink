package migration1559081901

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	null "gopkg.in/guregu/null.v3"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(
		`DROP INDEX idx_txes_from;
		 DROP INDEX idx_txes_nonce;
		 DROP INDEX idx_tx_attempts_tx_id;
		 DROP INDEX idx_tx_attempts_created_at;
		 ALTER TABLE txes RENAME TO txes_archive;
		 ALTER TABLE tx_attempts RENAME TO tx_attempts_archive;
	`,
	).Error; err != nil {
		return errors.Wrap(err, "failed to drop txes and txattempts")
	}
	if err := tx.AutoMigrate(&Tx{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Tx")
	}
	if err := tx.AutoMigrate(&TxAttempt{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TxAttempt")
	}
	if err := tx.Exec(
		`INSERT INTO txes (
			"id", "from", "to", "data", "nonce", "value", "gas_limit", "hash", "gas_price", "confirmed", "sent_at", "signed_raw_tx"
		 )
		 SELECT
			"id", "from", "to", "data", "nonce", "value", "gas_limit", "hash", "gas_price", "confirmed", "sent_at", "hex"
		 FROM txes_archive;
		 INSERT INTO tx_attempts (
			"hash", "tx_id", "gas_price", "confirmed", "sent_at", "created_at", "signed_raw_tx"
		 )
		 SELECT
			"hash", "tx_id", "gas_price", "confirmed", "sent_at", "created_at", "hex"
		 FROM tx_attempts_archive;
		 DROP TABLE txes_archive;
		 DROP TABLE tx_attempts_archive;
		 `).Error; err != nil {
		return errors.Wrap(err, "failed to migrate old Txes, TxAttempts")
	}
	return nil
}

// Tx is a capture of the model representing Txes before migration1586369235
// Let's please not use gorm automigrate ever again
type Tx struct {
	ID uint64 `gorm:"primary_key;auto_increment"`

	// SurrogateID is used to look up a transaction using a secondary ID, used to
	// associate jobs with transactions so that we don't double spend in certain
	// failure scenarios
	SurrogateID null.String `gorm:"index;unique"`

	Attempts []*TxAttempt `json:"-"`

	From     common.Address `gorm:"index;not null"`
	To       common.Address `gorm:"not null"`
	Data     []byte         `gorm:"not null"`
	Nonce    uint64         `gorm:"index;not null"`
	Value    *utils.Big     `gorm:"type:varchar(78);not null"`
	GasLimit uint64         `gorm:"not null"`

	// TxAttempt fields manually included; can't embed another primary_key
	Hash        common.Hash `gorm:"not null"`
	GasPrice    *utils.Big  `gorm:"type:varchar(78);not null"`
	Confirmed   bool        `gorm:"not null"`
	SentAt      uint64      `gorm:"not null"`
	SignedRawTx string      `gorm:"type:text;not null"`
}

// TxAttempt is a capture of the model representing TxAttempts before migration1586369235
type TxAttempt struct {
	ID uint64 `gorm:"primary_key;auto_increment"`

	TxID uint64 `gorm:"index;type:bigint REFERENCES txes(id) ON DELETE CASCADE"`
	Tx   *Tx    `json:"-" gorm:"PRELOAD:false;foreignkey:TxID"`

	CreatedAt time.Time `gorm:"index;not null"`

	Hash        common.Hash `gorm:"index;not null"`
	GasPrice    *utils.Big  `gorm:"type:varchar(78);not null"`
	Confirmed   bool        `gorm:"not null"`
	SentAt      uint64      `gorm:"not null"`
	SignedRawTx string      `gorm:"type:text;not null"`
}
