package ocr2key

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/keystore"

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
	return internal.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(export EncryptedOCRKeyExport, rawPrivKey internal.Raw) (KeyBundle, error) {
			var kb KeyBundle
			switch export.ChainType {
			case chaintype.EVM:
				kb = newKeyBundle(new(evmKeyring))
			case chaintype.Cosmos:
				kb = newKeyBundle(new(cosmosKeyring))
			case chaintype.Solana:
				kb = newKeyBundle(new(solanaKeyring))
			case chaintype.StarkNet:
				kb = newKeyBundle(new(starkkey.OCR2Key))
			case chaintype.Aptos:
				kb = newKeyBundle(new(aptosKeyring))
			case chaintype.Tron:
				kb = newKeyBundle(new(evmKeyring))
			default:
				return nil, chaintype.NewErrInvalidChainType(export.ChainType)
			}
			if err := kb.Unmarshal(internal.Bytes(rawPrivKey)); err != nil {
				return nil, err
			}
			return kb, nil
		},
	)
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
