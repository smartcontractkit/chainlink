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

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	it_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"
)

func TestVRFv2Plus(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*config.VRFv2Plus.General.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	// default wallet address is used to test Withdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 2
	vrfv2PlusContracts, subIDs, vrfv2PlusData, nodesMap, err := vrfv2plus.SetupVRFV2_5Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		linkToken,
		mockETHLinkFeed,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subID := subIDs[0]

	subscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.CoordinatorV2Plus)

	t.Run("Link Billing", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = false
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData,
			subID,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		require.False(t, randomWordsFulfilledEvent.OnlyPremium, "RandomWordsFulfilled Event's `OnlyPremium` field should be false")
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment, "RandomWordsFulfilled Event's `NativePayment` field should be false")
		require.True(t, randomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.VRFV2PlusConsumer[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
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
		subNativeTokenBalanceBeforeRequest := subscription.NativeBalance

		jobRunsBeforeTest, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData,
			subID,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		require.False(t, randomWordsFulfilledEvent.OnlyPremium)
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment)
		require.True(t, randomWordsFulfilledEvent.Success)
		expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err)
		subBalanceAfterRequest := subscription.NativeBalance
		require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

		jobRuns, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.VRFV2PlusConsumer[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
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
			env,
			&configCopy,
			linkToken,
			mockETHLinkFeed,
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData.KeyHash,
			1,
		)
		require.NoError(t, err)

		t.Run("Link Billing", func(t *testing.T) {
			configCopy := config.MustCopy().(tc.TestConfig)
			testConfig := configCopy.VRFv2Plus.General
			var isNativeBilling = false

			wrapperConsumerJuelsBalanceBeforeRequest, err := linkToken.BalanceOf(testcontext.Get(t), wrapperContracts.LoadTestConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.Balance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.LoadTestConsumers[0],
				vrfv2PlusContracts.CoordinatorV2Plus,
				vrfv2PlusData,
				wrapperSubID,
				isNativeBilling,
				*configCopy.VRFv2Plus.General.MinimumConfirmations,
				*configCopy.VRFv2Plus.General.CallbackGasLimit,
				*configCopy.VRFv2Plus.General.NumberOfWords,
				*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
				*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
				configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.Balance
			require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)

			wrapperConsumerJuelsBalanceAfterRequest, err := linkToken.BalanceOf(testcontext.Get(t), wrapperContracts.LoadTestConsumers[0].Address())
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

			wrapperConsumerBalanceBeforeRequestWei, err := env.EVMClient.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()))
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.NativeBalance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.LoadTestConsumers[0],
				vrfv2PlusContracts.CoordinatorV2Plus,
				vrfv2PlusData,
				wrapperSubID,
				isNativeBilling,
				*configCopy.VRFv2Plus.General.MinimumConfirmations,
				*configCopy.VRFv2Plus.General.CallbackGasLimit,
				*configCopy.VRFv2Plus.General.NumberOfWords,
				*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
				*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
				configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceWei := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.NativeBalance
			require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerWeiBalance := new(big.Int).Sub(wrapperConsumerBalanceBeforeRequestWei, consumerStatus.Paid)

			wrapperConsumerBalanceAfterRequestWei, err := env.EVMClient.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()))
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
		subIDsForCancelling, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*configCopy.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative),
			big.NewFloat(*configCopy.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink),
			linkToken,
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusContracts.VRFV2PlusConsumer,
			1,
		)
		require.NoError(t, err)
		subIDForCancelling := subIDsForCancelling[0]

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceNativeBeforeSubCancelling, err := env.EVMClient.BalanceAt(testcontext.Get(t), testWalletAddress)
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subIDForCancelling.String()).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")
		tx, err := vrfv2PlusContracts.CoordinatorV2Plus.CancelSubscription(subIDForCancelling, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2PlusContracts.CoordinatorV2Plus.WaitForSubscriptionCanceledEvent(subIDForCancelling, time.Second*30)
		require.NoError(t, err, "error waiting for subscription canceled event")

		cancellationTxReceipt, err := env.EVMClient.GetTxReceipt(tx.Hash())
		require.NoError(t, err, "error getting tx cancellation Tx Receipt")

		txGasUsed := new(big.Int).SetUint64(cancellationTxReceipt.GasUsed)
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

		testWalletBalanceNativeAfterSubCancelling, err := env.EVMClient.BalanceAt(testcontext.Get(t), testWalletAddress)
		require.NoError(t, err)

		testWalletBalanceLinkAfterSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForCancelling)
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
		testConfig.SubscriptionFundingAmountNative = ptr.Ptr(float64(0.000000000000000001)) //1 Wei
		testConfig.SubscriptionFundingAmountLink = ptr.Ptr(float64(0.000000000000000001))   //1 Juels

		subIDsForCancelling, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*configCopy.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative),
			big.NewFloat(*configCopy.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink),
			linkToken,
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusContracts.VRFV2PlusConsumer,
			1,
		)
		require.NoError(t, err)

		subIDForCancelling := subIDsForCancelling[0]

		subscriptionForCancelling, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscriptionForCancelling, subIDForCancelling, vrfv2PlusContracts.CoordinatorV2Plus)

		activeSubscriptionIdsBeforeSubCancellation, err := vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err)

		require.True(t, it_utils.BigIntSliceContains(activeSubscriptionIdsBeforeSubCancellation, subIDForCancelling))

		pendingRequestsExist, err := vrfv2PlusContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		randomWordsFulfilledEventTimeout := 5 * time.Second
		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData,
			subIDForCancelling,
			false,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			randomWordsFulfilledEventTimeout,
			l,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData,
			subIDForCancelling,
			true,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			randomWordsFulfilledEventTimeout,
			l,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfv2PlusContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfulfilled rand requests due to low sub balance")

		walletBalanceNativeBeforeSubCancelling, err := env.EVMClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		walletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		subscriptionForCancelling, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subIDForCancelling.String()).
			Str("Returning funds to", defaultWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		tx, err := vrfv2PlusContracts.CoordinatorV2Plus.OwnerCancelSubscription(subIDForCancelling)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2PlusContracts.CoordinatorV2Plus.WaitForSubscriptionCanceledEvent(subIDForCancelling, time.Second*30)
		require.NoError(t, err, "error waiting for subscription canceled event")

		cancellationTxReceipt, err := env.EVMClient.GetTxReceipt(tx.Hash())
		require.NoError(t, err, "error getting tx cancellation Tx Receipt")

		txGasUsed := new(big.Int).SetUint64(cancellationTxReceipt.GasUsed)
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

		walletBalanceNativeAfterSubCancelling, err := env.EVMClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		walletBalanceLinkAfterSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForCancelling)
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

		activeSubscriptionIdsAfterSubCancellation, err := vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error getting active subscription ids")

		require.False(
			t,
			it_utils.BigIntSliceContains(activeSubscriptionIdsAfterSubCancellation, subIDForCancelling),
			"Active subscription ids should not contain sub id after sub cancellation",
		)
	})
	t.Run("Owner Withdraw", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		subIDsForWithdraw, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*configCopy.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative),
			big.NewFloat(*configCopy.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink),
			linkToken,
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusContracts.VRFV2PlusConsumer,
			1,
		)
		require.NoError(t, err)
		subIDForWithdraw := subIDsForWithdraw[0]

		fulfilledEventLink, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData,
			subIDForWithdraw,
			false,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)
		require.NoError(t, err)

		fulfilledEventNative, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData,
			subIDForWithdraw,
			true,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)
		require.NoError(t, err)
		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceNativeBeforeWithdraw, err := env.EVMClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		defaultWalletBalanceLinkBeforeWithdraw, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfv2PlusContracts.CoordinatorV2Plus.Withdraw(
			common.HexToAddress(defaultWalletAddress),
		)
		require.NoError(t, err, "error withdrawing LINK from coordinator to default wallet")
		amountToWithdrawNative := fulfilledEventNative.Payment

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawNative.String()).
			Msg("Invoking Oracle Withdraw for Native")

		err = vrfv2PlusContracts.CoordinatorV2Plus.WithdrawNative(
			common.HexToAddress(defaultWalletAddress),
		)
		require.NoError(t, err, "error withdrawing Native tokens from coordinator to default wallet")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		defaultWalletBalanceNativeAfterWithdraw, err := env.EVMClient.BalanceAt(testcontext.Get(t), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		defaultWalletBalanceLinkAfterWithdraw, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		//not possible to verify exact amount of Native/LINK returned as defaultWallet is used in other tests in parallel which might affect the balance
		require.Equal(t, 1, defaultWalletBalanceNativeAfterWithdraw.Cmp(defaultWalletBalanceNativeBeforeWithdraw), "Native funds were not returned after oracle withdraw native")
		require.Equal(t, 1, defaultWalletBalanceLinkAfterWithdraw.Cmp(defaultWalletBalanceLinkBeforeWithdraw), "LINK funds were not returned after oracle withdraw")
	})
}

