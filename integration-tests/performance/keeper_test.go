package performance

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = Describe("Keeper suite @keeper", func() {
	var (
		err              error
		chainClient      blockchain.EVMClient
		contractDeployer contracts.ContractDeployer
		registry         contracts.KeeperRegistry
		consumer         contracts.KeeperConsumer
		linkToken        contracts.LinkToken
		chainlinkNodes   []client.Chainlink
		testEnvironment  *environment.Environment
		profileTest      *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			testEnvironment = environment.New(&environment.Config{NamespacePrefix: "performance-keeper"}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(eth.New(nil)).
				AddHelm(chainlink.New(0, map[string]interface{}{
					"replicas": "5",
					"env": map[string]interface{}{
						"MIN_INCOMING_CONFIRMATIONS": "1",
						"KEEPER_TURN_FLAG_ENABLED":   "true",
						"HTTP_SERVER_WRITE_TIMEOUT":  "300s",
					},
				}))
			err = testEnvironment.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			chainClient, err = blockchain.NewEthereumMultiNodeClientSetup(blockchain.SimulatedEVMNetwork)(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			contractDeployer, err = contracts.NewContractDeployer(chainClient)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			chainClient.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := chainClient.EstimateCostForChainlinkOperations(10)
			Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
			Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")
		})

		By("Deploy Keeper Contracts", func() {
			linkToken, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

			r, _, consumers, _ := actions.DeployKeeperContracts(
				ethereum.RegistryVersion_1_1,
				contracts.KeeperRegistrySettings{
					PaymentPremiumPPB:    uint32(200000000),
					FlatFeeMicroLINK:     uint32(0),
					BlockCountPerTurn:    big.NewInt(3),
					CheckGasLimit:        uint32(2500000),
					StalenessSeconds:     big.NewInt(90000),
					GasCeilingMultiplier: uint16(1),
					MinUpkeepSpend:       big.NewInt(0),
					MaxPerformGas:        uint32(5000000),
					FallbackGasPrice:     big.NewInt(2e11),
					FallbackLinkPrice:    big.NewInt(2e18),
				},
				1,
				uint32(2500000), //upkeepGasLimit
				linkToken,
				contractDeployer,
				chainClient,
				big.NewInt(9e18),
			)
			consumer = consumers[0]
			registry = r
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
					// Not the last node, hence not all nodes started profiling yet.
					return
				}

				actions.CreateKeeperJobs(chainlinkNodes, registry)
				err = chainClient.WaitForEvents()
				Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")

				Eventually(func(g Gomega) {
					cnt, err := consumer.Counter(context.Background())
					g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
					g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)), "Expected consumer counter to be greater than 0, but got %d", cnt.Int64())
					log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
				}, "2m", "1s").Should(Succeed())
			}

			profileTest = testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
				ProfileFunction: profileFunction,
				ProfileDuration: 10 * time.Second,
				ChainlinkNodes:  chainlinkNodes,
			})
			profileTest.Setup(testEnvironment)
		})
	})

	Describe("with Keeper job", func() {
		It("performs upkeep of a target contract", func() {
			profileTest.Run()
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			chainClient.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, &profileTest.TestReporter, chainClient)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
