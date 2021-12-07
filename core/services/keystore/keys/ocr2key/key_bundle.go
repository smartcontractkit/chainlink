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

type KeyBundle interface {
	ocrtypes.OnchainKeyring
	ocrtypes.OffchainKeyring
	ID() string
	ChainType() chaintype.ChainType
	Marshal() ([]byte, error)
	Unmarshal(b []byte) (err error)
	Raw() Raw
	PublicKeyAddressOnChain() string
}

var curve = secp256k1.S256()

func New(chainType chaintype.ChainType) (KeyBundle, error) {
	switch chainType {
	case chaintype.EVM:
		return newEVMKeyBundle()
	case chaintype.Solana:
		return newSolanaKeyBundle()
	}
	return nil, chaintype.NewErrInvalidChainType(chainType)
}

func MustNewInsecure(reader io.Reader, chainType chaintype.ChainType) KeyBundle {
	switch chainType {
	case chaintype.EVM:
		return mustNewEVMKeyBundleInsecure(reader)
	case chaintype.Solana:
		return mustNewSolanaKeyBundleInsecure(reader)
	}
	panic(chaintype.NewErrInvalidChainType(chainType))
}

func NewKeyBundleFromOCR1Key(v1key ocrkey.KeyV2) (evmKeyBundle, error) {
	onChainKeyRing := evmKeyring{
		privateKey: ecdsa.PrivateKey(*v1key.OnChainSigning),
	}
	offChainKeyRing := OffchainKeyring{
		signingKey:    ed25519.PrivateKey(*v1key.OffChainSigning),
		encryptionKey: *v1key.OffChainEncryption,
	}
	k := evmKeyBundle{
		keyBundleBase: keyBundleBase{
			chainType:       chaintype.EVM,
			OffchainKeyring: offChainKeyRing,
		},
		evmKeyring: onChainKeyRing,
	}
	marshalledPrivK, err := k.Marshal()
	if err != nil {
		return evmKeyBundle{}, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return k, nil
}

type keyBundleBase struct {
	OffchainKeyring
	id        models.Sha256Hash
	chainType chaintype.ChainType
}

func (kb keyBundleBase) ID() string {
	return hex.EncodeToString(kb.id[:])
}

func (kb keyBundleBase) ChainType() chaintype.ChainType {
	return kb.chainType
}

// String reduces the risk of accidentally logging the private key
func (kb keyBundleBase) String() string {
	return fmt.Sprintf("KeyBundle{chainType: %s, id: %s}", kb.ChainType(), kb.ID())
}

// GoStringer reduces the risk of accidentally logging the private key
func (kb keyBundleBase) GoStringer() string {
	return kb.String()
}

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
