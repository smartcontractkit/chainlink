package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmMocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmConfigMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	ksMocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
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

	_, count, err := txStore.TransactionsWithAttempts(txmgr.TransactionsWithAttemptsSelector{
		Offset: 0,
		Limit:  100,
	})
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

func TestTransactionsController_Index_Success_IdempotencyKey(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	db := app.GetSqlxDB()
	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.GetConfig().Database())
	ethKeyStore := cltest.NewKeyStore(t, db, app.Config.Database()).Eth()
	client := app.NewHTTPClient(nil)
	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)
	tx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 2, from)

	// add second tx attempt for tx
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(&attempt))

	_, count, err := txStore.TransactionsWithAttempts(txmgr.TransactionsWithAttemptsSelector{
		Offset: 0,
		Limit:  100,
	})
	require.NoError(t, err)
	require.Equal(t, count, 2)

	resp, cleanup := client.Get(fmt.Sprintf("/v2/transactions?size=%d&idempotencyKey=%s", 10, *tx.IdempotencyKey))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	var txs []presenters.EthTxResource
	body := cltest.ParseResponseBody(t, resp)
	require.NoError(t, web.ParsePaginatedResponse(body, &txs, &links))

	require.Len(t, txs, 1)
	require.Equal(t, attempt.Hash, txs[0].Hash, "expected tx attempts order by sentAt descending")
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

