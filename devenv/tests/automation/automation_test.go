package automation

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
)

const (
	automationDefaultUpkeepGasLimit = uint32(2500000)
	automationDefaultLinkFunds      = int64(9e18)
	automationExpectedData          = "abcdef"
	defaultAmountOfUpkeeps          = 2
)

func TestAutomationBasic(t *testing.T) {
	SetupAutomationBasic(t, false)
}

func SetupAutomationBasic(t *testing.T, nodeUpgrade bool) {
	// t.Parallel()

	// native, mercury_v02, mercury_v03 and logtrigger are reserved keywords, use them with caution
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		// "registry_2_0": ethereum.RegistryVersion_2_0,
		// "registry_2_1_conditional":      ethereum.RegistryVersion_2_1,
		// "registry_2_1_logtrigger":       ethereum.RegistryVersion_2_1,
		// "registry_2_1_with_mercury_v02": ethereum.RegistryVersion_2_1,
		// "registry_2_1_with_mercury_v03":                     ethereum.RegistryVersion_2_1,
		// "registry_2_1_with_logtrigger_and_mercury_v02": ethereum.RegistryVersion_2_1,
		// "registry_2_2_conditional":      ethereum.RegistryVersion_2_2,
		// "registry_2_2_logtrigger":       ethereum.RegistryVersion_2_2,
		// "registry_2_2_with_mercury_v02": ethereum.RegistryVersion_2_2,
		"registry_2_2_with_mercury_v03": ethereum.RegistryVersion_2_2,
		// "registry_2_2_with_logtrigger_and_mercury_v02": ethereum.RegistryVersion_2_2,
		// "registry_2_3_conditional_native":                   ethereum.RegistryVersion_2_3,
		// "registry_2_3_conditional_link":                     ethereum.RegistryVersion_2_3,
		// "registry_2_3_logtrigger_native": ethereum.RegistryVersion_2_3,
		// "registry_2_3_logtrigger_link":                      ethereum.RegistryVersion_2_3,
		// "registry_2_3_with_mercury_v03_link":                ethereum.RegistryVersion_2_3,
		// "registry_2_3_with_logtrigger_and_mercury_v02_link": ethereum.RegistryVersion_2_3,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			// t.Parallel()
			l := framework.L

			// TODO run CL node log scanner

			outputFile := "../../env-out.toml"
			in, err := de.LoadOutput[de.Cfg](outputFile)
			require.NoError(t, err)
			pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
			require.NoError(t, err)

			// if nodeUpgrade {
			// 	if cfg.GetChainlinkUpgradeImageConfig() == nil {
			// 		t.Fatal("[ChainlinkUpgradeImage] must be set in TOML config to upgrade nodes")
			// 	}
			// }

			// Use the name to determine if this is a log trigger or mercury or billing token is native
			isBillingTokenNative := strings.Contains(name, "native")
			isLogTrigger := strings.Contains(name, "logtrigger")
			isMercuryV02 := strings.Contains(name, "mercury_v02")
			isMercuryV03 := strings.Contains(name, "mercury_v03")
			isMercury := isMercuryV02 || isMercuryV03

			var config *automation.Automation
			for _, candidate := range pdConfig.Config {
				if candidate.MustGetRegistryVersion() == registryVersion {
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
			require.NotNil(t, config, "failed to find matching config with registry version %v; mercury v2: %b, mercury v3: %b", registryVersion, isMercuryV02, isMercuryV03)

			pks := []string{products.NetworkPrivateKey()}

			// on simulated network create new ephemeral addresses if insufficient private keys were provided
			// we require +1 private keys, because key at index 0 is the root key, which is not used during the test
			// for contract deployment and interaction
			// we create new addresses only on the simulated network to protect against fund loss
			if in.Blockchains[0].ChainID == "1337" && len(pks) != defaultAmountOfUpkeeps+1 {
				bcNode := in.Blockchains[0].Out.Nodes[0]
				c, _, _, err := products.ETHClient(
					t.Context(),
					bcNode.ExternalWSUrl,
					config.GasSettings.FeeCapMultiplier,
					config.GasSettings.TipCapMultiplier,
				)

				newPks, err := products.FundNewAddresses(t.Context(), defaultAmountOfUpkeeps, c, config.TestKeysMinFundingEth)
				require.NoError(t, err, "Failed to fund new addresses")
				pks = append(pks, newPks...)
			}
			require.GreaterOrEqual(t, len(pks), defaultAmountOfUpkeeps+1, "you must provide at least %d private keys", defaultAmountOfUpkeeps+1)

			chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
			require.NoError(t, err, "Failed to parse chain ID")

			chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
			require.NoError(t, err, "Failed to create chain client")

			sb, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get start block")

			a := AutomationTest{
				ChainClient:            chainClient,
				Config:                 *config,
				RegistrySettings:       automation.ReadRegistryConfig(config),
				PublicConfig:           automation.ReadPublicConfig(config.PublicConfig),
				PluginConfig:           automation.ReadPluginConfig(config.PluginConfig),
				UpkeepPrivilegeManager: chainClient.MustGetRootKeyAddress(),
				Logger:                 framework.L,
			}

			err = a.LoadContracts()
			require.NoError(t, err, "Failed to load contracts")

			consumers, upkeepIDs := automation.DeployConsumers(
				t,
				a.ChainClient,
				a.Registry,
				a.Registrar,
				a.LinkToken,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
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
					if err := a.Registry.SetUpkeepPrivilegeConfig(upkeepIDs[i], privilegeConfigBytes); err != nil {
						l.Error().Msg("Error when setting upkeep privilege config")
						return
					}
					l.Info().Int("Upkeep index", i).Msg("Upkeep privilege config set")
				}
			}

			if isLogTrigger || isMercuryV02 {
				for i := range upkeepIDs {
					if err := consumers[i].Start(); err != nil {
						l.Error().Msg("Error when starting consumer")
						return
					}
					l.Info().Int("Consumer index", i).Msg("Consumer started")
				}
			}

			l.Info().Msg("Waiting 5m for all upkeeps to be performed")
			gom := gomega.NewGomegaWithT(t)
			startTime := time.Now()

			t.Cleanup(func() {
				automation.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters
				for i := range upkeepIDs {
					counter, err := consumers[i].Counter(t.Context())
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 5
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			l.Info().Msgf("Total time taken to get 5 performs for each upkeep: %s", time.Since(startTime))

			// if nodeUpgrade {
			// 	// TODO: update ref
			// 	require.NotNil(t, cfg.GetChainlinkImageConfig(), "unable to upgrade node version, [ChainlinkUpgradeImage] was not set, must both a new image or a new version")
			// 	expect := 5
			// 	// Upgrade the nodes one at a time and check that the upkeeps are still being performed
			// 	for i := range 5 {
			// 		err = upgradeChainlinkNodeVersionsLocal(*cfg.GetChainlinkUpgradeImageConfig().Image, *cfg.GetChainlinkUpgradeImageConfig().Version, a.DockerEnv.ClCluster.Nodes[i])
			// 		require.NoError(t, err, "Error when upgrading node %d", i)
			// 		time.Sleep(time.Second * 10)
			// 		expect += 5
			// 		gom.Eventually(func(g gomega.Gomega) {
			// 			// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are increasing by 5 in each step within 5 minutes
			// 			for i := range upkeepIDs {
			// 				counter, err := consumers[i].Counter(t.Context())
			// 				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			// 				l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep index", i).Msg("Number of upkeeps performed")
			// 				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
			// 					"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
			// 			}
			// 		}, "5m", "1s").Should(gomega.Succeed())
			// 	}
			// }

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
		})
	}
}
