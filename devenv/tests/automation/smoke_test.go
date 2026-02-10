package automation

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
)

const (
	defaultUpkeepGasLimit           = uint32(2500000)
	defaultLinkFunds                = 9
	defaultEthFunds                 = 10.0
	defaultAmountOfUpkeeps          = 2
	defaultUpkeepExecutionTimeout   = "10m" // ~1m for cluster setup, ~7m for performing each upkeep 5 times, ~2m buffer
	defaultExpectedUpkeepExecutions = 5
)

func TestRegistry_2_0(t *testing.T) {
	basicAutomationTest(t, Testcase{
		RegistryVersion:          contracts.RegistryVersion_2_0,
		Name:                     "registry_2_0",
		UpkeepCount:              defaultAmountOfUpkeeps,
		UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
		ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
		TestKeyFundingEth:        defaultEthFunds,
		UpkeepFundingLink:        defaultLinkFunds,
		upgradeImage:             os.Getenv("CHAINLINK_UPGRADE_IMAGE"),
	})
}

func TestRegistry_2_1(t *testing.T) {
	// testNames := []string{"registry_2_1_conditional", "registry_2_1_logtrigger", "registry_2_1_with_mercury_v02", "registry_2_1_with_mercury_v03"}
	testNames := []string{"registry_2_1_logtrigger"}
	for _, tc := range testNames {
		basicAutomationTest(t, Testcase{
			RegistryVersion:          contracts.RegistryVersion_2_1,
			Name:                     tc,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			TestKeyFundingEth:        defaultEthFunds,
			UpkeepFundingLink:        defaultLinkFunds,
			upgradeImage:             os.Getenv("CHAINLINK_UPGRADE_IMAGE"),
		})
	}
}

func TestRegistry_2_2(t *testing.T) {
	testNames := []string{"registry_2_2_conditional", "registry_2_2_logtrigger", "registry_2_2_with_mercury_v02", "registry_2_2_with_mercury_v03", "registry_2_1_with_logtrigger_and_mercury_v02"}
	for _, tc := range testNames {
		basicAutomationTest(t, Testcase{
			RegistryVersion:          contracts.RegistryVersion_2_2,
			Name:                     tc,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			TestKeyFundingEth:        defaultEthFunds,
			UpkeepFundingLink:        defaultLinkFunds,
			upgradeImage:             os.Getenv("CHAINLINK_UPGRADE_IMAGE"),
		})
	}
}

func TestRegistry_2_3(t *testing.T) {
	testNames := []string{"registry_2_3_conditional_native", "registry_2_3_conditional_link", "registry_2_3_logtrigger_native", "registry_2_3_logtrigger_link", "registry_2_3_with_logtrigger_and_mercury_v02_link"}
	for _, tc := range testNames {
		basicAutomationTest(t, Testcase{
			RegistryVersion:          contracts.RegistryVersion_2_3,
			Name:                     tc,
			UpkeepCount:              defaultAmountOfUpkeeps,
			UpkeepExecutionTimeout:   defaultUpkeepExecutionTimeout,
			ExpectedUpkeepExecutions: defaultExpectedUpkeepExecutions,
			UpkeepFundingLink:        defaultLinkFunds,
			TestKeyFundingEth:        defaultEthFunds,
			upgradeImage:             os.Getenv("CHAINLINK_UPGRADE_IMAGE"),
		})
	}
}

