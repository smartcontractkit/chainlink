package vrfv2plus

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFV2PlusWithBHF(t *testing.T) {
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
	require.NotEmpty(t, c.VRFKeyData.BHFJobID, "BHF job ID must not be empty — was env-vrf2plus-bhX.toml used?")
	require.NotEmpty(t, c.DeployedContracts.BHS, "BHS contract address must not be empty")
	require.NotEmpty(t, c.DeployedContracts.BatchBHS, "BatchBHS contract address must not be empty")

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

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to create CL clients")

	t.Run("BHF Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
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

		// Wait at least 257 blocks so the BHF job can store the blockhash and
		// the coordinator can verify it against the BatchBHS.
		products.WaitUntilChainHead(ctx, t, chainClient, requestBlock, 257, chainID, 5*time.Minute)

		// Fund the subscription so the stuck request can be fulfilled.
		nativeWei := products.EtherToWei(big.NewFloat(c.SubFundingAmountNative))
		fErr := coord.FundSubscriptionWithNative(subID, nativeWei)
		require.NoError(t, fErr, "failed to fund subscription with native")

		linkJuels := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
		encodedSubID, eErr := encodeSubID(subID)
		require.NoError(t, eErr)
		_, fErr = linkToken.TransferAndCall(coord.Address(), linkJuels, encodedSubID)
		require.NoError(t, fErr, "failed to fund subscription with LINK")

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			event, wErr := coord.FilterRandomWordsFulfilled(
				&bind.FilterOpts{Context: ctx},
				requested.RequestId,
			)
			if wErr != nil {
				return false
			}
			return event.Success
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"stuck VRF request should be fulfilled after funding and BHF blockhash storage")

		status, sErr := consumer.GetRequestStatus(ctx, requested.RequestId)
		require.NoError(t, sErr, "error getting rand request status")
		require.True(t, status.Fulfilled, "request should be fulfilled")

		// Verify the BHF node sent a transaction to the BatchBHS contract.
		txs, _, txErr := cl[1].ReadTransactions()
		require.NoError(t, txErr, "error fetching txns from BHF node")
		batchBHSTxFound := false
		for _, tx := range txs.Data {
			if strings.EqualFold(tx.Attributes.To, c.DeployedContracts.BatchBHS) {
				batchBHSTxFound = true
				break
			}
		}
		require.True(t, batchBHSTxFound, "BHF node should have sent a tx to BatchBHS")

		// Verify BHS stored the correct blockhash.
		storedHash, hErr := bhs.GetBlockhash(ctx, requestBlock)
		require.NoError(t, hErr, "error getting blockhash for request block")
		require.Equal(t, 0, requested.Raw.BlockHash.Cmp(storedHash),
			"BHS stored blockhash should match RandomWordsRequested blockhash")
	})
}
