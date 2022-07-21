package dkgencryptkey

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
)

var suite pairing.Suite = &altbn_128.PairingSuite{}
var g1 = suite.G1()

type Raw []byte

func (r Raw) Key() Key {
	scalar := g1.Scalar()
	err := scalar.UnmarshalBinary(r)
	if err != nil {
		panic(err) // should never happen6
	}
	key, err := keyFromScalar(scalar)
	if err != nil {
		panic(err) // should never happen
	}
	return key
}

func (r Raw) String() string {
	return "<DKGEncrypt Raw Private Key>"
}

func (r Raw) GoString() string {
	return r.String()
}

type Key struct {
	privateKey     kyber.Scalar
	publicKeyBytes []byte
	PublicKey      kyber.Point
}

// New returns a new dkgencryptkey key
func New() (Key, error) {
	return keyFromScalar(g1.Scalar().Pick(suite.RandomStream()))
}

// MustNewXXXTestingOnly creates a new DKGEncrypt key from the given secret key.
// NOTE: for testing only.
func MustNewXXXTestingOnly(sk *big.Int) Key {
	key, err := keyFromScalar(g1.Scalar().SetInt64(sk.Int64()))
	if err != nil {
		panic(err)
	}
	return key
}

var _ fmt.GoStringer = &Key{}

// GoString implements fmt.GoStringer
func (k Key) GoString() string {
	return k.String()
}

// String returns the string representation of this key
func (k Key) String() string {
	return fmt.Sprintf("DKGEncryptKey{PrivateKey: <redacted>, PublicKey: %s", k.PublicKeyString())
}

// ID returns the ID of this key
func (k Key) ID() string {
	return k.PublicKeyString()
}

// PublicKeyString returns the hex representation of this key's public key
func (k Key) PublicKeyString() string {
	return hex.EncodeToString(k.publicKeyBytes)
}

// Raw returns the key raw data
func (k Key) Raw() Raw {
	raw, err := k.privateKey.MarshalBinary()
	if err != nil {
		panic(err) // should never happen
	}
	return Raw(raw)
}

// KyberScalar returns the private key as a kyber.Scalar object
func (k Key) KyberScalar() kyber.Scalar {
	return g1.Scalar().Set(k.privateKey)
}

// KyberPoint returns the public key as a kyber.Point object
func (k Key) KyberPoint() kyber.Point {
	return g1.Point().Base().Mul(k.privateKey, nil)
}

// keyFromScalar creates a new dkgencryptkey key from the given scalar.
// the given scalar must be a scalar of the g1 group in the altbn_128 pairing.
func keyFromScalar(k kyber.Scalar) (Key, error) {
	publicKey := g1.Point().Base().Mul(k, nil)
	publicKeyBytes, err := publicKey.MarshalBinary()
	if err != nil {
		return Key{}, errors.Wrap(err, "kyber point MarshalBinary")
	}
	return Key{
		privateKey:     k,
		PublicKey:      publicKey,
		publicKeyBytes: publicKeyBytes,
	}, nil
}
