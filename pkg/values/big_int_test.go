package values

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BigIntUnwrapTo(t *testing.T) {
	expected := big.NewInt(100)
	v := NewBigInt(expected)

	got := new(big.Int)
	err := v.UnwrapTo(got)
	require.NoError(t, err)

	assert.Equal(t, expected, got)

	gotInt := (*big.Int)(nil)
	err = v.UnwrapTo(gotInt)
	assert.ErrorContains(t, err, "cannot unwrap to nil pointer")

	var varAny any
	err = v.UnwrapTo(&varAny)
	require.NoError(t, err)
	assert.Equal(t, expected, varAny)

	var varStr string
	err = v.UnwrapTo(&varStr)
	assert.ErrorContains(t, err, "cannot unwrap to value of type: *string")
}

func Test_BigInt(t *testing.T) {
	testCases := []struct {
		name string
		bi   *big.Int
	}{
		{
			name: "positive",
			bi:   big.NewInt(100),
		},
		{
			name: "0",
			bi:   big.NewInt(0),
		},
		{
			name: "negative",
			bi:   big.NewInt(-1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := NewBigInt(tc.bi)

			vp := Proto(v)
			got := FromProto(vp)

			assert.Equal(t, tc.bi, got.(*BigInt).Underlying)
		})
	}
}
