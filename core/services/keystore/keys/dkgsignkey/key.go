package dkgsignkey

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

var suite = edwards25519.NewBlakeSHA256Ed25519()

// Raw represents a raw dkgsign secret key in little-endian byte order.
type Raw []byte

// Key returns a Key object from this raw data.
func (r Raw) Key() Key {
	privKey := suite.Scalar().SetBytes(r)
	key, err := keyFromScalar(privKey)
	if err != nil {
		panic(err) // should never happen
	}
	return key
}

func (r Raw) String() string {
	return "<DKGSign Raw Private Key>"
}

func (r Raw) GoString() string {
	return r.String()
}

// Key is DKG signing key that conforms to the keystore.Key interface
type Key struct {
	privateKey     kyber.Scalar
	publicKeyBytes []byte
	PublicKey      kyber.Point
}

// New creates a new DKGSign key
func New() (Key, error) {
	privateKey := suite.Scalar().Pick(suite.RandomStream())
	return keyFromScalar(privateKey)
}

// MustNewXXXTestingOnly creates a new DKGSign key from the given secret key.
// NOTE: for testing only.
func MustNewXXXTestingOnly(sk *big.Int) Key {
	key, err := keyFromScalar(scalarFromBig(sk))
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
	return fmt.Sprintf("DKGSignKey{PrivateKey: <redacted>, PublicKey: %s", k.PublicKey)
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
	return suite.Scalar().Set(k.privateKey)
}
