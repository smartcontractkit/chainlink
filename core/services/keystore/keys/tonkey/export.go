package tonkey

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const keyTypeIdentifier = "TON"

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	builder := func(_ internal.EncryptedKeyExport, rawPrivKey internal.Raw) (Key, error) {
		return KeyFor(rawPrivKey), nil
	}
	return internal.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		builder,
	)
}

// ToEncryptedJSON returns encrypted JSON representing key
func (key Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) ([]byte, error) {
	exporter := func(id string, key Key, cryptoJSON keystore.CryptoJSON) internal.EncryptedKeyExport {
		return internal.EncryptedKeyExport{
			KeyType:   id,
			PublicKey: hex.EncodeToString(key.pubKey),
			Crypto:    cryptoJSON,
		}
	}
	return internal.ToEncryptedJSON(
		keyTypeIdentifier,
		key,
		password,
		scryptParams,
		adulteratedPassword,
		exporter,
	)
}

func adulteratedPassword(password string) string {
	return "tonkey" + password
}
