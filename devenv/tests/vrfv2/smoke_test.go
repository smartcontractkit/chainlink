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
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrfv2-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)

	cfg, err := products.LoadOutput[productvrfv2.Configurator](outputFile)
	require.NoError(t, err)
	require.NotEmpty(t, cfg.Config)
	c := cfg.Config[0]

	// Unlike vrfv2plus smoke (reconcileConfiguredFunding), we do not top up between
	// subtests: each subtest deploys fresh load-test consumers and funded subs; Direct
	// Funding is self-contained. Node TX keys are funded once in ConfigureJobsAndContracts.

	keyHash := mustKeyHash(c)

	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err)
	bcNode := in.Blockchains[0].Out.Nodes[0]
	ctx := t.Context()
	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err)

	coord, err := contracts.LoadVRFCoordinatorV2(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err)

	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err)

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	fulfillTimeout := parseFulfillTimeout(c.RandomWordsFulfilledEventTimeout)

	t.Run("Request Randomness", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]
		sub, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err)
		balBefore := new(big.Int).Set(sub.Balance)

		_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err)

		expectedBal := new(big.Int).Sub(balBefore, fulfilled.Payment)
		subAfter, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err)
		require.Equal(t, 0, expectedBal.Cmp(subAfter.Balance))

		status, err := consumers[0].GetRequestStatus(ctx, fulfilled.RequestID)
		require.NoError(t, err)
		require.True(t, status.Fulfilled)
		require.Len(t, status.RandomWords, int(c.NumberOfWords))
		for _, w := range status.RandomWords {
			require.Equal(t, 1, w.Cmp(big.NewInt(0)))
		}
	})

	t.Run("VRF Node waits block confirmation number specified by the consumer before sending fulfilment on-chain", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		const expectedBlockWait = uint16(10)
		req, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			expectedBlockWait, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err)
		require.GreaterOrEqual(t, fulfilled.Raw.BlockNumber, req.Raw.BlockNumber+uint64(expectedBlockWait))
	})

	t.Run("CL Node VRF Job Runs", func(t *testing.T) {
		runsBefore, err := cl[0].MustReadRunsByJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, err)
		beforeN := len(runsBefore.Data)

		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		_, _, err = requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err)

		runsAfter, err := cl[0].MustReadRunsByJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, err)
		require.Len(t, runsAfter.Data, beforeN+1)
	})

	t.Run("Direct Funding", func(t *testing.T) {
		wrapper, err := contracts.DeployVRFV2Wrapper(chainClient,
			c.DeployedContracts.LinkToken, c.DeployedContracts.MockFeed, c.DeployedContracts.Coordinator)
		require.NoError(t, err)
		err = wrapper.SetConfig(c.WrapperGasOverhead, c.CoordinatorGasOverhead, c.WrapperPremiumPercentage, keyHash, c.WrapperMaxNumberOfWords)
		require.NoError(t, err)
		wrapperSubID, err := wrapper.GetSubID(ctx)
		require.NoError(t, err)

		// Fund the wrapper's coordinator subscription with LINK via the coordinator
		// (same as integration-tests FundSubscriptionWithLink — not TransferAndCall to the wrapper).
		amount := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
		enc, err := utilsABIEncodeUint64(wrapperSubID)
		require.NoError(t, err)
		_, err = linkToken.TransferAndCall(coord.Address(), amount, enc)
		require.NoError(t, err)

		wConsumer, err := contracts.DeployVRFV2WrapperLoadTestConsumer(chainClient, c.DeployedContracts.LinkToken, wrapper.Address())
		require.NoError(t, err)
		fundJuels := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(int64(c.WrapperConsumerFundingAmountLink)))
		err = linkToken.Transfer(wConsumer.Address(), fundJuels)
		require.NoError(t, err)

		balBefore, err := linkToken.BalanceOf(ctx, wConsumer.Address())
		require.NoError(t, err)
		wSub, err := coord.GetSubscription(ctx, wrapperSubID)
		require.NoError(t, err)
		subBalBefore := new(big.Int).Set(wSub.Balance)

		fulfilled, err := directFundingRequestAndWait(ctx, wConsumer, coord, wrapperSubID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, fulfillTimeout)
		require.NoError(t, err)

		expSub := new(big.Int).Sub(subBalBefore, fulfilled.Payment)
		wSubAfter, err := coord.GetSubscription(ctx, wrapperSubID)
		require.NoError(t, err)
		require.Equal(t, 0, expSub.Cmp(wSubAfter.Balance))

		consStatus, err := wConsumer.GetRequestStatus(ctx, fulfilled.RequestID)
		require.NoError(t, err)
		require.True(t, consStatus.Fulfilled)
		balAfter, err := linkToken.BalanceOf(ctx, wConsumer.Address())
		require.NoError(t, err)
		expConsumerBal := new(big.Int).Sub(balBefore, consStatus.Paid)
		require.Equal(t, 0, expConsumerBal.Cmp(balAfter))
		require.Len(t, consStatus.RandomWords, int(c.NumberOfWords))
		for _, w := range consStatus.RandomWords {
			require.Equal(t, 1, w.Cmp(big.NewInt(0)))
		}
		t.Logf("Consumer balance before %s after %s paid %s",
			(*commonassets.Link)(balBefore).Link(), (*commonassets.Link)(balAfter).Link(), (*commonassets.Link)(consStatus.Paid).Link())
	})

	t.Run("Oracle Withdraw", func(t *testing.T) {
		tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-527")

		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err)

		root := chainClient.MustGetRootKeyAddress()
		balBefore, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err)
		err = coord.OracleWithdraw(root, fulfilled.Payment)
		require.NoError(t, err)
		balAfter, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err)
		require.Equal(t, 1, balAfter.Cmp(balBefore))
	})

	t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
		_, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		toAddr, err := randomWalletAddress()
		require.NoError(t, err)
		balBefore, err := linkToken.BalanceOf(ctx, toAddr.Hex())
		require.NoError(t, err)

		sub, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err)
		subBal := new(big.Int).Set(sub.Balance)

		_, cancelEv, err := coord.CancelSubscription(subID, toAddr)
		require.NoError(t, err)
		require.Equal(t, 0, subBal.Cmp(cancelEv.Amount))

		_, err = coord.GetSubscription(ctx, subID)
		require.Error(t, err)

		balAfter, err := linkToken.BalanceOf(ctx, toAddr.Hex())
		require.NoError(t, err)
		returned := new(big.Int).Sub(balAfter, balBefore)
		require.Equal(t, 0, subBal.Cmp(returned))
	})

	t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, 0, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		pending, err := coord.PendingRequestsExist(ctx, subID)
		require.NoError(t, err)
		require.False(t, pending)

		_, _, err = requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
			5*time.Second, 0)
		require.Error(t, err)

		pending, err = coord.PendingRequestsExist(ctx, subID)
		require.NoError(t, err)
		require.True(t, pending)

		root := chainClient.MustGetRootKeyAddress()
		wBalBefore, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err)
		sub, err := coord.GetSubscription(ctx, subID)
		require.NoError(t, err)
		subBal := new(big.Int).Set(sub.Balance)

		_, cancelEv, err := coord.OwnerCancelSubscription(subID)
		require.NoError(t, err)
		require.Equal(t, 0, subBal.Cmp(cancelEv.Amount))

		_, err = coord.GetSubscription(ctx, subID)
		require.Error(t, err)

		wBalAfter, err := linkToken.BalanceOf(ctx, root.Hex())
		require.NoError(t, err)
		returned := new(big.Int).Sub(wBalAfter, wBalBefore)
		require.Equal(t, 0, subBal.Cmp(returned))
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
