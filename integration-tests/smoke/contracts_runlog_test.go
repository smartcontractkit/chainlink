//go:build smoke

package integration

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Direct request suite @runlog", func() {
	var (
		suiteSetup    actions.SuiteSetup
		networkInfo   actions.NetworkInfo
		adapter       environment.ExternalAdapter
		nodes         []client.Chainlink
		nodeAddresses []common.Address
		oracle        contracts.Oracle
		consumer      contracts.APIConsumer
		jobUUID       uuid.UUID
		err           error
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
			networkInfo = suiteSetup.DefaultNetwork()
			adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Funding Chainlink nodes", func() {
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())
			ethAmount, err := networkInfo.Deployer.CalculateETHForTXs(networkInfo.Wallets.Default(), networkInfo.Network.Config(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(nodes, networkInfo.Client, networkInfo.Wallets.Default(), ethAmount, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying and funding the contracts", func() {
			oracle, err = networkInfo.Deployer.DeployOracle(networkInfo.Wallets.Default(), networkInfo.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = networkInfo.Deployer.DeployAPIConsumer(networkInfo.Wallets.Default(), networkInfo.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.Fund(networkInfo.Wallets.Default(), nil, big.NewFloat(2))
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Permitting node to fulfill request", func() {
			err = oracle.SetFulfillmentPermission(networkInfo.Wallets.Default(), nodeAddresses[0].Hex(), true)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating directrequest job", func() {
			jobUUID = uuid.NewV4()

			bta := client.BridgeTypeAttributes{
				Name: "five",
				URL:  fmt.Sprintf("%s/five", adapter.ClusterURL()),
			}
			err = nodes[0].CreateBridge(&bta)
			Expect(err).ShouldNot(HaveOccurred())

			os := &client.DirectRequestTxPipelineSpec{
				BridgeTypeAttributes: bta,
				DataPath:             "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred())

			_, err = nodes[0].CreateJob(&client.DirectRequestJobSpec{
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
				networkInfo.Wallets.Default(),
				oracle.Address(),
				jobID,
				big.NewInt(1e18),
				fmt.Sprintf("%s/five", adapter.ClusterURL()),
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
		By("Calculating gas costs", func() {
			networkInfo.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
