package soak

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestVRFV2Soak(t *testing.T) {
	linkEthFeedResponse := big.NewInt(1e18)
	minimumConfirmations := 3
	subID := uint64(1)
	numberOfWords := uint32(3)
	maxGasPriceGWei := 1000
	callbackGasLimit := uint32(1000000)
	waitForRandRequestStatusToBeFulfilledTimeout := time.Second * 40
	l := utils.GetTestLogger(t)

	var testInputs testsetups.VRFV2SoakTestInputs
	err := envconfig.Process("VRFV2", &testInputs)
	require.NoError(t, err, "Error reading VRFV2 soak test inputs")
	testInputs.SetForRemoteRunner()
	testNetwork := networks.SelectedNetwork // Environment currently being used to soak test on

	testEnvironment := setupVRFV2Environment(t, testNetwork, config.BaseVRFV2NetworkDetailTomlConfig, "")
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err)
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err)
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err)
	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err)
	bhs, err := contractDeployer.DeployBlockhashStore()
	require.NoError(t, err)
	mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
	require.NoError(t, err)
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkToken.Address(), bhs.Address(), mf.Address())
	require.NoError(t, err)

	consumer, err := contractDeployer.DeployVRFv2Consumer(coordinator.Address())
	require.NoError(t, err)
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, testInputs.ChainlinkNodeFunding)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	err = coordinator.SetConfig(
		uint16(minimumConfirmations),
		2.5e6,
		86400,
		33825,
		linkEthFeedResponse,
		vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: 1,
			FulfillmentFlatFeeLinkPPMTier2: 1,
			FulfillmentFlatFeeLinkPPMTier3: 1,
			FulfillmentFlatFeeLinkPPMTier4: 1,
			FulfillmentFlatFeeLinkPPMTier5: 1,
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40)},
	)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	err = coordinator.CreateSubscription()
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	err = coordinator.AddConsumer(subID, consumer.Address())
	require.NoError(t, err, "Error adding a consumer to a subscription in VRFCoordinator contract")

	actions.FundVRFCoordinatorV2Subscription(
		t,
		linkToken,
		coordinator,
		chainClient,
		subID,
		testInputs.SubscriptionFunding,
	)

	var (
		encodedProvingKeys    = make([][2]*big.Int, 0)
		nativeTokenKeyAddress string
	)
	for _, chainlinkNode := range chainlinkNodes {
		vrfKey, err := chainlinkNode.MustCreateVRFKey()
		require.NoError(t, err)
		l.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.New()
		os := &client.VRFV2TxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		require.NoError(t, err)
		nativeTokenKeyAddress, err = chainlinkNode.PrimaryEthAddress()
		require.NoError(t, err)
		_, err = chainlinkNode.MustCreateJob(&client.VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddresses:            []string{nativeTokenKeyAddress},
			EVMChainID:               fmt.Sprint(chainClient.GetNetworkConfig().ChainID),
			MinIncomingConfirmations: minimumConfirmations,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		require.NoError(t, err)
		provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
		require.NoError(t, err)
		err = coordinator.RegisterProvingKey(
			nativeTokenKeyAddress,
			provingKey,
		)
		require.NoError(t, err)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)
	}

	keyHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
	require.NoError(t, err)

	evmKeySpecificConfigTemplate := `
[[EVM.KeySpecific]]
Key = '%s'

[EVM.KeySpecific.GasEstimator]
PriceMax = '%d gwei'
`
	//todo - make evmKeySpecificConfigTemplate for multiple eth keys
	evmKeySpecificConfig := fmt.Sprintf(evmKeySpecificConfigTemplate, nativeTokenKeyAddress, maxGasPriceGWei)
	tomlConfigWithUpdates := fmt.Sprintf("%s\n%s", config.BaseVRFV2NetworkDetailTomlConfig, evmKeySpecificConfig)

	newTestEnvironment := setupVRFV2Environment(t, testNetwork, tomlConfigWithUpdates, testEnvironment.Cfg.Namespace)

	err = newTestEnvironment.RolloutStatefulSets()
	require.NoError(t, err, "Error performing rollout restart for test environment")

	err = newTestEnvironment.Run()
	require.NoError(t, err, "Error running test environment")

	//need to get node's urls again since port changed after redeployment
	chainlinkNodes, err = client.ConnectChainlinkNodes(newTestEnvironment)
	require.NoError(t, err)

	vrfV2SoakTest := testsetups.NewVRFV2SoakTest(&testsetups.VRFV2SoakTestInputs{
		BlockchainClient:     chainClient,
		TestDuration:         testInputs.TestDuration,
		ChainlinkNodeFunding: testInputs.ChainlinkNodeFunding,
		SubscriptionFunding:  testInputs.SubscriptionFunding,
		StopTestOnError:      testInputs.StopTestOnError,
		RequestsPerMinute:    testInputs.RequestsPerMinute,
		TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int, wg *sync.WaitGroup) error {

			concurrentEVMClient, err := blockchain.ConcurrentEVMClient(testNetwork, newTestEnvironment, chainClient)
			consumer.ChangeEVMClient(concurrentEVMClient)
			if err != nil {
				return fmt.Errorf("error occurred creating ConcurrentEVMClient, error: %w", err)
			}

			wg.Add(1)
			// request randomness
			err = consumer.RequestRandomness(keyHash, subID, uint16(minimumConfirmations), callbackGasLimit, numberOfWords)
			if err != nil {
				return fmt.Errorf("error occurred Requesting Randomness, error: %w", err)
			}

			err = concurrentEVMClient.WaitForEvents()
			if err != nil {
				return fmt.Errorf("error occurred waiting on chain events, error: %w", err)
			}

			lastRequestID, err := consumer.GetLastRequestId(context.Background())
			if err != nil {
				return fmt.Errorf("error occurred getting Last Request ID, error: %w", err)
			}

			l.Info().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

			_, err = WaitForRandRequestToBeFulfilled(
				consumer,
				lastRequestID,
				waitForRandRequestStatusToBeFulfilledTimeout,
				wg,
				t,
			)

			if err != nil {
				return fmt.Errorf("error occurred waiting for Randomness Request Status with request ID: %v, error: %w", lastRequestID.String(), err)
			}

			l.Info().
				Int("Request Number", requestNumber).
				Str("RequestID", lastRequestID.String()).
				Msg("Randomness fulfilled")
			return nil
		},
	},
		chainlinkNodes)

	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(vrfV2SoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	vrfV2SoakTest.Setup(t, testEnvironment)
	l.Info().Msg("Set up soak test")
	vrfV2SoakTest.Run(t)
}

