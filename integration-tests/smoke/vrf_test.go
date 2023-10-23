package smoke

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv1"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestVRFBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(big.NewFloat(.1)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)
	env.ParallelTransactions(true)

	lt, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	contracts, err := vrfv1.DeployVRFContracts(env.ContractDeployer, env.EVMClient, lt)
	require.NoError(t, err, "Deploying VRF Contracts shouldn't fail")

	err = lt.Transfer(contracts.Consumer.Address(), big.NewInt(2e18))
	require.NoError(t, err, "Funding consumer contract shouldn't fail")
	_, err = env.ContractDeployer.DeployVRFContract()
	require.NoError(t, err, "Deploying VRF contract shouldn't fail")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")

	for _, n := range env.ClCluster.Nodes {
		nodeKey, err := n.API.MustCreateVRFKey()
		require.NoError(t, err, "Creating VRF key shouldn't fail")
		l.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
		pubKeyCompressed := nodeKey.Data.ID
		jobUUID := uuid.New()
		os := &client.VRFTxPipelineSpec{
			Address: contracts.Coordinator.Address(),
		}
		ost, err := os.String()
		require.NoError(t, err, "Building observation source spec shouldn't fail")
		job, err := n.API.MustCreateJob(&client.VRFJobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       contracts.Coordinator.Address(),
			MinIncomingConfirmations: 1,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			EVMChainID:               env.EVMClient.GetChainID().String(),
			ObservationSource:        ost,
		})
		require.NoError(t, err, "Creating VRF Job shouldn't fail")

		oracleAddr, err := n.API.PrimaryEthAddress()
		require.NoError(t, err, "Getting primary ETH address of chainlink node shouldn't fail")
		provingKey, err := actions.EncodeOnChainVRFProvingKey(*nodeKey)
		require.NoError(t, err, "Encoding on-chain VRF Proving key shouldn't fail")
		err = contracts.Coordinator.RegisterProvingKey(
			big.NewInt(1),
			oracleAddr,
			provingKey,
			actions.EncodeOnChainExternalJobID(jobUUID),
		)
		require.NoError(t, err, "Registering the on-chain VRF Proving key shouldn't fail")
		encodedProvingKeys := make([][2]*big.Int, 0)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)

		requestHash, err := contracts.Coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
		require.NoError(t, err, "Getting Hash of encoded proving keys shouldn't fail")
		err = contracts.Consumer.RequestRandomness(requestHash, big.NewInt(1))
		require.NoError(t, err, "Requesting randomness shouldn't fail")

		gom := gomega.NewGomegaWithT(t)
		timeout := time.Minute * 2
		gom.Eventually(func(g gomega.Gomega) {
			jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(job.Data.ID)
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Job execution shouldn't fail")

			out, err := contracts.Consumer.RandomnessOutput(context.Background())
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
