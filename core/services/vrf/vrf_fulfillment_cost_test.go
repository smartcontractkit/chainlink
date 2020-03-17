package vrf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureFulfillmenttGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	log := requestRandomness(t, coordinator, keyHash, fee, seed)
	proof, err := generateProofWithNonce(secretKey, log.Seed, one /* nonce */)
	require.NoError(t, err, "could not generate VRF proof!")
	proofBlob, err := proof.MarshalForSolidityVerifier()
	require.NoError(t, err, "could not marshal VRF proof for VRFCoordinator!")

	estimate := estimateGas(t, coordinator.backend, coordinator.neil.From,
		coordinator.rootContractAddress, coordinator.coordinatorABI,
		"fulfillRandomnessRequest", proofBlob[:])

}
