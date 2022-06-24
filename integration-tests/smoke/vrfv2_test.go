package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	it "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
)

var _ = Describe("VRFv2 suite @v2vrf", func() {
	DescribeTable("VRFv2 suite on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		networkChart environment.ConnectedChart,
		clChart environment.ConnectedChart,
	) {
		var (
			err                error
			c                  blockchain.EVMClient
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
		By("Deploying the environment", func() {
			e = environment.New(nil).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(eth.New(nil)).
				AddHelm(chainlink.New(0, nil))
			err = e.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			log.Trace().Msg("JUST A TRACE")
			c, err = blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings)(e)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(c)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.ConnectChainlinkNodes(e)
			Expect(err).ShouldNot(HaveOccurred())
			c.ParallelTransactions(true)
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
			err = actions.FundChainlinkNodes(cls, c, big.NewFloat(10))
			Expect(err).ShouldNot(HaveOccurred())
			err = c.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// https://docs.chain.link/docs/chainlink-vrf/#subscription-limits
			linkFunding := big.NewInt(100)

			err = lt.Transfer(consumer.Address(), big.NewInt(0).Mul(linkFunding, big.NewInt(1e18)))
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
			err = c.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			err = consumer.CreateFundedSubscription(big.NewInt(0).Mul(linkFunding, big.NewInt(1e18)))
			Expect(err).ShouldNot(HaveOccurred())
			err = c.WaitForEvents()
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
					EVMChainID:               fmt.Sprint(c.GetNetworkConfig().ChainID),
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
		By("randomness is fulfilled", func() {
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
		By("Printing gas stats", func() {
			c.GasStats().PrintStats()
		})

		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, utils.ProjectRoot, nil, nil, c)
			Expect(err).ShouldNot(HaveOccurred())
		})
	},
		Entry("VRFv2 on Geth @geth",
			blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings),
			eth.New(nil),
			chainlink.New(0, nil),
		),
	)
})
