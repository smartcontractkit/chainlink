package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestVRFv2PlusBilling(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithCLNodes(1).
		WithFunding(vrfv2plus_constants.ChainlinkNodeFundingAmountEth).
		Build()
	require.NoError(t, err, "error creating test env")
	t.Cleanup(func() {
		if err := env.Cleanup(); err != nil {
			l.Error().Err(err).Msg("Error cleaning up test environment")
		}
	})

	env.ParallelTransactions(true)

	mockETHLinkFeedAddress, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, vrfv2plus_constants.LinkEthFeedResponse)
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkAddress, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	vrfv2PlusContracts, subID, vrfv2PlusData, err := vrfv2plus.SetupVRFV2PlusEnvironment(env, linkAddress, mockETHLinkFeedAddress)
	require.NoError(t, err, "error setting up VRF v2 Plus env")

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	l.Debug().
		Interface("Juels Balance", subscription.Balance).
		Interface("Native Token Balance", subscription.EthBalance).
		Interface("Subscription ID", subID).
		Msg("Subscription Data")

	t.Run("VRFV2 Plus With Link Billing", func(t *testing.T) {
		var isNativeBilling = false
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		err = vrfv2PlusContracts.LoadTestConsumer.RequestRandomness(
			vrfv2PlusData.KeyHash,
			subID,
			vrfv2plus_constants.MinimumConfirmations,
			vrfv2plus_constants.CallbackGasLimit,
			isNativeBilling,
			vrfv2plus_constants.NumberOfWords,
			vrfv2plus_constants.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsFulfilledEvent, err := vrfv2PlusContracts.Coordinator.WaitForRandomWordsFulfilledEvent([]*big.Int{subID}, nil, time.Minute*2)
		require.NoError(t, err, "error waiting for RandomWordsFulfilled event")

		l.Debug().
			Interface("Total Payment in Juels", randomWordsFulfilledEvent.Payment).
			Interface("TX Hash", randomWordsFulfilledEvent.Raw.TxHash).
			Interface("Subscription ID", randomWordsFulfilledEvent.SubID).
			Interface("Request ID", randomWordsFulfilledEvent.RequestId).
			Bool("Success", randomWordsFulfilledEvent.Success).
			Msg("Randomness Fulfillment TX metadata")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.LoadTestConsumer.GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Interface("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, vrfv2plus_constants.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, w.Cmp(big.NewInt(0)), 1, "Expected the VRF job give an answer bigger than 0")
		}
	})

	t.Run("VRFV2 Plus With Native Billing", func(t *testing.T) {
		var isNativeBilling = true
		subNativeTokenBalanceBeforeRequest := subscription.EthBalance

		jobRunsBeforeTest, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)

		// test and assert
		err = vrfv2PlusContracts.LoadTestConsumer.RequestRandomness(
			vrfv2PlusData.KeyHash,
			subID,
			vrfv2plus_constants.MinimumConfirmations,
			vrfv2plus_constants.CallbackGasLimit,
			isNativeBilling,
			vrfv2plus_constants.NumberOfWords,
			vrfv2plus_constants.RandomnessRequestCountPerRequest,
		)
		require.NoError(t, err, "error requesting randomness")

		randomWordsFulfilledEvent, err := vrfv2PlusContracts.Coordinator.WaitForRandomWordsFulfilledEvent([]*big.Int{subID}, nil, time.Minute*2)
		require.NoError(t, err, "error waiting for RandomWordsFulfilled event")

		l.Debug().
			Interface("Total Payment in Wei", randomWordsFulfilledEvent.Payment).
			Interface("TX Hash", randomWordsFulfilledEvent.Raw.TxHash).
			Interface("Subscription ID", randomWordsFulfilledEvent.SubID).
			Interface("Request ID", randomWordsFulfilledEvent.RequestId).
			Bool("Success", randomWordsFulfilledEvent.Success).
			Msg("Randomness Fulfillment TX metadata")

		expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err)
		subBalanceAfterRequest := subscription.EthBalance
		require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.LoadTestConsumer.GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Debug().Interface("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, vrfv2plus_constants.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, w.Cmp(big.NewInt(0)), 1, "Expected the VRF job give an answer bigger than 0")
		}
	})

}

