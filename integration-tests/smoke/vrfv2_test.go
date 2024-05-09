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

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
)

func TestVRFv2Basic(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	vrfv2Config := config.VRFv2
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
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

	testEnv, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")
	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")

	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("Request Randomness", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDsForRequestRandomness, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForRequestRandomness := subIDsForRequestRandomness[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscription, subIDForRequestRandomness, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)

		subBalanceBeforeRequest := subscription.Balance

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForRequestRandomness,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})

	t.Run("CL Node VRF Job Runs", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDsForJobRuns, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")

		subIDForJobRuns := subIDsForJobRuns[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForJobRuns)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscription, subIDForJobRuns, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForJobRuns...)

		jobRunsBeforeTest, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		_, err = vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForJobRuns,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		jobRuns, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))
	})

	t.Run("Direct Funding (VRFV2Wrapper)", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		wrapperContracts, wrapperSubID, err := vrfv2.SetupVRFV2WrapperEnvironment(
			testcontext.Get(t),
			testEnv,
			chainID,
			&configCopy,
			vrfContracts.LinkToken,
			vrfContracts.MockETHLINKFeed,
			vrfContracts.CoordinatorV2,
			vrfKey.KeyHash,
			1,
		)
		require.NoError(t, err)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, *wrapperSubID)

		wrapperConsumer := wrapperContracts.LoadTestConsumers[0]

		wrapperConsumerJuelsBalanceBeforeRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
		require.NoError(t, err, "Error getting wrapper consumer balance")

		wrapperSubscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceBeforeRequest := wrapperSubscription.Balance

		// Request Randomness and wait for fulfillment event
		randomWordsFulfilledEvent, err := vrfv2.DirectFundingRequestRandomnessAndWaitForFulfillment(
			l,
			wrapperConsumer,
			vrfContracts.CoordinatorV2,
			*wrapperSubID,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err, "Error requesting randomness and waiting for fulfilment")

		// Check wrapper subscription balance
		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		wrapperSubscription, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceAfterRequest := wrapperSubscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		// Check status of randomness request within the wrapper consumer contract
		consumerStatus, err := wrapperConsumer.GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "Error getting randomness request status")
		require.True(t, consumerStatus.Fulfilled)

		// Check wrapper consumer LINK balance
		expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)
		wrapperConsumerJuelsBalanceAfterRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
		require.NoError(t, err, "Error getting wrapper consumer balance")
		require.Equal(t, expectedWrapperConsumerJuelsBalance, wrapperConsumerJuelsBalanceAfterRequest)

		// Check random word count
		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
		for _, w := range consumerStatus.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}

		l.Info().
			Str("Consumer Balance Before Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceBeforeRequest).Link()).
			Str("Consumer Balance After Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceAfterRequest).Link()).
			Bool("Fulfilment Status", consumerStatus.Fulfilled).
			Str("Paid by Consumer Contract (Link)", (*commonassets.Link)(consumerStatus.Paid).Link()).
			Str("Paid by Coordinator Sub (Link)", (*commonassets.Link)(randomWordsFulfilledEvent.Payment).Link()).
			Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
			Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
			Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
			Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
			Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
			Msg("Random Words Fulfilment Details For Link Billing")
	})

	t.Run("Oracle Withdraw", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		consumers, subIDsForOracleWithDraw, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")

		subIDForOracleWithdraw := subIDsForOracleWithDraw[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOracleWithdraw)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscription, subIDForOracleWithdraw, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForOracleWithDraw...)

		fulfilledEventLink, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForOracleWithdraw,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err)

		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceLinkBeforeOracleWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfContracts.CoordinatorV2.OracleWithdraw(common.HexToAddress(defaultWalletAddress), amountToWithdrawLink)
		require.NoError(t, err, "Error withdrawing LINK from coordinator to default wallet")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		defaultWalletBalanceLinkAfterOracleWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		require.Equal(
			t,
			1,
			defaultWalletBalanceLinkAfterOracleWithdraw.Cmp(defaultWalletBalanceLinkBeforeOracleWithdraw),
			"LINK funds were not returned after oracle withdraw",
		)
	})

	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		_, subIDsForCancelling, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForCancelling := subIDsForCancelling[0]
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscription, subIDForCancelling, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForCancelling...)

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForCancelling).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")

		tx, err := vrfContracts.CoordinatorV2.CancelSubscription(subIDForCancelling, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfContracts.CoordinatorV2.WaitForSubscriptionCanceledEvent([]uint64{subIDForCancelling}, time.Second*30)
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
			Str("Returned Subscription Amount Link", subscriptionCanceledEvent.Amount.String()).
			Uint64("SubID", subscriptionCanceledEvent.SubId).
			Str("Returned to", subscriptionCanceledEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceLink, subscriptionCanceledEvent.Amount, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		testWalletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

		subFundsReturnedLinkActual := new(big.Int).Sub(testWalletBalanceLinkAfterSubCancelling, testWalletBalanceLinkBeforeSubCancelling)

		l.Info().
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Msg("Sub funds returned")

		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")
	})

	t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		// Underfund subscription to force fulfillments to fail
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))

		consumers, subIDsForOwnerCancelling, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForOwnerCancelling := subIDsForOwnerCancelling[0]
		subscriptionForCancelling, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscriptionForCancelling, subIDForOwnerCancelling, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForOwnerCancelling...)

		// No GetActiveSubscriptionIds function available - skipping check

		pendingRequestsExist, err := vrfContracts.CoordinatorV2.PendingRequestsExist(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		// Request randomness - should fail due to underfunded subscription
		randomWordsFulfilledEventTimeout := 5 * time.Second
		_, err = vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForOwnerCancelling,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			randomWordsFulfilledEventTimeout,
		)
		require.Error(t, err, "Error should occur while waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfContracts.CoordinatorV2.PendingRequestsExist(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfilfulled requests due to low sub balance")

		walletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		subscriptionForCancelling, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOwnerCancelling)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForOwnerCancelling).
			Str("Returning funds to", defaultWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")

		// Call OwnerCancelSubscription
		tx, err := vrfContracts.CoordinatorV2.OwnerCancelSubscription(subIDForOwnerCancelling)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfContracts.CoordinatorV2.WaitForSubscriptionCanceledEvent([]uint64{subIDForOwnerCancelling}, time.Second*30)
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
			Str("Returned Subscription Amount Link", subscriptionCanceledEvent.Amount.String()).
			Uint64("SubID", subscriptionCanceledEvent.SubId).
			Str("Returned to", subscriptionCanceledEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceLink, subscriptionCanceledEvent.Amount, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		walletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		// Verify that subscription was deleted from Coordinator contract
		_, err = vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForOwnerCancelling)
		l.Info().
			Str("Expected error message", err.Error())
		require.Error(t, err, "Error did not occur when fetching deleted subscription from the Coordinator after owner cancelation")

		subFundsReturnedLinkActual := new(big.Int).Sub(walletBalanceLinkAfterSubCancelling, walletBalanceLinkBeforeSubCancelling)
		l.Info().
			Str("Wallet Balance Before Owner Cancelation", walletBalanceLinkBeforeSubCancelling.String()).
			Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
			Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
			Str("Sub Balance - Link", subBalanceLink.String()).
			Str("Wallet Balance After Owner Cancelation", walletBalanceLinkAfterSubCancelling.String()).
			Msg("Sub funds returned")

		require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")

		// Again, there is no GetActiveSubscriptionIds method on the v2 Coordinator contract, so we can't double check that the cancelled
		// subID is no longer in the list of active subs
	})
}

