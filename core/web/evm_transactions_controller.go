package web

import (
	"database/sql"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

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
	txs, count, err := tc.App.TxmORM().EthTransactionsWithAttempts(offset, size)
	ptxs := make([]presenters.EthTxResource, len(txs))
	for i, tx := range txs {
		tx.EthTxAttempts[0].EthTx = tx
		ptxs[i] = presenters.NewEthTxResourceFromAttempt(tx.EthTxAttempts[0])
	}
	paginatedResponse(c, "transactions", size, page, ptxs, count, err)
}

// Show returns the details of a Ethereum Transaction details.
// Example:
//  "<application>/transactions/:TxHash"
func (tc *TransactionsController) Show(c *gin.Context) {
	hash := common.HexToHash(c.Param("TxHash"))

	ethTxAttempt, err := tc.App.TxmORM().FindEthTxAttempt(hash)
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
