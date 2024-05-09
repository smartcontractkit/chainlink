package smoke

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	it_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"
)

func TestVRFv2Plus(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFv2Plus universe")

	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("Link Billing", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = false
		consumers, subIDsForRequestRandomness, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForRequestRandomness := subIDsForRequestRandomness[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subIDForRequestRandomness, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)

		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subIDForRequestRandomness,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		require.False(t, randomWordsFulfilledEvent.OnlyPremium, "RandomWordsFulfilled Event's `OnlyPremium` field should be false")
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment, "RandomWordsFulfilled Event's `NativePayment` field should be false")
		require.True(t, randomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2Plus.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})
	t.Run("Native Billing", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		testConfig := configCopy.VRFv2Plus.General
		var isNativeBilling = true

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		subNativeTokenBalanceBeforeRequest := subscription.NativeBalance

		jobRunsBeforeTest, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		require.False(t, randomWordsFulfilledEvent.OnlyPremium)
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment)
		require.True(t, randomWordsFulfilledEvent.Success)
		expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err)
		subBalanceAfterRequest := subscription.NativeBalance
		require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

		jobRuns, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *testConfig.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})
	t.Run("Direct Funding (VRFV2PlusWrapper)", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperEnvironment(
			testcontext.Get(t),
			env,
			chainID,
			&configCopy,
			vrfContracts.LinkToken,
			vrfContracts.MockETHLINKFeed,
			vrfContracts.CoordinatorV2Plus,
			vrfKey.KeyHash,
			1,
		)
		require.NoError(t, err)

		t.Run("Link Billing", func(t *testing.T) {
			configCopy := config.MustCopy().(tc.TestConfig)
			testConfig := configCopy.VRFv2Plus.General
			var isNativeBilling = false

			wrapperConsumerJuelsBalanceBeforeRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperContracts.LoadTestConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.Balance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.LoadTestConsumers[0],
				vrfContracts.CoordinatorV2Plus,
				vrfKey,
				wrapperSubID,
				isNativeBilling,
				configCopy.VRFv2Plus.General,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.Balance
			require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)

			wrapperConsumerJuelsBalanceAfterRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperContracts.LoadTestConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")
			require.Equal(t, expectedWrapperConsumerJuelsBalance, wrapperConsumerJuelsBalanceAfterRequest)

			//todo: uncomment when VRF-651 will be fixed
			//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
			vrfv2plus.LogFulfillmentDetailsLinkBilling(l, wrapperConsumerJuelsBalanceBeforeRequest, wrapperConsumerJuelsBalanceAfterRequest, consumerStatus, randomWordsFulfilledEvent)

			require.Equal(t, *testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
			for _, w := range consumerStatus.RandomWords {
				l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
				require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
			}
		})
		t.Run("Native Billing", func(t *testing.T) {
			configCopy := config.MustCopy().(tc.TestConfig)
			testConfig := configCopy.VRFv2Plus.General
			var isNativeBilling = true

			wrapperConsumerBalanceBeforeRequestWei, err := evmClient.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()))
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.NativeBalance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.LoadTestConsumers[0],
				vrfContracts.CoordinatorV2Plus,
				vrfKey,
				wrapperSubID,
				isNativeBilling,
				configCopy.VRFv2Plus.General,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceWei := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.NativeBalance
			require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerWeiBalance := new(big.Int).Sub(wrapperConsumerBalanceBeforeRequestWei, consumerStatus.Paid)

			wrapperConsumerBalanceAfterRequestWei, err := evmClient.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()))
			require.NoError(t, err, "error getting wrapper consumer balance")
			require.Equal(t, expectedWrapperConsumerWeiBalance, wrapperConsumerBalanceAfterRequestWei)

			//todo: uncomment when VRF-651 will be fixed
			//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
			vrfv2plus.LogFulfillmentDetailsNativeBilling(l, wrapperConsumerBalanceBeforeRequestWei, wrapperConsumerBalanceAfterRequestWei, consumerStatus, randomWordsFulfilledEvent)

			require.Equal(t, *testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
			for _, w := range consumerStatus.RandomWords {
				l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
				require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
			}
		})
	})
	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		_, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceNativeBeforeSubCancelling, err := evmClient.BalanceAt(testcontext.Get(t), testWalletAddress)
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subID.String()).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")
		tx, err := vrfContracts.CoordinatorV2Plus.CancelSubscription(subID, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfContracts.CoordinatorV2Plus.WaitForSubscriptionCanceledEvent(subID, time.Second*30)
		require.NoError(t, err, "error waiting for subscription canceled event")

		cancellationTxReceipt, err := evmClient.GetTxReceipt(tx.Hash())
		require.NoError(t, err, "error getting tx cancellation Tx Receipt")

		txGasUsed := new(big.Int).SetUint64(cancellationTxReceipt.GasUsed)
		// we don't have that information for older Geth versions
		if cancellationTxReceipt.EffectiveGasPrice == nil {
			cancellationTxReceipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
		}
		cancellationTxFeeWei := new(big.Int).Mul(txGasUsed, cancellationTxReceipt.EffectiveGasPrice)

		l.Info().
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Effective Gas Price", cancellationTxReceipt.EffectiveGasPrice.String()).
			Uint64("Gas Used", cancellationTxReceipt.GasUsed).
			Msg("Cancellation TX Receipt")

		l.Info().
			Str("Returned Subscription Amount Native", subscriptionCanceledEvent.AmountNative.String()).
			Str("Returned Subscription Amount Link", subscriptionCanceledEvent.AmountLink.String()).
			Str("SubID", subscriptionCanceledEvent.SubId.String()).
			Str("Returned to", subscriptionCanceledEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceNative, subscriptionCanceledEvent.AmountNative, "SubscriptionCanceled event native amount is not equal to sub amount while canceling subscription")
		require.Equal(t, subBalanceLink, subscriptionCanceledEvent.AmountLink, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		testWalletBalanceNativeAfterSubCancelling, err := evmClient.BalanceAt(testcontext.Get(t), testWalletAddress)
		require.NoError(t, err)

		testWalletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		subFundsReturnedNativeActual := new(big.Int).Sub(testWalletBalanceNativeAfterSubCancelling, testWalletBalanceNativeBeforeSubCancelling)
		subFundsReturnedLinkActual := new(big.Int).Sub(testWalletBalanceLinkAfterSubCancelling, testWalletBalanceLinkBeforeSubCancelling)

		subFundsReturnedNativeExpected := new(big.Int).Sub(subBalanceNative, cancellationTxFeeWei)
		deltaSpentOnCancellationTxFee := new(big.Int).Sub(subBalanceNative, subFundsReturnedNativeActual)
		l.Info().
			Str("Sub Balance - Native", subBalanceNative.String()).
			Str("Delta Spent On Cancellation Tx Fee - `NativeBalance - subFundsReturnedNativeActual`", deltaSpentOnCancellationTxFee.String()).
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Native", subFundsReturnedNativeActual.String()).
			Str("Sub Funds Returned Expected - `NativeBalance - cancellationTxFeeWei`", subFundsReturnedNativeExpected.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Msg("Sub funds returned")

		//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
		//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
		require.Equal(t, 1, testWalletBalanceNativeAfterSubCancelling.Cmp(testWalletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")
		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

	})
	t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		testConfig := configCopy.VRFv2Plus.General

		//underfund subs in order rand fulfillments to fail
		testConfig.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))
		testConfig.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
		activeSubscriptionIdsBeforeSubCancellation, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err)

		require.True(t, it_utils.BigIntSliceContains(activeSubscriptionIdsBeforeSubCancellation, subID))

		pendingRequestsExist, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout = ptr.Ptr(blockchain.StrDuration{Duration: 5 * time.Second})
		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			false,
			configCopy.VRFv2Plus.General,
			l,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			true,
			configCopy.VRFv2Plus.General,
			l,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfulfilled rand requests due to low sub balance")

		walletBalanceNativeBeforeSubCancelling, err := evmClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		walletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subID.String()).
			Str("Returning funds to", defaultWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		tx, err := vrfContracts.CoordinatorV2Plus.OwnerCancelSubscription(subID)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfContracts.CoordinatorV2Plus.WaitForSubscriptionCanceledEvent(subID, time.Second*30)
		require.NoError(t, err, "error waiting for subscription canceled event")

		cancellationTxReceipt, err := evmClient.GetTxReceipt(tx.Hash())
		require.NoError(t, err, "error getting tx cancellation Tx Receipt")

		txGasUsed := new(big.Int).SetUint64(cancellationTxReceipt.GasUsed)
		// we don't have that information for older Geth versions
		if cancellationTxReceipt.EffectiveGasPrice == nil {
			cancellationTxReceipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
		}
		cancellationTxFeeWei := new(big.Int).Mul(txGasUsed, cancellationTxReceipt.EffectiveGasPrice)

		l.Info().
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Effective Gas Price", cancellationTxReceipt.EffectiveGasPrice.String()).
			Uint64("Gas Used", cancellationTxReceipt.GasUsed).
			Msg("Cancellation TX Receipt")

		l.Info().
			Str("Returned Subscription Amount Native", subscriptionCanceledEvent.AmountNative.String()).
			Str("Returned Subscription Amount Link", subscriptionCanceledEvent.AmountLink.String()).
			Str("SubID", subscriptionCanceledEvent.SubId.String()).
			Str("Returned to", subscriptionCanceledEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceNative, subscriptionCanceledEvent.AmountNative, "SubscriptionCanceled event native amount is not equal to sub amount while canceling subscription")
		require.Equal(t, subBalanceLink, subscriptionCanceledEvent.AmountLink, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		walletBalanceNativeAfterSubCancelling, err := evmClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		walletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		subFundsReturnedNativeActual := new(big.Int).Sub(walletBalanceNativeAfterSubCancelling, walletBalanceNativeBeforeSubCancelling)
		subFundsReturnedLinkActual := new(big.Int).Sub(walletBalanceLinkAfterSubCancelling, walletBalanceLinkBeforeSubCancelling)

		subFundsReturnedNativeExpected := new(big.Int).Sub(subBalanceNative, cancellationTxFeeWei)
		deltaSpentOnCancellationTxFee := new(big.Int).Sub(subBalanceNative, subFundsReturnedNativeActual)
		l.Info().
			Str("Sub Balance - Native", subBalanceNative.String()).
			Str("Delta Spent On Cancellation Tx Fee - `NativeBalance - subFundsReturnedNativeActual`", deltaSpentOnCancellationTxFee.String()).
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Native", subFundsReturnedNativeActual.String()).
			Str("Sub Funds Returned Expected - `NativeBalance - cancellationTxFeeWei`", subFundsReturnedNativeExpected.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Str("walletBalanceNativeBeforeSubCancelling", walletBalanceNativeBeforeSubCancelling.String()).
			Str("walletBalanceNativeAfterSubCancelling", walletBalanceNativeAfterSubCancelling.String()).
			Msg("Sub funds returned")

		//todo - need to use different wallet for each test to verify exact amount of Native/LINK returned
		//todo - as defaultWallet is used in other tests in parallel which might affect the balance - TT-684
		//require.Equal(t, 1, walletBalanceNativeAfterSubCancelling.Cmp(walletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")

		//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
		//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

		activeSubscriptionIdsAfterSubCancellation, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error getting active subscription ids")

		require.False(
			t,
			it_utils.BigIntSliceContains(activeSubscriptionIdsAfterSubCancellation, subID),
			"Active subscription ids should not contain sub id after sub cancellation",
		)
	})
	t.Run("Owner Withdraw", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		fulfilledEventLink, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			false,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err)

		fulfilledEventNative, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			true,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err)
		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceNativeBeforeWithdraw, err := evmClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		defaultWalletBalanceLinkBeforeWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfContracts.CoordinatorV2Plus.Withdraw(
			common.HexToAddress(defaultWalletAddress),
		)
		require.NoError(t, err, "error withdrawing LINK from coordinator to default wallet")
		amountToWithdrawNative := fulfilledEventNative.Payment

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawNative.String()).
			Msg("Invoking Oracle Withdraw for Native")

		err = vrfContracts.CoordinatorV2Plus.WithdrawNative(
			common.HexToAddress(defaultWalletAddress),
		)
		require.NoError(t, err, "error withdrawing Native tokens from coordinator to default wallet")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		defaultWalletBalanceNativeAfterWithdraw, err := evmClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		defaultWalletBalanceLinkAfterWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		//not possible to verify exact amount of Native/LINK returned as defaultWallet is used in other tests in parallel which might affect the balance
		require.Equal(t, 1, defaultWalletBalanceNativeAfterWithdraw.Cmp(defaultWalletBalanceNativeBeforeWithdraw), "Native funds were not returned after oracle withdraw native")
		require.Equal(t, 1, defaultWalletBalanceLinkAfterWithdraw.Cmp(defaultWalletBalanceLinkBeforeWithdraw), "LINK funds were not returned after oracle withdraw")
	})
}

