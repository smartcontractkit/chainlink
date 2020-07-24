package web_test

import (
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

	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()

	from := cltest.GetAccountAddress(t, store)
	tx1 := cltest.CreateTxWithNonceAndGasPrice(t, store, from, 1, 0, 1)
	transaction := cltest.NewTransaction(0)
	require.NoError(t, utils.JustError(store.AddTxAttempt(tx1, transaction)))
	cltest.CreateTxWithNonceAndGasPrice(t, store, from, 3, 1, 1)
	cltest.CreateTxWithNonceAndGasPrice(t, store, from, 4, 2, 1)
	_, count, err := store.Transactions(0, 100)
	require.NoError(t, err)
	require.Equal(t, count, 3)

	resp, cleanup := client.Get("/v2/transactions?size=2")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	var txs []presenters.Tx
	body := cltest.ParseResponseBody(t, resp)
	require.NoError(t, web.ParsePaginatedResponse(body, &txs, &links))
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	require.Len(t, txs, 2)
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
	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})
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
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(t, store)

	tx := cltest.CreateTx(t, store, from, 1)
	tx1 := *tx

	transaction := cltest.NewTransaction(2)
	require.NoError(t, utils.JustError(store.AddTxAttempt(tx, transaction)))
	tx2 := *tx

	tests := []struct {
		name string
		hash string
		want models.Tx
	}{
		{"old hash", tx1.Hash.String(), tx1},
		{"current hash", tx2.Hash.String(), tx2},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			resp, cleanup := client.Get("/v2/transactions/" + test.hash)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, http.StatusOK)

			ptx := presenters.Tx{}
			require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &ptx))

			txp := presenters.NewTx(&test.want)
			assert.Equal(t, txp.Confirmed, ptx.Confirmed)
			assert.Equal(t, txp.Data, ptx.Data)
			assert.Equal(t, txp.GasLimit, ptx.GasLimit)
			assert.Equal(t, txp.GasPrice, ptx.GasPrice)
			assert.Equal(t, txp.Hash, ptx.Hash)
			assert.Equal(t, txp.SentAt, ptx.SentAt)
			assert.Equal(t, txp.To, ptx.To)
			assert.Equal(t, txp.Value, ptx.Value)
		})
	}
}

func TestTransactionsController_Show_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)

	resp, cleanup := client.Get("/v2/transactions/" + (tx.Hash.String() + "1"))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}
