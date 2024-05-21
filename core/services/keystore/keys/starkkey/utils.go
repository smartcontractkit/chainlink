package starkkey

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

// constants
var (
	byteLen = 32
)

// reimplements parts of
// https://github.com/NethermindEth/starknet.go/blob/0bdaab716ce24a521304744a8fbd8e01800c241d/curve/curve.go#L702
// generate the PK as a pseudo-random number in the interval [1, CurveOrder - 1]
// using io.Reader, and Key struct
func GenerateKey(material io.Reader) (k Key, err error) {
	max := new(big.Int).Sub(curve.Curve.N, big.NewInt(1))

	k.priv, err = rand.Int(material, max)
	if err != nil {
		return k, err
	}

	k.pub.X, k.pub.Y, err = curve.Curve.PrivateToPoint(k.priv)
	if err != nil {
		return k, err
	}

	if !curve.Curve.IsOnCurve(k.pub.X, k.pub.Y) {
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
