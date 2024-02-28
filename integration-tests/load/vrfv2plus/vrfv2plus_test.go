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
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	env              *test_env.CLClusterTestEnv
	vrfContracts     *vrfcommon.VRFContracts
	vrfKey           *vrfcommon.VRFKeyData
	subIDs           []*big.Int
	eoaWalletAddress string

	labels = map[string]string{
		"branch": "vrfv2Plus_healthcheck",
		"commit": "vrfv2Plus_healthcheck",
	}
)

func TestVRFV2PlusPerformance(t *testing.T) {
	l := logging.GetTestLogger(t)

	testType, err := tc.GetConfigurationNameFromEnv()
	require.NoError(t, err)
	testConfig, err := tc.GetConfig(testType, tc.VRFv2Plus)
	require.NoError(t, err)
	cfgl := testConfig.Logging.Loki

	vrfv2PlusConfig := testConfig.VRFv2Plus
	testReporter := &testreporters.VRFV2PlusTestReporter{}

	lc, err := wasp.NewLokiClient(wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken))
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

	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2PlusConsumer[0], lc, updatedLabels, testReporter, testType, &testConfig)

		if env.EVMClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", env.EVMClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *testConfig.VRFv2Plus.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, eoaWalletAddress, subIDs, l)
			}
		}
		if !*testConfig.VRFv2Plus.General.UseExistingEnv {
			if err := env.Cleanup(); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: *vrfv2PlusConfig.General.NumberOfSendingKeysToCreate,
		NumberOfConsumers:      1,
		NumberOfSubToCreate:    *vrfv2PlusConfig.General.NumberOfSubToCreate,
	}

	env, vrfContracts, subIDs, vrfKey, _, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, testConfig, cleanupFn, newEnvConfig, l)
	require.NoError(t, err)
	eoaWalletAddress = env.EVMClient.GetDefaultWallet().Address()

	l.Debug().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")
	for _, subID := range subIDs {
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information for subscription %s", subID.String())
		vrfv2plus.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2Plus)
	}

	singleFeedConfig := &wasp.Config{
		T:                     t,
		LoadType:              wasp.RPS,
		GenName:               "gun",
		RateLimitUnitDuration: vrfv2PlusConfig.Performance.RateLimitUnitDuration.Duration,
		Gun: NewSingleHashGun(
			vrfContracts,
			vrfKey.KeyHash,
			subIDs,
			&testConfig,
			l,
		),
		Labels:      labels,
		LokiConfig:  wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
		CallTimeout: 2 * time.Minute,
	}
	require.Len(t, vrfContracts.VRFV2PlusConsumer, 1, "only one consumer should be created for Load Test")
	consumer := vrfContracts.VRFV2PlusConsumer[0]
	err = consumer.ResetMetrics()
	require.NoError(t, err)
	MonitorLoadStats(lc, consumer, updatedLabels)

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2plus performance test", func(t *testing.T) {
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
	metrics := GetLoadTestMetrics(consumer)
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
