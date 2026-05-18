package ocrkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const keyTypeIdentifier = "OCR"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	return keys.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ EncryptedOCRKeyExport, rawPrivKey []byte) (KeyV2, error) {
			return Raw(rawPrivKey).Key(), nil
		},
	)
}

type EncryptedOCRKeyExport struct {
	KeyType               string                `json:"keyType"`
	ID                    string                `json:"id"`
	OnChainSigningAddress OnChainSigningAddress `json:"onChainSigningAddress"`
	OffChainPublicKey     OffChainPublicKey     `json:"offChainPublicKey"`
	ConfigPublicKey       ConfigPublicKey       `json:"configPublicKey"`
	Crypto                keystore.CryptoJSON   `json:"crypto"`
}

func (x EncryptedOCRKeyExport) GetCrypto() keystore.CryptoJSON {
	return x.Crypto
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key KeyV2, cryptoJSON keystore.CryptoJSON) EncryptedOCRKeyExport {
			return EncryptedOCRKeyExport{
				KeyType:               id,
				ID:                    key.ID(),
				OnChainSigningAddress: key.OnChainSigning.Address(),
				OffChainPublicKey:     key.OffChainSigning.PublicKey(),
				ConfigPublicKey:       key.PublicKeyConfig(),
				Crypto:                cryptoJSON,
			}
		},
	)
}

func adulteratedPassword(password string) string {
	return "ocrkey" + password
}
