package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr/txmgrtest"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsController_Index_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	db := app.GetDB()
	txStore := txmgrtest.NewTestTxStore(t, app.GetDB())
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)        // tx1
	tx2 := txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 2, from) // tx2
	txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 4, from)        // tx3

	// add second tx attempt for tx2
	blockNum := int64(3)
	attempt := txmgrtest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

	_, count, err := txStore.TransactionsWithAttempts(ctx, 0, 100)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	size := 2
	resp, cleanup := client.Get(fmt.Sprintf("/v2/transactions?size=%d", size))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	var txs []presenters.EthTxResource
	body := cltest.ParseResponseBody(t, resp)
	require.NoError(t, web.ParsePaginatedResponse(body, &txs, &links))
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	require.Len(t, txs, size)
	require.Equal(t, "4", txs[0].SentAt, "expected tx attempts order by sentAt descending")
	require.Equal(t, "3", txs[1].SentAt, "expected tx attempts order by sentAt descending")
}

func TestTransactionsController_Index_Error(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/transactions?size=TrainingDay")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, 422)
}

func TestTransactionsController_Show_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	txStore := txmgrtest.NewTestTxStore(t, app.GetDB())
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	tx := txmgrtest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, from)
	require.Len(t, tx.TxAttempts, 1)
	attempt := tx.TxAttempts[0]
	attempt.Tx = tx

	resp, cleanup := client.Get("/v2/transactions/" + attempt.Hash.String())
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	ptx := presenters.EthTxResource{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &ptx))
	txp := presenters.NewEthTxResourceFromAttempt(attempt)

	assert.Equal(t, txp.State, ptx.State)
	assert.Equal(t, txp.Data, ptx.Data)
	assert.Equal(t, txp.GasLimit, ptx.GasLimit)
	assert.Equal(t, txp.GasPrice, ptx.GasPrice)
	assert.Equal(t, txp.Hash, ptx.Hash)
	assert.Equal(t, txp.SentAt, ptx.SentAt)
	assert.Equal(t, txp.To, ptx.To)
	assert.Equal(t, txp.Value, ptx.Value)
}

func TestTransactionsController_Show_NotFound(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	txStore := txmgrtest.NewTestTxStore(t, app.GetDB())
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	tx := txmgrtest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, from)
	require.Len(t, tx.TxAttempts, 1)
	attempt := tx.TxAttempts[0]

	resp, cleanup := client.Get("/v2/transactions/" + (attempt.Hash.String() + "1"))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}
