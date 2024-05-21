//go:build integration

package cmd_test

import (
	"flag"
	"os"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	cosmosdb "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/denom"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/params"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

var nativeToken = "cosm"

func TestMain(m *testing.M) {

	params.InitCosmosSdk(
		/* bech32Prefix= */ "wasm",
		/* token= */ nativeToken,
	)

	code := m.Run()
	os.Exit(code)
}

func TestShell_SendCosmosCoins(t *testing.T) {
	ctx := testutils.Context(t)
	// TODO(BCI-978): cleanup once SetupLocalCosmosNode is updated
	chainID := cosmostest.RandomChainID()
	cosmosChain := coscfg.Chain{}
	cosmosChain.SetDefaults()
	accounts, _, url := cosmosclient.SetupLocalCosmosNode(t, chainID, *cosmosChain.GasToken)
	require.Greater(t, len(accounts), 1)
	nodes := coscfg.Nodes{
		&coscfg.Node{
			Name:          ptr("random"),
			TendermintURL: config.MustParseURL(url),
		},
	}
	chainConfig := coscfg.TOMLConfig{ChainID: &chainID, Enabled: ptr(true), Chain: cosmosChain, Nodes: nodes}
	app := cosmosStartNewApplication(t, &chainConfig)

	from := accounts[0]
	to := accounts[1]
	require.NoError(t, app.GetKeyStore().Cosmos().Add(ctx, cosmoskey.Raw(from.PrivateKey.Bytes()).Key()))
	chain, err := app.GetRelayers().LegacyCosmosChains().Get(chainID)
	require.NoError(t, err)

	reader, err := chain.Reader("")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		coin, err := reader.Balance(from.Address, *cosmosChain.GasToken)
		if !assert.NoError(t, err) {
			return false
		}
		return coin.IsPositive()
	}, time.Minute, 5*time.Second)

	client, r := app.NewShellAndRenderer()
	cliapp := cli.NewApp()

	for _, tt := range []struct {
		amount string
		expErr string
	}{
		{amount: "0.000001"},
		{amount: "1"},
		{amount: "30.000001"},
		{amount: "1000", expErr: "is too low for this transaction to be executed:"},
		{amount: "0", expErr: "amount must be greater than zero:"},
		{amount: "asdf", expErr: "invalid coin: failed to set decimal string"},
	} {
		tt := tt
		t.Run(tt.amount, func(t *testing.T) {
			startBal, err := reader.Balance(from.Address, *cosmosChain.GasToken)
			require.NoError(t, err)

			set := flag.NewFlagSet("sendcosmoscoins", 0)
			flagSetApplyFromAction(client.CosmosSendNativeToken, set, "cosmos")

			require.NoError(t, set.Set("id", chainID))
			require.NoError(t, set.Parse([]string{nativeToken, tt.amount, from.Address.String(), to.Address.String()}))

			c := cli.NewContext(cliapp, set, nil)
			err = client.CosmosSendNativeToken(c)
			if tt.expErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErr)
				return
			}

			// Check CLI output
			require.Greater(t, len(r.Renders), 0)
			renderer := r.Renders[len(r.Renders)-1]
			renderedMsg := renderer.(*cmd.CosmosMsgPresenter)
			require.NotEmpty(t, renderedMsg.ID)
			assert.Equal(t, string(cosmosdb.Unstarted), renderedMsg.State)
			assert.Nil(t, renderedMsg.TxHash)

			// Check balance
			sent, err := denom.ConvertDecCoinToDenom(sdk.NewDecCoinFromDec(nativeToken, sdk.MustNewDecFromStr(tt.amount)), *cosmosChain.GasToken)
			require.NoError(t, err)
			expBal := startBal.Sub(sent)

			testutils.AssertEventually(t, func() bool {
				endBal, err := reader.Balance(from.Address, *cosmosChain.GasToken)
				require.NoError(t, err)
				t.Logf("%s <= %s", endBal, expBal)
				return endBal.IsLTE(expBal)
			})
		})
	}
}
