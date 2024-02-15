package binary_test

import (
	rawbinary "encoding/binary"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestFloat32(t *testing.T) {
	t.Parallel()
	f := bi.Float32()
	fVal := float32(1.234)
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		encoded, err := f.Encode(fVal, nil)

		require.NoError(t, err)
		expected := make([]byte, 4)
		rawbinary.BigEndian.PutUint32(expected, math.Float32bits(fVal))
		assert.Equal(t, expected, encoded)

		decoded, remaining, err := f.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, fVal, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		prefix := []byte{1, 2, 3}

		encoded, err := f.Encode(fVal, prefix)
		require.NoError(t, err)

		onlyFencoded, err := f.Encode(fVal, nil)
		require.NoError(t, err)

		assert.Equal(t, append([]byte{1, 2, 3}, onlyFencoded...), encoded)
	})

	t.Run("Encode returns an error if input is not a Float32", func(t *testing.T) {
		_, err := f.Encode("not a Float32", nil)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})
	t.Run("Decode leaves a suffix", func(t *testing.T) {
		onlyFencoded, err := f.Encode(fVal, nil)
		require.NoError(t, err)

		suffix := []byte{1, 2, 3}

		decoded, remaining, err := f.Decode(append(onlyFencoded, suffix...))
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, fVal, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		_, _, err := f.Decode([]byte{1, 2, 3})
		assert.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns the correct type", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(fVal), f.GetType())
	})

	t.Run("Size returns the correct size", func(t *testing.T) {
		size, err := f.Size(100)
		require.NoError(t, err)
		assert.Equal(t, 4, size)
	})

	t.Run("FixedSize returns the correct size", func(t *testing.T) {
		size, err := f.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, 4, size)
	})
}
