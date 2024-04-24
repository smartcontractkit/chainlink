package loadvrfv2plus

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	labels = map[string]string{
		"branch": "vrfv2Plus_healthcheck",
		"commit": "vrfv2Plus_healthcheck",
	}
)

func TestVRFV2PlusPerformance(t *testing.T) {
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
	)
	l := logging.GetTestLogger(t)
	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2Plus)
	require.NoError(t, err)
	cfgl := testConfig.Logging.Loki

	vrfv2PlusConfig := testConfig.VRFv2Plus
	testReporter := &testreporters.VRFV2PlusTestReporter{}

	lokiConfig := wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken)
	lc, err := wasp.NewLokiClient(lokiConfig)
	if err != nil {
		l.Error().Err(err).Msg(ErrLokiClient)
		return
	}
	chainID := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0].ChainID
	updatedLabels := UpdateLabels(labels, t)

	l.Info().
		Str("Test Type", testType).
		Str("Test Duration", vrfv2PlusConfig.Performance.TestDuration.Duration.Truncate(time.Second).String()).
		Int64("RPS", *vrfv2PlusConfig.Performance.RPS).
		Str("RateLimitUnitDuration", vrfv2PlusConfig.Performance.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", *vrfv2PlusConfig.General.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", *vrfv2PlusConfig.General.RandomnessRequestCountPerRequestDeviation).
		Bool("UseExistingEnv", *vrfv2PlusConfig.General.UseExistingEnv).
		Msg("Performance Test Configuration")

	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2PlusConsumer[0], lc, updatedLabels, testReporter, testType, &testConfig)

		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "error getting EVM client")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *testConfig.VRFv2Plus.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*testConfig.VRFv2Plus.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: *vrfv2PlusConfig.General.NumberOfSendingKeysToCreate,
	}

	testEnv, vrfContracts, vrfKey, _, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, testConfig, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "Getting EVM client shouldn't fail")

	var consumers []contracts.VRFv2PlusLoadTestConsumer
	subIDs, consumers, err := vrfv2plus.SetupSubsAndConsumersForExistingEnv(
		testEnv,
		chainID,
		vrfContracts.CoordinatorV2Plus,
		vrfContracts.LinkToken,
		1,
		*vrfv2PlusConfig.General.NumberOfSubToCreate,
		testConfig,
		l,
	)
	require.NoError(t, err, "error setting up new consumers and subs")
	for _, subID := range subIDs {
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information for subscription %s", subID.String())
		vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	}
	subIDsForCancellingAfterTest = subIDs
	l.Info().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")

	vrfContracts.VRFV2PlusConsumer = consumers
	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2plus performance test", func(t *testing.T) {
		require.Len(t, vrfContracts.VRFV2PlusConsumer, 1, "only one consumer should be created for Load Test")
		consumer := vrfContracts.VRFV2PlusConsumer[0]
		err = consumer.ResetMetrics()
		require.NoError(t, err)
		MonitorLoadStats(testcontext.Get(t), lc, consumer, updatedLabels)

		singleFeedConfig := &wasp.Config{
			T:                     t,
			LoadType:              wasp.RPS,
			GenName:               "gun",
			RateLimitUnitDuration: vrfv2PlusConfig.Performance.RateLimitUnitDuration.Duration,
			Gun: NewSingleHashGun(
				vrfContracts,
				vrfKey.KeyHash,
				subIDs,
				vrfv2PlusConfig,
				l,
			),
			Labels:      labels,
			LokiConfig:  wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
			CallTimeout: 2 * time.Minute,
		}

		singleFeedConfig.Schedule = wasp.Plain(
			*vrfv2PlusConfig.Performance.RPS,
			vrfv2PlusConfig.Performance.TestDuration.Duration,
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)
		//todo - timeout should be configurable depending on the perf test type
		requestCount, fulfilmentCount, err := vrfcommon.WaitForRequestCountEqualToFulfilmentCount(testcontext.Get(t), consumer, 2*time.Minute, &wg)
		require.NoError(t, err)
		wg.Wait()

		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Final Request/Fulfilment Stats")
	})
}