func TestVRFv2PlusMultipleSendingKeys(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: 2,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = true

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		txKeys, _, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, newEnvConfig.NumberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < newEnvConfig.NumberOfTxKeysToCreate+1; i++ {
			randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
				consumers[0],
				vrfContracts.CoordinatorV2Plus,
				vrfKey,
				subID,
				isNativeBilling,
				configCopy.VRFv2Plus.General,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			//todo - move TransactionByHash to EVMClient in CTF
			fulfillmentTx, _, err := actions.GetTxByHash(testcontext.Get(t), evmClient, randomWordsFulfilledEvent.Raw.TxHash)
			require.NoError(t, err, "error getting tx from hash")
			fulfillmentTxFromAddress, err := actions.GetTxFromAddress(fulfillmentTx)
			require.NoError(t, err, "error getting tx from address")
			fulfillmentTxFromAddresses = append(fulfillmentTxFromAddresses, fulfillmentTxFromAddress)
		}
		require.Equal(t, newEnvConfig.NumberOfTxKeysToCreate+1, len(fulfillmentTxFromAddresses))
		var txKeyAddresses []string
		for _, txKey := range txKeys.Data {
			txKeyAddresses = append(txKeyAddresses, txKey.Attributes.Address)
		}
		less := func(a, b string) bool { return a < b }
		equalIgnoreOrder := cmp.Diff(txKeyAddresses, fulfillmentTxFromAddresses, cmpopts.SortSlices(less)) == ""
		require.True(t, equalIgnoreOrder)
	})
}

