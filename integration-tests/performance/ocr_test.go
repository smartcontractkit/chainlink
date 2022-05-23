package performance

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

var _ = Describe("OCR Feed @ocr", func() {
	var (
		err               error
		env               *environment.Environment
		networks          *blockchain.Networks
		contractDeployer  contracts.ContractDeployer
		linkTokenContract contracts.LinkToken
		chainlinkNodes    []client.Chainlink
		mockserver        *client.MockserverClient
		ocrInstances      []contracts.OffchainAggregator
		profileTest       *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			// Increase HTTP_SERVER_WRITE_TIMEOUT to be larger than profile duration.
			config.ProjectConfig.FrameworkConfig.ChainlinkEnvValues["HTTP_SERVER_WRITE_TIMEOUT"] = "300s"

			env, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(
					environment.ChainlinkReplicas(6, config.ChainlinkVals()),
					"chainlink-ocr-profiling",
					config.GethNetworks()...,
				),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			err = env.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
		})

		By("Connecting to launched resources", func() {
			// Load Networks
			networkRegistry := blockchain.NewDefaultNetworkRegistry()
			var err error
			Expect(networkRegistry).ShouldNot(BeNil(), "nil NetworkRegistry")

			networks, err = networkRegistry.GetNetworks(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			contractDeployer, err = contracts.NewContractDeployer(networks.Default)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

			chainlinkNodes, err = client.ConnectChainlinkNodes(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			mockserver, err = client.ConnectMockServer(env)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")

			networks.Default.ParallelTransactions(true)
			Expect(err).ShouldNot(HaveOccurred())

			linkTokenContract, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
		})

		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(chainlinkNodes, networks.Default, big.NewFloat(.01))
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying OCR contracts", func() {
			ocrInstances = actions.DeployOCRContracts(1, linkTokenContract, contractDeployer, chainlinkNodes, networks)
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
					// Not the last node, hence not all nodes started profiling yet.
					return
				}
				By("Setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, mockserver))
				By("Creating OCR jobs", actions.CreateOCRJobs(ocrInstances, chainlinkNodes, mockserver))
				By("Starting new round", actions.StartNewRound(1, ocrInstances, networks))
				answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Getting latest answer from OCR contract shouldn't fail")
				Expect(answer.Int64()).Should(Equal(int64(5)), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

				By("setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, chainlinkNodes, mockserver))
				By("starting new round", actions.StartNewRound(2, ocrInstances, networks))

				answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(answer.Int64()).Should(Equal(int64(10)), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
			}

			profileTest = testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
				ProfileFunction: profileFunction,
				ProfileDuration: time.Minute,
				ChainlinkNodes:  chainlinkNodes,
			})
			profileTest.Setup(env)
		})
	})

	Describe("With a single OCR contract", func() {
		It("performs two rounds", func() {
			profileTest.Run()
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			networks.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(env, networks, utils.ProjectRoot, chainlinkNodes, &profileTest.TestReporter)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
