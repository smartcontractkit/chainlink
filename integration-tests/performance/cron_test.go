package performance

//revive:disable:dot-imports
import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

var _ = Describe("Cronjob suite @cron", func() {
	var (
		err           error
		job           *client.Job
		chainlinkNode client.Chainlink
		mockserver    *client.MockserverClient
		e             *environment.Environment
		profileTest   *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			// Increase HTTP_SERVER_WRITE_TIMEOUT to be larger than profile duration.
			config.ProjectConfig.FrameworkConfig.ChainlinkEnvValues["HTTP_SERVER_WRITE_TIMEOUT"] = "300s"

			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(
					config.ChainlinkVals(),
					"chainlink-cron-profiling",
					config.GethNetworks()...,
				),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
		})

		By("Connecting to launched resources", func() {
			cls, err := client.ConnectChainlinkNodes(e)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			mockserver, err = client.ConnectMockServer(e)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver client shouldn't fail")
			chainlinkNode = cls[0]
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				// initial value set is performed before jobs creation
				Eventually(func(g Gomega) {
					err = mockserver.SetValuePath("/variable", 5)
					Expect(err).ShouldNot(HaveOccurred(), "Setting value path in mockserver shouldn't fail")

					bta := client.BridgeTypeAttributes{
						Name:        fmt.Sprintf("variable-%s", uuid.NewV4().String()),
						URL:         fmt.Sprintf("%s/variable", mockserver.Config.ClusterURL),
						RequestData: "{}",
					}
					err = chainlinkNode.CreateBridge(&bta)
					Expect(err).ShouldNot(HaveOccurred(), "Creating bridge in chainlink node shouldn't fail")

					job, err = chainlinkNode.CreateJob(&client.CronJobSpec{
						Schedule:          "CRON_TZ=UTC * * * * * *",
						ObservationSource: client.ObservationSourceSpecBridge(bta),
					})
					Expect(err).ShouldNot(HaveOccurred(), "Creating Cron Job in chainlink node shouldn't fail")

					jobRuns, err := chainlinkNode.ReadRunsByJob(job.Data.ID)
					g.Expect(err).ShouldNot(HaveOccurred(), "Reading Job run data shouldn't fail")
					g.Expect(len(jobRuns.Data)).Should(BeNumerically(">=", 5), "Expected number of job runs to be greater than 5, but got %d", len(jobRuns.Data))

					for _, jr := range jobRuns.Data {
						g.Expect(jr.Attributes.Errors).Should(Equal([]interface{}{nil}), "Job run %s shouldn't have errors", jr.ID)
					}
				}, "2m", "1s").Should(Succeed())
			}

			profileTest = testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
				ProfileFunction: profileFunction,
				ProfileDuration: 30 * time.Second,
				ChainlinkNodes:  []client.Chainlink{chainlinkNode},
			})
			profileTest.Setup(e)
		})
	})

	Describe("with Cron job", func() {
		It("runs 5 or more times with no errors", func() {
			profileTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			networkRegistry := blockchain.NewDefaultNetworkRegistry()
			networks, err := networkRegistry.GetNetworks(e)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			err = actions.TeardownSuite(e, networks, utils.ProjectRoot, []client.Chainlink{chainlinkNode}, &profileTest.TestReporter)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