func TestVRFv2PlusMigration(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	// Migrate subscription from old coordinator to new coordinator, verify if balances
	// are moved correctly and requests can be made successfully in the subscription in
	// new coordinator
	t.Run("Test migration of Subscription Billing subID", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			2,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		activeSubIdsOldCoordinatorBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
		require.Equal(t, subID, activeSubIdsOldCoordinatorBeforeMigration[0])

		oldSubscriptionBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		//Migration Process
		newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(vrfContracts.BHS.Address())
		require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfKey.VRFKey, newCoordinator)
		require.NoError(t, err, fmt.Errorf("%s, err: %w", vrfcommon.ErrRegisteringProvingKey, err))

		err = newCoordinator.SetConfig(
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.MaxGasLimitCoordinatorConfig,
			*configCopy.VRFv2Plus.General.StalenessSeconds,
			*configCopy.VRFv2Plus.General.GasAfterPaymentCalculation,
			big.NewInt(*configCopy.VRFv2Plus.General.LinkNativeFeedResponse),
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeNativePPM,
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeLinkDiscountPPM,
			*configCopy.VRFv2Plus.General.NativePremiumPercentage,
			*configCopy.VRFv2Plus.General.LinkPremiumPercentage,
		)
		require.NoError(t, err)

		err = newCoordinator.SetLINKAndLINKNativeFeed(vrfContracts.LinkToken.Address(), vrfContracts.MockETHLINKFeed.Address())
		require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)
		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            newCoordinator.Address(),
			FromAddresses:                 nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.VRFKey.Data.ID,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
		}

		_, err = vrfv2plus.CreateVRFV2PlusJob(
			nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

		err = vrfContracts.CoordinatorV2Plus.RegisterMigratableCoordinator(newCoordinator.Address())
		require.NoError(t, err, "error registering migratable coordinator")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		err = vrfContracts.CoordinatorV2Plus.Migrate(subID, newCoordinator.Address())

		require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfContracts.CoordinatorV2Plus.Address(), " to new Coordinator address ", newCoordinator.Address())
		migrationCompletedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForMigrationCompletedEvent(time.Minute * 1)
		require.NoError(t, err, "error waiting for MigrationCompleted event")
		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfContracts)

		oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceAfterMigration, migratedCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		migratedSubscription, err := newCoordinator.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetailsAfterMigration(l, newCoordinator, subID, migratedSubscription)

		//Verify that Coordinators were updated in Consumers
		for _, consumer := range consumers {
			coordinatorAddressInConsumerAfterMigration, err := consumer.GetCoordinator(testcontext.Get(t))
			require.NoError(t, err, "error getting Coordinator from Consumer contract")
			require.Equal(t, newCoordinator.Address(), coordinatorAddressInConsumerAfterMigration.String())
			l.Debug().
				Str("Consumer", consumer.Address()).
				Str("Coordinator", coordinatorAddressInConsumerAfterMigration.String()).
				Msg("Coordinator Address in Consumer After Migration")
		}

		//Verify old and migrated subs
		require.Equal(t, oldSubscriptionBeforeMigration.NativeBalance, migratedSubscription.NativeBalance)
		require.Equal(t, oldSubscriptionBeforeMigration.Balance, migratedSubscription.Balance)
		require.Equal(t, oldSubscriptionBeforeMigration.SubOwner, migratedSubscription.SubOwner)
		require.Equal(t, oldSubscriptionBeforeMigration.Consumers, migratedSubscription.Consumers)

		//Verify that old sub was deleted from old Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		_, err = vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		// If (subscription billing), numActiveSub should be 0 after migration in oldCoordinator
		require.Error(t, err, "error not occurred getting active sub ids. Should occur since it should revert when sub id array is empty")

		activeSubIdsMigratedCoordinator, err := newCoordinator.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsMigratedCoordinator, 1, "Active Sub Ids length is not equal to 1 for Migrated Coordinator after migration")
		require.Equal(t, subID, activeSubIdsMigratedCoordinator[0])

		//Verify that total balances changed for Link and Eth for new and old coordinator
		expectedLinkTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.Balance, migratedCoordinatorLinkTotalBalanceBeforeMigration)
		expectedEthTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.NativeBalance, migratedCoordinatorEthTotalBalanceBeforeMigration)

		expectedLinkTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorLinkTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.Balance)
		expectedEthTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorEthTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.NativeBalance)
		require.Equal(t, 0, expectedLinkTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorEthTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedLinkTotalBalanceForOldCoordinator.Cmp(oldCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForOldCoordinator.Cmp(oldCoordinatorEthTotalBalanceAfterMigration))

		//Verify rand requests fulfills with Link Token billing
		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
			consumers[0],
			newCoordinator,
			vrfKey,
			subID,
			false,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		//Verify rand requests fulfills with Native Token billing
		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
			consumers[1],
			newCoordinator,
			vrfKey,
			subID,
			true,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	})

	// Migrate wrapper subscription from old coordinator to new coordinator, verify if balances
	// are moved correctly and requests can be made successfully in the subscription in
	// new coordinator
	t.Run("Test migration of direct billing using VRFV2PlusWrapper subID", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperEnvironment(
			testcontext.Get(t),
			env,
			chainID,
			&configCopy,
			vrfContracts.LinkToken,
			vrfContracts.MockETHLINKFeed,
			vrfContracts.CoordinatorV2Plus,
			vrfKey.KeyHash,
			1,
		)
		require.NoError(t, err)
		subID := wrapperSubID

		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)

		activeSubIdsOldCoordinatorBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
		activeSubID := activeSubIdsOldCoordinatorBeforeMigration[0]
		require.Equal(t, subID, activeSubID)

		oldSubscriptionBeforeMigration, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		//Migration Process
		newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(vrfContracts.BHS.Address())
		require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfKey.VRFKey, newCoordinator)
		require.NoError(t, err, fmt.Errorf("%s, err: %w", vrfcommon.ErrRegisteringProvingKey, err))

		err = newCoordinator.SetConfig(
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.MaxGasLimitCoordinatorConfig,
			*configCopy.VRFv2Plus.General.StalenessSeconds,
			*configCopy.VRFv2Plus.General.GasAfterPaymentCalculation,
			big.NewInt(*configCopy.VRFv2Plus.General.LinkNativeFeedResponse),
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeNativePPM,
			*configCopy.VRFv2Plus.General.FulfillmentFlatFeeLinkDiscountPPM,
			*configCopy.VRFv2Plus.General.NativePremiumPercentage,
			*configCopy.VRFv2Plus.General.LinkPremiumPercentage,
		)
		require.NoError(t, err)

		err = newCoordinator.SetLINKAndLINKNativeFeed(vrfContracts.LinkToken.Address(), vrfContracts.MockETHLINKFeed.Address())
		require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)
		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            newCoordinator.Address(),
			FromAddresses:                 nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.VRFKey.Data.ID,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
		}

		_, err = vrfv2plus.CreateVRFV2PlusJob(
			nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

		err = vrfContracts.CoordinatorV2Plus.RegisterMigratableCoordinator(newCoordinator.Address())
		require.NoError(t, err, "error registering migratable coordinator")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		// Migrate wrapper's sub using coordinator's migrate method
		err = vrfContracts.CoordinatorV2Plus.Migrate(subID, newCoordinator.Address())

		require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfContracts.CoordinatorV2Plus.Address(), " to new Coordinator address ", newCoordinator.Address())
		migrationCompletedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForMigrationCompletedEvent(time.Minute * 1)
		require.NoError(t, err, "error waiting for MigrationCompleted event")
		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfContracts)

		oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceAfterMigration, migratedCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		migratedSubscription, err := newCoordinator.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetailsAfterMigration(l, newCoordinator, subID, migratedSubscription)

		// Verify that Coordinators were updated in Consumers- Consumer in this case is the VRFV2PlusWrapper
		coordinatorAddressInConsumerAfterMigration, err := wrapperContracts.VRFV2PlusWrapper.Coordinator(testcontext.Get(t))
		require.NoError(t, err, "error getting Coordinator from Consumer contract- VRFV2PlusWrapper")
		require.Equal(t, newCoordinator.Address(), coordinatorAddressInConsumerAfterMigration.String())
		l.Debug().
			Str("Consumer-VRFV2PlusWrapper", wrapperContracts.VRFV2PlusWrapper.Address()).
			Str("Coordinator", coordinatorAddressInConsumerAfterMigration.String()).
			Msg("Coordinator Address in VRFV2PlusWrapper After Migration")

		//Verify old and migrated subs
		require.Equal(t, oldSubscriptionBeforeMigration.NativeBalance, migratedSubscription.NativeBalance)
		require.Equal(t, oldSubscriptionBeforeMigration.Balance, migratedSubscription.Balance)
		require.Equal(t, oldSubscriptionBeforeMigration.SubOwner, migratedSubscription.SubOwner)
		require.Equal(t, oldSubscriptionBeforeMigration.Consumers, migratedSubscription.Consumers)

		//Verify that old sub was deleted from old Coordinator
		_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		_, err = vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		// If (subscription billing) or (direct billing and numActiveSubs is 0 before this test) -> numActiveSub should be 0 after migration in oldCoordinator
		require.Error(t, err, "error not occurred getting active sub ids. Should occur since it should revert when sub id array is empty")

		activeSubIdsMigratedCoordinator, err := newCoordinator.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsMigratedCoordinator, 1, "Active Sub Ids length is not equal to 1 for Migrated Coordinator after migration")
		require.Equal(t, subID, activeSubIdsMigratedCoordinator[0])

		//Verify that total balances changed for Link and Eth for new and old coordinator
		expectedLinkTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.Balance, migratedCoordinatorLinkTotalBalanceBeforeMigration)
		expectedEthTotalBalanceForMigratedCoordinator := new(big.Int).Add(oldSubscriptionBeforeMigration.NativeBalance, migratedCoordinatorEthTotalBalanceBeforeMigration)

		expectedLinkTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorLinkTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.Balance)
		expectedEthTotalBalanceForOldCoordinator := new(big.Int).Sub(oldCoordinatorEthTotalBalanceBeforeMigration, oldSubscriptionBeforeMigration.NativeBalance)
		require.Equal(t, 0, expectedLinkTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForMigratedCoordinator.Cmp(migratedCoordinatorEthTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedLinkTotalBalanceForOldCoordinator.Cmp(oldCoordinatorLinkTotalBalanceAfterMigration))
		require.Equal(t, 0, expectedEthTotalBalanceForOldCoordinator.Cmp(oldCoordinatorEthTotalBalanceAfterMigration))

		// Verify rand requests fulfills with Link Token billing
		isNativeBilling := false
		randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillmentUpgraded(
			wrapperContracts.LoadTestConsumers[0],
			newCoordinator,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, consumerStatus.Fulfilled)

		// Verify rand requests fulfills with Native Token billing
		isNativeBilling = true
		randomWordsFulfilledEvent, err = vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillmentUpgraded(
			wrapperContracts.LoadTestConsumers[0],
			newCoordinator,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		consumerStatus, err = wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, consumerStatus.Fulfilled)
	})
}

