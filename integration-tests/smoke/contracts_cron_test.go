//go:build integration

package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Cronjob suite @cron", func() {
	var (
		suiteSetup actions.SuiteSetup
		adapter    environment.ExternalAdapter
		nodes      []client.Chainlink
		job        *client.Job
		err        error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkCluster(1),
				actions.EVMNetworkFromConfigHook,
				actions.EthereumDeployerHook,
				actions.EthereumClientHook,
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Adding cron job to a node", func() {
			bta := client.BridgeTypeAttributes{
				Name:        "five",
				URL:         fmt.Sprintf("%s/five", adapter.ClusterURL()),
				RequestData: "{}",
			}
			err = nodes[0].CreateBridge(&bta)
			Expect(err).ShouldNot(HaveOccurred())

			job, err = nodes[0].CreateJob(&client.CronJobSpec{
				Schedule:          "CRON_TZ=UTC * * * * * *",
				ObservationSource: client.ObservationSourceSpecBridge(bta),
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with Cron job", func() {
		It("runs 5 times with no errors", func() {
			Eventually(func(g Gomega) {
				jobRuns, err := nodes[0].ReadRunsByJob(job.Data.ID)
				g.Expect(err).ShouldNot(HaveOccurred())

				g.Expect(len(jobRuns.Data)).Should(BeNumerically("==", 5))

				for _, jr := range jobRuns.Data {
					g.Expect(jr.Attributes.Errors).Should(Equal([]interface{}{nil}))
				}
			}, "2m", "1s").Should(Succeed())

		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
