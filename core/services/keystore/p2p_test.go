package keystore_test

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/stretchr/testify/require"
)

func Test_P2PKeyStore_E2E(t *testing.T) {
	db := pgtest.NewGormDB(t)
	keyStore := keystore.ExposedNewMaster(db)
	keyStore.Unlock(cltest.Password)
	ks := keyStore.P2P()
	reset := func() {
		require.NoError(t, db.Exec("DELETE FROM encrypted_key_rings").Error)
		keyStore.ResetXXXTestOnly()
		keyStore.Unlock(cltest.Password)
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

	t.Run("creates a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create()
		require.NoError(t, err)
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, key, retrievedKey)
	})

	t.Run("imports and exports a key", func(t *testing.T) {
		defer reset()
		key, err := ks.Create()
		require.NoError(t, err)
		exportJSON, err := ks.Export(key.ID(), cltest.Password)
		require.NoError(t, err)
		_, err = ks.Delete(key.ID())
		require.NoError(t, err)
		_, err = ks.Get(key.ID())
		require.Error(t, err)
		importedKey, err := ks.Import(exportJSON, cltest.Password)
		require.NoError(t, err)
		require.Equal(t, key.ID(), importedKey.ID())
		retrievedKey, err := ks.Get(key.ID())
		require.NoError(t, err)
		require.Equal(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		newKey, err := p2pkey.NewV2()
		require.NoError(t, err)
		err = ks.Add(newKey)
		require.NoError(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
		_, err = ks.Delete(newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(keys))
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)
	})

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		_, didExist, err := ks.EnsureKey()
		require.NoError(t, err)
		require.False(t, didExist)
		_, didExist, err = ks.EnsureKey()
		require.NoError(t, err)
		require.True(t, didExist)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
	})
}

func Test_P2PKeyStore_GetOrFirst(t *testing.T) {
	db := pgtest.NewGormDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	cfg := configtest.NewTestGeneralConfig(t)
	ks := keyStore.P2P()
	_, err := ks.GetOrFirst("")
	require.Contains(t, err.Error(), "no p2p keys exist")
	id := cfg.P2PPeerID().Raw()
	_, err = ks.GetOrFirst(id)
	require.Contains(t, err.Error(), fmt.Sprintf("unable to find P2P key with id %s", id))
	k1, err := keyStore.P2P().Create()
	require.NoError(t, err)
	k2, err := ks.GetOrFirst("")
	require.NoError(t, err)
	require.Equal(t, k1, k2)
	k3, err := ks.GetOrFirst(k1.ID())
	require.NoError(t, err)
	require.Equal(t, k1, k3)
	_, err = keyStore.P2P().Create()
	require.NoError(t, err)
	_, err = ks.GetOrFirst("")
	require.Contains(t, err.Error(), "multiple p2p keys found")
	k4, err := ks.GetOrFirst(k1.ID())
	require.NoError(t, err)
	require.Equal(t, k1, k4)
}
