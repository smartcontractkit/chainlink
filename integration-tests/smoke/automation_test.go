package smoke

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fxamacker/cbor/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/automationv2"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/gasprice"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/streams"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
)

const (
	automationDefaultUpkeepGasLimit = uint32(2500000)
	automationDefaultLinkFunds      = int64(9e18)
	automationExpectedData          = "abcdef"
	defaultAmountOfUpkeeps          = 2
)

// UpgradeChainlinkNodeVersions upgrades all Chainlink nodes to a new version, and then runs the test environment
// to apply the upgrades
func upgradeChainlinkNodeVersionsLocal(
	newImage, newVersion string,
	nodes ...*test_env.ClNode,
) error {
	for _, node := range nodes {
		if err := node.UpgradeVersion(newImage, newVersion); err != nil {
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	logging.Init()
	// config, err := tc.GetConfig(tc.NoTest, tc.Smoke, tc.Automation)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Running Smoke Test on %s\n", networks.MustGetSelectedNetworkConfig(config.Network)[0].Name) // Print to get around disabled logging
	// fmt.Printf("Chainlink Image %v\n", config.ChainlinkImage.Image)                                         // Print to get around disabled logging
	// fmt.Printf("Chainlink Version %v\n", config.ChainlinkImage.Version)                                     // Print to get around disabled logging
	os.Exit(m.Run())
}

func TestAutomationBasic(t *testing.T) {
	SetupAutomationBasic(t, false)
}

func SetupAutomationBasic(t *testing.T, nodeUpgrade bool) {
	t.Parallel()

	// native, mercury_v02, mercury_v03 and logtrigger are reserved keywords, use them with caution
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0":                                      ethereum.RegistryVersion_2_0,
		"registry_2_1_conditional":                          ethereum.RegistryVersion_2_1,
		"registry_2_1_logtrigger":                           ethereum.RegistryVersion_2_1,
		"registry_2_1_with_mercury_v02":                     ethereum.RegistryVersion_2_1,
		"registry_2_1_with_mercury_v03":                     ethereum.RegistryVersion_2_1,
		"registry_2_1_with_logtrigger_and_mercury_v02":      ethereum.RegistryVersion_2_1,
		"registry_2_2_conditional":                          ethereum.RegistryVersion_2_2,
		"registry_2_2_logtrigger":                           ethereum.RegistryVersion_2_2,
		"registry_2_2_with_mercury_v02":                     ethereum.RegistryVersion_2_2,
		"registry_2_2_with_mercury_v03":                     ethereum.RegistryVersion_2_2,
		"registry_2_2_with_logtrigger_and_mercury_v02":      ethereum.RegistryVersion_2_2,
		"registry_2_3_conditional_native":                   ethereum.RegistryVersion_2_3,
		"registry_2_3_conditional_link":                     ethereum.RegistryVersion_2_3,
		"registry_2_3_logtrigger_native":                    ethereum.RegistryVersion_2_3,
		"registry_2_3_logtrigger_link":                      ethereum.RegistryVersion_2_3,
		"registry_2_3_with_mercury_v03_link":                ethereum.RegistryVersion_2_3,
		"registry_2_3_with_logtrigger_and_mercury_v02_link": ethereum.RegistryVersion_2_3,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)

			cfg, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			if nodeUpgrade {
				if cfg.GetChainlinkUpgradeImageConfig() == nil {
					t.Fatal("[ChainlinkUpgradeImage] must be set in TOML config to upgrade nodes")
				}
			}

			// Use the name to determine if this is a log trigger or mercury or billing token is native
			isBillingTokenNative := strings.Contains(name, "native")
			isLogTrigger := strings.Contains(name, "logtrigger")
			isMercuryV02 := strings.Contains(name, "mercury_v02")
			isMercuryV03 := strings.Contains(name, "mercury_v03")
			isMercury := isMercuryV02 || isMercuryV03

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(cfg), isMercuryV02, isMercuryV03, &cfg,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(
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
				&cfg,
			)

			// Do it in two separate loops, so we don't end up setting up one upkeep, but starting the consumer for another one
			// since we cannot be sure that consumers and upkeeps at the same index are related
			if isMercury {
				for i := 0; i < len(upkeepIDs); i++ {
					// Set privilege config to enable mercury
					privilegeConfigBytes, _ := json.Marshal(streams.UpkeepPrivilegeConfig{
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
				for i := 0; i < len(upkeepIDs); i++ {
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
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 5
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			l.Info().Msgf("Total time taken to get 5 performs for each upkeep: %s", time.Since(startTime))

			if nodeUpgrade {
				require.NotNil(t, cfg.GetChainlinkImageConfig(), "unable to upgrade node version, [ChainlinkUpgradeImage] was not set, must both a new image or a new version")
				expect := 5
				// Upgrade the nodes one at a time and check that the upkeeps are still being performed
				for i := 0; i < 5; i++ {
					err = upgradeChainlinkNodeVersionsLocal(*cfg.GetChainlinkUpgradeImageConfig().Image, *cfg.GetChainlinkUpgradeImageConfig().Version, a.DockerEnv.ClCluster.Nodes[i])
					require.NoError(t, err, "Error when upgrading node %d", i)
					time.Sleep(time.Second * 10)
					expect += 5
					gom.Eventually(func(g gomega.Gomega) {
						// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are increasing by 5 in each step within 5 minutes
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(testcontext.Get(t))
							require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
							l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep index", i).Msg("Number of upkeeps performed")
							g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
								"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
						}
					}, "5m", "1s").Should(gomega.Succeed())
				}
			}

			// Cancel all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := a.Registry.CancelUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not cancel upkeep at index %d", i)
			}

			var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterCancellation[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int64("Upkeep Count", countersAfterCancellation[i].Int64()).Int("Upkeep Index", i).Msg("Cancelled upkeep")
			}

			l.Info().Msg("Making sure the counter stays consistent")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// Expect the counter to remain constant (At most increase by 1 to account for stale performs) because the upkeep was cancelled
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterCancellation[i].Int64()+1),
						"Expected consumer counter to remain less than or equal to %d, but got %d",
						countersAfterCancellation[i].Int64()+1, latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestSetUpkeepTriggerConfig(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, true, false, false, nil, &config)

			// Start log trigger based upkeeps for all consumers
			for i := 0; i < len(consumers); i++ {
				err := consumers[i].Start()
				if err != nil {
					return
				}
			}

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			l.Info().Msg("Waiting for all upkeeps to perform")
			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 5
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			topic0InBytesMatch := [32]byte{
				61, 83, 163, 149, 80, 224, 70, 136,
				6, 88, 39, 243, 187, 134, 88, 76,
				176, 7, 171, 158, 188, 167, 235,
				213, 40, 231, 48, 28, 156, 49, 235, 93,
			} // bytes representation of 0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d

			topic0InBytesNoMatch := [32]byte{
				62, 83, 163, 149, 80, 224, 70, 136,
				6, 88, 39, 243, 187, 134, 88, 76,
				176, 7, 171, 158, 188, 167, 235,
				213, 40, 231, 48, 28, 156, 49, 235, 93,
			} // changed the first byte from 61 to 62 to make it not match

			bytes0 := [32]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			} // bytes representation of 0x0000000000000000000000000000000000000000000000000000000000000000

			// Update the trigger config so no upkeeps are triggered
			for i := 0; i < len(consumers); i++ {
				upkeepAddr := consumers[i].Address()

				logTriggerConfigStruct := ac.IAutomationV21PlusCommonLogTriggerConfig{
					ContractAddress: common.HexToAddress(upkeepAddr),
					FilterSelector:  0,
					Topic0:          topic0InBytesNoMatch,
					Topic1:          bytes0,
					Topic2:          bytes0,
					Topic3:          bytes0,
				}
				encodedLogTriggerConfig, err := core.CompatibleUtilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
				if err != nil {
					return
				}

				err = a.Registry.SetUpkeepTriggerConfig(upkeepIDs[i], encodedLogTriggerConfig)
				require.NoError(t, err, "Could not set upkeep trigger config at index %d", i)
			}

			var countersAfterSetNoMatch = make([]*big.Int, len(upkeepIDs))

			// Wait for 10 seconds to let in-flight upkeeps finish
			time.Sleep(10 * time.Second)
			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterSetNoMatch[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int64("Upkeep Count", countersAfterSetNoMatch[i].Int64()).Int("Upkeep Index", i).Msg("Upkeep")
			}

			l.Info().Msg("Making sure the counter stays consistent")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// Expect the counter to remain constant (At most increase by 2 to account for stale performs) because the upkeep trigger config is not met
					bufferCount := int64(2)
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterSetNoMatch[i].Int64()+bufferCount),
						"Expected consumer counter to remain less than or equal to %d, but got %d",
						countersAfterSetNoMatch[i].Int64()+bufferCount, latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// Update the trigger config, so upkeeps start performing again
			for i := 0; i < len(consumers); i++ {
				upkeepAddr := consumers[i].Address()

				logTriggerConfigStruct := ac.IAutomationV21PlusCommonLogTriggerConfig{
					ContractAddress: common.HexToAddress(upkeepAddr),
					FilterSelector:  0,
					Topic0:          topic0InBytesMatch,
					Topic1:          bytes0,
					Topic2:          bytes0,
					Topic3:          bytes0,
				}
				encodedLogTriggerConfig, err := core.CompatibleUtilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
				if err != nil {
					return
				}

				err = a.Registry.SetUpkeepTriggerConfig(upkeepIDs[i], encodedLogTriggerConfig)
				require.NoError(t, err, "Could not set upkeep trigger config at index %d", i)
			}

			var countersAfterSetMatch = make([]*big.Int, len(upkeepIDs))

			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterSetMatch[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int64("Upkeep Count", countersAfterSetMatch[i].Int64()).Int("Upkeep Index", i).Msg("Upkeep")
			}

			// Wait for 30 seconds to make sure backend is ready
			time.Sleep(30 * time.Second)
			// Start the consumers again
			for i := 0; i < len(consumers); i++ {
				err := consumers[i].Start()
				if err != nil {
					return
				}
			}

			l.Info().Msg("Making sure the counter starts increasing again")
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := int64(5)
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", countersAfterSetMatch[i].Int64()+expect),
						"Expected consumer counter to be greater than %d, but got %d", countersAfterSetMatch[i].Int64()+expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer
		})
	}
}

