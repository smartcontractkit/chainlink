package binary_test

import (
	"errors"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/smartcontractkit/libocr/bigbigendian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestBigInteger(t *testing.T) {
	t.Run("NewBigInt returns error if numBytes is too large", func(t *testing.T) {
		_, err := bi.BigInt(bigbigendian.MaxSize+1, false)
		assert.True(t, errors.Is(err, types.ErrInvalidConfig))
	})

	signed, createErr := bi.BigInt(4, true)
	require.NoError(t, createErr)
	unsigned, createErr := bi.BigInt(4, false)
	require.NoError(t, createErr)

	t.Run("Encode and Decode work together on signed values", func(t *testing.T) {
		anyVal := big.NewInt(-100)
		bytes, err := signed.Encode(anyVal, nil)
		require.NoError(t, err)
		expected, err := bi.Int32().Encode(int32(anyVal.Int64()), nil)
		require.NoError(t, err)
		assert.Equal(t, expected, bytes)

		decoded, remaining, err := signed.Decode(bytes)
		require.NoError(t, err)
		assert.Equal(t, 0, decoded.(*big.Int).Cmp(anyVal))
		assert.Empty(t, remaining)
	})

	t.Run("Encode and Decode work together on unsigned values", func(t *testing.T) {
		anyVal := big.NewInt(100)
		bytes, err := unsigned.Encode(anyVal, nil)
		require.NoError(t, err)
		expected, err := bi.Uint32().Encode(uint32(anyVal.Uint64()), nil)
		require.NoError(t, err)
		assert.Equal(t, expected, bytes)

		decoded, remaining, err := signed.Decode(bytes)
		require.NoError(t, err)
		assert.Equal(t, 0, decoded.(*big.Int).Cmp(anyVal))
		assert.Empty(t, remaining)
	})

	t.Run("Encoding out of range values return an error", func(t *testing.T) {
		bi := big.NewInt(math.MaxInt32 + 1)
		_, err := signed.Encode(bi, nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
		bi = big.NewInt(math.MinInt32 - 1)
		_, err = signed.Encode(bi, nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))

		bi = big.NewInt(math.MaxUint32 + 1)
		_, err = unsigned.Encode(bi, nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
		bi = big.NewInt(-1)
		_, err = unsigned.Encode(bi, nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Encode returns an error if input is not a *big.Int", func(t *testing.T) {
		_, err := signed.Encode(100, nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("GetType returns *big.Int", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(&big.Int{}), signed.GetType())
	})

	t.Run("Size returns the number of bytes", func(t *testing.T) {
		size, err := signed.Size(100)
		require.NoError(t, err)
		assert.Equal(t, 4, size)
	})

	t.Run("FixedSize returns the number of bytes", func(t *testing.T) {
		size, err := signed.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, 4, size)
	})

	tests := []struct {
		name  string
		codec encodings.TypeCodec
	}{{"signed", signed}, {"unsigned", unsigned}}

	t.Run("Encode appends to prefix", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				anyVal := big.NewInt(100)
				bytes, err := test.codec.Encode(anyVal, []byte{1, 2, 3})
				require.NoError(t, err)
				expected, err := bi.Int32().Encode(int32(anyVal.Int64()), nil)
				require.NoError(t, err)
				assert.Equal(t, append([]byte{1, 2, 3}, expected...), bytes)
			})
		}
	})

	t.Run("Decode leaves suffix", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				anyVal := big.NewInt(100)
				bytes, err := test.codec.Encode(anyVal, nil)
				require.NoError(t, err)

				val, remaining, err := test.codec.Decode(append(bytes, 1, 2, 3))
				require.NoError(t, err)
				assert.Equal(t, anyVal, val)
				assert.Equal(t, []byte{1, 2, 3}, remaining)
			})
		}
	})

	t.Run("Decode returns an error when there are not enough bytes", func(t *testing.T) {
		_, _, err := signed.Decode([]byte{1, 2})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})
}
