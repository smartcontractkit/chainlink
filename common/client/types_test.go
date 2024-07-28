package client

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxDifficulty(t *testing.T) {
	cases := []struct {
		A, B, Result *big.Int
	}{
		{
			A: nil, B: nil, Result: nil,
		},
		{
			A: nil, B: big.NewInt(1), Result: big.NewInt(1),
		},
		{
			A: big.NewInt(1), B: big.NewInt(1), Result: big.NewInt(1),
		},
		{
			A: big.NewInt(1), B: big.NewInt(2), Result: big.NewInt(2),
		},
	}

	for _, test := range cases {
		actualResult := MaxTotalDifficulty(test.A, test.B)
		assert.Equal(t, test.Result, actualResult, "expected max(%v, %v) to produce %v", test.A, test.B, test.Result)
		inverted := MaxTotalDifficulty(test.B, test.A)
		assert.Equal(t, actualResult, inverted, "expected max(%v, %v) == max(%v, %v)", test.A, test.B, test.B, test.A)
	}
}
