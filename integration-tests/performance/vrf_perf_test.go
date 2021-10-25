package performance

import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("VRF perf test @perf-vrf", func() {
	var (
		suiteSetup actions.SuiteSetup
		nodes      []client.Chainlink
		adapter    environment.ExternalAdapter
		perfTest   Test
		err        error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkCluster(1),
				client.NewNetworkFromConfigWithDefault(client.NetworkGethPerformance),
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			suiteSetup.DefaultNetwork().Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err := actions.FundChainlinkNodes(
				nodes,
				suiteSetup.DefaultNetwork().Client,
				suiteSetup.DefaultNetwork().Wallets.Default(),
				big.NewFloat(10),
				big.NewFloat(10),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the VRF perf test", func() {
			perfTest = NewVRFTest(
				VRFTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts:    50,
						NumberOfRounds:       5,
						RoundTimeout:         60 * time.Second,
						GracefulStopDuration: 10 * time.Second,
					},
				},
				suiteSetup.Environment(),
				suiteSetup.DefaultNetwork().Link,
				suiteSetup.DefaultNetwork().Client,
				suiteSetup.DefaultNetwork().Wallets,
				suiteSetup.DefaultNetwork().Deployer,
				adapter,
			)
			err = perfTest.Setup()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("VRF perf test", func() {
		Measure("Measure VRF request latency", func(b Benchmarker) {
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