func TestAutomationAddFunds(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")
			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(1), automationDefaultUpkeepGasLimit, false, false, false, nil, &config)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			l.Info().Msg("Making sure for 2m no upkeeps are performed")
			gom := gomega.NewGomegaWithT(t)
			// Since the upkeep is currently underfunded, check that it doesn't get executed
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
						"Expected consumer counter to remain zero, but got %d", counter.Int64())
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			// Grant permission to the registry to fund the upkeep
			err = a.LinkToken.Approve(a.Registry.Address(), big.NewInt(0).Mul(big.NewInt(9e18), big.NewInt(int64(len(upkeepIDs)))))
			require.NoError(t, err, "Could not approve permissions for the registry on the link token contract")

			l.Info().Msg("Adding funds to the upkeeps")
			for i := 0; i < len(upkeepIDs); i++ {
				// Add funds to the upkeep whose ID we know from above
				err = a.Registry.AddUpkeepFunds(upkeepIDs[i], big.NewInt(9e18))
				require.NoError(t, err, "Unable to add upkeep")
			}

			l.Info().Msg("Waiting for 2m for all contracts to perform at least one upkeep")
			// Now the new upkeep should be performing because we added enough funds
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[0].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion
		})
	}
}

func TestAutomationPauseUnPause(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false, false, nil, &config)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(5)),
						"Expected consumer counter to be greater than 5, but got %d", counter.Int64())
					l.Info().Int("Upkeep Index", i).Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			// pause all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := a.Registry.PauseUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not pause upkeep at index %d", i)
			}

			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterPause[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int("Upkeep Index", i).Int64("Upkeeps Performed", countersAfterPause[i].Int64()).Msg("Paused Upkeep")
			}

			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// In most cases counters should remain constant, but there might be a straggling perform tx which
					// gets committed later and increases counter by 1
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterPause[i].Int64()+1),
						"Expected consumer counter not have increased more than %d, but got %d",
						countersAfterPause[i].Int64()+1, latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// unpause all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := a.Registry.UnpauseUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not unpause upkeep at index %d", i)
			}

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5 + numbers of performing before pause
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", countersAfterPause[i].Int64()+1),
						"Expected consumer counter to be greater than %d, but got %d", countersAfterPause[i].Int64()+1, counter.Int64())
					l.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m to perform, 1m buffer
		})
	}
}

