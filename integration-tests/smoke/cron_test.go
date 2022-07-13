package smoke

//revive:disable:dot-imports
import (
	"fmt"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	networks "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Cronjob suite @cron", func() {
	var (
		testScenarios = []TableEntry{
			Entry("Cronjob suite on Simulated Network @simulated",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, nil),
			),
			Entry("Cronjob suite on Metis Stardust @metis",
				blockchain.NewMetisMultiNodeClientSetup(networks.MetisStardust),
				ethereum.New(&ethereum.Props{
					NetworkName: networks.MetisStardust.Name,
					Simulated:   networks.MetisStardust.Simulated,
					WsURLs:      networks.MetisStardust.URLs,
				}),
				chainlink.New(0, map[string]interface{}{
					"env": networks.MetisStardust.ChainlinkValuesMap(),
				}),
			),
			Entry("Cronjob suite on Sepolia Testnet @sepolia",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SepoliaTestnet),
				ethereum.New(&ethereum.Props{
					NetworkName: networks.SepoliaTestnet.Name,
					Simulated:   networks.SepoliaTestnet.Simulated,
					WsURLs:      networks.SepoliaTestnet.URLs,
				}),
				chainlink.New(0, map[string]interface{}{
					"env": networks.SepoliaTestnet.ChainlinkValuesMap(),
				}),
			),
		}

		err             error
		job             *client.Job
		chainlinkNode   client.Chainlink
		mockServer      *client.MockserverClient
		testEnvironment *environment.Environment
	)

	AfterEach(func() {
		By("Tearing down the environment")
		err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, []client.Chainlink{chainlinkNode}, nil, nil)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("Cronjob suite on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		evmChart environment.ConnectedChart,
		chainlinkCharts ...environment.ConnectedChart,
	) {
		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-cron"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(evmChart)
		for _, chainlinkChart := range chainlinkCharts {
			testEnvironment.AddHelm(chainlinkChart)
		}
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred(), "Error deploying test environment")

		By("Connecting to launched resources")
		cls, err := client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		mockServer, err = client.ConnectMockServer(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver client shouldn't fail")
		chainlinkNode = cls[0]

		By("Adding cron job to a node")
		err = mockServer.SetValuePath("/variable", 5)
		Expect(err).ShouldNot(HaveOccurred(), "Setting value path in mockserver shouldn't fail")

		bta := client.BridgeTypeAttributes{
			Name:        fmt.Sprintf("variable-%s", uuid.NewV4().String()),
			URL:         fmt.Sprintf("%s/variable", mockServer.Config.ClusterURL),
			RequestData: "{}",
		}
		err = chainlinkNode.CreateBridge(&bta)
		Expect(err).ShouldNot(HaveOccurred(), "Creating bridge in chainlink node shouldn't fail")

		job, err = chainlinkNode.CreateJob(&client.CronJobSpec{
			Schedule:          "CRON_TZ=UTC * * * * * *",
			ObservationSource: client.ObservationSourceSpecBridge(bta),
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating Cron Job in chainlink node shouldn't fail")

		Eventually(func(g Gomega) {
			jobRuns, err := chainlinkNode.ReadRunsByJob(job.Data.ID)
			g.Expect(err).ShouldNot(HaveOccurred(), "Reading Job run data shouldn't fail")

			g.Expect(len(jobRuns.Data)).Should(BeNumerically(">=", 5), "Expected number of job runs to be greater than 5, but got %d", len(jobRuns.Data))

			for _, jr := range jobRuns.Data {
				g.Expect(jr.Attributes.Errors).Should(Equal([]interface{}{nil}), "Job run %s shouldn't have errors", jr.ID)
			}
		}, "2m", "1s").Should(Succeed())
	},
		testScenarios,
	)
})
