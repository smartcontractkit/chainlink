package keystore_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/caigo"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"

	starktxm "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_StarkNetKeyStore_E2E(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	keyStore := keystore.ExposedNewMaster(t, db, cfg.Database())
	require.NoError(t, keyStore.Unlock(cltest.Password))
	ks := keyStore.StarkNet()
	reset := func() {
		require.NoError(t, utils.JustError(db.Exec("DELETE FROM encrypted_key_rings")))
		keyStore.ResetXXXTestOnly()
		require.NoError(t, keyStore.Unlock(cltest.Password))
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
		newKey, err := starkkey.New()
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
		err := ks.EnsureKey()
		assert.NoError(t, err)

		err = ks.EnsureKey()
		assert.NoError(t, err)

		keys, err := ks.GetAll()
		require.NoError(t, err)
		require.Equal(t, 1, len(keys))
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
		signed, err := lk.Sign(context.Background(), starknetSenderAddr, nil)
		require.Nil(t, signed)
		require.NoError(t, err)
	})
	t.Run("key doesn't exists", func(t *testing.T) {
		baseKs.On("Get", mock.Anything).Return(starkkey.Key{}, fmt.Errorf("key doesn't exist"))
		signed, err := lk.Sign(context.Background(), "not an address", nil)
		require.Nil(t, signed)
		require.Error(t, err)
	})

	// TODO BCF-2242 remove this test once we have starknet smoke/integration tests
	// that exercise the transaction signing.
	t.Run("keystore adapter integration", func(t *testing.T) {

		adapter := starktxm.NewKeystoreAdapter(lk)
		baseKs.On("Get", starknetSenderAddr).Return(starkKey, nil)
		hash, err := caigo.Curve.PedersenHash([]*big.Int{big.NewInt(42)})
		require.NoError(t, err)
		r, s, err := adapter.Sign(context.Background(), starknetSenderAddr, hash)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.NotNil(t, s)

		pubx, puby, err := caigo.Curve.PrivateToPoint(starkKey.ToPrivKey())
		require.NoError(t, err)
		require.True(t, caigo.Curve.Verify(hash, r, s, pubx, puby))
	})
}
