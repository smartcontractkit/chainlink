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

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Direct request suite @runlog", func() {
	var (
		testScenarios = []TableEntry{
			Entry("Runlog suite on Simulated Network @simulated",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, nil),
			),
			Entry("Runlog suite on Metis Stardust @metis",
				blockchain.NewMetisMultiNodeClientSetup(networks.MetisStardust),
				ethereum.New(&ethereum.Props{
					NetworkName: networks.MetisStardust.Name,
					Simulated:   networks.MetisStardust.Simulated,
				}),
				chainlink.New(0, map[string]interface{}{
					"env": networks.MetisStardust.ChainlinkValuesMap(),
				}),
			),
		}

		err              error
		chainClient      blockchain.EVMClient
		contractDeployer contracts.ContractDeployer
		chainlinkNodes   []client.Chainlink
		oracle           contracts.Oracle
		consumer         contracts.APIConsumer
		jobUUID          uuid.UUID
		mockServer       *client.MockserverClient
		testEnvironment  *environment.Environment
	)

	AfterEach(func() {
		By("Tearing down the environment")
		chainClient.GasStats().PrintStats()
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("Direct request suite on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		evmChart environment.ConnectedChart,
		chainlinkCharts ...environment.ConnectedChart,
	) {

		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-runlog"}).
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
		mockServer, err = client.ConnectMockServer(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred())

		By("Funding Chainlink nodes")
		ethAmount, err := chainClient.EstimateCostForChainlinkOperations(1)
		Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, ethAmount)
		Expect(err).ShouldNot(HaveOccurred(), "Funding chainlink nodes with ETH shouldn't fail")

		By("Deploying contracts")
		lt, err := contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
		oracle, err = contractDeployer.DeployOracle(lt.Address())
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Oracle Contract shouldn't fail")
		consumer, err = contractDeployer.DeployAPIConsumer(lt.Address())
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Consumer Contract shouldn't fail")
		err = chainClient.SetDefaultWallet(0)
		Expect(err).ShouldNot(HaveOccurred(), "Setting default wallet shouldn't fail")
		err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
		Expect(err).ShouldNot(HaveOccurred(), "Transferring %d to consumer contract shouldn't fail", big.NewInt(2e18))

		By("Creating directrequest job")
		err = mockServer.SetValuePath("/variable", 5)
		Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")

		jobUUID = uuid.NewV4()

		bta := client.BridgeTypeAttributes{
			Name: fmt.Sprintf("five-%s", jobUUID.String()),
			URL:  fmt.Sprintf("%s/variable", mockServer.Config.ClusterURL),
		}
		err = chainlinkNodes[0].CreateBridge(&bta)
		Expect(err).ShouldNot(HaveOccurred(), "Creating bridge shouldn't fail")

		os := &client.DirectRequestTxPipelineSpec{
			BridgeTypeAttributes: bta,
			DataPath:             "data,result",
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred(), "Building observation source spec shouldn't fail")

		_, err = chainlinkNodes[0].CreateJob(&client.DirectRequestJobSpec{
			Name:                     "direct_request",
			MinIncomingConfirmations: "1",
			ContractAddress:          oracle.Address(),
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating direct_request job shouldn't fail")

		By("Calling oracle contract")
		jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
		var jobID [32]byte
		copy(jobID[:], jobUUIDReplaces)
		err = consumer.CreateRequestTo(
			oracle.Address(),
			jobID,
			big.NewInt(1e18),
			fmt.Sprintf("%s/variable", mockServer.Config.ClusterURL),
			"data,result",
			big.NewInt(100),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Calling oracle contract shouldn't fail")

		By("receives API call data on-chain")
		Eventually(func(g Gomega) {
			d, err := consumer.Data(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred(), "Getting data from consumer contract shouldn't fail")
			g.Expect(d).ShouldNot(BeNil(), "Expected the initial on chain data to be nil")
			log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
			g.Expect(d.Int64()).Should(BeNumerically("==", 5), "Expected the on-chain data to be 5, but found %d", d.Int64())
		}, "2m", "1s").Should(Succeed())
	},
		testScenarios,
	)
})
