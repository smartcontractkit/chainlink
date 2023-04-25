package vrf_test

import (
	"math/big"
	"testing"

	proof2 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

// TestMeasureFulfillmentGasCost establishes rough bounds on the cost of
// providing a proof to the VRF coordinator.
func TestMeasureFulfillmentGasCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := newVRFCoordinatorUniverse(t, key)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	// Set up a request to fulfill
	log := requestRandomness(t, coordinator, keyHash, fee)
	preSeed, err := proof2.BigToSeed(log.Seed)
	require.NoError(t, err, "pre-seed %x out of range", preSeed)
	s := proof2.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: log.Raw.Raw.BlockHash,
		BlockNum:  log.Raw.Raw.BlockNumber,
	}
	seed := proof2.FinalSeed(s)
	proof, err := secretKey.GenerateProofWithNonce(seed, big.NewInt(1) /* nonce */)
	require.NoError(t, err)
	proofBlob, err := GenerateProofResponseFromProof(proof, s)
	require.NoError(t, err, "could not generate VRF proof!")
	coordinator.backend.Commit() // Work around simbackend/EVM block number bug
	estimate := estimateGas(t, coordinator.backend, coordinator.neil.From,
		coordinator.rootContractAddress, coordinator.coordinatorABI,
		"fulfillRandomnessRequest", proofBlob[:])

	assert.Greater(t, estimate, uint64(108000),
		"fulfillRandomness tx cost less gas than expected")
	t.Log("estimate", estimate)
	// Note that this is probably a very loose upper bound on gas usage.
	// TODO:https://www.pivotaltracker.com/story/show/175040572
	assert.Less(t, estimate, uint64(500000),
		"fulfillRandomness tx cost more gas than expected")
}
