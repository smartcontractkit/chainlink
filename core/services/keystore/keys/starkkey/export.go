package starkkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	stark "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "StarkNet"

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON(keyJSON []byte, password string) (stark.StarkKey, error) {
	return keys.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ keys.EncryptedKeyExport, rawPrivKey []byte) (stark.StarkKey, error) {
			return stark.StarkRaw(rawPrivKey).Key(), nil
		},
	)
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON(key stark.StarkKey, password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key stark.StarkKey, cryptoJSON keystore.CryptoJSON) (keys.EncryptedKeyExport, error) {
			return keys.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: key.PublicKeyStr(),
				Crypto:    cryptoJSON,
			}, nil
		},
	)
}

func adulteratedPassword(password string) string {
	return "starkkey" + password
}