func TestVRFv2MultipleSendingKeys(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	if err != nil {
		t.Fatal(err)
	}
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	vrfv2Config := config.VRFv2
	cleanupFn := func() {
		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
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

	testEnv, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDsForMultipleSendingKeys, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForMultipleSendingKeys := subIDsForMultipleSendingKeys[0]
		subscriptionForMultipleSendingKeys, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForMultipleSendingKeys)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscriptionForMultipleSendingKeys, subIDForMultipleSendingKeys, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForMultipleSendingKeys...)

		txKeys, _, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, newEnvConfig.NumberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < newEnvConfig.NumberOfTxKeysToCreate+1; i++ {
			randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
				l,
				consumers[0],
				vrfContracts.CoordinatorV2,
				subIDForMultipleSendingKeys,
				vrfKey,
				*configCopy.VRFv2.General.MinimumConfirmations,
				*configCopy.VRFv2.General.CallbackGasLimit,
				*configCopy.VRFv2.General.NumberOfWords,
				*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
				*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
				configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
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

func TestVRFOwner(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	vrfv2Config := config.VRFv2
	cleanupFn := func() {
		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            true,
		UseTestCoordinator:     true,
	}

	testEnv, vrfContracts, vrfKey, _, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("Request Randomness With Force-Fulfill", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		consumers, subIDsForForceFulfill, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForForceFulfill := subIDsForForceFulfill[0]
		subscriptionForMultipleSendingKeys, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForForceFulfill)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscriptionForMultipleSendingKeys, subIDForForceFulfill, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForForceFulfill...)

		vrfCoordinatorOwner, err := vrfContracts.CoordinatorV2.GetOwner(testcontext.Get(t))
		require.NoError(t, err)
		require.Equal(t, vrfContracts.VRFOwner.Address(), vrfCoordinatorOwner.String())

		err = vrfContracts.LinkToken.Transfer(
			consumers[0].Address(),
			conversions.EtherToWei(big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink)),
		)
		require.NoError(t, err, "error transferring link to consumer contract")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		consumerLinkBalance, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), consumers[0].Address())
		require.NoError(t, err, "error getting consumer link balance")
		l.Info().
			Str("Balance", conversions.WeiToEther(consumerLinkBalance).String()).
			Str("Consumer", consumers[0].Address()).
			Msg("Consumer Link Balance")

		err = vrfContracts.MockETHLINKFeed.SetBlockTimestampDeduction(big.NewInt(3))
		require.NoError(t, err)
		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		// test and assert
		_, randFulfilledEvent, _, err := vrfv2.RequestRandomnessWithForceFulfillAndWaitForFulfillment(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			vrfContracts.VRFOwner,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			conversions.EtherToWei(big.NewFloat(5)),
			common.HexToAddress(vrfContracts.LinkToken.Address()),
			time.Minute*2,
		)
		require.NoError(t, err, "error requesting randomness with force-fulfillment and waiting for fulfilment")
		require.Equal(t, 0, randFulfilledEvent.Payment.Cmp(big.NewInt(0)), "Forced Fulfilled Randomness's Payment should be 0")

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}

		coordinatorConfig, err := vrfContracts.CoordinatorV2.GetConfig(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator config")

		coordinatorFeeConfig, err := vrfContracts.CoordinatorV2.GetFeeConfig(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator fee config")

		coordinatorFallbackWeiPerUnitLinkConfig, err := vrfContracts.CoordinatorV2.GetFallbackWeiPerUnitLink(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator FallbackWeiPerUnitLink")

		require.Equal(t, *configCopy.VRFv2.General.StalenessSeconds, coordinatorConfig.StalenessSeconds)
		require.Equal(t, *configCopy.VRFv2.General.GasAfterPaymentCalculation, coordinatorConfig.GasAfterPaymentCalculation)
		require.Equal(t, *configCopy.VRFv2.General.MinimumConfirmations, coordinatorConfig.MinimumRequestConfirmations)
		require.Equal(t, *configCopy.VRFv2.General.FulfillmentFlatFeeLinkPPMTier1, coordinatorFeeConfig.FulfillmentFlatFeeLinkPPMTier1)
		require.Equal(t, *configCopy.VRFv2.General.ReqsForTier2, coordinatorFeeConfig.ReqsForTier2.Int64())
		require.Equal(t, *configCopy.VRFv2.General.FallbackWeiPerUnitLink, coordinatorFallbackWeiPerUnitLinkConfig.Int64())
	})
}