func TestVRFv2PlusMultipleSendingKeys(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*config.VRFv2Plus.General.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	numberOfTxKeysToCreate := 2
	vrfv2PlusContracts, subIDs, vrfv2PlusData, nodesMap, err := vrfv2plus.SetupVRFV2_5Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		linkToken,
		mockETHLinkFeed,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subID := subIDs[0]

	subscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.CoordinatorV2Plus)

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = true
		txKeys, _, err := nodesMap[vrfcommon.VRF].CLNode.API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, numberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < numberOfTxKeysToCreate+1; i++ {
			randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
				vrfv2PlusContracts.VRFV2PlusConsumer[0],
				vrfv2PlusContracts.CoordinatorV2Plus,
				vrfv2PlusData,
				subID,
				isNativeBilling,
				*configCopy.VRFv2Plus.General.MinimumConfirmations,
				*configCopy.VRFv2Plus.General.CallbackGasLimit,
				*configCopy.VRFv2Plus.General.NumberOfWords,
				*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
				*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
				configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			//todo - move TransactionByHash to EVMClient in CTF
			fulfillmentTx, _, err := actions.GetTxByHash(testcontext.Get(t), env.EVMClient, randomWordsFulfilledEvent.Raw.TxHash)
			require.NoError(t, err, "error getting tx from hash")
			fulfillmentTxFromAddress, err := actions.GetTxFromAddress(fulfillmentTx)
			require.NoError(t, err, "error getting tx from address")
			fulfillmentTxFromAddresses = append(fulfillmentTxFromAddresses, fulfillmentTxFromAddress)
		}
		require.Equal(t, numberOfTxKeysToCreate+1, len(fulfillmentTxFromAddresses))
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
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")
	env.ParallelTransactions(true)

	mockETHLinkFeedAddress, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*config.VRFv2Plus.General.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkAddress, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	vrfv2PlusContracts, subIDs, vrfv2PlusData, nodesMap, err := vrfv2plus.SetupVRFV2_5Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		linkAddress,
		mockETHLinkFeedAddress,
		0,
		2,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	// Migrate subscription from old coordinator to new coordinator, verify if balances
	// are moved correctly and requests can be made successfully in the subscription in
	// new coordinator
	t.Run("Test migration of Subscription Billing subID", func(t *testing.T) {
		subID := subIDs[0]

		subscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.CoordinatorV2Plus)

		activeSubIdsOldCoordinatorBeforeMigration, err := vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
		require.Equal(t, subID, activeSubIdsOldCoordinatorBeforeMigration[0])

		oldSubscriptionBeforeMigration, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		//Migration Process
		newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(vrfv2PlusContracts.BHS.Address())
		require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfv2PlusData.VRFKey, newCoordinator)
		require.NoError(t, err, fmt.Errorf("%s, err: %w", vrfcommon.ErrRegisteringProvingKey, err))

		vrfv2PlusConfig := config.VRFv2Plus.General
		err = newCoordinator.SetConfig(
			*vrfv2PlusConfig.MinimumConfirmations,
			*vrfv2PlusConfig.MaxGasLimitCoordinatorConfig,
			*vrfv2PlusConfig.StalenessSeconds,
			*vrfv2PlusConfig.GasAfterPaymentCalculation,
			big.NewInt(*vrfv2PlusConfig.LinkNativeFeedResponse),
			*vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
			*vrfv2PlusConfig.FulfillmentFlatFeeLinkDiscountPPM,
			*vrfv2PlusConfig.NativePremiumPercentage,
			*vrfv2PlusConfig.LinkPremiumPercentage,
		)
		require.NoError(t, err)

		err = newCoordinator.SetLINKAndLINKNativeFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
		require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)
		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *vrfv2PlusConfig.VRFJobForwardingAllowed,
			CoordinatorAddress:            newCoordinator.Address(),
			FromAddresses:                 nodesMap[vrfcommon.VRF].TXKeyAddressStrings,
			EVMChainID:                    env.EVMClient.GetChainID().String(),
			MinIncomingConfirmations:      int(*vrfv2PlusConfig.MinimumConfirmations),
			PublicKey:                     vrfv2PlusData.VRFKey.Data.ID,
			EstimateGasMultiplier:         *vrfv2PlusConfig.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       *vrfv2PlusConfig.VRFJobBatchFulfillmentEnabled,
			BatchFulfillmentGasMultiplier: *vrfv2PlusConfig.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    vrfv2PlusConfig.VRFJobPollPeriod.Duration,
			RequestTimeout:                vrfv2PlusConfig.VRFJobRequestTimeout.Duration,
		}

		_, err = vrfv2plus.CreateVRFV2PlusJob(
			nodesMap[vrfcommon.VRF].CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

		err = vrfv2PlusContracts.CoordinatorV2Plus.RegisterMigratableCoordinator(newCoordinator.Address())
		require.NoError(t, err, "error registering migratable coordinator")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfv2PlusContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		err = vrfv2PlusContracts.CoordinatorV2Plus.Migrate(subID, newCoordinator.Address())

		require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfv2PlusContracts.CoordinatorV2Plus.Address(), " to new Coordinator address ", newCoordinator.Address())
		migrationCompletedEvent, err := vrfv2PlusContracts.CoordinatorV2Plus.WaitForMigrationCompletedEvent(time.Minute * 1)
		require.NoError(t, err, "error waiting for MigrationCompleted event")
		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfv2PlusContracts)

		oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfv2PlusContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceAfterMigration, migratedCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		migratedSubscription, err := newCoordinator.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetailsAfterMigration(l, newCoordinator, subID, migratedSubscription)

		//Verify that Coordinators were updated in Consumers
		for _, consumer := range vrfv2PlusContracts.VRFV2PlusConsumer {
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
		require.Equal(t, oldSubscriptionBeforeMigration.Owner, migratedSubscription.Owner)
		require.Equal(t, oldSubscriptionBeforeMigration.Consumers, migratedSubscription.Consumers)

		//Verify that old sub was deleted from old Coordinator
		_, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		_, err = vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
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
			vrfv2PlusContracts.VRFV2PlusConsumer[0],
			newCoordinator,
			vrfv2PlusData,
			subID,
			false,
			*config.VRFv2Plus.General.MinimumConfirmations,
			*config.VRFv2Plus.General.CallbackGasLimit,
			*config.VRFv2Plus.General.NumberOfWords,
			*config.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*config.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			config.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		//Verify rand requests fulfills with Native Token billing
		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
			vrfv2PlusContracts.VRFV2PlusConsumer[1],
			newCoordinator,
			vrfv2PlusData,
			subID,
			true,
			*config.VRFv2Plus.General.MinimumConfirmations,
			*config.VRFv2Plus.General.CallbackGasLimit,
			*config.VRFv2Plus.General.NumberOfWords,
			*config.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*config.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			config.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
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
			env,
			&configCopy,
			linkAddress,
			mockETHLinkFeedAddress,
			vrfv2PlusContracts.CoordinatorV2Plus,
			vrfv2PlusData.KeyHash,
			1,
		)
		require.NoError(t, err)
		subID := wrapperSubID

		subscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.CoordinatorV2Plus)

		activeSubIdsOldCoordinatorBeforeMigration, err := vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error occurred getting active sub ids")
		require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
		activeSubID := activeSubIdsOldCoordinatorBeforeMigration[0]
		require.Equal(t, subID, activeSubID)

		oldSubscriptionBeforeMigration, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		//Migration Process
		newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(vrfv2PlusContracts.BHS.Address())
		require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfv2PlusData.VRFKey, newCoordinator)
		require.NoError(t, err, fmt.Errorf("%s, err: %w", vrfcommon.ErrRegisteringProvingKey, err))

		vrfv2PlusConfig := config.VRFv2Plus.General
		err = newCoordinator.SetConfig(
			*vrfv2PlusConfig.MinimumConfirmations,
			*vrfv2PlusConfig.MaxGasLimitCoordinatorConfig,
			*vrfv2PlusConfig.StalenessSeconds,
			*vrfv2PlusConfig.GasAfterPaymentCalculation,
			big.NewInt(*vrfv2PlusConfig.LinkNativeFeedResponse),
			*vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
			*vrfv2PlusConfig.FulfillmentFlatFeeLinkDiscountPPM,
			*vrfv2PlusConfig.NativePremiumPercentage,
			*vrfv2PlusConfig.LinkPremiumPercentage,
		)
		require.NoError(t, err)

		err = newCoordinator.SetLINKAndLINKNativeFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
		require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)
		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
			ForwardingAllowed:             *vrfv2PlusConfig.VRFJobForwardingAllowed,
			CoordinatorAddress:            newCoordinator.Address(),
			FromAddresses:                 nodesMap[vrfcommon.VRF].TXKeyAddressStrings,
			EVMChainID:                    env.EVMClient.GetChainID().String(),
			MinIncomingConfirmations:      int(*vrfv2PlusConfig.MinimumConfirmations),
			PublicKey:                     vrfv2PlusData.VRFKey.Data.ID,
			EstimateGasMultiplier:         *vrfv2PlusConfig.VRFJobEstimateGasMultiplier,
			BatchFulfillmentEnabled:       *vrfv2PlusConfig.VRFJobBatchFulfillmentEnabled,
			BatchFulfillmentGasMultiplier: *vrfv2PlusConfig.VRFJobBatchFulfillmentGasMultiplier,
			PollPeriod:                    vrfv2PlusConfig.VRFJobPollPeriod.Duration,
			RequestTimeout:                vrfv2PlusConfig.VRFJobRequestTimeout.Duration,
		}

		_, err = vrfv2plus.CreateVRFV2PlusJob(
			nodesMap[vrfcommon.VRF].CLNode.API,
			vrfJobSpecConfig,
		)
		require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

		err = vrfv2PlusContracts.CoordinatorV2Plus.RegisterMigratableCoordinator(newCoordinator.Address())
		require.NoError(t, err, "error registering migratable coordinator")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfv2PlusContracts.CoordinatorV2Plus)
		require.NoError(t, err)

		migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
		require.NoError(t, err)

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		// Migrate sub using VRFV2PlusWrapper's migrate method
		err = wrapperContracts.VRFV2PlusWrapper.Migrate(common.HexToAddress(newCoordinator.Address()))

		require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfv2PlusContracts.CoordinatorV2Plus.Address(), " to new Coordinator address ", newCoordinator.Address())
		migrationCompletedEvent, err := vrfv2PlusContracts.CoordinatorV2Plus.WaitForMigrationCompletedEvent(time.Minute * 1)
		require.NoError(t, err, "error waiting for MigrationCompleted event")
		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfv2PlusContracts)

		oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfv2PlusContracts.CoordinatorV2Plus)
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
		require.Equal(t, oldSubscriptionBeforeMigration.Owner, migratedSubscription.Owner)
		require.Equal(t, oldSubscriptionBeforeMigration.Consumers, migratedSubscription.Consumers)

		//Verify that old sub was deleted from old Coordinator
		_, err = vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		_, err = vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
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
			vrfv2PlusData,
			subID,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
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
			vrfv2PlusData,
			subID,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
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
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(2).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := env.ContractDeployer.DeployVRFMockETHLINKFeed(big.NewInt(*config.VRFv2Plus.General.LinkNativeFeedResponse))

	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	//Underfund Subscription
	config.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0.000000000000000001)) // 1 Juel

	//decrease default span for checking blockhashes for unfulfilled requests
	config.VRFv2Plus.General.BHSJobWaitBlocks = ptr.Ptr(2)
	config.VRFv2Plus.General.BHSJobLookBackBlocks = ptr.Ptr(20)

	numberOfTxKeysToCreate := 0
	vrfContracts, subIDs, vrfKeyData, nodesMap, err := vrfv2plus.SetupVRFV2_5Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHS},
		&config,
		linkToken,
		mockETHLinkFeed,
		numberOfTxKeysToCreate,
		1,
		2,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	var isNativeBilling = true
	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		t.Skip("Skipped since should be run on-demand on live testnet due to long execution time")

		subID := subIDs[0]

		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)

		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		_, err = vrfContracts.VRFV2PlusConsumer[0].RequestRandomness(
			vrfKeyData.KeyHash,
			subID,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfKeyData.KeyHash},
			[]*big.Int{subID},
			[]common.Address{common.HexToAddress(vrfContracts.VRFV2PlusConsumer[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2plus.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsRequestedEvent, isNativeBilling)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)
		//Wait at least 256 blocks
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(257), env.EVMClient, &wg, time.Second*260, t)
		wg.Wait()
		require.NoError(t, err)
		err = vrfv2plus.FundSubscriptions(
			env,
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative),
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink),
			linkToken,
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
		status, err := vrfContracts.VRFV2PlusConsumer[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
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
		subID := subIDs[1]

		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)

		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		_, err = vrfContracts.VRFV2PlusConsumer[0].RequestRandomness(
			vrfKeyData.KeyHash,
			subID,
			*configCopy.VRFv2Plus.General.MinimumConfirmations,
			*configCopy.VRFv2Plus.General.CallbackGasLimit,
			isNativeBilling,
			*configCopy.VRFv2Plus.General.NumberOfWords,
			*configCopy.VRFv2Plus.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfContracts.CoordinatorV2Plus.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfKeyData.KeyHash},
			[]*big.Int{subID},
			[]common.Address{common.HexToAddress(vrfContracts.VRFV2PlusConsumer[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2plus.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2Plus, randomWordsRequestedEvent, isNativeBilling)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		_, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.Error(t, err, "error not occurred when getting blockhash for a blocknumber which was not stored in BHS contract")

		var wg sync.WaitGroup
		wg.Add(1)
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(*config.VRFv2Plus.General.BHSJobWaitBlocks+10), env.EVMClient, &wg, time.Minute*1, t)
		wg.Wait()
		require.NoError(t, err, "error waiting for blocknumber to be")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		var clNodeTxs *client.TransactionsData
		var txHash string
		gom := gomega.NewGomegaWithT(t)
		gom.Eventually(func(g gomega.Gomega) {
			clNodeTxs, _, err = nodesMap[vrfcommon.BHS].CLNode.API.ReadTransactions()
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting CL Node transactions")
			l.Debug().Int("Number of TXs", len(clNodeTxs.Data)).Msg("BHS Node txs")
			g.Expect(len(clNodeTxs.Data)).Should(gomega.BeNumerically("==", 1), "Expected 1 tx posted by BHS Node, but found %d", len(clNodeTxs.Data))
			txHash = clNodeTxs.Data[0].Attributes.Hash
		}, "2m", "1s").Should(gomega.Succeed())

		require.Equal(t, strings.ToLower(vrfContracts.BHS.Address()), strings.ToLower(clNodeTxs.Data[0].Attributes.To))

		bhsStoreTx, _, err := actions.GetTxByHash(testcontext.Get(t), env.EVMClient, common.HexToHash(txHash))
		require.NoError(t, err, "error getting tx from hash")

		bhsStoreTxInputData, err := actions.DecodeTxInputData(blockhash_store.BlockhashStoreABI, bhsStoreTx.Data())
		l.Info().
			Str("Block Number", bhsStoreTxInputData["n"].(*big.Int).String()).
			Msg("BHS Node's Store Blockhash for Blocknumber Method TX")
		require.Equal(t, randRequestBlockNumber, bhsStoreTxInputData["n"].(*big.Int).Uint64())

		err = env.EVMClient.WaitForEvents()
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

func TestVRFv2PlusPendingBlockSimulationAndZeroConfirmationDelays(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	if err != nil {
		t.Fatal(err)
	}

	// override config with minConf = 0 and use pending block for simulation
	config.VRFv2Plus.General.MinimumConfirmations = ptr.Ptr[uint16](0)
	config.VRFv2Plus.General.VRFJobSimulationBlock = ptr.Ptr[string]("pending")

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*config.VRFv2Plus.General.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	numberOfTxKeysToCreate := 2
	vrfv2PlusContracts, subIDs, vrfv2PlusData, nodesMap, err := vrfv2plus.SetupVRFV2_5Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		linkToken,
		mockETHLinkFeed,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subID := subIDs[0]

	subscription, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.CoordinatorV2Plus)

	var isNativeBilling = true

	jobRunsBeforeTest, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
	require.NoError(t, err, "error reading job runs")

	l.Info().Uint16("minimumConfirmationDelay", *config.VRFv2Plus.General.MinimumConfirmations).Msg("Minimum Confirmation Delay")

	// test and assert
	randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
		vrfv2PlusContracts.VRFV2PlusConsumer[0],
		vrfv2PlusContracts.CoordinatorV2Plus,
		vrfv2PlusData,
		subID,
		isNativeBilling,
		*config.VRFv2Plus.General.MinimumConfirmations,
		*config.VRFv2Plus.General.CallbackGasLimit,
		*config.VRFv2Plus.General.NumberOfWords,
		*config.VRFv2Plus.General.RandomnessRequestCountPerRequest,
		*config.VRFv2Plus.General.RandomnessRequestCountPerRequestDeviation,
		config.VRFv2Plus.General.RandomWordsFulfilledEventTimeout.Duration,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

	jobRuns, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
	require.NoError(t, err, "error reading job runs")
	require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

	status, err := vrfv2PlusContracts.VRFV2PlusConsumer[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	require.NoError(t, err, "error getting rand request status")
	require.True(t, status.Fulfilled)
	l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
}
