package vrf_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureFulfillmenttGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	// Set up a request to fulfill
	log := requestRandomness(t, coordinator, keyHash, fee, seed)
	proof, err := vrf.GenerateProofWithNonce(rawSecretKey, log.Seed,
		big.NewInt(1) /* nonce */)
	require.NoError(t, err, "could not generate VRF proof!")
	// Set up the proof with which to fulfill request
	proofBlob, err := proof.MarshalForSolidityVerifier()
	require.NoError(t, err, "could not marshal VRF proof for VRFCoordinator!")

	estimate := estimateGas(t, coordinator.backend, coordinator.neil.From,
		coordinator.rootContractAddress, coordinator.coordinatorABI,
		"fulfillRandomnessRequest", proofBlob[:])

	assert.Greater(t, estimate, uint64(145000),
		"fulfillRandomness tx cost less gas than expected")
	assert.Less(t, estimate, uint64(300000),
		"fulfillRandomness tx cost more gas than expected")
}
