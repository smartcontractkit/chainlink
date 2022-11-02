package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var _ = Describe("VRF suite @vrf", func() {
	var (
		testScenarios = []TableEntry{
			Entry("VRF suite on default environment @default", defaultVRFEnv()),
		}

		testEnvironment *environment.Environment
		chainClient     blockchain.EVMClient
		chainlinkNodes  []*client.Chainlink
	)

	AfterEach(func() {
		By("Tearing env down")
		chainClient.GasStats().PrintStats()
		err := actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("VRF suite on different EVM networks", func(
		testInputs *smokeTestInputs,
	) {
		By("Deploying the environment")
		testEnvironment = testInputs.environment
		testNetwork := testInputs.network
		err := testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = blockchain.NewEVMClient(testNetwork, testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting client shouldn't fail")
		cd, err := contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		chainClient.ParallelTransactions(true)

		By("Funding Chainlink nodes")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.01))
		Expect(err).ShouldNot(HaveOccurred(), "Funding chainlink nodes with ETH shouldn't fail")

		By("Deploying VRF contracts")
		lt, err := cd.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
		bhs, err := cd.DeployBlockhashStore()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Blockhash store shouldn't fail")
		coordinator, err := cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
		Expect(err).ShouldNot(HaveOccurred(), "Deploying VRF coordinator shouldn't fail")
		consumer, err := cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
		Expect(err).ShouldNot(HaveOccurred(), "Deploying VRF consumer contract shouldn't fail")
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for VRF setup contracts to deploy")

		err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
		Expect(err).ShouldNot(HaveOccurred(), "Funding consumer contract shouldn't fail")
		_, err = cd.DeployVRFContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying VRF contract shouldn't fail")
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")

		for _, n := range chainlinkNodes {
			nodeKey, err := n.MustCreateVRFKey()
			Expect(err).ShouldNot(HaveOccurred(), "Creating VRF key shouldn't fail")
			log.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
			pubKeyCompressed := nodeKey.Data.ID
			jobUUID := uuid.NewV4()
			os := &client.VRFTxPipelineSpec{
				Address: coordinator.Address(),
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred(), "Building observation source spec shouldn't fail")
			job, err := n.MustCreateJob(&client.VRFJobSpec{
				Name:                     fmt.Sprintf("vrf-%s", jobUUID),
				CoordinatorAddress:       coordinator.Address(),
				MinIncomingConfirmations: 1,
				PublicKey:                pubKeyCompressed,
				ExternalJobID:            jobUUID.String(),
				ObservationSource:        ost,
			})
			Expect(err).ShouldNot(HaveOccurred(), "Creating VRF Job shouldn't fail")

			oracleAddr, err := n.PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred(), "Getting primary ETH address of chainlink node shouldn't fail")
			provingKey, err := actions.EncodeOnChainVRFProvingKey(*nodeKey)
			Expect(err).ShouldNot(HaveOccurred(), "Encoding on-chain VRF Proving key shouldn't fail")
			err = coordinator.RegisterProvingKey(
				big.NewInt(1),
				oracleAddr,
				provingKey,
				actions.EncodeOnChainExternalJobID(jobUUID),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Registering the on-chain VRF Proving key shouldn't fail")
			encodedProvingKeys := make([][2]*big.Int, 0)
			encodedProvingKeys = append(encodedProvingKeys, provingKey)

			requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
			Expect(err).ShouldNot(HaveOccurred(), "Getting Hash of encoded proving keys shouldn't fail")
			err = consumer.RequestRandomness(requestHash, big.NewInt(1))
			Expect(err).ShouldNot(HaveOccurred(), "Requesting randomness shouldn't fail")

			By("Checking that randomness fulfilled")
			timeout := time.Minute * 2
			Eventually(func(g Gomega) {
				jobRuns, err := chainlinkNodes[0].MustReadRunsByJob(job.Data.ID)
				g.Expect(err).ShouldNot(HaveOccurred(), "Job execution shouldn't fail")

				out, err := consumer.RandomnessOutput(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Getting the randomness output of the consumer shouldn't fail")
				// Checks that the job has actually run
				g.Expect(len(jobRuns.Data)).Should(BeNumerically(">=", 1),
					fmt.Sprintf("Expected the VRF job to run once or more after %s", timeout))

				// TODO: This is an imperfect check, given it's a random number, it CAN be 0, but chances are unlikely.
				// So we're just checking that the answer has changed to something other than the default (0)
				// There's a better formula to ensure that VRF response is as expected, detailed under Technical Walkthrough.
				// https://blog.chain.link/chainlink-vrf-on-chain-verifiable-randomness/
				g.Expect(out.Uint64()).Should(Not(BeNumerically("==", 0)), "Expected the VRF job give an answer other than 0")
				log.Debug().Uint64("Output", out.Uint64()).Msg("Randomness fulfilled")
			}, timeout, "1s").Should(Succeed())
		}
	},
		testScenarios,
	)
})

func defaultVRFEnv() *smokeTestInputs {
	network := networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !network.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	env := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-vrf-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
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
