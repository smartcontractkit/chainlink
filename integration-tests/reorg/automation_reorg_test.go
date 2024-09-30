package reorg

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	sethUtils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/automationv2"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	reorgBlockCount       = 10 // Number of blocks to reorg (should be less than finalityDepth)
	upkeepCount           = 2
	nodeCount             = 6
	defaultUpkeepGasLimit = uint32(2500000)
	defaultLinkFunds      = int64(9e18)
	finalityDepth         int
	historyDepth          int
)

var logScannerSettings = test_env.GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(testreporters.NewAllowedLogMessage(
	"Got very old block.",
	"It is expected, because we are causing reorgs",
	zapcore.DPanicLevel,
	testreporters.WarnAboutAllowedMsgs_No,
))

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

	finalityDepth, err = findIntValue(c.NodeConfig.ChainConfigTOMLByChainID["1337"], "FinalityDepth")
	require.NoError(t, err, "Error getting finality depth")

	historyDepth, err = findIntValue(c.NodeConfig.ChainConfigTOMLByChainID["1337"], "HistoryDepth")
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

			privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
			require.NoError(t, err, "Error building ethereum network config")

			env, err := test_env.NewCLTestEnvBuilder().
				WithTestInstance(t).
				WithTestConfig(&config).
				WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
				WithMockAdapter().
				WithCLNodes(nodeCount).
				WithStandardCleanup().
				WithChainlinkNodeLogScanner(logScannerSettings).
				Build()
			require.NoError(t, err, "Error deploying test environment")

			nodeClients := env.ClCluster.NodeAPIs()

			evmNetwork, err := env.GetFirstEvmNetwork()
			require.NoError(t, err, "Error getting first evm network")

			sethClient, err := sethUtils.GetChainClient(&config, *evmNetwork)
			require.NoError(t, err, "Error getting seth client")

			err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(*config.GetCommonConfig().ChainlinkNodeFunding))
			require.NoError(t, err, "Failed to fund the nodes")

			gethRPCClient := ctfClient.NewRPCClient(evmNetwork.HTTPURLs[0], nil)

			a := automationv2.NewAutomationTestDocker(l, sethClient, nodeClients, &config)
			a.SetMercuryCredentialName("cred1")
			a.RegistrySettings = actions.ReadRegistryConfig(config)
			a.RegistrySettings.RegistryVersion = registryVersion
			a.PluginConfig = actions.ReadPluginConfig(config)
			a.PublicConfig = actions.ReadPublicConfig(config)
			a.RegistrarSettings = contracts.KeeperRegistrarSettings{
				AutoApproveConfigType: uint8(2),
				AutoApproveMaxAllowed: 1000,
				MinLinkJuels:          big.NewInt(0),
			}

			a.SetupAutomationDeployment(t)
			a.SetDockerEnv(env)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			// Use the name to determine if this is a log trigger or not
			isLogTrigger := strings.Contains(name, "logtrigger")
			consumers, upkeepIDs := actions.DeployConsumers(
				t,
				sethClient,
				a.Registry,
				a.Registrar,
				a.LinkToken,
				upkeepCount,
				big.NewInt(defaultLinkFunds),
				defaultUpkeepGasLimit,
				isLogTrigger,
				false,
				false,
				a.WETHToken,
				&config,
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