func TestVRFV2PlusWithBHS(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	//decrease default span for checking blockhashes for unfulfilled requests
	vrfv2PlusConfig.General.BHSJobWaitBlocks = ptr.Ptr(2)
	vrfv2PlusConfig.General.BHSJobLookBackBlocks = ptr.Ptr(20)

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHS},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	var isNativeBilling = true
	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		t.Skip("Skipped since should be run on-demand on live testnet due to long execution time")
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		_, err = consumers[0].RequestRandomness(
			vrfKey.KeyHash,
			subID,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfKey.KeyHash},
			[]*big.Int{subID},
			[]common.Address{common.HexToAddress(consumers[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2plus.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsRequestedEvent, isNativeBilling)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)
		//Wait at least 256 blocks
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(257), evmClient, &wg, time.Second*260, t)
		wg.Wait()
		require.NoError(t, err)
		err = vrfv2plus.FundSubscriptions(
			env,
			chainID,
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative),
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			subIDs,
		)
		require.NoError(t, err, "error funding subscriptions")
		randomWordsFulfilledEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsFulfilledEvent(
			[]*big.Int{subID},
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			time.Second*30,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfv2plus.LogRandomWordsFulfilledEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsFulfilledEvent, isNativeBilling)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		randRequestBlockHash, err := vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.NoError(t, err, "error getting blockhash for a blocknumber which was stored in BHS contract")

		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		_, err = consumers[0].RequestRandomness(
			vrfKey.KeyHash,
			subID,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfKey.KeyHash},
			[]*big.Int{subID},
			[]common.Address{common.HexToAddress(consumers[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2plus.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsRequestedEvent, isNativeBilling)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		_, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.Error(t, err, "error not occurred when getting blockhash for a blocknumber which was not stored in BHS contract")

		var wg sync.WaitGroup
		wg.Add(1)
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(*configCopy.VRFv2Plus.General.BHSJobWaitBlocks+10), evmClient, &wg, time.Minute*1, t)
		wg.Wait()
		require.NoError(t, err, "error waiting for blocknumber to be")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		var clNodeTxs *client.TransactionsData
		var txHash string
		gom := gomega.NewGomegaWithT(t)
		gom.Eventually(func(g gomega.Gomega) {
			clNodeTxs, _, err = nodeTypeToNodeMap[vrfcommon.BHS].CLNode.API.ReadTransactions()
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting CL Node transactions")
			l.Debug().Int("Number of TXs", len(clNodeTxs.Data)).Msg("BHS Node txs")
			g.Expect(len(clNodeTxs.Data)).Should(gomega.BeNumerically("==", 1), "Expected 1 tx posted by BHS Node, but found %d", len(clNodeTxs.Data))
			txHash = clNodeTxs.Data[0].Attributes.Hash
		}, "2m", "1s").Should(gomega.Succeed())

		require.Equal(t, strings.ToLower(vrfContracts.BHS.Address()), strings.ToLower(clNodeTxs.Data[0].Attributes.To))

		bhsStoreTx, _, err := actions.GetTxByHash(testcontext.Get(t), evmClient, common.HexToHash(txHash))
		require.NoError(t, err, "error getting tx from hash")

		bhsStoreTxInputData, err := actions.DecodeTxInputData(blockhash_store.BlockhashStoreABI, bhsStoreTx.Data())
		l.Info().
			Str("Block Number", bhsStoreTxInputData["n"].(*big.Int).String()).
			Msg("BHS Node's Store Blockhash for Blocknumber Method TX")
		require.Equal(t, randRequestBlockNumber, bhsStoreTxInputData["n"].(*big.Int).Uint64())

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		var randRequestBlockHash [32]byte
		gom.Eventually(func(g gomega.Gomega) {
			randRequestBlockHash, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting blockhash for a blocknumber which was stored in BHS contract")
		}, "2m", "1s").Should(gomega.Succeed())
		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})
}

