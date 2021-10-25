//go:build integration

package integration

import (
	"context"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("VRF suite @vrf", func() {

	var (
		suiteSetup         actions.SuiteSetup
		networkInfo        actions.NetworkInfo
		nodes              []client.Chainlink
		consumer           contracts.VRFConsumer
		coordinator        contracts.VRFCoordinator
		encodedProvingKeys = make([][2]*big.Int, 0)
		err                error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkCluster(1),
				client.DefaultNetworkFromConfig,
				"../",
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())
			networkInfo = suiteSetup.DefaultNetwork()

			networkInfo.Client.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			ethAmount, err := networkInfo.Deployer.CalculateETHForTXs(networkInfo.Wallets.Default(), networkInfo.Network.Config(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(nodes, networkInfo.Client, networkInfo.Wallets.Default(), ethAmount, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying VRF contracts", func() {
			bhs, err := networkInfo.Deployer.DeployBlockhashStore(networkInfo.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err = networkInfo.Deployer.DeployVRFCoordinator(networkInfo.Wallets.Default(), networkInfo.Link.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = networkInfo.Deployer.DeployVRFConsumer(networkInfo.Wallets.Default(), networkInfo.Link.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.Fund(networkInfo.Wallets.Default(), big.NewFloat(0), big.NewFloat(2))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = networkInfo.Deployer.DeployVRFContract(networkInfo.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = networkInfo.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating jobs and registering proving keys", func() {
			for _, n := range nodes {
				nodeKeys, err := n.ReadVRFKeys()
				Expect(err).ShouldNot(HaveOccurred())
				log.Debug().Interface("Key JSON", nodeKeys).Msg("Created proving key")
				pubKeyCompressed := nodeKeys.Data[0].ID
				jobUUID := uuid.NewV4()
				os := &client.VRFTxPipelineSpec{
					Address: coordinator.Address(),
				}
				ost, err := os.String()
				Expect(err).ShouldNot(HaveOccurred())
				_, err = n.CreateJob(&client.VRFJobSpec{
					Name:               "vrf",
					CoordinatorAddress: coordinator.Address(),
					PublicKey:          pubKeyCompressed,
					Confirmations:      1,
					ExternalJobID:      jobUUID.String(),
					ObservationSource:  ost,
				})
				Expect(err).ShouldNot(HaveOccurred())

				oracleAddr, err := n.PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				provingKey, err := actions.EncodeOnChainVRFProvingKey(nodeKeys.Data[0])
				Expect(err).ShouldNot(HaveOccurred())
				err = coordinator.RegisterProvingKey(
					networkInfo.Wallets.Default(),
					big.NewInt(1),
					oracleAddr,
					provingKey,
					actions.EncodeOnChainExternalJobID(jobUUID),
				)
				Expect(err).ShouldNot(HaveOccurred())
				encodedProvingKeys = append(encodedProvingKeys, provingKey)
			}
		})
	})

	Describe("with VRF job", func() {
		It("fulfills randomness", func() {
			requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.RequestRandomness(networkInfo.Wallets.Default(), requestHash, big.NewInt(1))
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				out, err := consumer.RandomnessOutput(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(out.Uint64()).Should(Not(BeNumerically("==", 0)))
				log.Debug().Uint64("Output", out.Uint64()).Msg("Randomness fulfilled")
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
