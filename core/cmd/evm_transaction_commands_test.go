package cmd_test

import (
	"flag"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestShell_IndexTransactions(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	_, from := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth())

	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.GetConfig().Database())
	tx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)
	attempt := tx.TxAttempts[0]

	// page 1
	set := flag.NewFlagSet("test transactions", 0)
	cltest.FlagSetApplyFromAction(client.IndexTransactions, set, "")

	require.NoError(t, set.Set("page", "1"))

	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs := *r.Renders[0].(*cmd.EthTxPresenters)
	assert.Equal(t, 1, len(renderedTxs))
	assert.Equal(t, attempt.Hash.String(), renderedTxs[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test txattempts", 0)
	cltest.FlagSetApplyFromAction(client.IndexTransactions, set, "")

	require.NoError(t, set.Set("page", "2"))

	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs = *r.Renders[1].(*cmd.EthTxPresenters)
	assert.Equal(t, 0, len(renderedTxs))
}

func TestShell_ShowTransaction(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	db := app.GetSqlxDB()
	_, from := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth())

	txStore := cltest.NewTestTxStore(t, db, app.GetConfig().Database())
	tx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)
	attempt := tx.TxAttempts[0]

	set := flag.NewFlagSet("test get tx", 0)
	cltest.FlagSetApplyFromAction(client.ShowTransaction, set, "")

	require.NoError(t, set.Parse([]string{attempt.Hash.String()}))

	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ShowTransaction(c))

	renderedTx := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &tx.FromAddress, renderedTx.From)
}

func TestShell_IndexTxAttempts(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	_, from := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth())

	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.GetConfig().Database())
	tx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)

	// page 1
	set := flag.NewFlagSet("test txattempts", 0)
	cltest.FlagSetApplyFromAction(client.IndexTxAttempts, set, "")

	require.NoError(t, set.Set("page", "1"))

	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	require.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts := *r.Renders[0].(*cmd.EthTxPresenters)
	require.Len(t, tx.TxAttempts, 1)
	assert.Equal(t, tx.TxAttempts[0].Hash.String(), renderedAttempts[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test transactions", 0)
	cltest.FlagSetApplyFromAction(client.IndexTxAttempts, set, "")

	require.NoError(t, set.Set("page", "2"))

	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts = *r.Renders[1].(*cmd.EthTxPresenters)
	assert.Equal(t, 0, len(renderedAttempts))
}

func TestShell_SendEther_From_Txm(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)
	fromAddress := key.Address

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethMock := newEthMockWithTransactionsOnBlocksAssertions(t)

	ethMock.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethMock, key),
	)
	client, r := app.NewShellAndRenderer()
	db := app.GetSqlxDB()

	set := flag.NewFlagSet("sendether", 0)
	cltest.FlagSetApplyFromAction(client.SendEther, set, "")

	amount := "100.5"
	to := "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"
	require.NoError(t, set.Parse([]string{amount, fromAddress.Hex(), to}))

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))

	dbEvmTx := txmgr.DbEthTx{}
	require.NoError(t, db.Get(&dbEvmTx, `SELECT * FROM eth_txes`))
	require.Equal(t, "100.500000000000000000", dbEvmTx.Value.String())
	require.Equal(t, fromAddress, dbEvmTx.FromAddress)
	require.Equal(t, to, dbEvmTx.ToAddress.String())

	output := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &dbEvmTx.FromAddress, output.From)
	assert.Equal(t, &dbEvmTx.ToAddress, output.To)
	assert.Equal(t, dbEvmTx.Value.String(), output.Value)
}

func TestShell_SendEther_From_Txm_WEI(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)
	fromAddress := key.Address

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethMock := newEthMockWithTransactionsOnBlocksAssertions(t)

	ethMock.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethMock, key),
	)
	client, r := app.NewShellAndRenderer()
	db := app.GetSqlxDB()

	set := flag.NewFlagSet("sendether", 0)
	cltest.FlagSetApplyFromAction(client.SendEther, set, "")

	require.NoError(t, set.Set("wei", "false"))

	amount := "1000000000000000000"
	to := "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"
	err = set.Parse([]string{amount, fromAddress.Hex(), to})
	require.NoError(t, err)

	err = set.Set("wei", "true")
	require.NoError(t, err)

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))

	dbEvmTx := txmgr.DbEthTx{}
	require.NoError(t, db.Get(&dbEvmTx, `SELECT * FROM eth_txes`))
	require.Equal(t, "1.000000000000000000", dbEvmTx.Value.String())
	require.Equal(t, fromAddress, dbEvmTx.FromAddress)
	require.Equal(t, to, dbEvmTx.ToAddress.String())

	output := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &dbEvmTx.FromAddress, output.From)
	assert.Equal(t, &dbEvmTx.ToAddress, output.To)
	assert.Equal(t, dbEvmTx.Value.String(), output.Value)
}