func TestVRFV2WithBHS(t *testing.T) {
	t.Parallel()
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")
	vrfv2Config := config.VRFv2
	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "Getting EVM client shouldn't fail")
		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	//decrease default span for checking blockhashes for unfulfilled requests
	vrfv2Config.General.BHSJobWaitBlocks = ptr.Ptr(2)
	vrfv2Config.General.BHSJobLookBackBlocks = ptr.Ptr(20)

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHS},
		NumberOfTxKeysToCreate: 0,
		UseVRFOwner:            false,
		UseTestCoordinator:     false,
	}

	testEnv, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, config, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "Error setting up VRFV2 universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		t.Skip("Skipped since should be run on-demand on live testnet due to long execution time")
		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)

		//Underfund Subscription
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		consumers, subIDsForBHS, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForBHS := subIDsForBHS[0]
		subscriptionForBHS, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForBHS)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscriptionForBHS, subIDForBHS, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForBHS...)

		randomWordsRequestedEvent, err := vrfv2.RequestRandomness(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForBHS,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
		)
		require.NoError(t, err, "error requesting randomness")

		vrfv2.LogRandomnessRequestedEvent(l, vrfContracts.CoordinatorV2, randomWordsRequestedEvent)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)
		//Wait at least 256 blocks
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(257), evmClient, &wg, time.Second*260, t)
		wg.Wait()
		require.NoError(t, err)
		err = vrfv2.FundSubscriptions(testEnv, chainID, big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink), vrfContracts.LinkToken, vrfContracts.CoordinatorV2, subIDsForBHS)
		require.NoError(t, err, "error funding subscriptions")
		randomWordsFulfilledEvent, err := vrfContracts.CoordinatorV2.WaitForRandomWordsFulfilledEvent(
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			time.Second*30,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfv2.LogRandomWordsFulfilledEvent(l, vrfContracts.CoordinatorV2, randomWordsFulfilledEvent)
		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))

		consumers, subIDsForBHS, err := vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForBHS := subIDsForBHS[0]
		subscriptionForBHS, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForBHS)
		require.NoError(t, err, "error getting subscription information")
		vrfv2.LogSubDetails(l, subscriptionForBHS, subIDForBHS, vrfContracts.CoordinatorV2)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForBHS...)

		randomWordsRequestedEvent, err := vrfv2.RequestRandomness(
			l,
			consumers[0],
			vrfContracts.CoordinatorV2,
			subIDForBHS,
			vrfKey,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
		)
		require.NoError(t, err, "error requesting randomness")

		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber

		_, err = vrfContracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.Error(t, err, "error not occurred when getting blockhash for a blocknumber which was not stored in BHS contract")

		var wg sync.WaitGroup
		wg.Add(1)
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(*configCopy.VRFv2.General.BHSJobWaitBlocks), evmClient, &wg, time.Minute*1, t)
		wg.Wait()
		require.NoError(t, err, "error waiting for blocknumber to be")

		err = evmClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)
		metrics, err := consumers[0].GetLoadTestMetrics(testcontext.Get(t))
		require.Equal(t, 0, metrics.RequestCount.Cmp(big.NewInt(1)))
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)))

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
