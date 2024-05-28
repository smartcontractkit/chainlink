package keystore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

func Test_CosmosKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.Cosmos()
	reset := func() {
		ctx := context.Background() // Executed during cleanup
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(ctx, cltest.Password))
	}

	t.Run("initializes with an empty state", func(t *testing.T) {
		defer reset()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
	})

	t.Run("errors when getting non-existent ID", func(t *testing.T) {
		defer reset()
		_, err := ks.Get("non-existent-id")
		require.Error(t, err)
	})

	t.Run("creates a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		key, err := ks.Create(ctx)
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		key, err := ks.Create(ctx)
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Export("non-existent", cltest.Password)
		assert.Error(t, err)
		_, err = ks.Delete(ctx, key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(ctx, exportJSON, cltest.Password)
		require.NoError(t, err)
		_, err = ks.Import(ctx, exportJSON, cltest.Password)
		assert.Error(t, err)
		_, err = ks.Import(ctx, []byte(""), cltest.Password)
		assert.Error(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		newKey := cosmoskey.New()
		err := ks.Add(ctx, newKey)
		require.NoError(t, err)
		err = ks.Add(ctx, newKey)
		assert.Error(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
		_, err = ks.Delete(ctx, newKey.ID())
		require.NoError(t, err)
		_, err = ks.Delete(ctx, newKey.ID())
		assert.Error(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKey(ctx)
		assert.NoError(t, err)

		err = ks.EnsureKey(ctx)
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
	})
}
