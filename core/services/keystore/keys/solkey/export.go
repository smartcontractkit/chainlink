package solkey

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	solkeys "github.com/smartcontractkit/chainlink-solana/pkg/solana/keys"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "Solana"

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON(keyJSON []byte, password string) (solkeys.Key, error) {
	return keys.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ keys.EncryptedKeyExport, rawPrivKey []byte) (solkeys.Key, error) {
			return solkeys.Raw(rawPrivKey).Key(), nil
		},
	)
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON(key solkeys.Key, password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key solkeys.Key, cryptoJSON keystore.CryptoJSON) (keys.EncryptedKeyExport, error) {
			return keys.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: hex.EncodeToString(key.GetPublic()),
				Crypto:    cryptoJSON,
			}, nil
		},
	)
}

func adulteratedPassword(password string) string {
	return "solkey" + password
}
