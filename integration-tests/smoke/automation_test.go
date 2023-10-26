package smoke

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	cltypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	evm21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	it_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"
)

var utilsABI = cltypes.MustGetABI(automation_utils_2_1.AutomationUtilsABI)

const (
	automationDefaultUpkeepGasLimit  = uint32(2500000)
	automationDefaultLinkFunds       = int64(9e18)
	automationDefaultUpkeepsToDeploy = 10
	automationExpectedData           = "abcdef"
	defaultAmountOfUpkeeps           = 2
)

var (
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
	automationDefaultRegistryConfig = contracts.KeeperRegistrySettings{
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

func TestMain(m *testing.M) {
	logging.Init()
	fmt.Printf("Running Smoke Test on %s\n", networks.SelectedNetwork.Name) // Print to get around disabled logging
	fmt.Printf("Chainlink Image %s\n", os.Getenv("CHAINLINK_IMAGE"))        // Print to get around disabled logging
	fmt.Printf("Chainlink Version %s\n", os.Getenv("CHAINLINK_VERSION"))    // Print to get around disabled logging
	os.Exit(m.Run())
}

func TestAutomationBasic(t *testing.T) {
	SetupAutomationBasic(t, false)
}

func SetupAutomationBasic(t *testing.T, nodeUpgrade bool) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0":                                 ethereum.RegistryVersion_2_0,
		"registry_2_1_conditional":                     ethereum.RegistryVersion_2_1,
		"registry_2_1_logtrigger":                      ethereum.RegistryVersion_2_1,
		"registry_2_1_with_mercury_v02":                ethereum.RegistryVersion_2_1,
		"registry_2_1_with_mercury_v03":                ethereum.RegistryVersion_2_1,
		"registry_2_1_with_logtrigger_and_mercury_v02": ethereum.RegistryVersion_2_1,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)

			var (
				upgradeImage   string
				upgradeVersion string
				err            error
				testName       = "basic-upkeep"
			)
			if nodeUpgrade {
				upgradeImage = os.Getenv("UPGRADE_IMAGE")
				upgradeVersion = os.Getenv("UPGRADE_VERSION")
				if len(upgradeImage) == 0 || len(upgradeVersion) == 0 {
					t.Fatal("UPGRADE_IMAGE and UPGRADE_VERSION must be set to upgrade nodes")
				}
				testName = "node-upgrade"
			}

			// Use the name to determine if this is a log trigger or mercury
			isLogTrigger := name == "registry_2_1_logtrigger" || name == "registry_2_1_with_logtrigger_and_mercury_v02"
			isMercuryV02 := name == "registry_2_1_with_mercury_v02" || name == "registry_2_1_with_logtrigger_and_mercury_v02"
			isMercuryV03 := name == "registry_2_1_with_mercury_v03"
			isMercury := isMercuryV02 || isMercuryV03

			chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupAutomationTestDocker(
				t, testName, registryVersion, defaultOCRRegistryConfig, nodeUpgrade, isMercuryV02, isMercuryV03,
			)

			consumers, upkeepIDs := actions.DeployConsumers(
				t,
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				isLogTrigger,
				isMercury,
			)

			for i := 0; i < len(upkeepIDs); i++ {
				if isLogTrigger || isMercuryV02 {
					if err := consumers[i].Start(); err != nil {
						l.Error().Msg("Error when starting consumer")
						return
					}
				}

				if isMercury {
					// Set privilege config to enable mercury
					privilegeConfigBytes, _ := json.Marshal(evm21.UpkeepPrivilegeConfig{
						MercuryEnabled: true,
					})
					if err := registry.SetUpkeepPrivilegeConfig(upkeepIDs[i], privilegeConfigBytes); err != nil {
						l.Error().Msg("Error when setting upkeep privilege config")
						return
					}
				}
			}

			l.Info().Msg("Waiting for all upkeeps to be performed")
			gom := gomega.NewGomegaWithT(t)
			startTime := time.Now()
			// TODO Tune this timeout window after stress testing
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
					expect := 5
					l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			l.Info().Msgf("Total time taken to get 5 performs for each upkeep: %s", time.Since(startTime))

			if nodeUpgrade {
				expect := 5
				// Upgrade the nodes one at a time and check that the upkeeps are still being performed
				for i := 0; i < 5; i++ {
					actions.UpgradeChainlinkNodeVersionsLocal(upgradeImage, upgradeVersion, testEnv.ClCluster.Nodes[i])
					time.Sleep(time.Second * 10)
					expect = expect + 5
					gom.Eventually(func(g gomega.Gomega) {
						// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are increasing by 5 in each step within 5 minutes
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(context.Background())
							require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
							l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep ID", i).Msg("Number of upkeeps performed")
							g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", int64(expect)),
								"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
						}
					}, "5m", "1s").Should(gomega.Succeed())
				}
			}

			// Cancel all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.CancelUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not cancel upkeep at index %d", i)
			}

			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error encountered when waiting for upkeeps to be cancelled")

			var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterCancellation[i], err = consumers[i].Counter(context.Background())
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int64("Upkeep Count", countersAfterCancellation[i].Int64()).Int("Upkeep Index", i).Msg("Cancelled upkeep")
			}

			l.Info().Msg("Making sure the counter stays consistent")
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// Expect the counter to remain constant (At most increase by 1 to account for stale performs) because the upkeep was cancelled
					latestCounter, err := consumers[i].Counter(context.Background())
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

	chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
		t, "set-trigger-config", ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, false, false, false,
	)

	consumers, upkeepIDs := actions.DeployConsumers(
		t,
		registry,
		registrar,
		linkToken,
		contractDeployer,
		chainClient,
		defaultAmountOfUpkeeps,
		big.NewInt(automationDefaultLinkFunds),
		automationDefaultUpkeepGasLimit,
		true,
		false,
	)

	// Start log trigger based upkeeps for all consumers
	for i := 0; i < len(consumers); i++ {
		err := consumers[i].Start()
		if err != nil {
			return
		}
	}

	l.Info().Msg("Waiting for all upkeeps to perform")
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		// Check if the upkeeps are performing multiple times by analyzing their counters
		for i := 0; i < len(upkeepIDs); i++ {
			counter, err := consumers[i].Counter(context.Background())
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

		logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
			ContractAddress: common.HexToAddress(upkeepAddr),
			FilterSelector:  0,
			Topic0:          topic0InBytesNoMatch,
			Topic1:          bytes0,
			Topic2:          bytes0,
			Topic3:          bytes0,
		}
		encodedLogTriggerConfig, err := utilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
		if err != nil {
			return
		}

		err = registry.SetUpkeepTriggerConfig(upkeepIDs[i], encodedLogTriggerConfig)
		require.NoError(t, err, "Could not set upkeep trigger config at index %d", i)
	}

	err := chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	var countersAfterSetNoMatch = make([]*big.Int, len(upkeepIDs))

	// Wait for 10 seconds to let in-flight upkeeps finish
	time.Sleep(10 * time.Second)
	for i := 0; i < len(upkeepIDs); i++ {
		// Obtain the amount of times the upkeep has been executed so far
		countersAfterSetNoMatch[i], err = consumers[i].Counter(context.Background())
		require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
		l.Info().Int64("Upkeep Count", countersAfterSetNoMatch[i].Int64()).Int("Upkeep Index", i).Msg("Upkeep")
	}

	l.Info().Msg("Making sure the counter stays consistent")
	gom.Consistently(func(g gomega.Gomega) {
		for i := 0; i < len(upkeepIDs); i++ {
			// Expect the counter to remain constant (At most increase by 2 to account for stale performs) because the upkeep trigger config is not met
			bufferCount := int64(2)
			latestCounter, err := consumers[i].Counter(context.Background())
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
			g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterSetNoMatch[i].Int64()+bufferCount),
				"Expected consumer counter to remain less than or equal to %d, but got %d",
				countersAfterSetNoMatch[i].Int64()+bufferCount, latestCounter.Int64())
		}
	}, "1m", "1s").Should(gomega.Succeed())

	// Update the trigger config, so upkeeps start performing again
	for i := 0; i < len(consumers); i++ {
		upkeepAddr := consumers[i].Address()

		logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
			ContractAddress: common.HexToAddress(upkeepAddr),
			FilterSelector:  0,
			Topic0:          topic0InBytesMatch,
			Topic1:          bytes0,
			Topic2:          bytes0,
			Topic3:          bytes0,
		}
		encodedLogTriggerConfig, err := utilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
		if err != nil {
			return
		}

		err = registry.SetUpkeepTriggerConfig(upkeepIDs[i], encodedLogTriggerConfig)
		require.NoError(t, err, "Could not set upkeep trigger config at index %d", i)
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	var countersAfterSetMatch = make([]*big.Int, len(upkeepIDs))

	for i := 0; i < len(upkeepIDs); i++ {
		// Obtain the amount of times the upkeep has been executed so far
		countersAfterSetMatch[i], err = consumers[i].Counter(context.Background())
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
			counter, err := consumers[i].Counter(context.Background())
			require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			expect := int64(5)
			l.Info().Int64("Upkeeps Performed", counter.Int64()).Int("Upkeep Index", i).Msg("Number of upkeeps performed")
			g.Expect(counter.Int64()).Should(gomega.BeNumerically(">=", countersAfterSetMatch[i].Int64()+expect),
				"Expected consumer counter to be greater than %d, but got %d", countersAfterSetMatch[i].Int64()+expect, counter.Int64())
		}
	}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer
}

