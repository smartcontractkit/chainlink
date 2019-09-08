package vrf

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestVRF_IsSquare(t *testing.T) {
	assert.True(t, IsSquare(big.NewInt(4)))
	minusOneModP := new(big.Int).Sub(P, big.NewInt(1))
	assert.False(t, IsSquare(minusOneModP))
}

func TestVRF_SquareRoot(t *testing.T) {
	assert.Equal(t, big.NewInt(2), SquareRoot(big.NewInt(4)))
}

func TestVRF_YSquared(t *testing.T) {
	assert.Equal(t, big.NewInt(15), YSquared(two))
}

func TestVRF_IsCurveXOrdinate(t *testing.T) {
	assert.True(t, IsCurveXOrdinate(big.NewInt(1)))
	assert.False(t, IsCurveXOrdinate(big.NewInt(5)))
}

func TestVRF_CoordsFromPoint(t *testing.T) {
	x, y := CoordsFromPoint(Generator)
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
	if err != nil {
		panic(err)
	}
	assert.Equal(
		t,
		bigFromHex("1ae61e33ec9365756efc1436222a72df7fdb74651e25c38bde613482291a0c69"),
		zqHash,
	)
}

func TestVRF_HashToCurve(t *testing.T) {
	reHashTriggeringInput := one
	iHash, err := utils.Keccak256(append(secp256k1.LongMarshal(Generator),
		asUint256(reHashTriggeringInput)...))
	require.NoError(t, err)
	x, err := ZqHash(P, iHash)
	require.NoError(t, err)
	assert.False(t, IsCurveXOrdinate(x),
		`need an example where first hash is not an x-ordinate for any point on the curve, to exercise rehash logic.`)
	p, err := HashToCurve(Generator, reHashTriggeringInput)
	if err != nil {
		panic(err)
	}
	x, y := CoordsFromPoint(p)
	// See 'Hashes to the curve with the same results as the golang code' in Curve.js
	eX := "530fddd863609aa12030a07c5fdb323bb392a88343cea123b7f074883d2654c4"
	eY := "6fd4ee394bf2a3de542c0e5f3c86fc8f75b278a017701a59d69bdf5134dd6b70"
	assert.Equal(t, bigFromHex(eX), x)
	assert.Equal(t, bigFromHex(eY), y)
}

func TestVRF_ScalarFromCurvePoints(t *testing.T) {
	g := Generator
	ga, err := secp256k1.EthereumAddress(g)
	require.NoError(t, err)
	s := ScalarFromCurvePoints(g, g, g, ga, g)
	eS := "2b1049accb1596a24517f96761b22600a690ee5c6b6cadae3fa522e7d95ba338"
	// See 'Computes the same hashed scalar from curve points as the golang code' in Curve.js
	assert.Equal(t, bigFromHex(eS), s)
}

func TestVRF_GenerateProof(t *testing.T) {
	secretKeyHaHaNeverDoThis := big.NewInt(1)
	seed := one
	nonce := one
	// Can't test c & s: They vary from run to run.
	proof, err := GenerateProof(secretKeyHaHaNeverDoThis, seed, nonce)
	require.NoError(t, err)
	publicKey := rcurve.Point().Mul(
		secp256k1.IntToScalar(secretKeyHaHaNeverDoThis), Generator)
	assert.True(t, publicKey.Equal(proof.PublicKey))
	gammaX, gammaY := CoordsFromPoint(proof.Gamma)
	// See 'Accepts a valid VRF proof' in VRF.js. These outputs are used there
	fmt.Printf("cGamma %+v\n", rcurve.Point().Mul(secp256k1.IntToScalar(proof.C), proof.Gamma))
	h, err := HashToCurve(publicKey, seed)
	require.NoError(t, err)
	fmt.Printf("sHash %+v\n", rcurve.Point().Mul(secp256k1.IntToScalar(proof.S), h))
	gX := "530fddd863609aa12030a07c5fdb323bb392a88343cea123b7f074883d2654c4"
	gY := "6fd4ee394bf2a3de542c0e5f3c86fc8f75b278a017701a59d69bdf5134dd6b70"
	assert.Equal(t, bigFromHex(gX), gammaX)
	assert.Equal(t, bigFromHex(gY), gammaY)
	verification, err := proof.Verify()
	require.NoError(t, err)
	assert.True(t, verification, "proof verification failed")
}
