package starkkey

import (
	crypto_rand "crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"

	adapters "github.com/smartcontractkit/chainlink-common/pkg/loop/adapters/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func KeyFor(raw internal.Raw) Key {
	k := Key{raw: raw}
	var err error

	priv := new(big.Int).SetBytes(internal.Bytes(raw))
	k.signFn = func(hash *big.Int) (x, y *big.Int, err error) {
		return curve.Curve.Sign(hash, priv)
	}
	k.pub.X, k.pub.Y, err = curve.Curve.PrivateToPoint(priv)
	if err != nil {
		panic(err) // key not generated
	}
	return k
}

type PublicKey struct {
	X, Y *big.Int
}

// Key represents Starknet key
type Key struct {
	raw    internal.Raw
	signFn func(*big.Int) (x, y *big.Int, err error)
	pub    PublicKey
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
	return GenerateKey(reader)
}

// ID gets Key ID
func (key Key) ID() string {
	return key.StarkKeyStr()
}

// StarkKeyStr is the starknet public key associated to the private key
// it is the X component of the ECDSA pubkey and used in the deployment of the account contract
// this func is used in exporting it via CLI and API
func (key Key) StarkKeyStr() string {
	return new(felt.Felt).SetBytes(key.pub.X.Bytes()).String()
}

// Raw from private key
func (key Key) Raw() internal.Raw { return key.raw }

func (key Key) Sign(hash []byte) ([]byte, error) {
	starkHash := new(big.Int).SetBytes(hash)
	x, y, err := key.signFn(starkHash)
	if err != nil {
		return nil, fmt.Errorf("error signing data with curve: %w", err)
	}
	sig, err := adapters.SignatureFromBigInts(x, y)
	if err != nil {
		return nil, err
	}
	return sig.Bytes()
}

// PublicKey copies public key object
func (key Key) PublicKey() PublicKey {
	return key.pub
}
