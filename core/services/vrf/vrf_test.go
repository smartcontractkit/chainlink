package vrf

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVRF_IsSquare(t *testing.T) {
	assert.True(t, IsSquare(four))
	minusOneModP := i().Sub(FieldSize, one)
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

func TestVRF_VerifyProof(t *testing.T) {
	sk, seed, nonce := big.NewInt(1), big.NewInt(2), big.NewInt(3)
	p, err := generateProofWithNonce(sk, seed, nonce)
	require.NoError(t, err, "could not generate proof")
	p.Seed = big.NewInt(0).Add(seed, big.NewInt(1))
	valid, err := p.VerifyVRFProof()
	require.NoError(t, err, "could not validate proof")
	require.False(t, valid, "invalid proof was found valid")
}
