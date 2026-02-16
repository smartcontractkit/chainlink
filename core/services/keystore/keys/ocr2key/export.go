package ocr2key

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

const keyTypeIdentifier = "OCR2"

// EncryptedOCRKeyExport represents encrypted OCR key export
type EncryptedOCRKeyExport struct {
	KeyType           string              `json:"keyType"`
	ChainType         corekeys.ChainType  `json:"chainType"`
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
	return internal.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(export EncryptedOCRKeyExport, rawPrivKey internal.Raw) (KeyBundle, error) {
			var kb KeyBundle
			switch export.ChainType {
			case corekeys.EVM:
				kb = newKeyBundle(new(evmKeyring))
			case corekeys.Cosmos:
				kb = newKeyBundle(new(cosmosKeyring))
			case corekeys.Solana:
				kb = newKeyBundle(new(solanaKeyring))
			case corekeys.StarkNet:
				kb = newKeyBundle(new(starkkey.OCR2Key))
			case corekeys.Aptos:
				kb = newKeyBundle(new(ed25519Keyring))
			case corekeys.Tron:
				kb = newKeyBundle(new(evmKeyring))
			case corekeys.TON:
				kb = newKeyBundle(new(tonKeyring))
			case corekeys.Sui:
				kb = newKeyBundle(new(ed25519Keyring))
			default:
				return nil, corekeys.NewErrInvalidChainType(export.ChainType)
			}
			if err := kb.Unmarshal(internal.Bytes(rawPrivKey)); err != nil {
				return nil, err
			}
			return kb, nil
		},
	)
}

// ToEncryptedJSON returns encrypted JSON representing key
func ToEncryptedJSON(key KeyBundle, password string, scryptParams commonkeystore.ScryptParams) (export []byte, err error) {
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
