package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimal(t *testing.T) {
	var tt = []struct {
		v           interface{}
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
		{uint64(1), false},
		{float64(1.1), false},
		{float32(1.1), false},
		{float64(-1.1), false},
		{math.Inf(1), true},
		{math.Inf(-1), true},
		{float32(math.Inf(-1)), true},
		{float32(math.Inf(1)), true},
		{math.NaN(), true},
		{float32(math.NaN()), true},
	}
	for _, tc := range tt {
		tc := tc
		_, err := ToDecimal(tc.v)
		if tc.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
