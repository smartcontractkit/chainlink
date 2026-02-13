package workflowkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

const keyTypeIdentifier = "Workflow"

func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	return internal.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ internal.EncryptedKeyExport, rawPrivKey internal.Raw) (Key, error) {
			return KeyFor(rawPrivKey), nil
		},
	)
}

func (k Key) ToEncryptedJSON(password string, scryptParams commonkeystore.ScryptParams) (export []byte, err error) {
	return internal.ToEncryptedJSON(
		keyTypeIdentifier,
		k,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key Key, cryptoJSON keystore.CryptoJSON) internal.EncryptedKeyExport {
			return internal.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: key.PublicKeyString(),
				Crypto:    cryptoJSON,
			}
		},
	)
}

func adulteratedPassword(password string) string {
	return "workflowkey" + password
}
