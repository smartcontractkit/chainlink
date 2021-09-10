package ocr2key

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "OCR2"

type EncryptedOCR2KeyExport struct {
	ID     models.Sha256Hash   `json:"id" gorm:"primary_key"`
	Crypto keystore.CryptoJSON `json:"crypto"`
}

func FromEncryptedJSON(keyJSON []byte, password string) (*KeyBundle, error) {
	var export EncryptedOCR2KeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return nil, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt OCR key")
	}
	key := Raw(privKey).Key()
	return &key, nil
}

func (key *KeyBundle) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	marshalledPrivK, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}
	cryptoJSON, err := keystore.EncryptDataV3(
		marshalledPrivK,
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt OCR2 key")
	}

	encryptedOCR2KExport := EncryptedOCR2KeyExport{
		ID:     key.id,
		Crypto: cryptoJSON,
	}
	return json.Marshal(encryptedOCR2KExport)
}
