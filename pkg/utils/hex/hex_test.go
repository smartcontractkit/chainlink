package hex_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
)

func TestParseBig(t *testing.T) {
	t.Parallel()

	t.Run("parses successfully", func(t *testing.T) {
		t.Parallel()

		_, err := hex.ParseBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
		assert.NoError(t, err)
	})

	t.Run("returns error", func(t *testing.T) {
		t.Parallel()

		_, err := hex.ParseBig("0xabc")
		assert.Error(t, err)
	})
}

func TestTrimPrefix(t *testing.T) {
	t.Parallel()

	t.Run("trims prefix", func(t *testing.T) {
		t.Parallel()

		s := hex.TrimPrefix("0xabc")
		assert.Equal(t, "abc", s)
	})

	t.Run("returns the same string if it doesn't have prefix", func(t *testing.T) {
		t.Parallel()

		s := hex.TrimPrefix("defg")
		assert.Equal(t, "defg", s)
	})
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	t.Run("has prefix", func(t *testing.T) {
		t.Parallel()

		r := hex.HasPrefix("0xabc")
		assert.True(t, r)
	})

	t.Run("doesn't have prefix", func(t *testing.T) {
		t.Parallel()

		r := hex.HasPrefix("abc")
		assert.False(t, r)
	})

	t.Run("has 0x suffix", func(t *testing.T) {
		t.Parallel()

		r := hex.HasPrefix("abc0x")
		assert.False(t, r)
	})
}

func TestDecodeString(t *testing.T) {
	t.Parallel()

	t.Run("0x prefix missing", func(t *testing.T) {
		t.Parallel()

		_, err := hex.DecodeString("abcd")
		assert.Error(t, err)
	})

	t.Run("wrong hex characters", func(t *testing.T) {
		t.Parallel()

		_, err := hex.DecodeString("0xabcdzzz")
		assert.Error(t, err)
	})

	t.Run("valid hex string", func(t *testing.T) {
		t.Parallel()

		b, err := hex.DecodeString("0x1234")
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x12, 0x34}, b)
	})

	t.Run("prepend odd length with zero", func(t *testing.T) {
		t.Parallel()

		b, err := hex.DecodeString("0x123")
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x1, 0x23}, b)
	})
}
