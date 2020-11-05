package web_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsController_Index_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(t, store)

	cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 3, 2, from) // tx2
	cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 4, 4, from)        // tx3

	// add second tx attempt for tx2
	blockNum := int64(3)
	attempt := cltest.NewEthTxAttempt(t, tx2.ID)
	attempt.State = models.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, store.DB.Create(&attempt).Error)

	_, count, err := store.EthTransactionsWithOrderedAttempts(0, 100)
	require.NoError(t, err)
	require.Equal(t, count, 3)

	size := 2
	resp, cleanup := client.Get(fmt.Sprintf("/v2/transactions?size=%d", size))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	var txs []presenters.EthTx
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

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/transactions?size=TrainingDay")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)
}

func TestTransactionsController_Show_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()

	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_chainId", app.Store.Config.ChainID())
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(t, store)

	tx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1, from)
	require.Len(t, tx.EthTxAttempts, 1)
	attempt := tx.EthTxAttempts[0]
	attempt.EthTx = tx

	resp, cleanup := client.Get("/v2/transactions/" + attempt.Hash.Hex())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	ptx := presenters.EthTx{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &ptx))
	txp := presenters.NewEthTxWithAttempt(&tx, &attempt)

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

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(t, store)
	tx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1, from)
	require.Len(t, tx.EthTxAttempts, 1)
	attempt := tx.EthTxAttempts[0]

	resp, cleanup := client.Get("/v2/transactions/" + (attempt.Hash.String() + "1"))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}
