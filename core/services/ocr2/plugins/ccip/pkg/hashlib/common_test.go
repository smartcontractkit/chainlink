package hashlib

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func TestBytesOfBytesKeccak(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		h, err := BytesOfBytesKeccak(nil)
		assert.NoError(t, err)
		assert.Equal(t, [32]byte{}, h)

		h1, err := BytesOfBytesKeccak([][]byte{{0x1}, {0x1}})
		assert.NoError(t, err)
		h2, err := BytesOfBytesKeccak([][]byte{{0x1, 0x1}})
		assert.NoError(t, err)
		assert.NotEqual(t, h1, h2)
	})

	t.Run("should not have collision", func(t *testing.T) {
		s1 := utils.RandomBytes32()
		s2 := utils.RandomBytes32()

		hs1, err := BytesOfBytesKeccak([][]byte{s1[:]})
		assert.NoError(t, err)

		h1, err := BytesOfBytesKeccak([][]byte{s1[:], s2[:]})
		assert.NoError(t, err)

		h2, err := BytesOfBytesKeccak([][]byte{append(hs1[:], s2[:]...)})
		assert.NoError(t, err)

		assert.NotEqual(t, h1, h2)
	})

}

func TestBytesOfBytesEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]byte
		expected []byte
	}{
		{
			name:  "0 sub array, empty slice",
			input: [][]byte{},
			expected: []byte{
				0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		{
			name: "1 sub array",
			input: [][]byte{
				{0, 1, 2, 3},
			},
			expected: []byte{
				0, 0, 0, 0, 0, 0, 0, 1,
				0, 0, 0, 0, 0, 0, 0, 4, 0, 1, 2, 3,
			},
		},
		{
			name: "3 sub array",
			input: [][]byte{
				{0, 1, 2, 3},
				{0, 0, 0},
				{7, 8},
			},
			expected: []byte{
				0, 0, 0, 0, 0, 0, 0, 3,
				0, 0, 0, 0, 0, 0, 0, 4, 0, 1, 2, 3,
				0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 2, 7, 8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := encodeBytesOfBytes(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, encoded)
		})
	}
}
