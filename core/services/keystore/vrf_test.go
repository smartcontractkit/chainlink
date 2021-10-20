package keystore_test

import (
	"bytes"
	"math/big"
	"testing"

	proof2 "github.com/smartcontractkit/chainlink/core/services/vrf/proof"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/smartcontractkit/chainlink/core/services/keystore"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_verifier_wrapper"
)

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func vrfVerifier(t *testing.T) *solidity_vrf_verifier_wrapper.VRFTestHelper {
	ethereumKey, _ := crypto.GenerateKey()
	auth := cltest.MustNewSimulatedBackendKeyedTransactor(t, ethereumKey)
	genesisData := core.GenesisAlloc{auth.From: {Balance: assets.Ether(100)}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	_, _, verifier, err := solidity_vrf_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
	backend.Commit()
	return verifier
}

var phrase = "engelbert humperdinck is the greatest musician of all time"

func TestKeyStoreEndToEnd(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ks := cltest.NewKeyStore(t, store.DB).VRF()
	ks.Unlock(phrase)
	key, err := ks.CreateKey() // NB: Varies from run to run. Shouldn't matter, though
	require.NoError(t, err, "could not create encrypted key")
	require.NoError(t, ks.Forget(key), "could not forget a created key from in-memory store")

	keys, err := ks.Get() // Test generic Get
	require.NoError(t, err, "failed to retrieve expected key from db")
	assert.True(t, len(keys) == 1 && keys[0].PublicKey == key, "did not get back the expected key from db retrial")

	ophrase := phrase + "corruption" // Cannot unlock with the wrong phrase
	_, err = ks.Unlock(ophrase)
	require.Error(t, err)

	keys, err = ks.Get(key) // Test targeted Get
	require.NoError(t, err, "key database retrieval failed")
	require.Equal(t, keys[0].PublicKey, key, "retrieved wrong key from db")
	require.Len(t, keys, 1, "retrieved more keys than expected from db")

	keys, err = ks.Get() // Verify both keys are present in the db
	require.NoError(t, err, "could not retrieve keys from db")
	require.Len(t, keys, 1, "failed to remember the key just created")

	unlockedKeys, err := ks.Unlock(phrase) // Unlocking enables generation of proofs
	require.NoError(t, err)
	assert.Len(t, unlockedKeys, 1, "should have only unlocked one key")
	assert.Equal(t, unlockedKeys[0], key, "should have only unlocked the key with the offered password")

	blockHash := common.Hash{}
	blockNum := 0
	preSeed := big.NewInt(10)
	seed := proof2.TestXXXSeedData(t, preSeed, blockHash, blockNum)

	proof, err := proof2.GenerateProofResponse(ks, key, seed)
	require.NoError(t, err, "failed to generate proof response")

	// ...but only for unlocked keys
	randomKey := vrfkey.CreateKey()
	_, err = proof2.GenerateProofResponse(ks, randomKey.PublicKey, seed)
	require.Error(t, err, "should not be able to generate VRF proofs unless key has been unlocked")
	require.Contains(t, err.Error(), "has not been unlocked", "complaint when attempting to generate VRF proof with unclocked key should be that it's locked")

	encryptedKey, err := ks.GetSpecificKey(key) // Can export a key to bytes
	require.NoError(t, err, "should be able to get a specific key")
	assert.True(t, bytes.Equal(encryptedKey.PublicKey[:], key[:]), "should have recovered the encrypted key for the requested public key")

	verifier := vrfVerifier(t) // Generated proof is valid
	coordinatorProof, err := proof2.UnmarshalProofResponse(proof)
	require.NoError(t, err)

	verifierProof, err := coordinatorProof.CryptoProof(seed)
	require.NoError(t, err, "recovered bad VRF proof")

	wireProof, err := proof2.MarshalForSolidityVerifier(&verifierProof)
	require.NoError(t, err, "could not marshal vrf proof for on-chain verification")

	_, err = verifier.RandomValueFromVRFProof(nil, wireProof[:])
	require.NoError(t, err, "failed to get VRF proof output from solidity VRF contract")

	err = ks.Delete(key)
	require.NoError(t, err, "failed to delete VRF key")

	_, err = proof2.GenerateProofResponse(ks, key, seed)
	require.Error(t, err, "should not be able to generate VRF proofs with a deleted key")
	require.Contains(t, err.Error(), "has not been unlocked", "complaint when trying to prove with deleted key should be that it's locked")

	keys, err = ks.Get(key) // Deleted key is removed from DB
	require.NoError(t, err, "failed to query db for key")
	require.Len(t, keys, 0, "deleted key should not be retrieved by db query")

	keyjson, err := encryptedKey.JSON()
	require.NoError(t, err, "failed to serialize key to JSON")

	_, err = ks.Import(keyjson, phrase)
	require.NoError(t, err, "failed to import encrypted key to database")

	_, err = ks.Import(keyjson, phrase)
	require.Equal(t, keystore.ErrMatchingVRFKey, err, "should be prevented from importing a key with a public key already present in the DB")

	_, err = proof2.GenerateProofResponse(ks, key, seed)
	require.NoError(t, err, "should be able to generate proof with unlocked key")
}