func TestAutomationRegisterUpkeep(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false, false, nil, &config)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			var initialCounters = make([]*big.Int, len(upkeepIDs))
			gom := gomega.NewGomegaWithT(t)
			// Observe that the upkeeps which are initially registered are performing and
			// store the value of their initial counters in order to compare later on that the value increased.
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
					l.Info().
						Int64("Upkeep counter", counter.Int64()).
						Int64("Upkeep index", int64(i)).
						Msg("Number of upkeeps performed")
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			newConsumers, _ := actions.RegisterNewUpkeeps(t, a.ChainClient, a.LinkToken,
				a.Registry, a.Registrar, automationDefaultUpkeepGasLimit, 1)

			// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
			newUpkeep := newConsumers[0]

			// Test that the newly registered upkeep is also performing.
			gom.Eventually(func(g gomega.Gomega) {
				counter, err := newUpkeep.Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				l.Info().Int64("Upkeeps Performed", counter.Int64()).Msg("Newly Registered Upkeep")
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for upkeep to perform, 1m buffer

			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")

					l.Info().
						Int64("Upkeep index", int64(i)).
						Int64("Upkeep counter", currentCounter.Int64()).
						Int64("initial counter", initialCounters[i].Int64()).
						Msg("Number of upkeeps performed")

					g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for upkeeps to perform, 1m buffer
		})
	}
}

