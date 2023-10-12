package loadvrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
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
		requestCount, fulfilmentCount, err := WaitForRequestCountEqualToFulfilmentCount(vrfv2PlusContracts.LoadTestConsumers[0], 30*time.Second, &wg)
		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Final Request/Fulfilment Stats")
		require.NoError(t, err)

		wg.Wait()
	})

}

func WaitForRequestCountEqualToFulfilmentCount(consumer contracts.VRFv2PlusLoadTestConsumer, timeout time.Duration, wg *sync.WaitGroup) (*big.Int, *big.Int, error) {
	metricsChannel := make(chan *contracts.VRFLoadTestMetrics)
	metricsErrorChannel := make(chan error)

	testContext, testCancel := context.WithTimeout(context.Background(), timeout)
	defer testCancel()

	ticker := time.NewTicker(time.Second * 1)
	var metrics *contracts.VRFLoadTestMetrics
	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return metrics.RequestCount, metrics.FulfilmentCount,
				fmt.Errorf("timeout waiting for rand request and fulfilments to be equal AFTER performance test was executed. Request Count: %d, Fulfilment Count: %d",
					metrics.RequestCount.Uint64(), metrics.FulfilmentCount.Uint64())
		case <-ticker.C:
			go getLoadTestMetrics(consumer, metricsChannel, metricsErrorChannel)
		case metrics = <-metricsChannel:
			if metrics.RequestCount.Cmp(metrics.FulfilmentCount) == 0 {
				wg.Done()
				return metrics.RequestCount, metrics.FulfilmentCount, nil
			}
		case err := <-metricsErrorChannel:
			wg.Done()
			return nil, nil, err
		}
	}
}

func getLoadTestMetrics(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	metricsChannel chan *contracts.VRFLoadTestMetrics,
	metricsErrorChannel chan error,
) {
	metrics, err := consumer.GetLoadTestMetrics(context.Background())
	if err != nil {
		metricsErrorChannel <- err
	}
	metricsChannel <- metrics
}
