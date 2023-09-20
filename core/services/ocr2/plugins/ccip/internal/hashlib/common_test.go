package hashlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesOfBytesKeccak(t *testing.T) {
	h, err := BytesOfBytesKeccak(nil)
	assert.NoError(t, err)
	assert.Equal(t, [32]byte{}, h)

	h1, err := BytesOfBytesKeccak([][]byte{{0x1}, {0x1}})
	assert.NoError(t, err)
	h2, err := BytesOfBytesKeccak([][]byte{{0x1, 0x1}})
	assert.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}