func TestAutomationAddFunds(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "add-funds", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumers, upkeepIDs := actions.DeployConsumers(t, registry, registrar, linkToken, contractDeployer, chainClient, defaultAmountOfUpkeeps, big.NewInt(1), automationDefaultUpkeepGasLimit, false, false)

			gom := gomega.NewGomegaWithT(t)
			// Since the upkeep is currently underfunded, check that it doesn't get executed
			gom.Consistently(func(g gomega.Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
					"Expected consumer counter to remain zero, but got %d", counter.Int64())
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			// Grant permission to the registry to fund the upkeep
			err := linkToken.Approve(registry.Address(), big.NewInt(9e18))
			require.NoError(t, err, "Could not approve permissions for the registry on the link token contract")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Add funds to the upkeep whose ID we know from above
			err = registry.AddUpkeepFunds(upkeepIDs[0], big.NewInt(9e18))
			require.NoError(t, err, "Unable to add upkeep")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for events")

			// Now the new upkeep should be performing because we added enough funds
			gom.Eventually(func(g gomega.Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for perform, 1m buffer
		})
	}
}

func TestAutomationPauseUnPause(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "pause-unpause", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumers, upkeepIDs := actions.DeployConsumers(t, registry, registrar, linkToken, contractDeployer, chainClient, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false)

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analyzing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(5)),
						"Expected consumer counter to be greater than 5, but got %d", counter.Int64())
					l.Info().Int("Upkeep Index", i).Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "5m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			// pause all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.PauseUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not pause upkeep at index %d", i)
			}

			err := chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for upkeeps to be paused")

			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterPause[i], err = consumers[i].Counter(context.Background())
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int("Upkeep Index", i).Int64("Upkeeps Performed", countersAfterPause[i].Int64()).Msg("Paused Upkeep")
			}

			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// In most cases counters should remain constant, but there might be a straggling perform tx which
					// gets committed later and increases counter by 1
					latestCounter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(latestCounter.Int64()).Should(gomega.BeNumerically("<=", countersAfterPause[i].Int64()+1),
						"Expected consumer counter not have increased more than %d, but got %d",
						countersAfterPause[i].Int64()+1, latestCounter.Int64())
				}
			}, "1m", "1s").Should(gomega.Succeed())

			// unpause all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.UnpauseUpkeep(upkeepIDs[i])
				require.NoError(t, err, "Could not unpause upkeep at index %d", i)
			}

			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for upkeeps to be unpaused")

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5 + numbers of performing before pause
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
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
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "register-upkeep", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumers, upkeepIDs := actions.DeployConsumers(t, registry, registrar, linkToken, contractDeployer, chainClient, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false)

			var initialCounters = make([]*big.Int, len(upkeepIDs))
			gom := gomega.NewGomegaWithT(t)
			// Observe that the upkeeps which are initially registered are performing and
			// store the value of their initial counters in order to compare later on that the value increased.
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
					l.Info().
						Int64("Upkeep counter", counter.Int64()).
						Int64("Upkeep ID", int64(i)).
						Msg("Number of upkeeps performed")
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			newConsumers, _ := actions.RegisterNewUpkeeps(t, contractDeployer, chainClient, linkToken,
				registry, registrar, automationDefaultUpkeepGasLimit, 1)

			// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
			newUpkeep := newConsumers[0]

			// Test that the newly registered upkeep is also performing.
			gom.Eventually(func(g gomega.Gomega) {
				counter, err := newUpkeep.Counter(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				l.Info().Int64("Upkeeps Performed", counter.Int64()).Msg("Newly Registered Upkeep")
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for upkeep to perform, 1m buffer

			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")

					l.Info().
						Int64("Upkeep ID", int64(i)).
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
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "pause-registry", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumers, upkeepIDs := actions.DeployConsumers(t, registry, registrar, linkToken, contractDeployer, chainClient, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false)
			gom := gomega.NewGomegaWithT(t)

			// Observe that the upkeeps which are initially registered are performing
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d")
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// Pause the registry
			err := registry.Pause()
			require.NoError(t, err, "Error pausing registry")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for registry to pause")

			// Store how many times each upkeep performed once the registry was successfully paused
			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterPause[i], err = consumers[i].Counter(context.Background())
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
			}

			// After we paused the registry, the counters of all the upkeeps should stay constant
			// because they are no longer getting serviced
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(context.Background())
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
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			chainClient, chainlinkNodes, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "keeper-nodes-down", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumers, upkeepIDs := actions.DeployConsumers(t, registry, registrar, linkToken, contractDeployer, chainClient, defaultAmountOfUpkeeps, big.NewInt(automationDefaultLinkFunds), automationDefaultUpkeepGasLimit, false, false)
			gom := gomega.NewGomegaWithT(t)
			nodesWithoutBootstrap := chainlinkNodes[1:]

			var initialCounters = make([]*big.Int, len(upkeepIDs))

			// Watch upkeeps being performed and store their counters in order to compare them later in the test
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index %d", i)
					g.Expect(counter.Int64()).Should(gomega.BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
				}
			}, "4m", "1s").Should(gomega.Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// Take down 1 node. Currently, using 4 nodes so f=1 and is the max nodes that can go down.
			err := nodesWithoutBootstrap[0].MustDeleteJob("1")
			require.NoError(t, err, "Error deleting job from Chainlink node")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for blockchain events")

			l.Info().Msg("Successfully managed to take down the first half of the nodes")

			// Assert that upkeeps are still performed and their counters have increased
			gom.Eventually(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(context.Background())
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
				err = chainClient.WaitForEvents()
				require.NoError(t, err, "Error waiting for blockchain events")
			}
			l.Info().Msg("Successfully managed to take down the second half of the nodes")

			// See how many times each upkeep was executed
			var countersAfterNoMoreNodes = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterNoMoreNodes[i], err = consumers[i].Counter(context.Background())
				require.NoError(t, err, "Failed to retrieve consumer counter for upkeep at index %d", i)
				l.Info().Int("Upkeep Index", i).Int64("Performed", countersAfterNoMoreNodes[i].Int64()).Msg("Upkeeps Performed")
			}

			// Once all the nodes are taken down, there might be some straggling transactions which went through before
			// all the nodes were taken down
			gom.Consistently(func(g gomega.Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(context.Background())
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
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "perform-simulation", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumersPerformance, _ := actions.DeployPerformanceConsumers(
				t,
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
			gom := gomega.NewGomegaWithT(t)

			consumerPerformance := consumersPerformance[0]

			// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
			gom.Consistently(func(g gomega.Gomega) {
				// Consumer count should remain at 0
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
			err := consumerPerformance.SetPerformGasToBurn(context.Background(), big.NewInt(100000))
			require.NoError(t, err, "Perform gas should be set successfully on consumer")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for set perform gas tx")

			// Upkeep should now start performing
			gom.Eventually(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
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
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			chainClient, chainlinkNodes, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "gas-limit", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			consumersPerformance, upkeepIDs := actions.DeployPerformanceConsumers(
				t,
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
			gom := gomega.NewGomegaWithT(t)

			nodesWithoutBootstrap := chainlinkNodes[1:]
			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

			// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
			gom.Consistently(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					gomega.Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			// Increase gas limit for the upkeep, higher than the performGasBurn
			err := registry.SetUpkeepGasLimit(upkeepID, uint32(4500000))
			require.NoError(t, err, "Error setting upkeep gas limit")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for SetUpkeepGasLimit tx")

			// Upkeep should now start performing
			gom.Eventually(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m to perform once, 1m buffer

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			err = consumerPerformance.SetCheckGasToBurn(context.Background(), big.NewInt(3000000))
			require.NoError(t, err, "Check gas burn should be set successfully on consumer")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for SetCheckGasToBurn tx")

			// Get existing performed count
			existingCnt, err := consumerPerformance.GetUpkeepCount(context.Background())
			require.NoError(t, err, "Calling consumer's counter shouldn't fail")
			l.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when check gas increased")

			// In most cases count should remain constant, but it might increase by upto 1 due to pending perform
			gom.Consistently(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					gomega.BeNumerically("<=", existingCnt.Int64()+1),
					"Expected consumer counter to remain less than equal %d, but got %d", existingCnt.Int64()+1, cnt.Int64(),
				)
			}, "1m", "1s").Should(gomega.Succeed())

			existingCnt, err = consumerPerformance.GetUpkeepCount(context.Background())
			require.NoError(t, err, "Calling consumer's counter shouldn't fail")
			existingCntInt := existingCnt.Int64()
			l.Info().Int64("Upkeep counter", existingCntInt).Msg("Upkeep counter when consistently block finished")

			// Now increase checkGasLimit on registry
			highCheckGasLimit := automationDefaultRegistryConfig
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			highCheckGasLimit.RegistryVersion = registryVersion

			ocrConfig, err := actions.BuildAutoOCR2ConfigVarsLocal(l, nodesWithoutBootstrap, highCheckGasLimit, registrar.Address(), 30*time.Second, registry.RegistryOwnerAddress())
			require.NoError(t, err, "Error building OCR config")

			err = registry.SetConfig(highCheckGasLimit, ocrConfig)
			require.NoError(t, err, "Registry config should be set successfully!")
			err = chainClient.WaitForEvents()
			require.NoError(t, err, "Error waiting for set config tx")

			// Upkeep should start performing again, and it should get regularly performed
			gom.Eventually(func(g gomega.Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(gomega.BeNumerically(">", existingCntInt),
					"Expected consumer counter to be greater than %d, but got %d", existingCntInt, cnt.Int64(),
				)
			}, "3m", "1s").Should(gomega.Succeed()) // ~1m to setup cluster, 1m to perform once, 1m buffer
		})
	}
}

func TestUpdateCheckData(t *testing.T) {
	t.Parallel()
	registryVersions := map[string]ethereum.KeeperRegistryVersion{
		"registry_2_0": ethereum.RegistryVersion_2_0,
		"registry_2_1": ethereum.RegistryVersion_2_1,
	}

	for n, rv := range registryVersions {
		name := n
		registryVersion := rv
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			chainClient, _, contractDeployer, linkToken, registry, registrar, _ := setupAutomationTestDocker(
				t, "update-check-data", registryVersion, defaultOCRRegistryConfig, false, false, false,
			)

			performDataChecker, upkeepIDs := actions.DeployPerformDataCheckerConsumers(
				t,
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				defaultAmountOfUpkeeps,
				big.NewInt(automationDefaultLinkFunds),
				automationDefaultUpkeepGasLimit,
				[]byte(automationExpectedData),
			)
			gom := gomega.NewGomegaWithT(t)

			gom.Consistently(func(g gomega.Gomega) {
				// expect the counter to remain 0 because perform data does not match
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := performDataChecker[i].Counter(context.Background())
					g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to retrieve perform data checker"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(gomega.Equal(int64(0)),
						"Expected perform data checker counter to be 0, but got %d", counter.Int64())
					l.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(gomega.Succeed()) // ~1m for setup, 1m assertion

			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.UpdateCheckData(upkeepIDs[i], []byte(automationExpectedData))
				require.NoError(t, err, "Could not update check data for upkeep at index %d", i)
			}

			err := chainClient.WaitForEvents()
			require.NoError(t, err, "Error while waiting for check data update")

			// retrieve new check data for all upkeeps
			for i := 0; i < len(upkeepIDs); i++ {
				upkeep, err := registry.GetUpkeepInfo(context.Background(), upkeepIDs[i])
				require.NoError(t, err, "Failed to get upkeep info at index %d", i)
				require.Equal(t, []byte(automationExpectedData), upkeep.CheckData, "Upkeep data not as expected")
			}

			gom.Eventually(func(g gomega.Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := performDataChecker[i].Counter(context.Background())
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

type TestConfig struct {
	ChainlinkNodeFunding float64 `envconfig:"CHAINLINK_NODE_FUNDING" default:"1"`
}

func setupAutomationTestDocker(
	t *testing.T,
	testName string,
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
	statefulDb bool,
	isMercuryV02 bool,
	isMercuryV03 bool,
) (
	blockchain.EVMClient,
	[]*client.ChainlinkClient,
	contracts.ContractDeployer,
	contracts.LinkToken,
	contracts.KeeperRegistry,
	contracts.KeeperRegistrar,
	*test_env.CLClusterTestEnv,
) {
	require.False(t, isMercuryV02 && isMercuryV03, "Cannot run test with both Mercury V02 and V03 on")

	l := logging.GetTestLogger(t)
	// Add registry version to config
	registryConfig.RegistryVersion = registryVersion
	network := networks.SelectedNetwork

	// build the node config
	clNodeConfig := node.NewConfig(node.NewBaseConfig())
	syncInterval := models.MustMakeDuration(5 * time.Minute)
	clNodeConfig.Feature.LogPoller = it_utils.Ptr[bool](true)
	clNodeConfig.OCR2.Enabled = it_utils.Ptr[bool](true)
	clNodeConfig.Keeper.TurnLookBack = it_utils.Ptr[int64](int64(0))
	clNodeConfig.Keeper.Registry.SyncInterval = &syncInterval
	clNodeConfig.Keeper.Registry.PerformGasOverhead = it_utils.Ptr[uint32](uint32(150000))
	clNodeConfig.P2P.V2.AnnounceAddresses = &[]string{"0.0.0.0:6690"}
	clNodeConfig.P2P.V2.ListenAddresses = &[]string{"0.0.0.0:6690"}

	//launch the environment
	var env *test_env.CLClusterTestEnv
	var err error
	var testConfig TestConfig
	err = envconfig.Process("AUTOMATION", &testConfig)
	require.NoError(t, err)
	l.Debug().Msgf("Funding amount: %f", testConfig.ChainlinkNodeFunding)
	clNodesCount := 5
	if isMercuryV02 || isMercuryV03 {
		env, err = test_env.NewCLTestEnvBuilder().
			WithTestLogger(t).
			WithGeth().
			WithMockAdapter().
			WithFunding(big.NewFloat(testConfig.ChainlinkNodeFunding)).
			WithStandardCleanup().
			Build()
		require.NoError(t, err, "Error deploying test environment for Mercury")
		env.ParallelTransactions(true)

		secretsConfig := `
		[Mercury.Credentials.cred1]
		LegacyURL = '%s'
		URL = '%s'
		Username = 'node'
		Password = 'nodepass'`
		secretsConfig = fmt.Sprintf(secretsConfig, env.MockAdapter.InternalEndpoint, env.MockAdapter.InternalEndpoint)

		var httpUrls []string
		var wsUrls []string
		if network.Simulated {
			httpUrls = []string{env.Geth.InternalHttpUrl}
			wsUrls = []string{env.Geth.InternalWsUrl}
		} else {
			httpUrls = network.HTTPURLs
			wsUrls = network.URLs
		}

		node.SetChainConfig(clNodeConfig, wsUrls, httpUrls, network, false)

		err = env.StartClCluster(clNodeConfig, clNodesCount, secretsConfig)
		require.NoError(t, err, "Error starting CL nodes test environment for Mercury")
		err = env.FundChainlinkNodes(big.NewFloat(testConfig.ChainlinkNodeFunding))
		require.NoError(t, err, "Error funding CL nodes")

		if isMercuryV02 {
			output := `{"chainlinkBlob":"0x0001c38d71fed6c320b90e84b6f559459814d068e2a1700adc931ca9717d4fe70000000000000000000000000000000000000000000000000000000001a80b52b4bf1233f9cb71144a253a1791b202113c4ab4a92fa1b176d684b4959666ff8200000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000260000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001004254432d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000645570be000000000000000000000000000000000000000000000000000002af2b818dc5000000000000000000000000000000000000000000000000000002af2426faf3000000000000000000000000000000000000000000000000000002af32dc209700000000000000000000000000000000000000000000000000000000012130f8df0a9745bb6ad5e2df605e158ba8ad8a33ef8a0acf9851f0f01668a3a3f2b68600000000000000000000000000000000000000000000000000000000012130f60000000000000000000000000000000000000000000000000000000000000002c4a7958dce105089cf5edb68dad7dcfe8618d7784eb397f97d5a5fade78c11a58275aebda478968e545f7e3657aba9dcbe8d44605e4c6fde3e24edd5e22c94270000000000000000000000000000000000000000000000000000000000000002459c12d33986018a8959566d145225f0c4a4e61a9a3f50361ccff397899314f0018162cf10cd89897635a0bb62a822355bd199d09f4abe76e4d05261bb44733d"}`
			env.MockAdapter.SetStringValuePath("/client", []string{http.MethodGet, http.MethodPost}, map[string]string{"Content-Type": "application/json"}, output)
		}
		if isMercuryV03 {
			output := `{"reports":[{"feedID":"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000","validFromTimestamp":0,"observationsTimestamp":0,"fullReport":"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"}]}`
			env.MockAdapter.SetStringValuePath("/api/v1/reports/bulk", []string{http.MethodGet, http.MethodPost}, map[string]string{"Content-Type": "application/json"}, output)
		}
	} else {
		env, err = test_env.NewCLTestEnvBuilder().
			WithTestLogger(t).
			WithGeth().
			WithMockAdapter().
			WithCLNodes(clNodesCount).
			WithCLNodeConfig(clNodeConfig).
			WithFunding(big.NewFloat(testConfig.ChainlinkNodeFunding)).
			WithStandardCleanup().
			Build()
		require.NoError(t, err, "Error deploying test environment")
	}

	env.ParallelTransactions(true)
	nodeClients := env.ClCluster.NodeAPIs()
	workerNodes := nodeClients[1:]

	linkToken, err := env.ContractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying LINK token")

	registry, registrar := actions.DeployAutoOCRRegistryAndRegistrar(
		t,
		registryVersion,
		registryConfig,
		linkToken,
		env.ContractDeployer,
		env.EVMClient,
	)

	// Fund the registry with LINK
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(defaultAmountOfUpkeeps))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	err = actions.CreateOCRKeeperJobsLocal(l, nodeClients, registry.Address(), network.ChainID, 0, registryVersion)
	require.NoError(t, err, "Error creating OCR Keeper Jobs")
	ocrConfig, err := actions.BuildAutoOCR2ConfigVarsLocal(l, workerNodes, registryConfig, registrar.Address(), 30*time.Second, registry.RegistryOwnerAddress())
	require.NoError(t, err, "Error building OCR config vars")
	err = registry.SetConfig(automationDefaultRegistryConfig, ocrConfig)
	require.NoError(t, err, "Registry config should be set successfully")
	require.NoError(t, env.EVMClient.WaitForEvents(), "Waiting for config to be set")

	return env.EVMClient, nodeClients, env.ContractDeployer, linkToken, registry, registrar, env
}
