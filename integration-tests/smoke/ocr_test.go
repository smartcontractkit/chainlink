package smoke

//revive:disable:dot-imports
import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OCR Feed @ocr", func() {
	var (
		testScenarios = []TableEntry{
			Entry("OCR suite on Simulated Network @simulated",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, map[string]interface{}{
					"replicas": 6,
				}),
			),
			Entry("OCR suite on Metis Stardust @metis",
				blockchain.NewMetisMultiNodeClientSetup(networks.MetisStardust),
				ethereum.New(&ethereum.Props{
					NetworkName: networks.MetisStardust.Name,
					Simulated:   networks.MetisStardust.Simulated,
				}),
				chainlink.New(0, map[string]interface{}{
					"env":      networks.MetisStardust.ChainlinkValuesMap(),
					"replicas": 6,
				}),
			),
		}

		err               error
		testEnvironment   *environment.Environment
		chainClient       blockchain.EVMClient
		contractDeployer  contracts.ContractDeployer
		linkTokenContract contracts.LinkToken
		chainlinkNodes    []client.Chainlink
		mockServer        *ctfClient.MockserverClient
		ocrInstances      []contracts.OffchainAggregator
	)

	AfterEach(func() {
		By("Tearing down the environment")
		chainClient.GasStats().PrintStats()
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("OCR suite on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		evmChart environment.ConnectedChart,
		chainlinkCharts ...environment.ConnectedChart,
	) {
		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-ocr"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(evmChart)
		for _, chainlinkChart := range chainlinkCharts {
			testEnvironment.AddHelm(chainlinkChart)
		}
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = clientFunc(testEnvironment)
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

		By("Funding Chainlink nodes")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.01))
		Expect(err).ShouldNot(HaveOccurred())

		By("Deploying OCR contracts")
		ocrInstances = actions.DeployOCRContracts(1, linkTokenContract, contractDeployer, chainlinkNodes, chainClient)
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

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
	},
		testScenarios,
	)
})
