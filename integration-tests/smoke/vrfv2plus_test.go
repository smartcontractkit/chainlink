package smoke

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestVRFv2Plus(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	var vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig
	err := envconfig.Process("VRFV2PLUS", &vrfv2PlusConfig)
	require.NoError(t, err)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(big.NewFloat(vrfv2PlusConfig.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(vrfv2PlusConfig.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	vrfv2PlusContracts, subIDs, vrfv2PlusData, err := vrfv2plus.SetupVRFV2_5Environment(env, vrfv2PlusConfig, linkToken, mockETHLinkFeed, defaultWalletAddress, 1, 1, l)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subID := subIDs[0]

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.Coordinator)

	t.Run("Link Billing", func(t *testing.T) {
		testConfig := vrfv2PlusConfig
		var isNativeBilling = false
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subID,
			isNativeBilling,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, testConfig.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})
	t.Run("Native Billing", func(t *testing.T) {
		testConfig := vrfv2PlusConfig
		var isNativeBilling = true
		subNativeTokenBalanceBeforeRequest := subscription.NativeBalance

		jobRunsBeforeTest, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subID,
			isNativeBilling,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err)
		subBalanceAfterRequest := subscription.NativeBalance
		require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, testConfig.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})
	t.Run("Direct Funding (VRFV2PlusWrapper)", func(t *testing.T) {
		testConfig := vrfv2PlusConfig
		wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperEnvironment(
			env,
			testConfig,
			linkToken,
			mockETHLinkFeed,
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData.KeyHash,
			1,
		)
		require.NoError(t, err)

		t.Run("Link Billing", func(t *testing.T) {
			testConfig := vrfv2PlusConfig
			var isNativeBilling = false

			wrapperConsumerJuelsBalanceBeforeRequest, err := linkToken.BalanceOf(context.Background(), wrapperContracts.LoadTestConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.Balance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.LoadTestConsumers[0],
				vrfv2PlusContracts.Coordinator,
				vrfv2PlusData,
				wrapperSubID,
				isNativeBilling,
				testConfig,
				testConfig.RandomWordsFulfilledEventTimeout,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.Balance
			require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)

			wrapperConsumerJuelsBalanceAfterRequest, err := linkToken.BalanceOf(context.Background(), wrapperContracts.LoadTestConsumers[0].Address())
			require.NoError(t, err, "error getting wrapper consumer balance")
			require.Equal(t, expectedWrapperConsumerJuelsBalance, wrapperConsumerJuelsBalanceAfterRequest)

			//todo: uncomment when VRF-651 will be fixed
			//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
			vrfv2plus.LogFulfillmentDetailsLinkBilling(l, wrapperConsumerJuelsBalanceBeforeRequest, wrapperConsumerJuelsBalanceAfterRequest, consumerStatus, randomWordsFulfilledEvent)

			require.Equal(t, testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
			for _, w := range consumerStatus.RandomWords {
				l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
				require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
			}
		})
		t.Run("Native Billing", func(t *testing.T) {
			testConfig := vrfv2PlusConfig
			var isNativeBilling = true

			wrapperConsumerBalanceBeforeRequestWei, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()))
			require.NoError(t, err, "error getting wrapper consumer balance")

			wrapperSubscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceBeforeRequest := wrapperSubscription.NativeBalance

			randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
				wrapperContracts.LoadTestConsumers[0],
				vrfv2PlusContracts.Coordinator,
				vrfv2PlusData,
				wrapperSubID,
				isNativeBilling,
				testConfig,
				testConfig.RandomWordsFulfilledEventTimeout,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			expectedSubBalanceWei := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
			wrapperSubscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), wrapperSubID)
			require.NoError(t, err, "error getting subscription information")
			subBalanceAfterRequest := wrapperSubscription.NativeBalance
			require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

			consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
			require.NoError(t, err, "error getting rand request status")
			require.True(t, consumerStatus.Fulfilled)

			expectedWrapperConsumerWeiBalance := new(big.Int).Sub(wrapperConsumerBalanceBeforeRequestWei, consumerStatus.Paid)

			wrapperConsumerBalanceAfterRequestWei, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()))
			require.NoError(t, err, "error getting wrapper consumer balance")
			require.Equal(t, expectedWrapperConsumerWeiBalance, wrapperConsumerBalanceAfterRequestWei)

			//todo: uncomment when VRF-651 will be fixed
			//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
			vrfv2plus.LogFulfillmentDetailsNativeBilling(l, wrapperConsumerBalanceBeforeRequestWei, wrapperConsumerBalanceAfterRequestWei, consumerStatus, randomWordsFulfilledEvent)

			require.Equal(t, testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
			for _, w := range consumerStatus.RandomWords {
				l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
				require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
			}
		})
	})
	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		testConfig := vrfv2PlusConfig
		subIDsForCancelling, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			testConfig,
			linkToken,
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusContracts.LoadTestConsumers,
			1,
		)
		require.NoError(t, err)
		subIDForCancelling := subIDsForCancelling[0]

		testWalletAddress, err := actions.GenerateWallet()

		testWalletBalanceNativeBeforeSubCancelling, err := env.EVMClient.BalanceAt(context.Background(), testWalletAddress)
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(context.Background(), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subIDForCancelling.String()).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")
		tx, err := vrfv2PlusContracts.Coordinator.CancelSubscription(subIDForCancelling, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2PlusContracts.Coordinator.WaitForSubscriptionCanceledEvent(subIDForCancelling, time.Second*30)
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

		testWalletBalanceNativeAfterSubCancelling, err := env.EVMClient.BalanceAt(context.Background(), testWalletAddress)
		require.NoError(t, err)

		testWalletBalanceLinkAfterSubCancelling, err := linkToken.BalanceOf(context.Background(), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subIDForCancelling)
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
		testConfig := vrfv2PlusConfig
		//underfund subs in order rand fulfillments to fail
		testConfig.SubscriptionFundingAmountNative = float64(0.000000000000000001) //1 Wei
		testConfig.SubscriptionFundingAmountLink = float64(0.000000000000000001)   //1 Juels

		subIDsForCancelling, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			testConfig,
			linkToken,
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusContracts.LoadTestConsumers,
			1,
		)
		require.NoError(t, err)

		subIDForCancelling := subIDsForCancelling[0]

		subscriptionForCancelling, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		vrfv2plus.LogSubDetails(l, subscriptionForCancelling, subIDForCancelling, vrfv2PlusContracts.Coordinator)

		activeSubscriptionIdsBeforeSubCancellation, err := vrfv2PlusContracts.Coordinator.GetActiveSubscriptionIds(context.Background(), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err)

		require.True(t, utils.BigIntSliceContains(activeSubscriptionIdsBeforeSubCancellation, subIDForCancelling))

		pendingRequestsExist, err := vrfv2PlusContracts.Coordinator.PendingRequestsExist(context.Background(), subIDForCancelling)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subIDForCancelling,
			false,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			5*time.Second,
			l,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subIDForCancelling,
			true,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)

		require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfv2PlusContracts.Coordinator.PendingRequestsExist(context.Background(), subIDForCancelling)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfulfilled rand requests due to low sub balance")

		walletBalanceNativeBeforeSubCancelling, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		walletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(context.Background(), defaultWalletAddress)
		require.NoError(t, err)

		subscriptionForCancelling, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance
		subBalanceNative := subscriptionForCancelling.NativeBalance
		l.Info().
			Str("Subscription Amount Native", subBalanceNative.String()).
			Str("Subscription Amount Link", subBalanceLink.String()).
			Str("Returning funds from SubID", subIDForCancelling.String()).
			Str("Returning funds to", defaultWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		tx, err := vrfv2PlusContracts.Coordinator.OwnerCancelSubscription(subIDForCancelling)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2PlusContracts.Coordinator.WaitForSubscriptionCanceledEvent(subIDForCancelling, time.Second*30)
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

		walletBalanceNativeAfterSubCancelling, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		walletBalanceLinkAfterSubCancelling, err := linkToken.BalanceOf(context.Background(), defaultWalletAddress)
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subIDForCancelling)
		fmt.Println("err", err)
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
		//todo - as defaultWallet is used in other tests in parallel which might affect the balance
		//require.Equal(t, 1, walletBalanceNativeAfterSubCancelling.Cmp(walletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")

		//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
		//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

		activeSubscriptionIdsAfterSubCancellation, err := vrfv2PlusContracts.Coordinator.GetActiveSubscriptionIds(context.Background(), big.NewInt(0), big.NewInt(0))
		require.NoError(t, err, "error getting active subscription ids")

		require.False(
			t,
			utils.BigIntSliceContains(activeSubscriptionIdsAfterSubCancellation, subIDForCancelling),
			"Active subscription ids should not contain sub id after sub cancellation",
		)
	})
	t.Run("Oracle Withdraw", func(t *testing.T) {
		testConfig := vrfv2PlusConfig
		subIDsForOracleWithDraw, err := vrfv2plus.CreateFundSubsAndAddConsumers(
			env,
			testConfig,
			linkToken,
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusContracts.LoadTestConsumers,
			1,
		)
		require.NoError(t, err)
		subIDForOracleWithdraw := subIDsForOracleWithDraw[0]

		fulfilledEventLink, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subIDForOracleWithdraw,
			false,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err)

		fulfilledEventNative, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subIDForOracleWithdraw,
			true,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err)
		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceNativeBeforeOracleWithdraw, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		defaultWalletBalanceLinkBeforeOracleWithdraw, err := linkToken.BalanceOf(context.Background(), defaultWalletAddress)
		require.NoError(t, err)

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfv2PlusContracts.Coordinator.OracleWithdraw(
			common.HexToAddress(defaultWalletAddress),
			amountToWithdrawLink,
		)
		require.NoError(t, err, "error withdrawing LINK from coordinator to default wallet")
		amountToWithdrawNative := fulfilledEventNative.Payment

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawNative.String()).
			Msg("Invoking Oracle Withdraw for Native")

		err = vrfv2PlusContracts.Coordinator.OracleWithdrawNative(
			common.HexToAddress(defaultWalletAddress),
			amountToWithdrawNative,
		)
		require.NoError(t, err, "error withdrawing Native tokens from coordinator to default wallet")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

		defaultWalletBalanceNativeAfterOracleWithdraw, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(defaultWalletAddress))
		require.NoError(t, err)

		defaultWalletBalanceLinkAfterOracleWithdraw, err := linkToken.BalanceOf(context.Background(), defaultWalletAddress)
		require.NoError(t, err)

		//not possible to verify exact amount of Native/LINK returned as defaultWallet is used in other tests in parallel which might affect the balance
		require.Equal(t, 1, defaultWalletBalanceNativeAfterOracleWithdraw.Cmp(defaultWalletBalanceNativeBeforeOracleWithdraw), "Native funds were not returned after oracle withdraw native")
		require.Equal(t, 1, defaultWalletBalanceLinkAfterOracleWithdraw.Cmp(defaultWalletBalanceLinkBeforeOracleWithdraw), "LINK funds were not returned after oracle withdraw")
	})
}

