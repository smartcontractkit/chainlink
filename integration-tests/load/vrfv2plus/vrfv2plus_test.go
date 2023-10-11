package loadvrfv2plus

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync"
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

	l.Info().
		Str("Test Duration", vrfv2PlusConfig.TestDuration.Truncate(time.Second).String()).
		Int64("RPS", vrfv2PlusConfig.RPS).
		Uint16("RandomnessRequestCountPerRequest", vrfv2PlusConfig.RandomnessRequestCountPerRequest).
		Msg("Load Test Configs")

	singleFeedConfig := &wasp.Config{
		T:                     t,
		LoadType:              wasp.RPS,
		GenName:               "gun",
		RateLimitUnitDuration: vrfv2PlusConfig.RateLimitUnitDuration,
		Gun: NewSingleHashGun(
			vrfv2PlusContracts,
			vrfv2PlusData.KeyHash,
			subID,
			vrfv2PlusConfig,
			l),
		Labels:      labels,
		LokiConfig:  wasp.NewEnvLokiConfig(),
		CallTimeout: 2 * time.Minute,
	}

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

		var wg sync.WaitGroup

		wg.Add(1)
		requestCount, fulfilmentCount, err := vrfv2plus.WaitForRequestCountEqualToFulfilmentCount(vrfv2PlusContracts.LoadTestConsumers[0], 30*time.Second, &wg)
		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Final Request/Fulfilment Stats")
		require.NoError(t, err)

		wg.Wait()
	})

}
