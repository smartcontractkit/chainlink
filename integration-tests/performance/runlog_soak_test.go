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

var _ = Describe("Runlog soak test @soak-runlog", func() {
	var (
		suiteSetup     actions.SuiteSetup
		defaultNetwork actions.NetworkInfo
		nodes          []client.Chainlink
		adapter        environment.ExternalAdapter
		perfTest       Test
		err            error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				// no need more than one node for runlog test
				environment.NewChainlinkCluster(1),
				client.DefaultNetworkFromConfig,
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())
			defaultNetwork = suiteSetup.DefaultNetwork()
			adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			defaultNetwork.Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err := actions.FundChainlinkNodes(
				nodes,
				defaultNetwork.Client,
				defaultNetwork.Wallets.Default(),
				big.NewFloat(10),
				big.NewFloat(10),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the Runlog soak test", func() {
			perfTest = NewRunlogTest(
				RunlogTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts: 5,
						RoundTimeout:      180 * time.Second,
						TestDuration:      1 * time.Minute,
					},
					AdapterValue: 5,
				},
				suiteSetup.Environment(),
				defaultNetwork.Link,
				defaultNetwork.Client,
				defaultNetwork.Wallets,
				defaultNetwork.Deployer,
				adapter,
			)
			err = perfTest.Setup()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Runlog soak test", func() {
		Measure("Measure Runlog rounds", func(_ Benchmarker) {
			err = perfTest.Run()
			Expect(err).ShouldNot(HaveOccurred())
		}, 1)
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
