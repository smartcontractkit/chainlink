package store_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	tvrf "github.com/smartcontractkit/chainlink/core/internal/cltest/vrf"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_verifier_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
)

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func vrfVerifier() *solidity_vrf_verifier_wrapper.VRFTestHelper {
	ethereumKey, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(ethereumKey)
	genesisData := core.GenesisAlloc{auth.From: {Balance: big.NewInt(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	_, _, verifier, err := solidity_vrf_verifier_wrapper.DeployVRFTestHelper(auth, backend)
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
	require.NoError(t, err, "could not create encrypted key")
	require.NoError(t, ks.Forget(key),
		"could not forget a created key from in-memory store")
	keys, err := ks.Get() // Test generic Get
	require.NoError(t, err, "failed to retrieve expected key from db")
	assert.True(t, len(keys) == 1 && keys[0].PublicKey == key,
		"did not get back the expected key from  db retrial")
	ophrase := phrase + "corruption" // Extra key; make sure it's not returned by Get
	newKey, err := ks.CreateKey(ophrase, vrfkey.FastScryptParams)
	require.NoError(t, err, "could not create extra key")
	keys, err = ks.Get(key) // Test targeted Get
	require.NoError(t, err, "key databese retrieval failed")
	require.NoError(t, ks.Forget(newKey),
		"failed to forget in-memory copy of second key")
	require.Equal(t, keys[0].PublicKey, key, "retrieved wrong key from db")
	require.Len(t, keys, 1, "retrieved more keys than expected from db")
	keys, err = ks.Get() // Verify both keys are present in the db
	require.NoError(t, err, "could not retrieve keys from db")
	require.Len(t, keys, 2, "failed to remember both the keys just created")
	unlockedKeys, err := ks.Unlock(phrase) // Unlocking enables generation of proofs
	require.Contains(t, err.Error(), "could not decrypt key with given password",
		"should have a complaint about not being able to unlock the key with a different password")
	assert.Contains(t, err.Error(), newKey.String(),
		"complaint about inability to unlock should pertain to the key with a different password")
	assert.Len(t, unlockedKeys, 1, "should have only unlocked one key")
	assert.Equal(t, unlockedKeys[0], key,
		"should have only unlocked the key with the offered password")
	blockHash := common.Hash{}
	blockNum := 0
	preSeed := big.NewInt(10)
	s := tvrf.SeedData(t, preSeed, blockHash, blockNum)
	proof, err := ks.GenerateProof(key, s)
	assert.NoError(t, err,
		"should be able to generate VRF proofs with unlocked keys")
	// ...but only for unlocked keys
	_, err = ks.GenerateProof(newKey, s)
	require.Error(t, err,
		"should not be able to generate VRF proofs unless key has been unlocked")
	require.Contains(t, err.Error(), "has not been unlocked",
		"complaint when attempting to generate VRF proof with unclocked key should be that it's locked")
	encryptedKey, err := ks.GetSpecificKey(key) // Can export a key to bytes
	require.NoError(t, err, "should be able to get a specific key")
	assert.True(t, bytes.Equal(encryptedKey.PublicKey[:], key[:]),
		"should have recovered the encrypted key for the requested public key")
	verifier := vrfVerifier() // Generated proof is valid
	coordinatorProof, err := vrf.UnmarshalProofResponse(proof)
	require.NoError(t, err)
	verifierProof, err := coordinatorProof.ActualProof(s)
	require.NoError(t, err, "recovered bad VRF proof")
	wireProof, err := verifierProof.MarshalForSolidityVerifier()
	require.NoError(t, err, "could not marshal vrf proof for on-chain verification")
	_, err = verifier.RandomValueFromVRFProof(nil, wireProof[:])
	require.NoError(t, err,
		"failed to get VRF proof output from solidity VRF contract")
	require.NoError(t, ks.Delete(key), "failed to delete VRF key")
	_, err = ks.GenerateProof(key, s)
	require.Error(t, err,
		"should not be able to generate VRF proofs with a deleted key")
	require.Contains(t, err.Error(), "has not been unlocked",
		"complaint when trying to prove with deleted key should be that it's locked")
	keys, err = ks.Get(key) // Deleted key is removed from DB
	require.NoError(t, err, "failed to query db for key")
	require.Len(t, keys, 0, "deleted key should not be retrieved by db query")
	keyjson, err := encryptedKey.JSON()
	require.NoError(t, err, "failed to serialize key to JSON")
	require.NoError(t, ks.Import(keyjson, phrase),
		"failed to import encrypted key to database")
	err = ks.Import(keyjson, phrase)
	require.Equal(t, strpkg.MatchingVRFKeyError, err,
		"should be prevented from importing a key with a public key already "+
			"present in the DB")
	_, err = ks.GenerateProof(key, s)
	require.NoError(t, err, "should be able to generate proof with unlocked key")
}
