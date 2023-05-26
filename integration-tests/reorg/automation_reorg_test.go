package reorg

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

var (
	baseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
	simulatedEVMNonDevTOML = `
[[EVM]]
ChainID = 1337
MinContractPayment = '0'
Enabled = true
FinalityDepth = 50
LogPollInterval = '1s'

[EVM.HeadTracker]
HistoryDepth = 100

[EVM.GasEstimator]
Mode = 'FixedPrice'
LimitDefault = 5_000_000`
	networkTOML = `FinalityDepth = 200

[EVM.HeadTracker]
HistoryDepth = 400`
	activeEVMNetwork          = networks.SelectedNetwork
	defaultAutomationSettings = map[string]interface{}{
		"toml":     client.AddNetworkDetailedConfig(baseTOML+simulatedEVMNonDevTOML, networkTOML, activeEVMNetwork),
		"replicas": "6",
		"db": map[string]interface{}{
			"stateful": false,
			"capacity": "1Gi",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "250m",
					"memory": "256Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "250m",
					"memory": "256Mi",
				},
			},
		},
	}

	defaultReorgEthereumSettings = &reorg.Props{
		NetworkName: activeEVMNetwork.Name,
		NetworkType: "geth-reorg",
		Values: map[string]interface{}{
			"geth": map[string]interface{}{
				"genesis": map[string]interface{}{
					"networkId": "1337",
				},
				"miner": map[string]interface{}{
					"replicas": "2",
				},
			},
		},
	}

	defaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		FlatFeeMicroLINK:     uint32(0),
		BlockCountPerTurn:    big.NewInt(10),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		MinUpkeepSpend:       big.NewInt(0),
		MaxPerformGas:        uint32(5000000),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5000),
		MaxPerformDataSize:   uint32(5000),
	}
)

const (
	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
	numberOfUpkeeps       = 2
	automationReorgBlocks = 10 //TODO: Make this a flag
)

func TestAutomationReorg(t *testing.T) {
	network := networks.SelectedNetwork

	testEnvironment := environment.
		New(&environment.Config{
			NamespacePrefix: fmt.Sprintf("automation-reorg-%d", automationReorgBlocks),
			TTL:             time.Hour * 1}).
		AddHelm(reorg.New(defaultReorgEthereumSettings)).
		AddHelm(chainlink.New(0, defaultAutomationSettings)).
		AddChart(blockscout.New(&blockscout.Props{
			Name:    "geth-blockscout",
			WsURL:   activeEVMNetwork.URL,
			HttpURL: activeEVMNetwork.HTTPURLs[0]}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error setting up test environment")

	chainClient, err := blockchain.NewEVMClient(network, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Error building contract deployer")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	chainClient.ParallelTransactions(true)

	// Register cleanup for any test
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.PanicLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	txCost, err := chainClient.EstimateCostForChainlinkOperations(1000)
	require.NoError(t, err, "Error estimating cost for Chainlink Operations")
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
	require.NoError(t, err, "Error funding Chainlink nodes")

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying LINK token")

	registry, registrar := actions.DeployAutoOCRRegistryAndRegistrar(
		t,
		ethereum.RegistryVersion_2_0,
		defaultOCRRegistryConfig,
		numberOfUpkeeps,
		linkToken,
		contractDeployer,
		chainClient,
	)

	actions.CreateOCRKeeperJobs(t, chainlinkNodes, registry.Address(), network.ChainID, 0)
	nodesWithoutBootstrap := chainlinkNodes[1:]
	ocrConfig, err := actions.BuildAutoOCR2ConfigVars(t, nodesWithoutBootstrap, defaultOCRRegistryConfig, registrar.Address(), 5*time.Second)
	require.NoError(t, err, "OCR2 config should be built successfully")
	err = registry.SetConfig(defaultOCRRegistryConfig, ocrConfig)
	require.NoError(t, err, "Registry config should be be set successfully")
	require.NoError(t, chainClient.WaitForEvents(), "Waiting for config to be set")

	consumers, upkeepIDs := actions.DeployConsumers(
		t,
		registry,
		registrar,
		linkToken,
		contractDeployer,
		chainClient,
		numberOfUpkeeps,
		big.NewInt(defaultLinkFunds),
		defaultUpkeepGasLimit,
	)

	log.Info().Msg("Waiting for all upkeeps to be performed")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			expect := 5
			log.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
				"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
		}
	}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

	rc, err := NewReorgController(
		&ReorgConfig{
			FromPodLabel:            reorg.TXNodesAppLabel,
			ToPodLabel:              reorg.MinerNodesAppLabel,
			Network:                 chainClient,
			Env:                     testEnvironment,
			BlockConsensusThreshold: 3,
			Timeout:                 1800 * time.Second,
		},
	)

	require.NoError(t, err, "Error getting reorg controller")
	rc.ReOrg(automationReorgBlocks)
	rc.WaitReorgStarted()

	gom.Eventually(func(g gomega.Gomega) {
		// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			expect := 10
			log.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
				"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
		}
	}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

	err = rc.WaitDepthReached()

	gom.Eventually(func(g gomega.Gomega) {
		// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
			require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			expect := 20
			log.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
				"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
		}
	}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

}
