package ocr2key

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "OCR2"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyBundle, error) {
	var export EncryptedOCRKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return KeyBundle{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return KeyBundle{}, errors.Wrap(err, "failed to decrypt OCR key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

type EncryptedOCRKeyExport struct {
	KeyType               string              `json:"keyType"`
	ChainType             chaintype.ChainType `json:"chainType"`
	ID                    string              `json:"id"`
	OnChainSigningAddress string              `json:"onChainSigningAddress"`
	OffChainPublicKey     string              `json:"offchainPublicKey"`
	ConfigPublicKey       string              `json:"configPublicKey"`
	Crypto                keystore.CryptoJSON `json:"crypto"`
}

func (key KeyBundle) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt Eth key")
	}
	pubKeyConfig := key.PublicKeyConfig()
	encryptedOCRKExport := EncryptedOCRKeyExport{
		KeyType:               keyTypeIdentifier,
		ChainType:             key.ChainType,
		ID:                    key.ID(),
		OnChainSigningAddress: key.PublicKeyAddressOnChain(),
		OffChainPublicKey:     hex.EncodeToString(key.PublicKeyOffChain()),
		ConfigPublicKey:       hex.EncodeToString(pubKeyConfig[:]),
		Crypto:                cryptoJSON,
	}
	return json.Marshal(encryptedOCRKExport)
}