func setupVRFV2Environment(t *testing.T, testNetwork blockchain.EVMNetwork, networkDetailTomlConfig string, existingNamespace string) (testEnvironment *environment.Environment) {
	gethChartConfig := getGethChartConfig(testNetwork)

	if existingNamespace != "" {
		testEnvironment = environment.New(&environment.Config{
			Namespace: existingNamespace,
			Test:      t,
			TTL:       time.Hour * 1, // 30 days,
		})
	} else {
		testEnvironment = environment.New(&environment.Config{
			NamespacePrefix: fmt.Sprintf("soak-vrfv2-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
			Test:            t,
			TTL:             time.Hour * 1, // 30 days,
		})
	}

	testEnvironment = testEnvironment.
		AddHelm(gethChartConfig).
		AddHelm(chainlink.New(0, map[string]any{
			"toml": client.AddNetworkDetailedConfig("", networkDetailTomlConfig, testNetwork),
			//need to restart the node with updated eth key config
			"db": map[string]interface{}{
				"stateful": "true",
			},
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment
}

func getGethChartConfig(testNetwork blockchain.EVMNetwork) environment.ConnectedChart {

	evmConfig := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	return evmConfig
}

func WaitForRandRequestToBeFulfilled(
	consumer contracts.VRFv2Consumer,
	lastRequestID *big.Int,
	timeout time.Duration,
	wg *sync.WaitGroup,
	t *testsetups.VRFV2SoakTest,
) (contracts.RequestStatus, error) {
	requestStatusChannel := make(chan contracts.RequestStatus)
	requestStatusErrorsChannel := make(chan error)

	// set the requests to only run for a certain amount of time
	testContext, testCancel := context.WithTimeout(context.Background(), timeout)
	defer testCancel()

	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return contracts.RequestStatus{}, fmt.Errorf("timeout waiting for new transmission event")

		case <-ticker.C:
			go getRandomnessRequestStatus(
				consumer,
				lastRequestID,
				requestStatusChannel,
				requestStatusErrorsChannel,
			)

		case requestStatus := <-requestStatusChannel:
			if requestStatus.Fulfilled == true {
				t.NumberOfFulfillments++
				wg.Done()
				return requestStatus, nil
			}

		case err := <-requestStatusErrorsChannel:
			wg.Done()
			return contracts.RequestStatus{}, err
		}
	}
}

func getRandomnessRequestStatus(
	consumer contracts.VRFv2Consumer,
	lastRequestID *big.Int,
	requestStatusChannel chan contracts.RequestStatus,
	requestStatusErrorsChannel chan error,
) {
	requestStatus, err := consumer.GetRequestStatus(context.Background(), lastRequestID)

	if err != nil {
		requestStatusErrorsChannel <- fmt.Errorf("error occurred getting Request Status for requestID: %g", lastRequestID)
	}
	requestStatusChannel <- requestStatus
}
