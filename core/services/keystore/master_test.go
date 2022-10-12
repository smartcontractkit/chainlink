package keystore_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
)

func TestMasterKeystore_Unlock_Save(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	keyStore := keystore.ExposedNewMaster(t, db, cfg)
	const tableName = "encrypted_key_rings"
	reset := func() {
		keyStore.ResetXXXTestOnly()
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
		require.NoError(t, err)
	}

	t.Run("can be unlocked more than once, as long as the passwords match", func(t *testing.T) {
		defer reset()
		require.NoError(t, keyStore.Unlock(cltest.Password))
		require.NoError(t, keyStore.Unlock(cltest.Password))
		require.NoError(t, keyStore.Unlock(cltest.Password))
		require.Error(t, keyStore.Unlock("wrong password"))
	})

	t.Run("saves an empty keyRing", func(t *testing.T) {
		defer reset()
		require.NoError(t, keyStore.Unlock(cltest.Password))
		cltest.AssertCount(t, db, tableName, 1)
		require.NoError(t, keyStore.ExportedSave())
		cltest.AssertCount(t, db, tableName, 1)
	})

	t.Run("won't load a saved keyRing if the password is incorrect", func(t *testing.T) {
		defer reset()
		require.NoError(t, keyStore.Unlock(cltest.Password))
		cltest.MustAddRandomKeyToKeystore(t, keyStore.Eth()) // need at least 1 key to encrypt
		cltest.AssertCount(t, db, tableName, 1)
		keyStore.ResetXXXTestOnly()
		cltest.AssertCount(t, db, tableName, 1)
		require.Error(t, keyStore.Unlock("password2"))
		cltest.AssertCount(t, db, tableName, 1)
	})

	t.Run("loads a saved keyRing if the password is correct", func(t *testing.T) {
		defer reset()
		require.NoError(t, keyStore.Unlock(cltest.Password))
		require.NoError(t, keyStore.ExportedSave())
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(cltest.Password))
	})
}
