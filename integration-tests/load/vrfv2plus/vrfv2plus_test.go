package loadvrfv2plus

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kelseyhightower/envconfig"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestVRFV2PlusLoad(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	var vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig
	err = envconfig.Process("VRFV2PLUS", &vrfv2PlusConfig)
	require.NoError(t, err)

	l := logging.GetTestLogger(t)
	//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
	vrfv2PlusConfig.ChainlinkNodeFunding = cfg.Common.NodeFunds
	vrfv2PlusConfig.IsNativePayment = cfg.Common.IsNativePayment
	vrfv2PlusConfig.MinimumConfirmations = cfg.Common.MinimumConfirmations
	vrfv2PlusConfig.SubscriptionFundingAmountLink = cfg.Common.Funding.SubFundsLink
	vrfv2PlusConfig.SubscriptionFundingAmountNative = cfg.Common.Funding.SubFundsNative

	var env *test_env.CLClusterTestEnv
	var vrfv2PlusContracts *vrfv2plus.VRFV2_5Contracts
	var subID *big.Int
	var vrfv2PlusData *vrfv2plus.VRFV2PlusData
	if vrfv2PlusConfig.UseExistingEnv {
		//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
		vrfv2PlusConfig.CoordinatorAddress = cfg.ExistingEnvConfig.CoordinatorAddress
		vrfv2PlusConfig.ConsumerAddress = cfg.ExistingEnvConfig.ConsumerAddress
		vrfv2PlusConfig.SubID = cfg.ExistingEnvConfig.SubID
		vrfv2PlusConfig.KeyHash = cfg.ExistingEnvConfig.KeyHash

		env, err = test_env.NewCLTestEnvBuilder().
			WithTestLogger(t).
			Build()

		require.NoError(t, err, "error creating test env")

		coordinator, err := env.ContractLoader.LoadVRFCoordinatorV2_5(vrfv2PlusConfig.CoordinatorAddress)
		require.NoError(t, err)

		consumer, err := env.ContractLoader.LoadVRFv2PlusLoadTestConsumer(vrfv2PlusConfig.ConsumerAddress)
		require.NoError(t, err)

		vrfv2PlusContracts = &vrfv2plus.VRFV2_5Contracts{
			Coordinator:       coordinator,
			LoadTestConsumers: []contracts.VRFv2PlusLoadTestConsumer{consumer},
			BHS:               nil,
		}
		var ok bool
		subID, ok = new(big.Int).SetString(vrfv2PlusConfig.SubID, 10)
		require.True(t, ok)

		vrfv2PlusData = &vrfv2plus.VRFV2PlusData{
			VRFV2PlusKeyData: vrfv2plus.VRFV2PlusKeyData{
				VRFKey:            nil,
				EncodedProvingKey: [2]*big.Int{},
				KeyHash:           common.HexToHash(vrfv2PlusConfig.KeyHash),
			},
			VRFJob:            nil,
			PrimaryEthAddress: "",
			ChainID:           nil,
		}
	} else {
		env, err = test_env.NewCLTestEnvBuilder().
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

		vrfv2PlusContracts, subID, vrfv2PlusData, err = vrfv2plus.SetupVRFV2_5Environment(env, vrfv2PlusConfig, linkToken, mockETHLinkFeed, 1)
		require.NoError(t, err, "error setting up VRF v2_5 env")
	}

	subscription, err := vrfv2PlusContracts.Coordinator.GetSubscription(context.Background(), subID)
	require.NoError(t, err, "error getting subscription information")

	vrfv2plus.LogSubDetails(l, subscription, subID, vrfv2PlusContracts.Coordinator)

	labels := map[string]string{
		"branch": "vrfv2Plus_healthcheck",
		"commit": "vrfv2Plus_healthcheck",
	}

	l.Info().
		Str("Test Duration", vrfv2PlusConfig.TestDuration.Truncate(time.Second).String()).
		Int64("RPS", vrfv2PlusConfig.RPS).
		Str("RateLimitUnitDuration", vrfv2PlusConfig.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", vrfv2PlusConfig.RandomnessRequestCountPerRequest).
		Bool("UseExistingEnv", vrfv2PlusConfig.UseExistingEnv).
		Msg("Load Test Configs")

	lokiConfig := wasp.NewEnvLokiConfig()
	lc, err := wasp.NewLokiClient(lokiConfig)
	if err != nil {
		l.Error().Err(err).Msg(ErrLokiClient)
		return
	}

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
			l,
		),
		Labels:      labels,
		LokiConfig:  lokiConfig,
		CallTimeout: 2 * time.Minute,
	}
	consumer := vrfv2PlusContracts.LoadTestConsumers[0]
	err = consumer.ResetMetrics()
	require.NoError(t, err)
	updatedLabels := UpdateLabels(labels, t)
	MonitorLoadStats(lc, vrfv2PlusContracts, updatedLabels)

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2plus soak test", func(t *testing.T) {

		singleFeedConfig.Schedule = wasp.Plain(
			vrfv2PlusConfig.RPS,
			vrfv2PlusConfig.TestDuration,
		)
		_, err = wasp.NewProfile().
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
		//send final results
		SendLoadTestMetricsToLoki(vrfv2PlusContracts, lc, updatedLabels)
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
				ticker.Stop()
				wg.Done()
				return metrics.RequestCount, metrics.FulfilmentCount, nil
			}
		case err := <-metricsErrorChannel:
			ticker.Stop()
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
