package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_NonNegativeInt64ToUint64(t *testing.T) {
	tests := []struct {
		input    int64
		expected uint64
	}{
		{input: 42, expected: 42},
		{input: -1, expected: 0},
		{input: 0, expected: 0},
		{input: 9223372036854775807, expected: 9223372036854775807}, // Max int64
	}

	for _, test := range tests {
		result := utils.NonNegativeInt64ToUint64(test.input)
		assert.Equal(t, test.expected, result)
	}
}
