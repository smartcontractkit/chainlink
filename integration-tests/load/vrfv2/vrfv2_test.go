package loadvrfv2

import (
	"context"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
)

var (
	env              *test_env.CLClusterTestEnv
	vrfv2Contracts   *vrfv2_actions.VRFV2Contracts
	vrfv2Data        *vrfv2_actions.VRFV2Data
	subIDs           []uint64
	eoaWalletAddress string

	labels = map[string]string{
		"branch": "vrfv2_healthcheck",
		"commit": "vrfv2_healthcheck",
	}

	testType = os.Getenv("TEST_TYPE")
)

func TestVRFV2Performance(t *testing.T) {
	cfg, err := ReadConfig()
	require.NoError(t, err)
	var vrfv2Config vrfv2_config.VRFV2Config
	err = envconfig.Process("VRFV2", &vrfv2Config)
	require.NoError(t, err)

	testReporter := &testreporters.VRFV2TestReporter{}

	SetPerformanceTestConfig(testType, &vrfv2Config, cfg)

	l := logging.GetTestLogger(t)
	//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
	vrfv2Config.MinimumConfirmations = cfg.Common.MinimumConfirmations

	lokiConfig := wasp.NewEnvLokiConfig()
	lc, err := wasp.NewLokiClient(lokiConfig)
	if err != nil {
		l.Error().Err(err).Msg(ErrLokiClient)
		return
	}

	updatedLabels := UpdateLabels(labels, t)

	l.Info().
		Str("Test Type", testType).
		Str("Test Duration", vrfv2Config.TestDuration.Truncate(time.Second).String()).
		Int64("RPS", vrfv2Config.RPS).
		Str("RateLimitUnitDuration", vrfv2Config.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", vrfv2Config.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", vrfv2Config.RandomnessRequestCountPerRequestDeviation).
		Bool("UseExistingEnv", vrfv2Config.UseExistingEnv).
		Msg("Performance Test Configuration")

	if vrfv2Config.UseExistingEnv {
		//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
		vrfv2Config.CoordinatorAddress = cfg.ExistingEnvConfig.CoordinatorAddress
		vrfv2Config.ConsumerAddress = cfg.ExistingEnvConfig.ConsumerAddress
		vrfv2Config.LinkAddress = cfg.ExistingEnvConfig.LinkAddress
		vrfv2Config.SubscriptionFundingAmountLink = cfg.ExistingEnvConfig.SubFunding.SubFundsLink
		vrfv2Config.SubID = cfg.ExistingEnvConfig.SubID
		vrfv2Config.KeyHash = cfg.ExistingEnvConfig.KeyHash

		env, err = test_env.NewCLTestEnvBuilder().
			WithTestInstance(t).
			WithCustomCleanup(
				func() {
					teardown(t, vrfv2Contracts.LoadTestConsumers[0], lc, updatedLabels, testReporter, testType, vrfv2Config)
					if env.EVMClient.NetworkSimulated() {
						l.Info().
							Str("Network Name", env.EVMClient.GetNetworkName()).
							Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
					} else {
						if cfg.Common.CancelSubsAfterTestRun {
							//cancel subs and return funds to sub owner
							cancelSubsAndReturnFunds(subIDs, l)
						}
					}
				}).
			Build()

		require.NoError(t, err, "error creating test env")

		coordinator, err := env.ContractLoader.LoadVRFCoordinatorV2(vrfv2Config.CoordinatorAddress)
		require.NoError(t, err)

		var consumers []contracts.VRFv2LoadTestConsumer
		if cfg.ExistingEnvConfig.CreateFundSubsAndAddConsumers {
			linkToken, err := env.ContractLoader.LoadLINKToken(vrfv2Config.LinkAddress)
			require.NoError(t, err)
			consumers, err = vrfv2_actions.DeployVRFV2Consumers(env.ContractDeployer, coordinator, 1)
			require.NoError(t, err)
			err = env.EVMClient.WaitForEvents()
			require.NoError(t, err, vrfv2_actions.ErrWaitTXsComplete)
			l.Info().
				Str("Coordinator", cfg.ExistingEnvConfig.CoordinatorAddress).
				Int("Number of Subs to create", vrfv2Config.NumberOfSubToCreate).
				Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
			subIDs, err = vrfv2_actions.CreateFundSubsAndAddConsumers(
				env,
				vrfv2Config,
				linkToken,
				coordinator,
				consumers,
				vrfv2Config.NumberOfSubToCreate,
			)
			require.NoError(t, err)
		} else {
			consumer, err := env.ContractLoader.LoadVRFv2LoadTestConsumer(vrfv2Config.ConsumerAddress)
			require.NoError(t, err)
			consumers = append(consumers, consumer)
			subIDs = append(subIDs, vrfv2Config.SubID)
		}

		err = FundNodesIfNeeded(cfg, env.EVMClient, l)
		require.NoError(t, err)

		vrfv2Contracts = &vrfv2_actions.VRFV2Contracts{
			Coordinator:       coordinator,
			LoadTestConsumers: consumers,
			BHS:               nil,
		}

		vrfv2Data = &vrfv2_actions.VRFV2Data{
			VRFV2KeyData: vrfv2_actions.VRFV2KeyData{
				VRFKey:            nil,
				EncodedProvingKey: [2]*big.Int{},
				KeyHash:           common.HexToHash(vrfv2Config.KeyHash),
			},
			VRFJob:            nil,
			PrimaryEthAddress: "",
			ChainID:           nil,
		}

	} else {
		//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
		vrfv2Config.ChainlinkNodeFunding = cfg.NewEnvConfig.NodeSendingKeyFunding
		vrfv2Config.SubscriptionFundingAmountLink = cfg.NewEnvConfig.Funding.SubFundsLink
		env, err = test_env.NewCLTestEnvBuilder().
			WithTestInstance(t).
			WithGeth().
			WithCLNodes(1).
			WithFunding(big.NewFloat(vrfv2Config.ChainlinkNodeFunding)).
			WithCustomCleanup(
				func() {
					teardown(t, vrfv2Contracts.LoadTestConsumers[0], lc, updatedLabels, testReporter, testType, vrfv2Config)

					if env.EVMClient.NetworkSimulated() {
						l.Info().
							Str("Network Name", env.EVMClient.GetNetworkName()).
							Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
					} else {
						if cfg.Common.CancelSubsAfterTestRun {
							//cancel subs and return funds to sub owner
							cancelSubsAndReturnFunds(subIDs, l)
						}
					}
					if err := env.Cleanup(); err != nil {
						l.Error().Err(err).Msg("Error cleaning up test environment")
					}
				}).
			Build()

		require.NoError(t, err, "error creating test env")

		env.ParallelTransactions(true)

		mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(vrfv2Config.LinkNativeFeedResponse))
		require.NoError(t, err, "error deploying mock ETH/LINK feed")

		linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
		require.NoError(t, err, "error deploying LINK contract")

		vrfv2Contracts, subIDs, vrfv2Data, err = vrfv2_actions.SetupVRFV2Environment(
			env,
			vrfv2Config,
			linkToken,
			mockETHLinkFeed,
			//register proving key against EOA address in order to return funds to this address
			env.EVMClient.GetDefaultWallet().Address(),
			0,
			1,
			vrfv2Config.NumberOfSubToCreate,
			l,
		)
		require.NoError(t, err, "error setting up VRF v2 env")
	}
	eoaWalletAddress = env.EVMClient.GetDefaultWallet().Address()

	l.Debug().Int("Number of Subs", len(subIDs)).Msg("Subs involved in the test")
	for _, subID := range subIDs {
		subscription, err := vrfv2Contracts.Coordinator.GetSubscription(context.Background(), subID)
		require.NoError(t, err, "error getting subscription information for subscription %d", subID)
		vrfv2_actions.LogSubDetails(l, subscription, subID, vrfv2Contracts.Coordinator)
	}

	singleFeedConfig := &wasp.Config{
		T:                     t,
		LoadType:              wasp.RPS,
		GenName:               "gun",
		RateLimitUnitDuration: vrfv2Config.RateLimitUnitDuration,
		Gun: NewSingleHashGun(
			vrfv2Contracts,
			vrfv2Data.KeyHash,
			subIDs,
			vrfv2Config,
			l,
		),
		Labels:      labels,
		LokiConfig:  lokiConfig,
		CallTimeout: 2 * time.Minute,
	}
	require.Len(t, vrfv2Contracts.LoadTestConsumers, 1, "only one consumer should be created for Load Test")
	consumer := vrfv2Contracts.LoadTestConsumers[0]
	err = consumer.ResetMetrics()
	require.NoError(t, err)
	MonitorLoadStats(lc, consumer, updatedLabels)

	// is our "job" stable at all, no memory leaks, no flaking performance under some RPS?
	t.Run("vrfv2 performance test", func(t *testing.T) {

		singleFeedConfig.Schedule = wasp.Plain(
			vrfv2Config.RPS,
			vrfv2Config.TestDuration,
		)
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(singleFeedConfig)).
			Run(true)
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)
		//todo - timeout should be configurable depending on the perf test type
		requestCount, fulfilmentCount, err := vrfv2_actions.WaitForRequestCountEqualToFulfilmentCount(consumer, 2*time.Minute, &wg)
		require.NoError(t, err)
		wg.Wait()

		l.Info().
			Interface("Request Count", requestCount).
			Interface("Fulfilment Count", fulfilmentCount).
			Msg("Final Request/Fulfilment Stats")
	})

}

