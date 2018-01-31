package store_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
)

const passphrase = "p@ssword"

func TestCreateEthereumAccount(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	_, err := store.KeyStore.NewAccount(passphrase)
	assert.Nil(t, err)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Equal(t, 1, len(files))
}

func TestUnlockKey(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	store.KeyStore.NewAccount(passphrase)

	assert.NotNil(t, store.KeyStore.Unlock("wrong phrase"))
	assert.Nil(t, store.KeyStore.Unlock(passphrase))
}

func TestShowEthBalance(t *testing.T) {
	t.Parallel()
	cl_store, cleanup := cltest.NewStore()
	defer cleanup()
	ans := cl_store.KeyStore.ShowEthBalance(cl_store.TxManager)
	assert.Equal(t, ans, store.MissingWalletString)
}
