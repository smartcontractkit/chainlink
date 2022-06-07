package ocr2key

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

//nolint
type KeyBundle interface {
	// OnchainKeyring is used for signing reports (groups of observations, verified onchain)
	ocrtypes.OnchainKeyring
	// OffchainKeyring is used for signing observations
	ocrtypes.OffchainKeyring
	ID() string
	ChainType() chaintype.ChainType
	Marshal() ([]byte, error)
	Unmarshal(b []byte) (err error)
	Raw() Raw
	OnChainPublicKey() string
}

var curve = secp256k1.S256()

// New returns key bundle based on the chain type
func New(chainType chaintype.ChainType) (KeyBundle, error) {
	switch chainType {
	case chaintype.EVM:
		return newKeyBundle(chaintype.EVM, newEVMKeyring)
	case chaintype.Solana:
		return newKeyBundle(chaintype.Solana, newSolanaKeyring)
	case chaintype.Terra:
		return newKeyBundle(chaintype.Terra, newTerraKeyring)
	case chaintype.Starknet:
		return newKeyBundle(chaintype.Starknet, newStarknetKeyring)
	}
	return nil, chaintype.NewErrInvalidChainType(chainType)
}

// MustNewInsecure returns key bundle based on the chain type or panics
func MustNewInsecure(reader io.Reader, chainType chaintype.ChainType) KeyBundle {
	switch chainType {
	case chaintype.EVM:
		return mustNewKeyBundleInsecure(chaintype.EVM, newEVMKeyring, reader)
	case chaintype.Solana:
		return mustNewKeyBundleInsecure(chaintype.Solana, newSolanaKeyring, reader)
	case chaintype.Terra:
		return mustNewKeyBundleInsecure(chaintype.Terra, newTerraKeyring, reader)
	case chaintype.Starknet:
		return mustNewKeyBundleInsecure(chaintype.Starknet, newStarknetKeyring, reader)
	}
	panic(chaintype.NewErrInvalidChainType(chainType))
}

// NewKeyBundleFromOCR1Key gets the key bundle from an OCR1 key
func NewKeyBundleFromOCR1Key(v1key ocrkey.KeyV2) (keyBundle[*evmKeyring], error) {
	onChainKeyRing := evmKeyring{
		privateKey: ecdsa.PrivateKey(*v1key.OnChainSigning),
	}
	offChainKeyRing := OffchainKeyring{
		signingKey:    ed25519.PrivateKey(*v1key.OffChainSigning),
		encryptionKey: *v1key.OffChainEncryption,
	}
	k := keyBundle[*evmKeyring]{
		keyBundleBase: keyBundleBase{
			chainType:       chaintype.EVM,
			OffchainKeyring: offChainKeyRing,
		},
		keyring: &onChainKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return keyBundle[*evmKeyring]{}, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return k, nil
}

var _ fmt.GoStringer = &keyBundleBase{}

type keyBundleBase struct {
	OffchainKeyring
	id        models.Sha256Hash
	chainType chaintype.ChainType
}

func (kb keyBundleBase) ID() string {
	return hex.EncodeToString(kb.id[:])
}

// ChainType gets the chain type from the key bundle
func (kb keyBundleBase) ChainType() chaintype.ChainType {
	return kb.chainType
}

// String reduces the risk of accidentally logging the private key
func (kb keyBundleBase) String() string {
	return fmt.Sprintf("KeyBundle{chainType: %s, id: %s}", kb.ChainType(), kb.ID())
}

// GoString reduces the risk of accidentally logging the private key
func (kb keyBundleBase) GoString() string {
	return kb.String()
}

//nolint
type Raw []byte

func (raw Raw) Key() KeyBundle {
	var temp struct{ ChainType chaintype.ChainType }
	err := json.Unmarshal(raw, &temp)
	if err != nil {
		panic(err)
	}
	switch temp.ChainType {
	case chaintype.EVM:
		result := mustNewKeyFromRaw(raw, &evmKeyring{})
		return &result
	case chaintype.Solana:
		result := mustNewKeyFromRaw(raw, &solanaKeyring{})
		return &result
	case chaintype.Terra:
		result := mustNewKeyFromRaw(raw, &terraKeyring{})
		return &result
	case chaintype.Starknet:
		result := mustNewKeyFromRaw(raw, &starknetKeyring{})
		return &result
	default:
		panic(chaintype.NewErrInvalidChainType(temp.ChainType))
	}
}

// type is added to the beginning of the passwords for OCR key bundles,
// so that the keys can't accidentally be mis-used in the wrong place
func adulteratedPassword(auth string) string {
	s := "ocr2key" + auth
	return s
}
