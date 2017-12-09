package services_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/stretchr/testify/assert"
)

const passphrase = "p@ssword"

func TestCreateEthereumAccount(t *testing.T) {
	t.Parallel()
	store := cltest.Store()
	defer store.Close()

	_, err := store.KeyStore.NewAccount(passphrase)
	assert.Nil(t, err)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Equal(t, 1, len(files))
}

func TestUnlockKey(t *testing.T) {
	t.Parallel()
	store := cltest.Store()
	defer store.Close()

	account, err := store.KeyStore.NewAccount(passphrase)
	assert.Nil(t, err)

	assert.NotNil(t, store.KeyStore.Unlock(account, "wrong phrase"))
	err = store.KeyStore.Unlock(account, passphrase)
	assert.Nil(t, err)
}
