package reorg

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	baseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
	networkTOML = `Enabled = true
FinalityDepth = 200
LogPollInterval = '1s'

[EVM.HeadTracker]
HistoryDepth = 400

[EVM.GasEstimator]
Mode = 'FixedPrice'
LimitDefault = 5_000_000`

	defaultAutomationSettings = map[string]interface{}{
		"toml": "",
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
		NetworkName: "",
		NetworkType: "geth-reorg",
		Values: map[string]interface{}{
			"geth": map[string]interface{}{
				"genesis": map[string]interface{}{
					"networkId": "1337",
				},
				"miner": map[string]interface{}{
					"replicas": 2,
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
		MaxRevertDataSize:    uint32(5000),
	}
)

const (
	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
	numberOfUpkeeps       = 2
	automationReorgBlocks = 50
	numberOfNodes         = 6
)

/*
 * This test verifies that conditional upkeeps automatically recover from chain reorgs
 * The blockchain is configured to have two separate miners and one geth node. The test starts
 * with happy path where the two miners remain in sync and upkeeps are expected to be performed.
 * Then reorg starts and the connection between the two geth miners is severed. This makes the
 * chain unstable, however all the CL nodes get the same view of the unstable chain through the
 * same geth node.
 *
 * Upkeeps are expected to be performed during the reorg as there are only two versions of the
 * the chain, on average 1/2 performUpkeeps should go through.
 *
 * The miner nodes are synced back after automationReorgBlocks. The syncing event can cause a
 * large reorg from CL node perspective, causing existing performUpkeeps to become staleUpkeeps.
 * Automation should be able to recover from this and upkeeps should continue to occur at a
 * normal pace after the event.
 */
func TestAutomationReorg(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0":             ethereum.RegistryVersion_2_0,
		"registry_2_1_conditional": ethereum.RegistryVersion_2_1,
		"registry_2_1_logtrigger":  ethereum.RegistryVersion_2_1,
		"registry_2_2_conditional": ethereum.RegistryVersion_2_2,
		"registry_2_2_logtrigger":  ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			config, err := tc.GetConfig("Reorg", tc.Automation)
			if err != nil {
				t.Fatal(err)
			}

			network := networks.MustGetSelectedNetworkConfig(config.Network)[0]

			defaultAutomationSettings["replicas"] = numberOfNodes
			defaultAutomationSettings["toml"] = networks.AddNetworkDetailedConfig(baseTOML, config.Pyroscope, networkTOML, network)

			var overrideFn = func(_ interface{}, target interface{}) {
				ctf_config.MustConfigOverrideChainlinkVersion(config.GetChainlinkImageConfig(), target)
				ctf_config.MightConfigOverridePyroscopeKey(config.GetPyroscopeConfig(), target)
			}

			cd := chainlink.NewWithOverride(0, defaultAutomationSettings, config.ChainlinkImage, overrideFn)

			ethSetting := defaultReorgEthereumSettings
			ethSetting.NetworkName = network.Name

			testEnvironment := environment.
				New(&environment.Config{
					NamespacePrefix: fmt.Sprintf("automation-reorg-%d", automationReorgBlocks),
					TTL:             time.Hour * 1,
					Test:            t}).
				AddHelm(reorg.New(ethSetting)).
				AddChart(blockscout.New(&blockscout.Props{
					Name:    "geth-blockscout",
					WsURL:   network.URL,
					HttpURL: network.HTTPURLs[0]})).
				AddHelm(cd)
			err = testEnvironment.Run()
			require.NoError(t, err, "Error setting up test environment")

			if testEnvironment.WillUseRemoteRunner() {
				return
			}

			chainClient, err := blockchain.NewEVMClient(network, testEnvironment, l)
			require.NoError(t, err, "Error connecting to blockchain")
			contractDeployer, err := contracts.NewContractDeployer(chainClient, l)
			require.NoError(t, err, "Error building contract deployer")
			chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Error connecting to Chainlink nodes")
			chainClient.ParallelTransactions(true)

			// Register cleanup for any test
			t.Cleanup(func() {
				err := actions.TeardownSuite(t, testEnvironment, chainlinkNodes, nil, zapcore.PanicLevel, &config, chainClient)
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
				registryVersion,
				defaultOCRRegistryConfig,
				linkToken,
				contractDeployer,
				chainClient,
			)
			// Fund the registry with LINK
			err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfUpkeeps))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			actions.CreateOCRKeeperJobs(t, chainlinkNodes, registry.Address(), network.ChainID, 0, registryVersion)
			nodesWithoutBootstrap := chainlinkNodes[1:]
			defaultOCRRegistryConfig.RegistryVersion = registryVersion
			ocrConfig, err := actions.BuildAutoOCR2ConfigVars(t, nodesWithoutBootstrap, defaultOCRRegistryConfig, registrar.Address(), 5*time.Second, registry.ChainModuleAddress(), registry.ReorgProtectionEnabled())
			require.NoError(t, err, "OCR2 config should be built successfully")
			if registryVersion == ethereum.RegistryVersion_2_0 {
				err = registry.SetConfig(defaultOCRRegistryConfig, ocrConfig)
			} else {
				err = registry.SetConfigTypeSafe(ocrConfig)
			}
			require.NoError(t, err, "Registry config should be be set successfully")
			require.NoError(t, chainClient.WaitForEvents(), "Waiting for config to be set")

			// Use the name to determine if this is a log trigger or not
			isLogTrigger := name == "registry_2_1_logtrigger"
			consumers, upkeepIDs := actions.DeployConsumers(t, registry, registrar, linkToken, contractDeployer, chainClient, numberOfUpkeeps, big.NewInt(defaultLinkFunds), defaultUpkeepGasLimit, isLogTrigger, false)

			l.Info().Msg("Waiting for all upkeeps to be performed")

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 5
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "7m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~3m for performing each upkeep 5 times, ~3m buffer

			l.Info().Msg("All upkeeps performed under happy path. Starting reorg")

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

			l.Info().Msg("Reorg started. Expecting chain to become unstable and upkeeps to still getting performed")

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they reach 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 10
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed())

			l.Info().Msg("Upkeep performed during unstable chain, waiting for reorg to finish")
			err = rc.WaitDepthReached()
			require.NoError(t, err)

			l.Info().Msg("Reorg finished, chain should be stable now. Expecting upkeeps to keep getting performed")
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they reach 20
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 20
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "10m", "1s").Should(gomega.Succeed())
		})
	}
}
