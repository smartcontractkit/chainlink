package smoke

import (
	"context"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var defaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
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

var _ = Describe("Automation OCR Suite @automation", func() {
	numberOfUpkeeps := 2
	var (
		err                  error
		chainClient          blockchain.EVMClient
		contractDeployer     contracts.ContractDeployer
		registry             contracts.KeeperRegistry
		registrar            contracts.KeeperRegistrar
		consumers            []contracts.KeeperConsumer
		consumersPerformance []contracts.KeeperConsumerPerformance
		performDataChecker   []contracts.KeeperPerformDataChecker
		upkeepIDs            []*big.Int
		linkToken            contracts.LinkToken
		chainlinkNodes       []*client.Chainlink
		testEnvironment      *environment.Environment

		testScenarios = []TableEntry{
			Entry("v2.0 Basic smoke test @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, BasicCounter, BasicSmokeTest, big.NewInt(defaultLinkFunds), numberOfUpkeeps),
			Entry("v2.0 Add funds to upkeep test @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, BasicCounter, AddFundsToUpkeepTest, big.NewInt(1), numberOfUpkeeps),
			Entry("v2.0 Pause and unpause upkeeps @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, BasicCounter, PauseUnpauseUpkeepTest, big.NewInt(defaultLinkFunds), numberOfUpkeeps),
			Entry("v2.0 Register upkeep test @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, BasicCounter, RegisterUpkeepTest, big.NewInt(defaultLinkFunds), numberOfUpkeeps),
			Entry("v2.0 Pause registry test @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, BasicCounter, PauseRegistryTest, big.NewInt(defaultLinkFunds), numberOfUpkeeps),
			Entry("v2.0 Handle f keeper nodes going down @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, BasicCounter, HandleKeeperNodesGoingDown, big.NewInt(defaultLinkFunds), numberOfUpkeeps),
			Entry("v2.0 Perform simulation test @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, PerformanceCounter, PerformSimulationTest, big.NewInt(defaultLinkFunds), 1),
			Entry("v2.0 Check/Perform Gas limit test @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, PerformanceCounter, CheckPerformGasLimitTest, big.NewInt(defaultLinkFunds), 1),
			Entry("v2.0 Update check data @simulated", ethereum.RegistryVersion_2_0, defaultOCRRegistryConfig, PerformDataChecker, UpdateCheckDataTest, big.NewInt(defaultLinkFunds), numberOfUpkeeps),
		}
	)

	AfterEach(func() {
		By("Tearing down the environment")
		chainClient.GasStats().PrintStats()
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("Automation OCR Suite @automation", func(
		registryVersion ethereum.KeeperRegistryVersion,
		registryConfig contracts.KeeperRegistrySettings,
		consumerContract KeeperConsumerContracts,
		testToRun KeeperTests,
		linkFundsForEachUpkeep *big.Int,
		numberOfUpkeeps int,
	) {
		By("Deploying the environment")
		network := networks.SimulatedEVM
		baseTOML := `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[Keeper]
TurnFlagEnabled = true
TurnLookBack = 0

[Keeper.Registry]
SyncInterval = '5m'
PerformGasOverhead = 150_000

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-automation"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(eth.New(nil)).
			AddHelm(chainlink.New(0, map[string]interface{}{
				"replicas": "5",
				"toml":     client.AddNetworksConfig(baseTOML, network),
			}))
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(network, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		contractDeployer, err = contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		chainClient.ParallelTransactions(true)

		By("Funding Chainlink nodes")
		txCost, err := chainClient.EstimateCostForChainlinkOperations(1000)
		Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
		Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")

		By("Deploy Link Token Contract")
		linkToken, err = contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		By("Deploy Registry and Registrar")
		registry, registrar = actions.DeployAutoOCRRegistryAndRegistrar(
			registryVersion,
			registryConfig,
			numberOfUpkeeps,
			linkToken,
			contractDeployer,
			chainClient,
		)

		By("Create OCR Automation Jobs")
		actions.CreateOCRKeeperJobs(chainlinkNodes, registry.Address(), network.ChainID)
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig := actions.BuildAutoOCR2ConfigVars(nodesWithoutBootstrap, registryConfig, registrar.Address())
		err = registry.SetConfig(defaultRegistryConfig, ocrConfig)
		Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
		Expect(chainClient.WaitForEvents()).ShouldNot(HaveOccurred(), "Waiting for config to be set")

		By("Deploy Consumers")
		switch consumerContract {
		case BasicCounter:
			consumers, upkeepIDs = actions.DeployConsumers(
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				numberOfUpkeeps,
				linkFundsForEachUpkeep,
				defaultUpkeepGasLimit,
			)
		case PerformanceCounter:
			consumersPerformance, upkeepIDs = actions.DeployPerformanceConsumers(
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				numberOfUpkeeps,
				linkFundsForEachUpkeep,
				defaultUpkeepGasLimit,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
		case PerformDataChecker:
			performDataChecker, upkeepIDs = actions.DeployPerformDataCheckerConsumers(
				registry,
				registrar,
				linkToken,
				contractDeployer,
				chainClient,
				numberOfUpkeeps,
				linkFundsForEachUpkeep,
				defaultUpkeepGasLimit,
				[]byte(expectedData),
			)
		}

		// BasicCounter

		if testToRun == BasicSmokeTest {
			By("watches all the registered upkeeps perform and then cancels them from the registry")
			Eventually(func(g Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					expect := 5
					g.Expect(counter.Int64()).Should(BeNumerically(">=", int64(expect)),
						"Expected consumer counter to be greater than %d, but got %d", expect, counter.Int64())
					log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "5m", "1s").Should(Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			// Cancel all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.CancelUpkeep(upkeepIDs[i])
				Expect(err).ShouldNot(HaveOccurred(), "Could not cancel upkeep at index "+strconv.Itoa(i))
			}

			err := chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for upkeeps to be cancelled")

			var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterCancellation[i], err = consumers[i].Counter(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
				log.Info().Msg("Cancelled upkeep at index " + strconv.Itoa(i) + " which performed " +
					strconv.Itoa(int(countersAfterCancellation[i].Int64())) + " times")
			}

			Consistently(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// Expect the counter to remain constant (At most increase by 1 to account for stale performs) because the upkeep was cancelled
					latestCounter, err := consumers[i].Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterCancellation[i].Int64()+1),
						"Expected consumer counter to remain less than equal %d, but got %d", countersAfterCancellation[i].Int64()+1, latestCounter.Int64())
				}
			}, "1m", "1s").Should(Succeed())
		}

		if testToRun == PauseUnpauseUpkeepTest {
			By("watches all the registered upkeeps perform, pause and then unpause them from the registry")

			Eventually(func(g Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(5)),
						"Expected consumer counter to be greater than 5, but got %d", counter.Int64())
					log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "5m", "1s").Should(Succeed()) // ~1m for cluster setup, ~2m for performing each upkeep 5 times, ~2m buffer

			// pause all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.PauseUpkeep(upkeepIDs[i])
				Expect(err).ShouldNot(HaveOccurred(), "Could not pause upkeep at index "+strconv.Itoa(i))
			}

			err := chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for upkeeps to be paused")

			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				// Obtain the amount of times the upkeep has been executed so far
				countersAfterPause[i], err = consumers[i].Counter(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
				log.Info().Msg("Paused upkeep at index " + strconv.Itoa(i) + " which performed " +
					strconv.Itoa(int(countersAfterPause[i].Int64())) + " times")
			}

			Consistently(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					// In most cases counters should remain constant, but there might be a straggling perform tx which
					// gets committed later and increases counter by 1
					latestCounter, err := consumers[i].Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterPause[i].Int64()+1),
						"Expected consumer counter not have increased more than %d, but got %d",
						countersAfterPause[i].Int64()+1, latestCounter.Int64())
				}
			}, "1m", "1s").Should(Succeed())

			// unpause all the registered upkeeps via the registry
			for i := 0; i < len(upkeepIDs); i++ {
				err := registry.UnpauseUpkeep(upkeepIDs[i])
				Expect(err).ShouldNot(HaveOccurred(), "Could not unpause upkeep at index "+strconv.Itoa(i))
			}

			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for upkeeps to be unpaused")

			Eventually(func(g Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5 + numbers of performing before pause
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", countersAfterPause[i].Int64()+1),
						"Expected consumer counter to be greater than %d, but got %d", countersAfterPause[i].Int64()+1, counter.Int64())
					log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(Succeed()) // ~1m to perform, 1m buffer
		}

		if testToRun == RegisterUpkeepTest {
			By("registers a new upkeep after the initial one was already registered and watches all of them perform")
			var initialCounters = make([]*big.Int, len(upkeepIDs))

			// Observe that the upkeeps which are initially registered are performing and
			// store the value of their initial counters in order to compare later on that the value increased.
			Eventually(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
					log.Info().
						Int64("Upkeep counter", counter.Int64()).
						Int64("Upkeep ID", int64(i)).
						Msg("Number of upkeeps performed")
				}
			}, "4m", "1s").Should(Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			newConsumers, _ := actions.RegisterNewUpkeeps(contractDeployer, chainClient, linkToken,
				registry, registrar, defaultUpkeepGasLimit, 1)

			// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
			newUpkeep := newConsumers[0]

			// Test that the newly registered upkeep is also performing.
			Eventually(func(g Gomega) {
				counter, err := newUpkeep.Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
				log.Info().Msg("Newly registered upkeeps performed " + strconv.Itoa(int(counter.Int64())) + " times")
			}, "2m", "1s").Should(Succeed()) // ~1m for upkeep to perform, 1m buffer

			Eventually(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")

					log.Info().
						Int64("Upkeep ID", int64(i)).
						Int64("Upkeep counter", currentCounter.Int64()).
						Int64("initial counter", initialCounters[i].Int64()).
						Msg("Number of upkeeps performed")

					g.Expect(currentCounter.Int64()).Should(BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
			}, "2m", "1s").Should(Succeed()) // ~1m for upkeeps to perform, 1m buffer
		}

		if testToRun == AddFundsToUpkeepTest {
			By("adds funds to a new underfunded upkeep")
			// Since the upkeep is currently underfunded, check that it doesn't get executed
			Consistently(func(g Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(Equal(int64(0)),
					"Expected consumer counter to remain zero, but got %d", counter.Int64())
			}, "2m", "1s").Should(Succeed()) // ~1m for setup, 1m assertion

			// Grant permission to the registry to fund the upkeep
			err = linkToken.Approve(registry.Address(), big.NewInt(9e18))
			Expect(err).ShouldNot(HaveOccurred(), "Could not approve permissions for the registry "+
				"on the link token contract")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

			// Add funds to the upkeep whose ID we know from above
			err = registry.AddUpkeepFunds(upkeepIDs[0], big.NewInt(9e18))
			Expect(err).ShouldNot(HaveOccurred(), "Could not fund upkeep")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

			// Now the new upkeep should be performing because we added enough funds
			Eventually(func(g Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
					"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
			}, "2m", "1s").Should(Succeed()) // ~1m for perform, 1m buffer
		}

		if testToRun == PauseRegistryTest {
			By("pauses the registry and makes sure that the upkeeps are no longer performed")
			// Observe that the upkeeps which are initially registered are performing
			Eventually(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d")
				}
			}, "4m", "1s").Should(Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// Pause the registry
			err := registry.Pause()
			Expect(err).ShouldNot(HaveOccurred(), "Could not pause the registry")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

			// Store how many times each upkeep performed once the registry was successfully paused
			var countersAfterPause = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterPause[i], err = consumers[i].Counter(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
					" for upkeep at index "+strconv.Itoa(i))
			}

			// After we paused the registry, the counters of all the upkeeps should stay constant
			// because they are no longer getting serviced
			Consistently(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(Equal(countersAfterPause[i].Int64()),
						"Expected consumer counter to remain constant at %d, but got %d",
						countersAfterPause[i].Int64(), latestCounter.Int64())
				}
			}, "1m", "1s").Should(Succeed())
		}

		if testToRun == HandleKeeperNodesGoingDown {
			By("upkeeps are still performed if some keeper nodes go down")
			var initialCounters = make([]*big.Int, len(upkeepIDs))

			// Watch upkeeps being performed and store their counters in order to compare them later in the test
			Eventually(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					initialCounters[i] = counter
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
						"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
				}
			}, "4m", "1s").Should(Succeed()) // ~1m for cluster setup, ~1m for performing each upkeep once, ~2m buffer

			// Take down 1 node. Currently, using 4 nodes so f=1 and is the max nodes that can go down.
			err = nodesWithoutBootstrap[0].MustDeleteJob("1")
			Expect(err).ShouldNot(HaveOccurred(), "Could not delete the job from one of the nodes")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error deleting the Keeper job from the node")

			log.Info().Msg("Successfully managed to take down the first half of the nodes")

			// Assert that upkeeps are still performed and their counters have increased
			Eventually(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					currentCounter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(currentCounter.Int64()).Should(BeNumerically(">", initialCounters[i].Int64()),
						"Expected counter to have increased from initial value of %s, but got %s",
						initialCounters[i], currentCounter)
				}
			}, "2m", "1s").Should(Succeed()) // ~1m for each upkeep to perform once, 1m buffer

			// Take down the rest
			restOfNodesDown := nodesWithoutBootstrap[1:]
			for _, nodeToTakeDown := range restOfNodesDown {
				err = nodeToTakeDown.MustDeleteJob("1")
				Expect(err).ShouldNot(HaveOccurred(), "Could not delete the job from one of the nodes")
				err = chainClient.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Error deleting the Keeper job from the node")
			}
			log.Info().Msg("Successfully managed to take down the second half of the nodes")

			// See how many times each upkeep was executed
			var countersAfterNoMoreNodes = make([]*big.Int, len(upkeepIDs))
			for i := 0; i < len(upkeepIDs); i++ {
				countersAfterNoMoreNodes[i], err = consumers[i].Counter(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
				log.Info().Msg("Upkeep at index " + strconv.Itoa(i) + " performed " +
					strconv.Itoa(int(countersAfterNoMoreNodes[i].Int64())) + " times")
			}

			// Once all the nodes are taken down, there might be some straggling transactions which went through before
			// all the nodes were taken down
			Consistently(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterNoMoreNodes[i].Int64()+1),
						"Expected consumer counter to not have increased more than %d, but got %d",
						countersAfterNoMoreNodes[i].Int64()+1, latestCounter.Int64())
				}
			}, "2m", "1s").Should(Succeed())
		}

		// PerformanceCounter

		if testToRun == PerformSimulationTest {
			By("tests that performUpkeep simulation is run before tx is broadcast")
			consumerPerformance := consumersPerformance[0]

			// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
			Consistently(func(g Gomega) {
				// Consumer count should remain at 0
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)
			}, "2m", "1s").Should(Succeed()) // ~1m for setup, 1m assertion

			// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
			err = consumerPerformance.SetPerformGasToBurn(context.Background(), big.NewInt(100000))
			Expect(err).ShouldNot(HaveOccurred(), "Perform gas should be set successfully on consumer")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set perform gas tx")

			// Upkeep should now start performing
			Eventually(func(g Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
			}, "2m", "1s").Should(Succeed()) // ~1m to perform once, 1m buffer
		}

		if testToRun == CheckPerformGasLimitTest {
			By("tests that check/perform gas limits are respected for upkeeps")
			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

			// Initially performGas is set higher than defaultUpkeepGasLimit, so no upkeep should be performed
			Consistently(func(g Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)
			}, "2m", "1s").Should(Succeed()) // ~1m for setup, 1m assertion

			// Increase gas limit for the upkeep, higher than the performGasBurn
			err = registry.SetUpkeepGasLimit(upkeepID, uint32(4500000))
			Expect(err).ShouldNot(HaveOccurred(), "Upkeep gas limit should be set successfully")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetUpkeepGasLimit tx")

			// Upkeep should now start performing
			Eventually(func(g Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
				)
			}, "2m", "1s").Should(Succeed()) // ~1m to perform once, 1m buffer

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			err = consumerPerformance.SetCheckGasToBurn(context.Background(), big.NewInt(3000000))
			Expect(err).ShouldNot(HaveOccurred(), "Check gas burn should be set successfully on consumer")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetCheckGasToBurn tx")

			// Get existing performed count
			existingCnt, err := consumerPerformance.GetUpkeepCount(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
			log.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when check gas increased")

			// In most cases count should remain constant, but it might increase by upto 1 due to pending perform
			Consistently(func(g Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					BeNumerically("<=", existingCnt.Int64()+1),
					"Expected consumer counter to remain less than equal %d, but got %d", existingCnt.Int64()+1, cnt.Int64(),
				)
			}, "1m", "1s").Should(Succeed())

			existingCnt, err = consumerPerformance.GetUpkeepCount(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
			existingCntInt := existingCnt.Int64()
			log.Info().Int64("Upkeep counter", existingCntInt).Msg("Upkeep counter when consistently block finished")

			// Now increase checkGasLimit on registry
			highCheckGasLimit := defaultRegistryConfig
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			ocrConfig := actions.BuildAutoOCR2ConfigVars(nodesWithoutBootstrap, highCheckGasLimit, registrar.Address())
			err = registry.SetConfig(highCheckGasLimit, ocrConfig)
			Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set config tx")

			// Upkeep should start performing again, and it should get regularly performed
			Eventually(func(g Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(BeNumerically(">", existingCntInt),
					"Expected consumer counter to be greater than %d, but got %d", existingCntInt, cnt.Int64(),
				)
			}, "3m", "1s").Should(Succeed()) // ~1m to setup cluster, 1m to perform once, 1m buffer
		}

		// PerformDataChecker

		if testToRun == UpdateCheckDataTest {
			By("tests that counters will be updated after their check data is updated")

			Consistently(func(g Gomega) {
				// expect the counter to remain 0 because perform data does not match
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := performDataChecker[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve perform data checker"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(Equal(int64(0)),
						"Expected perform data checker counter to be 0, but got %d", counter.Int64())
					log.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(Succeed()) // ~1m for setup, 1m assertion

			for i := 0; i < len(upkeepIDs); i++ {
				err = registry.UpdateCheckData(upkeepIDs[i], []byte(expectedData))
				Expect(err).ShouldNot(HaveOccurred(), "Could not update check data for upkeep at index "+strconv.Itoa(i))
			}

			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for check data update")

			// retrieve new check data for all upkeeps
			for i := 0; i < len(upkeepIDs); i++ {
				upkeep, err := registry.GetUpkeepInfo(context.Background(), upkeepIDs[i])
				Expect(err).ShouldNot(HaveOccurred(), "Failed to get upkeep info at index "+strconv.Itoa(i))
				Expect(upkeep.CheckData).Should(Equal([]byte(expectedData)), "Expect the check data to be %s, but got %s", expectedData, string(upkeep.CheckData))
			}

			Eventually(func(g Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 5
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := performDataChecker[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve perform data checker counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
						"Expected perform data checker counter to be greater than 0, but got %d", counter.Int64())
					log.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "2m", "1s").Should(Succeed()) // ~1m to perform once, 1m buffer
		}

	},
		testScenarios,
	)
})
