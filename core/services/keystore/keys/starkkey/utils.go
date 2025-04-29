package starkkey

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

// reimplements parts of
// https://github.com/NethermindEth/starknet.go/blob/0bdaab716ce24a521304744a8fbd8e01800c241d/curve/curve.go#L702
// generate the PK as a pseudo-random number in the interval [1, CurveOrder - 1]
// using io.Reader, and Key struct
func GenerateKey(material io.Reader) (k Key, err error) {
	max := new(big.Int).Sub(curve.Curve.N, big.NewInt(1))

	priv, err := rand.Int(material, max)
	if err != nil {
		return k, err
	}
	k.signFn = func(hash *big.Int) (x, y *big.Int, err error) {
		return curve.Curve.Sign(hash, priv)
	}

	k.pub.X, k.pub.Y, err = curve.Curve.PrivateToPoint(priv)
	if err != nil {
		return k, err
	}

	if !curve.Curve.IsOnCurve(k.pub.X, k.pub.Y) {
		return k, errors.New("key gen is not on stark curve")
	}
	k.raw = internal.NewRaw(padBytes(priv.Bytes()))

	return k, nil
}

// pad bytes to privateKeyLen
func padBytes(a []byte) []byte {
	if len(a) < privateKeyLen {
		pad := make([]byte, privateKeyLen-len(a))
		return append(pad, a...)
	}

	// return original if length is >= to specified length
	return a
}
