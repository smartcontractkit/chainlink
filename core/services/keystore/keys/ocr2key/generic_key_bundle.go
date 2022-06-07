package ocr2key

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type (
	keyring interface {
		ocrtypes.OnchainKeyring
		marshal() ([]byte, error)
		unmarshal(in []byte) error
	}

	keyBundle[K keyring] struct {
		keyBundleBase
		keyring K
	}

	keyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		Keyring         []byte
	}
)

var _ KeyBundle = &keyBundle[*starknetKeyring]{}

func newKeyBundle[K keyring](chain chaintype.ChainType, newKeyring func(material io.Reader) (K, error)) (*keyBundle[K], error) {
	return newKeyBundleFrom(chain, newKeyring, cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func mustNewKeyBundleInsecure[K keyring](chain chaintype.ChainType, newKeyring func(material io.Reader) (K, error), reader io.Reader) *keyBundle[K] {
	key, err := newKeyBundleFrom(chain, newKeyring, reader, reader, reader)
	if err != nil {
		panic(errors.Wrapf(err, "failed to generate new OCR2-%s Key", chain))
	}
	return key
}

func newKeyBundleFrom[K keyring](chain chaintype.ChainType, newKeyring func(material io.Reader) (K, error), onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*keyBundle[K], error) {
	offchainKeyring, err := newOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	kr, err := newKeyring(onchainSigningKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := keyBundle[K]{
		keyBundleBase: keyBundleBase{
			chainType:       chain,
			OffchainKeyring: *offchainKeyring,
		},
		keyring: kr,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return &k, nil
}

func mustNewKeyFromRaw[K keyring](raw []byte, key K) keyBundle[K] {
	// offchain private key 64 bytes || offchain encryption key 32 bytes || onchain 32 bytes private key
	var kb keyBundle[K]
	kb.keyring = key
	err := kb.Unmarshal(raw)
	if err != nil {
		panic(err)
	}
	return kb
}

func (kb *keyBundle[K]) MaxSignatureLength() int {
	return kb.keyring.MaxSignatureLength()
}

func (kb *keyBundle[K]) PublicKey() ocrtypes.OnchainPublicKey {
	return kb.keyring.PublicKey()
}

func (kb *keyBundle[K]) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	return kb.keyring.Sign(reportCtx, report)
}

func (kb *keyBundle[K]) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	return kb.keyring.Verify(publicKey, reportCtx, report, signature)
}

// OnChainPublicKey returns public component of the keypair used on chain
func (kb *keyBundle[K]) OnChainPublicKey() string {
	return hex.EncodeToString(kb.keyring.PublicKey())
}

func (kb *keyBundle[K]) Marshal() ([]byte, error) {
	offchainKeyringBytes, err := kb.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	keyringBytes, err := kb.keyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := keyBundleRawData{
		ChainType:       kb.chainType,
		OffchainKeyring: offchainKeyringBytes,
		Keyring:         keyringBytes,
	}
	return json.Marshal(&rawKeyData)
}

func (kb *keyBundle[K]) Unmarshal(b []byte) (err error) {
	var rawKeyData keyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	err = kb.keyring.unmarshal(rawKeyData.Keyring)
	if err != nil {
		return err
	}
	kb.chainType = rawKeyData.ChainType
	kb.id = sha256.Sum256(b)
	return nil
}

func (kb *keyBundle[K]) Raw() Raw {
	b, err := kb.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}
