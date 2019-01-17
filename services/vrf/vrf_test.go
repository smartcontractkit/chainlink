package vrf

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, big.NewInt(2*2*2+3), YSquared(big.NewInt(2)))
}

func TestVRF_IsCurveXOrdinate(t *testing.T) {
	assert.True(t, IsCurveXOrdinate(big.NewInt(1)))
	assert.False(t, IsCurveXOrdinate(big.NewInt(4)))
}

func TestVRF_CoordsFromPoint(t *testing.T) {
	x, y := CoordsFromPoint(Generator)
	assert.Equal(t, x, big.NewInt(1))
	assert.Equal(t, y, big.NewInt(2))
}

func bigFromHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, s))
	}
	return n
}

func TestVRF_ZqHash(t *testing.T) {
	log2Mod := 5.0
	modulus := big.NewInt(int64(math.Pow(2, log2Mod)))
	bitMask := big.NewInt(int64(math.Pow(2, log2Mod+1) - 1))
	reHashTriggeringSeed := big.NewInt(0)
	hash, err := HashUint256s(reHashTriggeringSeed)
	if err != nil {
		panic(err)
	}
	hash.And(hash, bitMask)
	assert.Equal(t, 1, hash.Cmp(modulus),
		`need an example which hashes to something bigger than the
modulus, to test the rehash logic.`)
	zqHash, err := ZqHash(modulus, reHashTriggeringSeed)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, zqHash, big.NewInt(13))
}

func TestVRF_HashToCurve(t *testing.T) {
	reHashTriggeringInput := []*big.Int{
		big.NewInt(1), big.NewInt(2), big.NewInt(5)}
	x, err := ZqHash(P, reHashTriggeringInput...)
	if err != nil {
		panic(err)
	}
	assert.False(t, IsCurveXOrdinate(x),
		`need an example where first hash is not an x-ordinate for any
point on the curve, to exercise rehash logic.`)
	p, err := HashToCurve(reHashTriggeringInput[0],
		reHashTriggeringInput[1], reHashTriggeringInput[2])
	if err != nil {
		panic(err)
	}
	x, y := CoordsFromPoint(p)
	// See 'Hashes to the curve with the same results as the golang code' in Curve.js
	eX := "247154f2ce523897365341b03669e1061049e801e8750ae708e1cb02f36cb225"
	eY := "16e1157d5b94324127e094abe222a05a5c47be3124254a6aa047d5e1f2d864ea"
	assert.Equal(t, bigFromHex(eX), x)
	assert.Equal(t, bigFromHex(eY), y)
}

func TestVRF_ScalarFromCurve(t *testing.T) {
	g := Generator
	s, err := ScalarFromCurve(g, g, g, g, g)
	if err != nil {
		panic(err)
	}
	eS := "57bf013147ceec913f17ef97d3bcfad8315d99752af81f8913ad1c88493e669"
	// See 'Computes the same hashed scalar from curve points as the golang code' in Curve.js
	assert.Equal(t, bigFromHex(eS), s)
}

func pointsEqual(p1, p2 *bn256.G1) bool {
	s1, _ := ScalarFromCurve(p1)
	s2, _ := ScalarFromCurve(p2)
	return s1.Cmp(s2) == 0
}

func TestVRF_GenerateProof(t *testing.T) {
	secretKeyHaHaNeverDoThis := big.NewInt(2)
	seed := big.NewInt(0)
	// Can't test c & s: They vary from run to run.
	proof, err := GenerateProof(secretKeyHaHaNeverDoThis, seed)
	if err != nil {
		panic(err)
	}
	publicKey := new(bn256.G1).ScalarMult(
		Generator, secretKeyHaHaNeverDoThis)
	assert.True(t, pointsEqual(publicKey, proof.PublicKey))
	gammaX, gammaY := CoordsFromPoint(proof.Gamma)
	// See 'Accepts a valid VRF proof' in VRF.js
	gX := "26feb384a4a3f28742d0e0e0f5458474ba54ef9816d4d31f3bf538dfcf67cf3f"
	gY := "1eaed2431dd78ad75dd0c9f013cabff4f1d8c4c83cda79fff3855c988a3606d8"
	assert.Equal(t, bigFromHex(gX), gammaX)
	assert.Equal(t, bigFromHex(gY), gammaY)
	verification, err := proof.VerifyProof()
	if err != nil {
		panic(err)
	}
	if !verification {
		panic(err)
	}
}
