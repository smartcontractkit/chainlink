package keystore_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/stretchr/testify/require"
)

func TestMasterKeystore_Unlock_Save(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)

	keyStore := keystore.ExposedNewMaster(t, db)
	reset := func() {
		keyStore.ResetXXXTestOnly()
		err := db.Exec("DELETE FROM encrypted_key_rings").Error
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
		cltest.AssertCount(t, db, keystore.ExportedEncryptedKeyRing{}, 1)
		require.NoError(t, keyStore.ExportedSave())
		cltest.AssertCount(t, db, keystore.ExportedEncryptedKeyRing{}, 1)
	})

	t.Run("won't load a saved keyRing if the password is incorrect", func(t *testing.T) {
		defer reset()
		require.NoError(t, keyStore.Unlock(cltest.Password))
		cltest.MustAddRandomKeyToKeystore(t, keyStore.Eth()) // need at least 1 key to encrypt
		cltest.AssertCount(t, db, keystore.ExportedEncryptedKeyRing{}, 1)
		keyStore.ResetXXXTestOnly()
		cltest.AssertCount(t, db, keystore.ExportedEncryptedKeyRing{}, 1)
		require.Error(t, keyStore.Unlock("password2"))
		cltest.AssertCount(t, db, keystore.ExportedEncryptedKeyRing{}, 1)
	})

	t.Run("loads a saved keyRing if the password is correct", func(t *testing.T) {
		defer reset()
		require.NoError(t, keyStore.Unlock(cltest.Password))
		require.NoError(t, keyStore.ExportedSave())
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(cltest.Password))
	})
}
