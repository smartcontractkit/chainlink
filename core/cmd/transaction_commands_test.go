package cmd_test

import (
	"flag"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_IndexTransactions(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	store := app.GetStore()
	_, from := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth)

	tx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)
	attempt := tx.EthTxAttempts[0]

	// page 1
	set := flag.NewFlagSet("test transactions", 0)
	set.Int("page", 1, "doc")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs := *r.Renders[0].(*cmd.EthTxPresenters)
	assert.Equal(t, 1, len(renderedTxs))
	assert.Equal(t, attempt.Hash.Hex(), renderedTxs[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test txattempts", 0)
	set.Int("page", 2, "doc")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs = *r.Renders[1].(*cmd.EthTxPresenters)
	assert.Equal(t, 0, len(renderedTxs))
}

func TestClient_ShowTransaction(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	store := app.GetStore()
	_, from := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth)

	tx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)
	attempt := tx.EthTxAttempts[0]

	set := flag.NewFlagSet("test get tx", 0)
	set.Parse([]string{attempt.Hash.Hex()})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ShowTransaction(c))

	renderedTx := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &tx.FromAddress, renderedTx.From)
}

func TestClient_IndexTxAttempts(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	store := app.GetStore()
	_, from := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth)

	tx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)

	// page 1
	set := flag.NewFlagSet("test txattempts", 0)
	set.Int("page", 1, "doc")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts := *r.Renders[0].(*cmd.EthTxPresenters)
	require.Len(t, tx.EthTxAttempts, 1)
	assert.Equal(t, tx.EthTxAttempts[0].Hash.Hex(), renderedAttempts[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test transactions", 0)
	set.Int("page", 2, "doc")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts = *r.Renders[1].(*cmd.EthTxPresenters)
	assert.Equal(t, 0, len(renderedAttempts))
}

func TestClient_SendEther_From_BPTXM(t *testing.T) {
	t.Parallel()

	oca := common.HexToAddress("0xDEADB3333333F")
	app := startNewApplication(t,
		withKey(),
		withConfig(map[string]interface{}{
			"OPERATOR_CONTRACT_ADDRESS": &oca,
		}),
		withMocks(newEthMock(t)),
		startAndConnect(),
	)
	client, r := app.NewClientAndRenderer()
	s := app.GetStore()

	set := flag.NewFlagSet("sendether", 0)
	amount := "100.5"
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, app.KeyStore.Eth, 0)
	to := "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"
	set.Parse([]string{amount, fromAddress.Hex(), to})

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))

	etx := models.EthTx{}
	require.NoError(t, s.DB.First(&etx).Error)
	require.Equal(t, "100.500000000000000000", etx.Value.String())
	require.Equal(t, fromAddress, etx.FromAddress)
	require.Equal(t, to, etx.ToAddress.Hex())

	output := *r.Renders[0].(*cmd.EthTxPresenter)
	assert.Equal(t, &etx.FromAddress, output.From)
	assert.Equal(t, &etx.ToAddress, output.To)
	assert.Equal(t, etx.Value.String(), output.Value)
}
