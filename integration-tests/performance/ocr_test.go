package performance

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	it "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

var _ = Describe("OCR Feed @ocr", func() {
	var (
		err               error
		env               *environment.Environment
		c                 blockchain.EVMClient
		contractDeployer  contracts.ContractDeployer
		linkTokenContract contracts.LinkToken
		chainlinkNodes    []client.Chainlink
		ms                *client.MockserverClient
		ocrInstances      []contracts.OffchainAggregator
		profileTest       *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env = environment.New(nil).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, map[string]interface{}{
					"replicas": "6",
					"env": map[string]interface{}{
						"HTTP_SERVER_WRITE_TIMEOUT": "300s",
					},
				}))
			err = env.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			// Load Networks
			var err error
			c, err = blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings)(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			contractDeployer, err = contracts.NewContractDeployer(c)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

			chainlinkNodes, err = client.ConnectChainlinkNodes(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			ms, err = client.ConnectMockServer(env)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")

			c.ParallelTransactions(true)
			Expect(err).ShouldNot(HaveOccurred())

			linkTokenContract, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
		})

		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(chainlinkNodes, c, big.NewFloat(.01))
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying OCR contracts", func() {
			ocrInstances = actions.DeployOCRContracts(1, linkTokenContract, contractDeployer, chainlinkNodes, c)
			err = c.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
					// Not the last node, hence not all nodes started profiling yet.
					return
				}
				By("Setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, ms))
				By("Creating OCR jobs", actions.CreateOCRJobs(ocrInstances, chainlinkNodes, ms))
				By("Starting new round", actions.StartNewRound(1, ocrInstances, c))
				answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
				Expect(err).ShouldNot(HaveOccurred(), "Getting latest answer from OCR contract shouldn't fail")
				Expect(answer.Int64()).Should(Equal(int64(5)), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

				By("setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, chainlinkNodes, ms))
				By("starting new round", actions.StartNewRound(2, ocrInstances, c))

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
			c.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(env, utils.ProjectRoot, chainlinkNodes, &profileTest.TestReporter, c)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
