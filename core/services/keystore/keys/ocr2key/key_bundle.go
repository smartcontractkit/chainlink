package ocr2key

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	starknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// nolint
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

// check generic keybundle for each chain conforms to KeyBundle interface
var _ KeyBundle = &keyBundle[*evmKeyring]{}
var _ KeyBundle = &keyBundle[*cosmosKeyring]{}
var _ KeyBundle = &keyBundle[*solanaKeyring]{}
var _ KeyBundle = &keyBundle[*starknet.OCR2Key]{}

var curve = secp256k1.S256()

// New returns key bundle based on the chain type
func New(chainType chaintype.ChainType) (KeyBundle, error) {
	switch chainType {
	case chaintype.EVM:
		return newKeyBundleRand(chaintype.EVM, newEVMKeyring)
	case chaintype.Cosmos:
		return newKeyBundleRand(chaintype.Cosmos, newCosmosKeyring)
	case chaintype.Solana:
		return newKeyBundleRand(chaintype.Solana, newSolanaKeyring)
	case chaintype.StarkNet:
		return newKeyBundleRand(chaintype.StarkNet, starknet.NewOCR2Key)
	}
	return nil, chaintype.NewErrInvalidChainType(chainType)
}

// MustNewInsecure returns key bundle based on the chain type or panics
func MustNewInsecure(reader io.Reader, chainType chaintype.ChainType) KeyBundle {
	switch chainType {
	case chaintype.EVM:
		return mustNewKeyBundleInsecure(chaintype.EVM, newEVMKeyring, reader)
	case chaintype.Cosmos:
		return mustNewKeyBundleInsecure(chaintype.Cosmos, newCosmosKeyring, reader)
	case chaintype.Solana:
		return mustNewKeyBundleInsecure(chaintype.Solana, newSolanaKeyring, reader)
	case chaintype.StarkNet:
		return mustNewKeyBundleInsecure(chaintype.StarkNet, starknet.NewOCR2Key, reader)
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

// nolint
type Raw []byte

func (raw Raw) Key() (kb KeyBundle) {
	var temp struct{ ChainType chaintype.ChainType }
	err := json.Unmarshal(raw, &temp)
	if err != nil {
		panic(err)
	}
	switch temp.ChainType {
	case chaintype.EVM:
		kb = newKeyBundle(new(evmKeyring))
	case chaintype.Cosmos:
		kb = newKeyBundle(new(cosmosKeyring))
	case chaintype.Solana:
		kb = newKeyBundle(new(solanaKeyring))
	case chaintype.StarkNet:
		kb = newKeyBundle(new(starknet.OCR2Key))
	default:
		return nil
	}
	if err := kb.Unmarshal(raw); err != nil {
		panic(err)
	}
	return
}

// type is added to the beginning of the passwords for OCR key bundles,
// so that the keys can't accidentally be mis-used in the wrong place
func adulteratedPassword(auth string) string {
	s := "ocr2key" + auth
	return s
}
