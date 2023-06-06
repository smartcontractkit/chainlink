package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
)

func TestVRFv2Basic(t *testing.T) {

	t.Parallel()
	l := utils.GetTestLogger(t)

	testNetwork := networks.SelectedNetwork
	testEnvironment := vrfv2_actions.SetupVRFV2Environment(
		t,
		testNetwork,
		config.BaseVRFV2NetworkDetailTomlConfig,
		"",
		"smoke-vrfv2",
		"",
		time.Minute*20,
	)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err)
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err)
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err)

	chainClient.ParallelTransactions(true)

	mockETHLINKFeed, err := contractDeployer.DeployMockETHLINKFeed(vrfv2_constants.LinkEthFeedResponse)
	require.NoError(t, err)
	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err)

	vrfV2Contracts, chainlinkNodesAfterRedeployment, vrfV2jobs, testEnvironmentAfterRedeployment := vrfv2_actions.SetupVRFV2Universe(
		t,
		linkToken,
		mockETHLINKFeed,
		contractDeployer,
		chainClient,
		chainlinkNodes,
		testNetwork,
		testEnvironment,
		vrfv2_constants.ChainlinkNodeFundingAmountEth,
		vrfv2_constants.VRFSubscriptionFundingAmountLink,
		"smoke-vrfv2",
		time.Minute*20,
	)

	consumerContract := vrfV2Contracts.LoadTestConsumer

	t.Cleanup(func() {
		err := actions.TeardownSuite(
			t,
			testEnvironmentAfterRedeployment,
			utils.ProjectRoot,
			chainlinkNodesAfterRedeployment,
			nil,
			zapcore.ErrorLevel,
			chainClient,
		)
		require.NoError(t, err, "Error tearing down environment")
	})

	err = consumerContract.RequestRandomness(
		vrfV2jobs[0].KeyHash,
		vrfv2_constants.SubID,
		vrfv2_constants.MinimumConfirmations,
		vrfv2_constants.CallbackGasLimit,
		vrfv2_constants.NumberOfWords,
		vrfv2_constants.RandomnessRequestCountPerRequest,
	)
	require.NoError(t, err)

	gom := gomega.NewGomegaWithT(t)
	timeout := time.Minute * 2
	var lastRequestID *big.Int
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := chainlinkNodesAfterRedeployment[0].MustReadRunsByJob(vrfV2jobs[0].Job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
		lastRequestID, err = consumerContract.GetLastRequestId(context.Background())
		l.Debug().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		status, err := consumerContract.GetRequestStatus(context.Background(), lastRequestID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(status.Fulfilled).Should(gomega.BeTrue())
		l.Debug().Interface("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		for _, w := range status.RandomWords {
			l.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
			g.Expect(w.Uint64()).Should(gomega.BeNumerically(">", 0), "Expected the VRF job give an answer bigger than 0")
		}
	}, timeout, "1s").Should(gomega.Succeed())
}
