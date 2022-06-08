package keys

import (
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Encrypted interface {
	GetCrypto() keystore.CryptoJSON
}

// EncryptedKeyExport represents a chain specific encrypted key
type EncryptedKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

func (x EncryptedKeyExport) GetCrypto() keystore.CryptoJSON {
	return x.Crypto
}

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON[E Encrypted, K any](
	identifier string,
	keyJSON []byte,
	password string,
	passwordFunc func(string) string,
	privKeyToKey func(export E, rawPrivKey []byte) (K, error),
) (K, error) {
	var export E
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return *new(K), err
	}
	privKey, err := keystore.DecryptDataV3(export.GetCrypto(), passwordFunc(password))
	if err != nil {
		return *new(K), errors.Wrapf(err, "failed to decrypt %s key", identifier)
	}
	key, err := privKeyToKey(export, privKey)
	if err != nil {
		return *new(K), errors.Wrapf(err, "failed to convert %s key to key bundle", identifier)
	}

	return key, nil
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON[E Encrypted, K any](
	identifier string,
	raw []byte,
	key K,
	password string,
	scryptParams utils.ScryptParams,
	passwordFunc func(string) string,
	buildExport func(id string, key K, cryptoJSON keystore.CryptoJSON) (E, error),
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
	encryptedKeyExport, err := buildExport(identifier, key, cryptoJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "could not build encrypted export for %s key", identifier)
	}

	return json.Marshal(encryptedKeyExport)
}
