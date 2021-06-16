package vrf_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMeasureFulfillmentGasCost establishes rough bounds on the cost of
// providing a proof to the VRF coordinator.
func TestMeasureFulfillmentGasCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := newVRFCoordinatorUniverse(t, key)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	// Set up a request to fulfill
	log := requestRandomness(t, coordinator, keyHash, fee)
	preSeed, err := vrf.BigToSeed(log.Seed)
	require.NoError(t, err, "pre-seed %x out of range", preSeed)
	s := vrf.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: log.Raw.Raw.BlockHash,
		BlockNum:  log.Raw.Raw.BlockNumber,
	}
	seed := vrf.FinalSeed(s)
	proof, err := secretKey.GenerateProofWithNonce(seed, big.NewInt(1) /* nonce */)
	require.NoError(t, err)
	proofBlob, err := vrf.GenerateProofResponseFromProof(proof, s)
	require.NoError(t, err, "could not generate VRF proof!")
	coordinator.backend.Commit() // Work around simbackend/EVM block number bug
	estimate := estimateGas(t, coordinator.backend, coordinator.neil.From,
		coordinator.rootContractAddress, coordinator.coordinatorABI,
		"fulfillRandomnessRequest", proofBlob[:])

	assert.Greater(t, estimate, uint64(108000),
		"fulfillRandomness tx cost less gas than expected")
	// Note that this is probably a very loose upper bound on gas usage.
	// TODO:https://www.pivotaltracker.com/story/show/175040572
	assert.Less(t, estimate, uint64(500000),
		"fulfillRandomness tx cost more gas than expected")
}
