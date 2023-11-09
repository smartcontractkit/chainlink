package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsController_Index_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	db := app.GetSqlxDB()
	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.GetConfig().Database())
	ethKeyStore := cltest.NewKeyStore(t, db, app.Config.Database()).Eth()
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 2, from) // tx2
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 4, from)        // tx3

	// add second tx attempt for tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(&attempt))

	_, count, err := txStore.TransactionsWithAttempts(0, 100)
	require.NoError(t, err)
	require.Equal(t, count, 3)

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
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/transactions?size=TrainingDay")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, 422)
}

func TestTransactionsController_Show_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.GetConfig().Database())
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())

	tx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, from)
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
	require.NoError(t, app.Start(testutils.Context(t)))

	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.GetConfig().Database())
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	tx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, from)
	require.Len(t, tx.TxAttempts, 1)
	attempt := tx.TxAttempts[0]

	resp, cleanup := client.Get("/v2/transactions/" + (attempt.Hash.String() + "1"))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}

const txCreatePath = "/v2/transactions/evm"

// TestTransactionsController_Create_Stateless_Validations - tests Create endpoint of TestTransactionsController that
// do not require different state/configuration of the application.
func TestTransactionsController_Create_Stateless_Validations(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	t.Run("Fails on malformed json", func(t *testing.T) {
		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer([]byte("Hello")))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
	})
	t.Run("Fails on missing Idempotency key", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
		}

		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer(body))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "idempotencyKey must be set", respError.Error())
	})
	t.Run("Fails on malformed payload", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey:     "idempotency_key",
		}

		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer(body))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "encodedPayload is malformed: empty hex string", respError.Error())
	})
	t.Run("Fails if chain ID is not set", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey:     "idempotency_key",
			EncodedPayload:     "0x",
		}

		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer(body))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "chainID must be set", respError.Error())
	})
	t.Run("Fails on requesting chain that is not available", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			FromAddress:        common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey:     "idempotency_key",
			EncodedPayload:     "0x",
			ChainID:            utils.NewBigI(1),
		}

		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer(body))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, web.ErrMissingChainID.Error(), respError.Error())
	})
	t.Run("Fails when fromAddress is not specified and there are no available keys ", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			IdempotencyKey:     "idempotency_key",
			EncodedPayload:     "0x",
			ChainID:            utils.NewBigI(0),
		}

		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer(body))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "failed to get fromAddress: no sending keys available for chain 0", respError.Error())
	})
	t.Run("Fails when specified fromAddress is not available for the chain ", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			DestinationAddress: common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			FromAddress:        common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371"),
			IdempotencyKey:     "idempotency_key",
			EncodedPayload:     "0x",
			ChainID:            utils.NewBigI(0),
		}

		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		resp, cleanup := client.Post(txCreatePath, bytes.NewBuffer(body))
		t.Cleanup(cleanup)

		cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t,
			"fromAddress 0xfa01fA015c8A5332987319823728982379128371 is not available: no sending "+
				"keys available for chain 0 that match whitelist: [0xfa01fA015c8A5332987319823728982379128371]",
			respError.Error())
	})
}
