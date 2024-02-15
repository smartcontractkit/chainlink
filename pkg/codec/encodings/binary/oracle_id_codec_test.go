package binary_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/binary"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestOracleID(t *testing.T) {
	t.Parallel()
	t.Run("Encodes and decodes to the same value with correct encoding length", func(t *testing.T) {
		i := &binary.OracleID{}
		value := commontypes.OracleID(123)

		encoded, err := i.Encode(value, nil)

		require.NoError(t, err)
		assert.Equal(t, 1, len(encoded))

		decoded, remaining, err := i.Decode(encoded)

		require.NoError(t, err)
		assert.Equal(t, 0, len(remaining))
		assert.Equal(t, value, decoded)
	})

	t.Run("Encodes appends to prefix", func(t *testing.T) {
		i := &binary.OracleID{}
		value := commontypes.OracleID(123)
		prefix := []byte{1, 2, 3}

		encoded, err := i.Encode(value, prefix)

		require.NoError(t, err)
		assert.Equal(t, 4, len(encoded))
		expected, err := i.Encode(value, nil)
		require.NoError(t, err)
		assert.Equal(t, expected, encoded[3:])
	})

	t.Run("Decode leaves a suffix", func(t *testing.T) {
		i := &binary.OracleID{}
		value := commontypes.OracleID(123)
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
		i := &binary.OracleID{}
		_, _, err := i.Decode([]byte{})
		require.True(t, errors.Is(err, types.ErrInvalidEncoding))
	})

	t.Run("GetType returns correct type", func(t *testing.T) {
		i := &binary.OracleID{}
		assert.Equal(t, i.GetType(), reflect.TypeOf(commontypes.OracleID(0)))
	})

	t.Run("Size returns correct size", func(t *testing.T) {
		size, err := (&binary.OracleID{}).Size(100)
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("FixedSize returns correct size", func(t *testing.T) {
		size, err := (&binary.OracleID{}).FixedSize()
		require.NoError(t, err)
		assert.Equal(t, size, 1)
	})

	t.Run("returns an error if the input is not an OracleID", func(t *testing.T) {
		i := &binary.OracleID{}

		_, err := i.Encode("foo", nil)
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})
}
