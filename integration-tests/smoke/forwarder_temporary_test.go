package smoke

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
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
)

var _ = Describe("Forwarder suite", func() {
	var (
		testScenarios = []TableEntry{
			Entry("Operator Forwarder suite on Simulated Network @forwarder-smoke-simulated",
				networks.SimulatedEVM,
				big.NewFloat(10),
				environment.New(&environment.Config{}).
					AddHelm(mockservercfg.New(nil)).
					AddHelm(mockserver.New(nil)).
					AddHelm(ethereum.New(nil)).
					AddHelm(chainlink.New(0, map[string]interface{}{
						"env": map[string]interface{}{
							"ETH_USE_FORWARDERS": "true",
						},
						"replicas": 1,
					})),
			),
			Entry("Flux monitor suite on Görli Testnet @forwarder-smoke-goerli",
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
						"env": map[string]interface{}{
							"ETH_USE_FORWARDERS": "true",
							"ETH_CHAIN_ID":       os.Getenv("ETH_CHAIN_ID"),
							"ETH_URL":            os.Getenv("ETH_URL"),
						},
						"replicas": 1,
					})),
			),
		}

		err              error
		chainClient      blockchain.EVMClient
		contractDeployer contracts.ContractDeployer
		contractLoader   contracts.ContractLoader
		linkToken        contracts.LinkToken
		chainlinkNodes   []*client.Chainlink
		mockServer       *ctfClient.MockserverClient
		nodeAddresses    []common.Address
		adapterPath      string
		adapterUUID      string
		testEnvironment  *environment.Environment
	)

	AfterEach(func() {
		By("Tearing down the environment")
		chainClient.GasStats().PrintStats()
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("Forwarder smoke test suite", func(
		testNetwork *blockchain.EVMNetwork,
		funding *big.Float,
		env *environment.Environment,
	) {
		By("Deploying the environment")
		testEnvironment = env
		testEnvironment.Cfg.NamespacePrefix = fmt.Sprintf("smoke-forwarder-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"))

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
		Expect(err).ShouldNot(HaveOccurred(), "Creating mock server client shouldn't fail")
		fmt.Println(nodeAddresses)
		chainClient.ParallelTransactions(true)

		By("Setting initial adapter value")
		adapterUUID = uuid.NewV4().String()
		adapterPath = fmt.Sprintf("/variable-%s", adapterUUID)
		err = mockServer.SetValuePath(adapterPath, 1e5)
		Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")

		By("Deploying and funding contract")
		linkToken, err = contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		operator, authorizedForwarder, _ := actions.DeployForwarderContracts(contractDeployer, linkToken, chainClient)
		actions.AcceptAuthorizedReceiversOperator(operator, authorizedForwarder, nodeAddresses, chainClient, contractLoader)

		actions.TrackForwarder(chainClient, authorizedForwarder, chainlinkNodes)
	},
		testScenarios,
	)
})