func TestVRFV2PlusWithBHF(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus

	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	// BHF job config
	config.VRFv2Plus.General.BHFJobWaitBlocks = ptr.Ptr(260)
	config.VRFv2Plus.General.BHFJobLookBackBlocks = ptr.Ptr(500)
	config.VRFv2Plus.General.BHFJobPollPeriod = ptr.Ptr(blockchain.StrDuration{Duration: time.Second * 30})
	config.VRFv2Plus.General.BHFJobRunTimeout = ptr.Ptr(blockchain.StrDuration{Duration: time.Minute * 24})

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHF},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(
		testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err)
	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	var isNativeBilling = true
	t.Run("BHF Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		// t.Skip("Skipped since should be run on-demand on live testnet due to long execution time")
		configCopy := config.MustCopy().(tc.TestConfig)
		// Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		_, err = consumers[0].RequestRandomness(
			vrfKey.KeyHash,
			subID,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfKey.KeyHash},
			[]*big.Int{subID},
			[]common.Address{common.HexToAddress(consumers[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2plus.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsRequestedEvent, isNativeBilling)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)
		//Wait at least 260 blocks
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(260), evmClient, &wg, time.Second*262, t)
		wg.Wait()
		require.NoError(t, err)
		l.Info().Float64("SubscriptionFundingAmountNative", *configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative).
			Float64("SubscriptionFundingAmountLink", *configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink).
			Msg("Funding subscription")
		err = vrfv2plus.FundSubscriptions(
			env,
			chainID,
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative),
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			subIDs,
		)
		require.NoError(t, err, "error funding subscriptions")
		randomWordsFulfilledEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsFulfilledEvent(
			[]*big.Int{subID},
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			time.Minute*2,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfv2plus.LogRandomWordsFulfilledEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsFulfilledEvent, isNativeBilling)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		clNodeTxs, _, err := nodeTypeToNodeMap[vrfcommon.BHF].CLNode.API.ReadTransactions()
		require.NoError(t, err, "error fetching txns from BHF node")
		batchBHSTxFound := false
		for _, tx := range clNodeTxs.Data {
			if strings.EqualFold(tx.Attributes.To, vrfContracts.BatchBHS.Address()) {
				batchBHSTxFound = true
			}
		}
		require.True(t, batchBHSTxFound)

		randRequestBlockHash, err := vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.NoError(t, err, "error getting blockhash for a blocknumber which was stored in BHS contract")

		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})
}

