package store_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chainlink/core/internal/cltest"
	"chainlink/core/services/signatures/secp256k1"
	strpkg "chainlink/core/store"
	"chainlink/core/store/models/vrfkey"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
var phrase = "englebert humperdinck is the greatest musician of all time"

func TestKeyStoreEndToEnd(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ks := strpkg.NewVRFKeyStore(store)
	key, err := ks.CreateKey(phrase, vrfkey.FastScryptParams) // NB: Varies from run to run. Shouldn't matter, though
	require.NoError(t, err, "could not create encrypted key")
	require.NoError(t, ks.Forget(key),
		"could not forget a created key from in-memory store")
	keys, err := ks.Get(nil) // Test generic Get
	require.NoError(t, err, "failed to retrieve expected key from db")
	assert.True(t, len(keys) == 1 && keys[0].PublicKey == *key,
		"did not get back the expected key from  db retrial")
	ophrase := phrase + "corruption" // Extra key; make sure it's not returned by Get
	newKey, err := ks.CreateKey(ophrase, vrfkey.FastScryptParams)
	require.NoError(t, err, "could not create extra key")
	keys, err = ks.Get(key) // Test targeted Get
	require.NoError(t, err, "key databese retrieval failed")
	require.NoError(t, ks.Forget(newKey),
		"failed to forget in-memory copy of second key")
	require.Equal(t, keys[0].PublicKey, *key, "retrieved wrong key from db")
	require.Len(t, keys, 1, "retrieved more keys than expected from db")
	keys, err = ks.Get(nil) // Verify both keys are present in the db
	require.NoError(t, err, "could not retrieve keys from db")
	require.Len(t, keys, 2, "failed to remember both the keys just created")
	unlockedKeys, err := ks.Unlock(phrase) // Unlocking enables generation of proofs
	require.Contains(t, err.Error(), "could not decrypt key with given password",
		"should have a complaint about not being able to unlock the key with a different password")
	assert.Contains(t, err.Error(), newKey.String(),
		"complaint about inability to unlock should pertain to the key with a different password")
	assert.Len(t, unlockedKeys, 1, "should have only unlocked one key")
	assert.Equal(t, unlockedKeys[0], *key,
		"should have only unlocked the key with the offered password")
	encryptedKey, err := ks.GetSpecificKey(key) // Can export a key to bytes
	require.NoError(t, err, "should be able to get a specific key")
	assert.True(t, bytes.Equal(encryptedKey.PublicKey[:], key[:]),
		"should have recovered the encrypted key for the requested public key")
	require.NoError(t, ks.Delete(key), "failed to delete VRF key")
	keys, err = ks.Get(key) // Deleted key is removed from DB
	require.NoError(t, err, "failed to query db for key")
	require.Len(t, keys, 0, "deleted key should not be retrieved by db query")
	keyjson, err := encryptedKey.JSON()
	require.NoError(t, err, "failed to serialize key to JSON")
	require.NoError(t, ks.Import(keyjson, phrase),
		"failed to import encrypted key to database")
	err = ks.Import(keyjson, phrase)
	require.Equal(t, strpkg.MatchingVRFKeyError, err,
		"should be prevented from importing a key with a public key already present in the DB")
}
