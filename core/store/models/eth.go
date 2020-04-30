package models

import (
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	null "gopkg.in/guregu/null.v3"
)

// Tx contains fields necessary for an Ethereum transaction with
// an additional field for the TxAttempt.
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
	Value    *utils.Big     `gorm:"not null"`
	GasLimit uint64         `gorm:"not null"`

	// TxAttempt fields manually included; can't embed another primary_key
	Hash        common.Hash `gorm:"not null"`
	GasPrice    *utils.Big  `gorm:"not null"`
	Confirmed   bool        `gorm:"not null"`
	SentAt      uint64      `gorm:"not null"`
	SignedRawTx []byte      `gorm:"not null"`
	CreatedAt   time.Time   `json:"-"`
	UpdatedAt   time.Time   `json:"-"`
}

// String implements Stringer for Tx
func (tx *Tx) String() string {
	return fmt.Sprintf("Tx(ID: %d, From: %s, To: %s, Hash: %s, SentAt: %d)",
		tx.ID,
		tx.From.String(),
		tx.To.String(),
		tx.Hash.String(),
		tx.SentAt)
}

// EthTx creates a new Ethereum transaction with a given gasPrice in wei
// that is ready to be signed.
func (tx Tx) EthTx(gasPriceWei *big.Int) *types.Transaction {
	return types.NewTransaction(
		tx.Nonce,
		tx.To,
		tx.Value.ToInt(),
		tx.GasLimit,
		gasPriceWei,
		tx.Data,
	)
}

// TxAttempt is used for keeping track of transactions that
// have been written to the Ethereum blockchain. This makes
// it so that if the network is busy, a transaction can be
// resubmitted with a higher GasPrice.
type TxAttempt struct {
	ID uint64 `gorm:"primary_key;auto_increment"`

	TxID uint64 `gorm:"index;type:bigint REFERENCES txes(id) ON DELETE CASCADE"`
	Tx   *Tx    `json:"-" gorm:"PRELOAD:false;foreignkey:TxID"`

	CreatedAt time.Time `gorm:"index;not null"`

	Hash        common.Hash `gorm:"index;not null"`
	GasPrice    *utils.Big  `gorm:"type:varchar(78);not null"`
	Confirmed   bool        `gorm:"not null"`
	SentAt      uint64      `gorm:"not null"`
	SignedRawTx []byte      `gorm:"not null"`
	UpdatedAt   time.Time   `json:"-"`
}

// String implements Stringer for TxAttempt
func (txa *TxAttempt) String() string {
	return fmt.Sprintf("TxAttempt{ID: %d, TxID: %d, Hash: %s, SentAt: %d, Confirmed: %t}",
		txa.ID,
		txa.TxID,
		txa.Hash.String(),
		txa.SentAt,
		txa.Confirmed)
}

// GetID returns the ID of this structure for jsonapi serialization.
func (txa TxAttempt) GetID() string {
	return txa.Hash.Hex()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (txa TxAttempt) GetName() string {
	return "txattempts"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (txa *TxAttempt) SetID(value string) error {
	txa.Hash = common.HexToHash(value)
	return nil
}

func HighestPricedTxAttemptPerTx(items []TxAttempt) []TxAttempt {
	highestPricedSet := map[uint64]TxAttempt{}
	for _, item := range items {
		if currentHighest, ok := highestPricedSet[item.TxID]; ok {
			if currentHighest.GasPrice.ToInt().Cmp(item.GasPrice.ToInt()) == -1 {
				highestPricedSet[item.TxID] = item
			}
		} else {
			highestPricedSet[item.TxID] = item
		}
	}
	highestPriced := make([]TxAttempt, len(highestPricedSet))
	i := 0
	for _, attempt := range highestPricedSet {
		highestPriced[i] = attempt
		i++
	}
	return highestPriced
}

// Head represents a BlockNumber, BlockHash.
type Head struct {
	ID     uint64      `gorm:"primary_key;auto_increment"`
	Hash   common.Hash `gorm:"not null"`
	Number int64       `gorm:"index;not null"`
}

// AfterCreate is a gorm hook that trims heads after its creation
func (h Head) AfterCreate(scope *gorm.Scope) (err error) {
	scope.DB().Exec(`
	DELETE FROM heads
	WHERE id <= (
	  SELECT id
	  FROM (
		SELECT id
		FROM heads
		ORDER BY id DESC
		LIMIT 1 OFFSET 100
	  ) foo
	)`)
	if err != nil {
		return err
	}
	return nil
}

// NewHead returns a Head instance with a BlockNumber and BlockHash.
func NewHead(bigint *big.Int, hash common.Hash) *Head {
	if bigint == nil {
		return nil
	}

	return &Head{
		Number: bigint.Int64(),
		Hash:   hash,
	}
}

// String returns a string representation of this number.
func (l *Head) String() string {
	return l.ToInt().String()
}

// ToInt return the height as a *big.Int. Also handles nil by returning nil.
func (l *Head) ToInt() *big.Int {
	if l == nil {
		return nil
	}
	return big.NewInt(l.Number)
}

// GreaterThan compares BlockNumbers and returns true if the receiver BlockNumber is greater than
// the supplied BlockNumber
func (l *Head) GreaterThan(r *Head) bool {
	if l == nil {
		return false
	}
	if l != nil && r == nil {
		return true
	}
	return l.Number > r.Number
}

// NextInt returns the next BlockNumber as big.int, or nil if nil to represent latest.
func (l *Head) NextInt() *big.Int {
	if l == nil {
		return nil
	}
	return new(big.Int).Add(l.ToInt(), big.NewInt(1))
}
