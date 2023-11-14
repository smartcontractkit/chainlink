package web

import (
	"database/sql"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// TransactionsController displays Ethereum transactions requests.
type TransactionsController struct {
	App chainlink.Application
}

// Index returns paginated transactions
func (tc *TransactionsController) Index(c *gin.Context, size, page, offset int) {
	txs, count, err := tc.App.TxmStorageService().TransactionsWithAttempts(offset, size)
	ptxs := make([]presenters.EthTxResource, len(txs))
	for i, tx := range txs {
		tx.TxAttempts[0].Tx = tx
		ptxs[i] = presenters.NewEthTxResourceFromAttempt(tx.TxAttempts[0])
	}
	paginatedResponse(c, "transactions", size, page, ptxs, count, err)
}

// Show returns the details of a Ethereum Transaction details.
// Example:
//
//	"<application>/transactions/:TxHash"
func (tc *TransactionsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))

	ethTxAttempt, err := tc.App.TxmStorageService().FindTxAttempt(hash)
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("Transaction not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewEthTxResourceFromAttempt(*ethTxAttempt), "transaction")
}

type EvmTransactionController struct {
	Enabled     bool
	Logger      logger.SugaredLogger
	AuditLogger audit.AuditLogger
	Chains      evm.LegacyChainContainer
	KeyStore    keystore.Eth
}

func NewEVMTransactionController(app chainlink.Application) *EvmTransactionController {
	return &EvmTransactionController{
		Enabled:     app.GetConfig().TxmAsService().Enabled(),
		Logger:      app.GetLogger(),
		AuditLogger: app.GetAuditLogger(),
		Chains:      app.GetRelayers().LegacyEVMChains(),
		KeyStore:    app.GetKeyStore().Eth(),
	}
}

// Create signs and sends transaction from specified address. If address is not provided uses one of enabled keys for
// specified chain.
func (tc *EvmTransactionController) Create(c *gin.Context) {
	if !tc.Enabled {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("transactions creation disabled. To enable set TxmAsService.Enabled=true"))
		return
	}
	var tx models.CreateEVMTransactionRequest
	if err := c.ShouldBindJSON(&tx); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	if tx.IdempotencyKey == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("idempotencyKey must be set"))
		return
	}

	decoded, err := hexutil.Decode(tx.EncodedPayload)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("encodedPayload is malformed: %v", err))
		return
	}

	if tx.ChainID == nil {
		jsonAPIError(c, http.StatusBadRequest, errors.New("chainID must be set"))
		return
	}

	if tx.ToAddress == nil {
		jsonAPIError(c, http.StatusBadRequest, errors.New("toAddress must be set"))
		return
	}

	chain, err := getChain(tc.Chains, tx.ChainID.String())
	if err != nil {
		if errors.Is(err, ErrMissingChainID) {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}

		tc.Logger.Errorf("Failed to get chain", "err", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if tx.FromAddress == utils.ZeroAddress {
		tx.FromAddress, err = tc.KeyStore.GetRoundRobinAddress(tx.ChainID.ToInt())
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("failed to get fromAddress: %v", err))
			return
		}
	} else {
		err = tc.KeyStore.CheckEnabled(tx.FromAddress, tx.ChainID.ToInt())
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity,
				errors.Errorf("fromAddress %v is not available: %v", tx.FromAddress, err))
			return
		}
	}

	if tx.FeeLimit == 0 {
		tx.FeeLimit = chain.Config().EVM().GasEstimator().LimitDefault()
	}

	value := tx.Value.ToInt()
	if value == nil {
		value = big.NewInt(0)
	}
	etx, err := chain.TxManager().CreateTransaction(c, txmgr.TxRequest{
		IdempotencyKey:   &tx.IdempotencyKey,
		FromAddress:      tx.FromAddress,
		ToAddress:        *tx.ToAddress,
		EncodedPayload:   decoded,
		Value:            *value,
		FeeLimit:         tx.FeeLimit,
		ForwarderAddress: tx.ForwarderAddress,
		Strategy:         txmgrcommon.NewSendEveryStrategy(),
	})
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("transaction failed: %v", err))
		return
	}

	tc.AuditLogger.Audit(audit.EthTransactionCreated, map[string]interface{}{
		"ethTX": etx,
	})

	jsonAPIResponse(c, presenters.NewEthTxResource(etx), "eth_tx")
}
