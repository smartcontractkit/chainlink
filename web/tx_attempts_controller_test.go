package web_test

import (
	"math/big"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxAttemptsController_Index_Success(t *testing.T) {
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
	_, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), 2)
	require.NoError(t, err)
	_, err = store.AddTxAttempt(tx, tx.EthTx(big.NewInt(3)), 3)
	require.NoError(t, err)

	resp, cleanup := client.Get("/v2/txattempts?size=2")
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

func TestTxAttemptsController_Index_Error(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x100")
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Get("/v2/txattempts?size=TrainingDay")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)
}

func TestServiceAgreementsController_Show_Success(t *testing.T) {
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

	resp, cleanup := client.Get("/v2/txattempts/" + tx.Hash.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	var attempt models.TxAttempt
	err := web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &attempt, &links)
	require.NoError(t, err)

	assert.Equal(t, tx.TxAttempt, attempt)
}

func TestServiceAgreementsController_Show_NotFound(t *testing.T) {
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
