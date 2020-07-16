package models

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Key holds the private key metadata for a given address that is used to unlock
// said key when given a password.
//
// By default, a key is assumed to represent an ethereum account.
type Key struct {
	Address   EIP55Address `gorm:"primary_key;type:varchar(64)"`
	JSON      JSON         `gorm:"type:text"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
}

type EncryptedSecretVRFKey = vrfkey.EncryptedSecretKey
type PublicKey = vrfkey.PublicKey

// NewKeyFromFile creates an instance in memory from a key file on disk.
func NewKeyFromFile(path string) (*Key, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	js := gjson.ParseBytes(dat)
	address, err := NewEIP55Address(common.HexToAddress(js.Get("address").String()).Hex())
	if err != nil {
		return nil, multierr.Append(errors.New("unable to create Key model"), err)
	}

	return &Key{Address: address, JSON: JSON{Result: js}}, nil
}

// WriteToDisk writes this key to disk at the passed path.
func (k *Key) WriteToDisk(path string) error {
	return utils.WriteFileWithPerms(path, []byte(k.JSON.String()), 0700)
}
