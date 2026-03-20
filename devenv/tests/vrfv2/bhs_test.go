package vrfv2

import (
	"errors"
	"fmt"
	"math/big"
	"os"
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
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrfv2-bhs-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)

	cfg, err := products.LoadOutput[productvrfv2.Configurator](outputFile)
	require.NoError(t, err)
	c := cfg.Config[0]
	require.True(t, c.EnableBHSJob)
	require.NotEmpty(t, c.VRFKeyData.BHSJobID)

	keyHash := mustKeyHash(c)
	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err)
	ctx := t.Context()
	chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err)

	coord, err := contracts.LoadVRFCoordinatorV2(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err)
	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err)
	bhs, err := contracts.LoadBlockhashStore(chainClient, c.DeployedContracts.BHS)
	require.NoError(t, err)

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	t.Run("BHS Job with complete E2E - wait 256 blocks to see if Rand Request is fulfilled", func(t *testing.T) {
		if os.Getenv("TEST_UNSKIP") != "true" {
			t.Skip(`Skipped due to long execution time. Run on-demand with TEST_UNSKIP="true".`)
		}
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, 0, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		req, err := consumers[0].RequestRandomnessFromKey(coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, 0)
		require.NoError(t, err)
		reqBlock := req.Raw.BlockNumber

		err = waitUntilBlock(ctx, chainClient, reqBlock+257, time.Second, parseFulfillTimeout("280s"))
		require.NoError(t, err)

		amount := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
		enc, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID)
		require.NoError(t, err)
		_, err = linkToken.TransferAndCall(coord.Address(), amount, enc)
		require.NoError(t, err)

		fulfillTimeout := parseFulfillTimeout(c.RandomWordsFulfilledEventTimeout)
		_, err = contracts.WaitRandomWordsFulfilled(coord, req.RequestID, req.Raw.BlockNumber, fulfillTimeout)
		require.NoError(t, err)
	})

	t.Run("BHS Job should fill in blockhashes into BHS contract for unfulfilled requests", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, 0, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		req, err := consumers[0].RequestRandomnessFromKey(coord, keyHash, subID,
			c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
			c.RandomnessRequestCountPerRequest, 0)
		require.NoError(t, err)
		reqBlock := req.Raw.BlockNumber

		_, err = bhs.GetBlockhash(ctx, reqBlock)
		require.Error(t, err)

		blocks := c.BHSJobWaitBlocks
		if blocks < 0 {
			t.Fatalf("negative blocks: %d", blocks)
		}
		err = waitUntilBlock(ctx, chainClient, reqBlock+uint64(blocks), time.Second, time.Minute)
		require.NoError(t, err)

		metrics, err := consumers[0].GetLoadTestMetrics(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, metrics.RequestCount.Cmp(big.NewInt(1)))
		require.Equal(t, 0, metrics.FulfilmentCount.Cmp(big.NewInt(0)))

		gom := gomega.NewGomegaWithT(t)
		var bhsNodeTxHash string
		gom.Eventually(func(g gomega.Gomega) {
			txs, _, rErr := cl[1].ReadTransactions()
			g.Expect(rErr).ShouldNot(gomega.HaveOccurred())
			g.Expect(txs.Data).Should(gomega.HaveLen(1))
			g.Expect(strings.EqualFold(txs.Data[0].Attributes.To, c.DeployedContracts.BHS)).Should(gomega.BeTrue())
			bhsNodeTxHash = txs.Data[0].Attributes.Hash
		}, "2m", "1s").Should(gomega.Succeed())

		tx, _, err := chainClient.Client.TransactionByHash(ctx, common.HexToHash(bhsNodeTxHash))
		require.NoError(t, err)
		storedBlock, err := decodeBHSStoreBlockNumber(tx.Data())
		require.NoError(t, err)
		require.Equal(t, reqBlock, storedBlock)

		var storedHash [32]byte
		gom.Eventually(func(g gomega.Gomega) {
			h, hErr := bhs.GetBlockhash(ctx, reqBlock)
			g.Expect(hErr).ShouldNot(gomega.HaveOccurred())
			storedHash = h
		}, "2m", "1s").Should(gomega.Succeed())
		require.Equal(t, 0, req.Raw.BlockHash.Cmp(common.BytesToHash(storedHash[:])))
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