func TestVRFV2PlusBHSPerformance(t *testing.T) {
	var (
		testEnv                      *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		defaultWalletAddress         string
		vrfKey                       *vrfcommon.VRFKeyData
	)
	l := logging.GetTestLogger(t)

	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2Plus)
	require.NoError(t, err)
	vrfv2PlusConfig := testConfig.VRFv2Plus
	testReporter := &testreporters.VRFV2PlusTestReporter{}
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
		Str("Test Duration", vrfv2PlusConfig.Performance.TestDuration.Duration.Truncate(time.Second).String()).
		Int64("RPS", *vrfv2PlusConfig.Performance.RPS).
		Str("RateLimitUnitDuration", vrfv2PlusConfig.Performance.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", *vrfv2PlusConfig.General.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", *vrfv2PlusConfig.General.RandomnessRequestCountPerRequestDeviation).
		Bool("UseExistingEnv", *vrfv2PlusConfig.General.UseExistingEnv).
		Msg("Performance Test Configuration")

	chainID := networks.MustGetSelectedNetworkConfig(testConfig.GetNetworkConfig())[0].ChainID

	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2PlusConsumer[0], lc, updatedLabels, testReporter, testType, &testConfig)

		evmClient, err := testEnv.GetEVMClient(chainID)
		require.NoError(t, err, "error getting EVM client")

		if evmClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", evmClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *testConfig.VRFv2Plus.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, defaultWalletAddress, subIDsForCancellingAfterTest, l)
			}
		}
		if !*testConfig.VRFv2Plus.General.UseExistingEnv {
			if err := testEnv.Cleanup(test_env.CleanupOpts{}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: *vrfv2PlusConfig.General.NumberOfSendingKeysToCreate,
	}

	testEnv, vrfContracts, vrfKey, _, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, testConfig, chainID, cleanupFn, newEnvConfig, l)
	require.NoError(t, err, "error setting up VRFV2Plus universe")

	evmClient, err := testEnv.GetEVMClient(chainID)
	require.NoError(t, err, "error getting EVM client")

	defaultWalletAddress = evmClient.GetDefaultWallet().Address()

	t.Run("vrfv2plus and bhs performance test", func(t *testing.T) {
		configCopy := testConfig.MustCopy().(tc.TestConfig)
		//Underfund Subscription
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
		configCopy.VRFv2Plus.General.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))

		underfundedSubIDs, consumers, err := vrfv2plus.SetupSubsAndConsumersForExistingEnv(
			testEnv,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			vrfContracts.LinkToken,
			1,
			*vrfv2PlusConfig.General.NumberOfSubToCreate,
			configCopy,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs for Load Test")
		for _, subID := range underfundedSubIDs {
			subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
			require.NoError(t, err, "error getting subscription information for subscription %s", subID.String())
			vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
		}
		subIDsForCancellingAfterTest = underfundedSubIDs
		l.Debug().Int("Number of Subs", len(underfundedSubIDs)).Msg("Subs involved in the test")
		vrfContracts.VRFV2PlusConsumer = consumers
		require.Len(t, vrfContracts.VRFV2PlusConsumer, 1, "only one consumer should be created for Load Test")
		consumer := vrfContracts.VRFV2PlusConsumer[0]
		err = consumer.ResetMetrics()
		require.NoError(t, err)
		MonitorLoadStats(testcontext.Get(t), lc, consumer, updatedLabels)

		singleFeedConfig := &wasp.Config{
			T:                     t,
			LoadType:              wasp.RPS,
			GenName:               "gun",
			RateLimitUnitDuration: configCopy.VRFv2Plus.Performance.BHSTestRateLimitUnitDuration.Duration,
			Gun: NewBHSTestGun(
				vrfContracts,
				vrfKey.KeyHash,
				underfundedSubIDs,
				configCopy.VRFv2Plus,
				l,
			),
			Labels:      labels,
			LokiConfig:  lokiConfig,
			CallTimeout: 2 * time.Minute,
		}

		singleFeedConfig.Schedule = wasp.Plain(
			*configCopy.VRFv2Plus.Performance.BHSTestRPS,
			configCopy.VRFv2Plus.Performance.BHSTestDuration.Duration,
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
		_, err = actions.WaitForBlockNumberToBe(latestBlockNumber+uint64(256), evmClient, &wgBlockNumberTobe, configCopy.VRFv2Plus.General.WaitFor256BlocksTimeout.Duration, t)
		wgBlockNumberTobe.Wait()
		require.NoError(t, err, "error waiting for block number to be")

		metrics, err := consumers[0].GetLoadTestMetrics(testcontext.Get(t))
		require.NoError(t, err)
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)), "Fulfilment count should be 0 since sub is underfunded. Check if the sub is actually funded")

		var subIDsString []string
		for _, subID := range underfundedSubIDs {
			subIDsString = append(subIDsString, subID.String())
		}

		l.Info().
			Float64("SubscriptionRefundingAmountNative", *configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative).
			Float64("SubscriptionRefundingAmountLink", *configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink).
			Strs("SubIDs", subIDsString).
			Msg("Funding Subscriptions with Link and Native Tokens")
		err = vrfv2plus.FundSubscriptions(
			testEnv,
			chainID,
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountNative),
			big.NewFloat(*configCopy.VRFv2Plus.General.SubscriptionRefundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			underfundedSubIDs,
		)
		require.NoError(t, err, "error funding subscriptions")

		var wgAllRequestsFulfilled sync.WaitGroup
		wgAllRequestsFulfilled.Add(1)
		requestCount, fulfilmentCount, err := vrfcommon.WaitForRequestCountEqualToFulfilmentCount(testcontext.Get(t), consumer, 2*time.Minute, &wgAllRequestsFulfilled)
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
	consumer contracts.VRFv2PlusLoadTestConsumer,
	lc *wasp.LokiClient,
	updatedLabels map[string]string,
	testReporter *testreporters.VRFV2PlusTestReporter,
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
			RequestCount:                         metrics.RequestCount,
			FulfilmentCount:                      metrics.FulfilmentCount,
			AverageFulfillmentInMillions:         metrics.AverageFulfillmentInMillions,
			SlowestFulfillment:                   metrics.SlowestFulfillment,
			FastestFulfillment:                   metrics.FastestFulfillment,
			P90FulfillmentBlockTime:              metrics.P90FulfillmentBlockTime,
			P95FulfillmentBlockTime:              metrics.P95FulfillmentBlockTime,
			AverageResponseTimeInSecondsMillions: metrics.AverageResponseTimeInSecondsMillions,
			SlowestResponseTimeInSeconds:         metrics.SlowestResponseTimeInSeconds,
			FastestResponseTimeInSeconds:         metrics.FastestResponseTimeInSeconds,
		},
		testConfig,
	)

	// send Slack notification
	err := testReporter.SendSlackNotification(t, nil, testConfig)
	if err != nil {
		log.Warn().Err(err).Msg("Error sending Slack notification")
	}
}
