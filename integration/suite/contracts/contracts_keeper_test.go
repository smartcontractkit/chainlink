package contracts

import (
	"context"
	"errors"
	"math/big"

	"github.com/avast/retry-go"
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
		s                *actions.DefaultSuiteSetup
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
			s, err = actions.DefaultLocalSetup(
				// need to register at least 5 nodes to perform upkeep
				environment.NewChainlinkCluster(5),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())

			s.Client.ParallelTransactions(true)
		})
		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(
				nodes,
				s.Client,
				s.Wallets.Default(),
				big.NewFloat(2),
				big.NewFloat(1),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying Keeper contracts", func() {
			ef, err := s.Deployer.DeployMockETHLINKFeed(s.Wallets.Default(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := s.Deployer.DeployMockGasFeed(s.Wallets.Default(), big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			registry, err = s.Deployer.DeployKeeperRegistry(
				s.Wallets.Default(),
				&contracts.KeeperRegistryOpts{
					LinkAddr:             s.Link.Address(),
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
			err = registry.Fund(s.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = s.Deployer.DeployKeeperConsumer(s.Wallets.Default(), big.NewInt(5))
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.Fund(s.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Registering upkeep target", func() {
			registrar, err := s.Deployer.DeployUpkeepRegistrationRequests(
				s.Wallets.Default(),
				s.Link.Address(),
				big.NewInt(0),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = registry.SetRegistrar(s.Wallets.Default(), registrar.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = registrar.SetRegistrarConfig(
				s.Wallets.Default(),
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
				s.Wallets.Default().Address(),
				[]byte("0x"),
				big.NewInt(9e18),
				0,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Link.TransferAndCall(s.Wallets.Default(), registrar.Address(), big.NewInt(9e18), req)
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
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
			err = registry.SetKeepers(s.Wallets.Default(), nodeAddressesStr, payees)
			Expect(err).ShouldNot(HaveOccurred())
			_, err = nodes[0].CreateJob(&client.KeeperJobSpec{
				Name:            "keeper",
				ContractAddress: registry.Address(),
				FromAddress:     na,
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Describe("with Keeper job", func() {
		It("performs upkeep of a target contract", func() {
			err = retry.Do(func() error {
				cnt, err := consumer.Counter(context.Background())
				if err != nil {
					return err
				}
				if cnt.Int64() == 0 {
					return errors.New("awaiting for upkeep")
				}
				log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
				return nil
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	AfterEach(func() {
		By("Tearing down the environment", s.TearDown())
	})
})
