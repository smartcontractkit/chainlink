package vrfkey

import (
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "VRF"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	var export EncryptedVRFKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return KeyV2{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return KeyV2{}, errors.Wrap(err, "failed to decrypt VRF key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

type EncryptedVRFKeyExport struct {
	KeyType string              `json:"keyType"`
	ID      string              `json:"id"`
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt VRF key")
	}
	encryptedOCRKExport := EncryptedVRFKeyExport{
		KeyType: keyTypeIdentifier,
		ID:      key.ID(),
		Address: key.PublicKey.Address().Hex(),
		Crypto:  cryptoJSON,
	}
	return json.Marshal(encryptedOCRKExport)
}

func adulteratedPassword(password string) string {
	return "vrfkey" + password
}
