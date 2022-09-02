package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OCR Feed @ocr", func() {
	var (
		testScenarios = []TableEntry{
			Entry("OCR suite on Simulated Network @simulated",
				networks.SimulatedEVM,
				big.NewFloat(50),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(nil)).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env":      networks.SimulatedEVM.ChainlinkValuesMap(),
						"replicas": 6,
					})),
			),
			Entry("OCR suite on General EVM @general",
				networks.GeneralEVM(),
				big.NewFloat(1),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(&ethereum.Props{
						NetworkName: networks.GeneralEVM().Name,
						Simulated:   networks.GeneralEVM().Simulated,
						WsURLs:      networks.GeneralEVM().URLs,
					})).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env":      networks.GeneralEVM().ChainlinkValuesMap(),
						"replicas": 6,
					})),
			),
			Entry("OCR suite on Metis Stardust @metis",
				networks.MetisStardust,
				big.NewFloat(.01),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(&ethereum.Props{
						NetworkName: networks.MetisStardust.Name,
						Simulated:   networks.MetisStardust.Simulated,
						WsURLs:      networks.MetisStardust.URLs,
					})).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env":      networks.MetisStardust.ChainlinkValuesMap(),
						"replicas": 6,
					})),
			),
			Entry("OCR suite on Sepolia Testnet @sepolia",
				networks.SepoliaTestnet,
				big.NewFloat(.1),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(&ethereum.Props{
						NetworkName: networks.SepoliaTestnet.Name,
						Simulated:   networks.SepoliaTestnet.Simulated,
						WsURLs:      networks.SepoliaTestnet.URLs,
					})).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env":      networks.SepoliaTestnet.ChainlinkValuesMap(),
						"replicas": 6,
					})),
			),
			Entry("OCR suite on Görli Testnet @goerli",
				networks.GoerliTestnet,
				big.NewFloat(.1),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(&ethereum.Props{
						NetworkName: networks.GoerliTestnet.Name,
						Simulated:   networks.GoerliTestnet.Simulated,
						WsURLs:      networks.GoerliTestnet.URLs,
					})).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env":      networks.GoerliTestnet.ChainlinkValuesMap(),
						"replicas": 6,
					})),
			),
			Entry("OCR suite on Klaytn Baobab @klaytn",
				networks.KlaytnBaobab,
				big.NewFloat(1),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(&ethereum.Props{
						NetworkName: networks.KlaytnBaobab.Name,
						Simulated:   networks.KlaytnBaobab.Simulated,
						WsURLs:      networks.KlaytnBaobab.URLs,
					})).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env":      networks.KlaytnBaobab.ChainlinkValuesMap(),
						"replicas": 6,
					})),
			),
		}

		err               error
		testEnvironment   *environment.Environment
		chainClient       blockchain.EVMClient
		contractDeployer  contracts.ContractDeployer
		linkTokenContract contracts.LinkToken
		chainlinkNodes    []*client.Chainlink
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
		testNetwork *blockchain.EVMNetwork,
		funding *big.Float,
		env *environment.Environment,
	) {
		By("Deploying the environment")
		testEnvironment = env
		testEnvironment.Cfg.NamespacePrefix = fmt.Sprintf("smoke-ocr-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"))

		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(testNetwork, testEnvironment)
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
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, funding)
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
