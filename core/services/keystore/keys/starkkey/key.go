package starkkey

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	starksig "github.com/NethermindEth/juno/pkg/crypto/signature"
	"github.com/NethermindEth/juno/pkg/crypto/weierstrass"
	"github.com/ethereum/go-ethereum/common/math"
)

var curve = weierstrass.Stark()

// Raw represents the Stark private key
type Raw []byte

// Key gets the Key
func (raw Raw) Key() Key {
	privKey := starksig.PrivateKey{}
	privKey.D = new(big.Int).SetBytes(raw)
	privKey.PublicKey.Curve = curve
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(raw)

	return Key{
		privkey: privKey,
	}
}

// String returns description
func (raw Raw) String() string {
	return "<StarkNet Raw Private Key>"
}

// GoString wraps String()
func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

// Key represents StarkNet key
type Key struct {
	privkey starksig.PrivateKey
}

// New creates new Key
func New() (Key, error) {
	return newFrom(crypto_rand.Reader)
}

// MustNewInsecure return Key if no error
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

func newFrom(reader io.Reader) (Key, error) {
	privKey, err := starksig.GenerateKey(curve, reader)
	if err != nil {
		return Key{}, err
	}
	return Key{
		privkey: *privKey,
	}, nil
}

// ID gets Key ID
func (key Key) ID() string {
	return key.PublicKeyStr()
}

// PublicKeyStr
func (key Key) PublicKeyStr() string {
	return "0x" + hex.EncodeToString(PubToStarkKey(key.privkey.PublicKey))
}

// PubToStarkKey implements the pubkey to starkkey functionality: https://github.com/0xs34n/starknet.js/blob/cd61356974d355aa42f07a3d63f7ccefecbd913c/src/utils/ellipticCurve.ts#L49
func PubToStarkKey(pubkey starksig.PublicKey) []byte {
	return math.PaddedBigBytes(pubkey.X, 32)
}

// Raw from private key
func (key Key) Raw() Raw {
	return key.privkey.D.Bytes()
}

// String is the print-friendly format of the Key
func (key Key) String() string {
	return fmt.Sprintf("StarkNetKey{PrivateKey: <redacted>, Public Key: %s}", key.PublicKeyStr())
}

// GoString wraps String()
func (key Key) GoString() string {
	return key.String()
}

// ToPrivKey returns the key usable for signing.
func (key Key) ToPrivKey() starksig.PrivateKey {
	return key.privkey
}

// PublicKey copies public key object
func (key Key) PublicKey() starksig.PublicKey {
	return key.privkey.PublicKey
}
