package web_test

import (
	"math/big"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsController_Index_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()

	from := cltest.GetAccountAddress(store)
	tx := cltest.CreateTxAndAttempt(store, from, 1)
	_, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), 2)
	require.NoError(t, err)
	_, err = store.AddTxAttempt(tx, tx.EthTx(big.NewInt(3)), 3)
	require.NoError(t, err)

	resp, cleanup := client.Get("/v2/transactions?size=2")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	var attempts []models.TxAttempt
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &attempts, &links)
	require.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, attempts, 2)
	assert.Equal(t, uint64(3), attempts[0].SentAt, "expected tx attempts order by sentAt descending")
	assert.Equal(t, uint64(2), attempts[1].SentAt, "expected tx attempts order by sentAt descending")
}

func TestTransactionsController_Index_Error(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x100")
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/transactions?size=TrainingDay")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)
}

func TestTransactionsController_Show_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(store)
	tx := cltest.CreateTxAndAttempt(store, from, 1)
	tx1 := *tx
	_, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), 2)
	require.NoError(t, err)
	tx2 := *tx

	tests := []struct {
		name string
		hash string
		want models.Tx
	}{
		{"old hash", tx1.Hash.String(), tx1},
		{"current hash", tx2.Hash.String(), tx2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, cleanup := client.Get("/v2/transactions/" + test.hash)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, 200)

			ptx := presenters.Tx{}
			require.NoError(t, cltest.ParseJSONAPIResponse(resp, &ptx))

			test.want.ID = 0
			assert.Equal(t, presenters.NewTx(&test.want), ptx)
		})
	}
}

func TestTransactionsController_Show_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x100")
	})

	require.NoError(t, app.Start())
	store := app.GetStore()
	client := app.NewHTTPClient()
	from := cltest.GetAccountAddress(store)
	tx := cltest.CreateTxAndAttempt(store, from, 1)

	resp, cleanup := client.Get("/v2/transactions/" + (tx.Hash.String() + "1"))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)
}
