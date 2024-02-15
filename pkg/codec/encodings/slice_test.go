package encodings_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/binary"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/testutils"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestSlice(t *testing.T) {
	t.Parallel()
	anyValue := 0x34
	elementCodec := &testutils.TestTypeCodec{
		Value: anyValue,
		Bytes: []byte{0x03, 0x04},
	}
	sizeCodec := &testutils.TestTypeCodec{
		Value: 3,
		Bytes: []byte{3},
	}
	testSlice, sliceCreateErr := encodings.NewSlice(elementCodec, sizeCodec)
	require.NoError(t, sliceCreateErr)

	anyErr := fmt.Errorf("%w: testing", types.ErrInvalidType)
	errCodec := &testutils.TestTypeCodec{
		Value: 1,
		Err:   anyErr,
	}
	errorSlice, sliceCreateErr := encodings.NewSlice(errCodec, sizeCodec)
	require.NoError(t, sliceCreateErr)

	errorSizeSlice, sliceCreateErr := encodings.NewSlice(elementCodec, errCodec)
	require.NoError(t, sliceCreateErr)

	t.Run("NewSlice returns error if args are nil", func(t *testing.T) {
		_, err := encodings.NewSlice(nil, sizeCodec)
		require.True(t, errors.Is(err, types.ErrInvalidConfig))

		_, err = encodings.NewSlice(elementCodec, nil)
		require.True(t, errors.Is(err, types.ErrInvalidConfig))
	})

	t.Run("NewSlice returns error if size arg is not an int", func(t *testing.T) {
		_, err := encodings.NewSlice(nil, &binary.Int8{})
		require.True(t, errors.Is(err, types.ErrInvalidConfig))
	})

	t.Run("GetType returns slice of underlying type", func(t *testing.T) {
		assert.Equal(t, testSlice.GetType(), reflect.SliceOf(reflect.TypeOf(anyValue)))
	})

	t.Run("Encode prefixes encoded slice with length and encodes elements", func(t *testing.T) {
		encoded, err := testSlice.Encode([]int{anyValue, anyValue, anyValue}, []byte{})
		require.NoError(t, err)
		assert.Equal(t, []byte{0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04}, encoded)
	})

	t.Run("Encode prefixes encodes array", func(t *testing.T) {
		encoded, err := testSlice.Encode([3]int{anyValue, anyValue, anyValue}, []byte{0xFE, 0xED})
		require.NoError(t, err)
		assert.Equal(t, []byte{0xFE, 0xED, 0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04}, encoded)
	})

	t.Run("Encode returns an error if slice is too large for size bytes", func(t *testing.T) {
		tmp := make([]int, math.MaxUint16+1)
		_, err := testSlice.Encode(tmp, []byte{})
		require.True(t, errors.Is(err, types.ErrSliceWrongLen))
	})

	t.Run("Decode returns slice", func(t *testing.T) {
		actual, remaining, err := testSlice.Decode([]byte{0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04})
		require.NoError(t, err)
		assert.Equal(t, []int{anyValue, anyValue, anyValue}, actual)
		assert.Equal(t, []byte{}, remaining)
	})

	t.Run("Decode returns remaining bytes", func(t *testing.T) {
		actual, remaining, err := testSlice.Decode([]byte{0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04, 0xCA, 0xFE})
		require.NoError(t, err)
		assert.Equal(t, []int{anyValue, anyValue, anyValue}, actual)
		assert.Equal(t, []byte{0xCA, 0xFE}, remaining)
	})

	t.Run("Decode returns an error if the size codec returns a negative number", func(t *testing.T) {
		negSizeCodec := &testutils.TestTypeCodec{
			Value: 3,
			Bytes: []byte{3},
		}
		ts, err := encodings.NewSlice(elementCodec, negSizeCodec)
		require.NoError(t, err)
		_, _, err = ts.Decode([]byte{0xFF, 0xFF, 0xFF, 0xFF})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("Decode returns an error if the size codec non int type", func(t *testing.T) {
		negSizeCodec := &testutils.TestTypeCodec{
			Value: 2,
			Bytes: []byte{3},
		}
		ts, err := encodings.NewSlice(elementCodec, negSizeCodec)
		require.NoError(t, err)
		negSizeCodec.Value = "foo"
		_, _, err = ts.Decode([]byte{0xFF, 0xFF, 0xFF, 0xFF})
		require.True(t, errors.Is(err, types.ErrInternal))
	})

	t.Run("Encode returns an error if the underlying type returns an error", func(t *testing.T) {
		_, err := errorSlice.Encode([]int{anyValue}, []byte{})
		require.Equal(t, anyErr, err)
	})

	t.Run("Encode returns errors if the size codec returns an error", func(t *testing.T) {
		_, err := errorSizeSlice.Encode([]int{anyValue}, []byte{})
		require.Equal(t, anyErr, err)
	})

	t.Run("Encode returns an error if the value is not a slice or array", func(t *testing.T) {
		_, err := testSlice.Encode(anyValue, []byte{})
		require.True(t, errors.Is(err, types.ErrNotASlice))
	})

	t.Run("Decode returns an error if the underlying type returns an error", func(t *testing.T) {
		_, _, err := errorSlice.Decode([]byte{0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04})
		require.Equal(t, anyErr, err)
	})

	t.Run("Decode returns an error if the size returns an error", func(t *testing.T) {
		_, _, err := errorSizeSlice.Decode([]byte{0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04})
		require.Equal(t, anyErr, err)
	})

	t.Run("Decode returns an error if there are not enough bytes for the size", func(t *testing.T) {
		_, _, err := testSlice.Decode([]byte{1})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("Size returns numElements * the element size + the size of the size", func(t *testing.T) {
		size, err := testSlice.Size(3)
		require.NoError(t, err)
		assert.Equal(t, 3*len(elementCodec.Bytes)+len(sizeCodec.Bytes), size)
	})

	t.Run("Size returns errors from element codec", func(t *testing.T) {
		_, err := errorSlice.Size(3)
		require.Equal(t, anyErr, err)
	})

	t.Run("Size returns errors from size codec", func(t *testing.T) {
		_, err := errorSizeSlice.Size(3)
		require.Equal(t, anyErr, err)
	})

	t.Run("FixedSize returns an error", func(t *testing.T) {
		_, err := testSlice.FixedSize()
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}
