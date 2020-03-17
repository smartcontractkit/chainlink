package vrf

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
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
	rawData, err := coordinator.coordinatorABI.Pack("fulfillRandomnessRequest",
		proofBlob[:])
	require.NoError(t, err)
	callMsg := ethereum.CallMsg{From: coordinator.neil.From,
		To: &coordinator.rootContractAddress, Data: rawData}
	estimate, err := coordinator.backend.EstimateGas(context.TODO(), callMsg)
	require.NoError(t, err)
	assert.Greater(t, estimate, uint64(148000))
	assert.Less(t, estimate, uint64(200000))
}