func TestAutomationPauseRegistry(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false, false, nil, &config)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			// Observe that the upkeeps which are initially registered are performing
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d")
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// Pause the registry
			err = a.Registry.Pause()
			require.NoError(t, err, "Error pausing registry")

			// Store how many times each upkeep performed once the registry was successfully paused
			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterPause[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			}

			// After we paused the registry, the counters of all the upkeeps should stay constant
			// because they are no longer getting serviced
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.Equal(countersAfterPause[i].Int64()),
						"Expected consumer counter to remain constant at %d, but got %d",
						countersAfterPause[i].Int64(), latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestAutomationKeeperNodesDown(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false, false, nil, &config)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			nodesWithoutBootstrap := a.ChainlinkNodes[1:]

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			// Watch upkeeps being performed and store their counters in order to compare them later in the test
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					initialCounters[i] = counter
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// Take down 1 node. Currently, using 4 nodes so f=1 and is the max nodes that can go down.
			err = nodesWithoutBootstrap[0].MustDeleteJob("1")
			require.NoError(t, err, "Error deleting job from Chainlink node")

			l.Info().Msg("Successfully managed to take down the first half of the nodes")

			// Assert that upkeeps are still performed and their counters have increased
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(testcontext.Get(t))
					l.Info().Int64("Upkeeps Performed", currentCounter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(currentCounter.Int64()).Should(gomega.BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for each upkeep to perform once, 1m buffer

			// Take down the rest
			restOfNodesDown := nodesWithoutBootstrap[1:]
			for _, nodeToTakeDown := range restOfNodesDown {
				err = nodeToTakeDown.MustDeleteJob("1")
				require.NoError(t, err, "Error deleting job from Chainlink node")
			}
			l.Info().Msg("Successfully managed to take down the second half of the nodes")

			// See how many times each upkeep was executed
			var countersAfterNoMoreNodes = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterNoMoreNodes[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int("Upkeep Index", i).Int64("Performed", countersAfterNoMoreNodes[i].Int64()).Msg("Upkeeps Performed")
			}

			// Once all the nodes are taken down, there might be some straggling transactions which went through before
			// all the nodes were taken down
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterNoMoreNodes[i].Int64()+1),
						"Expected consumer counter to not have increased more than %d, but got %d",
						countersAfterNoMoreNodes[i].Int64()+1, latestCounter.Int64())
				}
			}, "2m", "1s").Should(gomega.Succeed())
		})
	}
}

func TestAutomationPerformSimulation(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumersPerformance, _ := actions.DeployPerformanceConsumers(
				t,
				a.ChainClient,
				a.Registry,
				a.Registrar,
				a.LinkToken,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
				&config,
			)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			consumerPerformance := consumersPerformance[0]

			// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
			gom.Consistently(func(g gomega.Gomega) {
				// Consumer count should remain at 0
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
			err = consumerPerformance.SetPerformGasToBurn(testcontext.Get(t), big.NewInt(100000))
			require.NoError(t, err, "Perform gas should be set successfully on consumer")

			// Upkeep should now start performing
			gom.Eventually(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m to perform once, 1m buffer
		})
	}
}

