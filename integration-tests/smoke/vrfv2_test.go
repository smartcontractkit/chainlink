package smoke

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_config"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
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
		WithLogStream().
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
}

func TestVRFv2MultipleSendingKeys(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	var vrfv2Config vrfv2_config.VRFV2Config
	err := envconfig.Process("VRFV2", &vrfv2Config)
	require.NoError(t, err)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(big.NewFloat(vrfv2Config.ChainlinkNodeFunding)).
		WithStandardCleanup().
		WithLogStream().
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
