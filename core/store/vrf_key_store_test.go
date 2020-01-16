package store_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"chainlink/core/internal/cltest"
	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/services/vrf/generated/solidity_verifier_wrapper"
	strpkg "chainlink/core/store"
	"chainlink/core/store/models/vrfkey"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func vRFVerifier() *solidity_verifier_wrapper.VRFTestHelper {
	ethereumKey, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(ethereumKey)
	genesisData := core.GenesisAlloc{auth.From: {Balance: big.NewInt(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	_, _, verifier, err := solidity_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
	backend.Commit()
	return verifier
}

var phrase = "englebert humperdinck is the greatest musician of all time"

func TestKeyStoreEndToEnd(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ks := strpkg.NewVRFKeyStore(store)
	key, err := ks.CreateKey(phrase, vrfkey.FastScryptParams) // NB: Varies from run to run. Shouldn't matter, though
	require.NoError(t, err)
	ks.Forget(key)
	require.NoError(t, err)
	keys, err := ks.Get(nil) // Test generic Get
	require.NoError(t, err)
	require.True(t, len(keys) == 1 && keys[0].PublicKey == *key)
	ophrase := phrase + "corruption" // Extra key; make sure it's not returned by Get
	newKey, err := ks.CreateKey(ophrase, vrfkey.FastScryptParams)
	require.NoError(t, err)
	keys, err = ks.Get(key) // Test targeted Get
	require.NoError(t, err)
	ks.Forget(newKey) // Remove second key from memory
	require.Equal(t, keys[0].PublicKey, *key)
	require.Len(t, keys, 1)
	keys, err = ks.Get(nil) // Verify both keys are present in the db
	require.NoError(t, err)
	require.Len(t, keys, 2, "failed to remember both the keys just created")
	unlockedKeys, err := ks.Unlock(phrase) // Unlocking enables generation of proofs
	require.Contains(t, err.Error(), "could not decrypt key with given password")
	require.Contains(t, err.Error(), newKey.String())
	require.Len(t, unlockedKeys, 1)
	require.Equal(t, unlockedKeys[0], *key)
	proof, err := ks.GenerateProof(key, big.NewInt(10))
	require.NoError(t, err)
	_, err = ks.GenerateProof(newKey, big.NewInt(10)) // ...but only for unlocked keys
	require.Error(t, err)
	require.Contains(t, err.Error(), "has not been unlocked")
	export, err := ks.Export(key) // Can export a key to bytes
	require.NoError(t, err)
	require.Len(t, export, 1)
	verifier := vRFVerifier() // Generated proof is valid
	_, err = verifier.RandomValueFromVRFProof(nil, proof[:])
	require.NoError(t, err)
	require.NoError(t, ks.Delete(key))             // Deleting actually deletes
	_, err = ks.GenerateProof(key, big.NewInt(10)) // Can't prove with deleted key
	require.Error(t, err)
	require.Contains(t, err.Error(), "has not been unlocked")
	keys, err = ks.Get(key) // Deleted key is removed from DB
	require.NoError(t, err)
	require.Len(t, keys, 0)
	require.NoError(t, ks.Import(export[0], phrase)) // Can re-use key after re-importing
	require.Error(t, ks.Import(export[0], phrase))   // Can't import over existing DB key
	_, err = ks.GenerateProof(key, big.NewInt(10))
	require.NoError(t, err)
}
