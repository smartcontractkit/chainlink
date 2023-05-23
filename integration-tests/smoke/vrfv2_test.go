package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	env_client "github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestVRFv2Basic(t *testing.T) {
	linkEthFeedResponse := big.NewInt(1e18)
	minimumConfirmations := 3
	subID := uint64(1)
	linkFundingAmount := big.NewInt(100)
	numberOfWords := uint32(3)
	maxGasPriceGWei := 1000
	callbackGasLimit := uint32(1000000)

	t.Parallel()
	l := utils.GetTestLogger(t)

	testNetwork := networks.SelectedNetwork
	testEnvironment := setupVRFV2Environment(t, testNetwork, config.BaseVRFV2NetworkDetailTomlConfig, "", "")
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err)
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err)
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})
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
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(1))
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
		linkFundingAmount,
	)

	var (
		job                   *client.Job
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
		job, err = chainlinkNode.MustCreateJob(&client.VRFV2JobSpec{
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

	newEnvLabel := "updatedWithRollout=true"
	newTestEnvironment := setupVRFV2Environment(t, testNetwork, tomlConfigWithUpdates, testEnvironment.Cfg.Namespace, newEnvLabel)

	err = newTestEnvironment.RolloutStatefulSets()
	require.NoError(t, err, "Error performing rollout restart for test environment")

	conds := &env_client.ReadyCheckData{
		ReadinessProbeCheckSelector: newEnvLabel,
		Timeout:                     5 * time.Minute,
	}

	err = newTestEnvironment.RunCustomReadyConditions(conds, 0)
	require.NoError(t, err, "Error running test environment")

	//need to get node's urls again since port changed after redeployment
	chainlinkNodes, err = client.ConnectChainlinkNodes(newTestEnvironment)
	require.NoError(t, err)

	err = consumer.RequestRandomness(keyHash, subID, uint16(minimumConfirmations), callbackGasLimit, numberOfWords)
	require.NoError(t, err)

	gom := gomega.NewGomegaWithT(t)
	timeout := time.Minute * 2
	var lastRequestID *big.Int
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := chainlinkNodes[0].MustReadRunsByJob(job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
		lastRequestID, err = consumer.GetLastRequestId(context.Background())
		l.Debug().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		status, err := consumer.GetRequestStatus(context.Background(), lastRequestID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(status.Fulfilled).Should(gomega.BeTrue())
		l.Debug().Interface("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		for _, w := range status.RandomWords {
			l.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
			g.Expect(w.Uint64()).Should(gomega.BeNumerically(">", 0), "Expected the VRF job give an answer bigger than 0")
		}
	}, timeout, "1s").Should(gomega.Succeed())
}

func setupVRFV2Environment(
	t *testing.T,
	testNetwork blockchain.EVMNetwork,
	networkDetailTomlConfig string,
	existingNamespace string,
	newEnvLabel string,
) (testEnvironment *environment.Environment) {
	gethChartConfig := getGethChartConfig(testNetwork)

	if existingNamespace != "" {
		testEnvironment = environment.New(&environment.Config{
			Namespace: existingNamespace,
			Test:      t,
			Labels:    []string{newEnvLabel},
		})
	} else {
		testEnvironment = environment.New(&environment.Config{
			NamespacePrefix: fmt.Sprintf("smoke-vrfv2-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
			Test:            t,
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

	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	return evmConfig
}
