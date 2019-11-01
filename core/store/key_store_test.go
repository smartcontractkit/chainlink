package store_test

import (
	"io/ioutil"
	"testing"

	"chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
)

const correctPassphrase = "p@ssword"

func TestCreateEthereumAccount(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, err := store.KeyStore.NewAccount(correctPassphrase)
	assert.NoError(t, err)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Equal(t, 1, len(files))
}

func TestUnlockKey_SingleAddress(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	store.KeyStore.NewAccount(correctPassphrase)

	assert.Error(t, store.KeyStore.Unlock("wrong phrase"))
	assert.NoError(t, store.KeyStore.Unlock(correctPassphrase))
}

func TestUnlockKey_MultipleAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		passphrase1, passphrase2 string
		wantErr                  bool
	}{
		{"correct", correctPassphrase, correctPassphrase, false},
		{"first wrong", "wrong", correctPassphrase, true},
		{"second wrong", correctPassphrase, "wrong", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			store.KeyStore.NewAccount(test.passphrase1)
			store.KeyStore.NewAccount(test.passphrase2)

			if test.wantErr {
				assert.Error(t, store.KeyStore.Unlock(correctPassphrase))
			} else {
				assert.NoError(t, store.KeyStore.Unlock(correctPassphrase))
			}
		})
	}
}

func TestKeyStore_SignHash_Success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, err := store.KeyStore.NewAccount(correctPassphrase)
	assert.NoError(t, err)
	assert.NoError(t, store.KeyStore.Unlock(correctPassphrase))

	_, err = store.KeyStore.SignHash(cltest.StringToHash("abc123"))
	assert.NoError(t, err)
}