func TestTransactionsController_Create(t *testing.T) {
	t.Parallel()
	const txCreatePath = "/v2/transactions/evm"
	t.Run("Returns error if endpoint is disabled", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, txCreatePath, nil)
		router := gin.New()
		controller := &web.EvmTransactionController{
			Enabled: false,
		}
		router.POST(txCreatePath, controller.Create)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		cltest.AssertServerResponse(t, resp.Result(), http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "transactions creation disabled. To enable set TxmAsService.Enabled=true", respError.Error())
	})

	createTx := func(controller *web.EvmTransactionController, request interface{}) *httptest.ResponseRecorder {
		body, err := json.Marshal(&request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, txCreatePath, bytes.NewBuffer(body))
		router := gin.New()
		controller.Enabled = true
		router.POST(txCreatePath, controller.Create)
		router.ServeHTTP(w, req)
		return w
	}

	t.Run("Fails on malformed json", func(t *testing.T) {
		resp := createTx(&web.EvmTransactionController{}, "Hello")

		cltest.AssertServerResponse(t, resp.Result(), http.StatusBadRequest)
	})
	t.Run("Fails on missing Idempotency key", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress: ptr(common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")),
		}

		resp := createTx(&web.EvmTransactionController{}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "idempotencyKey must be set", respError.Error())
	})
	t.Run("Fails on malformed payload", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress:      ptr(common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")),
			FromAddress:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey: "idempotency_key",
		}

		resp := createTx(&web.EvmTransactionController{}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "encodedPayload is malformed: empty hex string", respError.Error())
	})
	t.Run("Fails if chain ID is not set", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress:      ptr(common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")),
			FromAddress:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x",
		}

		resp := createTx(&web.EvmTransactionController{}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "chainID must be set", respError.Error())
	})
	t.Run("Fails if toAddress is not specified", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress:      nil,
			FromAddress:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x",
			ChainID:        utils.NewBigI(0),
		}

		resp := createTx(&web.EvmTransactionController{}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "toAddress must be set", respError.Error())
	})
	chainID := utils.NewBigI(673728)
	t.Run("Fails if requested chain that is not available", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress:      ptr(common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")),
			FromAddress:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x",
			ChainID:        chainID,
		}

		chainContainer := evmMocks.NewLegacyChainContainer(t)
		chainContainer.On("Get", chainID.String()).Return(nil, web.ErrMissingChainID).Once()
		controller := &web.EvmTransactionController{
			Chains: chainContainer,
		}
		resp := createTx(controller, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, web.ErrMissingChainID.Error(), respError.Error())
	})
	t.Run("Fails when fromAddress is not specified and there are no available keys ", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress:      ptr(common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")),
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x",
			ChainID:        chainID,
		}

		chainContainer := evmMocks.NewLegacyChainContainer(t)
		chain := evmMocks.NewChain(t)
		chainContainer.On("Get", chainID.String()).Return(chain, nil).Once()

		ethKeystore := ksMocks.NewEth(t)
		ethKeystore.On("GetRoundRobinAddress", chainID.ToInt()).
			Return(nil, errors.New("failed to get key")).Once()
		resp := createTx(&web.EvmTransactionController{
			Chains:   chainContainer,
			KeyStore: ethKeystore,
		}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, "failed to get fromAddress: failed to get key", respError.Error())
	})
	t.Run("Fails when specified fromAddress is not available for the chain", func(t *testing.T) {
		fromAddr := common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")
		request := models.CreateEVMTransactionRequest{
			ToAddress:      ptr(common.HexToAddress("0xFA01FA015C8A5332987319823728982379128371")),
			FromAddress:    fromAddr,
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x",
			ChainID:        chainID,
		}

		chainContainer := evmMocks.NewLegacyChainContainer(t)
		chain := evmMocks.NewChain(t)
		chainContainer.On("Get", chainID.String()).Return(chain, nil).Once()

		ethKeystore := ksMocks.NewEth(t)
		ethKeystore.On("CheckEnabled", fromAddr, chainID.ToInt()).
			Return(errors.New("no eth key exists with address")).Once()
		resp := createTx(&web.EvmTransactionController{
			Chains:   chainContainer,
			KeyStore: ethKeystore,
		}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t,
			"fromAddress 0xfa01fA015c8A5332987319823728982379128371 is not available: no eth key exists with address",
			respError.Error())
	})

	_, _, evmConfig := txmgr.MakeTestConfigs(t)
	feeLimit := evmConfig.GasEstimator().LimitDefault()
	newChain := func(t *testing.T, txm txmgr.TxManager) evm.Chain {
		chain := evmMocks.NewChain(t)
		chain.On("TxManager").Return(txm)

		config := evmConfigMocks.NewChainScopedConfig(t)
		config.On("EVM").Return(evmConfig).Maybe()
		chain.On("Config").Return(config).Maybe()

		return chain
	}
	t.Run("Correctly populates fields for TxRequest", func(t *testing.T) {
		payload := []byte("tx_payload")
		value := big.NewInt(rand.Int64())

		feeLimitOverride := feeLimit + 10
		request := models.CreateEVMTransactionRequest{
			ToAddress:        ptr(common.HexToAddress("0xEA746B853DcFFA7535C64882E191eE31BE8CE711")),
			FromAddress:      common.HexToAddress("0x39364605296d7c77e7C2089F0e48D527bb309d38"),
			IdempotencyKey:   "idempotency_key",
			EncodedPayload:   "0x" + fmt.Sprintf("%X", payload),
			ChainID:          chainID,
			Value:            utils.NewBig(value),
			ForwarderAddress: common.HexToAddress("0x59C2B3875797c521396e7575D706B9188894eAF2"),
			FeeLimit:         feeLimitOverride,
		}

		txm := txmMocks.NewTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee](t)
		expectedError := errors.New("stub error to shortcut execution")
		txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
			IdempotencyKey:   &request.IdempotencyKey,
			FromAddress:      request.FromAddress,
			ToAddress:        *request.ToAddress,
			EncodedPayload:   payload,
			Value:            *value,
			FeeLimit:         feeLimitOverride,
			ForwarderAddress: request.ForwarderAddress,
			Strategy:         txmgrcommon.NewSendEveryStrategy(),
		}).Return(txmgr.Tx{}, expectedError).Once()

		chainContainer := evmMocks.NewLegacyChainContainer(t)
		chain := newChain(t, txm)
		chainContainer.On("Get", chainID.String()).Return(chain, nil).Once()

		ethKeystore := ksMocks.NewEth(t)
		ethKeystore.On("CheckEnabled", request.FromAddress, chainID.ToInt()).Return(nil).Once()
		resp := createTx(&web.EvmTransactionController{
			Chains:   chainContainer,
			KeyStore: ethKeystore,
		}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusInternalServerError)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, fmt.Sprintf("transaction failed: %v", expectedError),
			respError.Error())
	})
	t.Run("Correctly populates fields for TxRequest with defaults", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			ToAddress:      ptr(common.HexToAddress("0xEA746B853DcFFA7535C64882E191eE31BE8CE711")),
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x",
			ChainID:        chainID,
		}

		expectedFromAddress := common.HexToAddress("0x59C2B3875797c521396e7575D706B9188894eAF2")
		ethKeystore := ksMocks.NewEth(t)
		ethKeystore.On("GetRoundRobinAddress", chainID.ToInt()).Return(expectedFromAddress, nil).Once()

		txm := txmMocks.NewTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee](t)
		expectedError := errors.New("stub error to shortcut execution")
		txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
			IdempotencyKey: &request.IdempotencyKey,
			FromAddress:    expectedFromAddress,
			ToAddress:      *request.ToAddress,
			EncodedPayload: []byte{},
			Value:          big.Int{},
			FeeLimit:       feeLimit,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		}).Return(txmgr.Tx{}, expectedError).Once()

		chainContainer := evmMocks.NewLegacyChainContainer(t)
		chain := newChain(t, txm)
		chainContainer.On("Get", chainID.String()).Return(chain, nil).Once()

		resp := createTx(&web.EvmTransactionController{
			Chains:   chainContainer,
			KeyStore: ethKeystore,
		}, request).Result()

		// we do not really care about results here as main check is during CreateTransaction call
		cltest.AssertServerResponse(t, resp, http.StatusInternalServerError)
		respError := cltest.ParseJSONAPIErrors(t, resp.Body)
		require.Equal(t, fmt.Sprintf("transaction failed: %v", expectedError),
			respError.Error())
	})

	payload := []byte("tx_payload")
	const txID = int64(54323)
	newTxFromRequest := func(request models.CreateEVMTransactionRequest) txmgr.Tx {
		return txmgr.Tx{
			ID:             txID,
			EncodedPayload: payload,
			FromAddress:    request.FromAddress,
			FeeLimit:       feeLimit,
			State:          txmgrcommon.TxInProgress,
			ToAddress:      *request.ToAddress,
			Value:          *request.Value.ToInt(),
			ChainID:        chainID.ToInt(),
		}
	}
	newTxManager := func(request models.CreateEVMTransactionRequest) txmgr.TxManager {
		txm := txmMocks.NewTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee](t)
		tx := newTxFromRequest(request)
		txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
			IdempotencyKey: &request.IdempotencyKey,
			FromAddress:    request.FromAddress,
			ToAddress:      *request.ToAddress,
			EncodedPayload: payload,
			Value:          *request.Value.ToInt(),
			FeeLimit:       feeLimit,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		}).Return(tx, nil).Once()
		return txm
	}
	t.Run("Happy path", func(t *testing.T) {
		request := models.CreateEVMTransactionRequest{
			FromAddress:    common.HexToAddress("0x59C2B3875797c521396e7575D706B9188894eAF2"),
			ToAddress:      ptr(common.HexToAddress("0xEA746B853DcFFA7535C64882E191eE31BE8CE711")),
			IdempotencyKey: "idempotency_key",
			EncodedPayload: "0x" + fmt.Sprintf("%X", payload),
			ChainID:        chainID,
			Value:          utils.NewBigI(6838712),
		}

		ethKeystore := ksMocks.NewEth(t)
		ethKeystore.On("CheckEnabled", request.FromAddress, chainID.ToInt()).Return(nil).Once()

		txm := newTxManager(request)

		chainContainer := evmMocks.NewLegacyChainContainer(t)
		chain := newChain(t, txm)
		chainContainer.On("Get", chainID.String()).Return(chain, nil).Once()
		block := int64(56345431)
		txWithAttempts := newTxFromRequest(request)
		txWithAttempts.TxAttempts = []txmgr.TxAttempt{
			{
				Hash:                    common.HexToHash("0xa1ce83ee556cbcfc6541d5909b0d7f28f6a77399d3bd4340246f684a0f25a7f5"),
				BroadcastBeforeBlockNum: &block,
			},
		}

		resp := createTx(&web.EvmTransactionController{
			AuditLogger: audit.NoopLogger,
			Chains:      chainContainer,
			KeyStore:    ethKeystore,
		}, request).Result()

		cltest.AssertServerResponse(t, resp, http.StatusAccepted)
	})
}
