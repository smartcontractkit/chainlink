package internal

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

// FromEncryptedJSON gets key [K] from keyJSON [E] and password
func FromEncryptedJSON[E Encrypted, K any](
	identifier string,
	keyJSON []byte,
	password string,
	passwordFunc func(string) string,
	privKeyToKey func(export E, rawPrivKey Raw) (K, error),
) (K, error) {
	// unmarshal byte data to [E] Encrypted key export
	var export E
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return *new(K), err
	}

	// decrypt data using prefixed password
	privKey, err := keystore.DecryptDataV3(export.GetCrypto(), passwordFunc(password))
	if err != nil {
		return *new(K), fmt.Errorf("failed to decrypt %s key: %w", identifier, err)
	}

	// convert unmarshalled data and decrypted key to [K] key format
	key, err := privKeyToKey(export, NewRaw(privKey))
	if err != nil {
		return *new(K), fmt.Errorf("failed to convert %s key to key bundle: %w", identifier, err)
	}

	return key, nil
}

type Key interface {
	Raw() Raw
}

// ToEncryptedJSON returns encrypted JSON [E] representing key [K]
func ToEncryptedJSON[E Encrypted, K Key](
	identifier string,
	key K,
	password string,
	scryptParams utils.ScryptParams,
	passwordFunc func(string) string,
	buildExport func(id string, key K, cryptoJSON keystore.CryptoJSON) E,
) (export []byte, err error) {
	// encrypt data using prefixed password
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw().bytes,
		[]byte(passwordFunc(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt %s key: %w", identifier, err)
	}

	// build [E] export struct using encrypted key, identifier, and original key [K]
	encryptedKeyExport := buildExport(identifier, key, cryptoJSON)

	return json.Marshal(encryptedKeyExport)
}
