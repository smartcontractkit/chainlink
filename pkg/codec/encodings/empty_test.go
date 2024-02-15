package encodings_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
)

func TestEmpty(t *testing.T) {
	t.Parallel()
	t.Run("Encode returns prefix", func(t *testing.T) {
		anyPrefix := []byte{0x01, 0x02, 0x03}
		actual, err := encodings.Empty{}.Encode(struct{}{}, anyPrefix)
		require.NoError(t, err)
		assert.Equal(t, anyPrefix, actual)
	})

	t.Run("Decode returns all bytes as remaining", func(t *testing.T) {
		anySuffix := []byte{0x01, 0x02, 0x03}
		actual, remaining, err := encodings.Empty{}.Decode(anySuffix)
		require.NoError(t, err)
		assert.Equal(t, anySuffix, remaining)
		assert.Equal(t, struct{}{}, actual)
	})

	t.Run("GetType returns empty struct", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(struct{}{}), encodings.Empty{}.GetType())
	})

	t.Run("Size returns 0", func(t *testing.T) {
		size, err := encodings.Empty{}.Size(0)
		require.NoError(t, err)
		assert.Equal(t, 0, size)
	})

	t.Run("FixedSize returns 0", func(t *testing.T) {
		size, err := encodings.Empty{}.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, 0, size)
	})
}
