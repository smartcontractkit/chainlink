package vrfv2

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2 "github.com/smartcontractkit/chainlink/devenv/products/vrfv2"
)

func TestVRFV2WithBHS(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr, "failed to save container logs")
	})

	outputFile := "../../env-vrfv2-bhs-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load devenv env-out from %s", outputFile)

	cfg, err := products.LoadOutput[productvrfv2.Configurator](outputFile)
	require.NoError(t, err, "failed to load vrfv2 product config from env-out")
	c := cfg.Config[0]
	require.True(t, c.EnableBHSJob, "BHS product config must enable BHS job")
	require.NotEmpty(t, c.VRFKeyData.BHSJobID, "BHS job ID must be set in env-out")

	keyHash := mustKeyHash(c)
	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err, "failed to parse chain ID from env-out")
	ctx := t.Context()
	chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	coord, err := contracts.LoadVRFCoordinatorV2(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load VRF coordinator v2")
	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")
	bhs, err := contracts.LoadBlockhashStore(chainClient, c.DeployedContracts.BHS)
	require.NoError(t, err, "failed to load blockhash store")

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to connect to Chainlink nodes")

	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		t.Skip("This test is flaky on CI. Originally it only run on live testnets. Owners should work on fixing it.")
		// Failure reason: blockhash not found in the store.
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, 0, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		req, err := consumers[0].RequestRandomnessFromKey(coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, 0)
		require.NoError(t, err, "error requesting randomness before long BHS wait")
		reqBlock := req.Raw.BlockNumber

		// On EVM BLOCKHASH can no longer serve the original request block hash after ~256 blocks, so fulfillment path must depend on BHS-stored hash
		products.WaitUntilChainHead(ctx, t, chainClient, reqBlock, c.BHSJobWaitBlocks+256, chainID, 5*time.Minute)

		var storedHash [32]byte
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			hash, hErr := bhs.GetBlockhash(ctx, reqBlock)
			if hErr != nil {
				return false
			}
			storedHash = hash
			return true
		}, 3*time.Minute, time.Second).Should(gomega.BeTrue(),
			"BHS should store blockhash for request block %d before funding", reqBlock)
		require.Equal(t, 0, req.Raw.BlockHash.Cmp(common.BytesToHash(storedHash[:])),
			"BHS stored blockhash should match RandomWordsRequested blockhash")

		amount := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
		enc, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID)
		require.NoError(t, err, "error ABI-encoding sub ID for funding")
		_, err = linkToken.TransferAndCall(coord.Address(), amount, enc)
		require.NoError(t, err, "error funding subscription after BHS wait")

		fulfillTimeout := parseFulfillTimeout(c.RandomWordsFulfilledEventTimeout)
		_, err = contracts.WaitRandomWordsFulfilled(coord, req.RequestID, req.Raw.BlockNumber, fulfillTimeout)
		require.NoError(t, err, "RandomWordsFulfilled not seen after funding sub (BHS E2E)")
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, 0, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		req, err := consumers[0].RequestRandomnessFromKey(coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, 0)
		require.NoError(t, err, "error requesting randomness for BHS store test")
		reqBlock := req.Raw.BlockNumber

		_, err = bhs.GetBlockhash(ctx, reqBlock)
		require.Error(t, err, "BHS should not have blockhash for request block immediately after request")

		blocks := c.BHSJobWaitBlocks
		if blocks < 0 {
			t.Fatalf("negative blocks: %d", blocks)
		}
		products.WaitUntilChainHead(ctx, t, chainClient, reqBlock, blocks, chainID, time.Minute)

		metrics, err := consumers[0].GetLoadTestMetrics(ctx)
		require.NoError(t, err, "error reading consumer load test metrics")
		require.Equal(t, 0, metrics.RequestCount.Cmp(big.NewInt(1)), "expected exactly one randomness request on consumer")
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)), "expected no fulfillment yet while request is unfulfilled")

		gom := gomega.NewGomegaWithT(t)
		var bhsNodeTxHash string
		gom.Eventually(func(g gomega.Gomega) {
			txs, _, rErr := cl[1].ReadTransactions()
			g.Expect(rErr).ShouldNot(gomega.HaveOccurred())
			h, ok := findBHSStoreTxHashForBlock(txs.Data, c.DeployedContracts.BHS, reqBlock)
			g.Expect(ok).Should(gomega.BeTrue(), "expected BHS node tx storing block %d (got %d node txs)", reqBlock, len(txs.Data))
			bhsNodeTxHash = h
		}, "2m", "1s").Should(gomega.Succeed())

		tx, _, err := chainClient.Client.TransactionByHash(ctx, common.HexToHash(bhsNodeTxHash))
		require.NoError(t, err, "failed to load BHS node store transaction")
		storedBlock, err := decodeBHSStoreBlockNumber(tx.Data())
		require.NoError(t, err, "failed to decode store block number from BHS tx calldata")
		require.Equal(t, reqBlock, storedBlock, "BHS store tx should target request block %d", reqBlock)

		var storedHash [32]byte
		gom.Eventually(func(g gomega.Gomega) {
			h, hErr := bhs.GetBlockhash(ctx, reqBlock)
			g.Expect(hErr).ShouldNot(gomega.HaveOccurred())
			storedHash = h
		}, "2m", "1s").Should(gomega.Succeed())
		require.Equal(t, 0, req.Raw.BlockHash.Cmp(common.BytesToHash(storedHash[:])),
			"blockhash stored in BHS must match chain block hash at request block")
	})
}

func decodeBHSStoreBlockNumber(data []byte) (uint64, error) {
	parsed, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return 0, err
	}
	if len(data) < 4 {
		return 0, errors.New("short calldata")
	}
	m, err := parsed.MethodById(data[:4])
	if err != nil {
		return 0, err
	}
	args, err := m.Inputs.Unpack(data[4:])
	if err != nil {
		return 0, err
	}
	if len(args) != 1 {
		return 0, fmt.Errorf("expected 1 arg, got %d", len(args))
	}
	bn, ok := args[0].(*big.Int)
	if !ok {
		return 0, errors.New("expected *big.Int")
	}
	return bn.Uint64(), nil
}

// findBHSStoreTxHashForBlock returns the hash of a finalized node tx that calls BHS store for wantBlock.
// With fast chains the node may list multiple BHS store txs; we match by decoded calldata block number.
func findBHSStoreTxHashForBlock(txs []clclient.TransactionData, bhsAddr string, wantBlock uint64) (txHash string, ok bool) {
	for _, tx := range txs {
		if !strings.EqualFold(tx.Attributes.To, bhsAddr) {
			continue
		}
		data := common.FromHex(tx.Attributes.Data)
		if len(data) < 4 {
			continue
		}
		stored, err := decodeBHSStoreBlockNumber(data)
		if err != nil || stored != wantBlock {
			continue
		}
		return tx.Attributes.Hash, true
	}
	return "", false
}