func TestVRFv2PlusReplayAfterTimeout(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	// 1. Add job spec with requestTimeout = 5 seconds
	timeout := time.Second * 5
	config.VRFv2Plus.General.VRFJobRequestTimeout = ptr.Ptr(blockchain.StrDuration{Duration: timeout})
	config.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
	config.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("Timed out request fulfilled after node restart with replay", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = false

		consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			2,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

		// 2. create request but without fulfilment - e.g. simulation failure (insufficient balance in the sub, )
		initialReqRandomWordsRequestedEvent, err := vrfv2plus.RequestRandomnessAndWaitForRequestedEvent(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subID,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for requested event")

		// 3. create new request in a subscription with balance and wait for fulfilment
		fundingLinkAmt := big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink)
		fundingNativeAmt := big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative)
		l.Info().
			Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).
			Int("Number of Subs to create", 1).
			Msg("Creating and funding subscriptions, adding consumers")
		fundedSubIDs, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			chainID,
			fundingLinkAmt,
			fundingNativeAmt,
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			[]contracts.VRFv2PlusLoadTestConsumer{consumers[1]},
			1,
		)
		require.NoError(t, err, "error creating funded sub in replay test")
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[1],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			fundedSubIDs[0],
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		require.True(t, randomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		// 4. wait for the request timeout (1s more) duration
		time.Sleep(timeout + 1*time.Second)

		// 5. fund sub so that node can fulfill request
		err = vrfv2plus.FundSubscriptions(
			env,
			chainID,
			fundingLinkAmt,
			fundingNativeAmt,
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			[]*big.Int{subID},
		)
		require.NoError(t, err, "error funding subs after request timeout")

		// 6. no fulfilment should happen since timeout+1 seconds passed in the job
		pendingReqExists, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
		require.NoError(t, err, "error fetching PendingRequestsExist from coordinator")
		require.True(t, pendingReqExists, "pendingRequest must exist since subID was underfunded till request timeout")

		// 7. remove job and add new job with requestTimeout = 1 hour
		vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]
		require.True(t, exists, "VRF Node does not exist")
		resp, err := vrfNode.CLNode.API.DeleteJob(vrfNode.Job.Data.ID)
		require.NoError(t, err, "error deleting job after timeout")
		require.Equal(t, resp.StatusCode, 204)

		configCopy.VRFv2Plus.General.VRFJobRequestTimeout = ptr.Ptr(blockchain.StrDuration{Duration: time.Duration(time.Hour * 1)})
		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *configCopy.VRFv2Plus.General.VRFJobForwardingAllowed,
			CoordinatorAddress:            vrfContracts.CoordinatorV2Plus.Address(),
			FromAddresses:                 vrfNode.TXKeyAddressStrings,
			EVMChainID:                    fmt.Sprint(chainID),
			MinIncomingConfirmations:      int(*configCopy.VRFv2Plus.General.MinimumConfirmations),
			PublicKey:                     vrfKey.PubKeyCompressed,
			EstimateGasMultiplier:         *configCopy.VRFv2Plus.General.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentEnabled,
			BatchFulfillmentGasMultiplier: *configCopy.VRFv2Plus.General.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    configCopy.VRFv2Plus.General.VRFJobPollPeriod.Duration,
			RequestTimeout:                configCopy.VRFv2Plus.General.VRFJobRequestTimeout.Duration,
			SimulationBlock:               configCopy.VRFv2Plus.General.VRFJobSimulationBlock,
			VRFOwnerConfig:                nil,
		}

		go func() {
			l.Info().
				Msg("Creating VRFV2 Plus Job with higher timeout (1hr)")
			job, err := vrfv2plus.CreateVRFV2PlusJob(
				vrfNode.CLNode.API,
				vrfJobSpecConfig,
			)
			require.NoError(t, err, "error creating job with higher timeout")
			vrfNode.Job = job
		}()

		// 8. Check if initial req in underfunded sub is fulfilled now, since it has been topped up and timeout increased
		l.Info().Str("reqID", initialReqRandomWordsRequestedEvent.RequestId.String()).
			Str("subID", subID.String()).
			Msg("Waiting for initalReqRandomWordsFulfilledEvent")
		initalReqRandomWordsFulfilledEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsFulfilledEvent(
			[]*big.Int{subID},
			[]*big.Int{initialReqRandomWordsRequestedEvent.RequestId},
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err, "error waiting for initial request RandomWordsFulfilledEvent")

		require.NoError(t, err, "error waiting for fulfilment of old req")
		require.False(t, initalReqRandomWordsFulfilledEvent.OnlyPremium, "RandomWordsFulfilled Event's `OnlyPremium` field should be false")
		require.Equal(t, isNativeBilling, initalReqRandomWordsFulfilledEvent.NativePayment, "RandomWordsFulfilled Event's `NativePayment` field should be false")
		require.True(t, initalReqRandomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		// Get request status
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), initalReqRandomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	})
}

func TestVRFv2PlusPendingBlockSimulationAndZeroConfirmationDelays(t *testing.T) {
	t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")
	vrfv2PlusConfig := config.VRFv2Plus
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := env.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	// override config with minConf = 0 and use pending block for simulation
	config.VRFv2Plus.General.MinimumConfirmations = ptr.Ptr[uint16](0)
	config.VRFv2Plus.General.VRFJobSimulationBlock = ptr.Ptr[string]("pending")

	env, vrfContracts, vrfKey, _, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := env.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
		env,
		chainID,
		vrfContracts.CoordinatorV2Plus,
		config,
		vrfContracts.LinkToken,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up new consumers and subs")
	subID := subIDs[0]
	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")
	vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)

	var isNativeBilling = true

	l.Info().Uint16("minimumConfirmationDelay", *config.VRFv2Plus.General.MinimumConfirmations).Msg("Minimum Confirmation Delay")

	// test and assert
	randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
		consumers[0],
		vrfContracts.CoordinatorV2Plus,
		vrfKey,
		subID,
		isNativeBilling,
		config.VRFv2Plus.General,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

	status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	require.NoError(t, err, "error getting rand request status")
	require.True(t, status.Fulfilled)
	l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
}
