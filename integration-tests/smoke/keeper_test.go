package smoke

//revive:disable:dot-imports
import (
	"context"
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

type KeeperTests int32

const (
	BasicSmokeTest KeeperTests = iota
	BcptTest
	PerformSimulationTest
	CheckPerformGasLimitTest
	RegisterUpkeepTest
)

type KeeperConsumerContracts int32

const (
	BasicCounter KeeperConsumerContracts = iota
	PerformanceCounter
)

var _ = Describe("Keeper v1.1 basic smoke test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, BasicSmokeTest))
var _ = Describe("Keeper v1.2 basic smoke test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, BasicSmokeTest))
var _ = Describe("Keeper v1.1 BCPT test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, highBCPTRegistryConfig, BasicCounter, BcptTest))
var _ = Describe("Keeper v1.2 BCPT test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, highBCPTRegistryConfig, BasicCounter, BcptTest))
var _ = Describe("Keeper v1.2 Perform simulation test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, PerformSimulationTest))
var _ = Describe("Keeper v1.2 Check/Perform Gas limit test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, CheckPerformGasLimitTest))
var _ = Describe("Keeper v1.1 Register upkeep test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest))
var _ = Describe("Keeper v1.2 Register upkeep test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest))

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

func getKeeperSuite(
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
	consumerContract KeeperConsumerContracts,
	testToRun KeeperTests,
) func() {
	return func() {
		var (
			err                  error
			networks             *blockchain.Networks
			contractDeployer     contracts.ContractDeployer
			registry             contracts.KeeperRegistry
			consumers            []contracts.KeeperConsumer
			consumersPerformance []contracts.KeeperConsumerPerformance
			upkeepIDs            []*big.Int
			linkToken            contracts.LinkToken
			chainlinkNodes       []client.Chainlink
			env                  *environment.Environment
		)

		BeforeEach(func() {
			By("Deploying the environment", func() {
				// Confirm all logs, txs after 1 block
				config.ProjectConfig.FrameworkConfig.ChainlinkEnvValues["MIN_INCOMING_CONFIRMATIONS"] = "1"
				// Turn on buddy turn taking algo
				config.ProjectConfig.FrameworkConfig.ChainlinkEnvValues["KEEPER_TURN_FLAG_ENABLED"] = "true"

				env, err = environment.DeployOrLoadEnvironment(
					environment.NewChainlinkConfig(
						environment.ChainlinkReplicas(6, config.ChainlinkVals()),
						"chainlink-keeper-core-ci",
						environment.PerformanceGeth,
					),
				)
				Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
				err = env.ConnectAll()
				Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
			})

			By("Connecting to launched resources", func() {
				networkRegistry := blockchain.NewDefaultNetworkRegistry()
				networks, err = networkRegistry.GetNetworks(env)
				Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
				contractDeployer, err = contracts.NewContractDeployer(networks.Default)
				Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
				chainlinkNodes, err = client.ConnectChainlinkNodes(env)
				Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
				networks.Default.ParallelTransactions(true)
			})

			By("Funding Chainlink nodes", func() {
				txCost, err := networks.Default.EstimateCostForChainlinkOperations(10)
				Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
				err = actions.FundChainlinkNodes(chainlinkNodes, networks.Default, txCost)
				Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")
			})

			By("Deploy Keeper Contracts", func() {
				linkToken, err = contractDeployer.DeployLinkTokenContract()
				Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

				switch consumerContract {
				case BasicCounter:
					registry, consumers, upkeepIDs = actions.DeployKeeperContracts(
						registryVersion,
						registryConfig,
						10,
						uint32(2500000), //upkeepGasLimit
						linkToken,
						contractDeployer,
						networks,
					)
				case PerformanceCounter:
					registry, consumersPerformance, upkeepIDs = actions.DeployPerformanceKeeperContracts(
						registryVersion,
						10,
						uint32(2500000), //upkeepGasLimit
						linkToken,
						contractDeployer,
						networks,
						&registryConfig,
						10000,   // How many blocks this upkeep will be eligible from first upkeep block
						5,       // Interval of blocks that upkeeps are expected to be performed
						100000,  // How much gas should be burned on checkUpkeep() calls
						4000000, // How much gas should be burned on performUpkeep() calls. Initially set higher than upkeepGasLimit
					)
				}
			})

			By("Register Keeper Jobs", func() {
				actions.CreateKeeperJobs(chainlinkNodes, registry)
				err = networks.Default.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")
			})
		})

		Describe("with Keeper job", func() {
			if testToRun == BasicSmokeTest {
				It("watches all the registered upkeeps perform and then cancels them from the registry", func() {
					Eventually(func(g Gomega) {
						// Check if the upkeeps are performing by analysing their counters and checking they are greater than 0
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(context.Background())
							g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
							g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
								"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
							log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Upkeeps performed")
						}
					}, "1m", "1s").Should(Succeed())

					// Cancel all the registered upkeeps via the registry
					for i := 0; i < len(upkeepIDs); i++ {
						err := registry.CancelUpkeep(upkeepIDs[i])
						Expect(err).ShouldNot(HaveOccurred(), "Upkeep should get cancelled successfully")
					}

					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error encountered when waiting for upkeeps to be cancelled")

					var countersAfterCancellation = make([]*big.Int, len(upkeepIDs))

					for i := 0; i < len(upkeepIDs); i++ {
						// Obtain the amount of times the upkeep has been executed so far
						countersAfterCancellation[i], err = consumers[i].Counter(context.Background())
						Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
						log.Info().Int64("Upkeep counter", countersAfterCancellation[i].Int64()).Msg("Upkeep cancelled")
					}

					Consistently(func(g Gomega) {
						for i := 0; i < len(upkeepIDs); i++ {
							// Expect the counter to remain constant because the upkeep was cancelled, so it shouldn't increase anymore
							latestCounter, err := consumers[i].Counter(context.Background())
							g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
							g.Expect(latestCounter.Int64()).Should(Equal(countersAfterCancellation[i].Int64()),
								"Expected consumer counter to remain constant at %d, but got %d",
								countersAfterCancellation[i].Int64(), latestCounter.Int64())
						}
					}, "1m", "1s").Should(Succeed())
				})
			}

			if testToRun == BcptTest {
				It("tests that keeper pairs change turn every blockCountPerTurn", func() {
					keepersPerformed := make([]string, 0)
					upkeepID := upkeepIDs[0]

					// Wait for upkeep to be performed twice by different keepers (buddies)
					Eventually(func(g Gomega) {
						upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
						g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

						latestKeeper := upkeepInfo.LastKeeper
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
					lowBcpt := defaultRegistryConfig
					lowBcpt.BlockCountPerTurn = big.NewInt(5)
					err = registry.SetConfig(lowBcpt)
					Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set config tx")

					// Expect a new keeper to perform
					Eventually(func(g Gomega) {
						upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
						g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")

						latestKeeper := upkeepInfo.LastKeeper
						g.Expect(latestKeeper).ShouldNot(Equal(actions.ZeroAddress.String()), "Last keeper should be non zero")
						g.Expect(latestKeeper).ShouldNot(BeElementOf(keepersPerformed), "A new keeper node should perform this upkeep")

						log.Info().Str("keeper", latestKeeper).Msg("New keeper performed upkeep")
						keepersPerformed = append(keepersPerformed, latestKeeper)
					}, "1m", "1s").Should(Succeed())
				})
			}

			if testToRun == PerformSimulationTest {
				It("tests that performUpkeep simulation is run before tx is broadcast", func() {
					consumerPerformance := consumersPerformance[0]
					upkeepID := upkeepIDs[0]

					// Initially performGas is set high, so performUpkeep reverts and no upkeep should be performed
					Consistently(func(g Gomega) {
						// Consumer count should remain at 0
						cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(
							Equal(int64(0)),
							"Expected consumer counter to to remain constant at %d, but got %d", 0, cnt.Int64(),
						)

						// Not even reverted upkeeps should be performed. Last keeper for the upkeep should be 0 address
						upkeepInfo, err := registry.GetUpkeepInfo(context.Background(), upkeepID)
						g.Expect(err).ShouldNot(HaveOccurred(), "Registry's getUpkeep shouldn't fail")
						g.Expect(upkeepInfo.LastKeeper).Should(Equal(actions.ZeroAddress.String()), "Last keeper should be zero address")
					}, "1m", "1s").Should(Succeed())

					// Set performGas on consumer to be low, so that performUpkeep starts becoming successful
					err = consumerPerformance.SetPerformGasToBurn(context.Background(), big.NewInt(100000))
					Expect(err).ShouldNot(HaveOccurred(), "Perform gas should be set successfully on consumer")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set perform gas tx")

					// Upkeep should now start performing
					Eventually(func(g Gomega) {
						cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
							"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
						)
					}, "1m", "1s").Should(Succeed())

				})
			}

			if testToRun == CheckPerformGasLimitTest {
				It("tests that check/perform gas limits are respected for upkeeps", func() {
					consumerPerformance := consumersPerformance[0]
					upkeepID := upkeepIDs[0]

					// Initially performGas is set higher than upkeepGasLimit, so no upkeep should be performed
					Consistently(func(g Gomega) {
						cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(
							Equal(int64(0)),
							"Expected consumer counter to remain constant at %d, but got %d", 0, cnt.Int64(),
						)
					}, "1m", "1s").Should(Succeed())

					// Increase gas limit for the upkeep, higher than the performGasBurn
					err = registry.SetUpkeepGasLimit(upkeepID, uint32(4500000))
					Expect(err).ShouldNot(HaveOccurred(), "upkeep gas limit should be set successfully")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetUpkeepGasLimit tx")

					// Upkeep should now start performing
					Eventually(func(g Gomega) {
						cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
							"Expected consumer counter to be greater than 0, but got %d", cnt.Int64(),
						)
					}, "1m", "1s").Should(Succeed())

					// Now increase the checkGasBurn on consumer, upkeep should stop performing
					err = consumerPerformance.SetCheckGasToBurn(context.Background(), big.NewInt(3000000))
					Expect(err).ShouldNot(HaveOccurred(), "Check gas burn should be set successfully on consumer")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error waiting for SetCheckGasToBurn tx")

					// Get existing performed count, expect it to remain constant
					existingCnt, err := consumerPerformance.GetUpkeepCount(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					log.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep counter when check gas increased")
					Consistently(func(g Gomega) {
						cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(
							Equal(existingCnt.Int64()),
							"Expected consumer counter to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
						)
					}, "1m", "1s").Should(Succeed())

					// Now increase checkGasLimit on registry
					highCheckGasLimit := defaultRegistryConfig
					highCheckGasLimit.CheckGasLimit = uint32(5000000)
					err = registry.SetConfig(highCheckGasLimit)
					Expect(err).ShouldNot(HaveOccurred(), "Registry config should be be set successfully")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error waiting for set config tx")

					// Upkeep should start performing again
					Eventually(func(g Gomega) {
						cnt, err := consumerPerformance.GetUpkeepCount(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(BeNumerically(">", existingCnt.Int64()),
							"Expected consumer counter to be greater than %d, but got %d", existingCnt.Int64(), cnt.Int64(),
						)
					}, "1m", "1s").Should(Succeed())
				})
			}

			if testToRun == RegisterUpkeepTest {
				It("registers a new upkeep after the initial one was already registered and watches both perform", func() {
					var oldestUpkeepCounter *big.Int

					// Test that the upkeep which is registered in the BeforeEach function is executed
					Eventually(func(g Gomega) {
						oldestUpkeepCounter, err = consumer.Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(oldestUpkeepCounter.Int64()).Should(BeNumerically(">", int64(0)),
							"Expected consumer counter to be greater than 0, but got %d", oldestUpkeepCounter.Int64())
						log.Info().Int64("Upkeep counter", oldestUpkeepCounter.Int64()).Msg("Upkeeps performed")
					}, "1m", "1s").Should(Succeed())

					// Now register a new upkeep
					registry, consumers, _ := actions.DeployKeeperContracts(
						registryVersion,
						registryConfig,
						1,
						uint32(2500000), //upkeepGasLimit
						linkToken,
						contractDeployer,
						networks,
					)

					// Register the Keeper job responsible for the new upkeep.
					actions.CreateKeeperJobs(chainlinkNodes, registry)
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")

					// Test that the newly registered upkeep is also getting performed
					Eventually(func(g Gomega) {
						cnt, err := consumers[0].Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)),
							"Expected consumer counter to be greater than 0, but got %d", cnt.Int64())
						log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
					}, "1m", "1s").Should(Succeed())

					// Get the current counter for the upkeep which was registered first (in the BeforeEach function).
					existingCnt, err := consumer.Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					log.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep cancelled")

					// While we registered the second upkeep, the first upkeep should have continued performing and
					// therefore the counter should have increased in the meantime.
					Expect(oldestUpkeepCounter.Int64() <= existingCnt.Int64()).To(BeTrue())
				})
			}
		})

		AfterEach(func() {
			By("Printing gas stats", func() {
				networks.Default.GasStats().PrintStats()
			})
			By("Tearing down the environment", func() {
				err = actions.TeardownSuite(env, networks, utils.ProjectRoot, chainlinkNodes, nil)
				Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
			})
		})
	}
}
