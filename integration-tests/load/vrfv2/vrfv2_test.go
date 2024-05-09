package loadvrfv2

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	testEnv          *test_env.CLClusterTestEnv
	vrfContracts     *vrfcommon.VRFContracts
	vrfKey           *vrfcommon.VRFKeyData
	subIDs           []uint64
	eoaWalletAddress string

	labels = map[string]string{
		"branch": "vrfv2_healthcheck",
		"commit": "vrfv2_healthcheck",
	}
)

func TestVRFV2Performance(t *testing.T) {
	l := logging.GetTestLogger(t)

	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2)
	require.NoError(t, err)

	testReporter := &testreporters.VRFV2TestReporter{}
	vrfv2Config := testConfig.VRFv2

	cfgl := testConfig.Logging.Loki
	lokiConfig := wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken)
	lc, err := wasp.NewLokiClient(lokiConfig)
	if err != nil {
		l.Error().Err(err).Msg(ErrLokiClient)
		return
	}

	updatedLabels := UpdateLabels(labels, t)

	l.Info().
		Str("Test Type", testType).
		Str("Test Duration", vrfv2Config.Performance.TestDuration.Duration.Truncate(time.Second).String()).
		Int64("RPS", *vrfv2Config.Performance.RPS).
		Str("RateLimitUnitDuration", vrfv2Config.Performance.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", *vrfv2Config.General.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", *vrfv2Config.General.RandomnessRequestCountPerRequestDeviation).
		Bool("UseExistingEnv", *vrfv2Config.General.UseExistingEnv).
		Msg("Performance Test Configuration")

	chainID := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2Consumers[0], lc, updatedLabels, testReporter, testType, &testConfig)

		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "error getting EVM client")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, eoaWalletAddress, subIDs, l)
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
		NumberOfTxKeysToCreate: *vrfv2Config.General.NumberOfSendingKeysToCreate,
		UseVRFOwner:            true,
		UseTestCoordinator:     true,
	}

	testEnv, vrfContracts, vrfKey, _, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, testConfig, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2 universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "error getting EVM client")

	var consumers []contracts.VRFv2LoadTestConsumer
	subIDs, consumers, err = vrfv2.SetupSubsAndConsumersForExistingEnv(
		testEnv,
		chainID,
		vrfContracts.CoordinatorV2,
		vrfContracts.LinkToken,
		1,
		*vrfv2Config.General.NumberOfSubToCreate,
		testConfig,
		l,
	)
	vrfContracts.VRFV2Consumers = consumers

	eoaWalletAddress = evmClient.GetDefaultWallet().Address()

	l.Debug().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")
	for _, subID := range subIDs {
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information for subscription %d", subID)
		vrfv2.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2)
	}

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2 performance test", func(t *testing.T) {
		require.Len(t, vrfContracts.VRFV2Consumers, 1, "only one consumer should be created for Load Test")
		err = vrfContracts.VRFV2Consumers[0].ResetMetrics()
		require.NoError(t, err)
		MonitorLoadStats(testcontext.Get(t), lc, vrfContracts.VRFV2Consumers[0], updatedLabels)

		singleFeedConfig := &wasp.Config{
			T:                     t,
			LoadType:              wasp.RPS,
			GenName:               "gun",
			RateLimitUnitDuration: vrfv2Config.Performance.RateLimitUnitDuration.Duration,
			Gun: NewSingleHashGun(
				vrfContracts,
				vrfKey.KeyHash,
				subIDs,
				vrfv2Config,
				l,
			),
			Labels:      labels,
			LokiConfig:  lokiConfig,
			CallTimeout: 2 * time.Minute,
		}

		singleFeedConfig.Schedule = wasp.Plain(
			*vrfv2Config.Performance.RPS,
			vrfv2Config.Performance.TestDuration.Duration,
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)
		//todo - timeout should be configurable depending on the perf test type
		requestCount, fulfilmentCount, err := vrfcommon.WaitForRequestCountEqualToFulfilmentCount(testcontext.Get(t), vrfContracts.VRFV2Consumers[0], 2*time.Minute, &wg)
		require.NoError(t, err)
		wg.Wait()

		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Final Request/Fulfilment Stats")
	})
}

