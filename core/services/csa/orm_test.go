package csa_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/csa"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ORM_CreateCSAKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	var (
		passphrase = "passphrase"
		orm        = csa.NewORM(store.DB)
	)

	key, err := csa.NewCSAKey(passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	count, err := orm.CountCSAKeys()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateCSAKey(context.Background(), key)
	require.NoError(t, err)

	count, err = orm.CountCSAKeys()
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)
}

func Test_ORM_ListCSAKeys(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	var (
		passphrase = "passphrase"
		orm        = csa.NewORM(store.DB)
	)

	key, err := csa.NewCSAKey(passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	id, err := orm.CreateCSAKey(context.Background(), key)
	require.NoError(t, err)

	mgrs, err := orm.ListCSAKeys(context.Background())
	require.NoError(t, err)
	require.Len(t, mgrs, 1)

	actual := mgrs[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, key.PublicKey, actual.PublicKey)
	expectedPrivKey, err := key.EncryptedPrivateKey.Decrypt("passphrase")
	require.NoError(t, err)
	actualPrivKey, err := actual.EncryptedPrivateKey.Decrypt("passphrase")
	require.NoError(t, err)
	assert.Equal(t, expectedPrivKey, actualPrivKey)
}

func Test_ORM_GetCSAKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	var (
		passphrase = "passphrase"
		orm        = csa.NewORM(store.DB)
	)

	key, err := csa.NewCSAKey(passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	id, err := orm.CreateCSAKey(context.Background(), key)
	require.NoError(t, err)

	actual, err := orm.GetCSAKey(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, key.PublicKey, actual.PublicKey)
	expectedPrivKey, err := key.EncryptedPrivateKey.Decrypt("passphrase")
	require.NoError(t, err)
	actualPrivKey, err := actual.EncryptedPrivateKey.Decrypt("passphrase")
	require.NoError(t, err)
	assert.Equal(t, expectedPrivKey, actualPrivKey)
}
