package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/core/store"

	"github.com/stretchr/testify/require"
)

const ocrPassphrase = "davie bowie is the greatest musician of all time"

func assertKeyCount(t *testing.T, store *strpkg.Store, n int) {
	encryptedKeys, err := store.FindEncryptedOCRKeys()
	require.NoError(t, err, "failed to retrieve keys from db")
	require.Len(t, encryptedKeys, 1)
}

func TestOCRKeyStoreEndToEnd(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ks := strpkg.NewOCRKeyStore(store)

	assertKeyCount := func(n int) {
		encryptedKeys, err := store.FindEncryptedOCRKeys()
		require.NoError(t, err, "failed to retrieve keys from db")
		require.Len(t, encryptedKeys, n)
	}

	createdKey, err := ks.CreateWeakKeyXXXTestingOnly(ocrPassphrase)
	require.NoError(t, err, "could not create OCR key")
	assertKeyCount(1)
	// retrieve key from in-memory keystore
	retreivedKey, err := ks.Get(createdKey.ID)
	require.NoError(t, err, "could not retrieve key from in-memory keystore")
	require.Equal(t, createdKey, retreivedKey)
	// forget key from in-memory keystore
	err = ks.Forget(retreivedKey)
	require.NoError(t, err, "could not forget key from in-memory keystore")
	require.False(t, ks.Has(createdKey.ID))
	assertKeyCount(1) // still in DB
	// load key back into keystore
	ids, err := ks.Unlock(ocrPassphrase)
	require.NoError(t, err)
	require.Len(t, ids, 1)
	require.Equal(t, ids[0], createdKey.ID)
	require.True(t, ks.Has(createdKey.ID))
	require.Equal(t, createdKey.ID, retreivedKey.ID)
	require.Equal(t, createdKey.PublicKeyOffChain(), retreivedKey.PublicKeyOffChain())
	require.Equal(t, createdKey.PublicKeyAddressOnChain(), retreivedKey.PublicKeyAddressOnChain())
	assertKeyCount(1)
	// delete key
	err = ks.Delete(createdKey)
	require.NoError(t, err, "could not delete key from keystore & DB")
	require.False(t, ks.Has(createdKey.ID))
	assertKeyCount(0)
}
