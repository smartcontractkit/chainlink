package ocr2key

import (
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type OCR3SignerVerifier interface {
	SignBlob(b []byte) (sig []byte, err error)
	VerifyBlob(publicKey ocrtypes.OnchainPublicKey, b []byte, sig []byte) bool
	Sign3(digest ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error)
	Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool
}

type KeyBundle interface {
	// OnchainKeyring is used for signing reports (groups of observations, verified onchain)
	ocrtypes.OnchainKeyring
	// offchainKeyring is used for signing observations
	ocrtypes.OffchainKeyring

	OCR3SignerVerifier

	ID() string
	ChainType() chaintype.ChainType
	Marshal() ([]byte, error)
	Unmarshal(b []byte) (err error)
	Raw() internal.Raw
	OnChainPublicKey() string
	// Decrypts ciphertext using the encryptionKey from an OCR2 offchainKeyring
	NaclBoxOpenAnonymous(ciphertext []byte) (plaintext []byte, err error)
}

// check generic keybundle for each chain conforms to KeyBundle interface
var _ KeyBundle = &keyBundle[*evmKeyring]{}
var _ KeyBundle = &keyBundle[*cosmosKeyring]{}
var _ KeyBundle = &keyBundle[*solanaKeyring]{}
var _ KeyBundle = &keyBundle[*starkkey.OCR2Key]{}
var _ KeyBundle = &keyBundle[*aptosKeyring]{}

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
		return newKeyBundleRand(chaintype.StarkNet, starkkey.NewOCR2Key)
	case chaintype.Aptos:
		return newKeyBundleRand(chaintype.Aptos, newAptosKeyring)
	case chaintype.Tron:
		return newKeyBundleRand(chaintype.Tron, newEVMKeyring)
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
		return mustNewKeyBundleInsecure(chaintype.StarkNet, starkkey.NewOCR2Key, reader)
	case chaintype.Aptos:
		return mustNewKeyBundleInsecure(chaintype.Aptos, newAptosKeyring, reader)
	case chaintype.Tron:
		return mustNewKeyBundleInsecure(chaintype.Tron, newEVMKeyring, reader)
	}
	panic(chaintype.NewErrInvalidChainType(chainType))
}

type keyBundleBase struct {
	offchainKeyring
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

func KeyFor(raw internal.Raw) (kb KeyBundle) {
	var temp struct{ ChainType chaintype.ChainType }
	err := json.Unmarshal(internal.Bytes(raw), &temp)
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
		kb = newKeyBundle(new(starkkey.OCR2Key))
	case chaintype.Aptos:
		kb = newKeyBundle(new(aptosKeyring))
	case chaintype.Tron:
		kb = newKeyBundle(new(evmKeyring))
	default:
		return nil
	}
	if err := kb.Unmarshal(internal.Bytes(raw)); err != nil {
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
