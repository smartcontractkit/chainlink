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
	// terraKeyBundle represents the bundle of keys needed for OCR
	terraKeyBundle struct {
		keyBundleBase
		terraKeyring
	}

	terraKeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		TerraKeyring    []byte
	}
)

var _ KeyBundle = &terraKeyBundle{}

func newTerraKeyBundle() (*terraKeyBundle, error) {
	return newTerraKeyBundleFrom(cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func mustNewTerraKeyBundleInsecure(reader io.Reader) *terraKeyBundle {
	key, err := newTerraKeyBundleFrom(reader, reader, reader)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate new OCR2-Terra Key"))
	}
	return key
}

func newTerraKeyBundleFrom(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*terraKeyBundle, error) {
	offchainKeyring, err := newOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	terraKeyRing, err := newTerraKeyring(onchainSigningKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := terraKeyBundle{
		keyBundleBase: keyBundleBase{
			chainType:       chaintype.Terra,
			OffchainKeyring: *offchainKeyring,
		},
		terraKeyring: *terraKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return &k, nil
}

func mustNewTerraKeyFromRaw(raw []byte) terraKeyBundle {
	// offchain private key 64 bytes || offchain encryption key 32 bytes || onchain 32 bytes private key
	var kb terraKeyBundle
	err := kb.Unmarshal(raw)
	if err != nil {
		panic(err)
	}
	return kb
}

// OnChainPublicKey returns public component of the keypair used on chain
func (kb *terraKeyBundle) OnChainPublicKey() string {
	return hex.EncodeToString(kb.terraKeyring.PublicKey())
}

func (kb *terraKeyBundle) Marshal() ([]byte, error) {
	offchainKeyringBytes, err := kb.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	terraKeyringBytes, err := kb.terraKeyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := terraKeyBundleRawData{
		ChainType:       kb.chainType,
		OffchainKeyring: offchainKeyringBytes,
		TerraKeyring:    terraKeyringBytes,
	}
	return json.Marshal(&rawKeyData)
}

func (kb *terraKeyBundle) Unmarshal(b []byte) (err error) {
	var rawKeyData terraKeyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	err = kb.terraKeyring.unmarshal(rawKeyData.TerraKeyring)
	if err != nil {
		return err
	}
	kb.chainType = rawKeyData.ChainType
	kb.id = sha256.Sum256(b)
	return nil
}

func (kb *terraKeyBundle) Raw() Raw {
	b, err := kb.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}
