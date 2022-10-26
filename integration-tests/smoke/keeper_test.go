package smoke

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

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

type KeeperTests int32

const (
	BasicSmokeTest KeeperTests = iota
	BcptTest
	PerformSimulationTest
	CheckPerformGasLimitTest
	RegisterUpkeepTest
	AddFundsToUpkeepTest
	RemovingKeeperTest
	PauseRegistryTest
	MigrateUpkeepTest
	HandleKeeperNodesGoingDown
	PauseUnpauseUpkeepTest
	UpdateCheckDataTest
)

type KeeperConsumerContracts int32

const (
	BasicCounter KeeperConsumerContracts = iota
	PerformanceCounter
	PerformDataChecker

	defaultUpkeepGasLimit             = uint32(2500000)
	defaultLinkFunds                  = int64(9e18)
	defaultUpkeepsToDeploy            = 10
	numUpkeepsAllowedForStragglingTxs = 6
	expectedData                      = "abcdef"
)

var defaultRegistryConfig = contracts.KeeperRegistrySettings{
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
}

var lowBCPTRegistryConfig = contracts.KeeperRegistrySettings{
	PaymentPremiumPPB:    uint32(200000000),
	FlatFeeMicroLINK:     uint32(0),
	BlockCountPerTurn:    big.NewInt(4),
	CheckGasLimit:        uint32(2500000),
	StalenessSeconds:     big.NewInt(90000),
	GasCeilingMultiplier: uint16(1),
	MinUpkeepSpend:       big.NewInt(0),
	MaxPerformGas:        uint32(5000000),
	FallbackGasPrice:     big.NewInt(2e11),
	FallbackLinkPrice:    big.NewInt(2e18),
}

var highBCPTRegistryConfig = contracts.KeeperRegistrySettings{
	PaymentPremiumPPB:    uint32(200000000),
	FlatFeeMicroLINK:     uint32(0),
	BlockCountPerTurn:    big.NewInt(10000),
	CheckGasLimit:        uint32(2500000),
	StalenessSeconds:     big.NewInt(90000),
	GasCeilingMultiplier: uint16(1),
	MinUpkeepSpend:       big.NewInt(0),
	MaxPerformGas:        uint32(5000000),
	FallbackGasPrice:     big.NewInt(2e11),
	FallbackLinkPrice:    big.NewInt(2e18),
}

