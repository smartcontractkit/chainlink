package keystore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

func Test_StarkNetKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	keyStore := keystore.ExposedNewMaster(t, db)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password))
	ks := keyStore.StarkNet()
	reset := func() {
		ctx := context.Background() // Executed on cleanup
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
		requireEqualKeys(t, key, retrievedKey)
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
		requireEqualKeys(t, importedKey, retrievedKey)
	})

	t.Run("adds an externally created key / deletes a key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		newKey, err := starkkey.New()
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

	t.Run("ensures key", func(t *testing.T) {
		defer reset()
		ctx := testutils.Context(t)
		err := ks.EnsureKey(ctx)
		assert.NoError(t, err)

		err = ks.EnsureKey(ctx)
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})
}

func TestStarknetSigner(t *testing.T) {
	var (
		starknetSenderAddr = "legit"
	)
	baseKs := mocks.NewStarkNet(t)
	starkKey, err := starkkey.New()
	require.NoError(t, err)

	lk := &keystore.StarknetLooppSigner{baseKs}
	// test that we implementw the loopp spec. signing nil data should not error
	// on existing sender id
	t.Run("key exists", func(t *testing.T) {
		baseKs.On("Get", starknetSenderAddr).Return(starkKey, nil)
		signed, err := lk.Sign(testutils.Context(t), starknetSenderAddr, nil)
		require.Nil(t, signed)
		require.NoError(t, err)
	})
	t.Run("key doesn't exists", func(t *testing.T) {
		baseKs.On("Get", mock.Anything).Return(starkkey.Key{}, errors.New("key doesn't exist"))
		signed, err := lk.Sign(testutils.Context(t), "not an address", nil)
		require.Nil(t, signed)
		require.Error(t, err)
	})
}
