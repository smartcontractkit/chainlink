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

func TestStructCodec(t *testing.T) {
	t.Parallel()
	t.Run("NewStructCodec returns an error if names are repeated", func(t *testing.T) {
		_, err := encodings.NewStructCodec([]encodings.NamedTypeCodec{
			{Name: "Foo", Codec: &testutils.TestTypeCodec{Value: 1}},
			{Name: "Bar", Codec: &testutils.TestTypeCodec{Value: 1}},
			{Name: "Foo", Codec: &testutils.TestTypeCodec{Value: 1}},
		})
		require.True(t, errors.Is(err, types.ErrInvalidConfig))
	})

	t.Run("NewStructCodec returns an error if names are invalid", func(t *testing.T) {
		_, err := encodings.NewStructCodec([]encodings.NamedTypeCodec{
			{Name: "Foo", Codec: &testutils.TestTypeCodec{Value: 1}},
			{Name: "Bar", Codec: &testutils.TestTypeCodec{Value: 1}},
			{Name: "", Codec: &testutils.TestTypeCodec{Value: 1}},
		})
		require.True(t, errors.Is(err, types.ErrInvalidConfig))
	})

	structCodec, createErr := encodings.NewStructCodec([]encodings.NamedTypeCodec{
		// use a pointer for one field to make sure it's not converted to a pointer to a pointer
		{Name: "Foo", Codec: &testutils.TestTypeCodec{Value: toPointer(int32(1)), Bytes: []byte{0x01}}},
		{Name: "Bar", Codec: &testutils.TestTypeCodec{Value: uint64(2), Bytes: []byte{0x02}}},
		{Name: "Baz", Codec: &testutils.TestTypeCodec{Value: "foo", Bytes: []byte("foo")}},
	})
	require.NoError(t, createErr)

	errCodec := &testutils.TestTypeCodec{Value: 1, Bytes: []byte{1, 2, 3}, Err: fmt.Errorf("%w: nope", types.ErrInvalidEncoding)}
	structCodecWithErr, createErr := encodings.NewStructCodec([]encodings.NamedTypeCodec{
		{Name: "Foo", Codec: errCodec},
	})
	require.NoError(t, createErr)

	t.Run("Encode returns the encoding of each element in-order", func(t *testing.T) {
		encoded, err := structCodec.Encode(&struct {
			Foo *int32
			Bar *uint64
			Baz *string
		}{
			Foo: toPointer(int32(1)),
			Bar: toPointer(uint64(2)),
			Baz: toPointer("foo"),
		}, nil)
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x02, 'f', 'o', 'o'}, encoded)
	})

	t.Run("Encode respects prefix", func(t *testing.T) {
		encoded, err := structCodec.Encode(&struct {
			Foo *int32
			Bar *uint64
			Baz *string
		}{
			Foo: toPointer(int32(1)),
			Bar: toPointer(uint64(2)),
			Baz: toPointer("foo"),
		}, []byte{0x03})
		require.NoError(t, err)
		require.Equal(t, []byte{0x03, 0x01, 0x02, 'f', 'o', 'o'}, encoded)
	})

	t.Run("Encode returns an error if the type is wrong", func(t *testing.T) {
		_, err := structCodec.Encode(&struct {
			Foo *int32
			Bar *uint64
			Baz *int
		}{
			Foo: toPointer(int32(1)),
			Bar: toPointer(uint64(2)),
			Baz: toPointer(3),
		}, nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Encode returns an error fields return an error", func(t *testing.T) {
		_, err := structCodecWithErr.Encode(&struct{ Foo *int }{Foo: toPointer(1)}, nil)
		assert.Equal(t, errCodec.Err, err)
	})

	t.Run("Decode returns the decoding of each element in-order", func(t *testing.T) {
		decoded, remaining, err := structCodec.Decode([]byte{0x01, 0x02, 'f', '0', '0'})
		require.NoError(t, err)
		require.Empty(t, remaining)
		require.Equal(t, &struct {
			Foo *int32
			Bar *uint64
			Baz *string
		}{
			Foo: toPointer(int32(1)),
			Bar: toPointer(uint64(2)),
			Baz: toPointer("foo"),
		}, decoded)
	})

	t.Run("Decode leaves extra data", func(t *testing.T) {
		decoded, remaining, err := structCodec.Decode([]byte{0x01, 0x02, 'f', '0', '0', 0x03})
		require.NoError(t, err)
		require.Equal(t, []byte{0x03}, remaining)
		require.Equal(t, &struct {
			Foo *int32
			Bar *uint64
			Baz *string
		}{
			Foo: toPointer(int32(1)),
			Bar: toPointer(uint64(2)),
			Baz: toPointer("foo"),
		}, decoded)
	})

	t.Run("Decode returns an error if there are not enough bytes to decode", func(t *testing.T) {
		_, _, err := structCodec.Decode([]byte{0x01})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("Decode returns an error fields return an error", func(t *testing.T) {
		_, _, err := structCodecWithErr.Decode([]byte{1, 2, 3})
		assert.Equal(t, errCodec.Err, err)
	})

	t.Run("GetType returns a type with elements in order", func(t *testing.T) {
		tpe := structCodec.GetType()
		expectedType := reflect.PointerTo(reflect.StructOf([]reflect.StructField{
			{Name: "Foo", Type: reflect.PointerTo(reflect.TypeOf(int32(0)))},
			{Name: "Bar", Type: reflect.PointerTo(reflect.TypeOf(uint64(0)))},
			{Name: "Baz", Type: reflect.PointerTo(reflect.TypeOf(""))},
		}))
		require.Equal(t, expectedType, tpe)
	})

	t.Run("Size returns size of elements", func(t *testing.T) {
		size, err := structCodec.Size(100)
		require.NoError(t, err)
		require.Equal(t, 5, size)
	})

	t.Run("Size returns error if elements return error", func(t *testing.T) {
		_, err := structCodecWithErr.Size(100)
		assert.Equal(t, errCodec.Err, err)
	})

	t.Run("FixedSize returns size of elements", func(t *testing.T) {
		size, err := structCodec.FixedSize()
		require.NoError(t, err)
		require.Equal(t, 5, size)
	})

	t.Run("FixedSize returns error if elements return error", func(t *testing.T) {
		_, err := structCodecWithErr.FixedSize()
		assert.Equal(t, errCodec.Err, err)
	})

	t.Run("SizeAtTopLevel returns Size from each element", func(t *testing.T) {
		size, err := structCodec.SizeAtTopLevel(100)
		require.NoError(t, err)
		require.Equal(t, 500, size)
	})

	t.Run("SizeAtTopLevel returns error if elements return error", func(t *testing.T) {
		_, err := structCodecWithErr.SizeAtTopLevel(100)
		assert.Equal(t, errCodec.Err, err)
	})
}

func toPointer[T any](t T) *T {
	return &t
}
