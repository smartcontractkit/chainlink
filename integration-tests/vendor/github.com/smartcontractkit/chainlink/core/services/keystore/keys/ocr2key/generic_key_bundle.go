package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	keyring interface {
		ocrtypes.OnchainKeyring
		Marshal() ([]byte, error)
		Unmarshal(in []byte) error
	}

	keyBundle[K keyring] struct {
		keyBundleBase
		keyring K
	}

	keyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		Keyring         []byte
		ID              models.Sha256Hash // tracked to preserve bundle ID in case of migrations

		// old chain specific format for migrating
		EVMKeyring    []byte `json:",omitempty"`
		SolanaKeyring []byte `json:",omitempty"`
	}
)

func newKeyBundle[K keyring](key K) *keyBundle[K] {
	return &keyBundle[K]{keyring: key}
}

func newKeyBundleRand[K keyring](chain chaintype.ChainType, newKeyring func(material io.Reader) (K, error)) (*keyBundle[K], error) {
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
	keyringBytes, err := kb.keyring.Marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := keyBundleRawData{
		ChainType:       kb.chainType,
		OffchainKeyring: offchainKeyringBytes,
		Keyring:         keyringBytes,
		ID:              kb.id, // preserve bundle ID
	}
	return json.Marshal(&rawKeyData)
}

func (kb *keyBundle[K]) Unmarshal(b []byte) (err error) {
	var rawKeyData keyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	if err = rawKeyData.Migrate(b); err != nil {
		return err
	}

	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}

	err = kb.keyring.Unmarshal(rawKeyData.Keyring)
	if err != nil {
		return err
	}
	kb.chainType = rawKeyData.ChainType
	kb.id = rawKeyData.ID
	return nil
}

func (kb *keyBundle[K]) Raw() Raw {
	b, err := kb.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

// migration code
func (kbraw *keyBundleRawData) Migrate(b []byte) error {
	// if key is not stored in Keyring param, use EVM, Solana as Keyring
	// for migrating, key will only be marshalled into Keyring
	if len(kbraw.Keyring) == 0 {
		if len(kbraw.EVMKeyring) != 0 {
			kbraw.Keyring = kbraw.EVMKeyring
		} else if len(kbraw.SolanaKeyring) != 0 {
			kbraw.Keyring = kbraw.SolanaKeyring
		}
	}

	// if key does not have an ID associated with it (old formats),
	// derive the key ID and preserve it
	if bytes.Equal(kbraw.ID[:], models.EmptySha256Hash[:]) {
		kbraw.ID = sha256.Sum256(b)
	}

	return nil
}
