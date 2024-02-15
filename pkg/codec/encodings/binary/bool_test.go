package binary_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/binary"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestBool(t *testing.T) {
	t.Parallel()
	b := binary.Bool{}
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		tEncoded, err := b.Encode(true, nil)
		require.NoError(t, err)

		fEncoded, err := b.Encode(false, nil)
		require.NoError(t, err)

		assert.Equal(t, []byte{1}, tEncoded)
		assert.Equal(t, []byte{0}, fEncoded)

		tDecoded, remaining, err := b.Decode(tEncoded)
		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, true, tDecoded)

		fDecoded, remaining, err := b.Decode(fEncoded)
		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, false, fDecoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		prefix := []byte{1, 2, 3}

		encoded, err := b.Encode(true, prefix)

		require.NoError(t, err)
		assert.Equal(t, []byte{1, 2, 3, 1}, encoded)
	})

	t.Run("Encode returns an error if input is not a bool", func(t *testing.T) {
		_, err := b.Encode("not a bool", nil)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})
	t.Run("Decode leaves a suffix", func(t *testing.T) {
		suffix := []byte{1, 2, 3}

		decoded, remaining, err := b.Decode(append([]byte{1}, suffix...))
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, true, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		_, _, err := b.Decode([]byte{})
		assert.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns the correct type", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(true), b.GetType())
	})

	t.Run("Size returns the correct size", func(t *testing.T) {
		size, err := b.Size(100)
		require.NoError(t, err)
		assert.Equal(t, 1, size)
	})

	t.Run("FixedSize returns the correct size", func(t *testing.T) {
		size, err := b.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, 1, size)
	})
}
