package mercury

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Values(t *testing.T) {
	t.Run("serializes max int192", func(t *testing.T) {
		encoded, err := EncodeValueInt192(MaxInt192)
		require.NoError(t, err)
		decoded, err := DecodeValueInt192(encoded)
		require.NoError(t, err)
		assert.Equal(t, MaxInt192, decoded)
	})
}
