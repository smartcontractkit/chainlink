package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

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

	vrfv2PlusContracts, subIDs, vrfv2PlusData, err := vrfv2plus.SetupVRFV2_5Environment(env, &vrfv2PlusConfig, linkToken, mockETHLinkFeed, 1, 1)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subID := subIDs[0]

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.Coordinator)

	t.Run("VRFV2 Plus With Link Billing", func(t *testing.T) {
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
			vrfv2PlusConfig.RandomnessRequestCountPerRequest,
			&vrfv2PlusConfig,
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

		require.Equal(t, vrfv2PlusConfig.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})

	t.Run("VRFV2 Plus With Native Billing", func(t *testing.T) {
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
			vrfv2PlusConfig.RandomnessRequestCountPerRequest,
			&vrfv2PlusConfig,
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

		require.Equal(t, vrfv2PlusConfig.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})

	wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperEnvironment(
		env,
		&vrfv2PlusConfig,
		linkToken,
		mockETHLinkFeed,
		vrfv2PlusContracts.Coordinator,
		vrfv2PlusData.KeyHash,
		1,
	)
	require.NoError(t, err)

	t.Run("VRFV2 Plus With Direct Funding (VRFV2PlusWrapper) - Link Billing", func(t *testing.T) {
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
			&vrfv2PlusConfig,
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

		require.Equal(t, vrfv2PlusConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
		for _, w := range consumerStatus.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
	})

	t.Run("VRFV2 Plus With Direct Funding (VRFV2PlusWrapper) - Native Billing", func(t *testing.T) {
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
			&vrfv2PlusConfig,
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

		require.Equal(t, vrfv2PlusConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
		for _, w := range consumerStatus.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}
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

	vrfv2PlusContracts, subIDs, vrfv2PlusData, err := vrfv2plus.SetupVRFV2_5Environment(env, &vrfv2PlusConfig, linkAddress, mockETHLinkFeedAddress, 2, 1)
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
	require.Equal(t, expectedLinkTotalBalanceForMigratedCoordinator, migratedCoordinatorLinkTotalBalanceAfterMigration)
	require.Equal(t, expectedEthTotalBalanceForMigratedCoordinator, migratedCoordinatorEthTotalBalanceAfterMigration)
	require.Equal(t, expectedLinkTotalBalanceForOldCoordinator, oldCoordinatorLinkTotalBalanceAfterMigration)
	require.Equal(t, expectedEthTotalBalanceForOldCoordinator, oldCoordinatorEthTotalBalanceAfterMigration)

	//Verify rand requests fulfills with Link Token billing
	_, err = vrfv2plus.RequestRandomnessAndWaitForFulfillmentUpgraded(
		vrfv2PlusContracts.LoadTestConsumers[0],
		newCoordinator,
		vrfv2PlusData,
		subID,
		false,
		&vrfv2PlusConfig,
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
		&vrfv2PlusConfig,
		l,
	)
	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

}
