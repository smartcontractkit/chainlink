package pgtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueName(t *testing.T) {
	t.Parallel()

	first := UniqueName(t, "cid-123")
	assert.Equal(t, first, UniqueName(t, "cid-123"), "same test and base must produce the same name")
	assert.NotEqual(t, first, UniqueName(t, "cid-456"), "different base must produce a different name")

	t.Run("subtest", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, first, UniqueName(t, "cid-123"), "subtests must not reuse the parent's name")
	})
}

func TestUniqueChainID(t *testing.T) {
	t.Parallel()

	first := UniqueChainID(t)
	second := UniqueChainID(t)
	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.NotEqual(t, 0, first.Sign(), "chain ID should be positive")
	assert.NotEqual(t, first, second, "each call should produce a new chain ID")
}
