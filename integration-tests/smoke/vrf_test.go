//go:build smoke
package smoke

import (
	"context"
	"fmt"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/utils"
)

var _ = Describe("VRF suite @vrf", func() {
	var (
		err                error
		nets               *client.Networks
		cd                 contracts.ContractDeployer
		consumer           contracts.VRFConsumer
		coordinator        contracts.VRFCoordinator
		encodedProvingKeys = make([][2]*big.Int, 0)
		lt                 contracts.LinkToken
		cls                []client.Chainlink
		e                  *environment.Environment
		job                *client.Job
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(nil),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Getting the clients", func() {
			networkRegistry := client.NewNetworkRegistry()
			nets, err = networkRegistry.GetNetworks(e)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.NewChainlinkClients(e)
			Expect(err).ShouldNot(HaveOccurred())
			nets.Default.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := nets.Default.EstimateCostForChainlinkOperations(1)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(cls, nets.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying VRF contracts", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			bhs, err := cd.DeployBlockhashStore()
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err = cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = cd.DeployVRFContract()
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating jobs and registering proving keys", func() {
			for _, n := range cls {
				nodeKey, err := n.CreateVRFKey()
				Expect(err).ShouldNot(HaveOccurred())
				log.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
				pubKeyCompressed := nodeKey.Data.ID
				jobUUID := uuid.NewV4()
				os := &client.VRFTxPipelineSpec{
					Address: coordinator.Address(),
				}
				ost, err := os.String()
				Expect(err).ShouldNot(HaveOccurred())
				job, err = n.CreateJob(&client.VRFJobSpec{
					Name:                     fmt.Sprintf("vrf-%s", jobUUID),
					CoordinatorAddress:       coordinator.Address(),
					MinIncomingConfirmations: 1,
					PublicKey:                pubKeyCompressed,
					ExternalJobID:            jobUUID.String(),
					ObservationSource:        ost,
				})
				Expect(err).ShouldNot(HaveOccurred())

				oracleAddr, err := n.PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				provingKey, err := actions.EncodeOnChainVRFProvingKey(*nodeKey)
				Expect(err).ShouldNot(HaveOccurred())
				err = coordinator.RegisterProvingKey(
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
		It("randomness is fulfilled", func() {
			requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.RequestRandomness(requestHash, big.NewInt(1))
			Expect(err).ShouldNot(HaveOccurred())

			timeout := time.Minute * 2

			Eventually(func(g Gomega) {
				jobRuns, err := cls[0].ReadRunsByJob(job.Data.ID)
				g.Expect(err).ShouldNot(HaveOccurred())

				out, err := consumer.RandomnessOutput(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
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
			err = actions.TeardownSuite(e, nets, utils.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
