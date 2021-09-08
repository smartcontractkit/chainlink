package contracts

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/integration/suite/testcommon"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Alerts suite", func() {
	var (
		suiteSetup     *actions.DefaultSuiteSetup
		adapter        environment.ExternalAdapter
		chainlinkNodes []client.Chainlink
		explorer       *client.ExplorerClient
		err            error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkClusterForAlertsTesting(3),
				client.NewNetworkFromConfig,
				testcommon.ConfigLocation,
			)
			Expect(err).ShouldNot(HaveOccurred())

			explorer, err = environment.GetExplorerClientFromEnv(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(explorer.BaseURL)

			chainlinkNodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(chainlinkNodes[0].URL())

			adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(adapter.ClusterURL())
		})
	})

	Describe("Alerts", func() {
		It("Test 1", func() {

		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
