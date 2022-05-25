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
	// starknetKeyBundle represents the bundle of keys needed for OCR
	starknetKeyBundle struct {
		keyBundleBase
		starknetKeyring
	}

	starknetKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		StarknetKeyring []byte
	}
)

var _ KeyBundle = &starknetKeyBundle{}

func newStarknetKeyBundle() (*starknetKeyBundle, error) {
	return newStarknetKeyBundleFrom(cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func mustNewStarknetKeyBundleInsecure(reader io.Reader) *starknetKeyBundle {
	key, err := newStarknetKeyBundleFrom(reader, reader, reader)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate new OCR2-Starknet Key"))
	}
	return key
}

func newStarknetKeyBundleFrom(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*starknetKeyBundle, error) {
	offchainKeyring, err := newOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	starknetKeyRing, err := newStarknetKeyring(onchainSigningKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := starknetKeyBundle{
		keyBundleBase: keyBundleBase{
			chainType:       chaintype.Starknet,
			OffchainKeyring: *offchainKeyring,
		},
		starknetKeyring: *starknetKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return &k, nil
}

func mustNewStarknetKeyFromRaw(raw []byte) starknetKeyBundle {
	// offchain private key 64 bytes || offchain encryption key 32 bytes || onchain 32 bytes private key
	var kb starknetKeyBundle
	err := kb.Unmarshal(raw)
	if err != nil {
		panic(err)
	}
	return kb
}

// OnChainPublicKey returns public component of the keypair used on chain
func (kb *starknetKeyBundle) OnChainPublicKey() string {
	return hex.EncodeToString(kb.starknetKeyring.PublicKey())
}

func (kb *starknetKeyBundle) Marshal() ([]byte, error) {
	offchainKeyringBytes, err := kb.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	starknetKeyringBytes, err := kb.starknetKeyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := starknetKeyBundleRawData{
		ChainType:       kb.chainType,
		OffchainKeyring: offchainKeyringBytes,
		StarknetKeyring: starknetKeyringBytes,
	}
	return json.Marshal(&rawKeyData)
}

func (kb *starknetKeyBundle) Unmarshal(b []byte) (err error) {
	var rawKeyData starknetKeyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	err = kb.starknetKeyring.unmarshal(rawKeyData.StarknetKeyring)
	if err != nil {
		return err
	}
	kb.chainType = rawKeyData.ChainType
	kb.id = sha256.Sum256(b)
	return nil
}

func (kb *starknetKeyBundle) Raw() Raw {
	b, err := kb.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}
