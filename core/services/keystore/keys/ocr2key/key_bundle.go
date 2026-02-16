package ocr2key

import (
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
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
	ChainType() corekeys.ChainType
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
var _ KeyBundle = &keyBundle[*ed25519Keyring]{}
var _ KeyBundle = &keyBundle[*tonKeyring]{}
var _ KeyBundle = &keyBundle[*ed25519Keyring]{}

var curve = secp256k1.S256()

// New returns key bundle based on the chain type
func New(chainType corekeys.ChainType) (KeyBundle, error) {
	switch chainType {
	case corekeys.EVM:
		return newKeyBundleRand(corekeys.EVM, newEVMKeyring)
	case corekeys.Cosmos:
		return newKeyBundleRand(corekeys.Cosmos, newCosmosKeyring)
	case corekeys.Solana:
		return newKeyBundleRand(corekeys.Solana, newSolanaKeyring)
	case corekeys.StarkNet:
		return newKeyBundleRand(corekeys.StarkNet, starkkey.NewOCR2Key)
	case corekeys.Aptos:
		return newKeyBundleRand(corekeys.Aptos, newEd25519Keyring)
	case corekeys.Tron:
		return newKeyBundleRand(corekeys.Tron, newEVMKeyring)
	case corekeys.TON:
		return newKeyBundleRand(corekeys.TON, newTONKeyring)
	case corekeys.Sui:
		return newKeyBundleRand(corekeys.Sui, newEd25519Keyring)
	}
	return nil, corekeys.NewErrInvalidChainType(chainType)
}

// MustNewInsecure returns key bundle based on the chain type or panics
func MustNewInsecure(reader io.Reader, chainType corekeys.ChainType) KeyBundle {
	switch chainType {
	case corekeys.EVM:
		return mustNewKeyBundleInsecure(corekeys.EVM, newEVMKeyring, reader)
	case corekeys.Cosmos:
		return mustNewKeyBundleInsecure(corekeys.Cosmos, newCosmosKeyring, reader)
	case corekeys.Solana:
		return mustNewKeyBundleInsecure(corekeys.Solana, newSolanaKeyring, reader)
	case corekeys.StarkNet:
		return mustNewKeyBundleInsecure(corekeys.StarkNet, starkkey.NewOCR2Key, reader)
	case corekeys.Aptos:
		return mustNewKeyBundleInsecure(corekeys.Aptos, newEd25519Keyring, reader)
	case corekeys.Tron:
		return mustNewKeyBundleInsecure(corekeys.Tron, newEVMKeyring, reader)
	case corekeys.TON:
		return mustNewKeyBundleInsecure(corekeys.TON, newTONKeyring, reader)
	case corekeys.Sui:
		return mustNewKeyBundleInsecure(corekeys.Sui, newEd25519Keyring, reader)
	}
	panic(corekeys.NewErrInvalidChainType(chainType))
}

type keyBundleBase struct {
	offchainKeyring
	id        keys.Sha256Hash
	chainType corekeys.ChainType
}

func (kb keyBundleBase) ID() string {
	return hex.EncodeToString(kb.id[:])
}

// ChainType gets the chain type from the key bundle
func (kb keyBundleBase) ChainType() corekeys.ChainType {
	return kb.chainType
}

func KeyFor(raw internal.Raw) (kb KeyBundle) {
	var temp struct{ ChainType corekeys.ChainType }
	err := json.Unmarshal(internal.Bytes(raw), &temp)
	if err != nil {
		panic(err)
	}
	switch temp.ChainType {
	case corekeys.EVM:
		kb = newKeyBundle(new(evmKeyring))
	case corekeys.Cosmos:
		kb = newKeyBundle(new(cosmosKeyring))
	case corekeys.Solana:
		kb = newKeyBundle(new(solanaKeyring))
	case corekeys.StarkNet:
		kb = newKeyBundle(new(starkkey.OCR2Key))
	case corekeys.Aptos:
		kb = newKeyBundle(new(ed25519Keyring))
	case corekeys.Tron:
		kb = newKeyBundle(new(evmKeyring))
	case corekeys.TON:
		kb = newKeyBundle(new(tonKeyring))
	case corekeys.Sui:
		kb = newKeyBundle(new(ed25519Keyring))
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
