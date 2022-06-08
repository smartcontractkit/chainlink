package starkkey

import (
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "StarkNet"

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	var export EncryptedStarkNetKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return Key{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return Key{}, errors.Wrap(err, "failed to decrypt StarkNet key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

// EncryptedStarkNetKeyExport represents the StarkNet encrypted key
type EncryptedStarkNetKeyExport struct {
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
		return nil, errors.Wrapf(err, "could not encrypt StarkNet key")
	}
	encryptedStarkNetKeyExport := EncryptedStarkNetKeyExport{
		KeyType:   keyTypeIdentifier,
		PublicKey: key.PublicKeyStr(),
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedStarkNetKeyExport)
}

func adulteratedPassword(password string) string {
	return "starkkey" + password
}
