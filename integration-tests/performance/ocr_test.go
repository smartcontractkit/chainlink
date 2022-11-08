package performance

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OCR Feed @ocr", func() {
	var (
		err               error
		testEnvironment   *environment.Environment
		chainClient       blockchain.EVMClient
		contractDeployer  contracts.ContractDeployer
		linkTokenContract contracts.LinkToken
		chainlinkNodes    []*client.Chainlink
		mockServer        *ctfClient.MockserverClient
		ocrInstances      []contracts.OffchainAggregator
		profileTest       *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			testEnvironment = environment.New(&environment.Config{NamespacePrefix: "performance-ocr"}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, map[string]interface{}{
					"replicas": "6",
					"env": map[string]interface{}{
						"HTTP_SERVER_WRITE_TIMEOUT": "300s",
					},
				}))
			err = testEnvironment.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			// Load Networks
			var err error
			chainClient, err = blockchain.NewEVMClient(blockchain.SimulatedEVMNetwork, testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			contractDeployer, err = contracts.NewContractDeployer(chainClient)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

			chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			mockServer, err = ctfClient.ConnectMockServer(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")

			chainClient.ParallelTransactions(true)
			Expect(err).ShouldNot(HaveOccurred())

			linkTokenContract, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
		})

		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.01))
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying OCR contracts", func() {
			ocrInstances = actions.DeployOCRContracts(1, linkTokenContract, contractDeployer, chainlinkNodes, chainClient)
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode *client.Chainlink) {
				defer GinkgoRecover()
				if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
					// Not the last node, hence not all nodes started profiling yet.
					return
				}
				By("Setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, mockServer))
				By("Creating OCR jobs", actions.CreateOCRJobs(ocrInstances, chainlinkNodes, mockServer))
				By("Starting new round", actions.StartNewRound(1, ocrInstances, chainClient))
				answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Getting latest answer from OCR contract shouldn't fail")
				Expect(answer.Int64()).Should(Equal(int64(5)), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

				By("setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, chainlinkNodes, mockServer))
				By("starting new round", actions.StartNewRound(2, ocrInstances, chainClient))

				answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(answer.Int64()).Should(Equal(int64(10)), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
			}

			profileTest = testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
				ProfileFunction: profileFunction,
				ProfileDuration: time.Minute,
				ChainlinkNodes:  chainlinkNodes,
			})
			profileTest.Setup(testEnvironment)
		})
	})

	Describe("With a single OCR contract", func() {
		It("performs two rounds", func() {
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
