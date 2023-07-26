package smoke

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
)

func TestVRFBasic(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	env, err := docker.NewChainlinkCluster(t, 1)
	require.NoError(t, err)
	clients, err := docker.ConnectClients(env)
	require.NoError(t, err)
	clients.Networks[0].ParallelTransactions(true)

	err = actions.FundChainlinkNodes(clients.Chainlink, clients.Networks[0], big.NewFloat(.01))
	require.NoError(t, err, "Funding chainlink nodes with ETH shouldn't fail")

	lt, err := clients.NetworkDeployers[0].DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	bhs, err := clients.NetworkDeployers[0].DeployBlockhashStore()
	require.NoError(t, err, "Deploying Blockhash store shouldn't fail")
	coordinator, err := clients.NetworkDeployers[0].DeployVRFCoordinator(lt.Address(), bhs.Address())
	require.NoError(t, err, "Deploying VRF coordinator shouldn't fail")
	consumer, err := clients.NetworkDeployers[0].DeployVRFConsumer(lt.Address(), coordinator.Address())
	require.NoError(t, err, "Deploying VRF consumer contract shouldn't fail")
	err = clients.Networks[0].WaitForEvents()
	require.NoError(t, err, "Failed to wait for VRF setup contracts to deploy")

	err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
	require.NoError(t, err, "Funding consumer contract shouldn't fail")
	_, err = clients.NetworkDeployers[0].DeployVRFContract()
	require.NoError(t, err, "Deploying VRF contract shouldn't fail")
	err = clients.Networks[0].WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")

	for _, n := range clients.Chainlink {
		nodeKey, err := n.MustCreateVRFKey()
		require.NoError(t, err, "Creating VRF key shouldn't fail")
		l.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
		pubKeyCompressed := nodeKey.Data.ID
		jobUUID := uuid.New()
		os := &client.VRFTxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		require.NoError(t, err, "Building observation source spec shouldn't fail")
		job, err := n.MustCreateJob(&client.VRFJobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			MinIncomingConfirmations: 1,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
		})
		require.NoError(t, err, "Creating VRF Job shouldn't fail")

		oracleAddr, err := n.PrimaryEthAddress()
		require.NoError(t, err, "Getting primary ETH address of chainlink node shouldn't fail")
		provingKey, err := actions.EncodeOnChainVRFProvingKey(*nodeKey)
		require.NoError(t, err, "Encoding on-chain VRF Proving key shouldn't fail")
		err = coordinator.RegisterProvingKey(
			big.NewInt(1),
			oracleAddr,
			provingKey,
			actions.EncodeOnChainExternalJobID(jobUUID),
		)
		require.NoError(t, err, "Registering the on-chain VRF Proving key shouldn't fail")
		encodedProvingKeys := make([][2]*big.Int, 0)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)

		requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
		require.NoError(t, err, "Getting Hash of encoded proving keys shouldn't fail")
		err = consumer.RequestRandomness(requestHash, big.NewInt(1))
		require.NoError(t, err, "Requesting randomness shouldn't fail")

		gom := gomega.NewGomegaWithT(t)
		timeout := time.Minute * 2
		gom.Eventually(func(g gomega.Gomega) {
			jobRuns, err := clients.Chainlink[0].MustReadRunsByJob(job.Data.ID)
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Job execution shouldn't fail")

			out, err := consumer.RandomnessOutput(context.Background())
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Getting the randomness output of the consumer shouldn't fail")
			// Checks that the job has actually run
			g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 1),
				fmt.Sprintf("Expected the VRF job to run once or more after %s", timeout))

			// TODO: This is an imperfect check, given it's a random number, it CAN be 0, but chances are unlikely.
			// So we're just checking that the answer has changed to something other than the default (0)
			// There's a better formula to ensure that VRF response is as expected, detailed under Technical Walkthrough.
			// https://bl.chain.link/chainlink-vrf-on-chain-verifiable-randomness/
			g.Expect(out.Uint64()).ShouldNot(gomega.BeNumerically("==", 0), "Expected the VRF job give an answer other than 0")
			l.Debug().Uint64("Output", out.Uint64()).Msg("Randomness fulfilled")
		}, timeout, "1s").Should(gomega.Succeed())
	}
}