var _ = Describe("Keeper Suite @keeper", func() {
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
			Entry("v1.1 Basic smoke test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, BasicSmokeTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.2 Basic smoke test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, BasicSmokeTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Basic smoke test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, BasicCounter, BasicSmokeTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.1 BCPT test @simulated", ethereum.RegistryVersion_1_1, highBCPTRegistryConfig, BasicCounter, BcptTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.2 BCPT test @simulated", ethereum.RegistryVersion_1_2, highBCPTRegistryConfig, BasicCounter, BcptTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 BCPT test @simulated", ethereum.RegistryVersion_1_3, highBCPTRegistryConfig, BasicCounter, BcptTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.2 Perform simulation test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, PerformSimulationTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Perform simulation test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, PerformanceCounter, PerformSimulationTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.2 Check/Perform Gas limit test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, CheckPerformGasLimitTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Check/Perform Gas limit test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, PerformanceCounter, CheckPerformGasLimitTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.1 Register upkeep test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.2 Register upkeep test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Register upkeep test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.1 Add funds to upkeep test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, AddFundsToUpkeepTest, big.NewInt(1)),
			Entry("v1.2 Add funds to upkeep test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, AddFundsToUpkeepTest, big.NewInt(1)),
			Entry("v1.3 Add funds to upkeep test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, BasicCounter, AddFundsToUpkeepTest, big.NewInt(1)),

			Entry("v1.1 Removing one keeper test @simulated", ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, RemovingKeeperTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.2 Removing one keeper test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, RemovingKeeperTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Removing one keeper test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, BasicCounter, RemovingKeeperTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.2 Pause registry test @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, PauseRegistryTest, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Pause registry test @simulated", ethereum.RegistryVersion_1_3, defaultRegistryConfig, BasicCounter, PauseRegistryTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.2 Migrate upkeep from a registry to another @simulated", ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, MigrateUpkeepTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.1 Handle keeper nodes going down @simulated", ethereum.RegistryVersion_1_1, lowBCPTRegistryConfig, BasicCounter, HandleKeeperNodesGoingDown, big.NewInt(defaultLinkFunds)),
			Entry("v1.2 Handle keeper nodes going down @simulated", ethereum.RegistryVersion_1_2, lowBCPTRegistryConfig, BasicCounter, HandleKeeperNodesGoingDown, big.NewInt(defaultLinkFunds)),
			Entry("v1.3 Handle keeper nodes going down @simulated", ethereum.RegistryVersion_1_3, lowBCPTRegistryConfig, BasicCounter, HandleKeeperNodesGoingDown, big.NewInt(defaultLinkFunds)),

			Entry("v1.3 Pause and unpause upkeeps @simulated", ethereum.RegistryVersion_1_3, lowBCPTRegistryConfig, BasicCounter, PauseUnpauseUpkeepTest, big.NewInt(defaultLinkFunds)),

			Entry("v1.3 Update check data @simulated", ethereum.RegistryVersion_1_3, lowBCPTRegistryConfig, PerformDataChecker, UpdateCheckDataTest, big.NewInt(defaultLinkFunds)),
		}
	)

	DescribeTable("Keeper Suite @keeper", func(
		registryVersion ethereum.KeeperRegistryVersion,
		registryConfig contracts.KeeperRegistrySettings,
		consumerContract KeeperConsumerContracts,
		testToRun KeeperTests,
		linkFundsForEachUpkeep *big.Int,
	) {
		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-keeper"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(eth.New(nil)).
			AddHelm(chainlink.New(0, map[string]interface{}{
				"replicas": "5",
				"env": map[string]interface{}{
					"MIN_INCOMING_CONFIRMATIONS": "1",
					"KEEPER_TURN_FLAG_ENABLED":   "true",
					"KEEPER_TURN_LOOK_BACK":      "0",
				},
			}))
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(networks.SimulatedEVM, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		contractDeployer, err = contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		chainClient.ParallelTransactions(true)

		By("Funding Chainlink nodes")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.5))
		Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")

		By("Deploy Keeper Contracts")
		linkToken, err = contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		switch consumerContract {
		case BasicCounter:
			registry, registrar, consumers, upkeepIDs = actions.DeployKeeperContracts(
				registryVersion,
				registryConfig,
				defaultUpkeepsToDeploy,
				defaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				linkFundsForEachUpkeep,
			)
		case PerformanceCounter:
			registry, registrar, consumersPerformance, upkeepIDs = actions.DeployPerformanceKeeperContracts(
				registryVersion,
				defaultUpkeepsToDeploy,
				defaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				&registryConfig,
				linkFundsForEachUpkeep,
				10000,   // How many blocks this upkeep will be eligible from first upkeep block
				5,       // Interval of blocks that upkeeps are expected to be performed
				100000,  // How much gas should be burned on checkUpkeep() calls
				4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than defaultUpkeepGasLimit
			)
		case PerformDataChecker:
			registry, registrar, performDataChecker, upkeepIDs = actions.DeployPerformDataCheckerContracts(
				registryVersion,
				defaultUpkeepsToDeploy,
				defaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				&registryConfig,
				linkFundsForEachUpkeep,
				[]byte(expectedData),
			)
		}

		By("Register Keeper Jobs")
		actions.CreateKeeperJobs(chainlinkNodes, registry)
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")

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
			}, "2m", "1s").Should(Succeed())

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
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(5)),
						"Expected perform data checker counter to be greater than 5, but got %d", counter.Int64())
					log.Info().Int64("Upkeep perform data checker", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "3m", "1s").Should(Succeed())
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
			}, "3m", "1s").Should(Succeed())

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
					// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
					// is sufficient to check that the upkeep count does not increase by more than 6.
					latestCounter, err := consumers[i].Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterPause[i].Int64()+numUpkeepsAllowedForStragglingTxs),
						"Expected consumer counter not have increased more than %d, but got %d",
						countersAfterPause[i].Int64()+numUpkeepsAllowedForStragglingTxs, latestCounter.Int64())
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
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(5)+countersAfterPause[i].Int64()),
						"Expected consumer counter to be greater than %d, but got %d", int64(5)+countersAfterPause[i].Int64(), counter.Int64())
					log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "3m", "1s").Should(Succeed())
		}

		if testToRun == BasicSmokeTest {
			By("watches all the registered upkeeps perform and then cancels them from the registry")
			Eventually(func(g Gomega) {
				// Check if the upkeeps are performing multiple times by analysing their counters and checking they are greater than 10
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Int64()).Should(BeNumerically(">", int64(10)),
						"Expected consumer counter to be greater than 10, but got %d", counter.Int64())
					log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")
				}
			}, "5m", "1s").Should(Succeed())

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
					// Expect the counter to remain constant because the upkeep was cancelled, so it shouldn't increase anymore
					latestCounter, err := consumers[i].Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(Equal(countersAfterCancellation[i].Int64()),
						"Expected consumer counter to remain constant at %d, but got %d",
						countersAfterCancellation[i].Int64(), latestCounter.Int64())
				}
			}, "1m", "1s").Should(Succeed())
		}

		if testToRun == BcptTest {
			By("tests that keeper pairs change turn every blockCountPerTurn")
			keepersPerformed := make([]string, 0)
			upkeepID := upkeepIDs[0]

			// Wait for upkeep to be performed twice by different keepers (buddies)
			Eventually(func(g Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Number of upkeeps performed")

				upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
				g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				log.Info().Str("keeper", latestKeeper).Msg("last keeper to perform upkeep")
				g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

				log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
				keepersPerformed = append(keepersPerformed, latestKeeper)
			}, "1m", "1s").Should(Succeed())

			Eventually(func(g Gomega) {
				upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
				g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

				log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
				keepersPerformed = append(keepersPerformed, latestKeeper)
			}, "1m", "1s").Should(Succeed())

			// Expect no new keepers to perform for a while
			Consistently(func(g Gomega) {
				upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
				g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).Should(BeElementOf(keepersPerformed), "Existing keepers should alternate turns within BCPT")
			}, "1m", "1s").Should(Succeed())

			// Now set BCPT to be low, so keepers change turn frequently
			err = registry.SetConfig(lowBCPTRegistryConfig)
			Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set config tx")

			// Expect a new keeper to perform
			Eventually(func(g Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Num upkeeps performed")

				upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
				g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

				latestKeeper := upkeepInfo.LastKeeper
				log.Info().Str("keeper", latestKeeper).Msg("last keeper to perform upkeep")
				g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
				g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

				log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
				keepersPerformed = append(keepersPerformed, latestKeeper)
			}, "1m", "1s").Should(Succeed())
		}

		if testToRun == PerformSimulationTest {
			By("tests that performUpkeep simulation is run before tx is broadcast")
			consumerPerformance := consumersPerformance[0]
			upkeepID := upkeepIDs[0]

			// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
			Consistently(func(g Gomega) {
				// Consumer count should remain at 0
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					Equal(int64(0)),
					"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
				)

				// Not even reverted upkeeps should be performed. Last keeper for the upkeep should be 0 address
				upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
				g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")
				g.Expect(upkeepInfo.LastKeeper).Should(Equal(actions.ZeroAddress.String()), "Last keeper should be zero address")
			}, "1m", "1s").Should(Succeed())

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
			}, "1m", "1s").Should(Succeed())
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
			}, "1m", "1s").Should(Succeed())

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
			}, "1m", "1s").Should(Succeed())

			// Now increase the checkGasBurn on consumer, upkeep should stop performing
			err = consumerPerformance.SetCheckGasToBurn(context.Background(), big.NewInt(3000000))
			Expect(err).ShouldNot(HaveOccurred(), "Check gas burn should be set successfully on consumer")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetCheckGasToBurn tx")

			// Get existing performed count
			existingCnt, err := consumerPerformance.GetUpkeepCount(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
			log.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when check gas increased")

			// In most cases count should remain constant, but there might be a straggling perform tx which
			// gets committed later. Since every keeper node cannot have more than 1 straggling tx, it
			// is sufficient to check that the upkeep count does not increase by more than 6.
			Consistently(func(g Gomega) {
				cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(
					BeNumerically("<=", existingCnt.Int64()+numUpkeepsAllowedForStragglingTxs),
					"Expected consumer counter to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
				)
			}, "3m", "1s").Should(Succeed())

			existingCnt, err = consumerPerformance.GetUpkeepCount(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
			existingCntInt := existingCnt.Int64()
			log.Info().Int64("Upkeep counter", existingCntInt).Msg("Upkeep counter when consistently block finished")

			// Now increase checkGasLimit on registry
			highCheckGasLimit := defaultRegistryConfig
			highCheckGasLimit.CheckGasLimit = uint32(5000000)
			err = registry.SetConfig(highCheckGasLimit)
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
			}, "1m", "1s").Should(Succeed())
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
			}, "1m", "1s").Should(Succeed())

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
			}, "1m", "1s").Should(Succeed())

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
			}, "1m", "1s").Should(Succeed())
		}

		if testToRun == AddFundsToUpkeepTest {
			By("adds funds to a new underfunded upkeep")
			// Since the upkeep is currently underfunded, check that it doesn't get executed
			Consistently(func(g Gomega) {
				counter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counter.Int64()).Should(Equal(int64(0)),
					"Expected consumer counter to remain zero, but got %d", counter.Int64())
			}, "1m", "1s").Should(Succeed())

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
			}, "1m", "1s").Should(Succeed())
		}

		if testToRun == RemovingKeeperTest {
			By("removes one keeper and makes sure the upkeeps are not affected by this and still perform")
			var initialCounters = make([]*big.Int, len(upkeepIDs))
			// Make sure the upkeeps are running before we remove a keeper
			Eventually(func(g Gomega) {
				for upkeepID := 0; upkeepID < len(upkeepIDs); upkeepID++ {
					counter, err := consumers[upkeepID].Counter(context.Background())
					initialCounters[upkeepID] = counter
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep with ID "+strconv.Itoa(upkeepID))
					g.Expect(counter.Cmp(big.NewInt(0)) == 1, "Expected consumer counter to be greater than 0, but got %s", counter)
				}
			}, "1m", "1s").Should(Succeed())

			keepers, err := registry.GetKeeperList(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Encountered error when getting the list of Keepers")

			// Remove the first keeper from the list
			newKeeperList := keepers[1:]

			// Construct the addresses of the payees required by the SetKeepers function
			payees := make([]string, len(keepers)-1)
			for i := 0; i < len(payees); i++ {
				payees[i], err = chainlinkNodes[0].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred(), "Encountered error when building the payee list")
			}

			err = registry.SetKeepers(newKeeperList, payees)
			Expect(err).ShouldNot(HaveOccurred(), "Encountered error when setting the new list of Keepers")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")
			log.Info().Msg("Successfully removed keeper at address " + keepers[0] + " from the list of Keepers")

			// The upkeeps should still perform and their counters should have increased compared to the first check
			Eventually(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					counter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter"+
						" for upkeep at index "+strconv.Itoa(i))
					g.Expect(counter.Cmp(initialCounters[i]) == 1, "Expected consumer counter to be greater "+
						"than initial counter which was %s, but got %s", initialCounters[i], counter)
				}
			}, "1m", "1s").Should(Succeed())
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
			}, "1m", "1s").Should(Succeed())

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

		if testToRun == MigrateUpkeepTest {
			By("creates another registry and migrates one upkeep to the new registry")
			// Deploy the second registry, second registrar, and the same number of upkeeps as the first one
			secondRegistry, _, _, _ := actions.DeployKeeperContracts(
				ethereum.RegistryVersion_1_2,
				defaultRegistryConfig,
				defaultUpkeepsToDeploy,
				defaultUpkeepGasLimit,
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(defaultLinkFunds),
			)

			// Set the jobs for the second registry
			actions.CreateKeeperJobs(chainlinkNodes, secondRegistry)
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")

			err = registry.SetMigrationPermissions(common.HexToAddress(secondRegistry.Address()), 3)
			Expect(err).ShouldNot(HaveOccurred(), "Couldn't set bidirectional permissions for first registry")
			err := secondRegistry.SetMigrationPermissions(common.HexToAddress(registry.Address()), 3)
			Expect(err).ShouldNot(HaveOccurred(), "Couldn't set bidirectional permissions for second registry")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for setting the permissions")

			// Check that the first upkeep from the first registry is performing (before being migrated)
			Eventually(func(g Gomega) {
				counterBeforeMigration, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(counterBeforeMigration.Int64()).Should(BeNumerically(">", int64(0)),
					"Expected consumer counter to be greater than 0, but got %s", counterBeforeMigration)
			}, "1m", "1s").Should(Succeed())

			// Migrate the upkeep with index 0 from the first to the second registry
			err = registry.Migrate([]*big.Int{upkeepIDs[0]}, common.HexToAddress(secondRegistry.Address()))
			Expect(err).ShouldNot(HaveOccurred(), "Couldn't migrate the first upkeep")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for the migration")

			// Pause the first registry, in that way we make sure that the upkeep is being performed by the second one
			err = registry.Pause()
			Expect(err).ShouldNot(HaveOccurred(), "Could not pause the registry")
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for the pausing of the first registry")

			counterAfterMigration, err := consumers[0].Counter(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")

			// Check that once we migrated the upkeep, the counter has increased
			Eventually(func(g Gomega) {
				currentCounter, err := consumers[0].Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
				g.Expect(currentCounter.Int64()).Should(BeNumerically(">", counterAfterMigration.Int64()),
					"Expected counter to have increased, but stayed constant at %s", counterAfterMigration)
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
			}, "1m", "1s").Should(Succeed())

			// Take down half of the Keeper nodes by deleting the Keeper job registered above (after registry deployment)
			firstHalfToTakeDown := chainlinkNodes[:len(chainlinkNodes)/2+1]
			for _, nodeToTakeDown := range firstHalfToTakeDown {
				err = nodeToTakeDown.MustDeleteJob("1")
				Expect(err).ShouldNot(HaveOccurred(), "Could not delete the job from one of the nodes")
				err = chainClient.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Error deleting the Keeper job from the node")
			}
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
			}, "3m", "1s").Should(Succeed())

			// Take down the other half of the Keeper nodes
			secondHalfToTakeDown := chainlinkNodes[len(chainlinkNodes)/2+1:]
			for _, nodeToTakeDown := range secondHalfToTakeDown {
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
			// all the nodes were taken down. Every keeper node can have at most 1 straggling transaction per upkeep,
			// so a +6 on the upper limit side should be sufficient.
			Consistently(func(g Gomega) {
				for i := 0; i < len(upkeepIDs); i++ {
					latestCounter, err := consumers[i].Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer counter for upkeep at index "+strconv.Itoa(i))
					g.Expect(latestCounter.Int64()).Should(BeNumerically("<=", countersAfterNoMoreNodes[i].Int64()+numUpkeepsAllowedForStragglingTxs),
						"Expected consumer counter to not have increased more than %d, but got %d",
						countersAfterNoMoreNodes[i].Int64()+numUpkeepsAllowedForStragglingTxs, latestCounter.Int64())
				}
			}, "3m", "1s").Should(Succeed())
		}

		By("Printing gas stats")
		chainClient.GasStats().PrintStats()

		By("Tearing down the environment")
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	},
		testScenarios,
	)
})
