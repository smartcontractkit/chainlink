package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	ethdeploy "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("VRFv2 suite @v2vrf", func() {
	var (
		testScenarios = []TableEntry{
			Entry("VRFv2 suite on Simulated Network @simulated", defaultVRFv2Env()),
		}

		testEnvironment *environment.Environment
		chainClient     blockchain.EVMClient
		chainlinkNodes  []*client.Chainlink
		// used both as a feed and a fallback value
		linkEthFeedResponse = big.NewInt(1e18)
	)

	AfterEach(func() {
		By("Tearing env down")
		chainClient.GasStats().PrintStats()
		err := actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("VRFv2 suite on different EVM networks", func(
		testInputs *smokeTestInputs,
	) {
		By("Deploying the environment")
		testEnvironment = testInputs.environment
		testNetwork := testInputs.network
		err := testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(testNetwork, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred())
		contractDeployer, err := contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred())
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred())
		chainClient.ParallelTransactions(true)

		By("Deploying VRF contracts")
		linkToken, err := contractDeployer.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred())
		bhs, err := contractDeployer.DeployBlockhashStore()
		Expect(err).ShouldNot(HaveOccurred())
		mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
		Expect(err).ShouldNot(HaveOccurred())
		coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkToken.Address(), bhs.Address(), mf.Address())
		Expect(err).ShouldNot(HaveOccurred())
		consumer, err := contractDeployer.DeployVRFConsumerV2(linkToken.Address(), coordinator.Address())
		Expect(err).ShouldNot(HaveOccurred())
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(1))
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

		By("Creating jobs and registering proving keys")
		var (
			job                *client.Job
			encodedProvingKeys = make([][2]*big.Int, 0)
		)
		for _, n := range chainlinkNodes {
			vrfKey, err := n.MustCreateVRFKey()
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
			job, err = n.MustCreateJob(&client.VRFV2JobSpec{
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

		By("randomness is fulfilled")
		words := uint32(10)
		keyHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
		Expect(err).ShouldNot(HaveOccurred())
		err = consumer.RequestRandomness(keyHash, 1, 1, 300000, words)
		Expect(err).ShouldNot(HaveOccurred())

		timeout := time.Minute * 2

		Eventually(func(g Gomega) {
			jobRuns, err := chainlinkNodes[0].MustReadRunsByJob(job.Data.ID)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(len(jobRuns.Data)).Should(BeNumerically("==", 1))
			randomness, err := consumer.GetAllRandomWords(context.Background(), int(words))
			g.Expect(err).ShouldNot(HaveOccurred())
			for _, w := range randomness {
				log.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
				g.Expect(w.Uint64()).Should(Not(BeNumerically("==", 0)), "Expected the VRF job give an answer other than 0")
			}
		}, timeout, "1s").Should(Succeed())
	},
		testScenarios,
	)
})

func defaultVRFv2Env() *smokeTestInputs {
	network := networks.SelectedNetwork
	evmConfig := ethdeploy.New(nil)
	if !network.Simulated {
		evmConfig = ethdeploy.New(&ethdeploy.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	env := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-vrfv2-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
	}).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"env": network.ChainlinkValuesMap(),
		}))
	return &smokeTestInputs{
		network:     network,
		environment: env,
	}
}
