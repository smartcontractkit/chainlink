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
	// evmKeyBundle represents the bundle of keys needed for OCR
	evmKeyBundle struct {
		keyBundleBase
		evmKeyring
	}

	evmKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		EVMKeyring      []byte
	}
)

var _ KeyBundle = &evmKeyBundle{}

func newEVMKeyBundle() (*evmKeyBundle, error) {
	return newEVMKeyBundleFrom(cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func mustNewEVMKeyBundleInsecure(reader io.Reader) *evmKeyBundle {
	key, err := newEVMKeyBundleFrom(reader, reader, reader)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate new OCR2-EVM Key"))
	}
	return key
}

func newEVMKeyBundleFrom(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*evmKeyBundle, error) {
	offchainKeyring, err := newOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	evmKeyRing, err := newEVMKeyring(onchainSigningKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := evmKeyBundle{
		keyBundleBase: keyBundleBase{
			chainType:       chaintype.EVM,
			OffchainKeyring: *offchainKeyring,
		},
		evmKeyring: *evmKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return &k, nil
}

func mustNewEVMKeyFromRaw(raw []byte) evmKeyBundle {
	// offchain private key 64 bytes || offchain encryption key 32 bytes || onchain 32 bytes private key
	var kb evmKeyBundle
	err := kb.Unmarshal(raw)
	if err != nil {
		panic(err)
	}
	return kb
}

// OnChainPublicKey returns public component of the keypair used on chain
func (kb *evmKeyBundle) OnChainPublicKey() string {
	return hex.EncodeToString(kb.evmKeyring.PublicKey())
}

func (kb *evmKeyBundle) Marshal() ([]byte, error) {
	offchainKeyringBytes, err := kb.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	evmKeyringBytes, err := kb.evmKeyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := evmKeyBundleRawData{
		ChainType:       kb.chainType,
		OffchainKeyring: offchainKeyringBytes,
		EVMKeyring:      evmKeyringBytes,
	}
	return json.Marshal(&rawKeyData)
}

func (kb *evmKeyBundle) Unmarshal(b []byte) (err error) {
	var rawKeyData evmKeyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	err = kb.evmKeyring.unmarshal(rawKeyData.EVMKeyring)
	if err != nil {
		return err
	}
	kb.chainType = rawKeyData.ChainType
	kb.id = sha256.Sum256(b)
	return nil
}

func (kb *evmKeyBundle) Raw() Raw {
	b, err := kb.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}
