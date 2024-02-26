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
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")

	useVRFOwner := false
	useTestCoordinator := false
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

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*config.VRFv2.General.LinkNativeFeedResponse))
	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 1
	vrfv2Contracts, subIDs, vrfv2KeyData, nodesMap, err := vrfv2.SetupVRFV2Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		useVRFOwner,
		useTestCoordinator,
		linkToken,
		mockETHLinkFeed,
		defaultWalletAddress,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2 env")

	subID := subIDs[0]

	subscription, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2.LogSubDetails(l, subscription, subID, vrfv2Contracts.CoordinatorV2)

	t.Run("Request Randomness", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			vrfv2Contracts.VRFV2Consumer[0],
			vrfv2Contracts.CoordinatorV2,
			subID,
			vrfv2KeyData,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := nodesMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodesMap[vrfcommon.VRF].Job.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2Contracts.VRFV2Consumer[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *config.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})

	t.Run("Direct Funding (VRFV2Wrapper)", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		wrapperContracts, wrapperSubID, err := vrfv2.SetupVRFV2WrapperEnvironment(
			env,
			&configCopy,
			linkToken,
			mockETHLinkFeed,
			vrfv2Contracts.CoordinatorV2,
			vrfv2KeyData.KeyHash,
			1,
		)
		require.NoError(t, err)
		wrapperConsumer := wrapperContracts.LoadTestConsumers[0]

		wrapperConsumerJuelsBalanceBeforeRequest, err := linkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
		require.NoError(t, err, "Error getting wrapper consumer balance")

		wrapperSubscription, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceBeforeRequest := wrapperSubscription.Balance

		// Request Randomness and wait for fulfillment event
		randomWordsFulfilledEvent, err := vrfv2.DirectFundingRequestRandomnessAndWaitForFulfillment(
			l,
			wrapperConsumer,
			vrfv2Contracts.CoordinatorV2,
			*wrapperSubID,
			vrfv2KeyData,
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
		wrapperSubscription, err = vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceAfterRequest := wrapperSubscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		// Check status of randomness request within the wrapper consumer contract
		consumerStatus, err := wrapperConsumer.GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "Error getting randomness request status")
		require.True(t, consumerStatus.Fulfilled)

		// Check wrapper consumer LINK balance
		expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)
		wrapperConsumerJuelsBalanceAfterRequest, err := linkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
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
		subIDsForOracleWithDraw, err := vrfv2.CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink),
			linkToken,
			vrfv2Contracts.CoordinatorV2,
			vrfv2Contracts.VRFV2Consumer,
			1,
		)
		require.NoError(t, err)

		subIDForOracleWithdraw := subIDsForOracleWithDraw[0]

		fulfilledEventLink, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			vrfv2Contracts.VRFV2Consumer[0],
			vrfv2Contracts.CoordinatorV2,
			subIDForOracleWithdraw,
			vrfv2KeyData,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			configCopy.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
		)
		require.NoError(t, err)

		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceLinkBeforeOracleWithdraw, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfv2Contracts.CoordinatorV2.OracleWithdraw(common.HexToAddress(defaultWalletAddress), amountToWithdrawLink)
		require.NoError(t, err, "Error withdrawing LINK from coordinator to default wallet")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		defaultWalletBalanceLinkAfterOracleWithdraw, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
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
		subIDsForCancelling, err := vrfv2.CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink),
			linkToken,
			vrfv2Contracts.CoordinatorV2,
			vrfv2Contracts.VRFV2Consumer,
			1,
		)
		require.NoError(t, err)
		subIDForCancelling := subIDsForCancelling[0]

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForCancelling).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")

		tx, err := vrfv2Contracts.CoordinatorV2.CancelSubscription(subIDForCancelling, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2Contracts.CoordinatorV2.WaitForSubscriptionCanceledEvent([]uint64{subIDForCancelling}, time.Second*30)
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
			Str("Returned Subscription Amount Link", subscriptionCanceledEvent.Amount.String()).
			Uint64("SubID", subscriptionCanceledEvent.SubId).
			Str("Returned to", subscriptionCanceledEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceLink, subscriptionCanceledEvent.Amount, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		testWalletBalanceLinkAfterSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		//Verify that sub was deleted from Coordinator
		_, err = vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
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
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0.000000000000000001)) // 1 Juel

		subIDsForCancelling, err := vrfv2.CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink),
			linkToken,
			vrfv2Contracts.CoordinatorV2,
			vrfv2Contracts.VRFV2Consumer,
			1,
		)
		require.NoError(t, err)

		subIDForCancelling := subIDsForCancelling[0]

		subscriptionForCancelling, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "Error getting subscription information")

		vrfv2.LogSubDetails(l, subscriptionForCancelling, subIDForCancelling, vrfv2Contracts.CoordinatorV2)

		// No GetActiveSubscriptionIds function available - skipping check

		pendingRequestsExist, err := vrfv2Contracts.CoordinatorV2.PendingRequestsExist(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		// Request randomness - should fail due to underfunded subscription
		randomWordsFulfilledEventTimeout := 5 * time.Second
		_, err = vrfv2.RequestRandomnessAndWaitForFulfillment(
			l,
			vrfv2Contracts.VRFV2Consumer[0],
			vrfv2Contracts.CoordinatorV2,
			subIDForCancelling,
			vrfv2KeyData,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			randomWordsFulfilledEventTimeout,
		)
		require.Error(t, err, "Error should occur while waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfv2Contracts.CoordinatorV2.PendingRequestsExist(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfilfulled requests due to low sub balance")

		walletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		subscriptionForCancelling, err = vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForCancelling).
			Str("Returning funds to", defaultWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")

		// Call OwnerCancelSubscription
		tx, err := vrfv2Contracts.CoordinatorV2.OwnerCancelSubscription(subIDForCancelling)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2Contracts.CoordinatorV2.WaitForSubscriptionCanceledEvent([]uint64{subIDForCancelling}, time.Second*30)
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
			Str("Returned Subscription Amount Link", subscriptionCanceledEvent.Amount.String()).
			Uint64("SubID", subscriptionCanceledEvent.SubId).
			Str("Returned to", subscriptionCanceledEvent.To.String()).
			Msg("Subscription Canceled Event")

		require.Equal(t, subBalanceLink, subscriptionCanceledEvent.Amount, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")

		walletBalanceLinkAfterSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		// Verify that subscription was deleted from Coordinator contract
		_, err = vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subIDForCancelling)
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
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	if err != nil {
		t.Fatal(err)
	}

	useVRFOwner := false
	useTestCoordinator := false

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&config).
		WithTestInstance(t).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*config.VRFv2.General.LinkNativeFeedResponse))
	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 2
	vrfv2Contracts, subIDs, vrfv2KeyData, nodesMap, err := vrfv2.SetupVRFV2Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		useVRFOwner,
		useTestCoordinator,
		linkToken,
		mockETHLinkFeed,
		defaultWalletAddress,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2 env")

	subID := subIDs[0]

	subscription, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2.LogSubDetails(l, subscription, subID, vrfv2Contracts.CoordinatorV2)

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		txKeys, _, err := nodesMap[vrfcommon.VRF].CLNode.API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, numberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < numberOfTxKeysToCreate+1; i++ {
			randomWordsFulfilledEvent, err := vrfv2.RequestRandomnessAndWaitForFulfillment(
				l,
				vrfv2Contracts.VRFV2Consumer[0],
				vrfv2Contracts.CoordinatorV2,
				subID,
				vrfv2KeyData,
				*config.VRFv2.General.MinimumConfirmations,
				*config.VRFv2.General.CallbackGasLimit,
				*config.VRFv2.General.NumberOfWords,
				*config.VRFv2.General.RandomnessRequestCountPerRequest,
				*config.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
				config.VRFv2.General.RandomWordsFulfilledEventTimeout.Duration,
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

func TestVRFOwner(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")

	useVRFOwner := true
	useTestCoordinator := true
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

	mockETHLinkFeed, err := env.ContractDeployer.DeployVRFMockETHLINKFeed(big.NewInt(*config.VRFv2.General.LinkNativeFeedResponse))

	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 1
	vrfv2Contracts, subIDs, vrfv2Data, _, err := vrfv2.SetupVRFV2Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF},
		&config,
		useVRFOwner,
		useTestCoordinator,
		linkToken,
		mockETHLinkFeed,
		defaultWalletAddress,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2 env")

	subID := subIDs[0]

	subscription, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2.LogSubDetails(l, subscription, subID, vrfv2Contracts.CoordinatorV2)

	t.Run("Request Randomness With Force-Fulfill", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)

		vrfCoordinatorOwner, err := vrfv2Contracts.CoordinatorV2.GetOwner(testcontext.Get(t))
		require.NoError(t, err)
		require.Equal(t, vrfv2Contracts.VRFOwner.Address(), vrfCoordinatorOwner.String())

		err = linkToken.Transfer(
			vrfv2Contracts.VRFV2Consumer[0].Address(),
			conversions.EtherToWei(big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink)),
		)
		require.NoError(t, err, "error transferring link to consumer contract")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		consumerLinkBalance, err := linkToken.BalanceOf(testcontext.Get(t), vrfv2Contracts.VRFV2Consumer[0].Address())
		require.NoError(t, err, "error getting consumer link balance")
		l.Info().
			Str("Balance", conversions.WeiToEther(consumerLinkBalance).String()).
			Str("Consumer", vrfv2Contracts.VRFV2Consumer[0].Address()).
			Msg("Consumer Link Balance")

		err = mockETHLinkFeed.SetBlockTimestampDeduction(big.NewInt(3))
		require.NoError(t, err)
		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)

		// test and assert
		_, randFulfilledEvent, _, err := vrfv2.RequestRandomnessWithForceFulfillAndWaitForFulfillment(
			l,
			vrfv2Contracts.VRFV2Consumer[0],
			vrfv2Contracts.CoordinatorV2,
			vrfv2Contracts.VRFOwner,
			vrfv2Data,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
			conversions.EtherToWei(big.NewFloat(5)),
			common.HexToAddress(linkToken.Address()),
			time.Minute*2,
		)
		require.NoError(t, err, "error requesting randomness with force-fulfillment and waiting for fulfilment")
		require.Equal(t, 0, randFulfilledEvent.Payment.Cmp(big.NewInt(0)), "Forced Fulfilled Randomness's Payment should be 0")

		status, err := vrfv2Contracts.VRFV2Consumer[0].GetRequestStatus(testcontext.Get(t), randFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}

		coordinatorConfig, err := vrfv2Contracts.CoordinatorV2.GetConfig(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator config")

		coordinatorFeeConfig, err := vrfv2Contracts.CoordinatorV2.GetFeeConfig(testcontext.Get(t))
		require.NoError(t, err, "error getting coordinator fee config")

		coordinatorFallbackWeiPerUnitLinkConfig, err := vrfv2Contracts.CoordinatorV2.GetFallbackWeiPerUnitLink(testcontext.Get(t))
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
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.VRFv2)
	require.NoError(t, err, "Error getting config")

	useVRFOwner := true
	useTestCoordinator := true
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

	mockETHLinkFeed, err := env.ContractDeployer.DeployVRFMockETHLINKFeed(big.NewInt(*config.VRFv2.General.LinkNativeFeedResponse))

	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	//Underfund Subscription
	config.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0.000000000000000001)) // 1 Juel

	//decrease default span for checking blockhashes for unfulfilled requests
	config.VRFv2.General.BHSJobWaitBlocks = ptr.Ptr(2)
	config.VRFv2.General.BHSJobLookBackBlocks = ptr.Ptr(20)

	numberOfTxKeysToCreate := 0
	vrfv2Contracts, subIDs, vrfv2KeyData, nodesMap, err := vrfv2.SetupVRFV2Environment(
		env,
		[]vrfcommon.VRFNodeType{vrfcommon.VRF, vrfcommon.BHS},
		&config,
		useVRFOwner,
		useTestCoordinator,
		linkToken,
		mockETHLinkFeed,
		defaultWalletAddress,
		numberOfTxKeysToCreate,
		1,
		1,
		l,
	)
	require.NoError(t, err, "error setting up VRF v2 env")

	subID := subIDs[0]

	subscription, err := vrfv2Contracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2.LogSubDetails(l, subscription, subID, vrfv2Contracts.CoordinatorV2)

	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		t.Skip("Skipped since should be run on-demand on live testnet due to long execution time")
		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		_, err := vrfv2Contracts.VRFV2Consumer[0].RequestRandomness(
			vrfv2KeyData.KeyHash,
			subID,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfv2Contracts.CoordinatorV2.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfv2KeyData.KeyHash},
			[]uint64{subID},
			[]common.Address{common.HexToAddress(vrfv2Contracts.VRFV2Consumer[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2.LogRandomnessRequestedEvent(l, vrfv2Contracts.CoordinatorV2, randomWordsRequestedEvent)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber
		var wg sync.WaitGroup
		wg.Add(1)
		//Wait at least 256 blocks
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(257), env.EVMClient, &wg, time.Second*260, t)
		wg.Wait()
		require.NoError(t, err)
		err = vrfv2.FundSubscriptions(env, big.NewFloat(*configCopy.VRFv2.General.SubscriptionFundingAmountLink), linkToken, vrfv2Contracts.CoordinatorV2, subIDs)
		require.NoError(t, err, "error funding subscriptions")
		randomWordsFulfilledEvent, err := vrfv2Contracts.CoordinatorV2.WaitForRandomWordsFulfilledEvent(
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			time.Second*30,
		)
		require.NoError(t, err, "error waiting for randomness fulfilled event")
		vrfv2.LogRandomWordsFulfilledEvent(l, vrfv2Contracts.CoordinatorV2, randomWordsFulfilledEvent)
		status, err := vrfv2Contracts.VRFV2Consumer[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		//BHS node should fill in blockhashes into BHS contract depending on the waitBlocks and lookBackBlocks settings
		configCopy := config.MustCopy().(tc.TestConfig)
		_, err := vrfv2Contracts.VRFV2Consumer[0].RequestRandomness(
			vrfv2KeyData.KeyHash,
			subID,
			*configCopy.VRFv2.General.MinimumConfirmations,
			*configCopy.VRFv2.General.CallbackGasLimit,
			*configCopy.VRFv2.General.NumberOfWords,
			*configCopy.VRFv2.General.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsRequestedEvent, err := vrfv2Contracts.CoordinatorV2.WaitForRandomWordsRequestedEvent(
			[][32]byte{vrfv2KeyData.KeyHash},
			[]uint64{subID},
			[]common.Address{common.HexToAddress(vrfv2Contracts.VRFV2Consumer[0].Address())},
			time.Minute*1,
		)
		require.NoError(t, err, "error waiting for randomness requested event")
		vrfv2.LogRandomnessRequestedEvent(l, vrfv2Contracts.CoordinatorV2, randomWordsRequestedEvent)
		randRequestBlockNumber := randomWordsRequestedEvent.Raw.BlockNumber

		_, err = vrfv2Contracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
		require.Error(t, err, "error not occurred when getting blockhash for a blocknumber which was not stored in BHS contract")

		var wg sync.WaitGroup
		wg.Add(1)
		_, err = actions.WaitForBlockNumberToBe(randRequestBlockNumber+uint64(*config.VRFv2.General.BHSJobWaitBlocks), env.EVMClient, &wg, time.Minute*1, t)
		wg.Wait()
		require.NoError(t, err, "error waiting for blocknumber to be")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfcommon.ErrWaitTXsComplete)
		metrics, err := vrfv2Contracts.VRFV2Consumer[0].GetLoadTestMetrics(testcontext.Get(t))
		require.Equal(t, 0, metrics.RequestCount.Cmp(big.NewInt(1)))
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)))

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

		require.Equal(t, strings.ToLower(vrfv2Contracts.BHS.Address()), strings.ToLower(clNodeTxs.Data[0].Attributes.To))

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
			randRequestBlockHash, err = vrfv2Contracts.BHS.GetBlockHash(testcontext.Get(t), big.NewInt(int64(randRequestBlockNumber)))
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "error getting blockhash for a blocknumber which was stored in BHS contract")
		}, "2m", "1s").Should(gomega.Succeed())
		l.Info().
			Str("Randomness Request's Blockhash", randomWordsRequestedEvent.Raw.BlockHash.String()).
			Str("Block Hash stored by BHS contract", fmt.Sprintf("0x%x", randRequestBlockHash)).
			Msg("BHS Contract's stored Blockhash for Randomness Request")
		require.Equal(t, 0, randomWordsRequestedEvent.Raw.BlockHash.Cmp(randRequestBlockHash))
	})
}
