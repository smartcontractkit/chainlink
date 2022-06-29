package smoke

//revive:disable:dot-imports
import (
	"fmt"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	networks "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

var _ = Describe("Cronjob suite @cron", func() {
	var (
		testScenarios = []TableEntry{
			Entry("Cronjob suite on Simulated Network @simulated", blockchain.NewEthereumMultiNodeClientSetup, networks.SimulatedEVMNetwork, nil),
			Entry("Cronjob suite on Metis Stardust @metis", blockchain.NewMetisMultiNodeClientSetup, networks.MetisTestNetwork, map[string]interface{}{
				"env": map[string]interface{}{
					"eth_url":      networks.MetisTestNetwork.URLs[0],
					"eth_chain_id": fmt.Sprint(networks.MetisTestNetwork.ChainID),
				},
			}),
		}

		err             error
		job             *client.Job
		chainlinkNode   client.Chainlink
		mockServer      *client.MockserverClient
		testEnvironment *environment.Environment
	)

	AfterEach(func() {
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, []client.Chainlink{chainlinkNode}, nil, nil)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})

	DescribeTable("Cronjob suite on different EVM networks", func(
		clientFunc func(networkSettings *blockchain.EVMNetwork) func(*environment.Environment) (blockchain.EVMClient, error),
		evmNetwork *blockchain.EVMNetwork,
		chainlinkValues map[string]interface{},
	) {
		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-cron"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(ethereum.New(nil)).
			AddHelm(chainlink.New(0, chainlinkValues))
		err = testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

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
