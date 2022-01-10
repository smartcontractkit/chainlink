//go:build smoke

package smoke_test

//revive:disable:dot-imports
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
	"github.com/smartcontractkit/integrations-framework/utils"
)

var _ = Describe("Keeper suite @keeper", func() {
	var (
		err              error
		networks         *client.Networks
		contractDeployer contracts.ContractDeployer
		registry         contracts.KeeperRegistry
		consumer         contracts.KeeperConsumer
		checkGasLimit    = uint32(2500000)
		linkToken        contracts.LinkToken
		chainlinkNodes   []client.Chainlink
		nodeAddresses    []common.Address
		env              *environment.Environment
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(6, nil)),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = env.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			networkRegistry := client.NewNetworkRegistry()
			networks, err = networkRegistry.GetNetworks(env)
			Expect(err).ShouldNot(HaveOccurred())
			contractDeployer, err = contracts.NewContractDeployer(networks.Default)
			Expect(err).ShouldNot(HaveOccurred())
			chainlinkNodes, err = client.ConnectChainlinkNodes(env)
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(chainlinkNodes)
			Expect(err).ShouldNot(HaveOccurred())
			networks.Default.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := networks.Default.EstimateCostForChainlinkOperations(10)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(chainlinkNodes, networks.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred())
			// Edge case where simulated networks need some funds at the 0x0 address in order for keeper reads to work
			if networks.Default.GetNetworkType() == "eth_simulated" {
				actions.FundAddresses(networks.Default, big.NewFloat(1), "0x0")
			}
		})

		By("Deploying Keeper contracts", func() {
			linkToken, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			registry, err = contractDeployer.DeployKeeperRegistry(
				&contracts.KeeperRegistryOpts{
					LinkAddr:             linkToken.Address(),
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
			err = linkToken.Transfer(registry.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = contractDeployer.DeployKeeperConsumer(big.NewInt(5))
			Expect(err).ShouldNot(HaveOccurred())
			err = linkToken.Transfer(consumer.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred())
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Registering upkeep target", func() {
			registrar, err := contractDeployer.DeployUpkeepRegistrationRequests(
				linkToken.Address(),
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
			err = linkToken.TransferAndCall(registrar.Address(), big.NewInt(9e18), req)
			Expect(err).ShouldNot(HaveOccurred())
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Adding Keepers and a job", func() {
			primaryNode := chainlinkNodes[0]
			primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred())
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
			_, err = primaryNode.CreateJob(&client.KeeperJobSpec{
				Name:                     "keeper-test-job",
				ContractAddress:          registry.Address(),
				FromAddress:              primaryNodeAddress,
				MinIncomingConfirmations: 1,
				ObservationSource:        client.ObservationSourceKeeperDefault(),
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with Keeper job", func() {
		It("performs upkeep of a target contract", func() {
			Eventually(func(g Gomega) {
				cnt, err := consumer.Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)))
				log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
			}, "2m", "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			networks.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(env, networks, utils.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
