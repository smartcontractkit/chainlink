package vrfv2

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2 "github.com/smartcontractkit/chainlink/devenv/products/vrfv2"
)

func TestVRFv2MultipleSendingKeys(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrfv2-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)

	cfg, err := products.LoadOutput[productvrfv2.Configurator](outputFile)
	require.NoError(t, err)
	c := cfg.Config[0]
	require.GreaterOrEqual(t, c.NumTxKeys, 2)

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

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	t.Run("Request Randomness with multiple sending keys", func(t *testing.T) {
		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err)
		subID := subIDs[0]

		txKeys, _, err := cl[0].ReadTxKeys("evm")
		require.NoError(t, err)
		require.Len(t, txKeys.Data, c.NumTxKeys+1)

		fulfillTimeout := parseFulfillTimeout(c.RandomWordsFulfilledEventTimeout)
		// Match legacy smoke/vrfv2_test.go: always use Seth key 0 for consumer requests.
		// The assertion is on fulfillment tx senders — the VRF job rotates node fromAddresses.
		var fromAddrs []string
		for range c.NumTxKeys + 1 {
			_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
				c.MinimumConfirmations, c.CallbackGasLimit, c.NumberOfWords,
				c.RandomnessRequestCountPerRequest, c.RandomnessRequestCountPerRequestDeviation,
				fulfillTimeout, 0)
			require.NoError(t, err)
			tx, _, err := chainClient.Client.TransactionByHash(ctx, fulfilled.Raw.TxHash)
			require.NoError(t, err)
			from, err := getTxFromAddress(tx)
			require.NoError(t, err)
			fromAddrs = append(fromAddrs, strings.ToLower(from))
		}

		var keyAddrs []string
		for _, k := range txKeys.Data {
			keyAddrs = append(keyAddrs, strings.ToLower(k.Attributes.Address))
		}
		less := func(a, b string) bool { return a < b }
		require.Empty(t, cmp.Diff(keyAddrs, fromAddrs, cmpopts.SortSlices(less)))
	})
}
