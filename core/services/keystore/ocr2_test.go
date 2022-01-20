package keystore_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/stretchr/testify/require"
)

func Test_OCR2KeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	keyStore.Unlock(cltest.Password)
	ks := keyStore.OCR2()
	reset := func() {
		_, err := db.Exec("DELETE FROM encrypted_key_rings")
		require.NoError(t, err)
		keyStore.ResetXXXTestOnly()
		err = keyStore.Unlock(cltest.Password)
		require.NoError(t, err)
	}

	t.Run("initializes with an empty state", func(t *testing.T) {
		defer reset()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
	})

	t.Run("errors when getting non-existant ID", func(t *testing.T) {
		defer reset()
		_, err := ks.Get("non-existant-id")
		require.Error(t, err)
	})

	t.Run("creates a key with valid type", func(t *testing.T) {
		defer reset()
		key, err := ks.Create("evm")
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
		key, err = ks.Create("solana")
		require.NoError(t, err)
		retrievedKey, err = ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("gets keys by type", func(t *testing.T) {
		defer reset()
		keys, err := ks.GetAllOfType("evm")
		require.NoError(t, err)
		require.Len(t, keys, 0)
		keys, err = ks.GetAllOfType("solana")
		require.NoError(t, err)
		require.Len(t, keys, 0)
		_, err = ks.Create("evm")
		require.NoError(t, err)
		keys, err = ks.GetAllOfType("evm")
		require.NoError(t, err)
		require.Len(t, keys, 1)
		keys, err = ks.GetAllOfType("solana")
		require.NoError(t, err)
		require.Len(t, keys, 0)
		_, err = ks.Create("solana")
		require.NoError(t, err)
		keys, err = ks.GetAllOfType("evm")
		require.NoError(t, err)
		require.Len(t, keys, 1)
		keys, err = ks.GetAllOfType("solana")
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("errors when creating a key with an invalid type", func(t *testing.T) {
		defer reset()
		_, err := ks.Create("foobar")
		require.Error(t, err)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create("evm")
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		err = ks.Delete(key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(exportJSON, cltest.Password)
		require.NoError(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
		require.Equal(t, importedKey.ChainType(), retrievedKey.ChainType())
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		newKey, err := ocr2key.New("evm")
		require.NoError(t, err)
		err = ks.Add(newKey)
		require.NoError(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
		err = ks.Delete(newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		_, didExist, err := ks.EnsureKeys()
		require.NoError(t, err)
		require.False(t, didExist["evm"])
		require.False(t, didExist["solana"])
		require.False(t, didExist["terra"])
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 3, len(keys))
		_, didExist, err = ks.EnsureKeys()
		require.NoError(t, err)
		require.True(t, didExist["evm"])
		require.True(t, didExist["solana"])
		require.True(t, didExist["terra"])
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 3, len(keys))
	})
}
