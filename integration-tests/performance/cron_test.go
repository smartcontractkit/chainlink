package performance

//revive:disable:dot-imports
import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	it "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Cronjob suite @cron", func() {
	var (
		err             error
		job             *client.Job
		chainlinkNode   client.Chainlink
		mockServer      *client.MockserverClient
		testEnvironment *environment.Environment
		profileTest     *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			testEnvironment = environment.New(&environment.Config{NamespacePrefix: "performance-cron"}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, map[string]interface{}{
					"env": map[string]interface{}{
						"HTTP_SERVER_WRITE_TIMEOUT": "300s",
					},
				}))
			err = testEnvironment.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			cls, err := client.ConnectChainlinkNodes(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			mockServer, err = client.ConnectMockServer(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver client shouldn't fail")
			chainlinkNode = cls[0]
		})

		By("Setting up profiling", func() {
			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				// initial value set is performed before jobs creation
				Eventually(func(g Gomega) {
					err = mockServer.SetValuePath("/variable", 5)
					Expect(err).ShouldNot(HaveOccurred(), "Setting value path in mockserver shouldn't fail")

					bta := client.BridgeTypeAttributes{
						Name:        fmt.Sprintf("variable-%s", uuid.NewV4().String()),
						URL:         fmt.Sprintf("%s/variable", mockServer.Config.ClusterURL),
						RequestData: "{}",
					}
					err = chainlinkNode.CreateBridge(&bta)
					Expect(err).ShouldNot(HaveOccurred(), "Creating bridge in chainlink node shouldn't fail")
					// Getting issue here need to fix
					job, err = chainlinkNode.CreateJob(&it.CronJobSpec{
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
			profileTest.Setup(testEnvironment)
		})
	})

	Describe("with Cron job", func() {
		It("runs 5 or more times with no errors", func() {
			profileTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, []client.Chainlink{chainlinkNode}, &profileTest.TestReporter, nil)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
