package secp256k1

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/cryptotest"
)

var numFieldSamples = 10

var observedFieldElts map[string]bool

func init() {
	observedFieldElts = make(map[string]bool)
}

// observedFieldElt ensures that novel scalars are being picked.
func observedFieldElt(t *testing.T, s *fieldElt) {
	elt := s.Bytes()
	data := hex.Dump(elt[:])
	require.False(t, observedFieldElts[data])
	observedFieldElts[data] = true
}

var randomStream = cryptotest.NewStream(&testing.T{}, 0)

func TestField_SetIntAndEqual(t *testing.T) {
	tests := []int64{5, 67108864, 67108865, 4294967295}
	g := newFieldZero()
	for _, test := range tests {
		f := fieldEltFromInt(test)
		i := big.NewInt(test)
		g.SetInt(i)
		assert.Equal(t, f, g,
			"different values obtained for same input, using "+
				"SetInt vs fieldEltFromInt")
		i.Add(i, big.NewInt(1))
		assert.Equal(t, f, g,
			"SetInt should take a copy of the backing big.Int")
	}
}

func TestField_String(t *testing.T) {
	require.Equal(t, "fieldElt{0}", fieldZero.String())
}

func TestField_Equal(t *testing.T) {
	require.True(t, (*fieldElt)(nil).Equal((*fieldElt)(nil)))
	require.False(t, (*fieldElt)(nil).Equal(fieldZero))
	require.False(t, fieldZero.Equal((*fieldElt)(nil)))
}

func TestField_Set(t *testing.T) {
	f := fieldEltFromInt(1)
	g := newFieldZero()
	g.Set(f)
	g.Add(g, fieldEltFromInt(1))
	assert.Equal(t, f, fieldEltFromInt(1),
		"Set takes a copy of the backing big.Int")
}

func TestFieldEltFromInt(t *testing.T) {
	assert.Equal(t, fieldEltFromInt(1), // Also tests fieldElt.modQ
		fieldEltFromBigInt(new(big.Int).Add(q, big.NewInt(1))),
		"only one representation of a ℤ/qℤ element should be used")
}

func TestField_SmokeTestPick(t *testing.T) {
	f := newFieldZero()
	f.Pick(randomStream)
	observedFieldElt(t, f)
	assert.Equal(t, f.int().Cmp(big.NewInt(1000000000)), 1,
		"should be greater than 1000000000, with very high probability")
}

func TestField_Neg(t *testing.T) {
	f := newFieldZero()
	for i := 0; i < numFieldSamples; i++ {
		f.Pick(randomStream)
		observedFieldElt(t, f)
		g := f.Clone()
		g.Neg(g)
		require.True(t, g.Add(f, g).Equal(fieldZero),
			"adding something to its negative should give zero: "+
				"failed with %s", f)
	}
}

func TestField_Sub(t *testing.T) {
	f := newFieldZero()
	for i := 0; i < numFieldSamples; i++ {
		f.Pick(randomStream)
		observedFieldElt(t, f)
		require.True(t, f.Sub(f, f).Equal(fieldZero),
			"subtracting something from itself should give zero: "+
				"failed with %s", f)
	}
}

func TestField_Clone(t *testing.T) {
	f := fieldEltFromInt(1)
	g := f.Clone()
	h := f.Clone()
	assert.Equal(t, f, g, "clone output does not equal original")
	g.Add(f, f)
	assert.Equal(t, f, h, "clone does not make a copy")
}

func TestField_SetBytesAndBytes(t *testing.T) {
	f := newFieldZero()
	g := newFieldZero()
	for i := 0; i < numFieldSamples; i++ {
		f.Pick(randomStream)
		observedFieldElt(t, f)
		g.SetBytes(f.Bytes())
		require.True(t, g.Equal(f),
			"roundtrip through serialization should give same "+
				"result back: failed with %s", f)
	}
}

func TestField_MaybeSquareRootInField(t *testing.T) {
	f := newFieldZero()
	minusOne := fieldEltFromInt(-1)
	assert.Nil(t, maybeSqrtInField(minusOne), "-1 is not a square, in this field")
	for i := 0; i < numFieldSamples; i++ {
		f.Pick(randomStream)
		observedFieldElt(t, f)
		require.Equal(t, f.int().Cmp(q), -1, "picked larger value than q: %s", f)
		require.NotEqual(t, f.int().Cmp(big.NewInt(-1)), -1,
			"backing int must be non-negative")
		s := fieldSquare(f)
		g := maybeSqrtInField(s)
		require.NotEqual(t, (*fieldElt)(nil), g)
		ng := newFieldZero().Neg(g)
		require.True(t, f.Equal(g) || f.Equal(ng), "squaring something and "+
			"taking the square root should give ± the original: failed with %s", f)
		bigIntSqrt := newFieldZero() // Cross-check against big.ModSqrt
		rv := bigIntSqrt.int().ModSqrt(s.int(), q)
		require.NotNil(t, rv)
		require.True(t, bigIntSqrt.Equal(g) || bigIntSqrt.Equal(ng))
		nonSquare := newFieldZero().Neg(s)
		rv = bigIntSqrt.int().ModSqrt(nonSquare.int(), q)
		require.Nil(t, rv, "ModSqrt indicates nonSquare is square")
		require.Nil(t, maybeSqrtInField(nonSquare), "the negative of square "+
			"should not be a square")
	}
}

func TestField_RightHandSide(t *testing.T) {
	assert.Equal(t, rightHandSide(fieldEltFromInt(1)), fieldEltFromInt(8))
	assert.Equal(t, rightHandSide(fieldEltFromInt(2)), fieldEltFromInt(15))
}
