//go:build smoke
package smoke

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

var _ = Describe("Keeper suite @keeper", func() {
	var (
		err           error
		nets          *client.Networks
		cd            contracts.ContractDeployer
		registry      contracts.KeeperRegistry
		consumer      contracts.KeeperConsumer
		checkGasLimit = uint32(2500000)
		lt            contracts.LinkToken
		cls           []client.Chainlink
		nodeAddresses []common.Address
		e             *environment.Environment
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(6, nil)),
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
			nodeAddresses, err = actions.ChainlinkNodeAddresses(cls)
			Expect(err).ShouldNot(HaveOccurred())
			nets.Default.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := nets.Default.EstimateCostForChainlinkOperations(10)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(cls, nets.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying Keeper contracts", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			ef, err := cd.DeployMockETHLINKFeed(big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := cd.DeployMockGasFeed(big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			registry, err = cd.DeployKeeperRegistry(
				&contracts.KeeperRegistryOpts{
					LinkAddr:             lt.Address(),
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
			err = lt.Transfer(registry.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = cd.DeployKeeperConsumer(big.NewInt(5))
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.Transfer(consumer.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Registering upkeep target", func() {
			registrar, err := cd.DeployUpkeepRegistrationRequests(
				lt.Address(),
				big.NewInt(0),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = registry.SetRegistrar(registrar.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = registrar.SetRegistrarConfig(
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
				consumer.Address(),
				[]byte("0x"),
				big.NewInt(9e18),
				0,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.TransferAndCall(registrar.Address(), big.NewInt(9e18), req)
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Adding Keepers and a job", func() {
			keys, err := cls[0].ReadETHKeys()
			Expect(err).ShouldNot(HaveOccurred())
			na := keys.Data[0].Attributes.Address
			nodeAddressesStr := make([]string, 0)
			for _, cla := range nodeAddresses {
				nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
			}
			payees := []string{
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
			}
			err = registry.SetKeepers(nodeAddressesStr, payees)
			Expect(err).ShouldNot(HaveOccurred())
			_, err = cls[0].CreateJob(&client.KeeperJobSpec{
				Name:            "keeper",
				ContractAddress: registry.Address(),
				FromAddress:     na,
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
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
			nets.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nets, "../")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
