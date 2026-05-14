package utils

import (
	"encoding/json"
	"math"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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
		// json.Number is produced by Decoder.UseNumber, which the pipeline
		// uses when decoding JSON HTTP responses. Both integer and floating
		// point forms have to be accepted; see #8504.
		{json.Number("123"), false},
		{json.Number("-1.1"), false},
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