func TestVRFV2BHSPerformance(t *testing.T) {
	l := logging.GetTestLogger(t)

	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2)
	require.NoError(t, err)

	testReporter := &testreporters.VRFV2TestReporter{}
	vrfv2Config := testConfig.VRFv2

	cfgl := testConfig.Logging.Loki
	lokiConfig := wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken)
	lc, err := wasp.NewLokiClient(lokiConfig)
	if err != nil {
		l.Error().Err(err).Msg(ErrLokiClient)
		return
	}

	updatedLabels := UpdateLabels(labels, t)

	l.Info().
		Str("Test Type", testType).
		Str("Test Duration", vrfv2Config.Performance.TestDuration.Duration.Truncate(time.Second).String()).
		Int64("RPS", *vrfv2Config.Performance.RPS).
		Str("RateLimitUnitDuration", vrfv2Config.Performance.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", *vrfv2Config.General.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", *vrfv2Config.General.RandomnessRequestCountPerRequestDeviation).
		Bool("UseExistingEnv", *vrfv2Config.General.UseExistingEnv).
		Msg("Performance Test Configuration")

	chainID := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2Consumers[0], lc, updatedLabels, testReporter, testType, &testConfig)

		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "error getting EVM client")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, eoaWalletAddress, subIDs, l)
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
		NumberOfTxKeysToCreate: *vrfv2Config.General.NumberOfSendingKeysToCreate,
		UseVRFOwner:            true,
		UseTestCoordinator:     true,
	}

	testEnv, vrfContracts, vrfKey, _, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, testConfig, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2 universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "error getting EVM client")

	var consumers []contracts.VRFv2LoadTestConsumer
	subIDs, consumers, err = vrfv2.SetupSubsAndConsumersForExistingEnv(
		testEnv,
		chainID,
		vrfContracts.CoordinatorV2,
		vrfContracts.LinkToken,
		1,
		*vrfv2Config.General.NumberOfSubToCreate,
		testConfig,
		l,
	)
	vrfContracts.VRFV2Consumers = consumers

	eoaWalletAddress = evmClient.GetDefaultWallet().Address()

	l.Debug().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")
	for _, subID := range subIDs {
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information for subscription %d", subID)
		vrfv2.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2)
	}

	t.Run("vrfv2 and bhs performance test", func(t *testing.T) {
		configCopy := testConfig.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		consumers, subIDs, err = vrfv2.SetupNewConsumersAndSubs(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			configCopy,
			vrfContracts.LinkToken,
			1,
			*configCopy.VRFv2.General.NumberOfSubToCreate,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subscriptions")
		vrfContracts.VRFV2Consumers = consumers
		require.Len(t, vrfContracts.VRFV2Consumers, 1, "only one consumer should be created for Load Test")
		err = vrfContracts.VRFV2Consumers[0].ResetMetrics()
		require.NoError(t, err, "error resetting consumer metrics")
		MonitorLoadStats(testcontext.Get(t), lc, vrfContracts.VRFV2Consumers[0], updatedLabels)

		singleFeedConfig := &wasp.Config{
			T:                     t,
			LoadType:              wasp.RPS,
			GenName:               "gun",
			RateLimitUnitDuration: configCopy.VRFv2.Performance.BHSTestRateLimitUnitDuration.Duration,
			Gun: NewBHSTestGun(
				vrfContracts,
				vrfKey.KeyHash,
				subIDs,
				configCopy.VRFv2,
				l,
			),
			Labels:      labels,
			LokiConfig:  lokiConfig,
			CallTimeout: 2 * time.Minute,
		}

		singleFeedConfig.Schedule = wasp.Plain(
			*configCopy.VRFv2.Performance.BHSTestRPS,
			configCopy.VRFv2.Performance.BHSTestDuration.Duration,
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)

		var wgBlockNumberTobe sync.WaitGroup
		wgBlockNumberTobe.Add(1)
		//Wait at least 256 blocks
		latestBlockNumber, err := evmClient.LatestBlockNumber(testcontext.Get(t))
		require.NoError(t, err, "error getting latest block number")
		_, err = actions.WaitForBlockNumberToBe(latestBlockNumber+uint64(256), evmClient, &wgBlockNumberTobe, configCopy.VRFv2.General.WaitFor256BlocksTimeout.Duration, t)
		wgBlockNumberTobe.Wait()
		require.NoError(t, err, "error waiting for block number to be")
		err = vrfv2.FundSubscriptions(testEnv, chainID, big.NewFloat(*configCopy.VRFv2.General.SubscriptionRefundingAmountLink), vrfContracts.LinkToken, vrfContracts.CoordinatorV2, subIDs)
		require.NoError(t, err, "error funding subscriptions")
		var wgAllRequestsFulfilled sync.WaitGroup
		wgAllRequestsFulfilled.Add(1)
		requestCount, fulfilmentCount, err := vrfcommon.WaitForRequestCountEqualToFulfilmentCount(testcontext.Get(t), vrfContracts.VRFV2Consumers[0], 2*time.Minute, &wgAllRequestsFulfilled)
		require.NoError(t, err)
		wgAllRequestsFulfilled.Wait()

		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Final Request/Fulfilment Stats")
	})
}
func teardown(
	t *testing.T,
	consumer contracts.VRFv2LoadTestConsumer,
	lc *wasp.LokiClient,
	updatedLabels map[string]string,
	testReporter *testreporters.VRFV2TestReporter,
	testType string,
	testConfig *tc.TestConfig,
) {
	//send final results to Loki
	metrics := GetLoadTestMetrics(testcontext.Get(t), consumer)
	SendMetricsToLoki(metrics, lc, updatedLabels)
	//set report data for Slack notification
	testReporter.SetReportData(
		testType,
		testreporters.VRFLoadTestMetrics{
			RequestCount:                 metrics.RequestCount,
			FulfilmentCount:              metrics.FulfilmentCount,
			AverageFulfillmentInMillions: metrics.AverageFulfillmentInMillions,
			SlowestFulfillment:           metrics.SlowestFulfillment,
			FastestFulfillment:           metrics.FastestFulfillment,
		},
		testConfig,
	)

	// send Slack notification
	err := testReporter.SendSlackNotification(t, nil)
	if err != nil {
		log.Warn().Err(err).Msg("Error sending Slack notification")
	}
}