func TestAutomationCheckPerformGasLimit(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")
			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumersPerformance, upkeepIDs := actions.DeployPerformanceConsumers(
				t,
				a.ChainClient,
				a.Registry,
				a.Registrar,
				a.LinkToken,
				1, // It was impossible to investigate, why with multiple outputs it fails ONLY in CI and only for 2.1 and 2.2 versions
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
				&config,
			)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			nodesWithoutBootstrap := a.ChainlinkNodes[1:]

			// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
			l.Info().Msg("Making sure for 2m no upkeeps are performed")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					g.Expect(cnt.Int64()).Should(
						gomega.Equal(int64(0)),
						"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
					)
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			// Increase gas limit for the upkeep, higher than the performGasBurn
			l.Info().Msg("Increasing gas limit for upkeeps")
			for i := 0; i < len(upkeepIDs); i++ {
				err = a.Registry.SetUpkeepGasLimit(upkeepIDs[i], uint32(4500000))
				require.NoError(t, err, "Error setting upkeep gas limit")
			}

			// Upkeep should now start performing
			l.Info().Msg("Waiting for 4m for all contracts to perform at least one upkeep after gas limit increase")
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					l.Info().Int("Upkeep index", i).Int64("Upkeep counter", cnt.Int64()).Msg("Number of upkeeps performed")
					g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
					)
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m to perform once, 1m buffer

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			l.Info().Msg("Increasing check gas to burn for upkeeps")
			for i := 0; i < len(upkeepIDs); i++ {
				err = consumersPerformance[i].SetCheckGasToBurn(testcontext.Get(t), big.NewInt(3000000))
				require.NoError(t, err, "Check gas burn should be set successfully on consumer")
			}

			countPerID := make(map[*big.Int]*big.Int)

			// Get existing performed count
			l.Info().Msg("Getting existing performed count")
			for i := 0; i < len(upkeepIDs); i++ {
				existingCnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
				require.NoError(t, err, "Calling consumer's counter shouldn't fail")
				l.Info().
					Str("UpkeepID", upkeepIDs[i].String()).
					Int64("Upkeep counter", existingCnt.Int64()).
					Msg("Upkeep counter when check gas increased")
				countPerID[upkeepIDs[i]] = existingCnt
			}

			// In most cases count should remain constant, but it might increase by upto 1 due to pending perform
			l.Info().Msg("Waiting for 1m for all contracts to make sure they perform at maximum 1 upkeep")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
					existingCnt := countPerID[upkeepIDs[i]]
					g.Expect(cnt.Int64()).Should(
						gomega.BeNumerically("<=", existingCnt.Int64()+1),
						"Expected consumer counter to remain less than equal %d, but got %d", existingCnt.Int64()+1, cnt.Int64(),
					)
				}
			}, "1m", "1s").Should(gomega.Succeed())

			l.Info().Msg("Getting existing performed count")
			for i := 0; i < len(upkeepIDs); i++ {
				existingCnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
				require.NoError(t, err, "Calling consumer's counter shouldn't fail")
				l.Info().
					Str("UpkeepID", upkeepIDs[i].String()).
					Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when consistently block finished")
				countPerID[upkeepIDs[i]] = existingCnt
			}

			// Now increase checkGasLimit on registry
			highCheckGasLimit := actions.ReadRegistryConfig(config)
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			highCheckGasLimit.RegistryVersion = registryVersion

			ocrConfig, err := actions.BuildAutoOCR2ConfigVarsLocal(l, nodesWithoutBootstrap, highCheckGasLimit, a.Registrar.Address(), 30*time.Second, a.Registry.RegistryOwnerAddress(), a.Registry.ChainModuleAddress(), a.Registry.ReorgProtectionEnabled())
			require.NoError(t, err, "Error building OCR config")

			if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_0 {
				err = a.Registry.SetConfig(highCheckGasLimit, ocrConfig)
			} else {
				err = a.Registry.SetConfigTypeSafe(ocrConfig)
			}
			require.NoError(t, err, "Registry config should be set successfully!")

			l.Info().Msg("Waiting for 3m for all contracts to make sure they perform at maximum 1 upkeep after check gas limit increase")
			// Upkeep should start performing again, and it should get regularly performed
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					cnt, err := consumersPerformance[i].GetUpkeepCount(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					existingCnt := countPerID[upkeepIDs[i]]
					g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", existingCnt.Int64()),
						"Expected consumer counter to be greater than %d, but got %d", existingCnt.Int64(), cnt.Int64(),
					)
				}
			}, "3m", "1s").Should(gomega.Succeed()) // ~1m to setup cluster, 1m to perform once, 1m buffer
		})
	}
}

