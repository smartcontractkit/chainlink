package utils

import (
	"encoding/json"
	"math"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecimal(t *testing.T) {
	t.Parallel()

	dec := decimal.New(1, 0)
	big := big.NewInt(1)

	var tt = []struct {
		v           any
		expectedErr bool
	}{
		{"1.1", false},
		{int(1), false},
		{int(-1), false},
		{int8(1), false},
		{int16(1), false},
		{int32(1), false},
		{int64(-1), false},
		{int32(-1), false},
		{uint(1), false},
		{uint8(1), false},
		{uint16(1), false},
		{uint32(1), false},
		{uint64(1), false},
		{float64(1.1), false},
		{float32(1.1), false},
		{float64(-1.1), false},
		{dec, false},
		{&dec, false},
		{big, false},
		{*big, false},
		{math.Inf(1), true},
		{math.Inf(-1), true},
		{float32(math.Inf(-1)), true},
		{float32(math.Inf(1)), true},
		{math.NaN(), true},
		{float32(math.NaN()), true},
		{true, true},
		{json.Number("1"), false},
		{json.Number("1.1"), false},
		{json.Number("-42"), false},
		{json.Number("115792089237316195423570985008687907853269984665640564039457584007913129639935"), false},
		{json.Number(""), true},
		{json.Number("not-a-number"), true},
	}
	for _, tc := range tt {
		_, err := ToDecimal(tc.v)
		if tc.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestToDecimal_JSONNumber(t *testing.T) {
	t.Parallel()

	got, err := ToDecimal(json.Number("1024"))
	require.NoError(t, err)
	assert.True(t, decimal.NewFromInt(1024).Equal(got))

	got, err = ToDecimal(json.Number("3.50"))
	require.NoError(t, err)
	assert.True(t, decimal.RequireFromString("3.50").Equal(got))

	got, err = ToDecimal(json.Number("-7"))
	require.NoError(t, err)
	assert.True(t, decimal.NewFromInt(-7).Equal(got))

	_, err = ToDecimal(json.Number("1e2"))
	require.NoError(t, err)

	_, err = ToDecimal(json.Number(""))
	require.Error(t, err)

	_, err = ToDecimal(json.Number("xyz"))
	require.Error(t, err)
}
