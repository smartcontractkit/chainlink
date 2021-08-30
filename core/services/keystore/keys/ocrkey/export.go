package ocrkey

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "OCR"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	var export EncryptedOCRKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return KeyV2{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return KeyV2{}, errors.Wrap(err, "failed to decrypt OCR key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

type EncryptedOCRKeyExport struct {
	KeyType               string                `json:"keyType"`
	ID                    string                `json:"id"`
	OnChainSigningAddress OnChainSigningAddress `json:"onChainSigningAddress"`
	OffChainPublicKey     OffChainPublicKey     `json:"offChainPublicKey"`
	ConfigPublicKey       ConfigPublicKey       `json:"configPublicKey"`
	Crypto                keystore.CryptoJSON   `json:"crypto"`
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt Eth key")
	}
	encryptedOCRKExport := EncryptedOCRKeyExport{
		KeyType:               keyTypeIdentifier,
		ID:                    key.ID(),
		OnChainSigningAddress: key.OnChainSigning.Address(),
		OffChainPublicKey:     key.OffChainSigning.PublicKey(),
		ConfigPublicKey:       key.PublicKeyConfig(),
		Crypto:                cryptoJSON,
	}
	return json.Marshal(encryptedOCRKExport)
}

func adulteratedPassword(password string) string {
	return "ocrkey" + password
}
