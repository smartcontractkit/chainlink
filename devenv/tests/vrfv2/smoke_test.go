package vrfv2

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2 "github.com/smartcontractkit/chainlink/devenv/products/vrfv2"
)

func TestVRFv2Basic(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr, "failed to save container logs")
	})

	outputFile := "../../env-vrfv2-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load devenv env-out from %s", outputFile)

	cfg, err := products.LoadOutput[productvrfv2.Configurator](outputFile)
	require.NoError(t, err, "failed to load vrfv2 product config from env-out")
	require.NotEmpty(t, cfg.Config, "vrfv2 config must not be empty in env-out")
	c := cfg.Config[0]

	// Unlike vrfv2plus smoke (reconcileConfiguredFunding), we do not top up between
	// subtests: each subtest deploys fresh load-test consumers and funded subs; Direct
	// Funding is self-contained. Node TX keys are funded once in ConfigureJobsAndContracts.

	keyHash := mustKeyHash(c)

	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err, "failed to parse chain ID from env-out")
	bcNode := in.Blockchains[0].Out.Nodes[0]
	ctx := t.Context()
	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	coord, err := contracts.LoadVRFCoordinatorV2(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load VRF coordinator v2")

	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to connect to Chainlink nodes")

	fulfillTimeout := parseFulfillTimeout(c.RandomWordsFulfilledEventTimeout)

	t.Run("Request Randomness", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]
		sub, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err, "error getting subscription information before request")
		balBefore := new(big.Int).Set(sub.Balance)

		_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		expectedBal := new(big.Int).Sub(balBefore, fulfilled.Payment)
		subAfter, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err, "error getting subscription after fulfillment")
		require.Equal(t, 0, expectedBal.Cmp(subAfter.Balance),
			"subscription LINK balance should equal pre-fulfillment balance minus payment (expected %s got %s)", expectedBal.String(), subAfter.Balance.String())

		status, err := consumers[0].GetRequestStatus(ctx, fulfilled.RequestID)
		require.NoError(t, err, "error getting randomness request status")
		require.True(t, status.Fulfilled, "random words request should be fulfilled")
		require.Len(t, status.RandomWords, int(c.NumberOfWords), "wrong number of random words in consumer status")
		for i, w := range status.RandomWords {
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "random word %d should be non-zero", i)
		}
	})

	t.Run("VRF Node waits block confirmation number specified by the consumer before sending fulfilment on-chain", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		const expectedBlockWait = uint16(10)
		req, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			expectedBlockWait, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
		require.GreaterOrEqual(t, fulfilled.Raw.BlockNumber, req.Raw.BlockNumber+uint64(expectedBlockWait),
			"fulfillment block should be at least request block + minimum confirmations (req=%d fulfilled=%d minConf=%d)",
			req.Raw.BlockNumber, fulfilled.Raw.BlockNumber, expectedBlockWait)
	})

	t.Run("CL Node VRF Job Runs", func(t *testing.T) {
		runsBefore, err := cl[0].MustReadRunsByJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, err, "failed to read VRF job runs before request")
		beforeN := len(runsBefore.Data)

		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		_, _, err = requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		runsAfter, err := cl[0].MustReadRunsByJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, err, "failed to read VRF job runs after request")
		require.Len(t, runsAfter.Data, beforeN+1, "VRF job should have one new run after fulfillment (before=%d after=%d)", beforeN, len(runsAfter.Data))
	})

	t.Run("Direct Funding", func(t *testing.T) {
		wrapper, err := contracts.DeployVRFV2Wrapper(chainClient,
			c.DeployedContracts.LinkToken, c.DeployedContracts.MockFeed, c.DeployedContracts.Coordinator)
		require.NoError(t, err, "failed to deploy VRF v2 wrapper")
		err = wrapper.SetConfig(c.WrapperGasOverhead, c.CoordinatorGasOverhead, c.WrapperPremiumPercentage, keyHash, c.WrapperMaxNumberOfWords)
		require.NoError(t, err, "failed to set wrapper config")
		wrapperSubID, err := wrapper.GetSubID(ctx)
		require.NoError(t, err, "failed to read wrapper subscription ID")

		// Fund the wrapper's coordinator subscription with LINK via the coordinator
		// (same as integration-tests FundSubscriptionWithLink — not TransferAndCall to the wrapper).
		amount := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
		enc, err := utilsABIEncodeUint64(wrapperSubID)
		require.NoError(t, err, "failed to ABI-encode wrapper sub ID for funding")
		_, err = linkToken.TransferAndCall(coord.Address(), amount, enc)
		require.NoError(t, err, "failed to fund wrapper subscription with LINK via coordinator")

		wConsumer, err := contracts.DeployVRFV2WrapperLoadTestConsumer(chainClient, c.DeployedContracts.LinkToken, wrapper.Address())
		require.NoError(t, err, "failed to deploy wrapper load test consumer")
		fundJuels := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(int64(c.WrapperConsumerFundingAmountLink)))
		err = linkToken.Transfer(wConsumer.Address(), fundJuels)
		require.NoError(t, err, "failed to fund wrapper consumer with LINK")

		balBefore, err := linkToken.BalanceOf(ctx, wConsumer.Address())
		require.NoError(t, err, "failed to read wrapper consumer LINK balance before request")
		wSub, err := coord.GetSubscription(ctx, wrapperSubID)
		require.NoError(t, err, "failed to get wrapper subscription before request")
		subBalBefore := new(big.Int).Set(wSub.Balance)

		fulfilled, err := directFundingRequestAndWait(ctx, wConsumer, coord, wrapperSubID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, fulfillTimeout)
		require.NoError(t, err, "direct funding request did not fulfill in time")

		expSub := new(big.Int).Sub(subBalBefore, fulfilled.Payment)
		wSubAfter, err := coord.GetSubscription(ctx, wrapperSubID)
		require.NoError(t, err, "failed to get wrapper subscription after fulfillment")
		require.Equal(t, 0, expSub.Cmp(wSubAfter.Balance),
			"wrapper sub LINK balance should equal pre-fulfillment minus payment (expected %s got %s)", expSub.String(), wSubAfter.Balance.String())

		consStatus, err := wConsumer.GetRequestStatus(ctx, fulfilled.RequestID)
		require.NoError(t, err, "error getting wrapper consumer request status")
		require.True(t, consStatus.Fulfilled, "direct funding request should be fulfilled")
		balAfter, err := linkToken.BalanceOf(ctx, wConsumer.Address())
		require.NoError(t, err, "failed to read wrapper consumer LINK balance after request")
		expConsumerBal := new(big.Int).Sub(balBefore, consStatus.Paid)
		require.Equal(t, 0, expConsumerBal.Cmp(balAfter),
			"consumer LINK balance should equal pre-request minus paid amount (expected %s got %s)", expConsumerBal.String(), balAfter.String())
		require.Len(t, consStatus.RandomWords, int(c.NumberOfWords), "wrong number of random words from direct funding")
		for i, w := range consStatus.RandomWords {
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "random word %d should be non-zero", i)
		}
		t.Logf("Consumer balance before %s after %s paid %s",
			(*commonassets.Link)(balBefore).Link(), (*commonassets.Link)(balAfter).Link(), (*commonassets.Link)(consStatus.Paid).Link())
	})

	t.Run("Oracle Withdraw", func(t *testing.T) {
		tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-527")

		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		root := chainClient.MustGetRootKeyAddress()
		balBefore, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err, "failed to read oracle LINK balance before withdraw")
		err = coord.OracleWithdraw(root, fulfilled.Payment)
		require.NoError(t, err, "oracle withdraw failed")
		balAfter, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err, "failed to read oracle LINK balance after withdraw")
		require.Equal(t, 1, balAfter.Cmp(balBefore), "oracle LINK balance should increase after withdraw (before=%s after=%s)", balBefore.String(), balAfter.String())
	})

	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		_, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		toAddr, err := randomWalletAddress()
		require.NoError(t, err, "failed to generate cancel recipient address")
		balBefore, err := linkToken.BalanceOf(ctx, toAddr.Hex())
		require.NoError(t, err, "failed to read recipient LINK balance before cancel")

		sub, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err, "error getting subscription before cancel")
		subBal := new(big.Int).Set(sub.Balance)
		require.Equal(t, 1, subBal.Sign(), "subscription should be funded before cancel (expected positive balance from SubFundingAmountLink); got %s", subBal.String())

		_, cancelEv, err := coord.CancelSubscription(subID, toAddr)
		require.NoError(t, err, "cancel subscription failed")
		require.Equal(t, 0, subBal.Cmp(cancelEv.Amount),
			"SubscriptionCanceled amount should match sub balance (subBal=%s canceled=%s)", subBal.String(), cancelEv.Amount.String())

		_, err = coord.GetSubscription(ctx, subID)
		require.Error(t, err, "get subscription should fail after cancel")

		balAfter, err := linkToken.BalanceOf(ctx, toAddr.Hex())
		require.NoError(t, err, "failed to read recipient LINK balance after cancel")
		returned := new(big.Int).Sub(balAfter, balBefore)
		require.Equal(t, 0, subBal.Cmp(returned),
			"recipient should receive full sub LINK balance (expected return=%s got=%s)", subBal.String(), returned.String())
	})

	t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, 0, 1, 1)
		require.NoError(t, err, "error setting up unfunded consumer and sub")
		subID := subIDs[0]

		pending, err := coord.PendingRequestsExist(ctx, subID)
		require.NoError(t, err, "failed to check pending requests before stuck request")
		require.False(t, pending, "should have no pending requests before underfunded request")

		_, _, err = requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			5*time.Second, 0)
		require.Error(t, err, "underfunded request should not fulfill within short timeout")

		pending, err = coord.PendingRequestsExist(ctx, subID)
		require.NoError(t, err, "failed to check pending requests after stuck request")
		require.True(t, pending, "pending request should exist before owner cancel")

		root := chainClient.MustGetRootKeyAddress()
		wBalBefore, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err, "failed to read owner LINK balance before owner cancel")
		sub, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err, "error getting subscription before owner cancel")
		subBal := new(big.Int).Set(sub.Balance)

		_, cancelEv, err := coord.OwnerCancelSubscription(subID)
		require.NoError(t, err, "owner cancel subscription failed")
		require.Equal(t, 0, subBal.Cmp(cancelEv.Amount),
			"owner cancel amount should match sub balance (subBal=%s canceled=%s)", subBal.String(), cancelEv.Amount.String())

		_, err = coord.GetSubscription(ctx, subID)
		require.Error(t, err, "get subscription should fail after owner cancel")

		wBalAfter, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err, "failed to read owner LINK balance after owner cancel")
		returned := new(big.Int).Sub(wBalAfter, wBalBefore)
		require.Equal(t, 0, subBal.Cmp(returned),
			"owner should receive full sub LINK balance (expected=%s got=%s)", subBal.String(), returned.String())
	})
}

func utilsABIEncodeUint64(subID uint64) ([]byte, error) {
	return utils.ABIEncode(`[{"type":"uint64"}]`, subID)
}

func randomWalletAddress() (common.Address, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(key.PublicKey), nil
}
