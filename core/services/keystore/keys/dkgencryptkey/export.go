package dkgencryptkey

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "DKGEncrypt"

// FromEncryptedJSON returns a dkgencryptkey.KeyV2 from encrypted data in go-ethereum keystore format.
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	var export EncryptedDKGEncryptKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return Key{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return Key{}, errors.Wrap(err, "failed to decrypt DKGEncrypt key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

// EncryptedDKGEncryptKeyExport is an encrypted exported DKGEncrypt key
// that contains the DKGEncrypt key in go-ethereum keystore format.
type EncryptedDKGEncryptKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

// ToEncryptedJSON exports this key into a JSON object following the format of EncryptedDKGEncryptKeyExport
func (k Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		k.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt DKGEncrypt key")
	}
	encryptedDKGEncryptKeyExport := EncryptedDKGEncryptKeyExport{
		KeyType:   keyTypeIdentifier,
		PublicKey: k.PublicKeyString(),
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedDKGEncryptKeyExport)
}

func adulteratedPassword(password string) string {
	return "dkgencryptkey" + password
}
