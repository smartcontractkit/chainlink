package vrf

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
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
