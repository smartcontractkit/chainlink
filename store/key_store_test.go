package store_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

const passphrase = "p@ssword"

func TestCreateEthereumAccount(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	_, err := store.KeyStore.NewAccount(passphrase)
	assert.NoError(t, err)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Equal(t, 1, len(files))
}

func TestUnlockKey(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	store.KeyStore.NewAccount(passphrase)

	assert.Error(t, store.KeyStore.Unlock("wrong phrase"))
	assert.NoError(t, store.KeyStore.Unlock(passphrase))
}

func TestKeyStore_SignSuccess(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	_, err := store.KeyStore.NewAccount(passphrase)
	assert.NoError(t, err)
	assert.NoError(t, store.KeyStore.Unlock(passphrase))

	signature, err := store.KeyStore.Sign([]byte("abc123"))
	assert.NoError(t, err)
	assert.NotEqual(t, "", signature)
}

func TestKeyStore_SignAccountLocked(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	_, err := store.KeyStore.NewAccount(passphrase)
	assert.NoError(t, err)

	_, err = store.KeyStore.Sign([]byte("abc123"))
	assert.Error(t, err)
}
