package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
)

var (
	LinkEthFeedResponse              = big.NewInt(1e18)
	MinimumConfirmations             = uint16(3)
	RandomnessRequestCountPerRequest = uint16(1)
	//todo - get Sub id when creating subscription - need to listen for SubscriptionCreated Log
	SubID                            = uint64(1)
	VRFSubscriptionFundingAmountLink = big.NewInt(100)
	ChainlinkNodeFundingAmountEth    = big.NewFloat(1)
	NumberOfWords                    = uint32(3)
	MaxGasPriceGWei                  = 1000
	CallbackGasLimit                 = uint32(1000000)
	MaxGasLimitVRFCoordinatorConfig  = uint32(2.5e6)
	StalenessSeconds                 = uint32(86400)
	GasAfterPaymentCalculation       = uint32(33825)

	VRFCoordinatorV2FeeConfig = vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: 500,
		FulfillmentFlatFeeLinkPPMTier2: 500,
		FulfillmentFlatFeeLinkPPMTier3: 500,
		FulfillmentFlatFeeLinkPPMTier4: 500,
		FulfillmentFlatFeeLinkPPMTier5: 500,
		ReqsForTier2:                   big.NewInt(0),
		ReqsForTier3:                   big.NewInt(0),
		ReqsForTier4:                   big.NewInt(0),
		ReqsForTier5:                   big.NewInt(0),
	}
)

func TestVRFv2Basic(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithMockServer(1).
		WithCLNodes(1).
		WithFunding(big.NewFloat(1)).
		Build()
	require.NoError(t, err)
	env.ParallelTransactions(true)

	err = env.DeployMockETHLinkFeed(LinkEthFeedResponse)
	require.NoError(t, err)
	err = env.DeployLINKToken()
	require.NoError(t, err)
	err = env.DeployVRFV2Contracts()
	require.NoError(t, err)

	err = env.WaitForEvents()
	require.NoError(t, err)

	err = env.CoordinatorV2.SetConfig(
		MinimumConfirmations,
		MaxGasLimitVRFCoordinatorConfig,
		StalenessSeconds,
		GasAfterPaymentCalculation,
		LinkEthFeedResponse,
		VRFCoordinatorV2FeeConfig,
	)
	require.NoError(t, err)
	err = env.WaitForEvents()
	require.NoError(t, err)

	err = env.CoordinatorV2.CreateSubscription()
	require.NoError(t, err)
	err = env.WaitForEvents()
	require.NoError(t, err)

	err = env.CoordinatorV2.AddConsumer(SubID, env.LoadTestConsumer.Address())
	require.NoError(t, err)

	err = env.FundVRFCoordinatorV2Subscription(SubID, VRFSubscriptionFundingAmountLink)
	require.NoError(t, err)

	vrfV2jobs, err := env.CreateVRFv2Jobs(env.CoordinatorV2)
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
	err = env.LoadTestConsumer.RequestRandomness(
		vrfV2jobs[0].KeyHash,
		SubID,
		MinimumConfirmations,
		CallbackGasLimit,
		NumberOfWords,
		RandomnessRequestCountPerRequest,
	)
	require.NoError(t, err)

	gom := gomega.NewGomegaWithT(t)
	timeout := time.Minute * 1
	var lastRequestID *big.Int
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := env.CLNodes[0].API.MustReadRunsByJob(vrfV2jobs[0].Job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
		lastRequestID, err = env.LoadTestConsumer.GetLastRequestId(context.Background())
		l.Debug().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")

		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		status, err := env.LoadTestConsumer.GetRequestStatus(context.Background(), lastRequestID)
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
