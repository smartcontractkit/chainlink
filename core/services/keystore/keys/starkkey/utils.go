package starkkey

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/smartcontractkit/caigo"
)

// constants
var (
	byteLen = 32
)

// reimplements parts of https://github.com/smartcontractkit/caigo/blob/main/utils.go#L85
// generate the PK as a pseudo-random number in the interval [1, CurveOrder - 1]
// using io.Reader, and Key struct
func GenerateKey(material io.Reader) (k Key, err error) {
	max := new(big.Int).Sub(caigo.Curve.N, big.NewInt(1))

	k.priv, err = rand.Int(material, max)
	if err != nil {
		return k, err
	}

	k.pub.X, k.pub.Y, err = caigo.Curve.PrivateToPoint(k.priv)
	if err != nil {
		return k, err
	}

	if !caigo.Curve.IsOnCurve(k.pub.X, k.pub.Y) {
		return k, fmt.Errorf("key gen is not on stark curve")
	}

	return k, nil
}

// pad bytes to specific length
func padBytes(a []byte, length int) []byte {
	if len(a) < length {
		pad := make([]byte, length-len(a))
		return append(pad, a...)
	}

	// return original if length is >= to specified length
	return a
}
