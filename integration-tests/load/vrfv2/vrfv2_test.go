package loadvrfv2

import (
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	env              *test_env.CLClusterTestEnv
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

	cleanupFn := func() {
		teardown(t, vrfContracts.VRFV2Consumer[0], lc, updatedLabels, testReporter, testType, &testConfig)

		if env.EVMClient.NetworkSimulated() {
			l.Info().
				Str("Network Name", env.EVMClient.GetNetworkName()).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {
			if *vrfv2Config.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, eoaWalletAddress, subIDs, l)
			}
		}
		if !*vrfv2Config.General.UseExistingEnv {
			if err := env.Cleanup(); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:          []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate: *vrfv2Config.General.NumberOfSendingKeysToCreate,
		NumberOfConsumers:      1,
		NumberOfSubToCreate:    *vrfv2Config.General.NumberOfSubToCreate,
		UseVRFOwner:            true,
		UseTestCoordinator:     true,
	}

	env, vrfContracts, subIDs, vrfKey, _, err = vrfv2.SetupVRFV2Universe(testcontext.Get(t), t, testConfig, cleanupFn, newEnvConfig, l)
	require.NoError(t, err)

	eoaWalletAddress = env.EVMClient.GetDefaultWallet().Address()

	l.Debug().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")
	for _, subID := range subIDs {
		subscription, err := vrfContracts.CoordinatorV2.GetSubscription(testcontext.Get(t), subID)
		require.NoError(t, err, "error getting subscription information for subscription %d", subID)
		vrfv2.LogSubDetails(l, subscription, subID, vrfContracts.CoordinatorV2)
	}
	singleFeedConfig := &wasp.Config{
		T:                     t,
		LoadType:              wasp.RPS,
		GenName:               "gun",
		RateLimitUnitDuration: vrfv2Config.Performance.RateLimitUnitDuration.Duration,
		Gun: NewSingleHashGun(
			vrfContracts,
			vrfKey.KeyHash,
			subIDs,
			&testConfig,
			l,
		),
		Labels:      labels,
		LokiConfig:  lokiConfig,
		CallTimeout: 2 * time.Minute,
	}
	require.Len(t, vrfContracts.VRFV2Consumer, 1, "only one consumer should be created for Load Test")
	consumer := vrfContracts.VRFV2Consumer[0]
	err = consumer.ResetMetrics()
	require.NoError(t, err)
	MonitorLoadStats(lc, consumer, updatedLabels)

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2 performance test", func(t *testing.T) {

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
	consumer contracts.VRFv2LoadTestConsumer,
	lc *wasp.LokiClient,
	updatedLabels map[string]string,
	testReporter *testreporters.VRFV2TestReporter,
	testType string,
	testConfig *tc.TestConfig,
) {
	//send final results to Loki
	metrics := GetLoadTestMetrics(consumer)
	SendMetricsToLoki(metrics, lc, updatedLabels)
	//set report data for Slack notification
	testReporter.SetReportData(
		testType,
		metrics.RequestCount,
		metrics.FulfilmentCount,
		metrics.AverageFulfillmentInMillions,
		metrics.SlowestFulfillment,
		metrics.FastestFulfillment,
		testConfig,
	)

	// send Slack notification
	err := testReporter.SendSlackNotification(t, nil)
	if err != nil {
		log.Warn().Err(err).Msg("Error sending Slack notification")
	}
}
