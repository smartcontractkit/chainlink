package ocr2key

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
)

type (
	// KeyBundle represents the bundle of keys needed for OCR
	KeyBundle struct {
		id              models.Sha256Hash
		ChainType       chaintype.ChainType
		OffchainKeyring OffchainKeyring
		evmKeyring      EVMKeyring
		solanaKeyring   SolanaKeyring
	}

	KeyBundleRawData struct {
		ChainType       chaintype.ChainType
		OffchainKeyring []byte
		EVMKeyring      []byte
		SolanaKeyring   []byte
	}
)

var (
	curve = secp256k1.S256()
)

func (kb KeyBundle) GetID() string {
	return kb.ID()
}

func (kb *KeyBundle) SetID(value string) error {
	var result models.Sha256Hash
	decodedString, err := hex.DecodeString(value)

	if err != nil {
		return err
	}

	copy(result[:], decodedString[:32])
	kb.id = result
	return nil
}

// New makes a new set of OCR key bundles from cryptographically secure entropy
func New(chainType chaintype.ChainType) (*KeyBundle, error) {
	return newKeyBundleFrom(chainType, cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func NewFromOCR1Key(v1key ocrkey.KeyV2) (KeyBundle, error) {
	evmKeyring := EVMKeyring{
		privateKey: ecdsa.PrivateKey(*v1key.OnChainSigning),
	}
	offChainKeyRing := OffchainKeyring{
		signingKey:    ed25519.PrivateKey(*v1key.OffChainSigning),
		encryptionKey: *v1key.OffChainEncryption,
	}
	k := KeyBundle{
		ChainType:       chaintype.EVM,
		evmKeyring:      evmKeyring,
		OffchainKeyring: offChainKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return KeyBundle{}, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return k, nil
}

func MustNewInsecure(reader io.Reader, chainType chaintype.ChainType) KeyBundle {
	key, err := newKeyBundleFrom(chainType, reader, reader, reader)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate new OCR2 Key"))
	}
	return *key
}

func newKeyBundleFrom(chainType chaintype.ChainType, onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*KeyBundle, error) {
	offchainKeyring, err := newOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := &KeyBundle{
		ChainType:       chainType,
		OffchainKeyring: *offchainKeyring,
	}
	switch chainType {
	case chaintype.EVM:
		evmKeyRing, err2 := newEVMKeyring(onchainSigningKeyMaterial)
		if err2 != nil {
			return nil, err2
		}
		k.evmKeyring = *evmKeyRing
	case chaintype.Solana:
		solanaKeyRing, err2 := newSolanaKeyring(onchainSigningKeyMaterial)
		if err2 != nil {
			return nil, err2
		}
		k.solanaKeyring = *solanaKeyRing
	default:
		return nil, chaintype.NewErrInvalidChainType(chainType)
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return k, nil
}

func (kb KeyBundle) ID() string {
	return hex.EncodeToString(kb.id[:])
}

func (kb KeyBundle) OnchainKeyring() ocrtypes.OnchainKeyring {
	switch kb.ChainType {
	case chaintype.EVM:
		return &kb.evmKeyring
	case chaintype.Solana:
		return &kb.solanaKeyring
	default:
		panic(errors.Wrap(chaintype.NewErrInvalidChainType(kb.ChainType), "invariant"))
	}
}

// ConfigDiffieHellman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offChainEncryption key.)
func (kb *KeyBundle) ConfigDiffieHellman(base [curve25519.PointSize]byte) ([curve25519.PointSize]byte, error) {
	return kb.OffchainKeyring.ConfigDiffieHellman(base)
}

// PublicKeyAddressOnChain returns public component of the keypair used on chain
func (kb *KeyBundle) PublicKeyAddressOnChain() string {
	switch kb.ChainType {
	case chaintype.EVM:
		return kb.evmKeyring.SigningAddress().Hex()
	case chaintype.Solana:
		return kb.solanaKeyring.SigningAddress().Hex()
	default:
		panic(errors.Wrap(chaintype.NewErrInvalidChainType(kb.ChainType), "invariant"))
	}
}

// PublicKeyAddressOnChainRaw returns public component of the keypair used on chain
func (kb *KeyBundle) PublicKeyAddressOnChainRaw() []byte {
	switch kb.ChainType {
	case chaintype.EVM:
		result := kb.evmKeyring.SigningAddress()
		return result[:]
	case chaintype.Solana:
		result := kb.solanaKeyring.SigningAddress()
		return result[:]
	default:
		panic(errors.Wrap(chaintype.NewErrInvalidChainType(kb.ChainType), "invariant"))
	}
}

// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
func (kb *KeyBundle) PublicKeyOffChain() ocrtypes.OffchainPublicKey {
	return kb.OffchainKeyring.OffchainPublicKey()
}

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
func (kb *KeyBundle) PublicKeyConfig() [curve25519.PointSize]byte {
	return kb.OffchainKeyring.ConfigEncryptionPublicKey()
}

func (kb *KeyBundle) Marshal() ([]byte, error) {
	offchainKeyringBytes, err := kb.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := KeyBundleRawData{
		ChainType:       kb.ChainType,
		OffchainKeyring: offchainKeyringBytes,
	}
	switch kb.ChainType {
	case chaintype.EVM:
		evmKeyringBytes, err := kb.evmKeyring.marshal()
		if err != nil {
			return nil, err
		}
		rawKeyData.EVMKeyring = evmKeyringBytes
	case chaintype.Solana:
		solanaKeyringBytes, err := kb.solanaKeyring.marshal()
		if err != nil {
			return nil, err
		}
		rawKeyData.SolanaKeyring = solanaKeyringBytes
	default:
		panic(errors.Wrap(chaintype.NewErrInvalidChainType(kb.ChainType), "invariant"))
	}
	return json.Marshal(&rawKeyData)
}

func (kb *KeyBundle) Unmarshal(b []byte) (err error) {
	var rawKeyData KeyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = kb.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	switch rawKeyData.ChainType {
	case chaintype.EVM:
		err = kb.evmKeyring.unmarshal(rawKeyData.EVMKeyring)
		if err != nil {
			return err
		}
	case chaintype.Solana:
		err = kb.solanaKeyring.unmarshal(rawKeyData.SolanaKeyring)
		if err != nil {
			return err
		}
	default:
		panic(errors.Wrap(chaintype.NewErrInvalidChainType(kb.ChainType), "invariant"))
	}
	kb.ChainType = rawKeyData.ChainType
	kb.id = sha256.Sum256(b)
	return nil
}

// String reduces the risk of accidentally logging the private key
func (kb KeyBundle) String() string {
	addressOnChain := kb.PublicKeyAddressOnChain()
	return fmt.Sprintf(
		"KeyBundle{PublicKeyAddressOnChain: %s, PublicKeyOffChain: %s}",
		addressOnChain,
		hex.EncodeToString(kb.PublicKeyOffChain()),
	)
}

// GoStringer reduces the risk of accidentally logging the private key
func (kb KeyBundle) GoStringer() string {
	return kb.String()
}

// type is added to the beginning of the passwords for OCR key bundles,
// so that the keys can't accidentally be mis-used in the wrong place
func adulteratedPassword(auth string) string {
	s := "ocr2key" + auth
	return s
}
