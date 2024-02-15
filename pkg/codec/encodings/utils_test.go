package encodings_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/testutils"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestSafeDecode(t *testing.T) {
	t.Run("runs decoding function if there are enough bytes", func(t *testing.T) {
		anyValue := 100
		actual, remaining, err := encodings.SafeDecode[int]([]byte{1, 2, 3, 4}, 4, func(bytes []byte) int {
			return anyValue
		})
		require.NoError(t, err)
		assert.Equal(t, anyValue, actual)
		assert.Empty(t, remaining)
	})

	t.Run("runs returns remaining bytes", func(t *testing.T) {
		anyValue := 100
		actual, remaining, err := encodings.SafeDecode[int]([]byte{1, 2, 3, 4, 5, 6}, 4, func(bytes []byte) int {
			return anyValue
		})
		require.NoError(t, err)
		assert.Equal(t, anyValue, actual)
		assert.Equal(t, []byte{5, 6}, remaining)
	})

	t.Run("returns error if there are not enough bytes", func(t *testing.T) {
		_, _, err := encodings.SafeDecode[int]([]byte{1, 2, 3}, 4, func(bytes []byte) int {
			require.Fail(t, "method must not be called")
			return 0
		})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})
}

func TestEncodeEach(t *testing.T) {
	t.Run("encodes each element of a slice", func(t *testing.T) {
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		val := []any{codec.Value, codec.Value, codec.Value}

		encoded, err := encodings.EncodeEach(reflect.ValueOf(val), nil, codec)
		require.NoError(t, err)
		assert.Equal(t, []byte{1, 2, 1, 2, 1, 2}, encoded)
	})

	t.Run("encodes each element of an array", func(t *testing.T) {
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		val := [3]any{codec.Value, codec.Value, codec.Value}

		encoded, err := encodings.EncodeEach(reflect.ValueOf(val), nil, codec)
		require.NoError(t, err)
		assert.Equal(t, []byte{1, 2, 1, 2, 1, 2}, encoded)
	})

	t.Run("leading bytes are preserved", func(t *testing.T) {
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		val := []any{codec.Value, codec.Value, codec.Value}

		prefix := []byte{99, 100}
		encoded, err := encodings.EncodeEach(reflect.ValueOf(val), prefix, codec)
		require.NoError(t, err)
		assert.Equal(t, []byte{99, 100, 1, 2, 1, 2, 1, 2}, encoded)
	})

	t.Run("returns error if value is not a slice or an array", func(t *testing.T) {
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}

		_, err := encodings.EncodeEach(reflect.ValueOf(100), nil, codec)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("returns error if any encoding fails", func(t *testing.T) {
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
			Err:   fmt.Errorf("%w: some error", types.ErrInvalidEncoding),
		}
		val := []any{codec.Value, codec.Value, codec.Value}

		_, err := encodings.EncodeEach(reflect.ValueOf(val), nil, codec)
		assert.Equal(t, codec.Err, err)
	})
}

func TestDecodeEach(t *testing.T) {
	t.Run("decodes each element of a slice", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6}
		into := make([]int, 3)
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		val, remaining, err := encodings.DecodeEach(bytes, reflect.ValueOf(into), 3, codec)
		require.NoError(t, err)
		assert.Empty(t, remaining)
		assert.Equal(t, []int{100, 100, 100}, val)
	})

	t.Run("decodes each element of an array", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6}
		into := [3]int{}
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		val, remaining, err := encodings.DecodeEach(bytes, reflect.ValueOf(&into).Elem(), 3, codec)
		require.NoError(t, err)
		assert.Empty(t, remaining)
		assert.Equal(t, [3]int{100, 100, 100}, val)
	})

	t.Run("remaining bytes are returned", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		into := make([]int, 3)
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		val, remaining, err := encodings.DecodeEach(bytes, reflect.ValueOf(into), 3, codec)
		require.NoError(t, err)
		assert.Equal(t, []byte{7, 8}, remaining)
		assert.Equal(t, []int{100, 100, 100}, val)
	})

	t.Run("returns error if value is not a slice or an array", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6}
		into := "foo"
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		_, _, err := encodings.DecodeEach(bytes, reflect.ValueOf(&into), 3, codec)
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("returns error if size too small", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6}
		into := make([]int, 2)
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		_, _, err := encodings.DecodeEach(bytes, reflect.ValueOf(into), 3, codec)
		assert.True(t, errors.Is(err, types.ErrSliceWrongLen))
	})

	t.Run("returns error if any encoding fails", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6}
		into := make([]int, 3)
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
			Err:   fmt.Errorf("%w: some error", types.ErrInvalidEncoding),
		}
		_, _, err := encodings.DecodeEach(bytes, reflect.ValueOf(into), 3, codec)
		assert.Equal(t, codec.Err, err)
	})

	t.Run("returns error if index cannot be set", func(t *testing.T) {
		bytes := []byte{1, 2, 3, 4, 5, 6}
		into := [3]int{}
		codec := &testutils.TestTypeCodec{
			Bytes: []byte{1, 2},
			Value: 100,
		}
		_, _, err := encodings.DecodeEach(bytes, reflect.ValueOf(into), 3, codec)
		assert.True(t, errors.Is(err, types.ErrInternal))
	})
}

func TestIndirectIfPointer(t *testing.T) {
	t.Run("returns the value if it is not a pointer", func(t *testing.T) {
		i := 100
		v, err := encodings.IndirectIfPointer(reflect.ValueOf(i))
		require.NoError(t, err)
		assert.Equal(t, i, v.Interface())
	})

	t.Run("returns the dereferenced value if it is a pointer", func(t *testing.T) {
		i := 100
		v, err := encodings.IndirectIfPointer(reflect.ValueOf(&i))
		require.NoError(t, err)
		assert.Equal(t, i, v.Interface())
	})

	t.Run("returns error if the pointer is nil", func(t *testing.T) {
		_, err := encodings.IndirectIfPointer(reflect.ValueOf((*int)(nil)))
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})
}
