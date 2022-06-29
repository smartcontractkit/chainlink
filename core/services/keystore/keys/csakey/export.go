package csakey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "CSA"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	return keys.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ keys.EncryptedKeyExport, rawPrivKey []byte) (KeyV2, error) {
			return Raw(rawPrivKey).Key(), nil
		},
	)
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key KeyV2, cryptoJSON keystore.CryptoJSON) (keys.EncryptedKeyExport, error) {
			return keys.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: key.PublicKeyString(),
				Crypto:    cryptoJSON,
			}, nil
		},
	)
}

func adulteratedPassword(password string) string {
	return "csakey" + password
}
