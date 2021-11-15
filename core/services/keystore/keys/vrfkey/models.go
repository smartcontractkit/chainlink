package vrfkey

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// EncryptedVRFKey contains encrypted private key to be serialized to DB
//
// We could re-use geth's key handling, here, but this makes it much harder to
// misuse VRF proving keys as ethereum keys or vice versa.
type EncryptedVRFKey struct {
	PublicKey secp256k1.PublicKey
	VRFKey    gethKeyStruct `json:"vrf_key"`
	CreatedAt time.Time     `json:"-"`
	UpdatedAt time.Time     `json:"-"`
	DeletedAt *time.Time    `json:"-"`
}

// JSON returns the JSON representation of e, or errors
func (e *EncryptedVRFKey) JSON() ([]byte, error) {
	keyJSON, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal encrypted key to JSON")
	}
	return keyJSON, nil
}

// WriteToDisk writes the JSON representation of e to given file path, and
// ensures the file has appropriate access permissions
func (e *EncryptedVRFKey) WriteToDisk(path string) error {
	keyJSON, err := e.JSON()
	if err != nil {
		return errors.Wrapf(err, "while marshaling key to save to %s", path)
	}
	userReadWriteOtherNoAccess := os.FileMode(0600)
	return utils.WriteFileWithMaxPerms(path, keyJSON, userReadWriteOtherNoAccess)
}

// Copied from go-ethereum/accounts/keystore/key.go's encryptedKeyJSONV3
type gethKeyStruct struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Version int                 `json:"version"`
}

func (k gethKeyStruct) Value() (driver.Value, error) {
	return json.Marshal(&k)
}

func (k *gethKeyStruct) Scan(value interface{}) error {
	var toUnmarshal []byte
	switch s := value.(type) {
	case []byte:
		toUnmarshal = s
	case string:
		toUnmarshal = []byte(s)
	default:
		return errors.Wrap(
			fmt.Errorf("unable to convert %+v of type %T to gethKeyStruct",
				value, value), "scan failure")
	}
	return json.Unmarshal(toUnmarshal, k)
}
