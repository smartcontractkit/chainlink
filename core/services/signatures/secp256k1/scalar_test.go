package secp256k1

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.dedis.ch/kyber"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

var numScalarSamples = 10

var observedScalars map[string]bool

func init() {
	observedScalars = make(map[string]bool)
}

// observedScalar ensures that novel scalars are being picked.
func observedScalar(t *testing.T, s kyber.Scalar) {
	data, err := s.(*secp256k1Scalar).modG().MarshalBinary()
	require.NoError(t, err)
	scalar := hex.Dump(data)
	require.False(t, observedScalars[scalar])
	observedScalars[scalar] = true
}

var randomStreamScalar = cltest.NewStream(&testing.T{}, 0)

func TestScalar_SetAndEqual(t *testing.T) {
	tests := []int64{5, 67108864, 67108865, 4294967295}
	g := newScalar(scalarZero)
	for _, test := range tests {
		f := newScalar(big.NewInt(test))
		g.Set(f)
		assert.Equal(t, f, g,
			"the method Set should give the same value to receiver")
		f.Add(f, newScalar(big.NewInt(1)))
		assert.NotEqual(t, f, g,
			"SetInt should take a copy of the backing big.Int")
	}
}

func TestNewScalar(t *testing.T) {
	one := newScalar(big.NewInt(1))
	assert.Equal(t, ToInt(one),
		ToInt(newScalar(big.NewInt(0).Add(ToInt(one), GroupOrder))),
		"equivalence classes mod GroupOrder not equal")
}

func TestScalar_SmokeTestPick(t *testing.T) {
	f := newScalar(scalarZero).Clone()
	for i := 0; i < numScalarSamples; i++ {
		f.Pick(randomStreamScalar)
		observedScalar(t, f)
		require.True(t, ToInt(f).Cmp(big.NewInt(1000000000)) == 1,
			"implausibly low value returned from Pick: %v", f)
	}
}

func TestScalar_Neg(t *testing.T) {
	f := newScalar(scalarZero).Clone()
	for i := 0; i < numScalarSamples; i++ {
		f.Pick(randomStreamScalar)
		observedScalar(t, f)
		g := f.Clone()
		g.Neg(g)
		require.True(t, g.Add(f, g).Equal(newScalar(scalarZero)))
	}
}

func TestScalar_Sub(t *testing.T) {
	f := newScalar(scalarZero).Clone()
	for i := 0; i < numScalarSamples; i++ {
		f.Pick(randomStreamScalar)
		observedScalar(t, f)
		require.True(t, f.Sub(f, f).Equal(newScalar(scalarZero)),
			"subtracting something from itself should give zero")
	}
}

func TestScalar_Clone(t *testing.T) {
	f := newScalar(big.NewInt(1))
	g := f.Clone()
	h := f.Clone()
	assert.Equal(t, f, g, "clone output does not equal input")
	g.Add(f, f)
	assert.Equal(t, f, h, "clone does not make a copy")
}

func TestScalar_Marshal(t *testing.T) {
	f := newScalar(scalarZero)
	g := newScalar(scalarZero)
	for i := 0; i < numFieldSamples; i++ {
		f.Pick(randomStreamScalar)
		observedScalar(t, f)
		data, err := f.MarshalBinary()
		require.Nil(t, err)
		err = g.UnmarshalBinary(data)
		require.Nil(t, err)
		require.True(t, g.Equal(f),
			"roundtrip through serialization should give same "+
				"result back: failed with %s", f)
	}
}

func TestScalar_MulDivInv(t *testing.T) {
	f := newScalar(scalarZero)
	g := newScalar(scalarZero)
	h := newScalar(scalarZero)
	j := newScalar(scalarZero)
	k := newScalar(scalarZero)
	for i := 0; i < numFieldSamples; i++ {
		f.Pick(randomStreamScalar)
		observedScalar(t, f)
		g.Inv(f)
		h.Mul(f, g)
		require.True(t, h.Equal(newScalar(big.NewInt(1))))
		h.Div(f, f)
		require.True(t, h.Equal(newScalar(big.NewInt(1))))
		h.Div(newScalar(big.NewInt(1)), f)
		require.True(t, h.Equal(g))
		h.Pick(randomStreamScalar)
		observedScalar(t, h)
		j.Neg(j.Mul(h, f))
		k.Mul(h, k.Neg(f))
		require.True(t, j.Equal(k), "-(h*f) != h*(-f)")
	}
}
