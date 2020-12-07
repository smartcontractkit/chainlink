package store_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const correctPassphrase = "p@ssword"

func TestCreateEthereumAccount(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	require.NoError(t, store.KeyStore.Unlock(cltest.Password))

	_, err := store.KeyStore.NewAccount()
	assert.NoError(t, err)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Len(t, files, 2)
}

func TestUnlockKey_SingleAddress(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Verify the fixture account
	require.True(t, store.KeyStore.HasAccounts())
	require.Len(t, store.KeyStore.GetAccounts(), 1)

	assert.EqualError(t, store.KeyStore.Unlock("wrong phrase"), fmt.Sprintf("invalid password for account %s; could not decrypt key with given password", cltest.DefaultKey))
	assert.NoError(t, store.KeyStore.Unlock(cltest.Password))
}

func TestKeyStore_GetAccountByAddress(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	address := cltest.DefaultKeyAddress
	account, err := store.KeyStore.GetAccountByAddress(address)
	require.NoError(t, err)
	require.Equal(t, address, account.Address)

	missingAddress := cltest.NewAddress()
	account, err = store.KeyStore.GetAccountByAddress(missingAddress)
	require.EqualError(t, err, "no account found with that address")
}
