//go:build integration && wasmd

package cmd_test

import (
	"flag"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	cosmosdb "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/cosmostxm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/denom"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

func TestClient_SendCosmosCoins(t *testing.T) {
	// TODO(BCI-978): cleanup once SetupLocalCosmosNode is updated
	chainID := cosmostest.RandomChainID()
	accounts, _, _ := cosmosclient.SetupLocalCosmosNode(t, chainID)
	require.Greater(t, len(accounts), 1)
	app := cosmosStartNewApplication(t)

	from := accounts[0]
	to := accounts[1]
	require.NoError(t, app.GetKeyStore().Cosmos().Add(cosmoskey.Raw(from.PrivateKey.Bytes()).Key()))

	chain, err := app.GetChains().Cosmos.Chain(testutils.Context(t), chainID)
	require.NoError(t, err)

	reader, err := chain.Reader("")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		coin, err := reader.Balance(from.Address, "uatom")
		if !assert.NoError(t, err) {
			return false
		}
		return coin.IsPositive()
	}, time.Minute, 5*time.Second)

	db := app.GetSqlxDB()
	orm := cosmostxm.NewORM(chainID, db, logger.TestLogger(t), pgtest.NewQConfig(true))

	client, r := app.NewClientAndRenderer()
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
		{amount: "asdf", expErr: "invalid coin: failed to set decimal string:"},
	} {
		tt := tt
		t.Run(tt.amount, func(t *testing.T) {
			startBal, err := reader.Balance(from.Address, "uatom")
			require.NoError(t, err)

			set := flag.NewFlagSet("sendcosmoscoins", 0)
			cltest.FlagSetApplyFromAction(client.CosmosSendAtom, set, "cosmos")

			require.NoError(t, set.Set("id", chainID))
			require.NoError(t, set.Parse([]string{tt.amount, from.Address.String(), to.Address.String()}))

			c := cli.NewContext(cliapp, set, nil)
			err = client.CosmosSendAtom(c)
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
			id, err := strconv.ParseInt(renderedMsg.ID, 10, 64)
			require.NoError(t, err)
			msgs, err := orm.GetMsgs(id)
			require.NoError(t, err)
			require.Equal(t, 1, len(msgs))
			msg := msgs[0]
			assert.Equal(t, strconv.FormatInt(msg.ID, 10), renderedMsg.ID)
			assert.Equal(t, msg.ChainID, renderedMsg.ChainID)
			assert.Equal(t, msg.ContractID, renderedMsg.ContractID)
			require.NotEqual(t, cosmosdb.Errored, msg.State)
			switch msg.State {
			case cosmosdb.Unstarted:
				assert.Nil(t, msg.TxHash)
			case cosmosdb.Broadcasted, cosmosdb.Confirmed:
				assert.NotNil(t, msg.TxHash)
			}

			// Maybe wait for confirmation
			if msg.State != cosmosdb.Confirmed {
				require.Eventually(t, func() bool {
					msgs, err := orm.GetMsgs(id)
					if assert.NoError(t, err) && assert.NotEmpty(t, msgs) {
						if msg = msgs[0]; assert.Equal(t, msg.ID, id) {
							t.Log("State:", msg.State)
							return msg.State == cosmosdb.Confirmed
						}
					}
					return false
				}, testutils.WaitTimeout(t), time.Second)
				require.NotNil(t, msg.TxHash)
			}

			// Check balance
			endBal, err := reader.Balance(from.Address, "uatom")
			require.NoError(t, err)
			if assert.NotNil(t, startBal) && assert.NotNil(t, endBal) {
				diff := startBal.Sub(*endBal).Amount
				sent, err := denom.DecCoinToUAtom(sdk.NewDecCoinFromDec("atom", sdk.MustNewDecFromStr(tt.amount)))
				require.NoError(t, err)
				if assert.True(t, diff.IsInt64()) && assert.True(t, sent.Amount.IsInt64()) {
					require.Greater(t, diff.Int64(), sent.Amount.Int64())
				}
			}
		})
	}
}
