//go:build integration

package integration

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/environment"

	"github.com/smartcontractkit/integrations-framework/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Basic Contract Interactions @contract", func() {
	var (
		suiteSetup actions.SuiteSetup

		firstNetwork       actions.NetworkInfo
		firstNetworkWallet client.BlockchainWallet

		secondNetwork       actions.NetworkInfo
		secondNetworkWallet client.BlockchainWallet
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			var err error
			suiteSetup, err = actions.MultiNetworkSetup(
				environment.NewChainlinkCluster(0),
				client.DefaultNetworksFromConfig,
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())

			firstNetwork, err = suiteSetup.Network(0)
			Expect(err).ShouldNot(HaveOccurred())
			secondNetwork, err = suiteSetup.Network(1)
			Expect(err).ShouldNot(HaveOccurred())
			firstNetworkWallet = firstNetwork.Wallets.Default()
			secondNetworkWallet = secondNetwork.Wallets.Default()
		})
	})

	It("exercises basic contract usage", func() {
		By("deploying the storage contract", func() {
			// Deploy storage
			firstStoreInstance, err := firstNetwork.Deployer.DeployStorageContract(firstNetworkWallet)
			Expect(err).ShouldNot(HaveOccurred())
			secondStoreInstance, err := secondNetwork.Deployer.DeployStorageContract(secondNetworkWallet)
			Expect(err).ShouldNot(HaveOccurred())

			firstNetworkTestVal := big.NewInt(5)
			secondNetworkTestVal := big.NewInt(10)

			// Set both values
			err = firstStoreInstance.Set(firstNetworkTestVal)
			Expect(err).ShouldNot(HaveOccurred())
			err = secondStoreInstance.Set(secondNetworkTestVal)
			Expect(err).ShouldNot(HaveOccurred())

			// Check Answers
			val, err := firstStoreInstance.Get(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).To(Equal(firstNetworkTestVal))

			val, err = secondStoreInstance.Get(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).To(Equal(secondNetworkTestVal))
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
