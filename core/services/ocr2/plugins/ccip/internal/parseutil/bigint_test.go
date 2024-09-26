package parseutil

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestParseBigIntFromAny(t *testing.T) {
	decimalVal := decimal.New(123, 0)

	testCases := []struct {
		name   string
		val    any
		res    *big.Int
		expErr bool
	}{
		{name: "nil", val: nil, expErr: true},
		{name: "string", val: "123", res: big.NewInt(123)},
		{name: "decimal", val: decimal.New(123, 0), res: big.NewInt(123)},
		{name: "decimal pointer", val: &decimalVal, res: big.NewInt(123)},
		{name: "int64", val: int64(123), res: big.NewInt(123)},
		{name: "int", val: 123, res: big.NewInt(123)},
		{name: "float", val: 123.12, res: big.NewInt(123)},
		{name: "uint8", val: uint8(12), expErr: true},
		{name: "struct", val: struct{ name string }{name: "asd"}, expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ParseBigIntFromAny(tc.val)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.res, res)
		})
	}
}
