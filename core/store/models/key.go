package models

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// Key holds the private key metadata for a given address that is used to unlock
// said key when given a password.
//
// By default, a key is assumed to represent an ethereum account.
type Key struct {
	ID        int32 `gorm:"primary_key"`
	Address   EIP55Address
	JSON      JSON
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt null.Time `json:"-"`
	// This is the nonce that should be used for the next transaction.
	// Conceptually equivalent to geth's `PendingNonceAt` but more reliable
	// because we have a better view of our own transactions
	NextNonce *int64
	// LastUsed is the time that the address was last assigned to a transaction
	LastUsed *time.Time
	// IsFunding marks the address as being used for rescuing the  node and the pending transactions
	// Only one key can be IsFunding=true at a time.
	IsFunding bool
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

	return Key{Address: address, JSON: JSON{Result: js}}, nil
}

// WriteToDisk writes this key to disk at the passed path.
func (k *Key) WriteToDisk(path string) error {
	return utils.WriteFileWithMaxPerms(path, []byte(k.JSON.String()), 0600)
}
