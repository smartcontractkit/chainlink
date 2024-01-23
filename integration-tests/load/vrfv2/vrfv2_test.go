package loadvrfv2

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
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
		Str("Test Type", string(testType)).
		Str("Test Duration", vrfv2Config.Performance.TestDuration.Duration.Truncate(time.Second).String()).
		Int64("RPS", *vrfv2Config.Performance.RPS).
		Str("RateLimitUnitDuration", vrfv2Config.Performance.RateLimitUnitDuration.String()).
		Uint16("RandomnessRequestCountPerRequest", *vrfv2Config.General.RandomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", *vrfv2Config.General.RandomnessRequestCountPerRequestDeviation).
		Bool("UseExistingEnv", *vrfv2Config.Performance.UseExistingEnv).
		Msg("Performance Test Configuration")

	if *vrfv2Config.Performance.UseExistingEnv {
		//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
		cfg := testConfig.VRFv2

		vrfv2Config.Performance.CoordinatorAddress = cfg.ExistingEnvConfig.CoordinatorAddress
		vrfv2Config.Performance.ConsumerAddress = cfg.ExistingEnvConfig.ConsumerAddress
		vrfv2Config.Performance.LinkAddress = cfg.ExistingEnvConfig.LinkAddress
		vrfv2Config.General.SubscriptionFundingAmountLink = cfg.ExistingEnvConfig.SubFunding.SubFundsLink
		vrfv2Config.Performance.SubID = cfg.ExistingEnvConfig.SubID
		vrfv2Config.Performance.KeyHash = cfg.ExistingEnvConfig.KeyHash

		env, err = test_env.NewCLTestEnvBuilder().
			WithTestInstance(t).
			WithTestConfig(&testConfig).
			WithCustomCleanup(
				func() {
					teardown(t, vrfv2Contracts.LoadTestConsumers[0], lc, updatedLabels, testReporter, string(testType), &testConfig)
					if env.EVMClient.NetworkSimulated() {
						l.Info().
							Str("Network Name", env.EVMClient.GetNetworkName()).
							Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
					} else {
						if *vrfv2Config.Common.CancelSubsAfterTestRun {
							//cancel subs and return funds to sub owner
							cancelSubsAndReturnFunds(subIDs, l)
						}
					}
				}).
			Build()

		require.NoError(t, err, "error creating test env")

		coordinator, err := env.ContractLoader.LoadVRFCoordinatorV2(*vrfv2Config.Performance.CoordinatorAddress)
		require.NoError(t, err)

		var consumers []contracts.VRFv2LoadTestConsumer
		if *cfg.ExistingEnvConfig.CreateFundSubsAndAddConsumers {
			linkToken, err := env.ContractLoader.LoadLINKToken(*vrfv2Config.Performance.LinkAddress)
			require.NoError(t, err)
			consumers, err = vrfv2_actions.DeployVRFV2Consumers(env.ContractDeployer, coordinator.Address(), 1)
			require.NoError(t, err)
			err = env.EVMClient.WaitForEvents()
			require.NoError(t, err, vrfv2_actions.ErrWaitTXsComplete)
			l.Info().
				Str("Coordinator", *cfg.ExistingEnvConfig.CoordinatorAddress).
				Int("Number of Subs to create", *vrfv2Config.General.NumberOfSubToCreate).
				Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
			subIDs, err = vrfv2_actions.CreateFundSubsAndAddConsumers(
				env,
				big.NewFloat(*cfg.General.SubscriptionFundingAmountLink),
				linkToken,
				coordinator,
				consumers,
				*vrfv2Config.General.NumberOfSubToCreate,
			)
			require.NoError(t, err)
		} else {
			consumer, err := env.ContractLoader.LoadVRFv2LoadTestConsumer(*vrfv2Config.Performance.ConsumerAddress)
			require.NoError(t, err)
			consumers = append(consumers, consumer)
			subIDs = append(subIDs, *vrfv2Config.Performance.SubID)
		}

		err = FundNodesIfNeeded(&testConfig, env.EVMClient, l)
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
				KeyHash:           common.HexToHash(*vrfv2Config.Performance.KeyHash),
			},
			VRFJob:            nil,
			PrimaryEthAddress: "",
			ChainID:           nil,
		}

	} else {
		//todo: temporary solution with envconfig and toml config until VRF-662 is implemented
		testConfig.Common.ChainlinkNodeFunding = testConfig.VRFv2.NewEnvConfig.NodeSendingKeyFunding
		vrfv2Config.General.SubscriptionFundingAmountLink = testConfig.VRFv2.NewEnvConfig.Funding.SubFundsLink

		network, err := actions.EthereumNetworkConfigFromConfig(l, &testConfig)
		require.NoError(t, err, "Error building ethereum network config")
		env, err = test_env.NewCLTestEnvBuilder().
			WithTestInstance(t).
			WithTestConfig(&testConfig).
			WithPrivateEthereumNetwork(network).
			WithCLNodes(1).
			WithFunding(big.NewFloat(*testConfig.Common.ChainlinkNodeFunding)).
			WithCustomCleanup(
				func() {
					teardown(t, vrfv2Contracts.LoadTestConsumers[0], lc, updatedLabels, testReporter, string(testType), &testConfig)

					if env.EVMClient.NetworkSimulated() {
						l.Info().
							Str("Network Name", env.EVMClient.GetNetworkName()).
							Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
					} else {
						if *testConfig.VRFv2.Common.CancelSubsAfterTestRun {
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

		mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*vrfv2Config.General.LinkNativeFeedResponse))
		require.NoError(t, err, "error deploying mock ETH/LINK feed")

		linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
		require.NoError(t, err, "error deploying LINK contract")

		useVRFOwner := true
		useTestCoordinator := true

		vrfv2Contracts, subIDs, vrfv2Data, err = vrfv2_actions.SetupVRFV2Environment(
			env,
			&testConfig,
			useVRFOwner,
			useTestCoordinator,
			linkToken,
			mockETHLinkFeed,
			//register proving key against EOA address in order to return funds to this address
			env.EVMClient.GetDefaultWallet().Address(),
			0,
			1,
			*vrfv2Config.General.NumberOfSubToCreate,
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
		RateLimitUnitDuration: vrfv2Config.Performance.RateLimitUnitDuration.Duration,
		Gun: NewSingleHashGun(
			vrfv2Contracts,
			vrfv2Data.KeyHash,
			subIDs,
			&testConfig,
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

func FundNodesIfNeeded(vrfv2TestConfig tc.VRFv2TestConfig, client blockchain.EVMClient, l zerolog.Logger) error {
	cfg := vrfv2TestConfig.GetVRFv2Config()
	if cfg.ExistingEnvConfig.NodeSendingKeyFundingMin != nil && *cfg.ExistingEnvConfig.NodeSendingKeyFundingMin > 0 {
		for _, sendingKey := range cfg.ExistingEnvConfig.NodeSendingKeys {
			address := common.HexToAddress(sendingKey)
			sendingKeyBalance, err := client.BalanceAt(context.Background(), address)
			if err != nil {
				return err
			}
			fundingAtLeast := conversions.EtherToWei(big.NewFloat(*cfg.ExistingEnvConfig.NodeSendingKeyFundingMin))
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