func cancelSubsAndReturnFunds(subIDs []uint64, l zerolog.Logger) {
	for _, subID := range subIDs {
		l.Info().
			Uint64("Returning funds from SubID", subID).
			Str("Returning funds to", eoaWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		pendingRequestsExist, err := vrfv2Contracts.Coordinator.PendingRequestsExist(context.Background(), subID)
		if err != nil {
			l.Error().Err(err).Msg("Error checking if pending requests exist")
		}
		if !pendingRequestsExist {
			_, err := vrfv2Contracts.Coordinator.CancelSubscription(subID, common.HexToAddress(eoaWalletAddress))
			if err != nil {
				l.Error().Err(err).Msg("Error canceling subscription")
			}
		} else {
			l.Error().Uint64("Sub ID", subID).Msg("Pending requests exist for subscription, cannot cancel subscription and return funds")
		}
	}
}

func FundNodesIfNeeded(cfg *PerformanceConfig, client blockchain.EVMClient, l zerolog.Logger) error {
	if cfg.ExistingEnvConfig.NodeSendingKeyFundingMin > 0 {
		for _, sendingKey := range cfg.ExistingEnvConfig.NodeSendingKeys {
			address := common.HexToAddress(sendingKey)
			sendingKeyBalance, err := client.BalanceAt(context.Background(), address)
			if err != nil {
				return err
			}
			fundingAtLeast := conversions.EtherToWei(big.NewFloat(cfg.ExistingEnvConfig.NodeSendingKeyFundingMin))
			fundingToSendWei := new(big.Int).Sub(fundingAtLeast, sendingKeyBalance)
			fundingToSendEth := conversions.WeiToEther(fundingToSendWei)
			if fundingToSendWei.Cmp(big.NewInt(0)) == 1 {
				l.Info().
					Str("Sending Key", sendingKey).
					Str("Sending Key Current Balance", sendingKeyBalance.String()).
					Str("Should have at least", fundingAtLeast.String()).
					Str("Funding Amount in ETH", fundingToSendEth.String()).
					Msg("Funding Node's Sending Key")
				err := actions.FundAddress(client, sendingKey, fundingToSendEth)
				if err != nil {
					return err
				}
			} else {
				l.Info().
					Str("Sending Key", sendingKey).
					Str("Sending Key Current Balance", sendingKeyBalance.String()).
					Str("Should have at least", fundingAtLeast.String()).
					Msg("Skipping Node's Sending Key funding as it has enough funds")
			}
		}
	}
	return nil
}

func teardown(
	t *testing.T,
	consumer contracts.VRFv2LoadTestConsumer,
	lc *wasp.LokiClient,
	updatedLabels map[string]string,
	testReporter *testreporters.VRFV2TestReporter,
	testType string,
	vrfv2Config vrfv2_config.VRFV2Config,
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
		vrfv2Config,
	)

	// send Slack notification
	err := testReporter.SendSlackNotification(t, nil)
	if err != nil {
		log.Warn().Err(err).Msg("Error sending Slack notification")
	}
}
