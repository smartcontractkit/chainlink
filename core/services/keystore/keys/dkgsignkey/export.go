package dkgsignkey

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "DKGSign"

// FromEncryptedJSON returns a dkgsignkey.Key from encrypted data in go-ethereum keystore format.
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	var export EncryptedDKGSignKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return Key{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return Key{}, errors.Wrap(err, "failed to decrypt DKGSign key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

// EncryptedDKGSignKeyExport is an encrypted exported DKGSign
// that contains the DKGSign key in go-ethereum keystore format.
type EncryptedDKGSignKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

// ToEncryptedJSON exports this key into a JSON object following the format of EncryptedDKGSignKeyExport
func (key Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt DKGSign key")
	}
	encryptedDKGSignKeyExport := EncryptedDKGSignKeyExport{
		KeyType:   keyTypeIdentifier,
		PublicKey: key.PublicKeyString(),
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedDKGSignKeyExport)
}

func adulteratedPassword(password string) string {
	return "dkgsignkey" + password
}
