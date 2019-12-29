package cltest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"chainlink/core/store"
	strpkg "chainlink/core/store"
	"chainlink/core/store/models/vrfkey"
)

// StoredVRFKey creates a VRFKeyStore on store, imports a known VRF key into it,
// and returns the corresponding public key.
func StoredVRFKey(t *testing.T, store *strpkg.Store) *vrfkey.PublicKey {
	store.VRFKeyStore = strpkg.NewVRFKeyStore(store)
	keyFile, err := ioutil.ReadFile("../../tools/clroot/vrfkey.json")
	require.NoError(t, err)
	rawPassword, err := ioutil.ReadFile("../../tools/clroot/password.txt")
	require.NoError(t, err)
	password := strings.TrimSpace(string(rawPassword))
	require.NoError(t, store.VRFKeyStore.Import(keyFile, password))
	keys, err := store.VRFKeyStore.Unlock(password) // Extracts public key
	require.NoError(t, err)
	require.Len(t, keys, 1)
	return &keys[0]
}

func CreateVRFKey(t *testing.T, store *store.Store) vrfkey.PublicKey {
	key := vrfkey.CreateKey()
	store.VRFKeyStore.StoreInMemoryXXXTestingOnly(key)
	return key.PublicKey
}
