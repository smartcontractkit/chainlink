package keystore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

func Test_OCR2KeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.OCR2()
	reset := func() {
		ctx := context.Background() // Executed on cleanup
		_, err := db.Exec("DELETE FROM encrypted_key_rings")
		require.NoError(t, err)
		keyStore.ResetXXXTestOnly()
		err = keyStore.Unlock(ctx, cltest.Password)
		require.NoError(t, err)
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

	t.Run("creates a key with valid type", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		// lopp through different chain types
		for _, chain := range chaintype.SupportedChainTypes {
			key, err := ks.Create(ctx, chain)
			require.NoError(t, err)
			retrievedKey, err := ks.Get(key.ID())
			require.NoError(t, err)
			require.Equal(t, key, retrievedKey)
		}
	})

	t.Run("gets keys by type", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		created := map[chaintype.ChainType]bool{}
		for _, chain := range chaintype.SupportedChainTypes {
			// validate no keys exist for chain
			keys, err := ks.GetAllOfType(chain)
			require.NoError(t, err)
			require.Len(t, keys, 0)

			_, err = ks.Create(ctx, chain)
			require.NoError(t, err)
			created[chain] = true

			// validate that only 1 of each exists after creation
			for _, c := range chaintype.SupportedChainTypes {
				keys, err := ks.GetAllOfType(c)
				require.NoError(t, err)
				if created[c] {
					require.Len(t, keys, 1)
					continue
				}
				require.Len(t, keys, 0)
			}
		}
	})

	t.Run("errors when creating a key with an invalid type", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		_, err := ks.Create(ctx, "foobar")
		require.Error(t, err)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		for _, chain := range chaintype.SupportedChainTypes {
			key, err := ks.Create(ctx, chain)
			require.NoError(t, err)
			exportJSON, err := ks.Export(key.ID(), cltest.Password)
			require.NoError(t, err)
			_, err = ks.Export("non-existent", cltest.Password)
			assert.Error(t, err)
			err = ks.Delete(ctx, key.ID())
			require.NoError(t, err)
			_, err = ks.Get(key.ID())
			require.Error(t, err)
			importedKey, err := ks.Import(ctx, exportJSON, cltest.Password)
			require.NoError(t, err)
			_, err = ks.Import(ctx, []byte(""), cltest.Password)
			assert.Error(t, err)
			require.Equal(t, key.ID(), importedKey.ID())
			retrievedKey, err := ks.Get(key.ID())
			require.NoError(t, err)
			require.Equal(t, importedKey, retrievedKey)
			require.Equal(t, importedKey.ChainType(), retrievedKey.ChainType())
		}
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		for _, chain := range chaintype.SupportedChainTypes {
			newKey, err := ocr2key.New(chain)
			require.NoError(t, err)
			err = ks.Add(ctx, newKey)
			require.NoError(t, err)
			err = ks.Add(ctx, newKey)
			assert.Error(t, err)
			keys, err := ks.GetAll()
			require.NoError(t, err)
			require.Equal(t, 1, len(keys))
			err = ks.Delete(ctx, newKey.ID())
			require.NoError(t, err)
			err = ks.Delete(ctx, newKey.ID())
			assert.Error(t, err)
			keys, err = ks.GetAll()
			require.NoError(t, err)
			require.Equal(t, 0, len(keys))
			_, err = ks.Get(newKey.ID())
			require.Error(t, err)
		}
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKeys(ctx, chaintype.SupportedChainTypes...)
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		assert.NoError(t, err)
		require.Equal(t, len(chaintype.SupportedChainTypes), len(keys))

		err = ks.EnsureKeys(ctx, chaintype.SupportedChainTypes...)
		assert.NoError(t, err)

		// loop through different supported chain types
		for _, chain := range chaintype.SupportedChainTypes {
			keys, err := ks.GetAllOfType(chain)
			assert.NoError(t, err)
			require.Equal(t, 1, len(keys))
		}
	})

	t.Run("ensures key only for enabled chains", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKeys(ctx, chaintype.EVM)
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		assert.NoError(t, err)
		require.Equal(t, 1, len(keys))
		require.Equal(t, keys[0].ChainType(), chaintype.EVM)

		err = ks.EnsureKeys(ctx, chaintype.Cosmos)
		assert.NoError(t, err)

		keys, err = ks.GetAll()
		assert.NoError(t, err)
		require.Equal(t, 2, len(keys))

		cosmosKeys, err := ks.GetAllOfType(chaintype.Cosmos)
		assert.NoError(t, err)
		require.Equal(t, 1, len(cosmosKeys))
		require.Equal(t, cosmosKeys[0].ChainType(), chaintype.Cosmos)

		err = ks.EnsureKeys(ctx, chaintype.StarkNet)
		assert.NoError(t, err)

		keys, err = ks.GetAll()
		assert.NoError(t, err)
		require.Equal(t, 3, len(keys))

		straknetKeys, err := ks.GetAllOfType(chaintype.StarkNet)
		assert.NoError(t, err)
		require.Equal(t, 1, len(straknetKeys))
		require.Equal(t, straknetKeys[0].ChainType(), chaintype.StarkNet)
	})
}
