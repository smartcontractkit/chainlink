package web

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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
	txs, count, err := tc.App.TxmStorageService().TransactionsWithAttempts(c, offset, size)
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

	ethTxAttempt, err := tc.App.TxmStorageService().FindTxAttempt(c, hash)
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

// PurgeUntartedQueueRequest is a JSONAPI request to purge the unstarted transcation queue in the TXM.
type PurgeUnstartedQueueRequest struct {
	Subject string `json:"subject"`
}

// PurgeUnstartedQueueResponse is the JSONAPI response body returned after purging the unstarted transactions queue in the TXM.
type PurgeUnstartedQueueResponse struct {
	IDs []uint64 `json:"ids"`
}

// Purge clears the queue of unstarted transactions.
// Example:
//
// "<application/transactions/purge"
func (tc *TransactionsController) PurgeUnstartedQueue(c *gin.Context) {
	var req PurgeUnstartedQueueRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	subject := uuid.Nil
	if len(req.Subject) > 0 {
		subject, err = uuid.Parse(req.Subject)
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
	}

	println("XXXXXXXXXXX pruning unstarted tx queue with subject", subject.String())

	if _, err = tc.App.TxmStorageService().PruneUnstartedTxQueue(c, 0, subject); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
