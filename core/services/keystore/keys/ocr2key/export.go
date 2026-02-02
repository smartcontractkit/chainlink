package ocr2key

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

func (x EncryptedOCRKeyExport) GetCrypto() keystore.CryptoJSON {
	return x.Crypto
}

// FromEncryptedJSON returns key from encrypted json
func FromEncryptedJSON(keyJSON []byte, password string) (KeyBundle, error) {
	ocrBundle, err := corekeys.FromEncryptedOCRKeyBundle(keyJSON, password)
	switch {
	case errors.Is(err, corekeys.ErrInvalidExportFormat):
		return internal.FromEncryptedJSON(
			keyTypeIdentifier,
			keyJSON,
			password,
			adulteratedPassword,
			func(export EncryptedOCRKeyExport, rawPrivKey internal.Raw) (KeyBundle, error) {
				return unmarshalKeyBundle(export.ChainType, internal.Bytes(rawPrivKey))
			},
		)
	case err != nil:
		return nil, err
	}

	coreChainType := chaintype.ChainType(ocrBundle.ChainType)

	rawData := keyBundleRawData{
		ChainType:       coreChainType,
		OffchainKeyring: append(ocrBundle.OffchainSigningKey, ocrBundle.OffchainEncryptionKey...),
		Keyring:         ocrBundle.OnchainSigningKey,
	}
	rawBytes, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}

	return unmarshalKeyBundle(coreChainType, rawBytes)
}

func unmarshalKeyBundle(chainType chaintype.ChainType, rawKeyData []byte) (KeyBundle, error) {
	var kb KeyBundle
	switch chainType {
	case chaintype.EVM:
		kb = newKeyBundle(new(evmKeyring))
	case chaintype.Cosmos:
		kb = newKeyBundle(new(cosmosKeyring))
	case chaintype.Solana:
		kb = newKeyBundle(new(solanaKeyring))
	case chaintype.StarkNet:
		kb = newKeyBundle(new(starkkey.OCR2Key))
	case chaintype.Aptos:
		kb = newKeyBundle(new(ed25519Keyring))
	case chaintype.Tron:
		kb = newKeyBundle(new(evmKeyring))
	case chaintype.TON:
		kb = newKeyBundle(new(tonKeyring))
	case chaintype.Sui:
		kb = newKeyBundle(new(ed25519Keyring))
	default:
		return nil, chaintype.NewErrInvalidChainType(chainType)
	}
	if err := kb.Unmarshal(rawKeyData); err != nil {
		return nil, err
	}
	return kb, nil
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON(key KeyBundle, password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return internal.ToEncryptedJSON(
		keyTypeIdentifier,
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key KeyBundle, cryptoJSON keystore.CryptoJSON) EncryptedOCRKeyExport {
			pubKeyConfig := key.ConfigEncryptionPublicKey()
			pubKey := key.OffchainPublicKey()
			return EncryptedOCRKeyExport{
				KeyType:           id,
				ChainType:         key.ChainType(),
				ID:                key.ID(),
				OnchainPublicKey:  key.OnChainPublicKey(),
				OffChainPublicKey: hex.EncodeToString(pubKey[:]),
				ConfigPublicKey:   hex.EncodeToString(pubKeyConfig[:]),
				Crypto:            cryptoJSON,
			}
		},
	)
}
