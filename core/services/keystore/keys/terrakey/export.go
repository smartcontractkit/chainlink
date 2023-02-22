package terrakey

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "Terra"

// FromEncryptedJSON gets key from json and password
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

// ToEncryptedJSON returns encrypted JSON representing key
func (key Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key Key, cryptoJSON keystore.CryptoJSON) (keys.EncryptedKeyExport, error) {
			return keys.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: hex.EncodeToString(key.PublicKey().Bytes()),
				Crypto:    cryptoJSON,
			}, nil
		},
	)
}

func adulteratedPassword(password string) string {
	return "terrakey" + password
}
