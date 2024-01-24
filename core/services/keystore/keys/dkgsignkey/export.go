package dkgsignkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const keyTypeIdentifier = "DKGSign"

// FromEncryptedJSON returns a dkgsignkey.Key from encrypted data in go-ethereum keystore format.
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	return keys.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ keys.EncryptedKeyExport, rawPrivKey []byte) (Key, error) {
			return Raw(rawPrivKey).Key(), nil
		},
	)
}

// ToEncryptedJSON exports this key into a JSON object following the format of EncryptedDKGSignKeyExport
func (key Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key Key, cryptoJSON keystore.CryptoJSON) keys.EncryptedKeyExport {
			return keys.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: key.PublicKeyString(),
				Crypto:    cryptoJSON,
			}
		},
	)
}

func adulteratedPassword(password string) string {
	return "dkgsignkey" + password
}
