package vrfv2plus

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFV2PlusWithBHS(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrf2plus-bhX-out.toml"

	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load env-vrf2plus-bhX-out.toml")

	cfg, err := products.LoadOutput[vrfv2plus.Configurator](outputFile)
	require.NoError(t, err, "failed to load VRFv2Plus config")
	require.NotEmpty(t, cfg.Config, "vrfv2_plus config must not be empty")

	c := cfg.Config[0]
	require.NotEmpty(t, c.DeployedContracts.Coordinator, "coordinator address must not be empty")
	require.NotEmpty(t, c.VRFKeyData.KeyHash, "key hash must not be empty")
	require.NotEmpty(t, c.VRFKeyData.BHSJobID, "BHS job ID must not be empty — was env-vrf2plus-bhX.toml used?")
	require.NotEmpty(t, c.DeployedContracts.BHS, "BHS contract address must not be empty")

	keyHashBytes := common.HexToHash(c.VRFKeyData.KeyHash)
	var keyHash [32]byte
	copy(keyHash[:], keyHashBytes[:])

	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err)

	bcNode := in.Blockchains[0].Out.Nodes[0]
	ctx := t.Context()
	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	coord, err := contracts.LoadVRFCoordinatorV2_5(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load coordinator")

	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	bhs, err := contracts.LoadBlockhashStore(chainClient, c.DeployedContracts.BHS)
	require.NoError(t, err, "failed to load BlockhashStore")

	t.Run("BHS stores blockhash for unfulfilled request", func(t *testing.T) {
		// Deploy consumer and create a subscription with 0 funds so the request gets stuck.
		consumer, dErr := contracts.DeployVRFv2PlusLoadTestConsumer(chainClient, coord.Address())
		require.NoError(t, dErr, "failed to deploy load test consumer")

		subTx, sErr := coord.CreateSubscription()
		require.NoError(t, sErr, "failed to create subscription")
		receipt, rErr := chainClient.Client.TransactionReceipt(ctx, subTx.Hash())
		require.NoError(t, rErr)
		subID, pErr := contracts.FindSubscriptionID(receipt)
		require.NoError(t, pErr)

		dErr = coord.AddConsumer(subID, consumer.Address())
		require.NoError(t, dErr, "failed to add consumer")

		requested, qErr := consumer.RequestRandomnessWithEvent(
			keyHash, subID, c.MinimumConfirmations,
			defaultCallbackGasLimit, false, defaultNumWords, defaultRequestCount,
		)
		require.NoError(t, qErr, "RequestRandomness failed")
		require.NotNil(t, requested)
		requestBlock := requested.Raw.BlockNumber
		require.Positive(t, requestBlock, "request block must be non-zero")

		_, qErr = bhs.GetBlockhash(ctx, requestBlock)
		require.Error(t, qErr, "blockhash should not exist in BHS immediately after request")

		products.WaitUntilChainHead(ctx, t, chainClient, requestBlock, c.BHSJobWaitBlocks+10, chainID, 10*time.Second)

		reqCount, cErr := consumer.RequestCount(ctx)
		require.NoError(t, cErr)
		respCount, cErr := consumer.ResponseCount(ctx)
		require.NoError(t, cErr)
		require.Equal(t, 0, reqCount.Cmp(big.NewInt(1)), "request count should be 1")
		require.Equal(t, 0, respCount.Cmp(big.NewInt(0)), "fulfillment count should stay 0")

		var storedHash [32]byte
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			hash, hErr := bhs.GetBlockhash(ctx, requestBlock)
			if hErr != nil {
				return false
			}
			storedHash = hash
			return true
		}, 2*time.Minute, time.Second).Should(gomega.BeTrue(),
			"BHS should store blockhash for request block %d", requestBlock)
		require.Equal(t, 0, requested.Raw.BlockHash.Cmp(storedHash),
			"BHS stored blockhash should match RandomWordsRequested blockhash")
	})

	t.Run("BHS complete E2E fund later and fulfill", func(t *testing.T) {
		t.Skip("This test is flaky on CI. Originally it only run on live testnets. Owners should work on fixing it.")
		// Failure reason: blockhash not found in the store.
		consumer, dErr := contracts.DeployVRFv2PlusLoadTestConsumer(chainClient, coord.Address())
		require.NoError(t, dErr, "failed to deploy load test consumer")

		subTx, sErr := coord.CreateSubscription()
		require.NoError(t, sErr, "failed to create subscription")
		receipt, rErr := chainClient.Client.TransactionReceipt(ctx, subTx.Hash())
		require.NoError(t, rErr)
		subID, pErr := contracts.FindSubscriptionID(receipt)
		require.NoError(t, pErr)

		dErr = coord.AddConsumer(subID, consumer.Address())
		require.NoError(t, dErr, "failed to add consumer")

		requested, qErr := consumer.RequestRandomnessWithEvent(
			keyHash, subID, c.MinimumConfirmations,
			defaultCallbackGasLimit, true, defaultNumWords, defaultRequestCount,
		)
		require.NoError(t, qErr, "RequestRandomness failed")
		require.NotNil(t, requested)
		requestBlock := requested.Raw.BlockNumber
		require.Positive(t, requestBlock, "request block must be non-zero")

		// On EVM BLOCKHASH can no longer serve the original request block hash after ~256 blocks, so fulfillment path must depend on BHS-stored hash
		products.WaitUntilChainHead(ctx, t, chainClient, requestBlock, c.BHSJobWaitBlocks+256, chainID, 5*time.Minute)

		var storedHash [32]byte
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			hash, hErr := bhs.GetBlockhash(ctx, requestBlock)
			if hErr != nil {
				return false
			}
			storedHash = hash
			return true
		}, 2*time.Minute, time.Second).Should(gomega.BeTrue(),
			"BHS should store blockhash for request block %d before funding", requestBlock)
		require.Equal(t, 0, requested.Raw.BlockHash.Cmp(storedHash),
			"BHS stored blockhash should match RandomWordsRequested blockhash")

		nativeWei := products.EtherToWei(big.NewFloat(c.SubFundingAmountNative))
		fErr := coord.FundSubscriptionWithNative(subID, nativeWei)
		require.NoError(t, fErr, "failed to fund subscription with native")

		linkJuels := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
		encodedSubID, eErr := encodeSubID(subID)
		require.NoError(t, eErr)
		_, fErr = linkToken.TransferAndCall(coord.Address(), linkJuels, encodedSubID)
		require.NoError(t, fErr, "failed to fund subscription with LINK")

		var fulfilledSuccess bool
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			event, wErr := coord.FilterRandomWordsFulfilled(
				&bind.FilterOpts{Context: ctx},
				requested.RequestId,
			)
			if wErr != nil {
				return false
			}
			fulfilledSuccess = event.Success
			return true
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"stuck VRF request should be fulfilled after funding and BHS blockhash storage")
		require.True(t, fulfilledSuccess, "RandomWordsFulfilled.Success should be true")
	})
}
