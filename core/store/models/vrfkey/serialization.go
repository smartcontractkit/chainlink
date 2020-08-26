package vrfkey

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// EncryptedVRFKey contains encrypted private key to be serialized to DB
//
// We could re-use geth's key handling, here, but this makes it much harder to
// misuse VRF proving keys as ethereum keys or vice versa.
type EncryptedVRFKey struct {
	PublicKey PublicKey     `gorm:"primary_key"`
	VRFKey    gethKeyStruct `json:"vrf_key"`
	CreatedAt time.Time     `json:"-"`
	UpdatedAt time.Time     `json:"-"`
}

// passwordPrefix is added to the beginning of the passwords for
// EncryptedVRFKey's, so that VRF keys can't casually be used as ethereum
// keys, and vice-versa. If you want to do that, DON'T.
var passwordPrefix = "don't mix VRF and Ethereum keys!"

func adulteratedPassword(auth string) string {
	return passwordPrefix + auth
}

type ScryptParams struct{ N, P int }

var defaultScryptParams = ScryptParams{
	N: keystore.StandardScryptN, P: keystore.StandardScryptP}

// FastScryptParams is for use in tests, where you don't want to wear out your
// CPU with expensive key derivations, do not use it in production, or your
// encrypted VRF keys will be easy to brute-force!
var FastScryptParams = ScryptParams{N: 2, P: 1}

// Encrypt returns the key encrypted with passphrase auth
func (k *PrivateKey) Encrypt(auth string, p ...ScryptParams,
) (*EncryptedVRFKey, error) {
	var scryptParams ScryptParams
	switch len(p) {
	case 0:
		scryptParams = defaultScryptParams
	case 1:
		scryptParams = p[0]
	default:
		return nil, fmt.Errorf("can take at most one set of ScryptParams")
	}
	keyJSON, err := keystore.EncryptKey(k.gethKey(), adulteratedPassword(auth),
		scryptParams.N, scryptParams.P)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt vrf key")
	}
	rv := EncryptedVRFKey{}
	if e := json.Unmarshal(keyJSON, &rv.VRFKey); e != nil {
		return nil, errors.Wrapf(e, "geth returned unexpected key material")
	}
	rv.PublicKey = k.PublicKey
	roundTripKey, err := rv.Decrypt(auth)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt just-encrypted key!")
	}
	if !roundTripKey.k.Equal(k.k) || roundTripKey.PublicKey != k.PublicKey {
		panic(fmt.Errorf("roundtrip of key resulted in different value"))
	}
	return &rv, nil
}

// JSON returns the JSON representation of e, or errors
func (e *EncryptedVRFKey) JSON() ([]byte, error) {
	keyJSON, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal encrypted key to JSON")
	}
	return keyJSON, nil
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func (e *EncryptedVRFKey) Decrypt(auth string) (*PrivateKey, error) {
	keyJSON, err := json.Marshal(e.VRFKey)
	if err != nil {
		return nil, errors.Wrapf(err, "while marshaling key for decryption")
	}
	gethKey, err := keystore.DecryptKey(keyJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt key %s",
			e.PublicKey.String())
	}
	return fromGethKey(gethKey), nil
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

// MarshalText renders k as a text string
func (k PublicKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText reads a PublicKey into k from text, or errors
func (k *PublicKey) UnmarshalText(text []byte) error {
	if err := k.SetFromHex(string(text)); err != nil {
		return errors.Wrapf(err, "while parsing %s as public key", text)
	}
	return nil
}

// Value marshals PublicKey to be saved in the DB
func (k PublicKey) Value() (driver.Value, error) {
	return k.String(), nil
}

// Scan reconstructs a PublicKey from a DB record of it.
func (k *PublicKey) Scan(value interface{}) error {
	rawKey, ok := value.(string)
	if !ok {
		return errors.Wrap(fmt.Errorf("unable to convert %+v of type %T to PublicKey", value, value), "scan failure")
	}
	if err := k.SetFromHex(rawKey); err != nil {
		return errors.Wrapf(err, "while scanning %s as PublicKey", rawKey)
	}
	return nil
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
	// With sqlite gorm driver, we get a []byte, here. With postgres, a string!
	// https://github.com/jinzhu/gorm/issues/2276
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
