package integration

import (
	"context"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Basic Contract Interactions @contract", func() {
	var (
		suiteSetup    actions.SuiteSetup
		networkInfo   actions.NetworkInfo
		defaultWallet client.BlockchainWallet
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			var err error
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkCluster(0),
				client.DefaultNetworkFromConfig,
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())
			networkInfo = suiteSetup.DefaultNetwork()
			defaultWallet = networkInfo.Wallets.Default()
		})
	})

	It("can deploy all contracts", func() {
		By("basic interaction with a storage contract", func() {
			storeInstance, err := networkInfo.Deployer.DeployStorageContract(defaultWallet)
			Expect(err).ShouldNot(HaveOccurred())
			testVal := big.NewInt(5)
			err = storeInstance.Set(testVal)
			Expect(err).ShouldNot(HaveOccurred())
			val, err := storeInstance.Get(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).To(Equal(testVal))
		})

		By("deploying the flux monitor contract", func() {
			rac, err := networkInfo.Deployer.DeployReadAccessController(networkInfo.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			flags, err := networkInfo.Deployer.DeployFlags(networkInfo.Wallets.Default(), rac.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = networkInfo.Deployer.DeployDeviationFlaggingValidator(networkInfo.Wallets.Default(), flags.Address(), big.NewInt(0))
			Expect(err).ShouldNot(HaveOccurred())
			fluxOptions := contracts.DefaultFluxAggregatorOptions()
			_, err = networkInfo.Deployer.DeployFluxAggregatorContract(defaultWallet, fluxOptions)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("deploying the ocr contract", func() {
			ocrOptions := contracts.DefaultOffChainAggregatorOptions()
			_, err := networkInfo.Deployer.DeployOffChainAggregator(defaultWallet, ocrOptions)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("deploying keeper contracts", func() {
			ef, err := networkInfo.Deployer.DeployMockETHLINKFeed(networkInfo.Wallets.Default(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := networkInfo.Deployer.DeployMockGasFeed(networkInfo.Wallets.Default(), big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = networkInfo.Deployer.DeployKeeperRegistry(
				networkInfo.Wallets.Default(),
				&contracts.KeeperRegistryOpts{
					LinkAddr:             networkInfo.Link.Address(),
					ETHFeedAddr:          ef.Address(),
					GasFeedAddr:          gf.Address(),
					PaymentPremiumPPB:    uint32(200000000),
					BlockCountPerTurn:    big.NewInt(3),
					CheckGasLimit:        uint32(2500000),
					StalenessSeconds:     big.NewInt(90000),
					GasCeilingMultiplier: uint16(1),
					FallbackGasPrice:     big.NewInt(2e11),
					FallbackLinkPrice:    big.NewInt(2e18),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("deploying vrf contract", func() {
			bhs, err := networkInfo.Deployer.DeployBlockhashStore(networkInfo.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err := networkInfo.Deployer.DeployVRFCoordinator(networkInfo.Wallets.Default(), networkInfo.Link.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = networkInfo.Deployer.DeployVRFConsumer(networkInfo.Wallets.Default(), networkInfo.Link.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = networkInfo.Deployer.DeployVRFContract(networkInfo.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("deploying direct request contract", func() {
			_, err := networkInfo.Deployer.DeployOracle(networkInfo.Wallets.Default(), networkInfo.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = networkInfo.Deployer.DeployAPIConsumer(networkInfo.Wallets.Default(), networkInfo.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			networkInfo.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
