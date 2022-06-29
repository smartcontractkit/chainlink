package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	ethdeploy "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("VRFv2 suite @v2vrf", func() {
	DescribeTable("VRFv2 suite on different EVM networks", func(
		clientFunc func(networkSettings *blockchain.EVMNetwork) func(*environment.Environment) (blockchain.EVMClient, error),
		evmNetwork *blockchain.EVMNetwork,
		chainlinkValues map[string]interface{},
	) {
		var (
			err                error
			chainClient        blockchain.EVMClient
			contractDeployer   contracts.ContractDeployer
			consumer           contracts.VRFConsumerV2
			coordinator        contracts.VRFCoordinatorV2
			encodedProvingKeys = make([][2]*big.Int, 0)
			linkToken          contracts.LinkToken
			chainlinkNodes     []client.Chainlink
			testEnvironment    *environment.Environment
			vrfKey             *client.VRFKey
			job                *client.Job
			// used both as a feed and a fallback value
			linkEthFeedResponse = big.NewInt(1e18)
		)

		By("Deploying the environment", func() {
			testEnvironment = environment.New(&environment.Config{NamespacePrefix: "smoke-vrfv2"}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethdeploy.New(&ethdeploy.Props{
					NetworkName: evmNetwork.Name,
					Simulated:   evmNetwork.Simulated,
					WsURLs:      evmNetwork.URLs,
				})).
				AddHelm(chainlink.New(0, chainlinkValues))
			err = testEnvironment.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			log.Trace().Msg("JUST A TRACE")
			chainClient, err = clientFunc(evmNetwork)(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred())
			contractDeployer, err = contracts.NewContractDeployer(chainClient)
			Expect(err).ShouldNot(HaveOccurred())
			chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred())
			chainClient.ParallelTransactions(true)
		})

		By("Deploying VRF contracts", func() {
			linkToken, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			bhs, err := contractDeployer.DeployBlockhashStore()
			Expect(err).ShouldNot(HaveOccurred())
			mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err = contractDeployer.DeployVRFCoordinatorV2(linkToken.Address(), bhs.Address(), mf.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = contractDeployer.DeployVRFConsumerV2(linkToken.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(10))
			Expect(err).ShouldNot(HaveOccurred())
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// https://docs.chain.link/docs/chainlink-vrf/#subscription-limits
			linkFunding := big.NewInt(100)

			err = linkToken.Transfer(consumer.Address(), big.NewInt(0).Mul(linkFunding, big.NewInt(1e18)))
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
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			err = consumer.CreateFundedSubscription(big.NewInt(0).Mul(linkFunding, big.NewInt(1e18)))
			Expect(err).ShouldNot(HaveOccurred())
			err = chainClient.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating jobs and registering proving keys", func() {
			for _, n := range chainlinkNodes {
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
					EVMChainID:               fmt.Sprint(chainClient.GetNetworkConfig().ChainID),
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
				jobRuns, err := chainlinkNodes[0].ReadRunsByJob(job.Data.ID)
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
			chainClient.GasStats().PrintStats()
		})

		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(testEnvironment, utils.ProjectRoot, nil, nil, chainClient)
			Expect(err).ShouldNot(HaveOccurred())
		})
	},
		Entry("VRFv2 suite on Simulated Network @simulated", blockchain.NewEthereumMultiNodeClientSetup, networks.SimulatedEVMNetwork, nil),
		Entry("VRFv2 suite on Metis Stardust @metis", blockchain.NewMetisMultiNodeClientSetup, networks.MetisTestNetwork, map[string]interface{}{
			"env": map[string]interface{}{
				"eth_url":      networks.MetisTestNetwork.URLs[0],
				"eth_chain_id": fmt.Sprint(networks.MetisTestNetwork.ChainID),
			},
		}),
	)
})
