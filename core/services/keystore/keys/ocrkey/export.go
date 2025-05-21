package ocrkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const keyTypeIdentifier = "OCR"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	return internal.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ EncryptedOCRKeyExport, rawPrivKey internal.Raw) (KeyV2, error) {
			return KeyFor(rawPrivKey), nil
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
	return internal.ToEncryptedJSON(
		keyTypeIdentifier,
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
