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
