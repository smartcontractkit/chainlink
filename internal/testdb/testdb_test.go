package testdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrator_Hash(t *testing.T) {
	t.Parallel()

	t.Run("empty returns empty", func(t *testing.T) {
		m := migratorConfig(false)
		hash, err := m.Hash()
		require.NoError(t, err)
		require.Equal(t, "empty", hash)
	})

	t.Run("withTemplate hashes successfully", func(t *testing.T) {
		m := migratorConfig(true)
		hash, err := m.Hash()
		require.NoError(t, err)
		require.NotEmpty(t, hash)
		require.NotEqual(t, "empty", hash)
	})
}
