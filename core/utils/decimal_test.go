package utils

import (
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

func TestDecimalUint64NoOverflow(t *testing.T) {
	t.Parallel()

	// Values above math.MaxInt64 were silently corrupted: int64(math.MaxUint64) == -1.
	want, _ := decimal.NewFromString("18446744073709551615") // math.MaxUint64
	got, err := ToDecimal(uint64(math.MaxUint64))
	require.NoError(t, err)
	assert.True(t, got.Equal(want), "uint64(MaxUint64): got %s, want %s", got, want)

	got, err = ToDecimal(uint(math.MaxUint64))
	require.NoError(t, err)
	assert.True(t, got.Equal(want), "uint(MaxUint64): got %s, want %s", got, want)
}