func TestVRFv2PlusMigration(t *testing.T) {
	t.Parallel()
	//l := utils.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithCLNodes(1).
		WithFunding(vrfv2plus_constants.ChainlinkNodeFundingAmountEth).
		Build()

	require.NoError(t, err, "error creating test env")

	env.ParallelTransactions(true)

	mockETHLinkFeedAddress, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, vrfv2plus_constants.LinkEthFeedResponse)
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkAddress, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	//todo - add more consumers to the sub with diff eth and link balances
	oldVRFV2PlusContracts, subID, vrfv2PlusData, err := vrfv2plus.SetupVRFV2PlusEnvironment(env, linkAddress, mockETHLinkFeedAddress)
	require.NoError(t, err, "error setting up VRF v2 Plus env")

	newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(oldVRFV2PlusContracts.BHS.Address())
	require.NoError(t, err, vrfv2plus.ErrDeployCoordinator)

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	err = newCoordinator.SetConfig(
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.MaxGasLimitVRFCoordinatorConfig,
		vrfv2plus_constants.StalenessSeconds,
		vrfv2plus_constants.GasAfterPaymentCalculation,
		vrfv2plus_constants.LinkEthFeedResponse,
		vrfv2plus_constants.VRFCoordinatorV2PlusUpgradedVersionFeeConfig,
	)

	err = newCoordinator.SetLINKAndLINKETHFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
	require.NoError(t, err, vrfv2plus.ErrSetLinkETHLinkFeed)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	_, err = vrfv2plus.CreateVRFV2PlusJob(
		env.GetAPIs()[0],
		newCoordinator.Address(),
		vrfv2PlusData.PrimaryEthAddress,
		vrfv2PlusData.VRFKey.Data.ID,
		vrfv2PlusData.ChainID.String(),
		vrfv2plus_constants.MinimumConfirmations,
	)
	require.NoError(t, err, vrfv2plus.ErrCreateVRFV2PlusJobs)

	// test and assert
	err = oldVRFV2PlusContracts.LoadTestConsumer.RequestRandomness(
		vrfv2PlusData.KeyHash,
		subID,
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.CallbackGasLimit,
		true,
		vrfv2plus_constants.NumberOfWords,
		vrfv2plus_constants.RandomnessRequestCountPerRequest,
	)
	require.NoError(t, err, "error requesting randomness")

	_, err = oldVRFV2PlusContracts.Coordinator.WaitForRandomWordsFulfilledEvent([]*big.Int{subID}, nil, time.Minute*2)
	require.NoError(t, err, "error waiting for RandomWordsFulfilled event")

	err = oldVRFV2PlusContracts.Coordinator.RegisterMigratableCoordinator(newCoordinator.Address())

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	err = oldVRFV2PlusContracts.Coordinator.Migrate(subID, newCoordinator.Address())
	require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", oldVRFV2PlusContracts.Coordinator.Address(), " to new Coordinator address ", newCoordinator.Address())

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	// test and assert
	err = oldVRFV2PlusContracts.LoadTestConsumer.RequestRandomness(
		vrfv2PlusData.KeyHash,
		subID,
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.CallbackGasLimit,
		true,
		vrfv2plus_constants.NumberOfWords,
		vrfv2plus_constants.RandomnessRequestCountPerRequest,
	)
	require.NoError(t, err, "error requesting randomness")

	_, err = newCoordinator.WaitForRandomWordsFulfilledEvent([]*big.Int{subID}, nil, time.Minute*2)
	require.NoError(t, err, "error waiting for RandomWordsFulfilled event")

	//todo - check consumers, eth and link balances

}
