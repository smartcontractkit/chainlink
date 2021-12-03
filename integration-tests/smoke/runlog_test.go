//go:build smoke
package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/utils"
)

var _ = Describe("Direct request suite @runlog", func() {
	var (
		err        error
		nets       *client.Networks
		cd         contracts.ContractDeployer
		cls        []client.Chainlink
		oracle     contracts.Oracle
		consumer   contracts.APIConsumer
		jobUUID    uuid.UUID
		mockserver *client.MockserverClient
		e          *environment.Environment
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(3, nil)),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			networkRegistry := client.NewNetworkRegistry()
			nets, err = networkRegistry.GetNetworks(e)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.NewChainlinkClients(e)
			Expect(err).ShouldNot(HaveOccurred())
			mockserver, err = client.NewMockServerClientFromEnv(e)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Funding Chainlink nodes", func() {
			ethAmount, err := nets.Default.EstimateCostForChainlinkOperations(1)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(cls, nets.Default, ethAmount)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying contracts", func() {
			lt, err := cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			oracle, err = cd.DeployOracle(lt.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = cd.DeployAPIConsumer(lt.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.SetWallet(0)
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Creating directrequest job", func() {
			err = mockserver.SetValuePath("/variable", 5)
			Expect(err).ShouldNot(HaveOccurred())

			jobUUID = uuid.NewV4()

			bta := client.BridgeTypeAttributes{
				Name: fmt.Sprintf("five-%s", jobUUID.String()),
				URL:  fmt.Sprintf("%s/variable", mockserver.Config.ClusterURL),
			}
			err = cls[0].CreateBridge(&bta)
			Expect(err).ShouldNot(HaveOccurred())

			os := &client.DirectRequestTxPipelineSpec{
				BridgeTypeAttributes: bta,
				DataPath:             "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred())

			_, err = cls[0].CreateJob(&client.DirectRequestJobSpec{
				Name:              "direct_request",
				ContractAddress:   oracle.Address(),
				ExternalJobID:     jobUUID.String(),
				ObservationSource: ost,
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Calling oracle contract", func() {
			jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
			var jobID [32]byte
			copy(jobID[:], jobUUIDReplaces)
			err = consumer.CreateRequestTo(
				oracle.Address(),
				jobID,
				big.NewInt(1e18),
				fmt.Sprintf("%s/variable", mockserver.Config.ClusterURL),
				"data,result",
				big.NewInt(100),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Describe("with DirectRequest job", func() {
		It("receives API call data on-chain", func() {
			Eventually(func(g Gomega) {
				d, err := consumer.Data(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(d).ShouldNot(BeNil())
				log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
				g.Expect(d.Int64()).Should(BeNumerically("==", 5))
			}, "2m", "1s").Should(Succeed())
		})
	})
	AfterEach(func() {
		By("Tearing down the environment", func() {
			nets.Default.GasStats().PrintStats()
			err = actions.TeardownSuite(e, nets, utils.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
