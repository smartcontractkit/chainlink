package smoke

import (
	"context"
	"github.com/pkg/errors"
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

	vrfv2PlusContracts, subID, vrfv2PlusData, err := vrfv2plus.SetupVRFV2PlusEnvironment(env, linkAddress, mockETHLinkFeedAddress, 1)
	require.NoError(t, err, "error setting up VRF v2 Plus env")

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	l.Debug().
		Interface("Juels Balance", subscription.Balance).
		Interface("Native Token Balance", subscription.EthBalance).
		Interface("Subscription ID", subID).
		Interface("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")

	t.Run("VRFV2 Plus With Link Billing", func(t *testing.T) {
		var isNativeBilling = false
		subBalanceBeforeRequest := subscription.Balance

		jobRunsBeforeTest, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")

		// test and assert
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subID,
			isNativeBilling,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
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
		randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			vrfv2PlusContracts.LoadTestConsumers[0],
			vrfv2PlusContracts.Coordinator,
			vrfv2PlusData,
			subID,
			isNativeBilling,
			l,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err)
		subBalanceAfterRequest := subscription.EthBalance
		require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)

		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(vrfv2PlusData.VRFJob.Data.ID)
		require.NoError(t, err, "error reading job runs")
		require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))

		status, err := vrfv2PlusContracts.LoadTestConsumers[0].GetRequestStatus(context.Background(), randomWordsFulfilledEvent.RequestId)
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

	oldVRFV2PlusContracts, subID, vrfv2PlusData, err := vrfv2plus.SetupVRFV2PlusEnvironment(env, linkAddress, mockETHLinkFeedAddress, 2)
	require.NoError(t, err, "error setting up VRF v2 Plus env")

	for _, consumer := range oldVRFV2PlusContracts.LoadTestConsumers {
		coordinatorAddressInConsumer, err := consumer.GetCoordinator(context.Background())
		require.NoError(t, err, "error getting Coordinator from Consumer contract")
		require.Equal(t, oldVRFV2PlusContracts.Coordinator.Address(), coordinatorAddressInConsumer.String())
		l.Debug().
			Interface("Consumer", consumer.Address()).
			Interface("Coordinator", coordinatorAddressInConsumer.String()).
			Msg("Coordinator Address in Consumer Before Migration")
	}

	oldSubscription, err := oldVRFV2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	//Migration Process
	newCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2PlusUpgradedVersion(oldVRFV2PlusContracts.BHS.Address())
	require.NoError(t, err, vrfv2plus.ErrDeployCoordinator)

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	_, err = vrfv2plus.VRFV2PlusUpgradedVersionRegisterProvingKey(vrfv2PlusData.VRFKey, vrfv2PlusData.PrimaryEthAddress, newCoordinator)
	require.NoError(t, err, errors.Wrap(err, vrfv2plus.ErrRegisteringProvingKey))

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

	err = oldVRFV2PlusContracts.Coordinator.RegisterMigratableCoordinator(newCoordinator.Address())

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	err = oldVRFV2PlusContracts.Coordinator.Migrate(subID, newCoordinator.Address())
	require.NoError(t, err, "error migrating sub id ", subID.String(), " from ", oldVRFV2PlusContracts.Coordinator.Address(), " to new Coordinator address ", newCoordinator.Address())
	migrationCompletedEvent, err := oldVRFV2PlusContracts.Coordinator.WaitForMigrationCompletedEvent(time.Minute * 1)
	require.NoError(t, err, "error waiting for MigrationCompleted event")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, vrfv2plus.ErrWaitTXsComplete)

	l.Debug().
		Str("Subscription ID", migrationCompletedEvent.SubId.String()).
		Str("Migrated From Coordinator", oldVRFV2PlusContracts.Coordinator.Address()).
		Str("Migrated To Coordinator", migrationCompletedEvent.NewCoordinator.String()).
		Msg("MigrationCompleted Event")

	migratedSubscription, err := newCoordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	l.Debug().
		Interface("New Coordinator", newCoordinator.Address()).
		Interface("Subscription ID", subID).
		Interface("Juels Balance", migratedSubscription.Balance).
		Interface("Native Token Balance", migratedSubscription.EthBalance).
		Interface("Subscription Owner", migratedSubscription.Owner.String()).
		Interface("Subscription Consumers", migratedSubscription.Consumers).
		Msg("Subscription Data After Migration to New Coordinator")

	//Verify that Coordinators were updated in Consumers
	for _, consumer := range oldVRFV2PlusContracts.LoadTestConsumers {
		coordinatorAddressInConsumerAfterMigration, err := consumer.GetCoordinator(context.Background())
		require.NoError(t, err, "error getting Coordinator from Consumer contract")
		require.Equal(t, newCoordinator.Address(), coordinatorAddressInConsumerAfterMigration.String())
		l.Debug().
			Interface("Consumer", consumer.Address()).
			Interface("Coordinator", coordinatorAddressInConsumerAfterMigration).
			Msg("Coordinator Address in Consumer After Migration")
	}

	//Verify old and migrated subs
	require.Equal(t, oldSubscription.EthBalance, migratedSubscription.EthBalance)
	require.Equal(t, oldSubscription.Balance, migratedSubscription.Balance)
	require.Equal(t, oldSubscription.Owner, migratedSubscription.Owner)
	require.Equal(t, oldSubscription.Consumers, migratedSubscription.Consumers)

	//Verify that old sub was deleted from old Coordinator
	oldSubscription, err = oldVRFV2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")

	//Verify rand requests fulfills with Link Token billing
	_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
		oldVRFV2PlusContracts.LoadTestConsumers[0],
		newCoordinator,
		vrfv2PlusData,
		subID,
		false,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

	//Verify rand requests fulfills with Native Token billing
	_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
		oldVRFV2PlusContracts.LoadTestConsumers[1],
		newCoordinator,
		vrfv2PlusData,
		subID,
		true,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

}
