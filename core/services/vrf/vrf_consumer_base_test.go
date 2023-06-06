package vrf_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func TestConsumerBaseRejectsBadVRFCoordinator(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := newVRFCoordinatorUniverse(t, key)
	keyHash, _ /* jobID */, fee := registerProvingKey(t, coordinator)
	log := requestRandomness(t, coordinator, keyHash, fee)
	// Ensure that VRFConsumerBase.rawFulfillRandomness's check,
	// require(msg.sender==vrfCoordinator), by using the wrong sender address.
	_, err := coordinator.consumerContract.RawFulfillRandomness(coordinator.neil,
		keyHash, big.NewInt(0).SetBytes([]byte("a bad random value")))
	require.Error(t, err)
	// Verify that correct fulfilment is possible, in this setup
	_ = fulfillRandomnessRequest(t, coordinator, *log)
}
