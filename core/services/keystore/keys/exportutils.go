package keys

import (
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// EncryptedKeyExport represents a chain specific encrypted key
type EncryptedKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON[K any](
	identifier string,
	keyJSON []byte,
	password string,
	passwordFunc func(string) string,
	rawKeyToKey func([]byte) K,
) (K, error) {
	var export EncryptedKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return *new(K), err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, passwordFunc(password))
	if err != nil {
		return *new(K), errors.Wrapf(err, "failed to decrypt %s key", identifier)
	}
	key := rawKeyToKey(privKey)
	return key, nil
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON(
	identifier string,
	raw []byte,
	pubkey string,
	password string,
	scryptParams utils.ScryptParams,
	passwordFunc func(string) string,
) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		raw,
		[]byte(passwordFunc(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt %s key", identifier)
	}
	encryptedKeyExport := EncryptedKeyExport{
		KeyType:   identifier,
		PublicKey: pubkey,
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedKeyExport)
}
