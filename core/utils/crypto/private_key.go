package crypto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type EncryptedPrivateKey struct {
	keystore.CryptoJSON
}

func NewEncryptedPrivateKey(data []byte, passphrase string, scryptParams utils.ScryptParams) (*EncryptedPrivateKey, error) {
	cryptoJSON, err := keystore.EncryptDataV3(data, []byte(passphrase), scryptParams.N, scryptParams.P)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt key: %w", err)
	}

	return &EncryptedPrivateKey{CryptoJSON: cryptoJSON}, nil
}

// Decrypt returns the PrivateKey decrypted via auth, or an error
func (k EncryptedPrivateKey) Decrypt(passphrase string) (privkey []byte, err error) {
	privkey, err = keystore.DecryptDataV3(k.CryptoJSON, passphrase)
	if err != nil {
		return privkey, fmt.Errorf("could not decrypt private key: %w", err)
	}
	return privkey, nil
}

func (k *EncryptedPrivateKey) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &k)
}

func (k EncryptedPrivateKey) Value() (driver.Value, error) {
	return json.Marshal(k)
}
