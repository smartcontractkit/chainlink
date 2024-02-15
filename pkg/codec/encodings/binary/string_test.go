package binary_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/binary"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestString(t *testing.T) {
	s, createErr := binary.NewString(255, bi)
	require.NoError(t, createErr)
	s2, createErr := binary.NewString(256, bi)
	require.NoError(t, createErr)

	t.Run("Encode encodes as a slice with minimal leading bytes and can decode", func(t *testing.T) {
		encoded, err := s.Encode("foo", nil)
		require.NoError(t, err)

		size, err := bi.Int(1)
		require.NoError(t, err)
		bytesCodec, err := encodings.NewSlice(&binary.Uint8{}, size)
		require.NoError(t, err)
		expected, err := bytesCodec.Encode([]byte("foo"), nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded)

		encoded, err = s2.Encode("foo", nil)
		require.NoError(t, err)

		size, err = bi.Int(2)
		require.NoError(t, err)
		bytesCodec, err = encodings.NewSlice(&binary.Uint8{}, size)
		require.NoError(t, err)
		expected, err = bytesCodec.Encode([]byte("foo"), nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded)

		actual, remaining, err := s2.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, "foo", actual)
		assert.Empty(t, remaining)
	})

	t.Run("Encode respects prefix", func(t *testing.T) {
		prefix := []byte("bar")
		encoded, err := s.Encode("foo", prefix)
		require.NoError(t, err)

		rawEncoded, err := s.Encode("foo", nil)
		require.NoError(t, err)
		assert.Equal(t, append(prefix, rawEncoded...), encoded)
	})

	t.Run("Decode returns remaining", func(t *testing.T) {
		encoded, err := s.Encode("foo", nil)
		require.NoError(t, err)
		encoded = append(encoded, []byte("bar")...)

		actual, remaining, err := s.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, "foo", actual)
		assert.Equal(t, []byte("bar"), remaining)
	})

	t.Run("Encode returns an error if type is not a string", func(t *testing.T) {
		_, err := s.Encode(1, nil)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Encode returns an error if the value to encode is too long", func(t *testing.T) {
		_, err := s.Encode(string(make([]byte, 256)), nil)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
		_, err = s2.Encode(string(make([]byte, 257)), nil)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Decode returns an error if the encoded value is too long but fits in the buffer", func(t *testing.T) {
		s3, err := binary.NewString(258, bi)
		require.NoError(t, err)
		encoded, err := s3.Encode(string(make([]byte, 257)), nil)
		require.NoError(t, err)
		_, _, err = s2.Decode(encoded)
		assert.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("Decode returns an error if there are not enough bytes to decode", func(t *testing.T) {
		encoded, err := s.Encode("foo", nil)
		require.NoError(t, err)
		encoded = encoded[:len(encoded)-1]
		_, _, err = s.Decode(encoded)
		assert.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("Size returns an error", func(t *testing.T) {
		_, err := s.Size(10)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("FixedSize returns an error", func(t *testing.T) {
		_, err := s.FixedSize()
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("GetType returns string", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(""), s.GetType())
	})
}
