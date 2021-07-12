package cltest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/stretchr/testify/require"
)

// StoredVRFKey creates a VRFKeyStore on store, imports a known VRF key into it,
// and returns the corresponding public key.
func StoredVRFKey(t *testing.T, ks *keystore.VRF) *secp256k1.PublicKey {
	keyFile, err := ioutil.ReadFile("../../tools/clroot/vrfkey.json")
	require.NoError(t, err)
	rawPassword, err := ioutil.ReadFile("../../tools/clroot/password.txt")
	require.NoError(t, err)
	password := strings.TrimSpace(string(rawPassword))
	_, err = ks.Import(keyFile, password)
	require.NoError(t, err)
	keys, err := ks.Unlock(password) // Extracts public key
	require.NoError(t, err)
	require.Len(t, keys, 1)
	return &keys[0]
}
