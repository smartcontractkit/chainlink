package keystore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
)

func Test_CSAKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.CSA()
	reset := func() {
		ctx := context.Background() // Executed on cleanup
		_, err := db.Exec("DELETE FROM encrypted_key_rings")
		require.NoError(t, err)
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

		t.Run("prevents creating more than one key", func(t *testing.T) {
			ctx := testutils.Context(t)
			k, err2 := ks.Create(ctx)

			assert.Zero(t, k)
			assert.Error(t, err2)
			assert.True(t, errors.Is(err2, keystore.ErrCSAKeyExists))
		})
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		key, err := ks.Create(ctx)
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Delete(ctx, key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(ctx, exportJSON, cltest.Password)
		require.NoError(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)

		t.Run("prevents importing more than one key", func(t *testing.T) {
			k, err2 := ks.Import(testutils.Context(t), exportJSON, cltest.Password)

			assert.Zero(t, k)
			assert.Error(t, err2)
			assert.Equal(t, fmt.Sprintf("key with ID %s already exists", key.ID()), err2.Error())
		})

		t.Run("fails to import malformed key", func(t *testing.T) {
			k, err2 := ks.Import(testutils.Context(t), []byte(""), cltest.Password)

			assert.Zero(t, k)
			assert.Error(t, err2)
		})

		t.Run("fails to export non-existent key", func(t *testing.T) {
			exportJSON, err = ks.Export("non-existent", cltest.Password)

			assert.Error(t, err)
			assert.Empty(t, exportJSON)
		})
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		newKey, err := csakey.NewV2()
		require.NoError(t, err)
		err = ks.Add(ctx, newKey)
		require.NoError(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
		_, err = ks.Delete(ctx, newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)

		t.Run("prevents adding more than one key", func(t *testing.T) {
			ctx := testutils.Context(t)
			err = ks.Add(ctx, newKey)
			require.NoError(t, err)

			err = ks.Add(ctx, newKey)

			assert.Error(t, err)
			assert.True(t, errors.Is(err, keystore.ErrCSAKeyExists))
		})

		t.Run("fails to delete non-existent key", func(t *testing.T) {
			k, err2 := ks.Delete(testutils.Context(t), "non-existent")

			assert.Zero(t, k)
			assert.Error(t, err2)
		})
	})

	t.Run("adds an externally created key/ensures it already exists", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		newKey, err := csakey.NewV2()
		assert.NoError(t, err)
		err = ks.Add(ctx, newKey)
		assert.NoError(t, err)

		err = keyStore.CSA().EnsureKey(ctx)
		assert.NoError(t, err)
		keys, err2 := ks.GetAll()
		assert.NoError(t, err2)

		require.Equal(t, 1, len(keys))
		require.Equal(t, newKey.ID(), keys[0].ID())
		require.Equal(t, newKey.Version, keys[0].Version)
		require.Equal(t, newKey.PublicKey, keys[0].PublicKey)
	})

	t.Run("auto creates a key if it doesn't exists when trying to ensure it already exists", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		keys, err := ks.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(keys))

		err = keyStore.CSA().EnsureKey(ctx)
		assert.NoError(t, err)

		keys, err = ks.GetAll()
		assert.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
	})
}
