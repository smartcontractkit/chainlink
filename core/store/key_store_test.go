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

	_, err := store.KeyStore.NewAccount(correctPassphrase)
	assert.NoError(t, err)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Len(t, files, 1)
}

func TestUnlockKey_SingleAddress(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, address := cltest.MustAddRandomKeyToKeystore(t, store, 0)

	require.True(t, store.KeyStore.HasAccounts())
	require.Len(t, store.KeyStore.GetAccounts(), 1)

	assert.EqualError(t, store.KeyStore.Unlock("wrong phrase"), fmt.Sprintf("invalid password for account %s; could not decrypt key with given password", address.Hex()))
	assert.NoError(t, store.KeyStore.Unlock(cltest.Password))
}

func TestUnlockKey_MultipleAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		tryPassphrase string
		wantErr       bool
	}{
		{"correct", cltest.Password, false},
		{"wrong", "wrong", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			cltest.MustAddRandomKeyToKeystore(t, store)
			cltest.MustAddRandomKeyToKeystore(t, store)

			require.True(t, store.KeyStore.HasAccounts())
			require.Len(t, store.KeyStore.GetAccounts(), 2)
			defer cleanup()

			if test.wantErr {
				assert.Error(t, store.KeyStore.Unlock(test.tryPassphrase))
			} else {
				assert.NoError(t, store.KeyStore.Unlock(test.tryPassphrase))
			}
		})
	}
}

func TestKeyStore_GetAccountByAddress(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, address := cltest.MustAddRandomKeyToKeystore(t, store, 0)

	account, err := store.KeyStore.GetAccountByAddress(address)
	require.NoError(t, err)
	require.Equal(t, address, account.Address)

	missingAddress := cltest.NewAddress()
	account, err = store.KeyStore.GetAccountByAddress(missingAddress)
	require.EqualError(t, err, "no account found with that address")
}
