package vrf_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment_NoCancel(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.noCancelConsumers[0],
		uni.noCancelConsumerAddresses[0],
		uni.noCancelCoordinator,
		uni.noCancelAddress,
		uni.noCancelBatchCoordinatorAddress,
		5,     // number of requests to send
		false, // don't send big callback
		func(t *testing.T, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
			// all fulfillments should have the same flat cost
			feeConfig, err := coordinator.GetFeeConfig(nil)
			require.NoError(t, err)
			expectedPayment := big.NewInt(int64(feeConfig.FulfillmentFlatFeeLinkPPMTier1) * int64(1e12))
			require.Equal(t, expectedPayment, rwfe.Payment)
		},
	)
}

func TestVRFV2Integration_SingleConsumer_HappyPath_NoCancel(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPath(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.noCancelConsumers[0],
		uni.noCancelConsumerAddresses[0],
		uni.noCancelCoordinator,
		uni.noCancelAddress,
		uni.noCancelBatchCoordinatorAddress,
		func(t *testing.T, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
			// all fulfillments should have the same flat cost
			feeConfig, err := coordinator.GetFeeConfig(nil)
			require.NoError(t, err)
			expectedPayment := big.NewInt(int64(feeConfig.FulfillmentFlatFeeLinkPPMTier1) * int64(1e12))
			require.Equal(t, expectedPayment, rwfe.Payment)
		})
}

func TestVRFV2Integration_SingleConsumer_HappyPath_BatchFulfillment_BigGasCallback_NoCancel(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerHappyPathBatchFulfillment(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.noCancelConsumers[0],
		uni.noCancelConsumerAddresses[0],
		uni.noCancelCoordinator,
		uni.noCancelAddress,
		uni.noCancelBatchCoordinatorAddress,
		5,    // number of requests to send
		true, // send big callback
		func(t *testing.T, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
			// all fulfillments should have the same flat cost
			feeConfig, err := coordinator.GetFeeConfig(nil)
			require.NoError(t, err)
			expectedPayment := big.NewInt(int64(feeConfig.FulfillmentFlatFeeLinkPPMTier1) * int64(1e12))
			require.Equal(t, expectedPayment, rwfe.Payment)
		},
	)
}

func TestVRFV2Integration_SingleConsumer_NeedsBlockhashStore_NoCancel(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 2)
	testMultipleConsumersNeedBHS(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers,
		uni.noCancelConsumers,
		uni.noCancelConsumerAddresses,
		uni.noCancelCoordinator,
		uni.noCancelAddress,
		uni.noCancelBatchCoordinatorAddress,
		func(t *testing.T, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
			// all fulfillments should have the same flat cost
			feeConfig, err := coordinator.GetFeeConfig(nil)
			require.NoError(t, err)
			expectedPayment := big.NewInt(int64(feeConfig.FulfillmentFlatFeeLinkPPMTier1) * int64(1e12))
			require.Equal(t, expectedPayment, rwfe.Payment)
		},
	)
}

func TestVRFV2Integration_SingleConsumer_NeedsTopUp_NoCancel(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testSingleConsumerNeedsTopUp(
		t,
		ownerKey,
		uni,
		uni.vrfConsumers[0],
		uni.noCancelConsumers[0],
		uni.noCancelConsumerAddresses[0],
		uni.noCancelCoordinator,
		uni.noCancelAddress,
		uni.noCancelBatchCoordinatorAddress,
		big.NewInt(1),           // initial funding of 1 juel
		assets.Ether(1).ToInt(), // top up of 1 LINK
		func(t *testing.T, coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface, rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled) {
			// all fulfillments should have the same flat cost
			feeConfig, err := coordinator.GetFeeConfig(nil)
			require.NoError(t, err)
			expectedPayment := big.NewInt(int64(feeConfig.FulfillmentFlatFeeLinkPPMTier1) * int64(1e12))
			require.Equal(t, expectedPayment, rwfe.Payment)
		},
	)
}
