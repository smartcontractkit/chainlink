package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

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

var _ = Describe("OCR forwarder flow - each operator forwarder pair belongs to each node @ocr-forwarder", func() {
	var (
		testScenarios = []TableEntry{
			Entry("OCR with operator forwarder suite @default", forwarderOCREnv()),
		}

		err               error
		testEnvironment   *environment.Environment
		chainClient       blockchain.EVMClient
		contractDeployer  contracts.ContractDeployer
		contractLoader    contracts.ContractLoader
		linkTokenContract contracts.LinkToken
		chainlinkNodes    []*client.Chainlink
		nodeAddresses     []common.Address
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
		testInputs *smokeTestInputs,
	) {
		By("Deploying the environment")
		testEnvironment = testInputs.environment
		testNetwork := testInputs.network
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(testNetwork, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		contractDeployer, err = contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
		contractLoader, err = contracts.NewContractLoader(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Loading contracts shouldn't fail")
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		nodeAddresses, err = actions.ChainlinkNodeAddresses(chainlinkNodes)
		Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
		mockServer, err = ctfClient.ConnectMockServer(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")

		chainClient.ParallelTransactions(true)
		Expect(err).ShouldNot(HaveOccurred())

		linkTokenContract, err = contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		By("Funding Chainlink nodes")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.05))
		Expect(err).ShouldNot(HaveOccurred())

		By("Prepare forwarder contracts onchain")
		operators, authorizedForwarders, _ := actions.DeployForwarderContracts(contractDeployer, linkTokenContract, chainClient, len(chainlinkNodes[1:]))
		forwarderNodes := chainlinkNodes[1:]
		forwarderNodesAddresses := nodeAddresses[1:]
		for i := range forwarderNodes {
			actions.AcceptAuthorizedReceiversOperator(operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, chainClient, contractLoader)
			Expect(err).ShouldNot(HaveOccurred(), "Accepting Authorize Receivers on Operator shouldn't fail")
			By("Add forwarder track into DB")
			actions.TrackForwarder(chainClient, authorizedForwarders[i], forwarderNodes[i])
			err = chainClient.WaitForEvents()
		}
		By("Deploying OCR contracts")
		ocrInstances = actions.DeployOCRContractsForwarderFlow(1, linkTokenContract, contractDeployer, chainlinkNodes, authorizedForwarders, chainClient)
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

		By("Setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, mockServer))
		By("Creating OCR jobs", actions.CreateOCRJobsWithForwarder(ocrInstances, chainlinkNodes, mockServer))

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

func forwarderOCREnv() *smokeTestInputs {
	network := networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !network.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	envValueMap := network.ChainlinkValuesMap()
	envValueMap["ETH_USE_FORWARDERS"] = "true"
	env := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr-forwarder-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"env":      envValueMap,
			"replicas": 6,
		}))
	return &smokeTestInputs{
		environment: env,
		network:     network,
	}
}