func TestUpdateCheckData(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			require.NoError(t, err, "Failed to get config")

			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			performDataChecker, upkeepIDs := actions.DeployPerformDataCheckerConsumers(
				t,
				a.ChainClient,
				a.Registry,
				a.Registrar,
				a.LinkToken,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				[]byte(automationExpectedData),
				&config,
			)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			gom.Consistently(func(g gomega.Gomega) {
				// expect the counter to remain 0 because perform data does not match
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := performDataChecker[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
						"Expected perform data checker counter to be 0, but got %d", counter.Int64())
					l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			for i := 0; i < len(upkeepIDs); i++ {
				err := a.Registry.UpdateCheckData(upkeepIDs[i], []byte(automationExpectedData))
				require.NoError(t, err, "Could not update check data for upkeep at index %d", i)
			}

			// retrieve new check data for all upkeeps
			for i := 0; i < len(upkeepIDs); i++ {
				upkeep, err := a.Registry.GetUpkeepInfo(testcontext.Get(t), upkeepIDs[i])
				require.NoError(t, err, "Failed to get upkeep info at index %d", i)
				require.Equal(t, []byte(automationExpectedData), upkeep.CheckData, "Upkeep data not as expected")
			}

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := performDataChecker[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected perform data checker counter to be greater than 0, but got %d", counter.Int64())
					l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m to perform once, 1m buffer
		})
	}
}

func TestSetOffchainConfigWithMaxGasPrice(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		// registry20 also has upkeep offchain config but the max gas price check is not implemented
		"registry_2_1": ethereum.RegistryVersion_2_1,
		"registry_2_2": ethereum.RegistryVersion_2_2,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			config, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
			if err != nil {
				t.Fatal(err)
			}
			a := setupAutomationTestDocker(
				t, registryVersion, actions.ReadRegistryConfig(config), false, false, &config,
			)

			sb, err := a.ChainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get start block")

			consumers, upkeepIDs := actions.DeployConsumers(t, a.ChainClient, a.Registry, a.Registrar, a.LinkToken, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false, false, nil, &config)

			t.Cleanup(func() {
				actions.GetStalenessReportCleanupFn(t, a.Logger, a.ChainClient, sb, a.Registry, registryVersion)()
			})

			gom := gomega.NewGomegaWithT(t)
			l.Info().Msg("waiting for all upkeeps to be performed at least once")
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d")
				}
			}, "4m", "5s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// set the maxGasPrice to 1 wei
			uoc, _ := cbor.Marshal(gasprice.UpkeepOffchainConfig{MaxGasPrice: big.NewInt(1)})
			l.Info().Msgf("setting all upkeeps' offchain config to %s, which means maxGasPrice is 1 wei", hexutil.Encode(uoc))
			for _, uid := range upkeepIDs {
				err = a.Registry.SetUpkeepOffchainConfig(uid, uoc)
				require.NoError(t, err, "Error setting upkeep offchain config")
			}

			// Store how many times each upkeep performed once their offchain config is set with maxGasPrice = 1 wei
			var countersAfterSettingLowMaxGasPrice = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterSettingLowMaxGasPrice[i], err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int64("Upkeep Performed times", countersAfterSettingLowMaxGasPrice[i].Int64()).Int("Upkeep index", i).Msg("Number of upkeeps performed after setting low max gas price")
			}

			var latestCounter *big.Int
			// the upkeepsPerformed of all the upkeeps should stay constant because they are no longer getting serviced
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err = consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterSettingLowMaxGasPrice[i].Int64()+1),
						"Expected consumer counter to be less than %d, but got %d",
						countersAfterSettingLowMaxGasPrice[i].Int64()+1, latestCounter.Int64())
				}
			}, "2m", "5s").Should(gomega.Succeed())
			l.Info().Msg("no upkeeps is performed because their max gas price is only 1 wei")

			// setting offchain config with a high max gas price for the first upkeep, it should perform again while
			// other upkeeps should not perform
			// set the maxGasPrice to 500 gwei for the first upkeep
			uoc, _ = cbor.Marshal(gasprice.UpkeepOffchainConfig{MaxGasPrice: big.NewInt(500_000_000_000)})
			l.Info().Msgf("setting the first upkeeps' offchain config to %s, which means maxGasPrice is 500 gwei", hexutil.Encode(uoc))
			err = a.Registry.SetUpkeepOffchainConfig(upkeepIDs[0], uoc)
			require.NoError(t, err, "Error setting upkeep offchain config")

			upkeepsPerformedBefore := make(map[int]int64)
			upkeepsPerformedAfter := make(map[int]int64)
			for i := 0; i < len(upkeepIDs); i++ {
				latestCounter, err = consumers[i].Counter(testcontext.Get(t))
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				upkeepsPerformedBefore[i] = latestCounter.Int64()
				upkeepsPerformedAfter[i] = latestCounter.Int64()

				l.Info().Int64("No of Upkeep Performed", latestCounter.Int64()).Str("Consumer address", consumers[i].Address()).Msg("Number of upkeeps performed just after setting offchain config")
			}

			// the upkeepsPerformed of all other upkeeps should stay constant because their max gas price remains very low.
			// consumer at index N, might not be correlated with upkeep at index N, so instead of focusing on one of them
			// we iterate over all of them and make sure that at most only one is performing upkeeps
			gom.Consistently(func(g gomega.Gomega) {
				activeConsumers := 0
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err = consumers[i].Counter(testcontext.Get(t))
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					if latestCounter.Int64() != upkeepsPerformedAfter[i] {
						activeConsumers++
					}
					upkeepsPerformedAfter[i] = latestCounter.Int64()
				}
				// 0 is also okay, because it means that no upkeep was performed yet
				g.Expect(activeConsumers).Should(gomega.BeNumerically("<=", 1), "Only one consumer should have been performing upkeeps, but %d did", activeConsumers)
			}, "2m", "5s").Should(gomega.Succeed())

			performingConsumerIndex := -1
			onlyOneConsumerPerformed := false
			for i := 0; i < len(upkeepIDs); i++ {
				if upkeepsPerformedAfter[i] > upkeepsPerformedBefore[i] {
					onlyOneConsumerPerformed = true
					performingConsumerIndex = i
					break
				}
			}

			for i := 0; i < len(upkeepIDs); i++ {
				l.Info().Int64("No of Upkeep Performed", latestCounter.Int64()).Str("Consumer address", consumers[i].Address()).Msg("Number of upkeeps performed after waiting for the results of offchain config change")
			}

			require.True(t, onlyOneConsumerPerformed, "Only one consumer should have been performing upkeeps")

			l.Info().Msg("all the rest upkeeps did not perform again because their max gas price remains 1 wei")
			l.Info().Msg("making sure the consumer keeps performing upkeeps because its max gas price is 500 gwei")

			// the first upkeep should start performing again
			gom.Eventually(func(g gomega.Gomega) {
				latestCounter, err = consumers[performingConsumerIndex].Counter(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), fmt.Sprintf("Failed to retrieve consumer counter for upkeep at index %d", performingConsumerIndex))
				g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically(">", countersAfterSettingLowMaxGasPrice[0].Int64()),
					"Expected consumer counter to be greater than %d, but got %d",
					countersAfterSettingLowMaxGasPrice[0].Int64(), latestCounter.Int64())
			}, "2m", "5s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once
			l.Info().Int64("Upkeep Performed times", latestCounter.Int64()).Msg("the first upkeep performed again")
		})
	}
}

