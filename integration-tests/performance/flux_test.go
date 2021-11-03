//go:build performance

package performance

import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Performance tests @perf-flux", func() {
	var (
		suiteSetup     actions.SuiteSetup
		networkInfo    actions.NetworkInfo
		nodes          []client.Chainlink
		perfTest       Test
		err            error
		numberOfRounds int = 5
		numberOfNodes  int = 5
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkCluster(numberOfNodes),
				actions.EVMNetworkFromConfigHook,
				actions.EthereumDeployerHook,
				actions.EthereumClientHook,
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			networkInfo = suiteSetup.DefaultNetwork()

			networkInfo.Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(
				nodes,
				networkInfo.Client,
				networkInfo.Wallets.Default(),
				big.NewFloat(2),
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the FluxAggregator performance test", func() {
			perfTest = NewFluxTest(
				FluxTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts: 10,
						NumberOfRounds:    numberOfRounds,
					},
					RequiredSubmissions: numberOfNodes,
					RestartDelayRounds:  0,
					NodePollTimePeriod:  time.Second * 15,
				},
				contracts.DefaultFluxAggregatorOptions(),
				suiteSetup.Environment(),
				networkInfo.Client,
				networkInfo.Wallets,
				networkInfo.Deployer,
				nil,
			)
			err = perfTest.Setup()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("FluxMonitor", func() {
		Measure("Round latencies", func(b Benchmarker) {
			err = perfTest.Run()
			Expect(err).ShouldNot(HaveOccurred())
			err = perfTest.RecordValues(b)
			Expect(err).ShouldNot(HaveOccurred())
		}, 1)
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
