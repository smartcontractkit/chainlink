package smoke

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
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
)

var _ = Describe("VRF suite @vrf", func() {
	var (
		err                error
		nets               *blockchain.Networks
		cd                 contracts.ContractDeployer
		consumer           contracts.VRFConsumer
		coordinator        contracts.VRFCoordinator
		encodedProvingKeys = make([][2]*big.Int, 0)
		lt                 contracts.LinkToken
		chainlinkNodes     []client.Chainlink
		e                  *environment.Environment
		job                *client.Job
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(
					config.ChainlinkVals(),
					"chainlink-vrf-core-ci",
					config.GethNetworks()...,
				),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
		})

		By("Connecting to launched resources", func() {
			networkRegistry := blockchain.NewDefaultNetworkRegistry()
			nets, err = networkRegistry.GetNetworks(e)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(e)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			nets.Default.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := nets.Default.EstimateCostForChainlinkOperations(1)
			Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
			err = actions.FundChainlinkNodes(chainlinkNodes, nets.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred(), "Funding chainlink nodes with ETH shouldn't fail")
		})

		By("Deploying VRF contracts", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
			bhs, err := cd.DeployBlockhashStore()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Blockhash store shouldn't fail")
			coordinator, err = cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred(), "Deploying VRF coordinator shouldn't fail")
			consumer, err = cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred(), "Deploying VRF consumer contract shouldn't fail")
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for VRF setup contracts to deploy")

			err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred(), "Funding consumer contract shouldn't fail")
			_, err = cd.DeployVRFContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying VRF contract shouldn't fail")
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
		})

		By("Creating jobs and registering proving keys", func() {
			for _, n := range chainlinkNodes {
				nodeKey, err := n.CreateVRFKey()
				Expect(err).ShouldNot(HaveOccurred(), "Creating VRF key shouldn't fail")
				log.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
				pubKeyCompressed := nodeKey.Data.ID
				jobUUID := uuid.NewV4()
				os := &client.VRFTxPipelineSpec{
					Address: coordinator.Address(),
				}
				ost, err := os.String()
				Expect(err).ShouldNot(HaveOccurred(), "Building observation source spec shouldn't fail")
				job, err = n.CreateJob(&client.VRFJobSpec{
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
				encodedProvingKeys = append(encodedProvingKeys, provingKey)
			}
		})
	})

	Describe("with VRF job", func() {
		It("randomness is fulfilled", func() {
			requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
			Expect(err).ShouldNot(HaveOccurred(), "Getting Hash of encoded proving keys shouldn't fail")
			err = consumer.RequestRandomness(requestHash, big.NewInt(1))
			Expect(err).ShouldNot(HaveOccurred(), "Requesting randomness shouldn't fail")

			timeout := time.Minute * 2

			Eventually(func(g Gomega) {
				jobRuns, err := chainlinkNodes[0].ReadRunsByJob(job.Data.ID)
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
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			nets.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nets, utils.ProjectRoot, chainlinkNodes, nil)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
