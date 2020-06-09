package vrf_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureFulfillmentGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	// Set up a request to fulfill
	log := requestRandomness(t, coordinator, keyHash, fee, seed)
	preseed, err := vrf.BigToSeed(log.Seed)
	require.NoError(t, err, "preseed %x out of range", preseed)
	s := vrf.PreSeedData{
		PreSeed:   preseed,
		BlockHash: log.Raw.Raw.BlockHash,
		BlockNum:  log.Raw.Raw.BlockNumber,
	}
	proofBlob, err := vrf.GenerateProofResponseWithNonce(rawSecretKey, s,
		big.NewInt(1) /* nonce */)
	require.NoError(t, err, "could not generate VRF proof!")
	coordinator.backend.Commit() // Work around simbackend/EVM block number bug
	estimate := estimateGas(t, coordinator.backend, coordinator.neil.From,
		coordinator.rootContractAddress, coordinator.coordinatorABI,
		"fulfillRandomnessRequest", proofBlob[:])

	assert.Greater(t, estimate, uint64(108000),
		"fulfillRandomness tx cost less gas than expected")
	assert.Less(t, estimate, uint64(400000),
		"fulfillRandomness tx cost more gas than expected")
}
