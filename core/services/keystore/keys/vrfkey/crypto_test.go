package vrfkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	bm "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

func TestVRF_IsSquare(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	assert.True(t, IsSquare(bm.Four))
	minusOneModP := bm.I().Sub(FieldSize, bm.One)
	assert.False(t, IsSquare(minusOneModP))
}

func TestVRF_SquareRoot(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	assert.Equal(t, bm.Two, SquareRoot(bm.Four))
}

func TestVRF_YSquared(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	assert.Equal(t, bm.Add(bm.Mul(bm.Two, bm.Mul(bm.Two, bm.Two)), bm.Seven), YSquared(bm.Two)) // 2³+7
}

func TestVRF_IsCurveXOrdinate(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	assert.True(t, IsCurveXOrdinate(big.NewInt(1)))
	assert.False(t, IsCurveXOrdinate(big.NewInt(5)))
}
