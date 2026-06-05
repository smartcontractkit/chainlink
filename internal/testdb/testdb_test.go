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

	t.Run("withTemplate returns same hash from new migrator instances", func(t *testing.T) {
		hash1, err := migratorConfig(true).Hash()
		require.NoError(t, err)

		hash2, err := migratorConfig(true).Hash()
		require.NoError(t, err)
		require.Equal(t, hash1, hash2)
	})
}

func TestMigrator_HashCachedNoAllocs(t *testing.T) {
	m := migratorConfig(true)
	_, err := m.Hash()
	require.NoError(t, err)

	allocs := testing.AllocsPerRun(5, func() {
		_, err := m.Hash()
		require.NoError(t, err)
	})
	require.Zero(t, allocs)
}
