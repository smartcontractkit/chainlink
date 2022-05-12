package smoke_test

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
)

var _ = Describe("VRFv2 suite @v2vrf", func() {
	var (
		err                error
		nets               *blockchain.Networks
		cd                 contracts.ContractDeployer
		consumer           contracts.VRFConsumerV2
		coordinator        contracts.VRFCoordinatorV2
		encodedProvingKeys = make([][2]*big.Int, 0)
		lt                 contracts.LinkToken
		cls                []client.Chainlink
		e                  *environment.Environment
		vrfKey             *client.VRFKey
		job                *client.Job
		// used both as a feed and a fallback value
		linkEthFeedResponse = big.NewInt(1e18)
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(
					config.ChainlinkVals(),
					"chainlink-vrfv2-core-ci",
					environment.PerformanceGeth,
				),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			networkRegistry := blockchain.NewDefaultNetworkRegistry()
			nets, err = networkRegistry.GetNetworks(e)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.ConnectChainlinkNodes(e)
			Expect(err).ShouldNot(HaveOccurred())
			nets.Default.ParallelTransactions(true)
		})
		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(cls, nets.Default, big.NewFloat(3))
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying VRF contracts", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			bhs, err := cd.DeployBlockhashStore()
			Expect(err).ShouldNot(HaveOccurred())
			mf, err := cd.DeployMockETHLINKFeed(linkEthFeedResponse)
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err = cd.DeployVRFCoordinatorV2(lt.Address(), bhs.Address(), mf.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = cd.DeployVRFConsumerV2(lt.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			err = lt.Transfer(consumer.Address(), big.NewInt(0).Mul(big.NewInt(1e4), big.NewInt(1e18)))
			Expect(err).ShouldNot(HaveOccurred())
			err = coordinator.SetConfig(
				1,
				2.5e6,
				86400,
				33825,
				linkEthFeedResponse,
				ethereum.VRFCoordinatorV2FeeConfig{
					FulfillmentFlatFeeLinkPPMTier1: 1,
					FulfillmentFlatFeeLinkPPMTier2: 1,
					FulfillmentFlatFeeLinkPPMTier3: 1,
					FulfillmentFlatFeeLinkPPMTier4: 1,
					FulfillmentFlatFeeLinkPPMTier5: 1,
					ReqsForTier2:                   big.NewInt(10),
					ReqsForTier3:                   big.NewInt(20),
					ReqsForTier4:                   big.NewInt(30),
					ReqsForTier5:                   big.NewInt(40)},
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			err = consumer.CreateFundedSubscription(big.NewInt(0).Mul(big.NewInt(30), big.NewInt(1e18)))
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating jobs and registering proving keys", func() {
			for _, n := range cls {
				vrfKey, err = n.CreateVRFKey()
				Expect(err).ShouldNot(HaveOccurred())
				log.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
				pubKeyCompressed := vrfKey.Data.ID
				jobUUID := uuid.NewV4()
				os := &client.VRFV2TxPipelineSpec{
					Address: coordinator.Address(),
				}
				ost, err := os.String()
				Expect(err).ShouldNot(HaveOccurred())
				oracleAddr, err := n.PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				job, err = n.CreateJob(&client.VRFV2JobSpec{
					Name:                     fmt.Sprintf("vrf-%s", jobUUID),
					CoordinatorAddress:       coordinator.Address(),
					FromAddress:              oracleAddr,
					EVMChainID:               "1337",
					MinIncomingConfirmations: 1,
					PublicKey:                pubKeyCompressed,
					ExternalJobID:            jobUUID.String(),
					ObservationSource:        ost,
					BatchFulfillmentEnabled:  false,
				})
				Expect(err).ShouldNot(HaveOccurred())
				provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
				Expect(err).ShouldNot(HaveOccurred())
				err = coordinator.RegisterProvingKey(
					oracleAddr,
					provingKey,
				)
				Expect(err).ShouldNot(HaveOccurred())
				encodedProvingKeys = append(encodedProvingKeys, provingKey)
			}
		})
	})

	Describe("with VRF job", func() {
		It("randomness is fulfilled", func() {
			words := uint32(10)
			keyHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.RequestRandomness(keyHash, 1, 1, 300000, words)
			Expect(err).ShouldNot(HaveOccurred())

			timeout := time.Minute * 2

			Eventually(func(g Gomega) {
				jobRuns, err := cls[0].ReadRunsByJob(job.Data.ID)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(len(jobRuns.Data)).Should(BeNumerically("==", 1))
				randomness, err := consumer.GetAllRandomWords(context.Background(), int(words))
				g.Expect(err).ShouldNot(HaveOccurred())
				for _, w := range randomness {
					log.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
					g.Expect(w.Uint64()).Should(Not(BeNumerically("==", 0)), "Expected the VRF job give an answer other than 0")
				}
			}, timeout, "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			nets.Default.GasStats().PrintStats()
		})

		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nets, utils.ProjectRoot, nil, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
