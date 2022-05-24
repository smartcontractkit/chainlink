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
)

var _ = Describe("Keeper v1.1 basic smoke test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, defaultRegistryConfig, BasicSmokeTest))
var _ = Describe("Keeper v1.2 basic smoke test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, defaultRegistryConfig, BasicSmokeTest))
var _ = Describe("Keeper v1.1 BCPT test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_1, highBCPTRegistryConfig, BcptTest))
var _ = Describe("Keeper v1.2 BCPT test @keeper", getKeeperSuite(ethereum.RegistryVersion_1_2, highBCPTRegistryConfig, BcptTest))

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

func getKeeperSuite(registryVersion ethereum.KeeperRegistryVersion, registryConfig contracts.KeeperRegistrySettings, testToRun KeeperTests) func() {
	return func() {
		var (
			err              error
			networks         *blockchain.Networks
			contractDeployer contracts.ContractDeployer
			registry         contracts.KeeperRegistry
			consumer         contracts.KeeperConsumer
			upkeepID         *big.Int
			linkToken        contracts.LinkToken
			chainlinkNodes   []client.Chainlink
			env              *environment.Environment
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
						config.GethNetworks()...,
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

				r, consumers, upkeepIDs := actions.DeployKeeperContracts(
					registryVersion,
					registryConfig,
					1,
					linkToken,
					contractDeployer,
					networks,
				)
				consumer = consumers[0]
				registry = r
				upkeepID = upkeepIDs[0]
			})

			By("Register Keeper Jobs", func() {
				actions.CreateKeeperJobs(chainlinkNodes, registry)
				err = networks.Default.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")
			})
		})

		Describe("with Keeper job", func() {
			if testToRun == BasicSmokeTest {
				It("performs upkeep of a target contract, stops upon cancel", func() {
					// Let upkeep be performed atleast once
					Eventually(func(g Gomega) {
						cnt, err := consumer.Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)), "Expected consumer counter to be greater than 0, but got %d", cnt.Int64())
						log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
					}, "2m", "1s").Should(Succeed())

					// Now cancel the upkeep as registry owner, it should get immediately cancelled
					err := registry.CancelUpkeep(upkeepID)
					Expect(err).ShouldNot(HaveOccurred(), "Upkeep should get cancelled successfully")
					err = networks.Default.WaitForEvents()
					Expect(err).ShouldNot(HaveOccurred(), "Error waiting for cancel upkeep tx")

					// Get existing performed count
					existingCnt, err := consumer.Counter(context.Background())
					Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					log.Info().Int64("Upkeep counter", existingCnt.Int64()).Msg("Upkeep cancelled")

					// Expect count to be remain consistent
					Consistently(func(g Gomega) {
						cnt, err := consumer.Counter(context.Background())
						g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
						g.Expect(cnt.Int64()).Should(
							Equal(existingCnt.Int64()),
							"Expected consumer counter to to remain constant at %d, but got %d", existingCnt.Int64(), cnt.Int64(),
						)
					}, "1m", "1s").Should(Succeed())
				})
			}

			if testToRun == BcptTest {
				It("tests that keeper pairs change turn every blockCountPerTurn", func() {
					keepersPerformed := make([]string, 0)

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
						g.Expect(latestKeeper).Should(BeElementOf(keepersPerformed), "Existing keepers should alternate turns wihtin BCPT")
					}, "30s", "1s").Should(Succeed())

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
