package loadvrfv2plus

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestVRFV2PlusLoad(t *testing.T) {

	var vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig
	err := envconfig.Process("VRFV2PLUS", &vrfv2PlusConfig)
	require.NoError(t, err)

	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(big.NewFloat(vrfv2PlusConfig.ChainlinkNodeFunding)).
		WithLogWatcher().
		Build()
	require.NoError(t, err, "error creating test env")
	t.Cleanup(func() {
		if err := env.Cleanup(t); err != nil {
			l.Error().Err(err).Msg("Error cleaning up test environment")
		}
	})

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(vrfv2PlusConfig.LinkNativeFeedResponse))
	require.NoError(t, err, "error deploying mock ETH/LINK feed")

	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "error deploying LINK contract")

	vrfv2PlusContracts, subID, vrfv2PlusData, err := vrfv2plus.SetupVRFV2_5Environment(env, vrfv2PlusConfig, linkToken, mockETHLinkFeed, 1)
	require.NoError(t, err, "error setting up VRF v2_5 env")

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	l.Debug().
		Str("Juels Balance", subscription.Balance.String()).
		Str("Native Token Balance", subscription.NativeBalance.String()).
		Str("Subscription ID", subID.String()).
		Str("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")

	labels := map[string]string{
		"branch": "vrfv2Plus_healthcheck",
		"commit": "vrfv2Plus_healthcheck",
	}

	singleFeedConfig := &wasp.Config{
		T:        t,
		LoadType: wasp.RPS,
		GenName:  "gun",
		Gun: SingleFeedGun(
			vrfv2PlusContracts,
			vrfv2PlusData.KeyHash,
			subID,
			vrfv2PlusConfig,
			l),
		Labels:      labels,
		LokiConfig:  wasp.NewEnvLokiConfig(),
		CallTimeout: 2 * time.Minute,
	}

	//multiFeedConfig := &wasp.Config{
	//	T:          t,
	//	LoadType:   wasp.VU,
	//	GenName:    "vu",
	//	VU:         NewJobVolumeVU(cfg.SoakVolume.Pace.Duration(), 1, env.GetAPIs(), env.EVMClient, vrfv2PlusContracts),
	//	Labels:     labels,
	//	LokiConfig: wasp.NewEnvLokiConfig(),
	//}

	MonitorLoadStats(t, vrfv2PlusContracts, labels)

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2plus soak test", func(t *testing.T) {
		singleFeedConfig.Schedule = wasp.Plain(
			vrfv2PlusConfig.RPS,
			vrfv2PlusConfig.TestDuration,
		)
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)

		metrics, err := vrfv2PlusContracts.LoadTestConsumers[0].GetLoadTestMetrics(context.Background())
		require.NoError(t, err)
		if metrics.RequestCount.Cmp(metrics.FulfilmentCount) == 1 {
			fmt.Println("Waiting for all requests to be fulfilled")
			time.Sleep(10 * time.Second)
			newMetrics, err := vrfv2PlusContracts.LoadTestConsumers[0].GetLoadTestMetrics(context.Background())
			require.NoError(t, err)

			//requestId, err := vrfv2PlusContracts.LoadTestConsumers[0].GetLastRequestId(context.Background())
			//require.NoError(t, err)
			//
			//_, err = vrfv2PlusContracts.Coordinator.WaitForRandomWordsFulfilledEvent([]*big.Int{subID}, []*big.Int{requestId}, 10*time.Second)
			//require.NoError(t, err)

			fmt.Println("Request count:", newMetrics.RequestCount, "FulfilmentCount:", newMetrics.FulfilmentCount)

			require.Equal(t, newMetrics.RequestCount, newMetrics.FulfilmentCount)
		}
	})

	//// what are the limits for one "job", figuring out the max/optimal performance params by increasing RPS and varying configuration
	//t.Run("vrfv2Plus load test", func(t *testing.T) {
	//	singleFeedConfig.Schedule = wasp.Steps(
	//		cfg.Load.RPSFrom,
	//		cfg.Load.RPSIncrease,
	//		cfg.Load.RPSSteps,
	//		cfg.Load.Duration.Duration(),
	//	)
	//	_, err = wasp.NewProfile().
	//		Add(wasp.NewGenerator(singleFeedConfig)).
	//		Run(true)
	//	require.NoError(t, err)
	//})
	//
	//// how many "jobs" of the same type we can run at once at a stable load with optimal configuration?
	//t.Run("vrfv2Plus volume soak test", func(t *testing.T) {
	//	multiFeedConfig.Schedule = wasp.Plain(
	//		cfg.SoakVolume.Products,
	//		cfg.SoakVolume.Duration.Duration(),
	//	)
	//	_, err = wasp.NewProfile().
	//		Add(wasp.NewGenerator(multiFeedConfig)).
	//		Run(true)
	//	require.NoError(t, err)
	//})
	//
	//// what are the limits if we add more and more "jobs/products" of the same type, each "job" have a stable RPS we vary only amount of jobs
	//t.Run("vrfv2Plus volume load test", func(t *testing.T) {
	//	multiFeedConfig.Schedule = wasp.Steps(
	//		cfg.LoadVolume.ProductsFrom,
	//		cfg.LoadVolume.ProductsIncrease,
	//		cfg.LoadVolume.ProductsSteps,
	//		cfg.LoadVolume.Duration.Duration(),
	//	)
	//	_, err = wasp.NewProfile().
	//		Add(wasp.NewGenerator(multiFeedConfig)).
	//		Run(true)
	//	require.NoError(t, err)
	//})
}
