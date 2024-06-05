package ccipocr3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBytes32FromString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Bytes32
		expErr   bool
	}{
		{
			name:     "valid input",
			input:    "0x200000000000000000000000",
			expected: Bytes32{0x20, 0},
			expErr:   false,
		},
		{
			name:     "invalid input",
			input:    "lrfv",
			expected: Bytes32{},
			expErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NewBytes32FromString(tc.input)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
