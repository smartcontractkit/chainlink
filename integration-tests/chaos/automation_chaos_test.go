package chaos

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
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
	defaultAutomationSettings = map[string]interface{}{
		"toml":     client.AddNetworksConfig(baseTOML, networks.SelectedNetwork),
		"replicas": "6",
		"db": map[string]interface{}{
			"stateful": true,
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

	defaultEthereumSettings = &ethereum.Props{
		NetworkName: networks.SelectedNetwork.Name,
		Simulated:   networks.SelectedNetwork.Simulated,
		WsURLs:      networks.SelectedNetwork.URLs,
		Values: map[string]interface{}{
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "4000m",
					"memory": "4Gi",
				},
				"limits": map[string]interface{}{
					"cpu":    "4000m",
					"memory": "4Gi",
				},
			},
			"geth": map[string]interface{}{
				"blocktime": "1",
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

type KeeperConsumerContracts int32

const (
	BasicCounter KeeperConsumerContracts = iota

	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
	numberOfUpkeeps       = 2
)

func TestAutomationChaos(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	testCases := map[string]struct {
		networkChart environment.ConnectedChart
		clChart      environment.ConnectedChart
		chaosFunc    chaos.ManifestFunc
		chaosProps   *chaos.Props
	}{
		// see ocr_chaos.test.go for comments
		PodChaosFailMinorityNodes: {
			ethereum.New(defaultEthereumSettings),
			chainlink.New(0, defaultAutomationSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityNodes: {
			ethereum.New(defaultEthereumSettings),
			chainlink.New(0, defaultAutomationSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityDB: {
			ethereum.New(defaultEthereumSettings),
			chainlink.New(0, defaultAutomationSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
				DurationStr:    "1m",
				ContainerNames: &[]*string{a.Str("chainlink-db")},
			},
		},
		NetworkChaosFailMajorityNetwork: {
			ethereum.New(defaultEthereumSettings),
			chainlink.New(0, defaultAutomationSettings),
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{ChaosGroupMajority: a.Str("1")},
				ToLabels:    &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr: "1m",
			},
		},
		NetworkChaosFailBlockchainNode: {
			ethereum.New(defaultEthereumSettings),
			chainlink.New(0, defaultAutomationSettings),
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": a.Str("geth")},
				ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: a.Str("1")},
				DurationStr: "1m",
			},
		},
	}

	for n, tst := range testCases {
		name := n
		testCase := tst
		t.Run(fmt.Sprintf("Automation_%s", name), func(t *testing.T) {
			t.Parallel()
			network := networks.SelectedNetwork

			testEnvironment := environment.
				New(&environment.Config{
					NamespacePrefix: fmt.Sprintf("chaos-automation-%s", name),
					TTL:             time.Hour * 1,
					Test:            t,
				}).
				AddHelm(testCase.networkChart).
				AddHelm(testCase.clChart).
				AddChart(blockscout.New(&blockscout.Props{
					Name:    "geth-blockscout",
					WsURL:   network.URL,
					HttpURL: network.HTTPURLs[0],
				}))
			err := testEnvironment.Run()
			require.NoError(t, err, "Error setting up test environment")
			if testEnvironment.WillUseRemoteRunner() {
				return
			}

			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 2, ChaosGroupMinority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 3, 5, ChaosGroupMajority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 2, 5, ChaosGroupMajorityPlus)
			require.NoError(t, err)

			chainClient, err := blockchain.NewEVMClient(network, testEnvironment)
			require.NoError(t, err, "Error connecting to blockchain")
			contractDeployer, err := contracts.NewContractDeployer(chainClient)
			require.NoError(t, err, "Error building contract deployer")
			chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Error connecting to Chainlink nodes")
			chainClient.ParallelTransactions(true)

			// Register cleanup for any test
			t.Cleanup(func() {
				if chainClient != nil {
					chainClient.GasStats().PrintStats()
				}
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
				eth_contracts.RegistryVersion_2_0,
				defaultOCRRegistryConfig,
				numberOfUpkeeps,
				linkToken,
				contractDeployer,
				chainClient,
			)

			actions.CreateOCRKeeperJobs(t, chainlinkNodes, registry.Address(), network.ChainID, 0)
			nodesWithoutBootstrap := chainlinkNodes[1:]
			ocrConfig := actions.BuildAutoOCR2ConfigVars(t, nodesWithoutBootstrap, defaultOCRRegistryConfig, registrar.Address(), 5*time.Second)
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

			l.Info().Msg("Waiting for all upkeeps to be performed")

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 5
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			_, err = testEnvironment.Chaos.Run(testCase.chaosFunc(testEnvironment.Cfg.Namespace, testCase.chaosProps))
			require.NoError(t, err)

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 10
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "3m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer
		})
	}
}
