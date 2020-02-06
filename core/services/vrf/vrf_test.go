package vrf

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"chainlink/core/utils"
)

func TestVRF_IsSquare(t *testing.T) {
	assert.True(t, IsSquare(four))
	minusOneModP := i().Sub(fieldSize, one)
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

func TestVRF_fieldHash(t *testing.T) {
	utils.PanicsWithError(t, fmt.Sprintf(fieldHashPanicTemplate, 33*8),
		func() { fieldHash([]byte("much, much longer than 32 bytes!!")) })
}