func setupAutomationTestDocker(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
	isMercuryV02 bool,
	isMercuryV03 bool,
	automationTestConfig types.AutomationTestConfig,
) automationv2.AutomationTest {
	require.False(t, isMercuryV02 && isMercuryV03, "Cannot run test with both Mercury V02 and V03 on")

	l := logging.GetTestLogger(t)
	// Add registry version to config
	registryConfig.RegistryVersion = registryVersion

	// launch the environment
	var env *test_env.CLClusterTestEnv
	var err error
	require.NoError(t, err)
	l.Debug().Msgf("Funding amount: %f", *automationTestConfig.GetCommonConfig().ChainlinkNodeFunding)
	clNodesCount := 5

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, automationTestConfig)
	require.NoError(t, err, "Error building ethereum network config")

	if isMercuryV02 || isMercuryV03 {
		// start mock adapter only
		mockAdapterEnv, err := test_env.NewCLTestEnvBuilder().
			WithTestInstance(t).
			WithTestConfig(automationTestConfig).
			WithMockAdapter().
			WithoutCleanup().
			Build()
		require.NoError(t, err, "Error deploying test environment for Mercury")

		secretsConfig := `
		[Mercury.Credentials.cred1]
		LegacyURL = '%s'
		URL = '%s'
		Username = 'node'
		Password = 'nodepass'`
		secretsConfig = fmt.Sprintf(secretsConfig, mockAdapterEnv.MockAdapter.InternalEndpoint, mockAdapterEnv.MockAdapter.InternalEndpoint)

		builder, err := test_env.NewCLTestEnvBuilder().WithTestEnv(mockAdapterEnv)
		require.NoError(t, err, "Error building test environment for Mercury")

		env, err = builder.
			WithTestInstance(t).
			WithTestConfig(automationTestConfig).
			WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
			WithSecretsConfig(secretsConfig).
			WithCLNodes(clNodesCount).
			WithStandardCleanup().
			Build()
		require.NoError(t, err, "Error deploying test environment for Mercury")

		env.MockAdapter = mockAdapterEnv.MockAdapter
	} else {
		env, err = test_env.NewCLTestEnvBuilder().
			WithTestInstance(t).
			WithTestConfig(automationTestConfig).
			WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
			WithMockAdapter().
			WithCLNodes(clNodesCount).
			WithStandardCleanup().
			Build()
		require.NoError(t, err, "Error deploying test environment")
	}

	nodeClients := env.ClCluster.NodeAPIs()

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	sethClient, err := utils.TestAwareSethClient(t, automationTestConfig, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(*automationTestConfig.GetCommonConfig().ChainlinkNodeFunding))
	require.NoError(t, err, "Failed to fund the nodes")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()))
	})

	a := automationv2.NewAutomationTestDocker(l, sethClient, nodeClients, automationTestConfig)
	a.SetMercuryCredentialName("cred1")
	a.RegistrySettings = registryConfig
	a.RegistrarSettings = contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: uint8(2),
		AutoApproveMaxAllowed: 1000,
		MinLinkJuels:          big.NewInt(0),
	}
	a.PluginConfig = actions.ReadPluginConfig(automationTestConfig)
	a.PublicConfig = actions.ReadPublicConfig(automationTestConfig)

	a.SetupAutomationDeployment(t)
	a.SetDockerEnv(env)

	if isMercuryV02 || isMercuryV03 {
		require.False(t, a.IsOnk8s, "mercury mock is not supported on k8s")
		require.NotNil(t, a.DockerEnv, "docker env is not set")
		mercuryv03Mock200 := &parrot.Route{
			Method: http.MethodGet,
			Path:   "/api/v1/reports/bulk",
			ResponseBody: map[string]any{
				"reports": []map[string]any{
					{
						"feedID":                "0x00028c915d6af0fd66bba2d0fc9405226bca8d6806333121a7d9832103d1563c",
						"validFromTimestamp":    0,
						"observationsTimestamp": 0,
						"fullReport":            "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911",
					},
				},
			},
			ResponseStatusCode: http.StatusOK,
		}
		err = a.DockerEnv.MockAdapter.Client.RegisterRoute(mercuryv03Mock200)
		require.NoError(t, err, "Error setting mock adapter route")

		mercuryv02Mock200 := &parrot.Route{
			Method: http.MethodGet,
			Path:   "/client",
			ResponseBody: map[string]any{
				"chainlinkBlob": "0x0001c38d71fed6c320b90e84b6f559459814d068e2a1700adc931ca9717d4fe70000000000000000000000000000000000000000000000000000000001a80b52b4bf1233f9cb71144a253a1791b202113c4ab4a92fa1b176d684b4959666ff8200000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000260000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001004254432d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000645570be000000000000000000000000000000000000000000000000000002af2b818dc5000000000000000000000000000000000000000000000000000002af2426faf3000000000000000000000000000000000000000000000000000002af32dc209700000000000000000000000000000000000000000000000000000000012130f8df0a9745bb6ad5e2df605e158ba8ad8a33ef8a0acf9851f0f01668a3a3f2b68600000000000000000000000000000000000000000000000000000000012130f60000000000000000000000000000000000000000000000000000000000000002c4a7958dce105089cf5edb68dad7dcfe8618d7784eb397f97d5a5fade78c11a58275aebda478968e545f7e3657aba9dcbe8d44605e4c6fde3e24edd5e22c94270000000000000000000000000000000000000000000000000000000000000002459c12d33986018a8959566d145225f0c4a4e61a9a3f50361ccff397899314f0018162cf10cd89897635a0bb62a822355bd199d09f4abe76e4d05261bb44733d",
			},
			ResponseStatusCode: http.StatusOK,
		}
		err = a.DockerEnv.MockAdapter.Client.RegisterRoute(mercuryv02Mock200)
		require.NoError(t, err, "Error setting mock adapter route")
	}

	return *a
}
