package vrfkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	bm "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

func TestUnit_VRF_IsSquare(t *testing.T) {
	assert.True(t, IsSquare(bm.Four))
	minusOneModP := bm.I().Sub(FieldSize, bm.One)
	assert.False(t, IsSquare(minusOneModP))
}

func TestUnit_VRF_SquareRoot(t *testing.T) {
	assert.Equal(t, bm.Two, SquareRoot(bm.Four))
}

func TestUnit_VRF_YSquared(t *testing.T) {
	assert.Equal(t, bm.Add(bm.Mul(bm.Two, bm.Mul(bm.Two, bm.Two)), bm.Seven), YSquared(bm.Two)) // 2³+7
}

func TestUnit_VRF_IsCurveXOrdinate(t *testing.T) {
	assert.True(t, IsCurveXOrdinate(big.NewInt(1)))
	assert.False(t, IsCurveXOrdinate(big.NewInt(5)))
}