func TestVRFv2PlusMigration(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	var vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig
	err := envconfig.Process("VRFV2PLUS", &vrfv2PlusConfig)
	require.NoError(t, err)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(big.NewFloat(vrfv2PlusConfig.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")
	env.ParallelTransactions(true)

	mockETHLinkFeedAddress, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(vrfv2PlusConfig.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkAddress, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	nativeTokenPrimaryKeyAddress, err := env.ClCluster.NodeAPIs()[0].PrimaryEthAddress()
	require.NoError(t, err, "error getting primary eth address")

	vrfv2PlusContracts, subIDs, vrfv2PlusData, err := vrfv2plus.SetupVRFV2_5Environment(env, vrfv2PlusConfig, linkAddress, mockETHLinkFeedAddress, nativeTokenPrimaryKeyAddress, 2, 1, l)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subID := subIDs[0]

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.Coordinator)

	activeSubIdsOldCoordinatorBeforeMigration, err := vrfv2PlusContracts.Coordinator.GetActiveSubscriptionIds(context.Background(), big.NewInt(0), big.NewInt(0))
	require.NoError(t, err, "error occurred getting active sub ids")
	require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 1, "Active Sub Ids length is not equal to 1")
	require.Equal(t, subID, activeSubIdsOldCoordinatorBeforeMigration[0])

	oldSubscriptionBeforeMigration, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	//Migration Process
	newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(vrfv2PlusContracts.BHS.Address())
	require.NoError(t, err, "error deploying VRF CoordinatorV2PlusUpgradedVersion")

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfv2PlusData.VRFKey, vrfv2PlusData.PrimaryEthAddress, newCoordinator)
	require.NoError(t, err, errors.Wrap(err, vrfv2plus.ErrRegisteringProvingKey))

	err = newCoordinator.SetConfig(
		vrfv2PlusConfig.MinimumConfirmations,
		vrfv2PlusConfig.MaxGasLimitCoordinatorConfig,
		vrfv2PlusConfig.StalenessSeconds,
		vrfv2PlusConfig.GasAfterPaymentCalculation,
		big.NewInt(vrfv2PlusConfig.LinkNativeFeedResponse),
		vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionFeeConfig{
			FulfillmentFlatFeeLinkPPM:   vrfv2PlusConfig.FulfillmentFlatFeeLinkPPM,
			FulfillmentFlatFeeNativePPM: vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
		},
	)

	err = newCoordinator.SetLINKAndLINKNativeFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
	require.NoError(t, err, vrfv2plus.ErrSetLinkNativeLinkFeed)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	_, err = vrfv2plus.CreateVRFV2PlusJob(
		env.ClCluster.NodeAPIs()[0],
		newCoordinator.Address(),
		vrfv2PlusData.PrimaryEthAddress,
		vrfv2PlusData.VRFKey.Data.ID,
		vrfv2PlusData.ChainID.String(),
		vrfv2PlusConfig.MinimumConfirmations,
	)
	require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

	err = vrfv2PlusContracts.Coordinator.RegisterMigratableCoordinator(newCoordinator.Address())
	require.NoError(t, err, "error registering migratable coordinator")

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	oldCoordinatorLinkTotalBalanceBeforeMigration, oldCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfv2PlusContracts.Coordinator)
	require.NoError(t, err)

	migratedCoordinatorLinkTotalBalanceBeforeMigration, migratedCoordinatorEthTotalBalanceBeforeMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
	require.NoError(t, err)

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	err = vrfv2PlusContracts.Coordinator.Migrate(subID, newCoordinator.Address())
	require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", vrfv2PlusContracts.Coordinator.Address(), " to new Coordinator address ", newCoordinator.Address())
	migrationCompletedEvent, err := vrfv2PlusContracts.Coordinator.WaitForMigrationCompletedEvent(time.Minute * 1)
	require.NoError(t, err, "error waiting for MigrationCompleted event")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	vrfv2plus.LogMigrationCompletedEvent(l, migrationCompletedEvent, vrfv2PlusContracts)

	oldCoordinatorLinkTotalBalanceAfterMigration, oldCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetCoordinatorTotalBalance(vrfv2PlusContracts.Coordinator)
	require.NoError(t, err)

	migratedCoordinatorLinkTotalBalanceAfterMigration, migratedCoordinatorEthTotalBalanceAfterMigration, err := vrfv2plus.GetUpgradedCoordinatorTotalBalance(newCoordinator)
	require.NoError(t, err)

	migratedSubscription, err := newCoordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetailsAfterMigration(l, newCoordinator, subID, migratedSubscription)

	//Verify that Coordinators were updated in Consumers
	for _, consumer := range vrfv2PlusContracts.LoadTestConsumers {
		coordinatorAddressInConsumerAfterMigration, err := consumer.GetCoordinator(context.Background())
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
	_, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

	_, err = vrfv2PlusContracts.Coordinator.GetActiveSubscriptionIds(context.Background(), big.NewInt(0), big.NewInt(0))
	require.Error(t, err, "error not occurred getting active sub ids. Should occur since it should revert when sub id array is empty")

	activeSubIdsMigratedCoordinator, err := newCoordinator.GetActiveSubscriptionIds(context.Background(), big.NewInt(0), big.NewInt(0))
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
		vrfv2PlusContracts.LoadTestConsumers[0],
		newCoordinator,
		vrfv2PlusData,
		subID,
		false,
		vrfv2PlusConfig,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

	//Verify rand requests fulfills with Native Token billing
	_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
		vrfv2PlusContracts.LoadTestConsumers[1],
		newCoordinator,
		vrfv2PlusData,
		subID,
		true,
		vrfv2PlusConfig,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

}
