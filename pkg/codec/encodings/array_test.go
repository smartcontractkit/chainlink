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

func TestArray(t *testing.T) {
	t.Parallel()
	anyValue := 0x34
	elementCodec := &testutils.TestTypeCodec{
		Value: anyValue,
		Bytes: []byte{0x03, 0x04},
	}
	numElements := 3
	testArray, arrayCreateErr := encodings.NewArray(numElements, elementCodec)
	require.NoError(t, arrayCreateErr)

	anyErr := fmt.Errorf("%w: testing", types.ErrInvalidType)
	errCodec := &testutils.TestTypeCodec{
		Value: 1,
		Err:   anyErr,
	}
	errorArray, arrayCreateErr := encodings.NewArray(numElements, errCodec)
	require.NoError(t, arrayCreateErr)

	t.Run("NewArray returns error if elements is nil", func(t *testing.T) {
		_, err := encodings.NewArray(0, nil)
		require.True(t, errors.Is(err, types.ErrInvalidConfig))
	})

	t.Run("GetType returns Array of underlying type", func(t *testing.T) {
		assert.Equal(t, testArray.GetType(), reflect.ArrayOf(numElements, reflect.TypeOf(anyValue)))
	})

	t.Run("Encode prefixes encoded Array with length and encodes elements", func(t *testing.T) {
		encoded, err := testArray.Encode([]int{anyValue, anyValue, anyValue}, []byte{})
		require.NoError(t, err)
		assert.Equal(t, []byte{0x03, 0x04, 0x03, 0x04, 0x03, 0x04}, encoded)
	})

	t.Run("Encode prefixes encodes array", func(t *testing.T) {
		encoded, err := testArray.Encode([3]int{anyValue, anyValue, anyValue}, []byte{0xFE, 0xED})
		require.NoError(t, err)
		assert.Equal(t, []byte{0xFE, 0xED, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04}, encoded)
	})

	t.Run("Encode returns an error if Array is the wrong size", func(t *testing.T) {
		_, err := testArray.Encode([2]int{anyValue, anyValue}, []byte{})
		require.True(t, errors.Is(err, types.ErrSliceWrongLen))
	})

	t.Run("Decode returns Array", func(t *testing.T) {
		actual, remaining, err := testArray.Decode([]byte{0x03, 0x04, 0x03, 0x04, 0x03, 0x04})
		require.NoError(t, err)
		assert.Equal(t, [3]int{anyValue, anyValue, anyValue}, actual)
		assert.Equal(t, []byte{}, remaining)
	})

	t.Run("Decode returns remaining bytes", func(t *testing.T) {
		actual, remaining, err := testArray.Decode([]byte{0x03, 0x04, 0x03, 0x04, 0x03, 0x04, 0xCA, 0xFE})
		require.NoError(t, err)
		assert.Equal(t, [3]int{anyValue, anyValue, anyValue}, actual)
		assert.Equal(t, []byte{0xCA, 0xFE}, remaining)
	})

	t.Run("Encode returns an error if the underlying type returns an error", func(t *testing.T) {
		_, err := errorArray.Encode([3]int{anyValue, anyValue, anyValue}, []byte{})
		require.Equal(t, anyErr, err)
	})

	t.Run("Encode returns an error if the value is not a Array or array", func(t *testing.T) {
		_, err := testArray.Encode(anyValue, []byte{})
		require.True(t, errors.Is(err, types.ErrNotASlice))
	})

	t.Run("Decode returns an error if the underlying type returns an error", func(t *testing.T) {
		_, _, err := errorArray.Decode([]byte{0x03, 0x03, 0x04, 0x03, 0x04, 0x03, 0x04})
		require.Equal(t, anyErr, err)
	})

	t.Run("Decode returns an error if there are not enough bytes for the size", func(t *testing.T) {
		_, _, err := testArray.Decode([]byte{1})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("Size fixed size", func(t *testing.T) {
		size, err := testArray.Size(100)
		require.NoError(t, err)
		assert.Equal(t, numElements*len(elementCodec.Bytes), size)
	})

	t.Run("Size returns errors from element codec", func(t *testing.T) {
		_, err := errorArray.Size(100)
		require.Equal(t, anyErr, err)
	})

	t.Run("FixedSize returns num elements * size of element", func(t *testing.T) {
		size, err := testArray.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, numElements*len(elementCodec.Bytes), size)
	})
}
