package vrf_key

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
)

var passwordPrefix = "don't mix VRF and Ethereum keys!"

func adulteratedPassword(auth string) string {
	return passwordPrefix + auth
}

type ScryptParams struct{ N, P int }

var defaultScryptParams = ScryptParams{
	N: keystore.StandardScryptN, P: keystore.StandardScryptP}

// FastScryptParams is for use in tests, where you don't want to wear out your
// CPU with expensive key derivations
var FastScryptParams = ScryptParams{N: 2, P: 1}

// Encrypt returns the key encrypted with passphrase auth
func (k *PrivateKey) Encrypt(auth string, p ...ScryptParams,
) (*EncryptedSecretKey, error) {
	var scryptParams ScryptParams
	switch len(p) {
	case 0:
		scryptParams = defaultScryptParams
	case 1:
		scryptParams = p[0]
	default:
		return nil, fmt.Errorf("can take at most one set of ScryptParams")
	}
	keyJSON, err := keystore.EncryptKey(k.GethKey(), adulteratedPassword(auth),
		scryptParams.N, scryptParams.P)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt vrf key")
	}
	rv := EncryptedSecretKey{}
	if err := json.Unmarshal(keyJSON, &rv.VRFKey); err != nil {
		return nil, errors.Wrapf(err, "geth returned unexpected key material")
	}
	rv.PublicKey = k.PublicKey
	roundTripKey, err := rv.Decrypt(auth)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt just-encrypted key!")
	}
	if !roundTripKey.k.Equal(k.k) || roundTripKey.PublicKey != k.PublicKey ||
		!bytes.Equal(roundTripKey.ID, k.ID) {
		panic(fmt.Errorf("roundtrip of key resulted in different value"))
	}
	return &rv, nil
}

func (e *EncryptedSecretKey) JSON() ([]byte, error) {
	keyJSON, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal encrypted key to JSON")
	}
	return keyJSON, nil
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func (e *EncryptedSecretKey) Decrypt(auth string) (*PrivateKey, error) {
	keyJSON, err := json.Marshal(e.VRFKey)
	if err != nil {
		return nil, errors.Wrapf(err, "while marshaling key for decryption")
	}
	gethKey, err := keystore.DecryptKey(keyJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt key %s",
			e.PublicKey.String())
	}
	return FromGethKey(gethKey), nil
}

func (e *EncryptedSecretKey) WriteToDisk(path string) error {
	keyJSON, err := e.JSON()
	if err != nil {
		return errors.Wrapf(err, "while marshaling key to save to %s", path)
	}
	UserReadWriteOtherNoAccess := os.FileMode(0600)
	return ioutil.WriteFile(path, keyJSON, UserReadWriteOtherNoAccess)
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

func (k gethKeyStruct) Value() (driver.Value, error) {
	return json.Marshal(&k)
}

func (k *gethKeyStruct) Scan(value interface{}) error {
	rawKey, ok := value.([]byte)
	if !ok {
		return errors.Wrap(
			fmt.Errorf("unable to convert %+v of type %T to gethKeyStruct",
				value, value), "scan failure")
	}
	return json.Unmarshal(rawKey, k)
}
