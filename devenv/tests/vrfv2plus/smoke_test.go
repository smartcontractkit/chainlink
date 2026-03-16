package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFv2PlusSmoke(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrf2plus-out.toml"

	// Load environment outputs
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load env-out.toml")

	cfg, err := products.LoadOutput[vrfv2plus.Configurator](outputFile)
	require.NoError(t, err, "failed to load VRFv2Plus config from env-out.toml")
	require.NotEmpty(t, cfg.Config, "vrfv2_plus config must not be empty")

	c := cfg.Config[0]
	require.NotEmpty(t, c.DeployedContracts.Coordinator, "coordinator address must not be empty")
	require.NotEmpty(t, c.VRFKeyData.KeyHash, "key hash must not be empty")
	require.NotEmpty(t, c.VRFKeyData.VRFJobID, "VRF job ID must not be empty")

	// Parse key hash
	keyHashBytes := common.HexToHash(c.VRFKeyData.KeyHash)
	var keyHash [32]byte
	copy(keyHash[:], keyHashBytes[:])

	// Set up Seth chain client
	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err)
	bcNode := in.Blockchains[0].Out.Nodes[0]
	ctx := t.Context()
	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	// Load coordinator
	coord, err := contracts.LoadVRFCoordinatorV2_5(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load coordinator")

	// Load LINK token
	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	// Connect to CL nodes (used in Job Runs subtest)
	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to connect to CL nodes")

	runWithFunding := func(name string, fn func(t *testing.T)) {
		t.Run(name, func(t *testing.T) {
			reconcileConfiguredFunding(t, ctx, chainClient, coord, linkToken, c)
			fn(t)
		})
	}

	runWithFunding("Link Billing", func(t *testing.T) {
		consumer, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)

		fulfilled := requestAndWait(t, ctx, consumer, coord, keyHash, subID, false,
			c.MinimumConfirmations, defaultFulfillTimeout)

		require.True(t, fulfilled.Success, "RandomWordsFulfilled.Success should be true")
		require.False(t, fulfilled.NativePayment, "should be LINK payment")

		status, sErr := consumer.GetRequestStatus(ctx, fulfilled.RequestId)
		require.NoError(t, sErr)
		require.Len(t, status.RandomWords, 1, "expected 1 random word")
		for _, w := range status.RandomWords {
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "random word should be > 0")
		}
	})

	runWithFunding("Native Billing", func(t *testing.T) {
		consumer, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)

		fulfilled := requestAndWait(t, ctx, consumer, coord, keyHash, subID, true,
			c.MinimumConfirmations, defaultFulfillTimeout)

		require.True(t, fulfilled.Success, "RandomWordsFulfilled.Success should be true")
		require.True(t, fulfilled.NativePayment, "should be native payment")

		status, sErr := consumer.GetRequestStatus(ctx, fulfilled.RequestId)
		require.NoError(t, sErr)
		require.Len(t, status.RandomWords, 1, "expected 1 random word")
		for _, w := range status.RandomWords {
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "random word should be > 0")
		}
	})

	runWithFunding("Direct Funding — Link Billing", func(t *testing.T) {
		wrapperConsumer, err := contracts.LoadVRFV2PlusWrapperLoadTestConsumer(
			chainClient, c.DeployedContracts.WrapperConsumer)
		require.NoError(t, err)

		fulfilled := requestAndWaitWrapper(t, ctx, wrapperConsumer, coord, false,
			c.MinimumConfirmations, defaultFulfillTimeout)

		require.True(t, fulfilled.Success, "RandomWordsFulfilled.Success should be true")
		require.False(t, fulfilled.NativePayment, "should be LINK payment")

		status, sErr := wrapperConsumer.GetRequestStatus(ctx, fulfilled.RequestId)
		require.NoError(t, sErr)
		require.Len(t, status.RandomWords, 1, "expected 1 random word")
	})

	runWithFunding("Direct Funding — Native Billing", func(t *testing.T) {
		wrapperConsumer, err := contracts.LoadVRFV2PlusWrapperLoadTestConsumer(
			chainClient, c.DeployedContracts.WrapperConsumer)
		require.NoError(t, err)

		fulfilled := requestAndWaitWrapper(t, ctx, wrapperConsumer, coord, true,
			c.MinimumConfirmations, defaultFulfillTimeout)

		require.True(t, fulfilled.Success, "RandomWordsFulfilled.Success should be true")
		require.True(t, fulfilled.NativePayment, "should be native payment")

		status, sErr := wrapperConsumer.GetRequestStatus(ctx, fulfilled.RequestId)
		require.NoError(t, sErr)
		require.Len(t, status.RandomWords, 1, "expected 1 random word")
	})

	runWithFunding("Block Confirmation", func(t *testing.T) {
		consumer, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)

		const highConfirmations = uint16(10)
		fulfilled := requestAndWait(t, ctx, consumer, coord, keyHash, subID, false,
			highConfirmations, defaultFulfillTimeout)

		require.True(t, fulfilled.Success)
		// Verify the fulfillment happened (event existence proves confirmations were waited)
		require.NotNil(t, fulfilled.RequestId)
	})

	runWithFunding("Job Runs", func(t *testing.T) {
		runsBefore, rErr := cl[0].MustReadRunsByJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, rErr)
		beforeCount := len(runsBefore.Data)

		consumer, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)
		fulfilled := requestAndWait(t, ctx, consumer, coord, keyHash, subID, false,
			c.MinimumConfirmations, defaultFulfillTimeout)
		require.True(t, fulfilled.Success)

		runsAfter, rErr := cl[0].MustReadRunsByJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, rErr)
		require.Greater(t, len(runsAfter.Data), beforeCount,
			"job run count should increase after fulfillment")
	})

	runWithFunding("Cancel Sub", func(t *testing.T) {
		_, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)

		// Verify sub exists
		sub, sErr := coord.GetSubscription(ctx, subID)
		require.NoError(t, sErr)
		require.NotNil(t, sub.Owner)

		// Cancel subscription and return funds to our address
		owner := chainClient.MustGetRootKeyAddress()
		err = coord.CancelSubscription(subID, owner)
		require.NoError(t, err, "CancelSubscription should succeed")
	})

	runWithFunding("Owner Cancel", func(t *testing.T) {
		_, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)

		err = coord.OwnerCancelSubscription(subID)
		require.NoError(t, err, "OwnerCancelSubscription should succeed")
	})

	runWithFunding("Owner Withdraw", func(t *testing.T) {
		consumer, subID := newConsumerAndSub(t, ctx, chainClient, coord, linkToken, c)
		fulfilled := requestAndWait(t, ctx, consumer, coord, keyHash, subID, false,
			c.MinimumConfirmations, defaultFulfillTimeout)
		require.True(t, fulfilled.Success)

		owner := chainClient.MustGetRootKeyAddress()
		err = coord.Withdraw(owner)
		require.NoError(t, err, "Withdraw should succeed")

		err = coord.WithdrawNative(owner)
		require.NoError(t, err, "WithdrawNative should succeed")
	})
}

// newConsumerAndSub deploys a fresh load test consumer, creates and funds a subscription,
// and adds the consumer to it. Returns the consumer and the subscription ID.
func newConsumerAndSub(
	t *testing.T,
	ctx context.Context,
	chainClient *seth.Client,
	coord *contracts.EthereumVRFCoordinatorV2_5,
	linkToken contracts.LinkToken,
	c *vrfv2plus.VRFv2Plus,
) (*contracts.EthereumVRFv2PlusLoadTestConsumer, *big.Int) {
	t.Helper()

	consumer, err := contracts.DeployVRFv2PlusLoadTestConsumer(chainClient, coord.Address())
	require.NoError(t, err, "failed to deploy load test consumer")

	subID, err := createAndFundSub(ctx, chainClient, coord, linkToken,
		c.SubFundingAmountLink, c.SubFundingAmountNative)
	require.NoError(t, err, "failed to create and fund subscription")

	err = coord.AddConsumer(subID, consumer.Address())
	require.NoError(t, err, "failed to add consumer to subscription")

	return consumer, subID
}
