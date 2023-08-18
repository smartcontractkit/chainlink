package smoke

import (
	"context"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_constants"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

func TestVRFv2PlusBasic(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithCLNodes(1).
		WithFunding(vrfv2plus_constants.ChainlinkNodeFundingAmountEth).
		Build()
	require.NoError(t, err)
	env.ParallelTransactions(true)

	mockETHLinkFeedAddress, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, vrfv2plus_constants.LinkEthFeedResponse)
	require.NoError(t, err)
	linkAddress, err := actions.DeployLINKToken(env.ContractDeployer)
	require.NoError(t, err)
	vrfv2PlusContracts, err := vrfv2plus.DeployVRFV2PlusContracts(env.ContractDeployer, env.EVMClient)
	require.NoError(t, err)

	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err)

	err = vrfv2PlusContracts.Coordinator.SetLINKAndLINKETHFeed(linkAddress.Address(), mockETHLinkFeedAddress.Address())
	require.NoError(t, err)

	err = vrfv2PlusContracts.Coordinator.SetConfig(
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.MaxGasLimitVRFCoordinatorConfig,
		vrfv2plus_constants.StalenessSeconds,
		vrfv2plus_constants.GasAfterPaymentCalculation,
		vrfv2plus_constants.LinkEthFeedResponse,
		vrfv2plus_constants.VRFCoordinatorV2PlusFeeConfig,
	)
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err)

	err = vrfv2PlusContracts.Coordinator.CreateSubscription()
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err)

	subID, err := vrfv2PlusContracts.Coordinator.FindSubscriptionID()
	require.NoError(t, err)

	err = vrfv2PlusContracts.Coordinator.AddConsumer(subID, vrfv2PlusContracts.LoadTestConsumer.Address())
	require.NoError(t, err)

	err = vrfv2plus.FundVRFCoordinatorV2PlusSubscription(linkAddress, vrfv2PlusContracts.Coordinator, env.EVMClient, subID, vrfv2plus_constants.VRFSubscriptionFundingAmountLink)
	require.NoError(t, err)

	vrfV2jobs, err := vrfv2plus.CreateVRFV2PlusJobs(env.GetAPIs(), vrfv2PlusContracts.Coordinator, env.EVMClient, vrfv2plus_constants.MinimumConfirmations)
	require.NoError(t, err)

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	addr, err := env.CLNodes[0].API.PrimaryEthAddress()
	require.NoError(t, err)
	nodeConfig := node.NewConfig(env.CLNodes[0].NodeConfig,
		node.WithVRFv2EVMEstimator(addr),
	)
	err = env.CLNodes[0].Restart(nodeConfig)
	require.NoError(t, err)

	// test and assert
	err = vrfv2PlusContracts.LoadTestConsumer.RequestRandomness(
		vrfV2jobs[0].KeyHash,
		subID,
		vrfv2plus_constants.MinimumConfirmations,
		vrfv2plus_constants.CallbackGasLimit,
		false,
		vrfv2plus_constants.NumberOfWords,
		vrfv2plus_constants.RandomnessRequestCountPerRequest,
	)
	require.NoError(t, err)

	gom := gomega.NewGomegaWithT(t)
	timeout := time.Minute * 1
	var lastRequestID *big.Int
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(vrfV2jobs[0].Job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
		lastRequestID, err = vrfv2PlusContracts.LoadTestConsumer.GetLastRequestId(context.Background())
		l.Debug().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		status, err := vrfv2PlusContracts.LoadTestConsumer.GetRequestStatus(context.Background(), lastRequestID)
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
