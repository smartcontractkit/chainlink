package web_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsController_Show_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
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
	ta1 := tx.TxAttempt
	ta2, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), 2)
	require.NoError(t, err)
	txWithAttempt1 := *tx
	txWithAttempt1.TxAttempt = ta1
	txWithAttempt2 := *tx
	txWithAttempt2.TxAttempt = *ta2

	tests := []struct {
		name string
		hash string
		want models.Tx
	}{
		{"old hash", ta1.Hash.String(), txWithAttempt1},
		{"current hash", ta2.Hash.String(), txWithAttempt2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, cleanup := client.Get("/v2/transactions/" + test.hash)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, 200)

			ptx := presenters.Tx{}
			require.NoError(t, cltest.ParseJSONAPIResponse(resp, &ptx))

			test.want.ID = 0
			test.want.TxID = 0
			assert.Equal(t, &test.want, ptx.Tx)
		})
	}
}

func TestTransactionsController_Show_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
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

	resp, cleanup := client.Get("/v2/txattempts/" + (tx.Hash.String() + "1"))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)
}
