package keystore_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"

	"github.com/stretchr/testify/require"
)

func Test_VRFKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.VRF()
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
		require.Empty(t, keys)
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
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		newKey, err := vrfkey.NewV2()
		require.NoError(t, err)
		err = ks.Add(ctx, newKey)
		require.NoError(t, err)
		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
		_, err = ks.Delete(ctx, newKey.ID())
		require.NoError(t, err)
		keys, err = ks.GetAll()
		require.NoError(t, err)
		require.Empty(t, keys)
		_, err = ks.Get(newKey.ID())
		require.Error(t, err)
	})

	t.Run("fails to add an already added key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		k, err := vrfkey.NewV2()
		require.NoError(t, err)

		err = ks.Add(ctx, k)
		require.NoError(t, err)
		err = ks.Add(ctx, k)

		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("key with ID %s already exists", k.ID()), err.Error())
	})

	t.Run("fails to delete a key that doesn't exists", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		k, err := vrfkey.NewV2()
		require.NoError(t, err)

		err = ks.Add(ctx, k)
		require.NoError(t, err)

		fk, err := ks.Delete(ctx, "non-existent")

		assert.Zero(t, fk)
		assert.Error(t, err)
	})

	t.Run("imports a key exported from a v1 keystore", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		exportedKey := `{"PublicKey":"0xd2377bc6be8a2c5ce163e1867ee42ef111e320686f940a98e52e9c019ca0606800","vrf_key":{"address":"b94276ad4e5452732ec0cccf30ef7919b67844b6","crypto":{"cipher":"aes-128-ctr","ciphertext":"ff66d61d02dba54a61bab1ceb8414643f9e76b7351785d2959e2c8b50ee69a92","cipherparams":{"iv":"75705da271b11e330a27b8d593a3930c"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"efe5b372e4fe79d0af576a79d65a1ee35d0792d9c92b70107b5ada1817ea7c7b"},"mac":"e4d0bb08ffd004ab03aeaa42367acbd9bb814c6cfd981f5157503f54c30816e7"},"version":3}}`
		importedKey, err := ks.Import(ctx, []byte(exportedKey), "p4SsW0rD1!@#_")
		require.NoError(t, err)
		require.Equal(t, "0xd2377bc6be8a2c5ce163e1867ee42ef111e320686f940a98e52e9c019ca0606800", importedKey.ID())
	})

	t.Run("fails to import an already imported key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)

		exportedKey := `{"PublicKey":"0xd2377bc6be8a2c5ce163e1867ee42ef111e320686f940a98e52e9c019ca0606800","vrf_key":{"address":"b94276ad4e5452732ec0cccf30ef7919b67844b6","crypto":{"cipher":"aes-128-ctr","ciphertext":"ff66d61d02dba54a61bab1ceb8414643f9e76b7351785d2959e2c8b50ee69a92","cipherparams":{"iv":"75705da271b11e330a27b8d593a3930c"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"efe5b372e4fe79d0af576a79d65a1ee35d0792d9c92b70107b5ada1817ea7c7b"},"mac":"e4d0bb08ffd004ab03aeaa42367acbd9bb814c6cfd981f5157503f54c30816e7"},"version":3}}`
		importedKey, err := ks.Import(ctx, []byte(exportedKey), "p4SsW0rD1!@#_")
		require.NoError(t, err)

		keyStore.SetPassword("p4SsW0rD1!@#_")
		k, err := ks.Import(ctx, []byte(exportedKey), "p4SsW0rD1!@#_")

		assert.Zero(t, k)
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("key with ID %s already exists", importedKey.ID()), err.Error())
	})

	t.Run("fails to export non-existent key", func(t *testing.T) {
		k, err := ks.Export("non-existent", cltest.Password)

		assert.Zero(t, k)
		assert.Error(t, err)
	})

	t.Run("generate proof for keys", func(t *testing.T) {
		defer reset()

		t.Run("fails to generate proof for non-existent key", func(t *testing.T) {
			pf, err := ks.GenerateProof("non-existent", big.NewInt(int64(1)))

			assert.Zero(t, pf)
			assert.Error(t, err)
		})

		t.Run("generates a proof for a key", func(t *testing.T) {
			ctx := testutils.Context(t)
			k, err := vrfkey.NewV2()
			require.NoError(t, err)
			err = ks.Add(ctx, k)
			require.NoError(t, err)

			pf, err := ks.GenerateProof(k.ID(), big.NewInt(int64(1)))
			require.NoError(t, err)

			assert.NotZero(t, pf)
		})
	})
}
