package smoke

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/logging"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

const (
	automationDefaultUpkeepGasLimit  = uint32(2500000)
	automationDefaultLinkFunds       = int64(9e18)
	automationDefaultUpkeepsToDeploy = 10
	automationExpectedData           = "abcdef"
	defaultAmountOfUpkeeps           = 2
)

var (
	automationBaseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[Keeper]
TurnLookBack = 0

[Keeper.Registry]
SyncInterval = '5m'
PerformGasOverhead = 150_000

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`

	automationEnvVars = map[string]any{
		"FEATURE_LOG_POLLER":                   "true",
		"FEATURE_OFFCHAIN_REPORTING2":          "true",
		"FEATURE_OFFCHAIN_REPORTING":           "false",
		"KEEPER_TURN_LOOK_BACK":                "0",
		"KEEPER_REGISTRY_SYNC_INTERVAL":        "5m",
		"KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD": "150000",

		"P2PV2_ANNOUNCE_ADDRESSES": "0.0.0.0:6690",
		"P2PV2_LISTEN_ADDRESSES":   "0.0.0.0:6690",
		"P2P_ANNOUNCE_IP":          "",
		"P2P_ANNOUNCE_PORT":        "",
		"P2P_BOOTSTRAP_PEERS":      "",
		"P2P_LISTEN_IP":            "",
		"P2P_LISTEN_PORT":          "",
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

func TestAutomatedBasic(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "basic-upkeep", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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

	// Cancel all the registered upkeeps via the registry
	for i := 0; i < len(upkeepIDs); i++ {
		err := registry.CancelUpkeep(upkeepIDs[i])
		require.NoError(t, err, "Could not cancel upkeep at index %d", i)
	}

	err := chainClient.WaitForEvents()
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
}

func TestAutomatedAddFunds(t *testing.T) {
	t.Parallel()

	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "add-funds", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

	consumers, upkeepIDs := actions.DeployConsumers(
		t,
		registry,
		registrar,
		linkToken,
		contractDeployer,
		chainClient,
		defaultAmountOfUpkeeps,
		big.NewInt(1),
		automationDefaultUpkeepGasLimit,
	)

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
}

func TestAutomatedPauseUnPause(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "pause-unpause", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
	)

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
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
}

func TestAutomatedRegisterUpkeep(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "register-upkeep", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
	)

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
}

func TestAutomatedPauseRegistry(t *testing.T) {
	t.Parallel()

	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "pause-registry", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
	)
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
}

func TestAutomatedKeeperNodesDown(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	chainClient, chainlinkNodes, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "keeper-nodes-down", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
	)
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
}

func TestAutomatedPerformSimulation(t *testing.T) {
	t.Parallel()

	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "perform-simulation", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
}

func TestAutomatedCheckPerformGasLimit(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	chainClient, chainlinkNodes, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "gas-limit", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
	ocrConfig := actions.BuildAutoOCR2ConfigVars(t, nodesWithoutBootstrap, highCheckGasLimit, registrar.Address(), 5*time.Second)
	err = registry.SetConfig(highCheckGasLimit, ocrConfig)
	require.NoError(t, err, "Registry config should be be set successfully")
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
}

func TestUpdateCheckData(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	chainClient, _, contractDeployer, linkToken, registry, registrar, onlyStartRunner := setupAutomationTest(
		t, "update-check-data", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig,
	)
	if onlyStartRunner {
		return
	}

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
}

func setupAutomationTest(
	t *testing.T,
	testName string,
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
) (
	chainClient blockchain.EVMClient,
	chainlinkNodes []*client.Chainlink,
	contractDeployer contracts.ContractDeployer,
	linkToken contracts.LinkToken,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	onlyStartRunner bool,
) {
	network := networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !network.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
		// For if we end up using env vars
		automationEnvVars["ETH_URL"] = network.URLs[0]
		automationEnvVars["ETH_HTTP_URL"] = network.HTTPURLs[0]
		automationEnvVars["ETH_CHAIN_ID"] = fmt.Sprint(network.ChainID)
	}
	chainlinkChart := chainlink.New(0, map[string]any{
		"replicas": "5",
		"toml":     client.AddNetworksConfig(automationBaseTOML, network),
	})

	useEnvVars := strings.ToLower(os.Getenv("TEST_USE_ENV_VAR_CONFIG"))
	if useEnvVars == "true" {
		chainlinkChart = chainlink.NewVersioned(0, "0.0.11", map[string]any{
			"replicas": "5",
			"env":      automationEnvVars,
		})
	}

	testEnvironment := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-automation-%s-%s", testName, strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlinkChart)
	err := testEnvironment.Run()

	require.NoError(t, err, "Error setting up test environment")

	onlyStartRunner = testEnvironment.WillUseRemoteRunner()
	if !onlyStartRunner {
		chainClient, err = blockchain.NewEVMClient(network, testEnvironment)
		require.NoError(t, err, "Error connecting to blockchain")
		contractDeployer, err = contracts.NewContractDeployer(chainClient)
		require.NoError(t, err, "Error building contract deployer")
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		require.NoError(t, err, "Error connecting to Chainlink nodes")
		chainClient.ParallelTransactions(true)

		txCost, err := chainClient.EstimateCostForChainlinkOperations(1000)
		require.NoError(t, err, "Error estimating cost for Chainlink Operations")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
		require.NoError(t, err, "Error funding Chainlink nodes")

		linkToken, err = contractDeployer.DeployLinkTokenContract()
		require.NoError(t, err, "Error deploying LINK token")

		registry, registrar = actions.DeployAutoOCRRegistryAndRegistrar(
			t,
			registryVersion,
			registryConfig,
			defaultAmountOfUpkeeps,
			linkToken,
			contractDeployer,
			chainClient,
		)

		actions.CreateOCRKeeperJobs(t, chainlinkNodes, registry.Address(), network.ChainID, 0)
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig := actions.BuildAutoOCR2ConfigVars(t, nodesWithoutBootstrap, registryConfig, registrar.Address(), 5*time.Second)
		err = registry.SetConfig(automationDefaultRegistryConfig, ocrConfig)
		require.NoError(t, err, "Registry config should be be set successfully")
		require.NoError(t, chainClient.WaitForEvents(), "Waiting for config to be set")
		// Register cleanup for any test
		t.Cleanup(func() {
			err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
			require.NoError(t, err, "Error tearing down environment")
		})
	}

	return chainClient, chainlinkNodes, contractDeployer, linkToken, registry, registrar, onlyStartRunner
}
