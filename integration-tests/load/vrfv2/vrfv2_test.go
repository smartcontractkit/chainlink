package loadvrfv2

import (
	"math/big"
	"strconv"
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
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	labels = map[string]string{
		"branch": "vrfv2_healthcheck",
		"commit": "vrfv2_healthcheck",
	}
)

func TestVRFV2Performance(t *testing.T) {
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
	)
	l := logging.GetTestLogger(t)
	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2)
	require.NoError(t, err)
	cfgl := testConfig.Logging.Loki

	vrfv2Config := testConfig.VRFv2
	testReporter := &testreporters.VRFV2TestReporter{}

	lokiConfig := wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken)
	lc, err := wasp.NewLokiClient(lokiConfig)
	if err != nil {
		l.Error().Err(err).Msg(ErrLokiClient)
		return
	}
	network := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0]
	chainID := network.ChainID
	sethClient, err := actions_seth.GetChainClient(testConfig, network)
	require.NoError(t, err, "Error creating seth client")

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
	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2Consumers[0], lc, updatedLabels, testReporter, testType, &testConfig)

		require.NoError(t, err, "Getting Seth client shouldn't fail")

		if sethClient.Cfg.IsSimulatedNetwork() {
			l.Info().
				Str("Network Name", sethClient.Cfg.Network.Name).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, sethClient.MustGetRootKeyAddress().Hex(), subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: testConfig,
		ChainID:    chainID,
		CleanupFn:  cleanupFn,
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          *vrfv2Config.General.NumberOfSendingKeysToCreate,
		UseVRFOwner:                     true,
		UseTestCoordinator:              true,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	testEnv, vrfContracts, vrfKey, _, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2 universe")

	var consumers []contracts.VRFv2LoadTestConsumer
	subIDs, consumers, err := vrfv2.SetupSubsAndConsumersForExistingEnv(
		testEnv,
		chainID,
		vrfContracts.CoordinatorV2,
		vrfContracts.LinkToken,
		1,
		*vrfv2Config.General.NumberOfSubToCreate,
		testConfig,
		l,
	)
	require.NoError(t, err, "error setting up new consumers and subs")
	for _, subID := range subIDs {
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information for subscription %d", subID)
		vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subID, 10), vrfContracts.CoordinatorV2)
	}
	subIDsForCancellingAfterTest = subIDs
	l.Debug().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")

	vrfContracts.VRFV2Consumers = consumers

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
				sethClient,
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
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []uint64
		vrfKey                       *vrfcommon.VRFKeyData
	)
	l := logging.GetTestLogger(t)

	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2)
	require.NoError(t, err)
	vrfv2Config := testConfig.VRFv2
	testReporter := &testreporters.VRFV2TestReporter{}
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

	network := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0]
	chainID := network.ChainID
	sethClient, err := actions_seth.GetChainClient(testConfig, network)
	require.NoError(t, err, "Error creating seth client")
	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2Consumers[0], lc, updatedLabels, testReporter, testType, &testConfig)

		if sethClient.Cfg.IsSimulatedNetwork() {
			l.Info().
				Str("Network Name", sethClient.Cfg.Network.Name).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, sethClient.MustGetRootKeyAddress().Hex(), subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}
	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: testConfig,
		ChainID:    chainID,
		CleanupFn:  cleanupFn,
	}
	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          *vrfv2Config.General.NumberOfSendingKeysToCreate,
		UseVRFOwner:                     true,
		UseTestCoordinator:              true,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}
	testEnv, vrfContracts, vrfKey, _, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2 universe")

	t.Run("vrfv2 and bhs performance test", func(t *testing.T) {
		configCopy := testConfig.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		underfundedSubIDs, consumers, err := vrfv2.SetupSubsAndConsumersForExistingEnv(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2,
			vrfContracts.LinkToken,
			1,
			*vrfv2Config.General.NumberOfSubToCreate,
			configCopy,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		for _, subID := range underfundedSubIDs {
			subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
			require.NoError(t, err, "error getting subscription information for subscription %d", subID)
			vrfcommon.LogSubDetails(l, subscription, strconv.FormatUint(subID, 10), vrfContracts.CoordinatorV2)
		}
		subIDsForCancellingAfterTest = underfundedSubIDs
		l.Debug().Int("Number of Subs", len(underfundedSubIDs)).Msg("Subs involved in the test")
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
				underfundedSubIDs,
				configCopy.VRFv2,
				l,
				sethClient,
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
		latestBlockNumber, err := sethClient.Client.BlockNumber(testcontext.Get(t))
		require.NoError(t, err)
		_, err = actions.WaitForBlockNumberToBe(
			latestBlockNumber+uint64(257),
			sethClient,
			&wgBlockNumberTobe,
			configCopy.VRFv2.General.WaitFor256BlocksTimeout.Duration,
			t,
			l,
		)
		wgBlockNumberTobe.Wait()
		require.NoError(t, err, "error waiting for block number to be")

		metrics, err := consumers[0].GetLoadTestMetrics(testcontext.Get(t))
		require.NoError(t, err)
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)), "Fulfilment count should be 0 since sub is underfunded. Check if the sub is actually funded")

		var subIDsString []uint64
		subIDsString = append(subIDsString, underfundedSubIDs...)
		l.Info().
			Float64("SubscriptionRefundingAmountLink", *configCopy.VRFv2.General.SubscriptionRefundingAmountLink).
			Uints64("SubIDs", subIDsString).
			Msg("Funding Subscriptions with Link and Native Tokens")
		err = vrfv2.FundSubscriptions(
			big.NewFloat(*configCopy.VRFv2.General.SubscriptionRefundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2,
			underfundedSubIDs,
		)
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
