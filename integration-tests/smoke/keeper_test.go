package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strconv"

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
	AddFundsToUpkeepTest
	RemovingKeeperTest
	PauseRegistryTest
)

type KeeperConsumerContracts int32

const (
	BasicCounter KeeperConsumerContracts = iota
	PerformanceCounter
)

const upkeepGasLimit = uint32(2500000)

var _ = Describe("Keeper v1.1 basic smoke test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, BasicSmokeTest))
var _ = Describe("Keeper v1.2 basic smoke test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, BasicSmokeTest))
var _ = Describe("Keeper v1.1 BCPT test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, highBCPTRegistryConfig, BasicCounter, BcptTest))
var _ = Describe("Keeper v1.2 BCPT test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, highBCPTRegistryConfig, BasicCounter, BcptTest))
var _ = Describe("Keeper v1.2 Perform simulation test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, PerformSimulationTest))
var _ = Describe("Keeper v1.2 Check/Perform Gas limit test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, PerformanceCounter, CheckPerformGasLimitTest))
var _ = Describe("Keeper v1.1 Register upkeep test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest))
var _ = Describe("Keeper v1.2 Register upkeep test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, RegisterUpkeepTest))
var _ = Describe("Keeper v1.1 Add funds to upkeep test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, AddFundsToUpkeepTest))
var _ = Describe("Keeper v1.2 Add funds to upkeep test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, AddFundsToUpkeepTest))
var _ = Describe("Keeper v1.1 Removing one keeper test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicCounter, RemovingKeeperTest))
var _ = Describe("Keeper v1.2 Removing one keeper test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, RemovingKeeperTest))
var _ = Describe("Keeper v1.2 Pause registry test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicCounter, PauseRegistryTest))

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
			registrar            contracts.UpkeepRegistrar
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
				// Since this is a simulated chain, block numbers start from 0, we can't look back 
				config.ProjectConfig.FrameworkConfig.ChainlinkEnvValues["KEEPER_TURN_LOOK_BACK"] = "0"

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
					registry, registrar, consumers, upkeepIDs = actions.DeployKeeperContracts(
						registryVersion,
						registryConfig,
						10,
						upkeepGasLimit,
						linkToken,
						contractDeployer,
						networks,
					)
				case PerformanceCounter:
					registry, registrar, consumersPerformance, upkeepIDs = actions.DeployPerformanceKeeperContracts(
						registryVersion,
						10,
						upkeepGasLimit,
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
				It("registers a new upkeep after the initial one was already registered and watches all of them perform", func() {
					var initialCounters = make([]*big.Int, len(upkeepIDs))

					// Observe that the upkeeps which are initially registered are performing and
					// store the value of their initial counters in order to compare later on that the value increased.
					Eventually(func(g Gomega) {
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(context.Background())
							initialCounters[i] = counter
							g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
							g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
								"Expected consumer counter to be greater than 0, but got %d", counter.Int64())
							log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Upkeeps performed")
						}
					}, "1m", "1s").Should(Succeed())

					newConsumers, _ := actions.RegisterNewUpkeeps(contractDeployer, networks, linkToken,
						registry, registrar, upkeepGasLimit, 1)

					// We know that newConsumers has size 1, so we can just use the newly registered upkeep.
					newUpkeep := newConsumers[0]

					// Test that the newly registered upkeep is also performing.
					Eventually(func(g Gomega) {
						counter, err := newUpkeep.Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling newly deployed upkeep's counter shouldn't fail")
						g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
							"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
						log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Upkeeps performed")
					}, "1m", "1s").Should(Succeed())

					Eventually(func(g Gomega) {
						for i := 0; i < len(upkeepIDs); i++ {
							currentCounter, err := consumers[i].Counter(context.Background())
							Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
							Expect(initialCounters[i].Int64() < currentCounter.Int64()).To(BeTrue())
						}
					}, "1m", "1s").Should(Succeed())
				})
			}

			if testToRun == AddFundsToUpkeepTest {
				It("adds funds to a new underfunded upkeep", func() {
					listOfNewUpkeeps := actions.DeployKeeperConsumers(contractDeployer, networks, 1)
					newUpkeep := listOfNewUpkeeps[0]
					newUpkeepAddress := listOfNewUpkeeps[0].Address()

					req, err := registrar.EncodeRegisterRequest(
						fmt.Sprintf("upkeep_%d", len(upkeepIDs)),
						[]byte("0x1234"),
						newUpkeepAddress,
						upkeepGasLimit,
						networks.Default.GetDefaultWallet().Address(),
						[]byte("0x"),
						big.NewInt(1),
						0,
					)
					Expect(err).ShouldNot(HaveOccurred(), "Could not encode first register request")

					// We want the new upkeep to be initially underfunded, so just transfer a minuscule amount of LINK
					tx, err := linkToken.TransferAndCall(registrar.Address(), big.NewInt(1), req)
					Expect(err).ShouldNot(HaveOccurred(), "Could not transfer small amount of LINK")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

					receipt, err := networks.Default.GetTxReceipt(tx.Hash())
					Expect(err).ShouldNot(HaveOccurred(), "Could not obtain transaction receipt")

					var upkeepID *big.Int
					for _, rawLog := range receipt.Logs {
						parsedUpkeepId, err := registry.ParseUpkeepIdFromRegisteredLog(rawLog)
						if err == nil {
							upkeepID = parsedUpkeepId
							break
						}
					}
					Expect(upkeepID).ShouldNot(BeNil(), "Upkeep ID not found after registration")
					log.Info().Msg("Successfully registered new upkeep with ID " + upkeepID.String())

					// Since the upkeep is currently underfunded, check that it doesn't get executed for a while
					Consistently(func(g Gomega) {
						counter, err := newUpkeep.Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's counter shouldn't fail")
						g.Expect(counter.Int64()).Should(Equal(int64(0)),
							"Expected consumer counter to remain zero, but got %d", counter.Int64())
					}, "30s", "1s").Should(Succeed())

					// Create a new request for the register where we actually fund the upkeep with proper funds
					req, err = registrar.EncodeRegisterRequest(
						fmt.Sprintf("upkeep_%d", len(upkeepIDs)),
						[]byte("0x1234"),
						newUpkeepAddress,
						upkeepGasLimit,
						networks.Default.GetDefaultWallet().Address(),
						[]byte("0x"),
						big.NewInt(9e18),
						0,
					)
					Expect(err).ShouldNot(HaveOccurred(), "Could not encode second register request")

					// Transfer the funds to the newly registered upkeep
					tx, err = linkToken.TransferAndCall(registrar.Address(), big.NewInt(9e18), req)
					Expect(err).ShouldNot(HaveOccurred(), "Failed to fund the upkeep with LINK")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

					log.Info().Msg("Successfully funded the new upkeep")

					// Now the new upkeep should be performing because we added enough funds
					Eventually(func(g Gomega) {
						counter, err := newUpkeep.Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Couldn't retrieve the new upkeep's counter")
						g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
							"Expected newly registered upkeep's counter to be greater than 0, but got %d", counter.Int64())
						log.Info().Int64("Upkeep counter", counter.Int64()).Msg("Upkeeps performed")
					}, "30s", "1s").Should(Succeed())
				})
			}

			if testToRun == RemovingKeeperTest {
				It("removes one keeper and makes sure the upkeeps are not affected by this and still perform", func() {
					var initialCounters = make([]*big.Int, len(upkeepIDs))
					// Make sure the upkeeps are running before we remove a keeper
					Eventually(func(g Gomega) {
						for upkeepID := 0; upkeepID < len(upkeepIDs); upkeepID++ {
							counter, err := consumers[upkeepID].Counter(context.Background())
							initialCounters[upkeepID] = counter
							g.Expect(err).ShouldNot(HaveOccurred(), "Failed to get counter for upkeep "+strconv.Itoa(upkeepID))
							g.Expect(counter.Cmp(big.NewInt(0)) == 1, "Expected consumer counter to be greater than 0, but got %s", counter)
						}
					}, "1m", "1s").Should(Succeed())

					keepers, err := registry.GetKeeperList(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Encountered error when getting the list of keepers")

					// Remove the first keeper from the list
					newKeeperList := keepers[1:]

					// Construct the addresses of the payees required by the SetKeepers function
					payees := make([]string, len(keepers)-1)
					for i := 0; i < len(payees); i++ {
						payees[i], err = chainlinkNodes[0].PrimaryEthAddress()
						Expect(err).ShouldNot(HaveOccurred(), "Shouldn't encounter error when building the payee list")
					}

					err = registry.SetKeepers(newKeeperList, payees)
					Expect(err).ShouldNot(HaveOccurred(), "Encountered error when setting the new Keepers")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")
					log.Info().Msg("Successfully removed keeper at address " + keepers[0] + " from the list of Keepers")

					// The upkeeps should still perform and their counters should have increased compared to the first check
					Eventually(func(g Gomega) {
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(context.Background())
							g.Expect(err).ShouldNot(HaveOccurred(), "Failed to get counter for upkeep "+strconv.Itoa(i))
							g.Expect(counter.Cmp(initialCounters[i]) == 1, "Expected consumer counter to be greater "+
								"than initial counter which was %s, but got %s", initialCounters[i], counter)
						}
					}, "1m", "1s").Should(Succeed())
				})
			}

			if testToRun == PauseRegistryTest {
				It("pauses the registry and makes sure that the upkeeps are no longer performed", func() {
					// Observe that the upkeeps which are initially registered are performing
					Eventually(func(g Gomega) {
						for i := 0; i < len(upkeepIDs); i++ {
							counter, err := consumers[i].Counter(context.Background())
							g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer's counter")
							g.Expect(counter.Int64()).Should(BeNumerically(">", int64(0)),
								"Expected consumer counter to be greater than 0, but got %d")
						}
					}, "1m", "1s").Should(Succeed())

					// Pause the registry
					err := registry.Pause()
					Expect(err).ShouldNot(HaveOccurred(), "Could not pause the registry")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for events")

					// Store how many times each upkeep performed once the registry was successfully paused
					var countersAfterPause = make([]*big.Int, len(upkeepIDs))
					for i := 0; i < len(upkeepIDs); i++ {
						countersAfterPause[i], err = consumers[i].Counter(context.Background())
						Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer's counter")
					}

					// After we paused the registry, the counters of all the upkeeps should stay constant
					// because they are no longer getting serviced
					Consistently(func(g Gomega) {
						for i := 0; i < len(upkeepIDs); i++ {
							latestCounter, err := consumers[i].Counter(context.Background())
							g.Expect(err).ShouldNot(HaveOccurred(), "Failed to retrieve consumer's counter")
							g.Expect(latestCounter.Int64()).Should(Equal(countersAfterPause[i].Int64()),
								"Expected consumer counter to remain constant at %d, but got %d",
								countersAfterPause[i].Int64(), latestCounter.Int64())
						}
					}, "1m", "1s").Should(Succeed())
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
