package solkey

import (
	"encoding/hex"
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "Solana"

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	var export EncryptedSolanaKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return Key{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return Key{}, errors.Wrap(err, "failed to decrypt Solana key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

// EncryptedSolanaKeyExport represents the Solana encrypted key
type EncryptedSolanaKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

// ToEncryptedJSON returns encrypted JSON representing key
func (key Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt Solana key")
	}
	encryptedSolanaKeyExport := EncryptedSolanaKeyExport{
		KeyType:   keyTypeIdentifier,
		PublicKey: hex.EncodeToString(key.pubKey),
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedSolanaKeyExport)
}

func adulteratedPassword(password string) string {
	return "solkey" + password
}
