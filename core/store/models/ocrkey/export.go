package ocrkey

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type EncryptedOCRKeyExport struct {
	ID                    models.Sha256Hash     `json:"id" gorm:"primary_key"`
	OnChainSigningAddress OnChainSigningAddress `json:"onChainSigningAddress"`
	OffChainPublicKey     OffChainPublicKey     `json:"offChainPublicKey"`
	ConfigPublicKey       ConfigPublicKey       `json:"configPublicKey"`
	Crypto                keystore.CryptoJSON   `json:"crypto"`
}

func (pk *KeyBundle) ToEncryptedExport(auth string, scryptParams utils.ScryptParams) (export []byte, err error) {
	marshalledPrivK, err := json.Marshal(pk)
	if err != nil {
		return nil, err
	}
	cryptoJSON, err := keystore.EncryptDataV3(
		marshalledPrivK,
		[]byte(adulteratedPassword(auth)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt OCR key")
	}

	encryptedOCRKExport := EncryptedOCRKeyExport{
		ID:                    pk.ID,
		OnChainSigningAddress: pk.onChainSigning.Address(),
		OffChainPublicKey:     pk.offChainSigning.PublicKey(),
		ConfigPublicKey:       pk.PublicKeyConfig(),
		Crypto:                cryptoJSON,
	}
	return json.Marshal(encryptedOCRKExport)
}

// DecryptPrivateKey returns the PrivateKey in export, decrypted via auth, or an error
func (export EncryptedOCRKeyExport) DecryptPrivateKey(auth string) (*KeyBundle, error) {
	marshalledPrivK, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt key %s", export.ID.String())
	}
	var pk KeyBundle
	err = json.Unmarshal(marshalledPrivK, &pk)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal OCR private key %s", export.ID.String())
	}
	return &pk, nil
}
