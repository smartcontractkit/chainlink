package smoke

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/stretchr/testify/require"
)

func TestVRFv2PlusWrapperMigration(t *testing.T) {
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

	vrfv2PlusContracts, _, vrfv2PlusData, nodesMap, err := vrfv2plus.SetupVRFV2_5Environment(
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
	require.Len(t, activeSubIdsOldCoordinatorBeforeMigration, 2, "Active Sub Ids length is not equal to 2")
	require.Equal(t, subID, activeSubIdsOldCoordinatorBeforeMigration[1])

	oldSubscriptionBeforeMigration, err := vrfv2PlusContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	require.NoError(t, err, "error getting subscription information")

	// Deploy new coordinator, set configs and add VRF job spec in CL node
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

	// Register new coordinator as migratable coordinator with the old coordinator
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

	// Migrate coordinator using VRFV2PlusWrapper's migrate method
	err = wrapperContracts.VRFV2PlusWrapper.Migrate(testcontext.Get(t), common.HexToAddress(newCoordinator.Address()))
	// err = vrfv2PlusContracts.CoordinatorV2Plus.Migrate(subID, newCoordinator.Address())

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

	activeSubIdsOldCoordinator, err := vrfv2PlusContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
	require.NoError(t, err, "error occurred getting active sub ids")
	require.Len(t, activeSubIdsOldCoordinator, 1, "Active SubIds length is not 1 for old Coordinator after migration- wrapper subId should be removed")
	require.NotContains(t, activeSubIdsOldCoordinator, subID)

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
	consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	require.NoError(t, err, "error getting rand request status")
	require.True(t, consumerStatus.Fulfilled)

	// Verify rand requests fulfills with Native Token billing
	isNativeBilling = true
	randomWordsFulfilledEvent, err = vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillmentUpgraded(
		wrapperContracts.LoadTestConsumers[0],
		newCoordinator,
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
	consumerStatus, err = wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	require.NoError(t, err, "error getting rand request status")
	require.True(t, consumerStatus.Fulfilled)
}
