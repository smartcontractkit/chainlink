package vrfkey

import (
	"math/big"
	"testing"

	bm "github.com/smartcontractkit/chainlink/core/utils/big_math"
	"github.com/stretchr/testify/assert"
)

func TestVRF_IsSquare(t *testing.T) {
	assert.True(t, IsSquare(bm.Four))
	minusOneModP := bm.I().Sub(FieldSize, bm.One)
	assert.False(t, IsSquare(minusOneModP))
}

func TestVRF_SquareRoot(t *testing.T) {
	assert.Equal(t, bm.Two, SquareRoot(bm.Four))
}

func TestVRF_YSquared(t *testing.T) {
	assert.Equal(t, bm.Add(bm.Mul(bm.Two, bm.Mul(bm.Two, bm.Two)), bm.Seven), YSquared(bm.Two)) // 2³+7
}

func TestVRF_IsCurveXOrdinate(t *testing.T) {
	assert.True(t, IsCurveXOrdinate(big.NewInt(1)))
	assert.False(t, IsCurveXOrdinate(big.NewInt(5)))
}
