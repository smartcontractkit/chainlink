package vrf_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsumerBaseRejectsBadVRFCoordinator(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _ /* jobID */, fee := registerProvingKey(t, coordinator)
	log := requestRandomness(t, coordinator, keyHash, fee, big.NewInt(1) /* seed */)
	// Ensure that VRFConsumerBase.rawFulfillRandomness's check,
	// require(msg.sender==vrfCoordinator), by using the wrong sender address.
	_, err := coordinator.consumerContract.RawFulfillRandomness(coordinator.neil,
		keyHash, big.NewInt(0).SetBytes([]byte("a bad random value")))
	require.Error(t, err)
	// Verify that correct fulfillment is possible, in this setup
	_ = fulfillRandomnessRequest(t, coordinator, *log)
}
