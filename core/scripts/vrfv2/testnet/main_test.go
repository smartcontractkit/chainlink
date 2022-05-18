package main

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name        string
		top, bottom int64
		result      int64
	}{
		{
			name:   "zero 1",
			bottom: 0,
			top:    100,
			result: 0,
		},
		{
			name:   "zero 2",
			bottom: 0,
			top:    99,
			result: 0,
		},
		{
			name:   "one",
			bottom: 0,
			top:    100,
			result: 1,
		},
		{
			name:   "one2",
			bottom: 0,
			top:    99,
			result: 1,
		},
		{
			name:   "mid",
			bottom: 0,
			top:    159,
			result: 80,
		},
		{
			name:   "mid 2",
			bottom: 0,
			top:    159,
			result: 81,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			testFunc := func(val *big.Int) bool {
				return val.Cmp(big.NewInt(test.result)) < 1
			}

			result := binarySearch(big.NewInt(test.top), big.NewInt(test.bottom), testFunc)
			assert.Equal(t, test.result, result.Int64())
		})
	}
}
