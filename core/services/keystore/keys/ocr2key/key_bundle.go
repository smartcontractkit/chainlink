package ocr2key

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
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
		return newEVMKeyBundle()
	case chaintype.Solana:
		return newSolanaKeyBundle()
	case chaintype.Terra:
		return newTerraKeyBundle()
	}
	return nil, chaintype.NewErrInvalidChainType(chainType)
}

// MustNewInsecure returns key bundle based on the chain type or panics
func MustNewInsecure(reader io.Reader, chainType chaintype.ChainType) KeyBundle {
	switch chainType {
	case chaintype.EVM:
		return mustNewEVMKeyBundleInsecure(reader)
	case chaintype.Solana:
		return mustNewSolanaKeyBundleInsecure(reader)
	case chaintype.Terra:
		return mustNewTerraKeyBundleInsecure(reader)
	}
	panic(chaintype.NewErrInvalidChainType(chainType))
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
		result := mustNewEVMKeyFromRaw(raw)
		return &result
	case chaintype.Solana:
		result := mustNewSolanaKeyFromRaw(raw)
		return &result
	case chaintype.Terra:
		result := mustNewTerraKeyFromRaw(raw)
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
