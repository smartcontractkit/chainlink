package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_config"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestVRFv2Basic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	network, err := actions.EthereumNetworkConfigFromEnvOrDefault(l)
	require.NoError(t, err, "Error building ethereum network config")

	var vrfv2Config vrfv2_config.VRFV2Config
	err = envconfig.Process("VRFV2", &vrfv2Config)
	require.NoError(t, err)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(vrfv2Config.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(vrfv2Config.LinkNativeFeedResponse))
	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 1
	vrfv2Contracts, subIDs, vrfv2Data, err := vrfv2_actions.SetupVRFV2Environment(
		env,
		vrfv2Config,
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

	subscription, err := vrfv2Contracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2_actions.LogSubDetails(l, subscription, subID, vrfv2Contracts.Coordinator)

	t.Run("Request Randomness", func(t *testing.T) {
		testConfig := vrfv2Config
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2Data.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2_actions.RequestRandomnessAndWaitForFulfillment(
			vrfv2Contracts.LoadTestConsumers[0],
			vrfv2Contracts.Coordinator,
			vrfv2Data,
			subID,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2Contracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2Data.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2Contracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, testConfig.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})

	t.Run("Direct Funding (VRFV2Wrapper)", func(t *testing.T) {
		testConfig := vrfv2Config
		wrapperContracts, wrapperSubID, err := vrfv2_actions.SetupVRFV2WrapperEnvironment(
			env,
			testConfig,
			linkToken,
			mockETHLinkFeed,
			vrfv2Contracts.Coordinator,
			vrfv2Data.KeyHash,
			1,
		)
		require.NoError(t, err)
		wrapperConsumer := wrapperContracts.LoadTestConsumers[0]

		wrapperConsumerJuelsBalanceBeforeRequest, err := linkToken.BalanceOf(testcontext.Get(t), wrapperConsumer.Address())
		require.NoError(t, err, "Error getting wrapper consumer balance")

		wrapperSubscription, err := vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), *wrapperSubID)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceBeforeRequest := wrapperSubscription.Balance

		// Request Randomness and wait for fulfillment event
		randomWordsFulfilledEvent, err := vrfv2_actions.DirectFundingRequestRandomnessAndWaitForFulfillment(
			wrapperConsumer,
			vrfv2Contracts.Coordinator,
			vrfv2Data,
			*wrapperSubID,
			vrfv2Config.RandomnessRequestCountPerRequest,
			vrfv2Config,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err, "Error requesting randomness and waiting for fulfilment")

		// Check wrapper subscription balance
		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		wrapperSubscription, err = vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), *wrapperSubID)
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
		require.Equal(t, testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
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
		testConfig := vrfv2Config
		subIDsForOracleWithDraw, err := vrfv2_actions.CreateFundSubsAndAddConsumers(
			env,
			testConfig,
			linkToken,
			vrfv2Contracts.Coordinator,
			vrfv2Contracts.LoadTestConsumers,
			1,
		)
		require.NoError(t, err)

		subIDForOracleWithdraw := subIDsForOracleWithDraw[0]

		fulfilledEventLink, err := vrfv2_actions.RequestRandomnessAndWaitForFulfillment(
			vrfv2Contracts.LoadTestConsumers[0],
			vrfv2Contracts.Coordinator,
			vrfv2Data,
			subIDForOracleWithdraw,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			testConfig.RandomWordsFulfilledEventTimeout,
			l,
		)
		require.NoError(t, err)

		amountToWithdrawLink := fulfilledEventLink.Payment

		defaultWalletBalanceLinkBeforeOracleWithdraw, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		l.Info().
			Str("Returning to", defaultWalletAddress).
			Str("Amount", amountToWithdrawLink.String()).
			Msg("Invoking Oracle Withdraw for LINK")

		err = vrfv2Contracts.Coordinator.OracleWithdraw(common.HexToAddress(defaultWalletAddress), amountToWithdrawLink)
		require.NoError(t, err, "Error withdrawing LINK from coordinator to default wallet")

		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, vrfv2_actions.ErrWaitTXsComplete)

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
		testConfig := vrfv2Config
		subIDsForCancelling, err := vrfv2_actions.CreateFundSubsAndAddConsumers(
			env,
			testConfig,
			linkToken,
			vrfv2Contracts.Coordinator,
			vrfv2Contracts.LoadTestConsumers,
			1,
		)
		require.NoError(t, err)
		subIDForCancelling := subIDsForCancelling[0]

		testWalletAddress, err := actions.GenerateWallet()
		require.NoError(t, err)

		testWalletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
		require.NoError(t, err)

		subscriptionForCancelling, err := vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "error getting subscription information")

		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForCancelling).
			Str("Returning funds to", testWalletAddress.String()).
			Msg("Canceling subscription and returning funds to subscription owner")

		tx, err := vrfv2Contracts.Coordinator.CancelSubscription(subIDForCancelling, testWalletAddress)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2Contracts.Coordinator.WaitForSubscriptionCanceledEvent([]uint64{subIDForCancelling}, time.Second*30)
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
		_, err = vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), subIDForCancelling)
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
		testConfig := vrfv2Config

		// Underfund subscription to force fulfillments to fail
		testConfig.SubscriptionFundingAmountLink = float64(0.000000000000000001) // 1 Juel

		subIDsForCancelling, err := vrfv2_actions.CreateFundSubsAndAddConsumers(
			env,
			testConfig,
			linkToken,
			vrfv2Contracts.Coordinator,
			vrfv2Contracts.LoadTestConsumers,
			1,
		)
		require.NoError(t, err)

		subIDForCancelling := subIDsForCancelling[0]

		subscriptionForCancelling, err := vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "Error getting subscription information")

		vrfv2_actions.LogSubDetails(l, subscriptionForCancelling, subIDForCancelling, vrfv2Contracts.Coordinator)

		// No GetActiveSubscriptionIds function available - skipping check

		pendingRequestsExist, err := vrfv2Contracts.Coordinator.PendingRequestsExist(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err)
		require.False(t, pendingRequestsExist, "Pending requests should not exist")

		// Request randomness - should fail due to underfunded subscription
		randomWordsFulfilledEventTimeout := 5 * time.Second
		_, err = vrfv2_actions.RequestRandomnessAndWaitForFulfillment(
			vrfv2Contracts.LoadTestConsumers[0],
			vrfv2Contracts.Coordinator,
			vrfv2Data,
			subIDForCancelling,
			testConfig.RandomnessRequestCountPerRequest,
			testConfig,
			randomWordsFulfilledEventTimeout,
			l,
		)
		require.Error(t, err, "Error should occur while waiting for fulfilment due to low sub balance")

		pendingRequestsExist, err = vrfv2Contracts.Coordinator.PendingRequestsExist(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err)
		require.True(t, pendingRequestsExist, "Pending requests should exist after unfilfulled requests due to low sub balance")

		walletBalanceLinkBeforeSubCancelling, err := linkToken.BalanceOf(testcontext.Get(t), defaultWalletAddress)
		require.NoError(t, err)

		subscriptionForCancelling, err = vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), subIDForCancelling)
		require.NoError(t, err, "Error getting subscription information")
		subBalanceLink := subscriptionForCancelling.Balance

		l.Info().
			Str("Subscription Amount Link", subBalanceLink.String()).
			Uint64("Returning funds from SubID", subIDForCancelling).
			Str("Returning funds to", defaultWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")

		// Call OwnerCancelSubscription
		tx, err := vrfv2Contracts.Coordinator.OwnerCancelSubscription(subIDForCancelling)
		require.NoError(t, err, "Error canceling subscription")

		subscriptionCanceledEvent, err := vrfv2Contracts.Coordinator.WaitForSubscriptionCanceledEvent([]uint64{subIDForCancelling}, time.Second*30)
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
		_, err = vrfv2Contracts.Coordinator.GetSubscription(testcontext.Get(t), subIDForCancelling)
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

	network, err := actions.EthereumNetworkConfigFromEnvOrDefault(l)
	require.NoError(t, err, "Error building ethereum network config")

	var vrfv2Config vrfv2_config.VRFV2Config
	err = envconfig.Process("VRFV2", &vrfv2Config)
	require.NoError(t, err)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(vrfv2Config.ChainlinkNodeFunding)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(vrfv2Config.LinkNativeFeedResponse))
	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 2
	vrfv2Contracts, subIDs, vrfv2Data, err := vrfv2_actions.SetupVRFV2Environment(
		env,
		vrfv2Config,
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

	subscription, err := vrfv2Contracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2_actions.LogSubDetails(l, subscription, subID, vrfv2Contracts.Coordinator)

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		testConfig := vrfv2Config
		txKeys, _, err := env.ClCluster.Nodes[0].API.ReadTxKeys("evm")
		require.NoError(t, err, "error reading tx keys")

		require.Equal(t, numberOfTxKeysToCreate+1, len(txKeys.Data))

		var fulfillmentTxFromAddresses []string
		for i := 0; i < numberOfTxKeysToCreate+1; i++ {
			randomWordsFulfilledEvent, err := vrfv2_actions.RequestRandomnessAndWaitForFulfillment(
				vrfv2Contracts.LoadTestConsumers[0],
				vrfv2Contracts.Coordinator,
				vrfv2Data,
				subID,
				testConfig.RandomnessRequestCountPerRequest,
				testConfig,
				testConfig.RandomWordsFulfilledEventTimeout,
				l,
			)
			require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

			//todo - move TransactionByHash to EVMClient in CTF
			fulfillmentTx, _, err := env.EVMClient.(*blockchain.EthereumMultinodeClient).DefaultClient.(*blockchain.EthereumClient).
				Client.TransactionByHash(context.Background(), randomWordsFulfilledEvent.Raw.TxHash)
			require.NoError(t, err, "error getting tx from hash")
			fulfillmentTxFromAddress, err := actions.GetTxFromAddress(fulfillmentTx)
			require.NoError(t, err, "error getting tx from address")
			fulfillmentTxFromAddresses = append(fulfillmentTxFromAddresses, fulfillmentTxFromAddress)
		}
		require.Equal(t, numberOfTxKeysToCreate+1, len(fulfillmentTxFromAddresses))
		var txKeyAddresses []string
		for _, txKey := range txKeys.Data {
			txKeyAddresses = append(txKeyAddresses, txKey.ID)
		}
		less := func(a, b string) bool { return a < b }
		equalIgnoreOrder := cmp.Diff(txKeyAddresses, fulfillmentTxFromAddresses, cmpopts.SortSlices(less)) == ""
		require.True(t, equalIgnoreOrder)
	})
}
