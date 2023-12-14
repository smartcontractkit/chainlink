package smoke

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestVRFv2Basic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig(t.Name(), tc.Smoke, tc.VRFv2)
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
		WithLogStream().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(config.VRFv2.General.LinkNativeFeedResponse))
	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 1
	vrfv2Contracts, subIDs, vrfv2Data, err := vrfv2_actions.SetupVRFV2Environment(
		env,
		&config,
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
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfv2Data.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2_actions.RequestRandomnessAndWaitForFulfillment(
			vrfv2Contracts.LoadTestConsumers[0],
			vrfv2Contracts.Coordinator,
			vrfv2Data,
			subID,
			config.VRFv2.General.RandomnessRequestCountPerRequest,
			&config,
			config.VRFv2.General.RandomWordsFulfilledEventTimeout,
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

		require.Equal(t, config.VRFv2.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})
}

func TestVRFv2MultipleSendingKeys(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig(t.Name(), tc.Smoke, tc.VRFv2)
	if err != nil {
		t.Fatal(err)
	}

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
		WithStandardCleanup().
		WithLogStream().
		Build()
	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(config.VRFv2.General.LinkNativeFeedResponse))
	require.NoError(t, err)
	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)

	// register proving key against oracle address (sending key) in order to test oracleWithdraw
	defaultWalletAddress := env.EVMClient.GetDefaultWallet().Address()

	numberOfTxKeysToCreate := 2
	vrfv2Contracts, subIDs, vrfv2Data, err := vrfv2_actions.SetupVRFV2Environment(
		env,
		&config,
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
				config.VRFv2.General.RandomnessRequestCountPerRequest,
				&config,
				config.VRFv2.General.RandomWordsFulfilledEventTimeout,
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
