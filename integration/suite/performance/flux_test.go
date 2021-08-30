package performance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"time"
)

var _ = Describe("Performance tests", func() {
	var (
		s        *actions.DefaultSuiteSetup
		nodes    []client.Chainlink
		perfTest Test
		err      error
	)
	numberOfRounds := int64(5)
	numberOfNodes := 5

	BeforeEach(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(numberOfNodes),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())

			s.Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(
				nodes,
				s.Client,
				s.Wallets.Default(),
				big.NewFloat(2),
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the FluxAggregator performance test", func() {
			perfTest = NewFluxTest(
				FluxTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts: 100,
						NumberOfRounds:    numberOfRounds,
					},
					RequiredSubmissions: numberOfNodes,
					RestartDelayRounds:  0,
					NodePollTimePeriod:  time.Second * 15,
				},
				contracts.DefaultFluxAggregatorOptions(),
				s.Env,
				s.Client,
				s.Wallets,
				s.Deployer,
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
		By("Tearing down the environment", s.TearDown())
	})
})
