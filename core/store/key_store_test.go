package store_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	gethCommon "github.com/ethereum/go-ethereum/common"
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
	assert.Len(t, files, 2)
}

func TestUnlockKey_SingleAddress(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Verify the fixture account
	require.True(t, store.KeyStore.HasAccounts())
	require.Len(t, store.KeyStore.GetAccounts(), 1)

	assert.EqualError(t, store.KeyStore.Unlock("wrong phrase"), "invalid password for account 0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea; could not decrypt key with given password")
	assert.NoError(t, store.KeyStore.Unlock(cltest.Password))
}

func TestUnlockKey_MultipleAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		tryPassphrase        string
		secondAcctPassphrase string
		wantErr              bool
	}{
		{"correct", cltest.Password, cltest.Password, false},
		{"first wrong", "wrong", cltest.Password, true},
		{"second wrong", cltest.Password, "wrong", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			// Verify the fixture account
			require.True(t, store.KeyStore.HasAccounts())
			require.Len(t, store.KeyStore.GetAccounts(), 1)
			defer cleanup()

			_, err := store.KeyStore.NewAccount(test.secondAcctPassphrase)
			require.NoError(t, err)

			if test.wantErr {
				assert.Error(t, store.KeyStore.Unlock(test.tryPassphrase))
			} else {
				assert.NoError(t, store.KeyStore.Unlock(test.tryPassphrase))
			}
		})
	}
}

func TestKeyStore_SignHash_Success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	assert.NoError(t, store.KeyStore.Unlock(cltest.Password))

	_, err := store.KeyStore.SignHash(cltest.StringToHash("abc123"))
	assert.NoError(t, err)
}

func TestKeyStore_GetAccountByAddress(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	address := gethCommon.HexToAddress("0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea")
	account, err := store.KeyStore.GetAccountByAddress(address)
	require.NoError(t, err)
	require.Equal(t, address, account.Address)

	missingAddress := cltest.NewAddress()
	account, err = store.KeyStore.GetAccountByAddress(missingAddress)
	require.EqualError(t, err, "no account found with that address")
}
