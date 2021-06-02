package ethkey

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

// Key holds the private key metadata for a given address that is used to unlock
// said key when given a password.
//
// By default, a key is assumed to represent an ethereum account.
type Key struct {
	ID        int32 `gorm:"primary_key"`
	Address   EIP55Address
	JSON      postgres.Jsonb `json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	// This is the nonce that should be used for the next transaction.
	// Conceptually equivalent to geth's `PendingNonceAt` but more reliable
	// because we have a better view of our own transactions
	// NOTE: Be cautious about using this field, it is provided for convenience
	// only, can go out of date, and should not be relied upon. The source of
	// truth is always the database row for the key.
	NextNonce int64 `json:"-"`
	// IsFunding marks the address as being used for rescuing the  node and the pending transactions
	// Only one key can be IsFunding=true at a time.
	IsFunding bool
}

// Type returns type of key
func (k Key) Type() string {
	if k.IsFunding {
		return "funding"
	}
	return "sending"
}

// NewKeyFromFile creates an instance in memory from a key file on disk.
func NewKeyFromFile(path string) (Key, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return Key{}, err
	}

	js := gjson.ParseBytes(dat)
	address, err := NewEIP55Address(common.HexToAddress(js.Get("address").String()).Hex())
	if err != nil {
		return Key{}, multierr.Append(errors.New("unable to create Key model"), err)
	}

	return Key{Address: address, JSON: postgres.Jsonb{RawMessage: dat}}, nil
}
