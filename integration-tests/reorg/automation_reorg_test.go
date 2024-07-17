package reorg

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/client"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	geth_helm "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	reorgBlockCount       = 10 // Number of blocks to reorg (should be less than finalityDepth)
	upkeepCount           = 2
	nodeCount             = 6
	nodeFundsAmount       = new(big.Float).SetFloat64(2) // Each node will have 2 ETH
	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
	finalityDepth         int
	historyDepth          int

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

/*
 * This test verifies that conditional upkeeps automatically recover from chain reorgs.
 *
 * The test starts with happy path where upkeeps are expected to be performed.
 * Then reorg below finality depth happens which makes the chain unstable.
 *
 * Upkeeps are expected to be performed during the reorg.
 */
func TestAutomationReorg(t *testing.T) {
	c, err := tc.GetConfig([]string{"Reorg"}, tc.Automation)
	require.NoError(t, err, "Error getting config")

	findIntValue := func(text string, substring string) (int, error) {
		re := regexp.MustCompile(fmt.Sprintf(`%s\s*=\s*(\d+)`, substring))

		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			asInt, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, err
			}
			return asInt, nil
		}

		return 0, fmt.Errorf("no match found for %s", substring)
	}

	finalityDepth, err = findIntValue(c.NodeConfig.CommonChainConfigTOML, "FinalityDepth")
	require.NoError(t, err, "Error getting finality depth")

	historyDepth, err = findIntValue(c.NodeConfig.CommonChainConfigTOML, "HistoryDepth")
	require.NoError(t, err, "Error getting history depth")

	require.Less(t, reorgBlockCount, finalityDepth, "Reorg block count should be less than finality depth")

	t.Parallel()
	l := logging.GetTestLogger(t)

	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0":             ethereum.RegistryVersion_2_0,
		"registry_2_1_conditional": ethereum.RegistryVersion_2_1,
		"registry_2_1_logtrigger":  ethereum.RegistryVersion_2_1,
		"registry_2_2_conditional": ethereum.RegistryVersion_2_2, // Works only on Chainlink Node v2.10.0 or greater
		"registry_2_2_logtrigger":  ethereum.RegistryVersion_2_2, // Works only on Chainlink Node v2.10.0 or greater
		"registry_2_3_conditional": ethereum.RegistryVersion_2_3,
		"registry_2_3_logtrigger":  ethereum.RegistryVersion_2_3,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			config, err := tc.GetConfig([]string{"Reorg"}, tc.Automation)
			if err != nil {
				t.Fatal(err)
			}

			network := networks.MustGetSelectedNetworkConfig(config.Network)[0]

			tomlConfig, err := actions.BuildTOMLNodeConfigForK8s(&config, network)
			require.NoError(t, err, "Error building TOML config")

			defaultAutomationSettings["replicas"] = nodeCount
			defaultAutomationSettings["toml"] = tomlConfig

			var overrideFn = func(_ interface{}, target interface{}) {
				ctf_config.MustConfigOverrideChainlinkVersion(config.GetChainlinkImageConfig(), target)
				ctf_config.MightConfigOverridePyroscopeKey(config.GetPyroscopeConfig(), target)
			}

			cd := chainlink.NewWithOverride(0, defaultAutomationSettings, config.ChainlinkImage, overrideFn)

			testEnvironment := environment.
				New(&environment.Config{
					NamespacePrefix: fmt.Sprintf("automation-reorg-%d", reorgBlockCount),
					TTL:             time.Hour * 1,
					Test:            t}).
				// Use Geth blockchain to simulate reorgs
				AddHelm(geth_helm.New(&geth_helm.Props{
					NetworkName: network.Name,
					Simulated:   true,
					WsURLs:      network.URLs,
				})).
				AddHelm(cd)
			err = testEnvironment.Run()
			require.NoError(t, err, "Error setting up test environment")

			if testEnvironment.WillUseRemoteRunner() {
				return
			}
			if !testEnvironment.Cfg.InsideK8s {
				// Test is running locally, set forwarded URL of Geth blockchain node
				wsURLs := testEnvironment.URLs[network.Name+"_internal"]
				httpURLs := testEnvironment.URLs[network.Name+"_internal_http"]
				require.NotEmpty(t, wsURLs, "Forwarded Geth URLs should not be empty")
				require.NotEmpty(t, httpURLs, "Forwarded Geth URLs should not be empty")
				network.URLs = wsURLs
				network.HTTPURLs = httpURLs
			}

			gethRPCClient := ctf_client.NewRPCClient(network.HTTPURLs[0])
			chainClient, err := seth_utils.GetChainClient(config, network)
			require.NoError(t, err, "Error connecting to blockchain")
			chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Error connecting to Chainlink nodes")

			// Register cleanup for any test
			t.Cleanup(func() {
				err := actions.TeardownSuite(t, chainClient, testEnvironment, chainlinkNodes, nil, zapcore.PanicLevel, &config)
				require.NoError(t, err, "Error tearing down environment")
			})

			err = actions.FundChainlinkNodesFromRootAddress(l, chainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes), nodeFundsAmount)
			require.NoError(t, err, "Error funding Chainlink nodes")

			linkToken, err := contracts.DeployLinkTokenContract(l, chainClient)
			require.NoError(t, err, "Error deploying LINK token")

			wethToken, err := contracts.DeployWETHTokenContract(l, chainClient)
			require.NoError(t, err, "Error deploying weth token contract")

			// This feed is used for both eth/usd and link/usd
			ethUSDFeed, err := contracts.DeployMockETHUSDFeed(chainClient, defaultOCRRegistryConfig.FallbackLinkPrice)
			require.NoError(t, err, "Error deploying eth usd feed contract")

			defaultOCRRegistryConfig.RegistryVersion = registryVersion
			registry, registrar := actions.DeployAutoOCRRegistryAndRegistrar(
				t,
				chainClient,
				registryVersion,
				defaultOCRRegistryConfig,
				linkToken,
				wethToken,
				ethUSDFeed,
			)

			// Fund the registry with LINK
			err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(upkeepCount))))
			require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

			actions.CreateOCRKeeperJobs(t, chainlinkNodes, registry.Address(), network.ChainID, 0, registryVersion)
			nodesWithoutBootstrap := chainlinkNodes[1:]

			ocrConfig, err := actions.BuildAutoOCR2ConfigVars(t, nodesWithoutBootstrap, defaultOCRRegistryConfig, registrar.Address(), 5*time.Second, registry.ChainModuleAddress(), registry.ReorgProtectionEnabled(), linkToken, wethToken, ethUSDFeed)
			require.NoError(t, err, "OCR2 config should be built successfully")
			if registryVersion == ethereum.RegistryVersion_2_0 {
				err = registry.SetConfig(defaultOCRRegistryConfig, ocrConfig)
			} else {
				err = registry.SetConfigTypeSafe(ocrConfig)
			}
			require.NoError(t, err, "Registry config should be set successfully")

			// Use the name to determine if this is a log trigger or not
			isLogTrigger := strings.Contains(name, "logtrigger")
			consumers, upkeepIDs := actions.DeployConsumers(
				t,
				chainClient,
				registry,
				registrar,
				linkToken,
				upkeepCount,
				big.NewInt(defaultLinkFunds),
				defaultUpkeepGasLimit,
				isLogTrigger,
				false,
				false,
				wethToken,
			)

			if isLogTrigger {
				for i := 0; i < len(upkeepIDs); i++ {
					if err := consumers[i].Start(); err != nil {
						l.Error().Msg("Error when starting consumer")
						return
					}
					l.Info().Int("Consumer index", i).Msg("Consumer started")
				}
			}

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

			l.Info().
				Str("URL", gethRPCClient.URL).
				Int("BlocksBack", reorgBlockCount).
				Int("FinalityDepth", finalityDepth).
				Int("HistoryDepth", historyDepth).
				Msg("Rewinding blocks on chain below finality depth")
			err = gethRPCClient.GethSetHead(reorgBlockCount)
			require.NoError(t, err, "Error rewinding blocks on chain")

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
		})
	}
}
