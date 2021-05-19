package cltest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/require"

	strpkg "github.com/smartcontractkit/chainlink/core/store"
)

// StoredVRFKey creates a VRFKeyStore on store, imports a known VRF key into it,
// and returns the corresponding public key.
func StoredVRFKey(t *testing.T, store *strpkg.Store) *secp256k1.PublicKey {
	store.VRFKeyStore = vrf.NewVRFKeyStore(vrf.NewORM(store.DB), utils.GetScryptParams(store.Config))
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