func basicAutomationTest(t *testing.T, testcase Testcase) {
	l := framework.L
	l.Info().Msg("Running test " + testcase.Name + " with registry version " + testcase.RegistryVersion.String())

	t.Cleanup(func() {
		err := products.ScanLogs(l, products.DefaultSettings())
		require.NoError(t, err, "Found concerning logs in Chainlink Node logs")

		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
	require.NoError(t, err)
	// Use the name to determine if this is a log trigger or mercury or billing token is native
	isBillingTokenNative := strings.Contains(testcase.Name, "native")
	isLogTrigger := strings.Contains(testcase.Name, "logtrigger")
	isMercuryV02 := strings.Contains(testcase.Name, "mercury_v02")
	isMercuryV03 := strings.Contains(testcase.Name, "mercury_v03")
	isMercury := isMercuryV02 || isMercuryV03

	var config *automation.Automation
	for _, candidate := range pdConfig.Config {
		if candidate.MustGetRegistryVersion() == testcase.RegistryVersion {
			if !isMercury {
				config = candidate
				break
			}

			if isMercuryV02 && candidate.MercurySettings != nil && candidate.MercurySettings.Version == "v2" {
				config = candidate
				break
			}

			if isMercuryV03 && candidate.MercurySettings != nil && candidate.MercurySettings.Version == "v3" {
				config = candidate
				break
			}
		}
	}
	require.NotNil(t, config, "failed to find matching config with registry version %v; mercury v2: %v, mercury v3: %v", testcase.RegistryVersion.String(), isMercuryV02, isMercuryV03)

	pks := []string{products.NetworkPrivateKey()}

	require.Equal(t, "1337", in.Blockchains[0].ChainID, "automation smoke tests can only be run on simulated network. If do want to run on a live network, please read the code, understand the implications (e.g. potential fund loss) and adjust the test accordingly")

	// on simulated network create new ephemeral addresses if insufficient private keys were provided
	// we ignore key at index 0, because it is the root key, which is not used during the test
	// for contract deployment and interaction
	// we create new addresses only on the simulated network to protect against fund loss
	if in.Blockchains[0].ChainID == "1337" && len(pks)-1 != testcase.UpkeepCount {
		bcNode := in.Blockchains[0].Out.Nodes[0]
		c, _, _, err := products.ETHClient(
			t.Context(),
			bcNode.ExternalWSUrl,
			config.GasSettings.FeeCapMultiplier,
			config.GasSettings.TipCapMultiplier,
		)
		require.NoError(t, err, "Failed to create ETH client")

		newPks, err := products.FundNewAddresses(t.Context(), testcase.UpkeepCount, c, testcase.TestKeyFundingEth)
		require.NoError(t, err, "Failed to fund new addresses")
		pks = append(pks, newPks...)
	}
	require.GreaterOrEqual(t, len(pks), testcase.UpkeepCount+1, "you must provide at least %d private keys", testcase.UpkeepCount+1)
	chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
	require.NoError(t, err, "Failed to parse chain ID")

	chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
	require.NoError(t, err, "Failed to create chain client")

	sb, err := chainClient.Client.BlockNumber(t.Context())
	require.NoError(t, err, "Failed to get start block")

	a, err := NewTest(chainClient, config)
	require.NoError(t, err, "Failed to create automation test")

	consumers, upkeepIDs := automation.DeployConsumers(
		t,
		a.ChainClient,
		a.Registry,
		a.Registrar,
		a.LinkToken,
		testcase.UpkeepCount,
		big.NewInt(0).Mul(big.NewInt(testcase.UpkeepFundingLink), big.NewInt(1e18)),
		defaultUpkeepGasLimit,
		isLogTrigger,
		isMercury,
		isBillingTokenNative,
		a.WETHToken,
		*config,
	)

	// copied from core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/streams/streams.go to avoid depending on chainlink/v2
	type UpkeepPrivilegeConfig struct {
		MercuryEnabled bool `json:"mercuryEnabled"`
	}

	// Do it in two separate loops, so we don't end up setting up one upkeep, but starting the consumer for another one
	// since we cannot be sure that consumers and upkeeps at the same index are related
	if isMercury {
		for i := range upkeepIDs {
			// Set privilege config to enable mercury
			privilegeConfigBytes, _ := json.Marshal(UpkeepPrivilegeConfig{
				MercuryEnabled: true,
			})
			err := a.Registry.SetUpkeepPrivilegeConfig(upkeepIDs[i], privilegeConfigBytes)
			require.NoError(t, err, "error when setting upkeep privilege config")

			l.Info().Int("Upkeep index", i).Msg("Upkeep privilege config set")
		}
	}

	if isLogTrigger || isMercuryV02 {
		for i := range upkeepIDs {
			err := consumers[i].Start()
			require.NoError(t, err, "error when starting consumer")
			l.Info().Int("Consumer index", i).Msg("Consumer started")
		}
	}

	l.Info().Msgf("Waiting %s for %d upkeeps to be performed by %d contracts", testcase.UpkeepExecutionTimeout, testcase.ExpectedUpkeepExecutions, testcase.UpkeepCount)
	gom := gomega.NewGomegaWithT(t)
	startTime := time.Now()

	t.Cleanup(func() {
		GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, testcase.RegistryVersion)()
	})

	gom.Eventually(func(g gomega.Gomega) {
		// Check if the upkeeps are performing multiple times by analyzing their counters
		for i := range upkeepIDs {
			counter, err := consumers[i].Counter(t.Context())
			require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			expect := testcase.ExpectedUpkeepExecutions
			l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
				"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
		}
	}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())

	l.Info().Msgf("Total time taken to get %d performs for each upkeep: %s", testcase.ExpectedUpkeepExecutions, time.Since(startTime))

	if testcase.upgradeImage != "" {
		expect := testcase.ExpectedUpkeepExecutions
		// Upgrade the nodes one at a time and check that the upkeeps are still being performed
		for i := range 5 {
			in.NodeSets[0].NodeSpecs[i].Node.Image = testcase.upgradeImage
			l.Info().Msgf("Upgrading node %d to version %s", i, testcase.upgradeImage)
			err = products.RestartNodes(t.Context(), in.NodeSets[0], in.Blockchains[0], true, time.Minute)
			require.NoError(t, err, "Error when upgrading node %d", i)
			time.Sleep(time.Second * 10)
			expect += testcase.ExpectedUpkeepExecutions
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are increasing by 5 in each step within 5 minutes
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep index", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, testcase.UpkeepExecutionTimeout, "1s").Should(gomega.Succeed())
			l.Info().Msgf("All upkeeps performed after upgrading node %d", i)
		}
	}

	// Cancel all the registered upkeeps via the registry
	for i := range upkeepIDs {
		err := a.Registry.CancelUpkeep(upkeepIDs[i])
		require.NoError(t, err, "Could not cancel upkeep at index %d", i)
	}

	var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

	for i := range upkeepIDs {
		// Obtain the amount of times the upkeep has been executed so far
		countersAfterCancellation[i], err = consumers[i].Counter(t.Context())
		require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
		l.Info().Int64("Upkeep Count", countersAfterCancellation[i].Int64()).Int("Upkeep Index", i).Msg("Cancelled upkeep")
	}

	l.Info().Msg("Making sure the counter stays consistent")
	gom.Consistently(func(g gomega.Gomega) {
		for i := range upkeepIDs {
			// Expect the counter to remain constant (At most increase by 1 to account for stale performs) because the upkeep was cancelled
			latestCounter, err := consumers[i].Counter(t.Context())
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
			g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterCancellation[i].Int64()+1),
				"Expected consumer counter to remain less than or equal to %d, but got %d",
				countersAfterCancellation[i].Int64()+1, latestCounter.Int64())
		}
	}, "1m", "1s").Should(gomega.Succeed())
}
