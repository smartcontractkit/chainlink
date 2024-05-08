// DO NOT MODIFY: automatically generated from chainlink-common/pkg/codec/encodings/binary/gen/main.go using the template int_gen_test.go

package binary_test

import (
	rawbinary "encoding/binary"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/binary"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var bi = binary.BigEndian()

func TestInt8(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Int8()
		value := int8(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := []byte{123}

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Int8()
		value := int8(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 1+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Int8()
		value := int8(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Int8()
		bytes := make([]byte, 1-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Int8()
		assert.Equal(t, i.GetType(), reflect.TypeOf(int8(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Int8().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Int8().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("returns an error if the input is not an uint8", func(t *testing.T) {
		i := bi.Int8()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestUint8(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Uint8()
		value := uint8(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := []byte{123}

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Uint8()
		value := uint8(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 1+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Uint8()
		value := uint8(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Uint8()
		bytes := make([]byte, 1-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Uint8()
		assert.Equal(t, i.GetType(), reflect.TypeOf(uint8(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Uint8().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Uint8().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("returns an error if the input is not an uint8", func(t *testing.T) {
		i := bi.Uint8()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestInt16(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Int16()
		value := int16(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := rawbinary.BigEndian.AppendUint16(nil, 123)

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Int16()
		value := int16(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 2+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Int16()
		value := int16(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Int16()
		bytes := make([]byte, 2-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Int16()
		assert.Equal(t, i.GetType(), reflect.TypeOf(int16(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Int16().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Int16().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("returns an error if the input is not an uint16", func(t *testing.T) {
		i := bi.Int16()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestUint16(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Uint16()
		value := uint16(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := rawbinary.BigEndian.AppendUint16(nil, 123)

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Uint16()
		value := uint16(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 2+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Uint16()
		value := uint16(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Uint16()
		bytes := make([]byte, 2-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Uint16()
		assert.Equal(t, i.GetType(), reflect.TypeOf(uint16(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Uint16().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Uint16().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("returns an error if the input is not an uint16", func(t *testing.T) {
		i := bi.Uint16()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestInt32(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Int32()
		value := int32(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := rawbinary.BigEndian.AppendUint32(nil, 123)

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Int32()
		value := int32(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 4+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Int32()
		value := int32(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Int32()
		bytes := make([]byte, 4-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Int32()
		assert.Equal(t, i.GetType(), reflect.TypeOf(int32(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Int32().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Int32().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("returns an error if the input is not an uint32", func(t *testing.T) {
		i := bi.Int32()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestUint32(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Uint32()
		value := uint32(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := rawbinary.BigEndian.AppendUint32(nil, 123)

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Uint32()
		value := uint32(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 4+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Uint32()
		value := uint32(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Uint32()
		bytes := make([]byte, 4-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Uint32()
		assert.Equal(t, i.GetType(), reflect.TypeOf(uint32(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Uint32().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Uint32().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("returns an error if the input is not an uint32", func(t *testing.T) {
		i := bi.Uint32()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestInt64(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Int64()
		value := int64(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := rawbinary.BigEndian.AppendUint64(nil, 123)

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Int64()
		value := int64(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 8+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Int64()
		value := int64(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Int64()
		bytes := make([]byte, 8-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Int64()
		assert.Equal(t, i.GetType(), reflect.TypeOf(int64(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Int64().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Int64().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("returns an error if the input is not an uint64", func(t *testing.T) {
		i := bi.Int64()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestUint64(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := bi.Uint64()
		value := uint64(123)

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)

		expected := rawbinary.BigEndian.AppendUint64(nil, 123)

		assert.Equal(t, expected, encoded)

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := bi.Uint64()
		value := uint64(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 8+3, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := bi.Uint64()
		value := uint64(123)
		suffix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, nil)
		require.NoError(t, err)
		encoded = append(encoded, suffix...)

		decoded, remaining, err := i.Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, suffix, remaining)
		assert.Equal(t, value, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes", func(t *testing.T) {
		i := bi.Uint64()
		bytes := make([]byte, 8-1)
		_, _, err := i.Decode(bytes)
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := bi.Uint64()
		assert.Equal(t, i.GetType(), reflect.TypeOf(uint64(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := bi.Uint64().Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := bi.Uint64().FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("returns an error if the input is not an uint64", func(t *testing.T) {
		i := bi.Uint64()

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}

func TestBuilderInt(t *testing.T) {
	t.Run("Wraps encoding and decoding for 8 bytes", func(t *testing.T) {
		codec, err := bi.Int(1)
		require.NoError(t, err)
		anyValue := 123

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 1, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Int(1)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Int(1)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("Wraps encoding and decoding for 16 bytes", func(t *testing.T) {
		codec, err := bi.Int(2)
		require.NoError(t, err)
		anyValue := 123

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 2, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Int(2)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Int(2)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("Wraps encoding and decoding for 32 bytes", func(t *testing.T) {
		codec, err := bi.Int(4)
		require.NoError(t, err)
		anyValue := 123

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 4, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Int(4)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Int(4)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("Wraps encoding and decoding for 64 bytes", func(t *testing.T) {
		codec, err := bi.Int(8)
		require.NoError(t, err)
		anyValue := 123

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 8, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Int(8)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Int(8)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("Wraps encoding and decoding for other sized bytes", func(t *testing.T) {
		codec, err := bi.Int(10)
		require.NoError(t, err)
		anyValue := 123

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 10, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("returns an error if the input is not an int", func(t *testing.T) {
		codec, err := bi.Int(4)
		require.NoError(t, err)

		_, err = codec.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("GetType returns int", func(t *testing.T) {
		codec, err := bi.Int(4)
		require.NoError(t, err)

		assert.Equal(t, reflect.TypeOf(0), codec.GetType())
	})
}

func TestGetUintTypeCodecByByteSize(t *testing.T) {
	t.Run("Wraps encoding and decoding for 8 bytes", func(t *testing.T) {
		codec, err := bi.Uint(1)
		require.NoError(t, err)
		anyValue := uint(123)

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 1, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(1)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(1)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("Wraps encoding and decoding for 16 bytes", func(t *testing.T) {
		codec, err := bi.Uint(2)
		require.NoError(t, err)
		anyValue := uint(123)

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 2, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(2)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(2)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 2)
	})

	t.Run("Wraps encoding and decoding for 32 bytes", func(t *testing.T) {
		codec, err := bi.Uint(4)
		require.NoError(t, err)
		anyValue := uint(123)

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 4, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(4)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(4)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 4)
	})

	t.Run("Wraps encoding and decoding for 64 bytes", func(t *testing.T) {
		codec, err := bi.Uint(8)
		require.NoError(t, err)
		anyValue := uint(123)

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 8, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(8)
		require.NoError(t, err)
		size, err := codec.Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		codec, err := bi.Uint(8)
		require.NoError(t, err)
		size, err := codec.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 8)
	})

	t.Run("Wraps encoding and decoding for other sized bytes", func(t *testing.T) {
		codec, err := bi.Uint(10)
		require.NoError(t, err)
		anyValue := uint(123)

		encoded, err := codec.Encode(anyValue, nil)
		require.NoError(t, err)
		require.Equal(t, 10, len(encoded))

		decoded, remaining, err := codec.Decode(encoded)
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, anyValue, decoded)
	})

	t.Run("returns an error if the input is not an int", func(t *testing.T) {
		codec, err := bi.Uint(4)
		require.NoError(t, err)

		_, err = codec.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("GetType returns uint", func(t *testing.T) {
		codec, err := bi.Uint(4)
		require.NoError(t, err)

		assert.Equal(t, reflect.TypeOf(uint(0)), codec.GetType())
	})
}
