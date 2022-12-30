//go:build integration

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

	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra/denom"
	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
)

func TestClient_SendTerraCoins(t *testing.T) {
	t.Skip("requires terrad")
	chainID := terratest.RandomChainID()
	accounts, _, tendermintURL := terraclient.SetupLocalTerraNode(t, chainID)
	require.Greater(t, len(accounts), 1)
	app := terraStartNewApplication(t)

	from := accounts[0]
	to := accounts[1]
	require.NoError(t, app.GetKeyStore().Terra().Add(terrakey.Raw(from.PrivateKey.Bytes()).Key()))

	chains := app.GetChains()
	_, err := chains.Terra.Add(testutils.Context(t), chainID, nil)
	require.NoError(t, err)
	chain, err := chains.Terra.Chain(testutils.Context(t), chainID)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	_, err = chains.Terra.CreateNode(ctx, terradb.Node{
		Name:          t.Name(),
		TerraChainID:  chainID,
		TendermintURL: tendermintURL,
	})
	require.NoError(t, err)

	reader, err := chain.Reader("")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		coin, err := reader.Balance(from.Address, "uluna")
		if !assert.NoError(t, err) {
			return false
		}
		return coin.IsPositive()
	}, time.Minute, 5*time.Second)

	db := app.GetSqlxDB()
	orm := terratxm.NewORM(chainID, db, logger.TestLogger(t), pgtest.NewQConfig(true))

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
			startBal, err := reader.Balance(from.Address, "uluna")
			require.NoError(t, err)

			set := flag.NewFlagSet("sendterracoins", 0)
			cltest.CopyFlagSetFromAction(client.TerraSendLuna, set, "terra")

			require.NoError(t, set.Set("id", chainID))
			require.NoError(t, set.Parse([]string{tt.amount, from.Address.String(), to.Address.String()}))

			c := cli.NewContext(cliapp, set, nil)
			err = client.TerraSendLuna(c)
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
			renderedMsg := renderer.(*cmd.TerraMsgPresenter)
			require.NotEmpty(t, renderedMsg.ID)
			assert.Equal(t, string(terradb.Unstarted), renderedMsg.State)
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
			require.NotEqual(t, terradb.Errored, msg.State)
			switch msg.State {
			case terradb.Unstarted:
				assert.Nil(t, msg.TxHash)
			case terradb.Broadcasted, terradb.Confirmed:
				assert.NotNil(t, msg.TxHash)
			}

			// Maybe wait for confirmation
			if msg.State != terradb.Confirmed {
				require.Eventually(t, func() bool {
					msgs, err := orm.GetMsgs(id)
					if assert.NoError(t, err) && assert.NotEmpty(t, msgs) {
						if msg = msgs[0]; assert.Equal(t, msg.ID, id) {
							t.Log("State:", msg.State)
							return msg.State == terradb.Confirmed
						}
					}
					return false
				}, testutils.WaitTimeout(t), time.Second)
				require.NotNil(t, msg.TxHash)
			}

			// Check balance
			endBal, err := reader.Balance(from.Address, "uluna")
			require.NoError(t, err)
			if assert.NotNil(t, startBal) && assert.NotNil(t, endBal) {
				diff := startBal.Sub(*endBal).Amount
				sent, err := denom.ConvertToULuna(sdk.NewDecCoinFromDec("luna", sdk.MustNewDecFromStr(tt.amount)))
				require.NoError(t, err)
				if assert.True(t, diff.IsInt64()) && assert.True(t, sent.Amount.IsInt64()) {
					require.Greater(t, diff.Int64(), sent.Amount.Int64())
				}
			}
		})
	}
}
