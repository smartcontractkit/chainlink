package vrf

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"

	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/utils"
)

func TestVRF_IsSquare(t *testing.T) {
	assert.True(t, IsSquare(four))
	minusOneModP := new(big.Int).Sub(fieldSize, one)
	assert.False(t, IsSquare(minusOneModP))
}

func TestVRF_SquareRoot(t *testing.T) {
	assert.Equal(t, two, SquareRoot(four))
}

func TestVRF_YSquared(t *testing.T) {
	assert.Equal(t, add(mul(two, mul(two, two)), seven), YSquared(two)) // 2³+7
}

func TestVRF_IsCurveXOrdinate(t *testing.T) {
	assert.True(t, IsCurveXOrdinate(big.NewInt(1)))
	assert.False(t, IsCurveXOrdinate(big.NewInt(5)))
}

func TestVRF_CoordsFromPoint(t *testing.T) {
	x, y := secp256k1.Coordinates(Generator)
	assert.Equal(t, x, bigFromHex(
		"79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798"))
	assert.Equal(t, y, bigFromHex(
		"483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8"))
}

func TestVRF_ZqHash(t *testing.T) {
	var log2Mod uint = 256
	modulus := lsh(one, log2Mod-1)
	hash := sub(lsh(one, log2Mod), one)
	assert.Equal(t, 1, hash.Cmp(modulus),
		`need an example which hashes to something bigger than the modulus, to test the rehash logic.`)
	zqHash, err := ZqHash(modulus, hash.Bytes())
	require.NoError(t, err)
	assert.Equal(
		t,
		bigFromHex("1ae61e33ec9365756efc1436222a72df7fdb74651e25c38bde613482291a0c69"),
		zqHash,
	)
	utils.PanicsWithError(t, fmt.Sprintf(zqHashPanicTemplate, 1, 24),
		func() { ZqHash(one, []byte("foo")) })
	utils.PanicsWithError(t, fmt.Sprintf(zqHashPanicTemplate, 257, 24),
		func() { ZqHash(lsh(fieldSize, 1), []byte("foo")) })
	utils.PanicsWithError(t, fmt.Sprintf(zqHashPanicTemplate, 256, 33*8),
		func() { ZqHash(fieldSize, []byte("much, much longer than 32 bytes!!")) })
}

func address(t *testing.T, p kyber.Point) [20]byte {
	a, err := secp256k1.EthereumAddress(p)
	require.NoError(t, err)
	return a
}
