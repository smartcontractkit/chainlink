package cmd_test

import (
	"flag"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr/txmgrtest"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestShell_IndexTransactions(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	txStore := txmgrtest.NewTestTxStore(t, app.GetDB())
	tx := txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)
	attempt := tx.TxAttempts[0]

	// page 1
	set := flag.NewFlagSet("test transactions", 0)
	flagSetApplyFromAction(client.IndexTransactions, set, "")

	require.NoError(t, set.Set("page", "1"))

	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs := *r.Renders[0].(*cmd.EthTxPresenters)
	assert.Len(t, renderedTxs, 1)
	assert.Equal(t, attempt.Hash.String(), renderedTxs[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test txattempts", 0)
	flagSetApplyFromAction(client.IndexTransactions, set, "")

	require.NoError(t, set.Set("page", "2"))

	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs = *r.Renders[1].(*cmd.EthTxPresenters)
	assert.Empty(t, renderedTxs)
}

func TestShell_ShowTransaction(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	db := app.GetDB()
	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	txStore := txmgrtest.NewTestTxStore(t, db)
	tx := txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)
	attempt := tx.TxAttempts[0]

	set := flag.NewFlagSet("test get tx", 0)
	flagSetApplyFromAction(client.ShowTransaction, set, "")

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

	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	txStore := txmgrtest.NewTestTxStore(t, app.GetDB())
	tx := txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)

	// page 1
	set := flag.NewFlagSet("test txattempts", 0)
	flagSetApplyFromAction(client.IndexTxAttempts, set, "")

	require.NoError(t, set.Set("page", "1"))

	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	require.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts := *r.Renders[0].(*cmd.EthTxPresenters)
	require.Len(t, tx.TxAttempts, 1)
	assert.Equal(t, tx.TxAttempts[0].Hash.String(), renderedAttempts[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test transactions", 0)
	flagSetApplyFromAction(client.IndexTxAttempts, set, "")

	require.NoError(t, set.Set("page", "2"))

	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts = *r.Renders[1].(*cmd.EthTxPresenters)
	assert.Empty(t, renderedAttempts)
}

func TestShell_SendEther_From_Txm(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)
	fromAddress := key.Address

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethMock := newEthMockWithTransactionsOnBlocksAssertions(t)

	ethMock.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)
	ethMock.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)

		// NOTE: FallbackPollInterval is used in this test to quickly create TxAttempts
		// Testing triggers requires committing transactions and does not work with transactional tests
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(time.Second)
	},
		withKey(),
		withMocks(ethMock, key),
	)
	client, r := app.NewShellAndRenderer()
	db := app.GetDB()
	txStore := txmgr.NewTxStore(db, logger.TestLogger(t))
	set := flag.NewFlagSet("sendether", 0)
	flagSetApplyFromAction(client.SendEther, set, "")

	amount := "100.5"
	to := "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"
	require.NoError(t, set.Parse([]string{amount, fromAddress.Hex(), to}))
	require.NoError(t, set.Set("id", evmtest.MustGetDefaultChainID(t, app.Config.EVMConfigs()).String()))

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))

	evmTxes, err := txStore.GetAllTxes(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, evmTxes, 1)
	evmTx := evmTxes[0]
	value := assets.Eth(evmTx.Value)
	require.Equal(t, "100.500000000000000000", value.String())
	require.Equal(t, fromAddress, evmTx.FromAddress)
	require.Equal(t, to, evmTx.ToAddress.String())

	output := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &evmTx.FromAddress, output.From)
	assert.Equal(t, &evmTx.ToAddress, output.To)
	assert.Equal(t, value.String(), output.Value)
	assert.Equal(t, fmt.Sprintf("%d", *evmTx.Sequence), output.Nonce)

	attempts, err := txStore.GetAllTxAttempts(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, attempts, 1)
	assert.Equal(t, attempts[0].Hash, output.Hash)
}

func TestShell_SendEther_From_Txm_WEI(t *testing.T) {
	t.Parallel()

	key := cltest.MustGenerateRandomKey(t)
	fromAddress := key.Address

	balance, err := assets.NewEthValueS("200")
	require.NoError(t, err)

	ethMock := newEthMockWithTransactionsOnBlocksAssertions(t)

	ethMock.On("BalanceAt", mock.Anything, key.Address, (*big.Int)(nil)).Return(balance.ToInt(), nil)
	ethMock.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)

		// NOTE: FallbackPollInterval is used in this test to quickly create TxAttempts
		// Testing triggers requires committing transactions and does not work with transactional tests
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(time.Second)
	},
		withKey(),
		withMocks(ethMock, key),
	)
	client, r := app.NewShellAndRenderer()
	db := app.GetDB()
	txStore := txmgr.NewTxStore(db, logger.TestLogger(t))

	set := flag.NewFlagSet("sendether", 0)
	flagSetApplyFromAction(client.SendEther, set, "")

	require.NoError(t, set.Set("id", testutils.FixtureChainID.String()))
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

	evmTxes, err := txStore.GetAllTxes(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, evmTxes, 1)
	evmTx := evmTxes[0]
	value := assets.Eth(evmTx.Value)
	require.Equal(t, "1.000000000000000000", value.String())
	require.Equal(t, fromAddress, evmTx.FromAddress)
	require.Equal(t, to, evmTx.ToAddress.String())

	output := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &evmTx.FromAddress, output.From)
	assert.Equal(t, &evmTx.ToAddress, output.To)
	assert.Equal(t, value.String(), output.Value)
	assert.Equal(t, fmt.Sprintf("%d", *evmTx.Sequence), output.Nonce)

	attempts, err := txStore.GetAllTxAttempts(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, attempts, 1)
	assert.Equal(t, attempts[0].Hash, output.Hash)
}
