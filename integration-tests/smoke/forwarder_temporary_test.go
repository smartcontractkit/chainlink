package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/operator_factory"

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

		err                     error
		chainClient             blockchain.EVMClient
		contractDeployer        contracts.ContractDeployer
		contractLoader          contracts.ContractLoader
		linkToken               contracts.LinkToken
		operatorFactoryInstance contracts.OperatorFactory
		operatorInstance        contracts.Operator
		forwarderInstance       contracts.AuthorizedForwarder
		chainlinkNodes          []*client.Chainlink
		mockServer              *ctfClient.MockserverClient
		nodeAddresses           []common.Address
		adapterPath             string
		adapterUUID             string
		testEnvironment         *environment.Environment
	)

	AfterEach(func() {
		By("Tearing down the environment")
		chainClient.GasStats().PrintStats()
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

		By("Deploying OperatorFactory contract")
		operatorFactoryInstance, err = contractDeployer.DeployOperatorFactory(linkToken.Address())
		Expect(err).ShouldNot(HaveOccurred(), "Deploying OperatorFactory Contract shouldn't fail")
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for deployment of flux aggregator contract")

		By("Subscribe to Operator factory Events")
		operatorCreated := make(chan *operator_factory.OperatorFactoryOperatorCreated)
		authorizedForwarderCreated := make(chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated)
		actions.SubscribeOperatorFactoryEvents(authorizedForwarderCreated, operatorCreated, chainClient, operatorFactoryInstance)

		By("Create new operator and forwarder")
		_, err := operatorFactoryInstance.DeployNewOperatorAndForwarder()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying new operator with proposed ownership with forwarder shouldn't fail")
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Waiting for events in nodes shouldn't fail")
		eventDataAuthorizedForwarder, eventDataOperatorCreated := <-authorizedForwarderCreated, <-operatorCreated
		operator, authorizedForwarder := eventDataOperatorCreated.Operator, eventDataAuthorizedForwarder.Forwarder

		operatorInstance, err = contractLoader.LoadOperatorContract(operator)
		Expect(err).ShouldNot(HaveOccurred(), "Loading operator contract shouldn't fail")
		forwarderInstance, err = contractLoader.LoadAuthorizedForwarder(authorizedForwarder)
		Expect(err).ShouldNot(HaveOccurred(), "Loading authorized forwarder contract shouldn't fail")

		By("Accept authorized receivers")
		log.Info().Msg(fmt.Sprintf("Eoa addresses: %+v", nodeAddresses))
		log.Info().Msg(fmt.Sprintf("Forwarder addresses: %+v", authorizedForwarder.Hex()))
		err = operatorInstance.AcceptAuthorizedReceivers([]common.Address{authorizedForwarder}, nodeAddresses)
		Expect(err).ShouldNot(HaveOccurred(), "Accepting authorized receivers shouldn't fail")
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Waiting for events in nodes shouldn't fail")

		By("Verify authorized senders on forwarder")
		senders, err := forwarderInstance.GetAuthorizedSenders(context.Background())
		Expect(err).ShouldNot(HaveOccurred(), "Getting authorized senders shouldn't fail")
		log.Info().Msg(fmt.Sprintf("Authorized senders: %+v", senders))
		var sendersAddrs []string
		for _, o := range nodeAddresses {
			sendersAddrs = append(sendersAddrs, o.Hex())
		}
		Expect(senders).Should(Equal(sendersAddrs), "Eoa addresses should match node addresses")

		By("Verify forwarder Owner")
		owner, err := forwarderInstance.Owner(context.Background())
		Expect(err).ShouldNot(HaveOccurred(), "Getting authorized forwarder owner shouldn't fail")
		Expect(owner).Should(Equal(operator.Hex()), "Forwarder owner should match operator")

	},
		testScenarios,
	)
})
