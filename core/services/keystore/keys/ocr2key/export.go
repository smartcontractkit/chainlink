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

// EncryptedOCRKeyExport represents encrypted OCR key export
type EncryptedOCRKeyExport struct {
	KeyType           string              `json:"keyType"`
	ChainType         chaintype.ChainType `json:"chainType"`
	ID                string              `json:"id"`
	OnchainPublicKey  string              `json:"onchainPublicKey"`
	OffChainPublicKey string              `json:"offchainPublicKey"`
	ConfigPublicKey   string              `json:"configPublicKey"`
	Crypto            keystore.CryptoJSON `json:"crypto"`
}

// FromEncryptedJSON returns key from encrypted json
func FromEncryptedJSON(keyJSON []byte, password string) (KeyBundle, error) {
	var export EncryptedOCRKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return nil, err
	}
	rawKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt OCR key")
	}
	var kb KeyBundle
	switch export.ChainType {
	case chaintype.EVM:
		kb = newKeyBundle(new(evmKeyring))
	case chaintype.Solana:
		kb = newKeyBundle(new(solanaKeyring))
	case chaintype.Terra:
		kb = newKeyBundle(new(terraKeyring))
	case chaintype.Starknet:
		kb = newKeyBundle(new(starknetKeyring))
	default:
		return nil, chaintype.NewErrInvalidChainType(export.ChainType)
	}
	if err := kb.Unmarshal(rawKey); err != nil {
		return nil, err
	}
	return kb, nil
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON(key KeyBundle, password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt Eth key")
	}
	pubKeyConfig := key.ConfigEncryptionPublicKey()
	pubKey := key.OffchainPublicKey()
	encryptedOCRKExport := EncryptedOCRKeyExport{
		KeyType:           keyTypeIdentifier,
		ChainType:         key.ChainType(),
		ID:                key.ID(),
		OnchainPublicKey:  key.OnChainPublicKey(),
		OffChainPublicKey: hex.EncodeToString(pubKey[:]),
		ConfigPublicKey:   hex.EncodeToString(pubKeyConfig[:]),
		Crypto:            cryptoJSON,
	}
	return json.Marshal(encryptedOCRKExport)
}
