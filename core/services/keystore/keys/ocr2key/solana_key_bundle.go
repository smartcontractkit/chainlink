package ocr2key

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
)

type (
	// solanaKeyBundle represents the bundle of keys needed for OCR
	solanaKeyBundle struct {
		keyBundleBase
		solanaKeyring
	}

	solanaKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		SolanaKeyring   []byte
	}
)

var _ KeyBundle = &solanaKeyBundle{}

// New makes a new set of OCR key bundles from cryptographically secure entropy
func newSolanaKeyBundle() (*solanaKeyBundle, error) {
	return newSolanaKeyBundleFrom(cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func mustNewSolanaKeyBundleInsecure(reader io.Reader) *solanaKeyBundle {
	key, err := newSolanaKeyBundleFrom(reader, reader, reader)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate new OCR2-Solana Key"))
	}
	return key
}

func newSolanaKeyBundleFrom(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*solanaKeyBundle, error) {
	offchainKeyring, err := newOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	solanaKeyRing, err := newSolanaKeyring(onchainSigningKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := solanaKeyBundle{
		keyBundleBase: keyBundleBase{
			chainType:       chaintype.Solana,
			OffchainKeyring: *offchainKeyring,
		},
		solanaKeyring: *solanaKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return &k, nil
}

func mustNewSolanaKeyFromRaw(raw []byte) solanaKeyBundle {
	// offchain private key 64 bytes || offchain encryption key 32 bytes || onchain 32 bytes private key
	var kb solanaKeyBundle
	err := kb.Unmarshal(raw)
	if err != nil {
		panic(err)
	}
	return kb
}

// OnChainPublicKey returns public component of the keypair used on chain
func (kb *solanaKeyBundle) OnChainPublicKey() string {
	return hex.EncodeToString(kb.solanaKeyring.PublicKey())
}

func (kb *solanaKeyBundle) Marshal() ([]byte, error) {
	offchainKeyringBytes, err := kb.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	solanaKeyringBytes, err := kb.solanaKeyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := solanaKeyBundleRawData{
		ChainType:       kb.chainType,
		OffchainKeyring: offchainKeyringBytes,
		SolanaKeyring:   solanaKeyringBytes,
	}
	return json.Marshal(&rawKeyData)
}

func (kb *solanaKeyBundle) Unmarshal(b []byte) (err error) {
	var rawKeyData solanaKeyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	err = kb.solanaKeyring.unmarshal(rawKeyData.SolanaKeyring)
	if err != nil {
		return err
	}
	kb.chainType = rawKeyData.ChainType
	kb.id = sha256.Sum256(b)
	return nil
}

func (kb *solanaKeyBundle) Raw() Raw {
	b, err := kb.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}
