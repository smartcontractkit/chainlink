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

var _ = Describe("OCR soak test @soak-ocr", func() {
	var (
		suiteSetup  actions.SuiteSetup
		networkInfo actions.NetworkInfo
		nodes       []client.Chainlink
		adapter     environment.ExternalAdapter
		perfTest    Test
		err         error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkCluster(5),
				client.DefaultNetworkFromConfig,
				"./integration-tests",
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			networkInfo = suiteSetup.DefaultNetwork()

			networkInfo.Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err := actions.FundChainlinkNodes(
				nodes,
				networkInfo.Client,
				networkInfo.Wallets.Default(),
				big.NewFloat(10),
				big.NewFloat(10),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the OCR soak test", func() {
			perfTest = NewOCRTest(
				OCRTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts: 5,
					},
					RoundTimeout: 180 * time.Second,
					AdapterValue: 5,
					TestDuration: 10 * time.Minute,
				},
				contracts.DefaultOffChainAggregatorOptions(),
				suiteSetup.Environment(),
				networkInfo.Client,
				networkInfo.Wallets,
				networkInfo.Deployer,
				adapter,
			)
			err = perfTest.Setup()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("OCR Soak test", func() {
		Measure("Measure OCR rounds", func(_ Benchmarker) {
			err = perfTest.Run()
			Expect(err).ShouldNot(HaveOccurred())
		}, 1)
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
