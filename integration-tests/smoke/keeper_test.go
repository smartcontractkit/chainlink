package smoke

//revive:disable:dot-imports
import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
)

var _ = Describe("Keeper suite @keeper", func() {
	var (
		err              error
		networks         *blockchain.Networks
		contractDeployer contracts.ContractDeployer
		registry         contracts.KeeperRegistry
		consumer         contracts.KeeperConsumer
		linkToken        contracts.LinkToken
		chainlinkNodes   []client.Chainlink
		env              *environment.Environment
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(
					environment.ChainlinkReplicas(6, config.ChainlinkVals()),
					"chainlink-keeper-core-ci",
					config.GethNetworks()...,
				),
				tools.ChartsRoot,
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

			r, consumers := actions.DeployKeeperContracts(
				1,
				linkToken,
				contractDeployer,
				networks,
			)
			consumer = consumers[0]
			registry = r
		})

		By("Register Keeper Jobs", func() {
			actions.CreateKeeperJobs(chainlinkNodes, registry)
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Error creating keeper jobs")
		})
	})

	Describe("with Keeper job", func() {
		It("performs upkeep of a target contract", func() {
			Eventually(func(g Gomega) {
				cnt, err := consumer.Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)), "Expected consumer counter to be greater than 0, but got %d", cnt.Int64())
				log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
			}, "2m", "1s").Should(Succeed())
		})
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
})
