package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	vrfConst "github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

func TestVRFv2Basic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithCLNodes(1).
		WithFunding(vrfConst.ChainlinkNodeFundingAmountEth).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)
	env.ParallelTransactions(true)

	mockFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, vrfConst.LinkEthFeedResponse)
	require.NoError(t, err)
	lt, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)
	vrfv2Contracts, err := vrfv2_actions.DeployVRFV2Contracts(env.ContractDeployer, env.EVMClient, lt, mockFeed)
	require.NoError(t, err)

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err)

	err = vrfv2Contracts.Coordinator.SetConfig(
		vrfConst.MinimumConfirmations,
		vrfConst.MaxGasLimitVRFCoordinatorConfig,
		vrfConst.StalenessSeconds,
		vrfConst.GasAfterPaymentCalculation,
		vrfConst.LinkEthFeedResponse,
		vrfConst.VRFCoordinatorV2FeeConfig,
	)
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err)

	err = vrfv2Contracts.Coordinator.CreateSubscription()
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err)

	err = vrfv2Contracts.Coordinator.AddConsumer(vrfConst.SubID, vrfv2Contracts.LoadTestConsumer.Address())
	require.NoError(t, err)

	err = vrfv2_actions.FundVRFCoordinatorV2Subscription(lt, vrfv2Contracts.Coordinator, env.EVMClient, vrfConst.SubID, vrfConst.VRFSubscriptionFundingAmountLink)
	require.NoError(t, err)

	vrfV2jobs, err := vrfv2_actions.CreateVRFV2Jobs(env.ClCluster.NodeAPIs(), vrfv2Contracts.Coordinator, env.EVMClient, vrfConst.MinimumConfirmations)
	require.NoError(t, err)

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	addr, err := env.ClCluster.Nodes[0].API.PrimaryEthAddress()
	require.NoError(t, err)
	nodeConfig := node.NewConfig(env.ClCluster.Nodes[0].NodeConfig,
		node.WithVRFv2EVMEstimator(addr),
	)
	err = env.ClCluster.Nodes[0].Restart(nodeConfig)
	require.NoError(t, err)

	// test and assert
	err = vrfv2Contracts.LoadTestConsumer.RequestRandomness(
		vrfV2jobs[0].KeyHash,
		vrfConst.SubID,
		vrfConst.MinimumConfirmations,
		vrfConst.CallbackGasLimit,
		vrfConst.NumberOfWords,
		vrfConst.RandomnessRequestCountPerRequest,
	)
	require.NoError(t, err)

	gom := gomega.NewGomegaWithT(t)
	timeout := time.Minute * 2
	var lastRequestID *big.Int
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.ClCluster.Nodes[0].API.MustReadRunsByJob(vrfV2jobs[0].Job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
		lastRequestID, err = vrfv2Contracts.LoadTestConsumer.GetLastRequestId(context.Background())
		l.Debug().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		status, err := vrfv2Contracts.LoadTestConsumer.GetRequestStatus(context.Background(), lastRequestID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(status.Fulfilled).Should(gomega.BeTrue())
		l.Debug().Interface("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		for _, w := range status.RandomWords {
			l.Info().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
			g.Expect(w.Uint64()).Should(gomega.BeNumerically(">", 0), "Expected the VRF job give an answer bigger than 0")
		}
	}, timeout, "1s").Should(gomega.Succeed())
}
