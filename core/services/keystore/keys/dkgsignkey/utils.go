package dkgsignkey

import (
	"math/big"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
)

// scalarFromBig creates a kyber.Scalar belonging to the edwards25519
// kyber suite from a big integer. This is useful for testing.
func scalarFromBig(i *big.Int) kyber.Scalar {
	scalar := suite.Scalar()
	// big.Int.Bytes() returns a byte slice in big-endian order,
	// need to reverse the slice before we SetBytes since
	// SetBytes interprets it in little-endian order.
	b := i.Bytes()
	reverseSliceInPlace(b)
	return scalar.SetBytes(b)
}

func keyFromScalar(k kyber.Scalar) (Key, error) {
	publicKey := suite.Point().Base().Mul(k, nil)
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

func reverseSliceInPlace[T any](elems []T) {
	for i := 0; i < len(elems)/2; i++ {
		elems[i], elems[len(elems)-i-1] = elems[len(elems)-i-1], elems[i]
	}
}
