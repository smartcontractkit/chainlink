package integration

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("Keeper suite @keeper", func() {
	var (
		suiteSetup       actions.SuiteSetup
		networkInfo      actions.NetworkInfo
		nodes            []client.Chainlink
		nodeAddresses    []common.Address
		nodeAddressesStr = make([]string, 0)
		consumer         contracts.KeeperConsumer
		registry         contracts.KeeperRegistry
		checkGasLimit    = uint32(2500000)
		err              error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				// need to register at least 5 nodes to perform upkeep
				environment.NewChainlinkCluster(5),
				client.DefaultNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())
			networkInfo = suiteSetup.DefaultNetwork()

			networkInfo.Client.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			ethAmount, err := networkInfo.Deployer.CalculateETHForTXs(networkInfo.Wallets.Default(), networkInfo.Network.Config(), 10)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(
				nodes,
				networkInfo.Client,
				networkInfo.Wallets.Default(),
				ethAmount,
				big.NewFloat(1),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying Keeper contracts", func() {
			ef, err := networkInfo.Deployer.DeployMockETHLINKFeed(networkInfo.Wallets.Default(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := networkInfo.Deployer.DeployMockGasFeed(networkInfo.Wallets.Default(), big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			registry, err = networkInfo.Deployer.DeployKeeperRegistry(
				networkInfo.Wallets.Default(),
				&contracts.KeeperRegistryOpts{
					LinkAddr:             networkInfo.Link.Address(),
					ETHFeedAddr:          ef.Address(),
					GasFeedAddr:          gf.Address(),
					PaymentPremiumPPB:    uint32(200000000),
					BlockCountPerTurn:    big.NewInt(3),
					CheckGasLimit:        checkGasLimit,
					StalenessSeconds:     big.NewInt(90000),
					GasCeilingMultiplier: uint16(1),
					FallbackGasPrice:     big.NewInt(2e11),
					FallbackLinkPrice:    big.NewInt(2e18),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = registry.Fund(networkInfo.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = networkInfo.Deployer.DeployKeeperConsumer(networkInfo.Wallets.Default(), big.NewInt(5))
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.Fund(networkInfo.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Registering upkeep target", func() {
			registrar, err := networkInfo.Deployer.DeployUpkeepRegistrationRequests(
				networkInfo.Wallets.Default(),
				networkInfo.Link.Address(),
				big.NewInt(0),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = registry.SetRegistrar(networkInfo.Wallets.Default(), registrar.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = registrar.SetRegistrarConfig(
				networkInfo.Wallets.Default(),
				true,
				uint32(999),
				uint16(999),
				registry.Address(),
				big.NewInt(0),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req, err := registrar.EncodeRegisterRequest(
				"upkeep_1",
				[]byte("0x1234"),
				consumer.Address(),
				checkGasLimit,
				networkInfo.Wallets.Default().Address(),
				[]byte("0x"),
				big.NewInt(9e18),
				0,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Link.TransferAndCall(networkInfo.Wallets.Default(), registrar.Address(), big.NewInt(9e18), req)
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Adding Keepers and a job", func() {
			keys, err := nodes[0].ReadETHKeys()
			Expect(err).ShouldNot(HaveOccurred())
			na := keys.Data[0].Attributes.Address
			for _, cla := range nodeAddresses {
				nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
			}
			payees := []string{
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
			}
			err = registry.SetKeepers(networkInfo.Wallets.Default(), nodeAddressesStr, payees)
			Expect(err).ShouldNot(HaveOccurred())
			_, err = nodes[0].CreateJob(&client.KeeperJobSpec{
				Name:            "keeper",
				ContractAddress: registry.Address(),
				FromAddress:     na,
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with Keeper job", func() {
		It("performs upkeep of a target contract", func() {
			Eventually(func(g Gomega) {
				cnt, err := consumer.Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(cnt.Int64()).Should(BeNumerically(">", 0))
				log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
			}, "2m", "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			networkInfo.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
